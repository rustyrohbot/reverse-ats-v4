package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"reverse-ats/internal/models"
	"reverse-ats/internal/templates"
)

type InterviewsHandler struct {
	app *pocketbase.PocketBase
}

func NewInterviewsHandler(app *pocketbase.PocketBase) *InterviewsHandler {
	return &InterviewsHandler{app: app}
}

// convertTimeTo24Hour converts 12-hour time format (e.g., "12:00 PM") to 24-hour format (e.g., "12:00")
func convertTimeTo24Hour(timeStr string) string {
	if timeStr == "" {
		return ""
	}

	// Try to parse as 12-hour format first
	formats := []string{
		"3:04 PM",
		"03:04 PM",
		"15:04",    // Already 24-hour
		"3:04",     // Already 24-hour without leading zero
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			// Return in 24-hour format HH:MM
			return t.Format("15:04")
		}
	}

	// If all parsing fails, return original
	return timeStr
}

func recordToInterview(record *core.Record) models.Interview {
	// Extract date properly from DateField
	dateValue := ""
	if dt := record.GetDateTime("date"); !dt.IsZero() {
		// Format as YYYY-MM-DD for HTML date input
		dateValue = dt.Time().Format("2006-01-02")
	}

	// Convert times to 24-hour format for HTML time inputs
	startTime := convertTimeTo24Hour(record.GetString("start"))
	endTime := convertTimeTo24Hour(record.GetString("end"))

	interview := models.Interview{
		ID:        record.Id,
		RoleID:    record.GetString("role"),
		Date:      dateValue,
		Start:     startTime,
		End:       endTime,
		Notes:     record.GetString("notes"),
		Type:      record.GetString("type"),
		CreatedAt: record.GetDateTime("created").String(),
		UpdatedAt: record.GetDateTime("updated").String(),
	}

	// Get role and company info from expanded relation
	if roleRecord := record.ExpandedOne("role"); roleRecord != nil {
		interview.RoleName = roleRecord.GetString("name")
		// Get company from role's company relation
		companyID := roleRecord.GetString("company")
		interview.CompanyID = companyID
	}

	return interview
}

func sortInterviewsByCompanyName(interviews []models.Interview, order string) {
	sort.Slice(interviews, func(i, j int) bool {
		cmpResult := strings.Compare(
			strings.ToLower(interviews[i].CompanyName),
			strings.ToLower(interviews[j].CompanyName),
		)
		if order == "desc" {
			return cmpResult > 0
		}
		return cmpResult < 0
	})
}

func (h *InterviewsHandler) List(w http.ResponseWriter, r *http.Request) error {
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Validate sort field
	validSortFields := map[string]bool{
		"company_name": true,
		"date":         true,
	}

	if sortBy == "" || !validSortFields[sortBy] {
		sortBy = "date"
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	// For company_name sorting, we need to fetch all and sort in memory
	// For other fields, we can use PocketBase's native sorting
	var sortField string
	doInMemorySort := (sortBy == "company_name")

	if !doInMemorySort {
		sortField = sortBy
		if order == "desc" {
			sortField = "-" + sortBy
		}
	} else {
		// No sort for in-memory sorting
		sortField = ""
	}

	// Fetch interviews
	records, err := h.app.FindRecordsByFilter(
		"interviews",
		"",
		sortField,
		-1, // all records
		0,
	)
	if err != nil {
		http.Error(w, "Failed to fetch interviews", http.StatusInternalServerError)
		return err
	}

	// Convert records to Interview structs and fetch role/company info
	interviews := make([]models.Interview, len(records))
	for i, record := range records {
		interview := recordToInterview(record)

		// Fetch role to get role name and company info
		if roleID := record.GetString("role"); roleID != "" {
			if roleRecord, err := h.app.FindRecordById("roles", roleID); err == nil {
				interview.RoleName = roleRecord.GetString("name")
				interview.CompanyID = roleRecord.GetString("company")

				// Fetch company name
				if companyID := roleRecord.GetString("company"); companyID != "" {
					if companyRecord, err := h.app.FindRecordById("companies", companyID); err == nil {
						interview.CompanyName = companyRecord.GetString("name")
					}
				}
			}
		}
		interviews[i] = interview
	}

	// Sort in memory if needed
	if doInMemorySort {
		sortInterviewsByCompanyName(interviews, order)
	}

	// Fetch roles with company info for inline form dropdown
	roleRecords, err := h.app.FindRecordsByFilter("roles", "", "name", -1, 0)
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return err
	}

	// Convert to Role format with company names
	roles := make([]models.Role, len(roleRecords))
	for i, record := range roleRecords {
		role := models.Role{
			ID:        record.Id,
			Name:      record.GetString("name"),
			CompanyID: record.GetString("company"),
		}

		// Fetch company name
		if companyID := record.GetString("company"); companyID != "" {
			if companyRecord, err := h.app.FindRecordById("companies", companyID); err == nil {
				role.CompanyName = companyRecord.GetString("name")
			}
		}
		roles[i] = role
	}

	// Fetch companies for dropdown
	companyRecords, err := h.app.FindRecordsByFilter("companies", "", "name", -1, 0)
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return err
	}

	companies := make([]models.Company, len(companyRecords))
	for i, record := range companyRecords {
		companies[i] = models.Company{
			ID:   record.Id,
			Name: record.GetString("name"),
		}
	}

	return templates.InterviewsList(interviews, sortBy, order, companies).Render(r.Context(), w)
}

