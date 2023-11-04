package main

import (
	"golang.org/x/exp/slog"
	"os"
	"rest-api/internal/config"
	"rest-api/internal/lib/slogger/sl"
	"rest-api/internal/storage/mysql"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// config
	cfg := config.MustLoad()

	// log
	log := setupLogger(cfg.Env)
	log.Info("starting app", slog.String("env", cfg.Env))

	// db connect
	storage, err := mysql.New(cfg.AddressDB, cfg.Login, cfg.Pass, cfg.NameDB)
	if err != nil {
		log.Error("failed to init storage:", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
