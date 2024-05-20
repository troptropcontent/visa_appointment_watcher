package credentials

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Credentials struct {
	TWILIO_ACCOUNT_SID   string `validate:"required"`
	TWILIO_AUTH_TOKEN    string `validate:"required"`
	TWILIO_NUMBER        string `validate:"required"`
	META_ACCOUNT_ID      string `validate:"required"`
	META_TOKEN           string `validate:"required"`
	GOOGLE_CLIENT_ID     string `validate:"required"`
	GOOGLE_CLIENT_SECRET string `validate:"required"`
	GOOGLE_ACCESS_TOKEN  string `validate:"required"`
	GOOGLE_REFRESH_TOKEN string `validate:"required"`
	ADMIN_PHONE_NUMBER   string `validate:"required"`
	ADMIN_EMAIL          string `validate:"required"`
}

var Config Credentials

func Init() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("VAW")
	Config = Credentials{
		TWILIO_ACCOUNT_SID:   viper.GetString("TWILIO_ACCOUNT_SID"),
		TWILIO_AUTH_TOKEN:    viper.GetString("TWILIO_AUTH_TOKEN"),
		TWILIO_NUMBER:        viper.GetString("TWILIO_NUMBER"),
		META_ACCOUNT_ID:      viper.GetString("META_ACCOUNT_ID"),
		META_TOKEN:           viper.GetString("META_TOKEN"),
		GOOGLE_CLIENT_ID:     viper.GetString("GOOGLE_CLIENT_ID"),
		GOOGLE_CLIENT_SECRET: viper.GetString("GOOGLE_CLIENT_SECRET"),
		GOOGLE_ACCESS_TOKEN:  viper.GetString("GOOGLE_ACCESS_TOKEN"),
		GOOGLE_REFRESH_TOKEN: viper.GetString("GOOGLE_REFRESH_TOKEN"),
		ADMIN_PHONE_NUMBER:   viper.GetString("ADMIN_PHONE_NUMBER"),
		ADMIN_EMAIL:          viper.GetString("ADMIN_EMAIL"),
	}
	err := validator.New().Struct(Config)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		err_message := "Missing credential: "
		for i, err := range validationErrors {
			err_message += err.Field()
			if i < len(validationErrors)-1 {
				err_message += ", "
			}
		}

		return errors.New(err_message)
	}
	return nil
}

func MustInit() {
	err := Init()
	if err != nil {
		panic(err)
	}
}
