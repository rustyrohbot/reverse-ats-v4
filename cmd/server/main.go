package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"reverse-ats/internal/database"
	"reverse-ats/internal/db"
	"reverse-ats/internal/handlers"
)

func main() {
	// Database path
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data.db"
	}

	// Connect to database
	dbConn, err := database.New(dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Create queries instance
	queries := db.New(dbConn)

	// Create handlers
	companiesHandler := handlers.NewCompaniesHandler(queries, dbConn)
	rolesHandler := handlers.NewRolesHandler(queries, dbConn)
	interviewsHandler := handlers.NewInterviewsHandler(queries, dbConn)
	contactsHandler := handlers.NewContactsHandler(queries, dbConn)
	statsHandler := handlers.NewStatsHandler(queries, dbConn)
	exportHandler := handlers.NewExportHandler(dbConn)
	importHandler := handlers.NewImportHandler(queries, dbConn)

	// Setup routes
	mux := http.NewServeMux()

	// Static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Page routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/companies", http.StatusFound)
	})

	// Companies routes
	mux.HandleFunc("/companies", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			companiesHandler.Create(w, r)
		} else {
			companiesHandler.List(w, r)
		}
	})
	mux.HandleFunc("/companies/new", companiesHandler.New)
	mux.HandleFunc("/companies/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/edit") {
			companiesHandler.Edit(w, r)
		} else if r.Method == http.MethodPut || r.Method == http.MethodPost {
			companiesHandler.Update(w, r)
		} else if r.Method == http.MethodDelete {
			companiesHandler.Delete(w, r)
		}
	})

	// Roles routes
	mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			rolesHandler.Create(w, r)
		} else {
			rolesHandler.List(w, r)
		}
	})
	mux.HandleFunc("/roles/new", rolesHandler.New)
	mux.HandleFunc("/roles/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/edit") {
			rolesHandler.Edit(w, r)
		} else if r.Method == http.MethodPut || r.Method == http.MethodPost {
			rolesHandler.Update(w, r)
		} else if r.Method == http.MethodDelete {
			rolesHandler.Delete(w, r)
		}
	})

	// Interviews routes
	mux.HandleFunc("/interviews", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			interviewsHandler.Create(w, r)
		} else {
			interviewsHandler.List(w, r)
		}
	})
	mux.HandleFunc("/interviews/new", interviewsHandler.New)
	mux.HandleFunc("/interviews/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/edit") {
			interviewsHandler.Edit(w, r)
		} else if r.Method == http.MethodPut || r.Method == http.MethodPost {
			interviewsHandler.Update(w, r)
		} else if r.Method == http.MethodDelete {
			interviewsHandler.Delete(w, r)
		}
	})

	// Contacts routes
	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			contactsHandler.Create(w, r)
		} else {
			contactsHandler.List(w, r)
		}
	})
	mux.HandleFunc("/contacts/new", contactsHandler.New)
	mux.HandleFunc("/contacts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/edit") {
			contactsHandler.Edit(w, r)
		} else if r.Method == http.MethodPut || r.Method == http.MethodPost {
			contactsHandler.Update(w, r)
		} else if r.Method == http.MethodDelete {
			contactsHandler.Delete(w, r)
		}
	})

	// Stats route
	mux.HandleFunc("/stats", statsHandler.Show)

	// Export route
	mux.HandleFunc("/export", exportHandler.Export)

	// Import route
	mux.HandleFunc("/import", importHandler.Import)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
