package main

import (
	"flag"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	"github.com/troptropcontent/visa_appointment_watcher/internal/lib/appointment_date_scrapper"
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

func scrapeAppointmentDates() (current_date time.Time, next_date time.Time, err error) {
	current_date, next_date, err = appointment_date_scrapper.FindDates()
	if err != nil {
		return time.Time{}, time.Time{}, err
	}
	return current_date, next_date, nil
}

func shouldNotify(current_date, next_date time.Time) bool {
	log.Info().Msg("Checking if current date is before next date : " + current_date.Format("02-01-2006") + " < " + next_date.Format("02-01-2006"))
	if current_date.Before(next_date) {
		log.Info().Msg("Current date is before next date, not sending alert")
		return false
	}

	last_alert_sent_for_appointment_date_at, err := config.Get("last_alert_sent_for_appointment_date_at")
	if err != nil {
		return false
	}

	is_next_appointment_date_already_notified := last_alert_sent_for_appointment_date_at == next_date.Format("02-01-2006")

	if is_next_appointment_date_already_notified {
		log.Info().Msg("Next appointment date have already been notified, not sending alert")
		return false
	}
	log.Info().Msg("Next appointment date have not been notified, sending alert")
	return true
}

func notify(current_date, next_date time.Time) error {
	_, err := config.Get("alert_phone_number")
	if err != nil {
		return err
	}

	message := "Hello Amigo, an earlier appointment date is available on " + next_date.Format("02-01-2006") + " (for now your appointment is scheduled on " + current_date.Format("02-01-2006") + ")"
	log.Info().Msg("About to send the following message: " + message)
	// Twillio stuff will go here

	config.MustSet("last_alert_sent_for_appointment_date_at", next_date.Format("02-01-2006"))
	return nil
}

func scrapAndNotify() error {
	log.Info().Msg("Scraping appointment dates")
	current_date, next_date, err := scrapeAppointmentDates()
	if err != nil {
		log.Error().Msg("Failed to scrape appointment dates: " + err.Error())
		return err
	}

	log.Info().Msg("Checking if we should notify")
	should_notify := shouldNotify(current_date, next_date)

	if should_notify {
		return notify(current_date, next_date)
	}
	return nil
}

func startWatcher() *time.Ticker {
	appointment_date_ticker := time.NewTicker(45 * time.Second)
	go func() {
		for ; ; <-appointment_date_ticker.C {
			err := scrapAndNotify()
			if err != nil {
				log.Error().Msg(err.Error())
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
