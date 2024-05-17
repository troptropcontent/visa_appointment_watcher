package models

import (
	"time"

	"gorm.io/gorm"
)

type Watcher struct {
	gorm.Model
	User                   User
	UserId                 uint
	Type                   int
	NextAppointmentDate    time.Time
	CurrentAppointmentDate time.Time
	Phone                  string
	IsActive               bool
	LastRunAt              time.Time
}
