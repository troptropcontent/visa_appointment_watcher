package models

import "gorm.io/gorm"

type Credentials struct {
	gorm.Model
	WatcherId uint
	Watcher   Watcher
	Username  string
	Password  string
}
