package config

import (
	"errors"
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
	database_file := filepath.Join(ROOT_DIR, "storage", "config.db")
	db, err := gorm.Open(sqlite.Open(database_file), &gorm.Config{})
	if err != nil {
		error_msg := "database file : " + database_file + ", Error: " + err.Error()
		return errors.New(error_msg)
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
