package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"messaging/config"
	"messaging/internal/domain/email"
	"net/smtp"
)

type EmailSender struct {
	cfg config.EmailConfig
}

func NewEmailSender(cfg config.EmailConfig) *EmailSender {
	return &EmailSender{cfg: cfg}
}

type loginAuth struct{ username, password string }

func (a *loginAuth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		return []byte(a.password), nil
	}
	return nil, nil
}

func (s *EmailSender) Send(ctx context.Context, e *email.Email) error {
	// Mock implementation - replace with actual provider

	//// SMTP Server Details
	//host := "mail.netpardazco.com"
	//port := "465"
	//user := "notify@netpardazco.com"
	//password := "notify123*@!"
	//
	//// Email Details
	//from := "notify@netpardazco.com"
	//to := "recipient@example.com"
	//subject := "Test Email from Go Stdlib"
	//body := "This email was sent using only the Go standard library over Port 465!"

	host := s.cfg.Host
	port := s.cfg.Port
	user := "notify@netpardazco.com"
	password := "notify123*@!"
	from := "notify@netpardazco.com"

	message := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		from, e.Address, e.Subject, e.Body,
	)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%s", host, port), tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to dial TLS: %w", err)
	}

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	auth := &loginAuth{user, password}
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP auth failed: %w", err)
	}
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(e.Address); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data writer: %w", err)
	}
	if _, err = w.Write([]byte(message)); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	log.Printf("Email sent to %s", e.Address)
	return nil

}
