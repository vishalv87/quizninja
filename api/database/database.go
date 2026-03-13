package database

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	"quizninja-api/config"
	"quizninja-api/utils"

	"github.com/cenkalti/backoff/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/sirupsen/logrus"
)

var DB *sql.DB

func Connect(cfg *config.Config) {
	var dsn string

	if cfg.UseSupabase {
		// Resolve hostname to IPv4 to avoid IPv6 connectivity issues in Docker
		resolvedHost := resolveToIPv4(cfg.SupabaseDBHost)

		// Supabase PostgreSQL connection with SSL required and connection timeout
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require connect_timeout=10",
			resolvedHost, cfg.SupabaseDBPort, cfg.SupabaseDBUser,
			cfg.SupabaseDBPassword, cfg.SupabaseDBName)
		utils.WithFields(logrus.Fields{
			"host":          cfg.SupabaseDBHost,
			"resolved_host": resolvedHost,
			"port":          cfg.SupabaseDBPort,
			"db":            cfg.SupabaseDBName,
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

	// Parse DSN into pgx config and enable simple protocol to avoid
	// prepared statement collisions through PgBouncer/Supabase connection pooler
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		utils.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to parse database DSN")
	}
	connConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	connStr := stdlib.RegisterConnConfig(connConfig)

	// Configure exponential backoff for retry logic
	backoffPolicy := backoff.NewExponentialBackOff()
	backoffPolicy.InitialInterval = 100 * time.Millisecond
	backoffPolicy.MaxInterval = 10 * time.Second
	backoffPolicy.Multiplier = 1.5
	backoffPolicy.RandomizationFactor = 0.1

	var attemptCount int

	// Retry logic for database connection
	operation := func() (struct{}, error) {
		attemptCount++
		utils.WithFields(logrus.Fields{
			"attempt": attemptCount,
		}).Debug("Database connection attempt")

		DB, err = sql.Open("pgx", connStr)
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
	DB.SetConnMaxIdleTime(2 * time.Minute) // Maximum idle time of a connection

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

// resolveToIPv4 resolves a hostname to an IPv4 address only.
// This is useful in Docker environments where IPv6 may not be available.
// If resolution fails or no IPv4 address is found, returns the original hostname.
func resolveToIPv4(hostname string) string {
	// Create a context with timeout for DNS resolution
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use the default resolver with explicit IPv4-only lookup
	resolver := &net.Resolver{
		PreferGo: false, // Use system DNS resolver (more reliable in containers)
	}

	// Try to resolve to IPv4 addresses only using "ip4" network
	ips, err := resolver.LookupIP(ctx, "ip4", hostname)
	if err != nil {
		utils.WithFields(logrus.Fields{
			"hostname": hostname,
			"error":    err.Error(),
		}).Warn("Failed to resolve hostname to IPv4 using ip4 lookup, trying fallback")

		// Fallback: Try standard lookup and filter for IPv4
		ips, err = resolver.LookupIP(ctx, "ip", hostname)
		if err != nil {
			utils.WithFields(logrus.Fields{
				"hostname": hostname,
				"error":    err.Error(),
			}).Warn("Failed to resolve hostname, using original hostname")
			return hostname
		}
	}

	// Find the first IPv4 address
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			utils.WithFields(logrus.Fields{
				"hostname": hostname,
				"ipv4":     ipv4.String(),
			}).Info("Successfully resolved hostname to IPv4 address")
			return ipv4.String()
		}
	}

	// No IPv4 address found, return original hostname
	utils.WithFields(logrus.Fields{
		"hostname": hostname,
		"ip_count": len(ips),
	}).Warn("No IPv4 address found for hostname, using original hostname")
	return hostname
}
