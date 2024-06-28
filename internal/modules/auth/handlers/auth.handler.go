package auth_handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/troptropcontent/visa_appointment_watcher/internal/models"
	auth_service "github.com/troptropcontent/visa_appointment_watcher/internal/modules/auth/services"
)

func Signup(c echo.Context) error {
	user := models.User{}
	if err := c.Bind(&user); err != nil {
		// TODO: return signup page with error
		return c.JSON(http.StatusBadRequest, err)
	}

	err := auth_service.Signup(&user)
	if err != nil {
		// TODO: return signup page with error
		return c.JSON(http.StatusBadRequest, err)
	}

	// TODO: return signup page with success
	return c.JSON(http.StatusOK, user)
}
