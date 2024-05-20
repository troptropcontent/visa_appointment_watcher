package watcher

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/troptropcontent/visa_appointment_watcher/internal/credentials"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type GmailNotifier struct {
	client *gmail.Service
}

func NewGmailNotifier() (notifier *GmailNotifier, err error) {
	ctx := context.Background()

	config := oauth2.Config{
		ClientID:     credentials.Config.GOOGLE_CLIENT_ID,
		ClientSecret: credentials.Config.GOOGLE_CLIENT_SECRET,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send"},
	}
	token := oauth2.Token{
		AccessToken:  credentials.Config.GOOGLE_ACCESS_TOKEN,
		RefreshToken: credentials.Config.GOOGLE_REFRESH_TOKEN,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}
	var tokenSource = config.TokenSource(ctx, &token)
	srv, err := gmail.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, fmt.Errorf("failed to receive gmail client: %v", err)
	}

	return &GmailNotifier{
		client: srv,
	}, nil
}

func (n *GmailNotifier) Notify(current_date time.Time, next_date time.Time, destination string) error {
	to := destination
	var msgString string
	emailTo := "To: " + to + "\r\n"
	msgString = emailTo
	subject := "Subject: " + "Visa appointment available on " + next_date.Format(time.DateOnly) + "\n"
	msgString = msgString + subject
	msgString = msgString + "\n" + "You have an appointment available on " + next_date.Format(time.DateOnly)

	msg := []byte(msgString)

	//Stores the entire message
	message := gmail.Message{
		Raw: base64.URLEncoding.EncodeToString([]byte(msg)),
	}

	//"me" sets the sender email address, email that was used to create the crendentials
	_, err := n.client.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return err
	}
	return nil
}
