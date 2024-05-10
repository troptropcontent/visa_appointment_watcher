package watcher_handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	watcher_show_template "github.com/troptropcontent/visa_appointment_watcher/internal/views/watcher/show"
)

func parseLogs() (logs []watcher_show_template.Log) {
	rawLogs := loadLogs()

	for _, rawLog := range rawLogs {
		var log watcher_show_template.Log
		json.Unmarshal([]byte(rawLog), &log)
		logs = append(logs, log)
	}
	return logs
}

func loadLogs() []string {
	rawLogs, err := os.ReadFile("logs/watcher/watcher.log")
	if err != nil {
		return []string{}
	}

	return strings.Split(string(rawLogs), "\n")
}

// Show returns the status of the watcher and the recent logs
func Show(c echo.Context) error {
	is_activated, _ := config.Get("watcher_running")
	last_appointment_date_found, _ := config.Get("last_appointment_date_found")
	last_appointment_date_found_at, _ := config.Get("last_appointment_date_found_at")
	data := &watcher_show_template.Template{
		Title:                      "Watcher",
		Logs:                       parseLogs(),
		IsActivated:                is_activated == "true",
		LastAppointmentDateFound:   last_appointment_date_found,
		LastAppointmentDateFoundAt: last_appointment_date_found_at,
	}
	return c.Render(http.StatusOK, "watcher/show", data)
}
