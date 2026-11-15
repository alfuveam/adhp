// Package mail envia e-mails via SMTP (porta 465 com TLS implícito, demais com STARTTLS).
package mail

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"html/template"
	"mime"
	"mime/quotedprintable"
	"net"
	netmail "net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/alfuveam/adhp/backend/config"
)

var ErrNotConfigured = errors.New("SMTP não configurado")

// Configured informa se há SMTP configurado para envio.
func Configured() bool {
	return config.SmtpHost != "" && config.SmtpUser != "" && config.SmtpPass != ""
}

// IsMicrosoftAddress indica e-mails de domínios da Microsoft (Hotmail/Outlook/Live/MSN),
// para os quais a entrega está com problemas.
func IsMicrosoftAddress(email string) bool {
	i := strings.LastIndex(email, "@")
	if i < 0 {
		return false
	}
	domain := strings.ToLower(strings.TrimSpace(email[i+1:]))
	for _, prefix := range []string{"hotmail.", "outlook.", "live.", "msn.", "windowslive."} {
		if strings.HasPrefix(domain, prefix) {
			return true
		}
	}
	return false
}

func Send(to, subject, htmlBody string) error {
	if !Configured() {
		return ErrNotConfigured
	}
	host, port := config.SmtpHost, config.SmtpPort
	if port == "" {
		port = "587"
	}
	from := config.SmtpFrom
	if from == "" {
		from = config.SmtpUser
	}

	// o envelope precisa do endereço puro; o cabeçalho pode ter nome de exibição
	envelopeFrom := from
	if addr, err := netmail.ParseAddress(from); err == nil {
		envelopeFrom = addr.Address
	}
	domain := "localhost"
	if i := strings.LastIndex(envelopeFrom, "@"); i >= 0 {
		domain = envelopeFrom[i+1:]
	}
	idBytes := make([]byte, 16)
	_, _ = rand.Read(idBytes)

	// Date e Message-ID são exigidos por receptores rígidos (Outlook/Hotmail, Gmail);
	// assunto com acentos precisa ser codificado (RFC 2047).
	headers := [][2]string{
		{"From", from},
		{"To", to},
		{"Subject", mime.QEncoding.Encode("UTF-8", subject)},
		{"Date", time.Now().Format(time.RFC1123Z)},
		{"Message-ID", fmt.Sprintf("<%s@%s>", hex.EncodeToString(idBytes), domain)},
		{"MIME-Version", "1.0"},
		{"Content-Type", `text/html; charset="UTF-8"`},
		{"Content-Transfer-Encoding", "quoted-printable"},
	}
	var msg strings.Builder
	for _, h := range headers {
		fmt.Fprintf(&msg, "%s: %s\r\n", h[0], h[1])
	}
	msg.WriteString("\r\n")
	qp := quotedprintable.NewWriter(&msg)
	if _, err := qp.Write([]byte(htmlBody)); err != nil {
		return fmt.Errorf("codificar corpo: %w", err)
	}
	if err := qp.Close(); err != nil {
		return fmt.Errorf("codificar corpo: %w", err)
	}

	auth := smtp.PlainAuth("", config.SmtpUser, config.SmtpPass, host)
	return sendSMTP(net.JoinHostPort(host, port), host, port, auth, envelopeFrom, to, []byte(msg.String()))
}

// sendSMTP envia com timeouts. net/smtp.SendMail não tem timeout nem TLS implícito,
// então na porta 465 ficaria travado esperando a saudação em texto puro.
func sendSMTP(addr, host, port string, auth smtp.Auth, from, to string, msg []byte) error {
	dialer := &net.Dialer{Timeout: 15 * time.Second}
	tlsConfig := &tls.Config{ServerName: host}

	var conn net.Conn
	var err error
	if port == "465" {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	} else {
		conn, err = dialer.Dial("tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("smtp dial %s: %w", addr, err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(30 * time.Second))

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer c.Close()

	if port != "465" {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(tlsConfig); err != nil {
				return fmt.Errorf("smtp starttls: %w", err)
			}
		}
	}
	if ok, _ := c.Extension("AUTH"); ok {
		if err := c.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}
	if err := c.Mail(from); err != nil {
		return fmt.Errorf("smtp mail from: %w", err)
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("smtp rcpt to: %w", err)
	}
	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp data close: %w", err)
	}
	return c.Quit()
}

// template embutido: a imagem final é FROM scratch, sem arquivos externos
var recoveryTmpl = template.Must(template.New("recovery").Parse(`<!DOCTYPE html>
<html lang="pt-BR"><body style="margin:0;padding:24px;background:#f5f7f4;font-family:Arial,Helvetica,sans-serif;color:#1e293b">
<table role="presentation" width="100%" cellpadding="0" cellspacing="0"><tr><td align="center">
<table role="presentation" width="480" cellpadding="0" cellspacing="0" style="max-width:480px;background:#ffffff;border:1px solid #dce1d9;border-radius:8px">
<tr><td style="padding:28px">
<h1 style="margin:0 0 12px;font-size:20px;color:#1e293b">Recuperação de senha - ADHP</h1>
<p style="margin:0 0 16px;font-size:14px;line-height:1.5">Olá, {{.Name}}. Recebemos um pedido para redefinir a senha da sua conta. Use o código abaixo:</p>
<p style="margin:0 0 16px;font-size:30px;font-weight:bold;letter-spacing:8px;text-align:center;color:#285b75">{{.Code}}</p>
<p style="margin:0 0 8px;font-size:13px;line-height:1.5;color:#475569">O código vale por {{.Minutes}} minutos.</p>
<p style="margin:0;font-size:13px;line-height:1.5;color:#475569">Se você não pediu a recuperação, ignore este e-mail: sua senha continua a mesma.</p>
</td></tr></table></td></tr></table>
</body></html>`))

func SendRecoveryCode(to, name, code string, minutes int) error {
	var buf bytes.Buffer
	data := struct {
		Name, Code string
		Minutes    int
	}{name, code, minutes}
	if err := recoveryTmpl.Execute(&buf, data); err != nil {
		return err
	}
	return Send(to, "Recuperação de senha - ADHP", buf.String())
}
