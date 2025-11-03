package exporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// ExportAll exports all tables to CSV files in the specified directory
func ExportAll(app *pocketbase.PocketBase, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Export each table
	steps := []struct {
		name     string
		filename string
		fn       func(*pocketbase.PocketBase, string) error
	}{
		{"companies", "reverse-ats - Companies.csv", ExportCompanies},
		{"roles", "reverse-ats - Roles.csv", ExportRoles},
		{"contacts", "reverse-ats - Contacts.csv", ExportContacts},
		{"interviews", "reverse-ats - Interviews.csv", ExportInterviews},
		{"interviews-contacts", "reverse-ats - InterviewsContacts.csv", ExportInterviewsContacts},
	}

	for _, step := range steps {
		filepath := outputDir + "/" + step.filename
		fmt.Printf("Exporting %s to %s...\n", step.name, filepath)
		if err := step.fn(app, filepath); err != nil {
			return fmt.Errorf("failed to export %s: %w", step.name, err)
		}
	}

	fmt.Println("\n✅ All data exported successfully!")
	return nil
}

// ExportCompanies exports companies table to CSV
func ExportCompanies(app *pocketbase.PocketBase, filepath string) error {
	records, err := app.FindRecordsByFilter("companies", "", "created", -1, 0)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return WriteCompaniesCSV(file, records)
}

// ExportRoles exports roles table to CSV
func ExportRoles(app *pocketbase.PocketBase, filepath string) error {
	records, err := app.FindRecordsByFilter("roles", "", "created", -1, 0)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return WriteRolesCSV(file, records)
}

// ExportContacts exports contacts table to CSV
func ExportContacts(app *pocketbase.PocketBase, filepath string) error {
	records, err := app.FindRecordsByFilter("contacts", "", "created", -1, 0)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return WriteContactsCSV(file, records)
}

// ExportInterviews exports interviews table to CSV
func ExportInterviews(app *pocketbase.PocketBase, filepath string) error {
	records, err := app.FindRecordsByFilter("interviews", "", "created", -1, 0)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return WriteInterviewsCSV(file, records)
}

// ExportInterviewsContacts exports interview-contact relationships to CSV
func ExportInterviewsContacts(app *pocketbase.PocketBase, filepath string) error {
	// Fetch all interviews with contacts field
	records, err := app.FindRecordsByFilter("interviews", "", "created", -1, 0)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return WriteInterviewsContactsCSV(file, records)
}

// WriteCompaniesCSV writes companies data to CSV writer
func WriteCompaniesCSV(writer io.Writer, records []*core.Record) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	csvWriter.Write([]string{"companyID", "name", "description", "url", "linkedin", "hqCity", "hqState"})

	// Write data
	for _, record := range records {
		csvWriter.Write([]string{
			record.Id,
			record.GetString("name"),
			emptyToNull(record.GetString("description")),
			emptyToNull(record.GetString("url")),
			emptyToNull(record.GetString("linkedin")),
			emptyToNull(record.GetString("hq_city")),
			emptyToNull(record.GetString("hq_state")),
		})
	}

	return csvWriter.Error()
}

// WriteRolesCSV writes roles data to CSV writer
func WriteRolesCSV(writer io.Writer, records []*core.Record) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	csvWriter.Write([]string{
		"roleID", "companyID", "name", "url", "description", "coverLetter",
		"applicationLocation", "appliedDate", "closedDate", "postedRangeMin",
		"postedRangeMax", "equity", "workCity", "workState", "location",
		"status", "discovery", "referral", "notes",
	})

	// Write data
	for _, record := range records {
		csvWriter.Write([]string{
			record.Id,
			record.GetString("company"),
			record.GetString("name"),
			emptyToNull(record.GetString("url")),
			emptyToNull(record.GetString("description")),
			emptyToNull(record.GetString("cover_letter")),
			emptyToNull(record.GetString("application_location")),
			emptyToNull(record.GetString("applied_date")),
			emptyToNull(record.GetString("closed_date")),
			int64ToString(record.GetInt("posted_range_min")),
			int64ToString(record.GetInt("posted_range_max")),
			boolToString(record.GetBool("equity")),
			emptyToNull(record.GetString("work_city")),
			emptyToNull(record.GetString("work_state")),
			emptyToNull(record.GetString("location")),
			emptyToNull(record.GetString("status")),
			emptyToNull(record.GetString("discovery")),
			boolToString(record.GetBool("referral")),
			emptyToNull(record.GetString("notes")),
		})
	}

	return csvWriter.Error()
}

// WriteContactsCSV writes contacts data to CSV writer
func WriteContactsCSV(writer io.Writer, records []*core.Record) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	csvWriter.Write([]string{
		"contactID", "companyID", "firstName", "lastName", "role",
		"email", "phone", "linkedin", "notes",
	})

	// Write data
	for _, record := range records {
		csvWriter.Write([]string{
			record.Id,
			record.GetString("company"),
			record.GetString("first_name"),
			record.GetString("last_name"),
			emptyToNull(record.GetString("role")),
			emptyToNull(record.GetString("email")),
			emptyToNull(record.GetString("phone")),
			emptyToNull(record.GetString("linkedin")),
			emptyToNull(record.GetString("notes")),
		})
	}

	return csvWriter.Error()
}

// WriteInterviewsCSV writes interviews data to CSV writer
func WriteInterviewsCSV(writer io.Writer, records []*core.Record) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	csvWriter.Write([]string{
		"interviewID", "roleID", "date", "start", "end", "notes", "type",
	})

	// Write data
	for _, record := range records {
		csvWriter.Write([]string{
			record.Id,
			record.GetString("role"),
			record.GetString("date"),
			record.GetString("start"),
			record.GetString("end"),
			emptyToNull(record.GetString("notes")),
			record.GetString("type"),
		})
	}

	return csvWriter.Error()
}

// WriteInterviewsContactsCSV writes interview-contact relationships to CSV writer
func WriteInterviewsContactsCSV(writer io.Writer, records []*core.Record) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	csvWriter.Write([]string{
		"interviewsContactId", "interviewId", "contactId",
	})

	// Write data - flatten many-to-many relationships
	linkID := 1
	for _, record := range records {
		interviewID := record.Id
		contacts := record.GetStringSlice("contacts")

		for _, contactID := range contacts {
			csvWriter.Write([]string{
				strconv.Itoa(linkID),
				interviewID,
				contactID,
			})
			linkID++
		}
	}

	return csvWriter.Error()
}

// Helper functions
func emptyToNull(s string) string {
	if s == "" {
		return "NULL"
	}
	return s
}

func int64ToString(val int) string {
	if val == 0 {
		return "NULL"
	}
	return strconv.Itoa(val)
}

func boolToString(val bool) string {
	if val {
		return "true"
	}
	return "false"
}
