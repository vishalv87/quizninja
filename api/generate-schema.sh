#!/bin/bash

# Generate database schema.sql file
# This script connects to the database and generates an up-to-date schema.sql file

echo "Generating database schema.sql..."
go run cmd/generate-schema/main.go

if [ $? -eq 0 ]; then
    echo "✅ Schema.sql generated successfully at database/schema.sql"
else
    echo "❌ Failed to generate schema.sql"
    exit 1
fi