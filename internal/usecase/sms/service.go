package sms

import (
	"context"
	"fmt"
	"messaging/internal/domain/sms"
)

type Service struct {
	repo      sms.Repository
	eventRepo sms.EventRepository
	//logger    *logger.Logger
	sender sms.Sender
}

func NewService(repo sms.Repository, eventRepo sms.EventRepository, sender sms.Sender) *Service {
	return &Service{
		repo:      repo,
		eventRepo: eventRepo,
		sender:    sender,
	}
}

func (s *Service) Send(ctx context.Context, m *sms.Message) error {

	//We can consider to add transactional all of these steps to have more consistency for add and send a message or rollback
	err := s.eventRepo.SaveEvent(ctx, m)
	if err != nil {
		err = fmt.Errorf("failed to save sms event: %w", err)
		return err
	}

	err = s.repo.Save(ctx, m)
	if err != nil {
		err = fmt.Errorf("failed to save sms message: %w", err)
		return err
	}

	err = s.sender.Send(ctx, m.PhoneNumber, m.Text)
	if err != nil {
		err = fmt.Errorf("failed to send sms message: %w", err)
		return err
	}

	return nil
}
