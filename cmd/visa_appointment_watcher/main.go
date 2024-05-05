package main

import (
	"flag"
	"fmt"
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

func main() {
	username, password, alert_phone_number := mustGetAllParamsInFlags()
	config.MustInit()
	config.MustSetIfNotExists("watcher_running", "false")
	config.MustSetIfNotExists("username", username)
	config.MustSetIfNotExists("password", password)
	config.MustSetIfNotExists("alert_phone_number", alert_phone_number)
	config.MustSetIfNotExists("last_alert_sent_at", "")

	appointment_date_ticker := time.NewTicker(time.Minute * 20)
	defer appointment_date_ticker.Stop()

	go func() {
		for t := range appointment_date_ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	current_date, next_date, err := appointment_date_scrapper.FindDates()
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	log.Info().Msgf("Current appointment date: %s", current_date.Format("02-01-2006"))
	log.Info().Msgf("Next available date: %s", next_date.Format("02-01-2006"))
}
