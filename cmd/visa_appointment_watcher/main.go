package main

import (
	"flag"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	"github.com/troptropcontent/visa_appointment_watcher/internal/lib/appointment_date_scrapper"
)

const WATCHER_LOG_FILE = "logs/watcher.log"

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

func scrapeAppointmentDates(watcher_process_id string) (current_date time.Time, next_date time.Time, err error) {
	current_date, next_date, err = appointment_date_scrapper.FindDates(watcher_process_id)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return current_date, next_date, nil
}

func shouldNotify(current_date, next_date time.Time, watcher_logger *zerolog.Logger) bool {
	watcher_logger.Info().Msg("Checking if current date is before next date : " + current_date.Format("02-01-2006") + " < " + next_date.Format("02-01-2006"))
	if current_date.Before(next_date) {
		watcher_logger.Info().Msg("Current date is before next date, not sending alert")
		return false
	}

	last_alert_sent_for_appointment_date_at, err := config.Get("last_alert_sent_for_appointment_date_at")
	if err != nil {
		return false
	}

	is_next_appointment_date_already_notified := last_alert_sent_for_appointment_date_at == next_date.Format("02-01-2006")

	if is_next_appointment_date_already_notified {
		watcher_logger.Info().Msg("Next appointment date have already been notified, not sending alert")
		return false
	}
	watcher_logger.Info().Msg("Next appointment date have not been notified, sending alert")
	return true
}

func notify(current_date, next_date time.Time, watcher_logger *zerolog.Logger) error {
	_, err := config.Get("alert_phone_number")
	if err != nil {
		return err
	}

	message := "Hello Amigo, an earlier appointment date is available on " + next_date.Format("02-01-2006") + " (for now your appointment is scheduled on " + current_date.Format("02-01-2006") + ")"
	watcher_logger.Info().Msg("About to send the following message: " + message)
	// Twillio stuff will go here

	config.Set("last_alert_sent_for_appointment_date_at", next_date.Format("02-01-2006"))
	return nil
}

func scrapAndNotify(watcher_logger *zerolog.Logger, watcher_process_id string) error {
	watcher_logger.Info().Msg("Scraping appointment dates")
	current_date, next_date, err := scrapeAppointmentDates(watcher_process_id)
	if err != nil {
		watcher_logger.Error().Msg("Failed to scrape appointment dates: " + err.Error())
		return err
	}

	watcher_logger.Info().Msg("Checking if we should notify")
	should_notify := shouldNotify(current_date, next_date, watcher_logger)

	if should_notify {
		return notify(current_date, next_date, watcher_logger)
	}
	return nil
}

func newWatcherLogger(watcher_process_id string) (*zerolog.Logger, error) {
	log_file, err := os.OpenFile(WATCHER_LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	logger := zerolog.New(log_file).With().Timestamp().Caller().Str("watcher_process_id", watcher_process_id).Logger()

	return &logger, nil
}

func startWatcher() *time.Ticker {
	appointment_date_ticker := time.NewTicker(45 * time.Second)
	go func() {
		for ; ; <-appointment_date_ticker.C {
			watcher_process_id := uuid.New()
			watcher_logger, err := newWatcherLogger(watcher_process_id.String())
			if err != nil {
				log.Info().Msg("Failed to create watcher logger: " + err.Error())
				continue
			}
			watcher_logger.Info().Msg("Watcher starting")
			watcher_should_run, err := config.Get("watcher_running")
			if err != nil || watcher_should_run == "false" {
				watcher_logger.Info().Msg("Watcher is not running, not scraping appointment dates")
				continue
			}
			err = scrapAndNotify(watcher_logger, watcher_process_id.String())
			if err != nil {
				watcher_logger.Error().Msg(err.Error())
			}
		}
	}()
	return appointment_date_ticker
}

func main() {
	username, password, alert_phone_number := mustGetAllParamsInFlags()
	config.MustInit()
	config.MustSetIfNotExists("watcher_running", "false")
	config.MustSetIfNotExists("username", username)
	config.MustSetIfNotExists("password", password)
	config.MustSetIfNotExists("alert_phone_number", alert_phone_number)
	config.MustSetIfNotExists("last_alert_sent_for_appointment_date_at", "")

	watcher := startWatcher()
	defer watcher.Stop()

	time.Sleep(120 * time.Second)
	watcher.Stop()
}
