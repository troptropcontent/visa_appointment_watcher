package auth_api

import (
	"github.com/labstack/echo/v4"
	auth_handler "github.com/troptropcontent/visa_appointment_watcher/internal/modules/auth/handlers"
)

const (
	AuthGroup = "/auth"
)

func RegisterRoutes(server *echo.Echo) {
	group := server.Group(AuthGroup)
	group.POST("/signup", auth_handler.Signup)
}
