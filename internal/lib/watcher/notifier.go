package watcher

import (
	"os"
	"time"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Notifier interface {
	Notify(current_date time.Time, next_date time.Time, alert_phone_number string) error
}

type twilioNotifier struct {
	client *twilio.RestClient
}

func NewNotifier() Notifier {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: os.Getenv("TWILIO_ACCOUNT_SID"),
		Password: os.Getenv("TWILIO_AUTH_TOKEN"),
	})
	twilioNotifier := twilioNotifier{client: client}
	return &twilioNotifier
}

func (n *twilioNotifier) Notify(current_date time.Time, next_date time.Time, alert_phone_number string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetFrom(os.Getenv("TWILIO_NUMBER"))
	params.SetTo(alert_phone_number)
	message := "Hugh ✋\nHurry up, a new appointment date for a visa is available the " + next_date.Format("02-01-2006") + " (for now, your appointment is scheduled for " + current_date.Format("02-01-2006") + ").\nTom"
	params.SetBody(message)

	_, err := n.client.Api.CreateMessage(params)
	if err != nil {
		return err
	} else {
		return nil
	}
}
