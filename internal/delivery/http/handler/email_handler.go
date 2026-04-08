package handler

import (
	"encoding/json"
	"messaging/internal/domain/event"
	"net/http"
	"time"

	"messaging/internal/usecase/email"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type EmailHandler struct {
	service *email.Service
}

func NewEmailHandler(service *email.Service) *EmailHandler {
	return &EmailHandler{service: service}
}

type SendEmailRequest struct {
	Address string `json:"address" validate:"required,email"`
	Subject string `json:"subject" validate:"required"`
	Body    string `json:"body" validate:"required"`
}

func (h *EmailHandler) SendEmail(c echo.Context) error {
	var req SendEmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	data, _ := json.Marshal(map[string]string{
		"address": req.Address,
		"subject": req.Subject,
		"body":    req.Body,
	})

	ev := event.Event{
		EventID:       uuid.New(),
		AggregateID:   uuid.New(),
		AggregateType: "email",
		EventType:     "EmailSendRequested",
		Version:       1,
		Data:          data,
		Metadata:      json.RawMessage(`{}`),
		Timestamp:     time.Now(),
	}

	if err := h.service.ProcessAndSendEmail(c.Request().Context(), ev); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to send Email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "Email request accepted"})
}

func (h *EmailHandler) GetByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid id"})
	}
	msg, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, msg)
}
