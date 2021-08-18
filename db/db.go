package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	models "timage.flomas.net/model"
)

func ConnectDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("local.db"), &gorm.Config{
		//DisableAutomaticPing: true,
		//DisableForeignKeyConstraintWhenMigrating: true,
	})

	if err != nil {
		return nil, err
	}

	//if err := db.SetupJoinTable(&models.User{}, "Images", &UserImages{}); err != nil {
	//	return nil, err
	//}

	if err := db.AutoMigrate(&models.Image{}, &models.User{}, &models.ImageStore{}); err != nil {
		return nil, err
	}

	return db, nil
}