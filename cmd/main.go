package main

import (
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/sudo-abhinav/rms/database"
	"github.com/sudo-abhinav/rms/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutDownTime = 10 * time.Second

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	//fmt.Println(os.Getenv("DB_HOST"), "hello")
	server := routes.SetupRoutes()
	fmt.Printf(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"))

	if err := database.ConnectAndMigrate(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error = %+v", err)
	}
	logrus.Print("migration successfully..")

	go func() {
		if err := server.RUN(":3000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Panicf("Failed to run server with error: %+v", err)
		}
	}()
	logrus.Print("Server started at :3000")
	<-done

	logrus.Info("shutting down server")

	if err := database.ShutDownDB(); err != nil {
		logrus.WithError(err).Error("failed to close database connection")
	}

	if err := server.Shutdown(shutDownTime); err != nil {
		logrus.WithError(err).Panic("failed to gracefully shutdown server")
	}
}
