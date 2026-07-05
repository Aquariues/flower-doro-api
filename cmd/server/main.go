package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Aquariues/flower-doro-api/internal/admin"
	"github.com/Aquariues/flower-doro-api/internal/config"
	"github.com/Aquariues/flower-doro-api/internal/database"
	"github.com/Aquariues/flower-doro-api/internal/httpapi"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}

	if cfg.AutoMigrate {
		if err := database.AutoMigrate(db); err != nil {
			log.Fatalf("auto migrate: %v", err)
		}
	}
	if err := database.EnsureGoAdminSchema(db); err != nil {
		log.Fatalf("ensure GoAdmin schema: %v", err)
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.GET("/admin", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/login")
	})
	router.HEAD("/admin", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/admin/login")
	})
	httpapi.Register(router, db)

	if err := admin.Register(router, cfg); err != nil {
		log.Fatalf("register admin: %v", err)
	}

	server := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("FlowerDoro API listening on %s", cfg.HTTPAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown: %v", err)
	}
}
