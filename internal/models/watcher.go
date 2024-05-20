package models

import (
	"time"

	"gorm.io/gorm"
)

type Watcher struct {
	gorm.Model
	User                   User
	UserId                 uint
	NextAppointmentDate    time.Time
	CurrentAppointmentDate time.Time
	IsActive               bool
	LastRunAt              time.Time
	WatcherConfigId        uint
	WatcherConfig          WatcherConfig
	CredentialsId          uint
	Credentials            Credentials
}
