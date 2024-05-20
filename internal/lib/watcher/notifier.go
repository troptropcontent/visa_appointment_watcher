package watcher

import (
	"time"
)

type Notifier interface {
	Notify(current_date time.Time, next_date time.Time, alert_phone_number string) error
}

func NewNotifier() Notifier {
	client := NewWhatsappNotifier()

	return client
}
