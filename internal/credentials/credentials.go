package credentials

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Credentials struct {
	TWILIO_ACCOUNT_SID string `validate:"required"`
	TWILIO_AUTH_TOKEN  string `validate:"required"`
	TWILIO_NUMBER      string `validate:"required"`
	META_ACCOUNT_ID    string `validate:"required"`
	META_TOKEN         string `validate:"required"`
}

var Config Credentials

func Init() error {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("VAW")
	Config = Credentials{
		TWILIO_ACCOUNT_SID: viper.GetString("TWILIO_ACCOUNT_SID"),
		TWILIO_AUTH_TOKEN:  viper.GetString("TWILIO_AUTH_TOKEN"),
		TWILIO_NUMBER:      viper.GetString("TWILIO_NUMBER"),
		META_ACCOUNT_ID:    viper.GetString("META_ACCOUNT_ID"),
		META_TOKEN:         viper.GetString("META_TOKEN"),
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
