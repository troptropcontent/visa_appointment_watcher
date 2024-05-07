package watcher_show_template

import "time"

type Log struct {
	Level   string    `json:"level"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

type Template struct {
	Title                      string
	Logs                       []Log
	IsActivated                bool
	LastAppointmentDateFound   string
	LastAppointmentDateFoundAt string
}

func (t *Template) TemplateFiles() []string {
	return []string{
		"internal/views/layouts/application.html",
		"internal/views/watcher/show/template.html",
		"internal/views/watcher/show/_switcher/template.html",
	}
}
