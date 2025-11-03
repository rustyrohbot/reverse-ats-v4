package handlers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"reverse-ats/internal/exporter"
)

type ExportHandler struct {
	app *pocketbase.PocketBase
}

func NewExportHandler(app *pocketbase.PocketBase) *ExportHandler {
	return &ExportHandler{app: app}
}

func (h *ExportHandler) Export(w http.ResponseWriter, r *http.Request) error {
	// Create a buffer to write our archive to
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Export steps
	steps := []struct {
		name     string
		filename string
		fn       func(*pocketbase.PocketBase, string) error
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
			return err
		}

		// Fetch records
		var records []*core.Record
		switch step.name {
		case "companies":
			records, err = h.app.FindRecordsByFilter("companies", "", "created", -1, 0)
		case "roles":
			records, err = h.app.FindRecordsByFilter("roles", "", "created", -1, 0)
		case "contacts":
			records, err = h.app.FindRecordsByFilter("contacts", "", "created", -1, 0)
		case "interviews", "interviews-contacts":
			records, err = h.app.FindRecordsByFilter("interviews", "", "created", -1, 0)
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to query %s: %v", step.name, err), http.StatusInternalServerError)
			return err
		}

		// Write CSV directly to zip writer
		var writeErr error
		switch step.name {
		case "companies":
			writeErr = exporter.WriteCompaniesCSV(writer, records)
		case "roles":
			writeErr = exporter.WriteRolesCSV(writer, records)
		case "contacts":
			writeErr = exporter.WriteContactsCSV(writer, records)
		case "interviews":
			writeErr = exporter.WriteInterviewsCSV(writer, records)
		case "interviews-contacts":
			writeErr = exporter.WriteInterviewsContactsCSV(writer, records)
		}

		if writeErr != nil {
			http.Error(w, fmt.Sprintf("Failed to write %s: %v", step.name, writeErr), http.StatusInternalServerError)
			return writeErr
		}
	}

	// Close the zip writer
	if err := zipWriter.Close(); err != nil {
		http.Error(w, "Failed to finalize zip", http.StatusInternalServerError)
		return err
	}

	// Set headers for download
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("reverse-ats-export-%s.zip", timestamp)
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", strconv.Itoa(buf.Len()))

	// Write the zip file to the response
	w.Write(buf.Bytes())
	return nil
}
