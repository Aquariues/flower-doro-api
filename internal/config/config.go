package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv           string
	AppName          string
	HTTPAddr         string
	DatabaseURL      string
	AdminDatabaseURL string
	GoAdminPrefix    string
	GoAdminTitle     string
	GoAdminLanguage  string
	AutoMigrate      bool
}

func Load() Config {
	databaseURL := getEnv("DATABASE_URL", "postgres://flowerdoro:flowerdoro@localhost:5432/flowerdoro?sslmode=disable")
	return Config{
		AppEnv:           getEnv("APP_ENV", "development"),
		AppName:          getEnv("APP_NAME", "FlowerDoro API"),
		HTTPAddr:         getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:      databaseURL,
		AdminDatabaseURL: getEnv("ADMIN_DATABASE_URL", databaseURL),
		GoAdminPrefix:    getEnv("GOADMIN_PREFIX", "admin"),
		GoAdminTitle:     getEnv("GOADMIN_TITLE", "FlowerDoro Admin"),
		GoAdminLanguage:  getEnv("GOADMIN_LANGUAGE", "en"),
		AutoMigrate:      getBoolEnv("AUTO_MIGRATE", true),
	}
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getBoolEnv(key string, fallback bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
