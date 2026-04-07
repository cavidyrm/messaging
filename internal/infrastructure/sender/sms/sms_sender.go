package sms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"messaging/config"
	"net/http"
)

type SMSSender struct {
	cfg config.SMSConfig
}

func NewSMSSender(cfg config.SMSConfig) *SMSSender {
	return &SMSSender{cfg: cfg}
}

func (s *SMSSender) Send(ctx context.Context, phoneNumber string, text string) error {
	// Mock implementation - replace with actual provider like Kavehnegar

	//check validation for credentials then call api and send
	//if s.cfg.APIKey == "" {
	//	return errors.New("api key is empty")
	//}

	//build sender based on config for now I have added hardcoded stuffs

	body := map[string]any{
		"lineNumber":  30002108004976,
		"MessageText": text,
		"Mobiles":     []string{phoneNumber},
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.sms.ir/v1/send/bulk", bytes.NewBuffer(jsonBody))
	req.Header.Set("x-api-key", "jOCRVhdCeYPxunEEbaZhXzaxqBO8gRWs148V7vQ2VUv29AZQ")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println(string(respBody))
	return nil
}
