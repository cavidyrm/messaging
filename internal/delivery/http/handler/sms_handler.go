package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"messaging/internal/usecase/sms"
)

type SMSHandler struct {
	service *sms.SMSService
}

func NewSMSHandler(service *sms.SMSService) *SMSHandler {
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

	//matching sms sending fields based on interfaces
	if err := h.service.Send(c.Request().Context(), req.PhoneNumber, req.Text); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "sent"})
}
