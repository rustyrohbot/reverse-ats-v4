package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"reverse-ats/internal/db"
	"reverse-ats/internal/importer"
)

const (
	maxUploadSize = 10 << 20 // 10 MB
	maxMemory     = 5 << 20  // 5 MB for parsing multipart form
)

type ImportHandler struct {
	queries *db.Queries
	dbConn  *sql.DB
}

func NewImportHandler(queries *db.Queries, dbConn *sql.DB) *ImportHandler {
	return &ImportHandler{queries: queries, dbConn: dbConn}
}

// validateCSVFile performs security checks on uploaded files
func validateCSVFile(filename string, size int64) error {
	// Check file extension
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".csv" {
		return fmt.Errorf("invalid file type: %s (only .csv files allowed)", ext)
	}

	// Check file size
	if size > maxUploadSize {
		return fmt.Errorf("file too large: %d bytes (max %d bytes)", size, maxUploadSize)
	}

	if size == 0 {
		return fmt.Errorf("empty file")
	}

	return nil
}

// saveUploadedFile saves a file to a temporary location with security checks
func saveUploadedFile(file io.Reader, filename string, size int64) (string, error) {
	// Validate file
	if err := validateCSVFile(filename, size); err != nil {
		return "", err
	}

	// Create temp file with secure permissions
	tempFile, err := os.CreateTemp("", "import-*.csv")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// Set secure permissions (owner read/write only)
	if err := tempFile.Chmod(0600); err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to set file permissions: %w", err)
	}

	// Copy with size limit
	written, err := io.CopyN(tempFile, file, maxUploadSize+1)
	if err != nil && err != io.EOF {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	if written > maxUploadSize {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("file exceeds maximum size")
	}

	return tempFile.Name(), nil
}

func (h *ImportHandler) Import(w http.ResponseWriter, r *http.Request) {
	// Limit request size
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize*5) // Allow up to 5 files

	// Parse multipart form
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		http.Error(w, "Failed to parse form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Track temporary files for cleanup
	var tempFiles []string
	defer func() {
		for _, f := range tempFiles {
			os.Remove(f)
		}
	}()

	// Map form field names to step names
	fieldToStep := map[string]string{
		"companies":           "companies",
		"roles":               "roles",
		"contacts":            "contacts",
		"interviews":          "interviews",
		"interviews_contacts": "interviews-contacts",
	}

	// Get import steps in correct order
	steps := importer.GetImportSteps()

	// Process uploaded files and map them to steps
	uploadedSteps := make(map[string]string) // step name -> temp file path
	for fieldName, stepName := range fieldToStep {
		file, header, err := r.FormFile(fieldName)
		if err != nil {
			// File is optional, skip if not provided
			if err == http.ErrMissingFile {
				continue
			}
			http.Error(w, fmt.Sprintf("Failed to read %s file: %v", fieldName, err), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save file to temp location
		tempPath, err := saveUploadedFile(file, header.Filename, header.Size)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to save %s file: %v", fieldName, err), http.StatusBadRequest)
			return
		}

		tempFiles = append(tempFiles, tempPath)
		uploadedSteps[stepName] = tempPath
	}

	// Check if at least one file was uploaded
	if len(uploadedSteps) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	// Set filepaths for uploaded steps
	for i := range steps {
		if tempPath, ok := uploadedSteps[steps[i].Name]; ok {
			steps[i].Filepath = tempPath
		}
	}

	// Import using shared logic (skip missing files)
	errors := importer.ImportFromSteps(h.queries, steps, true)

	// Return response
	if len(errors) > 0 {
		var errorMessages []string
		for _, err := range errors {
			errorMessages = append(errorMessages, err.Error())
		}
		http.Error(w, "Import completed with errors:\n"+strings.Join(errorMessages, "\n"), http.StatusInternalServerError)
		return
	}

	// Redirect to home page on success
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}
