package gorm_validator_test

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/troptropcontent/visa_appointment_watcher/database"
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	"github.com/troptropcontent/visa_appointment_watcher/internal/models"
	"gorm.io/gorm"
)

func TestUniqueness(t *testing.T) {
	config.MustLoad()
	config.MustInit()
	database.Init()
	models.Init()

	t.Run("When email is already taken", func(t *testing.T) {
		database.DB.Transaction(func(transaction *gorm.DB) error {
			another_user := models.User{
				Email:             "test@example.com",
				EncryptedPassword: "password",
				SignedUpThrough:   "local",
			}
			transaction.Create(&another_user)
			user := models.User{
				Email:             "test@example.com",
				EncryptedPassword: "password",
				SignedUpThrough:   "local",
			}

			err := transaction.Create(&user).Error
			validationErrors := err.(validator.ValidationErrors)
			uniquessError := false
			for _, validationError := range validationErrors {
				if validationError.Tag() == "uniqueness" && validationError.Field() == "Email" {
					uniquessError = true
				}
			}
			if !uniquessError {
				t.Errorf("Expected uniqueness error")
			}

			return errors.New("test")
		})
	})
	t.Run("When email is not taken", func(t *testing.T) {
		database.DB.Transaction(func(transaction *gorm.DB) error {
			another_user := models.User{
				Email:             "test2@example.com",
				EncryptedPassword: "password",
				SignedUpThrough:   "local",
			}
			transaction.Create(&another_user)

			user := models.User{
				Email:             "test@example.com",
				EncryptedPassword: "password",
				SignedUpThrough:   "local",
			}

			err := transaction.Create(&user).Error

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			return errors.New("test")
		})
	})
}
