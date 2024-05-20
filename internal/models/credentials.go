package models

import "gorm.io/gorm"

type Credentials struct {
	gorm.Model
	Username string
	Password string
}
