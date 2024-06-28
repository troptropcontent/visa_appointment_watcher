package main

import (
	"crypto/subtle"
	"flag"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/troptropcontent/visa_appointment_watcher/database"
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	"github.com/troptropcontent/visa_appointment_watcher/internal/credentials"
	watcher_handler "github.com/troptropcontent/visa_appointment_watcher/internal/handler/watcher"
	"github.com/troptropcontent/visa_appointment_watcher/internal/lib/logging"
	"github.com/troptropcontent/visa_appointment_watcher/internal/lib/watcher"
	"github.com/troptropcontent/visa_appointment_watcher/internal/models"
	auth_api "github.com/troptropcontent/visa_appointment_watcher/internal/modules/auth/api"
	"github.com/troptropcontent/visa_appointment_watcher/internal/views"
)

func mustGetOptionsFromEnvOrFlags() (username string, password string, alert_phone_number string) {
	flag_username := flag.String("username", "", "your username")
	flag_password := flag.String("password", "", "your password")
	flag_alert_phone_number := flag.String("alert_phone_number", "", "the number to send the alerts to")
	flag.Parse()

	if *flag_username != "" {
		username = *flag_username
	}
	if *flag_password != "" {
		password = *flag_password
	}
	if *flag_alert_phone_number != "" {
		alert_phone_number = *flag_alert_phone_number
	}

	if env_username := os.Getenv("VISA_APPOINTMENT_WATCHER_USERNAME"); env_username != "" {
		username = env_username
	}

	if env_password := os.Getenv("VISA_APPOINTMENT_WATCHER_PASSWORD"); env_password != "" {
		password = env_password
	}
	if env_alert_phone_number := os.Getenv("VISA_APPOINTMENT_WATCHER_ALERT_PHONE_NUMBER"); env_alert_phone_number != "" {
		alert_phone_number = env_alert_phone_number
	}

	if username == "" {
		panic("username is empty")
	}
	if password == "" {
		panic("password is empty")
	}
	if alert_phone_number == "" {
		panic("alert_phone_number is empty")
	}

	return username, password, alert_phone_number
}

func startWatcherTicker() *time.Ticker {
	appointment_date_ticker := time.NewTicker(15 * time.Minute)
	go func() {
		for range appointment_date_ticker.C {
			logging.Logger.Info().Msg("Trigger a new watcher run")
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
	config.MustInit()
	logging.Init(logging.Config{})
	database.Init()
	models.Init()
	username, password, alert_phone_number := mustGetOptionsFromEnvOrFlags()
	config.MustSetIfNotExists("watcher_running", "false")
	config.MustSet("username", username)
	config.MustSet("password", password)
	config.MustSet("alert_phone_number", alert_phone_number)
	config.MustSetIfNotExists("last_alert_sent_for_appointment_date_at", "")
	config.MustSetIfNotExists("last_appointment_date_found", "")
	config.MustSetIfNotExists("last_alert_sent_for_appointment_date_at", "")

	credentials.MustInit()

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

	// Auth endpoints
	auth_api.RegisterRoutes(server)

	server.Start(":1234")
}
