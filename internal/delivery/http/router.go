package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"messaging/internal/delivery/http/handler"
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
	api.POST("/email", emailHandler.SendEmail)

	//e.GET("/health", healthHandler.Health)

	return e
}
