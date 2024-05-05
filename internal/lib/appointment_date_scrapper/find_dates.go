package appointment_date_scrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
)

const (
	LOG_FILE = "logs/appointment_date_scrapper.log"
)

type AppointmentDateResponse []struct {
	Date        string `json:"date"`
	BusinessDay bool   `json:"business_day"`
}

func setupBrowser(date *time.Time, logger *zerolog.Logger) *rod.Browser {
	launcher := launcher.MustNewManaged("http://headlessBrowser:7317")

	launcher.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

	client := launcher.MustClient()

	browser := rod.New().Client(client).MustConnect()

	router := browser.HijackRequests()

	router.MustAdd("*/en-fr/niv/schedule/*/appointment/days/*.json*", func(ctx *rod.Hijack) {
		if err := ctx.LoadResponse(http.DefaultClient, true); err != nil {
			return
		}
		responseBody := ctx.Response.Body()
		var parsedResponse AppointmentDateResponse
		json.Unmarshal([]byte(responseBody), &parsedResponse)
		rawDate := parsedResponse[0].Date
		parsedDate, err := time.Parse("2006-01-02", rawDate)
		if err != nil {
			return
		}
		logger.Info().Msg(fmt.Sprintf("Next appointment date found: %s", parsedDate))
		*date = parsedDate
	})

	go router.Run()

	return browser
}

func findCurrentDate(page *rod.Page) (time.Time, error) {
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

func newLogger(scraper_process_id string) (*zerolog.Logger, error) {
	log_file, err := os.OpenFile(LOG_FILE, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	logger := zerolog.New(log_file).With().Timestamp().Caller().Str("scraper_process_id", scraper_process_id).Logger()

	return &logger, nil
}

// FindDates triggers the scraping process to find the current and next appointment date. It follows the steps of a user browsing the visa website and find the dates in the browser.
func FindDates() (current_date time.Time, next_date time.Time, err error) {
	username := config.MustGet("username")
	password := config.MustGet("password")

	scraper_process_id := uuid.New()
	logger, _ := newLogger(scraper_process_id.String())

	logger.Info().Msg("Starting the scrapper")

	browser := setupBrowser(&next_date, logger)

	logger.Info().Msg("Browser setup complete")
	logger.Info().Msg("Navigating to login page")

	page := browser.MustPage("https://ais.usvisa-info.com/en-fr/niv/users/sign_in")

	logger.Info().Msg("Waiting for the page to be stable")

	page.MustWaitStable()

	current_url := page.MustInfo().URL
	if current_url != "https://ais.usvisa-info.com/en-fr/niv/users/sign_in" {
		errMessage := "Failed to reach the login page"
		logger.Error().Msg(errMessage)
		return time.Time{}, time.Time{}, errors.New(errMessage)
	}

	logger.Info().Msg("Login page reached")
	logger.Info().Msg("Filling the login form")

	page.MustElement("input[type='email']").MustInput(username)
	page.MustElement("input[type='password']").MustInput(password)
	page.MustElement("label[for='policy_confirmed']").MustClick()

	logger.Info().Msg("Clicking on the login button")

	page.MustElement("input[type='submit']").MustClick()

	logger.Info().Msg("Waiting for the page to be stable")

	page.MustWaitStable()

	current_url = page.MustInfo().URL
	regex := regexp.MustCompile(`https://ais.usvisa-info.com/en-fr/niv/groups/\d+`)
	if !regex.MatchString(current_url) {
		errMessage := "next page not reached, login might have failed, credentials might be incorrect"
		logger.Error().Msg(errMessage)
		return time.Time{}, time.Time{}, errors.New(errMessage)
	}

	logger.Info().Msg("Login successful, next page reached")

	logger.Info().Msg("Finding the current appointment date")

	current_date, err = findCurrentDate(page)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	logger.Info().Msg(fmt.Sprintf("Current appointment date found: %s", current_date))
	logger.Info().Msg("Clicking on the continue button")

	page.MustElementR("a", "Continue").MustClick()

	logger.Info().Msg("Waiting for the page to be stable")

	page.MustWaitStable()

	current_url = page.MustInfo().URL
	regex = regexp.MustCompile(`https://ais.usvisa-info.com/en-fr/niv/schedule/\d+/continue_actions`)
	if !regex.MatchString(current_url) {
		errMessage := "failed to reach the next page, current url: " + current_url
		logger.Error().Msg(errMessage)
		return time.Time{}, time.Time{}, errors.New(errMessage)
	}

	logger.Info().Msg("Next page reached")
	logger.Info().Msg("Clicking on the reschedule appointment accordion")

	page.MustElementR("h5", "Reschedule Appointment").MustClick()

	logger.Info().Msg("Reschedule appointment accordion clicked")
	logger.Info().Msg("Waiting for the page to be stable")

	page.MustWaitStable()

	logger.Info().Msg("Clicking on the reschedule appointment button")

	page.MustElementR("a[href*='appointment']", "Reschedule Appointment").MustClick()

	logger.Info().Msg("Reschedule appointment button clicked")
	logger.Info().Msg("Waiting for the page to be stable")

	page.MustWaitStable()

	current_url = page.MustInfo().URL
	regex = regexp.MustCompile(`https://ais.usvisa-info.com/en-fr/niv/schedule/\d+/appointment`)
	if !regex.MatchString(current_url) {
		errMessage := "failed to reach the next page, current url: " + current_url
		logger.Error().Msg(errMessage)
		return time.Time{}, time.Time{}, errors.New(errMessage)
	}

	logger.Info().Msg("Scrapper finished")

	return current_date, next_date, nil
}
