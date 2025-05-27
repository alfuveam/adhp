package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"time"

	"github.com/alfuveam/adhp/backend/config"
	"github.com/alfuveam/adhp/backend/generated"
	"github.com/alfuveam/adhp/backend/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type ErrorResponseCreateAccount struct {
	Completename             string `json:"error_completename"`
	Email                    string `json:"error_email"`
	Password                 string `json:"error_password"`
	RepeticaoEspacadaMinutos string `json:"error_repeticao_espacada_minutos"`
	OnCreateAccount          string `json:"error_on_create_account"`
}

type MyCustomClaims struct {
	User *models.UserJwt
	jwt.RegisteredClaims
}

type ErroOnLogin struct {
	Error string `json:"error_on_login"`
}

type DashboardResponse struct {
	Name string `json:"name"`
}

func ValidateJWT(tokenString string) (models.UserJwt, bool, error) {
	var user models.UserJwt
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.MySigningKey), nil
	}, jwt.WithLeeway(5*time.Second))
	if err != nil {
		// log.Fatal(err)
		return user, false, err
	} else if claims, ok := token.Claims.(*MyCustomClaims); ok {
		return *claims.User, true, nil
	} else {
		// log.Fatal("unknown claims type, cannot proceed")
		return user, false, errors.New("unknown claims type, cannot proceed")
	}
}

func CreateAndSignJWT(userJWT *models.UserJwt) (string, error) {
	// Create claims while leaving out some of the optional fields
	claims := MyCustomClaims{
		userJWT,
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString([]byte(config.MySigningKey))
	_, _, err := ValidateJWT(ss)
	if err != nil {
		log.Println("CreateAndSignJWT: ", err)
	}

	return ss, nil
}

func OnLogin(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req generated.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Requisição com payload inválido",
		})
		return
	}

	user, err := q.OnLoginCheckEmailAndPasswordUser(context.Background(), req.Email)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid user",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(config.SaltDB+req.Password))

	if err != nil {
		log.Println("bcrypt.CompareHashAndPassword: ", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid password",
		})
		return
	}

	userJWT := models.UserJwt{
		Id:       user.ID.String(),
		UserType: int(user.Usertype),
	}

	token, err := CreateAndSignJWT(&userJWT)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid token",
		})
		return
	}

	res := LoginResponse{
		Message: "Login successful",
		Token:   token,
	}
	json.NewEncoder(w).Encode(res)
}

func OnLogout(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	response := map[string]string{
		"message": "Logout successful",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func RegisterUser(w http.ResponseWriter, r *http.Request, q *generated.Queries) {
	var req generated.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Requisição com payload inválido",
		})
		return
	}

	var errorRes ErrorResponseCreateAccount
	hasError := false
	if len(req.Completename) > 50 || len(req.Completename) < 5 {
		errorRes.Completename = "O tamanho correto do nome é entre 50 e 5."
		hasError = true
	}

	if len(req.Email) > 100 {
		errorRes.Email = "Tamanho incorreto do Email."
		hasError = true
	}

	if req.RepeticaoEspacadaMinutos == 1 { // 1 Hora
		req.RepeticaoEspacadaMinutos = 60
	} else if req.RepeticaoEspacadaMinutos == 2 { // 9 Horas
		req.RepeticaoEspacadaMinutos = 60 * 9
	} else if req.RepeticaoEspacadaMinutos == 3 { // 1 Dia
		req.RepeticaoEspacadaMinutos = 60 * 24
	} else if req.RepeticaoEspacadaMinutos == 4 { // 6 Dias
		req.RepeticaoEspacadaMinutos = 60 * 24 * 6
	} else if req.RepeticaoEspacadaMinutos == 5 { // 31 Dias
		req.RepeticaoEspacadaMinutos = 60 * 24 * 31
	} else {
		errorRes.RepeticaoEspacadaMinutos = "Repetição incorreta."
		hasError = true
	}

	if len(req.Password) > 100 || len(req.Password) < 9 {
		errorRes.Password = "O tamanho correto para o password é entre 9 e 100."
		hasError = true
	}

	req.Password = config.SaltDB + req.Password
	bytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
	if err != nil {
		log.Println("bcrypt.GenerateFromPassword: %v", err)
		errorRes.OnCreateAccount = "Erro ao criar a conta, entre em contato com o administrador."
		hasError = true
	}

	req.Password = string(bytes)

	if _, err := mail.ParseAddress(req.Email); err != nil {
		errorRes.Email = err.Error()
		hasError = true
	}

	if hasError {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	req.Createat = pgtype.Timestamp{Time: time.Now(), Valid: true}
	req.Lastlogin = pgtype.Timestamp{Time: time.Now(), Valid: true}

	createUser := generated.CreateUserParams{
		Completename:             req.Completename,
		Email:                    req.Email,
		Password:                 req.Password,
		Createat:                 req.Createat,
		Lastlogin:                req.Lastlogin,
		Usertype:                 1,
		RepeticaoEspacadaMinutos: req.RepeticaoEspacadaMinutos,
	}

	user, err := q.CreateUser(context.Background(), createUser)
	if err != nil {
		log.Println("req.CreateUser: %v", err)
		errorRes.OnCreateAccount = "Erro ao criar a conta, entre em contato com o administrador.."
		hasError = true
	}

	if hasError {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorRes)
		return
	}

	userJWT := models.UserJwt{
		Id:       user.ID.String(),
		UserType: int(user.Usertype),
	}
	token, err := CreateAndSignJWT(&userJWT)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErroOnLogin{
			Error: "Invalid token",
		})
		return
	}

	res := LoginResponse{
		Message: "Login successful",
		Token:   token,
	}
	json.NewEncoder(w).Encode(res)
}
