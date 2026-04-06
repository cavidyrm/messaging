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
				return fmt.Errorf("unhandled SMS event type: %s", event.EventType)
			}
			return nil
		}
	case "Email":
		//if ev.EventType == "EmailSendRequested" { return email service }
		return fmt.Errorf("email routing not fully implemented yet")

	default:
		return fmt.Errorf("unknown aggregate type: %s", event.AggregateType)
	}
	return nil
}
