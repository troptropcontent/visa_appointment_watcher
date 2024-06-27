package database

import (
	"github.com/troptropcontent/visa_appointment_watcher/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {

	dbPath := config.ROOT_DIR + "/storage/database.db"
	if config.Constants.ENV == "test" {
		dbPath = config.ROOT_DIR + "/storage/test_database.db"
	}
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
}
