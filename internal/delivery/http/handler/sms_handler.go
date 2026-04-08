package handler

import (
	"encoding/json"
	"messaging/internal/domain/event"
	"net/http"
	"time"

	"messaging/internal/usecase/sms"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SMSHandler struct {
	service *sms.Service
}

func NewSMSHandler(service *sms.Service) *SMSHandler {
	return &SMSHandler{service: service}
}

type SendSMSRequest struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
	Text        string `json:"text" validate:"required"`
}

func (h *SMSHandler) SendSMS(c echo.Context) error {
	var req SendSMSRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	data, _ := json.Marshal(map[string]string{
		"phone_number": req.PhoneNumber,
		"text":         req.Text,
	})

	event := event.Event{
		EventID:       uuid.New(),
		AggregateID:   uuid.New(),
		AggregateType: "sms",
		EventType:     "SmsSendRequested",
		Version:       1,
		Data:          data,
		Metadata:      json.RawMessage(`{}`),
		Timestamp:     time.Now(),
	}

	if err := h.service.ProcessAndSendSMS(c.Request().Context(), event); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send SMS"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "SMS request accepted"})
}

func (h *SMSHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}

	msg, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, msg)
}
