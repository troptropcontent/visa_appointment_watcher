package config

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	gorm.Model
	Key   string `gorm:"unique"`
	Value string
}

// Init initializes the config database, it returns an error if it fails
func Init() error {
	db, err := gorm.Open(sqlite.Open("./storage/config.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Config{})
	if err != nil {
		return err
	}

	DB = db
	return nil
}

// MustInit initializes the config database, it panics if it fails
func MustInit() {
	err := Init()
	if err != nil {
		panic(err)
	}
}
