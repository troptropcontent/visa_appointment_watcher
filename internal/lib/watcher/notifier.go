package watcher

import "time"

type Notifier interface {
	Notify(current_date time.Time, next_date time.Time) error
}

type twilioNotifier struct{}

func NewNotifier() Notifier {
	twilioNotifier := twilioNotifier{}
	return &twilioNotifier
}

func (n *twilioNotifier) Notify(current_date time.Time, next_date time.Time) error {
	return nil
}
