package main

import (
	"crypto/subtle"
	"flag"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	watcher_handler "github.com/troptropcontent/visa_appointment_watcher/internal/handler/watcher"
	"github.com/troptropcontent/visa_appointment_watcher/internal/lib/watcher"
	"github.com/troptropcontent/visa_appointment_watcher/internal/views"
)

func mustGetAllParamsInFlags() (string, string, string) {
	username := flag.String("username", "", "your username")
	password := flag.String("password", "", "your password")
	alert_phone_number := flag.String("alert_phone_number", "", "the number to send the alerts to")
	flag.Parse()
	if *username == "" {
		panic("username is empty")
	}
	if *password == "" {
		panic("password is empty")
	}
	if *alert_phone_number == "" {
		panic("alert_phone_number is empty")
	}
	return *username, *password, *alert_phone_number
}

func startWatcherTicker() *time.Ticker {
	appointment_date_ticker := time.NewTicker(45 * time.Second)
	go func() {
		for ; ; <-appointment_date_ticker.C {
			w := watcher.New()
			w.Run()
		}
	}()
	return appointment_date_ticker
}

func createServer() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return e
}

func authMiddleware() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		username_in_config, err := config.Get("username")
		if err != nil {
			return false, nil
		}
		password_in_config, err := config.Get("password")
		if err != nil {
			return false, nil
		}
		if subtle.ConstantTimeCompare([]byte(username), []byte(username_in_config)) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(password_in_config)) == 1 {
			return true, nil
		}
		return false, nil
	})
}

func main() {
	username, password, alert_phone_number := mustGetAllParamsInFlags()
	config.MustInit()
	config.MustSetIfNotExists("watcher_running", "false")
	config.MustSetIfNotExists("username", username)
	config.MustSetIfNotExists("password", password)
	config.MustSetIfNotExists("alert_phone_number", alert_phone_number)
	config.MustSetIfNotExists("last_alert_sent_for_appointment_date_at", "")
	config.MustSetIfNotExists("last_appointment_date_found", "")
	config.MustSetIfNotExists("last_alert_sent_for_appointment_date_at", "")

	watcherTicker := startWatcherTicker()
	defer watcherTicker.Stop()

	server := createServer()

	server.Renderer = views.NewRenderer()

	// Static files
	server.Static("/public", "public")

	// Health check
	server.GET("/up", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	// Watcher web
	watcher_web := server.Group("/watcher")
	watcher_web.Use(authMiddleware())
	watcher_web.GET("", watcher_handler.Show)
	watcher_web.POST("/activate", watcher_handler.Activate)
	watcher_web.POST("/deactivate", watcher_handler.Deactivate)
	server.Start(":3000")
}
