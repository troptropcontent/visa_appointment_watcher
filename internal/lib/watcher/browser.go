package watcher

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const HEADLESS_BROWSER_URL = "http://headlessBrowser:7317"

type Browser struct {
	browser  *rod.Browser
	Scrapper *Scrapper
}

type Page struct {
	Page *rod.Page
}

type Element struct {
	Element *rod.Element
}

func NewBrowser(scrapper *Scrapper) Browser {
	launcher := launcher.MustNewManaged(HEADLESS_BROWSER_URL)

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

		scrapper.NextDate = parsedDate
	})

	go router.Run()

	return Browser{browser: browser, Scrapper: scrapper}
}

func (b Browser) Close() {
	b.browser.Close()
}

func (b Browser) MustOpenPage(url string) Page {
	page := b.browser.MustPage(url)
	return Page{
		Page: page,
	}
}

func (p Page) MustWaitStable() {
	p.Page.MustWaitStable()
}

func (p Page) URL() string {
	return p.Page.MustInfo().URL
}

func (p Page) MustFindElement(selector string) Element {
	return Element{Element: p.Page.MustElement(selector)}
}

func (p Page) MustTakeScreenshot(args struct {
	Path     string
	FullPage bool
}) {
	if args.FullPage {
		p.Page.MustScreenshotFullPage(args.Path)
	} else {
		p.Page.MustScreenshot(args.Path)
	}
}

func (p *Page) MustFindElementByText(selector string, text string) Element {
	return Element{Element: p.Page.MustElementR(selector, text)}
}

func (e Element) MustInput(value string) {
	e.Element.MustInput(value)
}

func (e Element) MustClick() {
	e.Element.MustClick()
}