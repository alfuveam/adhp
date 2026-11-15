package controller

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/internal/mail"
	"github.com/alfuveam/adhp/backend/internal/turnstile"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

const (
	recoveryCodeTTL      = 15 * time.Minute
	recoveryResendWait   = time.Minute
	recoveryMaxAttempts  = 5
	msgTurnstileRequired = "Confirme a verificação de segurança e tente novamente."
)

type RecoveryRequest struct {
	Email          string `json:"email"`
	TurnstileToken string `json:"turnstile_token"`
}

type RecoveryResetRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

type RecoveryResponse struct {
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

func writeRecovery(w http.ResponseWriter, status int, res RecoveryResponse) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(res)
}

func hashRecoveryCode(code string) string {
	sum := sha256.Sum256([]byte(config.SaltDB + code))
	return hex.EncodeToString(sum[:])
}

func newRecoveryCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// pgx grava timestamp sem fuso pelo relógio de parede; usar sempre UTC
func nowUTC() time.Time { return time.Now().UTC() }

// RequestPasswordRecovery envia por e-mail um código de 6 dígitos para redefinir a senha.
func RequestPasswordRecovery(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req RecoveryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeRecovery(w, http.StatusBadRequest, RecoveryResponse{Error: "Requisição com payload inválido"})
		return
	}
	if !turnstile.Verify(req.TurnstileToken) {
		writeRecovery(w, http.StatusBadRequest, RecoveryResponse{Error: msgTurnstileRequired})
		return
	}
	email := strings.TrimSpace(req.Email)
	if email == "" {
		writeRecovery(w, http.StatusBadRequest, RecoveryResponse{Error: "Informe o e-mail da conta."})
		return
	}
	if !mail.Configured() {
		log.Println("RequestPasswordRecovery: SMTP não configurado")
		writeRecovery(w, http.StatusServiceUnavailable, RecoveryResponse{Error: "O envio de e-mail não está disponível no momento. Procure o professor responsável."})
		return
	}

	// mesma resposta para e-mail inexistente, para não revelar quem tem conta
	sent := RecoveryResponse{Message: "Se o e-mail estiver cadastrado, enviamos um código de recuperação."}

	ctx := context.Background()
	user, err := q.OnLoginCheckEmailAndPasswordUser(ctx, email)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			log.Println("RequestPasswordRecovery: busca do usuário:", err)
		}
		writeRecovery(w, http.StatusOK, sent)
		return
	}

	now := nowUTC()
	if prev, err := q.GetPasswordRecovery(ctx, user.ID); err == nil && now.Sub(prev.Createat.Time) < recoveryResendWait {
		writeRecovery(w, http.StatusTooManyRequests, RecoveryResponse{Error: "Aguarde um minuto antes de pedir outro código."})
		return
	}

	code, err := newRecoveryCode()
	if err != nil {
		log.Println("RequestPasswordRecovery: gerar código:", err)
		writeRecovery(w, http.StatusInternalServerError, RecoveryResponse{Error: "Erro ao gerar o código. Tente novamente."})
		return
	}
	err = q.UpsertPasswordRecovery(ctx, generated.UpsertPasswordRecoveryParams{
		UserID:    user.ID,
		CodeHash:  hashRecoveryCode(code),
		Createat:  pgtype.Timestamp{Time: now, Valid: true},
		ExpiresAt: pgtype.Timestamp{Time: now.Add(recoveryCodeTTL), Valid: true},
	})
	if err != nil {
		log.Println("RequestPasswordRecovery: salvar código:", err)
		writeRecovery(w, http.StatusInternalServerError, RecoveryResponse{Error: "Erro ao gerar o código. Tente novamente."})
		return
	}

	if err := mail.SendRecoveryCode(user.Email, user.Completename, code, int(recoveryCodeTTL.Minutes())); err != nil {
		log.Printf("RequestPasswordRecovery: envio para %s falhou: %v", user.Email, err)
		_ = q.DeletePasswordRecovery(ctx, user.ID)
		writeRecovery(w, http.StatusBadGateway, RecoveryResponse{Error: "Não foi possível enviar o e-mail. Tente novamente mais tarde ou procure o professor responsável."})
		return
	}

	writeRecovery(w, http.StatusOK, sent)
}

// ResetPassword troca a senha usando o código recebido por e-mail.
func ResetPassword(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req RecoveryResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeRecovery(w, http.StatusBadRequest, RecoveryResponse{Error: "Requisição com payload inválido"})
		return
	}
	if len(req.NewPassword) > 100 || len(req.NewPassword) < 9 {
		writeRecovery(w, http.StatusBadRequest, RecoveryResponse{Error: "A nova senha deve ter entre 9 e 100 caracteres."})
		return
	}

	invalid := RecoveryResponse{Error: "Código inválido ou expirado."}
	ctx := context.Background()
	user, err := q.OnLoginCheckEmailAndPasswordUser(ctx, strings.TrimSpace(req.Email))
	if err != nil {
		writeRecovery(w, http.StatusBadRequest, invalid)
		return
	}
	rec, err := q.GetPasswordRecovery(ctx, user.ID)
	if err != nil || nowUTC().After(rec.ExpiresAt.Time) || rec.Attempts >= recoveryMaxAttempts {
		writeRecovery(w, http.StatusBadRequest, invalid)
		return
	}
	given := hashRecoveryCode(strings.TrimSpace(req.Code))
	if subtle.ConstantTimeCompare([]byte(given), []byte(rec.CodeHash)) != 1 {
		_ = q.IncrementPasswordRecoveryAttempts(ctx, user.ID)
		writeRecovery(w, http.StatusBadRequest, invalid)
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(config.SaltDB+req.NewPassword), 14)
	if err != nil {
		log.Println("ResetPassword: bcrypt:", err)
		writeRecovery(w, http.StatusInternalServerError, RecoveryResponse{Error: "Erro ao redefinir a senha."})
		return
	}
	if err := q.UpdateUserPassword(ctx, generated.UpdateUserPasswordParams{ID: user.ID, Password: string(hashed)}); err != nil {
		log.Println("ResetPassword: atualizar senha:", err)
		writeRecovery(w, http.StatusInternalServerError, RecoveryResponse{Error: "Erro ao redefinir a senha."})
		return
	}
	_ = q.DeletePasswordRecovery(ctx, user.ID)

	writeRecovery(w, http.StatusOK, RecoveryResponse{Message: "Senha alterada. Você já pode entrar com a nova senha."})
}
