package main

import (
	"log"

	"quizninja-api/config"
	"quizninja-api/database"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	database.Connect(cfg)
	defer database.Close()

	// Generate schema.sql
	schemaPath := "database/schema.sql"
	if err := database.GenerateSchemaSQL(database.DB, schemaPath); err != nil {
		log.Fatalf("Failed to generate schema.sql: %v", err)
	}

	log.Printf("Schema.sql generated successfully at: %s", schemaPath)
}
