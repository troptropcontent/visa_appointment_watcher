package watcher

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/rs/zerolog"
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	"github.com/troptropcontent/visa_appointment_watcher/internal/lib/logging"
)

const (
	SCRAPPER_LOG_FILE = "logs/scrapper/scrapper.log"
)

type AppointmentDateResponse []struct {
	Date        string `json:"date"`
	BusinessDay bool   `json:"business_day"`
}

type Scrapper struct {
	Watcher     *Watcher
	Logger      *zerolog.Logger
	Browser     *rod.Browser
	CurrentDate time.Time
	NextDate    time.Time
}

func NewScrapper(watcher *Watcher) Scrapper {
	logger := logging.New(logging.Config{
		Filename: SCRAPPER_LOG_FILE,
	}).With().Str("watcher_process_id", watcher.ID.String()).Logger()

	return Scrapper{Watcher: watcher, Logger: &logger}
}

func findCurrentDate(p *Page) (time.Time, error) {
	page := p.Page
	regex := regexp.MustCompile(`Consular Appointment: (\d+) (\w+), (\d+), (\d+:\d+) PARIS local time at Paris —  get directions`)
	current_date := page.MustElement(".consular-appt").MustText()
	matches := regex.FindStringSubmatch(current_date)
	if len(matches) > 1 {
		day, _ := strconv.Atoi(matches[1])
		month := matches[2]
		year, _ := strconv.Atoi(matches[3])
		date, err := time.Parse("02-January-2006", fmt.Sprintf("%02d-%s-%d", day, month, year))
		if err != nil {
			return time.Time{}, err
		}
		return date, nil
	}

	return time.Time{}, errors.New("no date found")
}

func (s *Scrapper) FindDates() (err error) {
	username := config.MustGet("username")
	password := config.MustGet("password")

	browser := NewBrowser(s)
	defer browser.Close()

	page := browser.MustOpenPage("https://ais.usvisa-info.com/en-fr/niv/users/sign_in")

	page.MustWaitStable()

	if current_url := page.URL(); current_url != "https://ais.usvisa-info.com/en-fr/niv/users/sign_in" {
		errMessage := "failed to reach the login page, current url: " + current_url
		s.Logger.Error().Msg(errMessage)
		return errors.New(errMessage)
	}

	s.Logger.Info().Msg("Login page reached")
	s.Logger.Info().Msg("Filling the login form")

	page.MustFindElement("input[type='email']").MustInput(username)
	page.MustFindElement("input[type='password']").MustInput(password)
	page.MustFindElement("label[for='policy_confirmed']").MustClick()
	page.MustFindElement("input[type='submit']").MustClick()
	page.MustWaitStable()

	regex := regexp.MustCompile(`https://ais.usvisa-info.com/en-fr/niv/groups/\d+`)
	if current_url := page.URL(); !regex.MatchString(current_url) {
		errMessage := "next page not reached, login might have failed, credentials might be incorrect"
		s.Logger.Error().Msg(errMessage)
		return errors.New(errMessage)
	}

	s.Logger.Info().Msg("Login successful, next page reached")
	s.Logger.Info().Msg("Finding the current appointment date")

	scrapped_current_date, err := findCurrentDate(&page)
	if err != nil {
		s.Logger.Error().Msg("Failed to find the current appointment date")
		return err
	}

	s.Logger.Info().Msg(fmt.Sprintf("Current appointment date found: %s", scrapped_current_date))
	s.CurrentDate = scrapped_current_date

	s.Logger.Info().Msg("Clicking on the continue button")

	page.MustFindElementByText("a", "Continue").MustClick()

	s.Logger.Info().Msg("Waiting for the page to be stable")

	page.MustWaitStable()

	regex = regexp.MustCompile(`https://ais.usvisa-info.com/en-fr/niv/schedule/\d+/continue_actions`)
	if current_url := page.URL(); !regex.MatchString(current_url) {
		errMessage := "failed to reach the next page, current url: " + current_url
		s.Logger.Error().Msg(errMessage)
		return errors.New(errMessage)
	}

	s.Logger.Info().Msg("Next page reached")
	s.Logger.Info().Msg("Clicking on the reschedule appointment accordion")

	page.MustFindElementByText("h5", "Reschedule Appointment").MustClick()

	s.Logger.Info().Msg("Reschedule appointment accordion clicked")
	s.Logger.Info().Msg("Waiting for the page to be stable")

	page.MustWaitStable()

	s.Logger.Info().Msg("Clicking on the reschedule appointment button")

	page.MustFindElementByText("a[href*='appointment']", "Reschedule Appointment").MustClick()

	s.Logger.Info().Msg("Reschedule appointment button clicked")
	s.Logger.Info().Msg("Waiting for the page to be stable")

	page.MustWaitStable()

	regex = regexp.MustCompile(`https://ais.usvisa-info.com/en-fr/niv/schedule/\d+/appointment`)
	if current_url := page.URL(); !regex.MatchString(current_url) {
		errMessage := "failed to reach the next page, current url: " + current_url
		s.Logger.Error().Msg(errMessage)
		return errors.New(errMessage)
	}

	s.Logger.Info().Msg("Scrapper finished")

	return nil
}
