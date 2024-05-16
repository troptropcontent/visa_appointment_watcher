package watcher

import (
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

const HEADLESS_BROWSER_URL_DEVELOPMENT = "http://headlessBrowser:7317"
const HEADLESS_BROWSER_URL_PRODUCTION = "http://visa_appointment_watcher-browser:7317"

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

func gettHeadlessBrowserURL() string {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return HEADLESS_BROWSER_URL_PRODUCTION
	}
	return HEADLESS_BROWSER_URL_DEVELOPMENT
}

func NewBrowser(scrapper *Scrapper) Browser {
	launcher := launcher.MustNewManaged(gettHeadlessBrowserURL())

	launcher.Headless(false).XVFB("--server-num=5", "--server-args=-screen 0 1600x900x16")

	client := launcher.MustClient()

	browser := rod.New().Client(client).MustConnect()

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
