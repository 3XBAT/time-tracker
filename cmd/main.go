package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/3XBAT/time-tracker/internal/config"
	"github.com/3XBAT/time-tracker/internal/handlers"
	"github.com/3XBAT/time-tracker/internal/service"
	"github.com/3XBAT/time-tracker/internal/storage"
	"github.com/3XBAT/time-tracker/server"
	"github.com/golang-migrate/migrate/v4"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	// Драйвер для выполнения миграций
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	// Драйвер для получения миграций из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title Time Tracker API
// @version 1.01
// @description API Server for TimeTracker Application

// @host localhost:8080
// @BasePath /

func main() {
	log := setupLogger(envLocal)

	log.Info("starting application")

	cfg := config.MustLoad()

	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.DB.Username,
			cfg.DB.Password,
			cfg.DB.Host,
			cfg.DB.Port,
			cfg.DB.DBName,
			cfg.DB.SSLMode),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")
		}
	}

	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		log.Warn(err.Error())
	}

	dataStorage := storage.NewStorage(db)
	services := service.NewService(log, dataStorage)
	handler := handlers.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(cfg.Port, handler.InitRoutes()); err != nil {
			log.Error("error occurred while running the server %s", err.Error())
		}
	}()

	log.Info("server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.ShutDown(context.Background()); err != nil {
		log.Error("error occurred while shutting down server: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		log.Error("error occurred while closing db: %s", err.Error())
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {

	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
