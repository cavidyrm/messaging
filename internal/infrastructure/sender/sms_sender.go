package sender

import (
	"context"
	"errors"
	"messaging/config"
)

type SMSSender struct {
	cfg config.SMSConfig
}

func NewSMSSender(cfg config.SMSConfig) *SMSSender {
	return &SMSSender{cfg: cfg}
}

func (s *SMSSender) Send(ctx context.Context, phoneNumber, text string) error {
	// Mock implementation - replace with actual provider like Kavehnegar

	//check validation for credentials then call api and send
	if s.cfg.APIKey == "" {
		return errors.New("api key is empty")
	}

	return nil
}
