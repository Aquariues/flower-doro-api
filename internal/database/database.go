package database

import (
	_ "embed"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed goadmin_bootstrap.sql
var goAdminBootstrapSQL string

func Connect(databaseURL string) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
}

func AutoMigrate(db *gorm.DB) error {
	return RunMigrations(db)
}

func EnsureGoAdminSchema(db *gorm.DB) error {
	return db.Exec(goAdminBootstrapSQL).Error
}
