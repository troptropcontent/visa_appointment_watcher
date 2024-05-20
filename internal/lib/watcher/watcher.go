package watcher

import (
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	"github.com/troptropcontent/visa_appointment_watcher/internal/lib/logging"
)

const WATCHER_LOG_FILE = "logs/watcher/watcher.log"

type Watcher struct {
	ID          uuid.UUID
	IsActivated bool
	Logger      *zerolog.Logger
}

func (w *Watcher) shouldNotify(current_date, next_date time.Time) bool {
	if next_date.IsZero() {
		return false
	}

	if current_date.Before(next_date) {
		return false
	}

	last_alert_sent_for_appointment_date_at, err := config.Get("last_alert_sent_for_appointment_date_at")
	if err != nil {
		return false
	}

	is_next_appointment_date_already_notified := last_alert_sent_for_appointment_date_at == next_date.Format("02-01-2006")

	return !is_next_appointment_date_already_notified
}

func New() Watcher {
	logger := logging.New(logging.Config{
		Filename: WATCHER_LOG_FILE,
	})
	isActivated, _ := config.Get("watcher_running")
	return Watcher{Logger: &logger, IsActivated: isActivated == "true"}
}

func (w *Watcher) Run() error {
	w.Logger.Info().Msg("Watcher starting")
	if !w.IsActivated {
		w.Logger.Info().Msg("Watcher is not activated, operation stopped")
		return nil
	}

	w.Logger.Info().Msg("Watcher is activated, scraping the appointment dates")
	scrapper := NewScrapper(w)
	err := scrapper.FindDates()
	if err != nil {
		w.Logger.Error().Msgf("Scraper failed with the following error : %s", err.Error())
		return err
	}

	w.Logger.Info().Msgf("Scraper successfully find the appointment dates, current date: %s, next date: %s", scrapper.CurrentDate, scrapper.NextDate)

	if !w.shouldNotify(scrapper.CurrentDate, scrapper.NextDate) {
		w.Logger.Info().Msg("Watcher should not notify, operation stopped")
		return nil
	}

	err = Notify(scrapper.CurrentDate, scrapper.NextDate)

	if err != nil {
		w.Logger.Error().Msgf("Failed to notify: %s", err.Error())
		return err
	}

	w.Logger.Info().Msg("Watcher successfully notified, operation finished")

	return nil
}
