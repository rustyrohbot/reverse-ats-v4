package main

import (
	"fmt"
	"log"

	"github.com/pocketbase/pocketbase"

	"reverse-ats/internal/exporter"
)

func main() {
	// Fixed paths
	exportDir := "./export"
	dbPath := "./pb_data"

	fmt.Printf("Exporting data to: %s\n", exportDir)
	fmt.Printf("PocketBase data directory: %s\n\n", dbPath)

	// Initialize PocketBase
	app := pocketbase.NewWithConfig(pocketbase.Config{
		DefaultDataDir: dbPath,
	})

	// Bootstrap PocketBase (loads collections schema)
	if err := app.Bootstrap(); err != nil {
		log.Fatalf("Failed to bootstrap PocketBase: %v", err)
	}

	// Export all tables
	if err := exporter.ExportAll(app, exportDir); err != nil {
		log.Fatalf("Export failed: %v", err)
	}
}
