package database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	DBconn *sqlx.DB
)

type SSLMode string

const (
	SSLModeDisable SSLMode = "disable"
)

func ConnectAndMigrate(host, port, databaseName, user, password string, sslMode SSLMode) error {
	connectionSTR := fmt.Sprintf("host=%s port=%s user=%s password=%s  dbname=%s sslmode=%s", host, port, user, password, databaseName, sslMode)
	db, err := sqlx.Open("postgres", connectionSTR)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {

		return err
	}
	DBconn = db
	return migrateUp(db)
}

func ShutDownDB() error {
	return DBconn.Close()
}

func migrateUp(db *sqlx.DB) error {
	driver, driErr := postgres.WithInstance(db.DB, &postgres.Config{})
	if driErr != nil {
		return driErr
	}
	///home/abhinav/Desktop/RMS/database/migrations == database/migrations/
	m, migErr := migrate.NewWithDatabaseInstance(
		"file://database/migrations/", // Path to migration files
		"postgres", driver)            // Database driver and name
	if migErr != nil {
		// Return error if migration instance creation fails
		return migErr
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		// If an error occurs, but it's not "ErrNoChange" (no changes detected), return the error
		return err
	}

	return nil
}

func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := DBconn.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %+v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				logrus.Errorf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logrus.Errorf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}
