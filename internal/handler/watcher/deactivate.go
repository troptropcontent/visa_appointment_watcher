package watcher_handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	watcher_show_switcher_template "github.com/troptropcontent/visa_appointment_watcher/internal/views/watcher/show/_switcher"
)

func Deactivate(c echo.Context) error {
	config.MustSet("watcher_running", "false")
	data := &watcher_show_switcher_template.Template{
		IsActivated: false,
	}
	return c.Render(http.StatusOK, "watcher/show/_switcher", data)
}
