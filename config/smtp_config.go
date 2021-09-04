package config

import (
	"fmt"
	"net/smtp"
)

// SMTPClient is the wrapper for smtp
type SMTPClient struct {
	host string
	auth smtp.Auth
}

// NewSMTPClient returns a new smtp client wrapper
func NewSMTPClient(smtpHost string, smtpPort int, userKey string, keyPass string) *SMTPClient {
	auth := smtp.PlainAuth(
		"",
		userKey,
		keyPass,
		smtpHost,
	)
	return &SMTPClient{
		host: fmt.Sprintf("%s:%d", smtpHost, smtpPort),
		auth: auth,
	}
}

// SendMail wraps smtp.SendMail
func (s SMTPClient) SendMail(from string, to []string, msg []byte) error {
	return smtp.SendMail(s.host, s.auth, from, to, msg)
}
