package models

import "gorm.io/gorm"

type WatcherConfig struct {
	gorm.Model
	AlertPhone string
	VisaType   string
}
