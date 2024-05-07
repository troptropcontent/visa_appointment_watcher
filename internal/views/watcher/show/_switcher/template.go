package watcher_show_switcher_template

type Template struct {
	IsActivated bool
}

func (t *Template) TemplateFiles() []string {
	return []string{
		"internal/views/watcher/show/_switcher/template.html",
	}
}
