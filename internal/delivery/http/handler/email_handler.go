package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"messaging/internal/usecase/email"
)

type EmailHandler struct {
	service *email.Service
}

func NewEmailHandler(service *email.Service) *EmailHandler {
	return &EmailHandler{service: service}
}

type SendEmailRequest struct {
	EmailAddress string `json:"email_address" validate:"required,email"`
	Title        string `json:"title" validate:"required"`
	Text         string `json:"text" validate:"required"`
}

func (h *EmailHandler) SendEmail(c echo.Context) error {
	var req SendEmailRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := h.service.SendEmail(c.Request().Context(), req.EmailAddress, req.Title, req.Text); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "sent"})
}
