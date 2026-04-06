package email

import (
	"messaging/internal/domain/email"
	"messaging/internal/domain/event"
)

type Service struct {
	eventRepo event.Repository
	emailRepo email.Repository
	sender    email.Sender
}

func NewEmailService(eventRepo event.Repository, emailRepo email.Repository, sender email.Sender) *Service {
	return &Service{
		eventRepo: eventRepo,
		emailRepo: emailRepo,
		sender:    sender,
	}
}
