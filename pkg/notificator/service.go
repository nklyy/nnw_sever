package notificator

import (
	"bytes"
	"context"
	"fmt"
	"nnw_s/pkg/errors"
	"nnw_s/pkg/smtp"
	"os"
	"path"
	"path/filepath"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	mimeHeaders      = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	sendEmailTimeout = 5 * time.Second
)

type Service interface {
	SendEmail(ctx context.Context, email *Email) error
}

type service struct {
	log        *logrus.Logger
	smtpClient *smtp.Client
}

func NewService(log *logrus.Logger, smtpClient *smtp.Client) (Service, error) {
	if log == nil {
		return nil, errors.NewInternal("invalid logger")
	}
	if smtpClient == nil {
		return nil, errors.NewInternal("invalid SMTP client")
	}
	return &service{log: log, smtpClient: smtpClient}, nil
}

func (svc *service) SendEmail(ctx context.Context, email *Email) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	root := filepath.Dir(dir)

	t, err := template.ParseFiles(path.Join(root, "templates/"+email.Template))
	if err != nil {
		return errors.NewInternal(err.Error())
	}

	var body bytes.Buffer

	_, err = body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", email.Subject, mimeHeaders)))
	if err != nil {
		svc.log.WithContext(ctx).Errorf("failed to write subject to email body: %v", err)
		return errors.NewInternal(err.Error())
	}

	if err := t.Execute(&body, &email.Data); err != nil {
		svc.log.WithContext(ctx).Errorf("failed to write template data to email body: %v", err)
		return errors.NewInternal(err.Error())
	}

	newCtx, cancel := context.WithTimeout(ctx, sendEmailTimeout)
	defer cancel()

	done := make(chan struct{}, 1)
	sendEmailErr := make(chan error, 1)

	go func() {
		select {
		case <-newCtx.Done():
			sendEmailErr <- errors.NewInternal("failed to send email due to timeout")
			return
		default:
			if err := svc.smtpClient.SendMail(email.Sender, []string{email.Recipient}, body.Bytes()); err != nil {
				svc.log.WithContext(ctx).Errorf("failed to send email: %v", err)
				sendEmailErr <- err
				return
			}
		}

		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case err = <-sendEmailErr:
		return err
	}
}
