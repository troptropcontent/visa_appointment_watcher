package views

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/labstack/echo/v4"
	watcher_show_template "github.com/troptropcontent/visa_appointment_watcher/internal/views/watcher/show"
	watcher_show_switcher_template "github.com/troptropcontent/visa_appointment_watcher/internal/views/watcher/show/_switcher"
)

type Template struct {
	templates map[string]*template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return fmt.Errorf("template %s not found", name)
	}

	is_partial := strings.HasPrefix(strings.Split(name, "/")[len(strings.Split(name, "/"))-1], "_")
	if is_partial {
		return tmpl.ExecuteTemplate(w, name, data)
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}

type Templatable interface {
	TemplateFiles() []string
}

func templateFor(t Templatable) *template.Template {
	return template.Must(template.ParseFiles(t.TemplateFiles()...))
}

// NewRenderer returns a new renderer, it is responsible of registering all the templates.
func NewRenderer() echo.Renderer {
	t := &Template{
		templates: map[string]*template.Template{
			"watcher/show":           templateFor(&watcher_show_template.Template{}),
			"watcher/show/_switcher": templateFor(&watcher_show_switcher_template.Template{}),
		},
	}

	return t
}
