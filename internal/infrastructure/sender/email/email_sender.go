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

func (s *EmailSender) Send(ctx context.Context, email *email.Email) error {
	// Mock implementation - replace with actual provider

	// SMTP Server Details
	host := "mail.netpardazco.com"
	port := "465"
	user := "notify@netpardazco.com"
	password := "notify123*@!"

	// Email Details
	from := "username@something.com"
	to := "recipient@example.com"
	subject := "Test Email from Go Stdlib"
	body := "This email was sent using only the Go standard library over Port 465!"

	// 1. Craft the raw SMTP message (Headers + Body)
	// Notice the \r\n (Carriage Return + Line Feed), which is required by the SMTP spec
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)

	// 2. Setup Authentication
	auth := smtp.PlainAuth("", user, password, host)

	// 3. Configure TLS (Port 465 requires Implicit TLS)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         host,
	}

	// 4. Dial the server directly via TLS
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		log.Fatalf("Failed to dial TLS: %v", err)
	}

	// 5. Create new SMTP client using the TLS connection
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		log.Fatalf("Failed to create SMTP client: %v", err)
	}
	defer client.Quit()

	// 6. Authenticate
	if err = client.Auth(auth); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// 7. Set the sender and recipient
	if err = client.Mail(from); err != nil {
		log.Fatalf("Failed to set sender: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		log.Fatalf("Failed to set recipient: %v", err)
	}

	// 8. Send the email body
	w, err := client.Data()
	if err != nil {
		log.Fatalf("Failed to create data writer: %v", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		log.Fatalf("Failed to write message: %v", err)
	}

	err = w.Close()
	if err != nil {
		log.Fatalf("Failed to close data writer: %v", err)
	}

	log.Println("Email sent successfully!")

	return nil
}