func (h *InterviewsHandler) New(w http.ResponseWriter, r *http.Request) error {
	// Fetch roles with company info for dropdown
	roleRecords, err := h.app.FindRecordsByFilter("roles", "", "name", -1, 0)
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return err
	}

	// Convert to RoleWithCompany format
	roles := make([]models.Role, len(roleRecords))
	for i, record := range roleRecords {
		role := models.Role{
			ID:   record.Id,
			Name: record.GetString("name"),
			CompanyID: record.GetString("company"),
		}

		// Fetch company name
		if companyID := record.GetString("company"); companyID != "" {
			if companyRecord, err := h.app.FindRecordById("companies", companyID); err == nil {
				role.CompanyName = companyRecord.GetString("name")
			}
		}
		roles[i] = role
	}

	return templates.InterviewFormNew(roles).Render(r.Context(), w)
}

func (h *InterviewsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return err
	}

	collection, err := h.app.FindCollectionByNameOrId("interviews")
	if err != nil {
		http.Error(w, "Failed to find collection", http.StatusInternalServerError)
		return err
	}

	record := core.NewRecord(collection)
	record.Set("role", r.FormValue("role"))
	record.Set("date", r.FormValue("date"))
	record.Set("start", r.FormValue("start"))
	record.Set("end", r.FormValue("end"))
	record.Set("notes", r.FormValue("notes"))
	record.Set("type", r.FormValue("type"))
	// TODO: Handle contacts many-to-many relationship

	if err := h.app.Save(record); err != nil {
		http.Error(w, "Failed to create interview", http.StatusInternalServerError)
		return err
	}

	// If HTMX request, return just the new row
	if r.Header.Get("HX-Request") == "true" {
		interview := recordToInterview(record)

		// Fetch role to get role name and company info
		if roleID := record.GetString("role"); roleID != "" {
			if roleRecord, err := h.app.FindRecordById("roles", roleID); err == nil {
				interview.RoleName = roleRecord.GetString("name")
				interview.CompanyID = roleRecord.GetString("company")

				// Fetch company name
				if companyID := roleRecord.GetString("company"); companyID != "" {
					if companyRecord, err := h.app.FindRecordById("companies", companyID); err == nil {
						interview.CompanyName = companyRecord.GetString("name")
					}
				}
			}
		}

		return templates.InterviewRow(interview).Render(r.Context(), w)
	}

	// Otherwise redirect
	http.Redirect(w, r, "/interviews", http.StatusSeeOther)
	return nil
}

func (h *InterviewsHandler) Edit(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL path parameter
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return fmt.Errorf("missing id parameter")
	}

	record, err := h.app.FindRecordById("interviews", id)
	if err != nil {
		http.Error(w, "Interview not found", http.StatusNotFound)
		return err
	}

	interview := recordToInterview(record)

	// Fetch role info for display
	if roleID := record.GetString("role"); roleID != "" {
		if roleRecord, err := h.app.FindRecordById("roles", roleID); err == nil {
			interview.RoleName = roleRecord.GetString("name")
		}
	}

	// Fetch roles with company info for dropdown
	roleRecords, err := h.app.FindRecordsByFilter("roles", "", "name", -1, 0)
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return err
	}

	// Convert to Role format with company names
	roles := make([]models.Role, len(roleRecords))
	for i, rec := range roleRecords {
		role := models.Role{
			ID:   rec.Id,
			Name: rec.GetString("name"),
			CompanyID: rec.GetString("company"),
		}

		// Fetch company name
		if companyID := rec.GetString("company"); companyID != "" {
			if companyRecord, err := h.app.FindRecordById("companies", companyID); err == nil {
				role.CompanyName = companyRecord.GetString("name")
			}
		}
		roles[i] = role
	}

	return templates.InterviewFormEdit(interview, roles).Render(r.Context(), w)
}

func (h *InterviewsHandler) Update(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL path parameter
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return fmt.Errorf("missing id parameter")
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return err
	}

	record, err := h.app.FindRecordById("interviews", id)
	if err != nil {
		http.Error(w, "Interview not found", http.StatusNotFound)
		return err
	}

	record.Set("role", r.FormValue("role_id"))
	record.Set("date", r.FormValue("date"))
	record.Set("start", r.FormValue("start"))
	record.Set("end", r.FormValue("end"))
	record.Set("notes", r.FormValue("notes"))
	record.Set("type", r.FormValue("type"))
	// TODO: Handle contacts many-to-many relationship

	if err := h.app.Save(record); err != nil {
		http.Error(w, "Failed to update interview", http.StatusInternalServerError)
		return err
	}

	http.Redirect(w, r, "/interviews", http.StatusSeeOther)
	return nil
}

func (h *InterviewsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL path parameter
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return fmt.Errorf("missing id parameter")
	}

	record, err := h.app.FindRecordById("interviews", id)
	if err != nil {
		http.Error(w, "Interview not found", http.StatusNotFound)
		return err
	}

	if err := h.app.Delete(record); err != nil {
		http.Error(w, "Failed to delete interview", http.StatusInternalServerError)
		return err
	}

	// If HTMX request, return empty response (row will be removed)
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	// Otherwise redirect
	http.Redirect(w, r, "/interviews", http.StatusSeeOther)
	return nil
}

// GetRolesByCompany returns role options for a selected company
func (h *InterviewsHandler) GetRolesByCompany(w http.ResponseWriter, r *http.Request) error {
	companyID := r.URL.Query().Get("company")
	if companyID == "" {
		w.Write([]byte(`<option value="">Select company first</option>`))
		return nil
	}

	// Fetch roles for this company
	filter := fmt.Sprintf("company='%s'", companyID)
	roleRecords, err := h.app.FindRecordsByFilter("roles", filter, "name", -1, 0)
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return err
	}

	// Build HTML options
	html := `<option value="">Select Role *</option>`
	for _, record := range roleRecords {
		html += fmt.Sprintf(`<option value="%s">%s</option>`, record.Id, record.GetString("name"))
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
	return nil
}
