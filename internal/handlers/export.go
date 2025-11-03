package handlers

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"reverse-ats/internal/exporter"
)

type ExportHandler struct {
	dbConn *sql.DB
}

func NewExportHandler(dbConn *sql.DB) *ExportHandler {
	return &ExportHandler{dbConn: dbConn}
}

func (h *ExportHandler) Export(w http.ResponseWriter, r *http.Request) {
	// Create a buffer to write our archive to
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Export steps
	steps := []struct {
		name     string
		filename string
		fn       func(*sql.DB, string) error
	}{
		{"companies", "reverse-ats - Companies.csv", exporter.ExportCompanies},
		{"roles", "reverse-ats - Roles.csv", exporter.ExportRoles},
		{"contacts", "reverse-ats - Contacts.csv", exporter.ExportContacts},
		{"interviews", "reverse-ats - Interviews.csv", exporter.ExportInterviews},
		{"interviews-contacts", "reverse-ats - InterviewsContacts.csv", exporter.ExportInterviewsContacts},
	}

	// Export each table to the zip
	for _, step := range steps {
		writer, err := zipWriter.Create(step.filename)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create %s in zip: %v", step.name, err), http.StatusInternalServerError)
			return
		}

		// Query data
		var query string
		switch step.name {
		case "companies":
			records, err = h.app.FindRecordsByFilter("companies", "", "id", -1, 0)
		case "roles":
			records, err = h.app.FindRecordsByFilter("roles", "", "id", -1, 0)
		case "contacts":
			records, err = h.app.FindRecordsByFilter("contacts", "", "id", -1, 0)
		case "interviews", "interviews-contacts":
			records, err = h.app.FindRecordsByFilter("interviews", "", "id", -1, 0)
		}

		rows, err := h.dbConn.Query(query)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to query %s: %v", step.name, err), http.StatusInternalServerError)
			return
		}

		// Write CSV directly to zip writer
		var writeErr error
		switch step.name {
		case "companies":
			writeErr = exporter.WriteCompaniesCSV(writer, rows)
		case "roles":
			writeErr = exporter.WriteRolesCSV(writer, rows)
		case "contacts":
			writeErr = exporter.WriteContactsCSV(writer, rows)
		case "interviews":
			writeErr = exporter.WriteInterviewsCSV(writer, rows)
		case "interviews-contacts":
			writeErr = exporter.WriteInterviewsContactsCSV(writer, rows)
		}
		rows.Close()

		if writeErr != nil {
			http.Error(w, fmt.Sprintf("Failed to write %s: %v", step.name, writeErr), http.StatusInternalServerError)
			return
		}
	}

	// Close the zip writer
	if err := zipWriter.Close(); err != nil {
		http.Error(w, "Failed to finalize zip", http.StatusInternalServerError)
		return
	}

	// Set headers for download
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("reverse-ats-export-%s.zip", timestamp)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

	// Write the zip file to the response
	w.Write(buf.Bytes())
}
