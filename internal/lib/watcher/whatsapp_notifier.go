package watcher

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/troptropcontent/visa_appointment_watcher/internal/credentials"
)

type WhatsappNotifier struct {
	metaAccountId string
	token         string
}

func NewWhatsappNotifier() *WhatsappNotifier {
	notifier := WhatsappNotifier{
		metaAccountId: credentials.Config.META_ACCOUNT_ID,
		token:         credentials.Config.META_TOKEN,
	}
	return &notifier
}

func (w *WhatsappNotifier) Notify(current_date time.Time, next_date time.Time, alert_phone_number string) error {
	url := "https://graph.facebook.com/v19.0/" + w.metaAccountId + "/messages"
	fmt.Println(url)
	bearerToken := w.token
	contentType := "application/json"
	jsonData := `{
		"messaging_product": "whatsapp",
		"recipient_type": "individual",
		"to": "` + alert_phone_number + `",
		"type": "template",
		"template": {
			"name": "vaw",
			"language": {
				"code": "en"
			},
			"components": [
				{
					"type": "body",
					"parameters": [
						{
							"type": "text",
							"text": "` + next_date.Format("02-01-2006") + `"
						},
						{
							"type": "text",
							"text": "` + current_date.Format("02-01-2006") + `"
						}
					]
				},
				{
					"type": "header",
					"parameters": [
						{
							"type": "text",
							"text": "` + current_date.Format("02-01-2006") + `"
						}
					]
				}
			]
		}
	}`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Content-Type", contentType)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("WhatsApp message sent successfully")
	fmt.Println(resp)

	return nil

}