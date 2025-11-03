package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pocketbase/pocketbase"

	"reverse-ats/internal/importer"
)

func main() {
	// Fixed paths
	csvDir := "./import"
	dbPath := "./pb_data"

	// Check if directory exists
	if _, err := os.Stat(csvDir); os.IsNotExist(err) {
		log.Fatalf("Directory does not exist: %s", csvDir)
	}

	fmt.Printf("Importing CSV files from: %s\n", csvDir)
	fmt.Printf("PocketBase data directory: %s\n\n", dbPath)

	// Initialize PocketBase
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: dbPath,
	})

	// Bootstrap PocketBase (loads collections schema)
	if err := app.Bootstrap(); err != nil {
		log.Fatalf("Failed to bootstrap PocketBase: %v", err)
	}

	// Import all CSV files
	if err := importer.ImportAll(app, csvDir); err != nil {
		log.Fatalf("Import failed: %v", err)
	}
}
