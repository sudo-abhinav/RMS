package main

import (
	"errors"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/sudo-abhinav/rms/database"
	"github.com/sudo-abhinav/rms/routes"
	"log/slog"
	_ "log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutDownTime = 10 * time.Second

func main() {
	//TODO :-  implementing logger
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	logger.Info("JSONHandler Example", "Content", "Logging in JSON format")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	server := routes.SetupRoutes()
	//fmt.Printf(
	//	os.Getenv("DB_HOST"),
	//	os.Getenv("DB_PORT"),
	//	os.Getenv("DB_NAME"),
	//	os.Getenv("DB_USER"),
	//	os.Getenv("DB_PASS"))

	if err := database.ConnectAndMigrate(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		database.SSLModeDisable); err != nil {
		logger.Error("Failed to initialize and migrate database with error = %+v", err)
	}
	logger.Info("migration successfully..")

	go func() {
		if err := server.RUN(":3000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Warn("Failed to run server with error: %+v", err)
		}
	}()
	logger.Info("Server started at :3000")
	<-done

	logger.Info("Server Shutdown", "Content", "shutting down server")
	//logrus.Info("shutting down server")

	if err := database.ShutDownDB(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}

	if err := server.Shutdown(shutDownTime); err != nil {
		logger.With(err).Warn("failed to gracefully shutdown server")
	}
}
