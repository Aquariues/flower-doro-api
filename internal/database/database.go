package database

import (
	_ "embed"

	"github.com/Aquariues/flower-doro-api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed goadmin_bootstrap.sql
var goAdminBootstrapSQL string

func Connect(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pgcrypto`).Error; err != nil {
		return err
	}
	return db.AutoMigrate(
		&models.User{},
		&models.Flower{},
		&models.GardenFlower{},
		&models.FocusSession{},
	)
}

func EnsureGoAdminSchema(db *gorm.DB) error {
	var exists bool
	err := db.Raw(`SELECT EXISTS (
		SELECT 1
		FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = 'goadmin_session'
	)`).Scan(&exists).Error
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return db.Exec(goAdminBootstrapSQL).Error
}
