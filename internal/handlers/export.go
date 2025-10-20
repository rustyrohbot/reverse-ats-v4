package handlers

import (
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"reverse-ats/internal/db"
)

type ExportHandler struct {
	queries *db.Queries
	dbConn  *sql.DB
}

func NewExportHandler(queries *db.Queries, dbConn *sql.DB) *ExportHandler {
	return &ExportHandler{queries: queries, dbConn: dbConn}
}

func (h *ExportHandler) Export(w http.ResponseWriter, r *http.Request) {
	// Create a buffer to write our archive to
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	// Export each table
	if err := h.exportCompanies(zipWriter); err != nil {
		http.Error(w, fmt.Sprintf("Failed to export companies: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.exportRoles(zipWriter); err != nil {
		http.Error(w, fmt.Sprintf("Failed to export roles: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.exportContacts(zipWriter); err != nil {
		http.Error(w, fmt.Sprintf("Failed to export contacts: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.exportInterviews(zipWriter); err != nil {
		http.Error(w, fmt.Sprintf("Failed to export interviews: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.exportInterviewsContacts(zipWriter); err != nil {
		http.Error(w, fmt.Sprintf("Failed to export interviews-contacts: %v", err), http.StatusInternalServerError)
		return
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

func (h *ExportHandler) exportCompanies(zipWriter *zip.Writer) error {
	// Query all companies
	query := "SELECT * FROM companies ORDER BY company_id"
	rows, err := h.dbConn.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	writer, err := zipWriter.Create("reverse-ats - Companies.csv")
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(writer)

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

	csvWriter.Flush()
	return csvWriter.Error()
}

func (h *ExportHandler) exportRoles(zipWriter *zip.Writer) error {
	// Query all roles
	query := "SELECT * FROM roles ORDER BY role_id"
	rows, err := h.dbConn.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	writer, err := zipWriter.Create("reverse-ats - Roles.csv")
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(writer)

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
			nullToString(role.Equity),
			nullToString(role.WorkCity),
			nullToString(role.WorkState),
			nullToString(role.Location),
			nullToString(role.Status),
			nullToString(role.Discovery),
			nullToString(role.Referral),
			nullToString(role.Notes),
		})
	}

	csvWriter.Flush()
	return csvWriter.Error()
}

func (h *ExportHandler) exportContacts(zipWriter *zip.Writer) error {
	// Query all contacts
	query := "SELECT * FROM contacts ORDER BY contact_id"
	rows, err := h.dbConn.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	writer, err := zipWriter.Create("reverse-ats - Contacts.csv")
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(writer)

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

	csvWriter.Flush()
	return csvWriter.Error()
}

func (h *ExportHandler) exportInterviews(zipWriter *zip.Writer) error {
	// Query all interviews
	query := "SELECT * FROM interviews ORDER BY interview_id"
	rows, err := h.dbConn.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	writer, err := zipWriter.Create("reverse-ats - Interviews.csv")
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(writer)

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

	csvWriter.Flush()
	return csvWriter.Error()
}

func (h *ExportHandler) exportInterviewsContacts(zipWriter *zip.Writer) error {
	// Query all interview-contact links
	query := "SELECT interviews_contact_id, interview_id, contact_id FROM interviews_contacts ORDER BY interviews_contact_id"
	rows, err := h.dbConn.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	writer, err := zipWriter.Create("reverse-ats - InterviewsContacts.csv")
	if err != nil {
		return err
	}

	csvWriter := csv.NewWriter(writer)

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

	csvWriter.Flush()
	return csvWriter.Error()
}

// Helper functions
func nullToString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func nullInt64ToString(ni sql.NullInt64) string {
	if ni.Valid {
		return strconv.FormatInt(ni.Int64, 10)
	}
	return ""
}
