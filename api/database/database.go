package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"quizninja-api/config"

	"github.com/cenkalti/backoff/v5"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(cfg *config.Config) {
	var dsn string

	if cfg.UseSupabase {
		// Supabase PostgreSQL connection with SSL required and connection timeout
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require connect_timeout=10",
			cfg.SupabaseDBHost, cfg.SupabaseDBPort, cfg.SupabaseDBUser,
			cfg.SupabaseDBPassword, cfg.SupabaseDBName)
		log.Println("Connecting to Supabase PostgreSQL database...")
	} else {
		// Traditional PostgreSQL connection with timeout
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable connect_timeout=10",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
		log.Println("Connecting to traditional PostgreSQL database...")
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
		log.Printf("Database connection attempt %d...", attemptCount)

		DB, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Attempt %d failed - error opening connection: %v", attemptCount, err)
			return struct{}{}, err
		}

		// Ping with context to ensure connection is actually established
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err = DB.PingContext(ctx); err != nil {
			log.Printf("Attempt %d failed - error pinging database: %v", attemptCount, err)
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
		log.Fatalf("Failed to connect to database after %d attempts: %v", attemptCount, err)
	}

	// Configure connection pool settings
	DB.SetMaxOpenConns(25)                 // Maximum number of open connections
	DB.SetMaxIdleConns(5)                  // Maximum number of idle connections
	DB.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection

	if cfg.UseSupabase {
		log.Printf("Successfully connected to Supabase database after %d attempt(s)", attemptCount)
	} else {
		log.Printf("Successfully connected to traditional database after %d attempt(s)", attemptCount)
	}

	// Run migrations - works for both traditional PostgreSQL and Supabase
	if err = RunMigrations(DB); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
