package watcher

import (
	"testing"
	"time"

	"github.com/troptropcontent/visa_appointment_watcher/internal/credentials"
)

func TestWhatsappNotifier_Notify(t *testing.T) {
	credentials.MustInit()
	whatsappNotifier := NewWhatsappNotifier()
	next_date := time.Now()
	current_date := next_date.AddDate(0, 0, 1)
	err := whatsappNotifier.Notify(current_date, next_date, "19176992382")
	if err != nil {
		t.Errorf("Error sending WhatsApp message: %v", err)
	}
}
