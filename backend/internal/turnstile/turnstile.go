// Package turnstile valida tokens do Cloudflare Turnstile.
package turnstile

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/alfuveam/adhp/backend/config"
)

// chave secreta de teste da Cloudflare que sempre aprova
const testSecret = "1x0000000000000000000000000000000AA"

type validationResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

var client = &http.Client{Timeout: 10 * time.Second}

// Enabled informa se a verificação está ativa (TURNSTILE_SECRET definido).
func Enabled() bool {
	return config.TurnstileSecret != "" && config.TurnstileSecret != testSecret
}

// Verify valida o token gerado pelo widget. Sem TURNSTILE_SECRET (dev) sempre aprova.
func Verify(token string) bool {
	if !Enabled() {
		return true
	}
	if token == "" {
		return false
	}

	form := url.Values{"secret": {config.TurnstileSecret}, "response": {token}}
	resp, err := client.PostForm("https://challenges.cloudflare.com/turnstile/v0/siteverify", form)
	if err != nil {
		log.Printf("turnstile: erro na requisição: %v", err)
		return false
	}
	defer resp.Body.Close()

	var v validationResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		log.Printf("turnstile: erro ao decodificar resposta: %v", err)
		return false
	}
	if !v.Success {
		log.Printf("turnstile: validação falhou: %v", v.ErrorCodes)
	}
	return v.Success
}
