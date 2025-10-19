package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"reverse-ats/internal/database"
	"reverse-ats/internal/db"
	"reverse-ats/internal/importer"
)

func main() {
	// Parse command line flags
	csvDir := flag.String("dir", ".", "Directory containing CSV files")
	dbPath := flag.String("db", "./data.db", "Path to SQLite database")
	flag.Parse()

	// Check if directory exists
	if _, err := os.Stat(*csvDir); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", *csvDir)
	}

	fmt.Printf("Importing CSV files from: %s\n", *csvDir)
	fmt.Printf("Database: %s\n\n", *dbPath)

	// Connect to database
	dbConn, err := database.New(*dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Create queries instance
	queries := db.New(dbConn)

	// Import all CSV files
	if err := importer.ImportAll(queries, *csvDir); err != nil {
		log.Fatalf("Import failed: %v", err)
	}
}
