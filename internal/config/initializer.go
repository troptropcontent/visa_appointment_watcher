package config

import (
	"path/filepath"
	"runtime"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var ROOT_DIR string

type Config struct {
	gorm.Model
	Key   string `gorm:"unique"`
	Value string
}

func findRootDir() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..")
}

// Init initializes the config database, it returns an error if it fails
func Init() error {
	ROOT_DIR = findRootDir()
	db, err := gorm.Open(sqlite.Open(ROOT_DIR+"/storage/config.db"), &gorm.Config{})
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
