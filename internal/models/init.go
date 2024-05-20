package models

import "github.com/troptropcontent/visa_appointment_watcher/database"

func Init() {
	database.DB.AutoMigrate(&Watcher{}, &Credentials{}, &User{})
}
