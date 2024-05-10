package config

import "errors"

// SetIfNotExists sets the value of a config key if it does not exist, it returns an error if it fails
func SetIfNotExists(key string, value string) error {
	config := Config{}
	DB.Limit(1).Where("key = ?", key).Find(&config)
	if config.ID == 0 {
		config.Key = key
		config.Value = value
		return DB.Save(&config).Error
	}
	return nil
}

// MustSetIfNotExists sets the value of a config key if it does not exist, it panics if it fails
func MustSetIfNotExists(key string, value string) {
	err := SetIfNotExists(key, value)
	if err != nil {
		panic(err)
	}
}

// Get returns the value of a config key, it returns an error if the key is not found
func Get(key string) (string, error) {
	config := Config{}
	DB.Limit(1).Where("key = ?", key).Find(&config)
	if config.ID == 0 {
		return "", errors.New("config key not found, " + key)
	}
	return config.Value, nil
}

// MustGet returns the value of a config key, it panics if the key is not found
func MustGet(key string) string {
	value, err := Get(key)
	if err != nil {
		panic(err)
	}
	return value
}

// Set sets the value of a config key, it returns an error if it fails
func Set(key string, value string) error {
	config := Config{}
	DB.Limit(1).Where("key = ?", key).Find(&config)
	if config.ID == 0 {
		config.Key = key
		config.Value = value
	}
	config.Value = value
	return DB.Save(&config).Error
}

// MustSet sets the value of a config key, it panics if it fails
func MustSet(key string, value string) {
	err := Set(key, value)
	if err != nil {
		panic(err)
	}
}
