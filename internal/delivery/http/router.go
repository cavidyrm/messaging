package http

import (
	"messaging/internal/delivery/http/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRouter(
	smsHandler *handler.SMSHandler,
	emailHandler *handler.EmailHandler,
) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	api := e.Group("/api/v1")

	api.POST("/sms", smsHandler.SendSMS)
	api.GET("/sms/:id", smsHandler.GetByID)

	api.POST("/email", emailHandler.SendEmail)
	api.GET("/email/:id", emailHandler.GetByID)

	//e.GET("/health", healthHandler.Health)

	return e
}
