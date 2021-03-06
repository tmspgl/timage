package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"timage.flomas.net/model"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("local.db"), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Image{}); err != nil {
		return nil, err
	}

	return db, nil
}