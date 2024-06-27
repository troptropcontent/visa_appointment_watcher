package models

import (
	"github.com/troptropcontent/visa_appointment_watcher/internal/lib/gorm_validator"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email             string `gorm:"unique" validate:"required,uniqueness"`
	EncryptedPassword string `validate:"required"`
	SignedUpThrough   string `validate:"required,oneof=local google"`
}

func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	return gorm_validator.New(tx).Struct(u)
}
