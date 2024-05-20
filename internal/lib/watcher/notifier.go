package watcher

import (
	"fmt"
	"sync"
	"time"

	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	"github.com/troptropcontent/visa_appointment_watcher/internal/credentials"
)

type Notifier interface {
	Notify(current_date time.Time, next_date time.Time, destination string) error
}

type Alert struct {
	Destination string
	Service     Notifier
}

func (a *Alert) Notify(current_date time.Time, next_date time.Time) error {
	return a.Service.Notify(current_date, next_date, a.Destination)
}

func Notify(current_date time.Time, next_date time.Time) error {
	whatsapp := NewWhatsappNotifier()
	gmail, err := NewGmailNotifier()
	if err != nil {
		return fmt.Errorf("failed to receive gmail client: %v", err)
	}
	alerts := []Alert{
		{
			Destination: credentials.Config.ADMIN_PHONE_NUMBER,
			Service:     whatsapp,
		},
		{
			Destination: credentials.Config.ADMIN_EMAIL,
			Service:     gmail,
		},
		{
			Destination: config.MustGet("alert_phone_number"),
			Service:     whatsapp,
		},
		{
			Destination: config.MustGet("username"),
			Service:     gmail,
		},
	}

	wg := sync.WaitGroup{}
	errChan := make(chan error)
	for _, alert := range alerts {
		wg.Add(1)
		go func(alert Alert) {
			defer wg.Done()
			if err := alert.Notify(current_date, next_date); err != nil {
				errChan <- err
			}
		}(alert)
	}
	wg.Wait()
	close(errChan)
	if len(errChan) > 0 {
		return fmt.Errorf("failed to send alert: %v", <-errChan)
	}
	return nil
}
