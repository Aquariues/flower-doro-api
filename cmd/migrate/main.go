package main

import (
	"log"

	"github.com/Aquariues/flower-doro-api/internal/config"
	"github.com/Aquariues/flower-doro-api/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	log.Println("database migrations applied")
}
