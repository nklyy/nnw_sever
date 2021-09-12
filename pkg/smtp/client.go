package smtp

import (
	"fmt"
	"net/smtp"
)

// Client is the wrapper for smtp
type Client struct {
	host string
	auth smtp.Auth
}

// NewClient returns a new smtp client wrapper
func NewClient(smtpHost string, smtpPort int, userKey string, keyPass string) *Client {
	auth := smtp.PlainAuth(
		"",
		userKey,
		keyPass,
		smtpHost,
	)
	return &Client{
		host: fmt.Sprintf("%s:%d", smtpHost, smtpPort),
		auth: auth,
	}
}

// SendMail wraps smtp.SendMail
func (s *Client) SendMail(from string, to []string, msg []byte) error {
	return smtp.SendMail(s.host, s.auth, from, to, msg)
}
