package main

import (
	"fmt"
	"log"

	"reverse-ats/internal/database"
	"reverse-ats/internal/exporter"
)

func main() {
	// Fixed paths
	exportDir := "./export"
	dbPath := "./data.db"

	fmt.Printf("Exporting data to: %s\n", exportDir)
	fmt.Printf("Database: %s\n\n", dbPath)

	// Connect to database
	dbConn, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Export all tables
	if err := exporter.ExportAll(dbConn, exportDir); err != nil {
		log.Fatalf("Export failed: %v", err)
	}
}
