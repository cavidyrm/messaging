package usecase

import (
	"context"
	"fmt"
	"messaging/internal/domain/event"
	emailUsecase "messaging/internal/usecase/email"
	smsUsecase "messaging/internal/usecase/sms"
)

type MessageRouter struct {
	smsService   *smsUsecase.Service
	emailService *emailUsecase.Service
}

func NewMessageRouter(smsService *smsUsecase.Service, emailService *emailUsecase.Service) *MessageRouter {
	return &MessageRouter{
		smsService:   smsService,
		emailService: emailService,
	}
}

func (r *MessageRouter) Route(ctx context.Context, event event.Event) error {
	// Route by Aggregate (Domain)
	switch event.AggregateType {
	case "SMS":
		// Route by specific Event
		if event.EventType == "SMSSendRequested" {
			// Pass the complete event down to the service layer
			err := r.smsService.ProcessAndSendSMS(ctx, event)
			if err != nil {
				return fmt.Errorf("unhandled SMS event type in router: %s , %v", event.EventType, err)
			}
			return nil
		}
	case "Email":
		if event.EventType == "EmailSendRequested" {
			if err := r.emailService.ProcessAndSendEmail(ctx, event); err != nil {
				return fmt.Errorf("unhandled email event type in router: %s , %v", event.EventType, err)
			}
		}
		return nil

	default:
		return fmt.Errorf("unknown aggregate type: %s", event.AggregateType)
	}
	return nil
}
