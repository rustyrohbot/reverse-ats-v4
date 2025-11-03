package exporter

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"reverse-ats/internal/db"
)

// ExportAll exports all tables to CSV files in the specified directory
func ExportAll(dbConn *sql.DB, outputDir string) error {
	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Export each table
	steps := []struct {
		name     string
		filename string
		fn       func(*sql.DB, string) error
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
		if err := step.fn(dbConn, filepath); err != nil {
			return fmt.Errorf("failed to export %s: %w", step.name, err)
		}
	}

	fmt.Println("\n✅ All data exported successfully!")
	return nil
}

// ExportCompanies exports companies table to CSV
func ExportCompanies(dbConn *sql.DB, filepath string) error {
	query := "SELECT * FROM companies ORDER BY company_id"
	rows, err := dbConn.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return WriteCompaniesCSV(file, rows)
}

// ExportRoles exports roles table to CSV
func ExportRoles(dbConn *sql.DB, filepath string) error {
	query := "SELECT * FROM roles ORDER BY role_id"
	rows, err := dbConn.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return WriteRolesCSV(file, rows)
}

// ExportContacts exports contacts table to CSV
func ExportContacts(dbConn *sql.DB, filepath string) error {
	query := "SELECT * FROM contacts ORDER BY contact_id"
	rows, err := dbConn.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return WriteContactsCSV(file, rows)
}

// ExportInterviews exports interviews table to CSV
func ExportInterviews(dbConn *sql.DB, filepath string) error {
	query := "SELECT * FROM interviews ORDER BY interview_id"
	rows, err := dbConn.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return WriteInterviewsCSV(file, rows)
}

// ExportInterviewsContacts exports interviews_contacts junction table to CSV
func ExportInterviewsContacts(dbConn *sql.DB, filepath string) error {
	query := "SELECT interviews_contact_id, interview_id, contact_id FROM interviews_contacts ORDER BY interviews_contact_id"
	rows, err := dbConn.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	return WriteInterviewsContactsCSV(file, rows)
}

// WriteCompaniesCSV writes companies data to CSV writer
func WriteCompaniesCSV(writer io.Writer, rows *sql.Rows) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	csvWriter.Write([]string{"companyID", "name", "description", "url", "linkedin", "hqCity", "hqState"})

	// Write data
	for rows.Next() {
		var company db.Company
		err := rows.Scan(
			&company.CompanyID,
			&company.Name,
			&company.Description,
			&company.Url,
			&company.Linkedin,
			&company.HqCity,
			&company.HqState,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			return err
		}

		csvWriter.Write([]string{
			strconv.FormatInt(company.CompanyID, 10),
			company.Name,
			nullToString(company.Description),
			nullToString(company.Url),
			nullToString(company.Linkedin),
			nullToString(company.HqCity),
			nullToString(company.HqState),
		})
	}

	return csvWriter.Error()
}

// WriteRolesCSV writes roles data to CSV writer
func WriteRolesCSV(writer io.Writer, rows *sql.Rows) error {
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
	for rows.Next() {
		var role db.Role
		err := rows.Scan(
			&role.RoleID,
			&role.CompanyID,
			&role.Name,
			&role.Url,
			&role.Description,
			&role.CoverLetter,
			&role.ApplicationLocation,
			&role.AppliedDate,
			&role.ClosedDate,
			&role.PostedRangeMin,
			&role.PostedRangeMax,
			&role.Equity,
			&role.WorkCity,
			&role.WorkState,
			&role.Location,
			&role.Status,
			&role.Discovery,
			&role.Referral,
			&role.Notes,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return err
		}

		csvWriter.Write([]string{
			strconv.FormatInt(role.RoleID, 10),
			strconv.FormatInt(role.CompanyID, 10),
			role.Name,
			nullToString(role.Url),
			nullToString(role.Description),
			nullToString(role.CoverLetter),
			nullToString(role.ApplicationLocation),
			nullToString(role.AppliedDate),
			nullToString(role.ClosedDate),
			nullInt64ToString(role.PostedRangeMin),
			nullInt64ToString(role.PostedRangeMax),
			nullBoolToString(role.Equity),
			nullToString(role.WorkCity),
			nullToString(role.WorkState),
			nullToString(role.Location),
			nullToString(role.Status),
			nullToString(role.Discovery),
			nullBoolToString(role.Referral),
			nullToString(role.Notes),
		})
	}

	return csvWriter.Error()
}

// WriteContactsCSV writes contacts data to CSV writer
func WriteContactsCSV(writer io.Writer, rows *sql.Rows) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	csvWriter.Write([]string{
		"contactID", "companyID", "firstName", "lastName", "role",
		"email", "phone", "linkedin", "notes",
	})

	// Write data
	for rows.Next() {
		var contact db.Contact
		err := rows.Scan(
			&contact.ContactID,
			&contact.CompanyID,
			&contact.FirstName,
			&contact.LastName,
			&contact.Role,
			&contact.Email,
			&contact.Phone,
			&contact.Linkedin,
			&contact.Notes,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		)
		if err != nil {
			return err
		}

		csvWriter.Write([]string{
			strconv.FormatInt(contact.ContactID, 10),
			strconv.FormatInt(contact.CompanyID, 10),
			contact.FirstName,
			contact.LastName,
			nullToString(contact.Role),
			nullToString(contact.Email),
			nullToString(contact.Phone),
			nullToString(contact.Linkedin),
			nullToString(contact.Notes),
		})
	}

	return csvWriter.Error()
}

// WriteInterviewsCSV writes interviews data to CSV writer
func WriteInterviewsCSV(writer io.Writer, rows *sql.Rows) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	csvWriter.Write([]string{
		"interviewID", "roleID", "date", "start", "end", "notes", "type",
	})

	// Write data
	for rows.Next() {
		var interview db.Interview
		err := rows.Scan(
			&interview.InterviewID,
			&interview.RoleID,
			&interview.Date,
			&interview.Start,
			&interview.End,
			&interview.Notes,
			&interview.Type,
			&interview.CreatedAt,
			&interview.UpdatedAt,
		)
		if err != nil {
			return err
		}

		csvWriter.Write([]string{
			strconv.FormatInt(interview.InterviewID, 10),
			strconv.FormatInt(interview.RoleID, 10),
			interview.Date,
			interview.Start,
			interview.End,
			nullToString(interview.Notes),
			interview.Type,
		})
	}

	return csvWriter.Error()
}

// WriteInterviewsContactsCSV writes interviews_contacts data to CSV writer
func WriteInterviewsContactsCSV(writer io.Writer, rows *sql.Rows) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Write header
	csvWriter.Write([]string{
		"interviewsContactId", "interviewId", "contactId",
	})

	// Write data
	for rows.Next() {
		var id, interviewID, contactID int64
		err := rows.Scan(&id, &interviewID, &contactID)
		if err != nil {
			return err
		}

		csvWriter.Write([]string{
			strconv.FormatInt(id, 10),
			strconv.FormatInt(interviewID, 10),
			strconv.FormatInt(contactID, 10),
		})
	}

	return csvWriter.Error()
}

// Helper functions
func nullToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return "NULL"
}

func nullInt64ToString(ni sql.NullInt64) string {
	if ni.Valid {
		return strconv.FormatInt(ni.Int64, 10)
	}
	return "NULL"
}

func nullBoolToString(nb sql.NullBool) string {
	if nb.Valid {
		if nb.Bool {
			return "true"
		}
		return "false"
	}
	return "NULL"
}
