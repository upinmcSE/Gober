package database

import (
	"Gober/internal/database"

	"gorm.io/gorm"
)

func DBMigrator(db *gorm.DB) error {
	return db.AutoMigrate(&database.User{})
}