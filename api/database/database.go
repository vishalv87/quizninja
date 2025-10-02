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
	var dsn string

	if cfg.UseSupabase {
		// Supabase PostgreSQL connection with SSL required
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			cfg.SupabaseDBHost, cfg.SupabaseDBPort, cfg.SupabaseDBUser,
			cfg.SupabaseDBPassword, cfg.SupabaseDBName)
		log.Println("Connecting to Supabase PostgreSQL database...")
	} else {
		// Traditional PostgreSQL connection
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
		log.Println("Connecting to traditional PostgreSQL database...")
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	if cfg.UseSupabase {
		log.Println("Successfully connected to Supabase database")
	} else {
		log.Println("Successfully connected to traditional database")
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