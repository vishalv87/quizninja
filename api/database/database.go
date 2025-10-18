package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"quizninja-api/config"
	"quizninja-api/utils"

	"github.com/cenkalti/backoff/v5"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

func Connect(cfg *config.Config) {
	var dsn string

	if cfg.UseSupabase {
		// Supabase PostgreSQL connection with SSL required and connection timeout
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require connect_timeout=10",
			cfg.SupabaseDBHost, cfg.SupabaseDBPort, cfg.SupabaseDBUser,
			cfg.SupabaseDBPassword, cfg.SupabaseDBName)
		utils.WithFields(logrus.Fields{
			"host": cfg.SupabaseDBHost,
			"port": cfg.SupabaseDBPort,
			"db":   cfg.SupabaseDBName,
		}).Info("Connecting to Supabase PostgreSQL database")
	} else {
		// Traditional PostgreSQL connection with timeout
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=10",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
		utils.WithFields(logrus.Fields{
			"host": cfg.DBHost,
			"port": cfg.DBPort,
			"db":   cfg.DBName,
		}).Info("Connecting to traditional PostgreSQL database")
	}

	// Configure exponential backoff for retry logic
	backoffPolicy := backoff.NewExponentialBackOff()
	backoffPolicy.InitialInterval = 100 * time.Millisecond
	backoffPolicy.MaxInterval = 10 * time.Second
	backoffPolicy.Multiplier = 1.5
	backoffPolicy.RandomizationFactor = 0.1

	var err error
	var attemptCount int

	// Retry logic for database connection
	operation := func() (struct{}, error) {
		attemptCount++
		utils.WithFields(logrus.Fields{
			"attempt": attemptCount,
		}).Debug("Database connection attempt")

		DB, err = sql.Open("postgres", dsn)
		if err != nil {
			utils.WithFields(logrus.Fields{
				"attempt": attemptCount,
				"error":   err.Error(),
			}).Warn("Database connection attempt failed - error opening connection")
			return struct{}{}, err
		}

		// Ping with context to ensure connection is actually established
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err = DB.PingContext(ctx); err != nil {
			utils.WithFields(logrus.Fields{
				"attempt": attemptCount,
				"error":   err.Error(),
			}).Warn("Database connection attempt failed - error pinging database")
			DB.Close() // Clean up failed connection
			return struct{}{}, err
		}

		return struct{}{}, nil
	}

	// Execute with retry and backoff
	_, err = backoff.Retry(
		context.Background(),
		operation,
		backoff.WithBackOff(backoffPolicy),
		backoff.WithMaxElapsedTime(30*time.Second),
	)
	if err != nil {
		utils.WithFields(logrus.Fields{
			"attempts": attemptCount,
			"error":    err.Error(),
		}).Fatal("Failed to connect to database after multiple attempts")
	}

	// Configure connection pool settings
	DB.SetMaxOpenConns(25)                 // Maximum number of open connections
	DB.SetMaxIdleConns(5)                  // Maximum number of idle connections
	DB.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	if cfg.UseSupabase {
		utils.WithFields(logrus.Fields{
			"attempts":   attemptCount,
			"database":   "supabase",
			"max_conns":  25,
			"idle_conns": 5,
		}).Info("Successfully connected to Supabase database")
	} else {
		utils.WithFields(logrus.Fields{
			"attempts":   attemptCount,
			"database":   "traditional",
			"max_conns":  25,
			"idle_conns": 5,
		}).Info("Successfully connected to traditional database")
	}

	// Run migrations - works for both traditional PostgreSQL and Supabase
	if err = RunMigrations(DB); err != nil {
		utils.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to run migrations")
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
