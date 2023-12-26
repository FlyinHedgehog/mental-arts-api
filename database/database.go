package database

import (
	"mentalartsapi_hw/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() (*gorm.DB, error) {
	dbURL := "postgres://root:root@localhost:5432/test_db"

	DB, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})

	if err != nil {
		return nil, err
	}

	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	return DB, nil
}
