package database

import (
	"database/sql"
	"fmt"
	"log"

	"quizninja-api/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect(cfg *config.Config) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to database")

	// Run migrations
	if err = RunMigrations(DB); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}