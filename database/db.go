package database

import (
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	db, err := gorm.Open(sqlite.Open(config.ROOT_DIR+"/storage/database.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
}
