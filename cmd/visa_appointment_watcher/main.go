package main

import (
	"flag"

	"github.com/rs/zerolog/log"
	"github.com/troptropcontent/visa_appointment_watcher/internal/lib/appointment_date_scrapper"
)

func getCredentials() (string, string) {
	username := flag.String("username", "", "your username")
	password := flag.String("password", "", "your password")
	flag.Parse()
	return *username, *password
}

func main() {
	username, password := getCredentials()

	if username == "" {
		panic("username is empty")
	}
	if password == "" {
		panic("password is empty")
	}

	current_date, next_date, err := appointment_date_scrapper.FindDates(username, password)
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	log.Info().Msgf("Current appointment date: %s", current_date.Format("02-01-2006"))
	log.Info().Msgf("Next available date: %s", next_date.Format("02-01-2006"))
}
