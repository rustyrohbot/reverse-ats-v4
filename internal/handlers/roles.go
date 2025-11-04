package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"reverse-ats/internal/models"
	"reverse-ats/internal/templates"
	"reverse-ats/internal/util"
)

type RolesHandler struct {
	app *pocketbase.PocketBase
}

func NewRolesHandler(app *pocketbase.PocketBase) *RolesHandler {
	return &RolesHandler{app: app}
}

func recordToRole(record *core.Record) models.Role {
	role := models.Role{
		ID:                  record.Id,
		CompanyID:           record.GetString("company"),
		Name:                record.GetString("name"),
		Url:                 record.GetString("url"),
		Description:         record.GetString("description"),
		CoverLetter:         record.GetString("cover_letter"),
		ApplicationLocation: record.GetString("application_location"),
		AppliedDate:         record.GetString("applied_date"),
		ClosedDate:          record.GetString("closed_date"),
		PostedRangeMin:      int64(record.GetInt("posted_range_min")),
		PostedRangeMax:      int64(record.GetInt("posted_range_max")),
		Equity:              record.GetBool("equity"),
		WorkCity:            record.GetString("work_city"),
		WorkState:           record.GetString("work_state"),
		Location:            record.GetString("location"),
		Status:              record.GetString("status"),
		Discovery:           record.GetString("discovery"),
		Referral:            record.GetBool("referral"),
		Notes:               record.GetString("notes"),
		CreatedAt:           record.GetDateTime("created").String(),
		UpdatedAt:           record.GetDateTime("updated").String(),
	}

	// Get company name from expanded relation
	if companyRecord := record.ExpandedOne("company"); companyRecord != nil {
		role.CompanyName = companyRecord.GetString("name")
	}

	return role
}

func sortRolesByCompanyName(roles []models.Role, order string) {
	sort.Slice(roles, func(i, j int) bool {
		cmpResult := strings.Compare(
			strings.ToLower(roles[i].CompanyName),
			strings.ToLower(roles[j].CompanyName),
		)
		if order == "desc" {
			return cmpResult > 0
		}
		return cmpResult < 0
	})
}

func (h *RolesHandler) List(w http.ResponseWriter, r *http.Request) error {
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Validate sort field
	validSortFields := map[string]bool{
		"company_name":     true,
		"applied_date":     true,
		"closed_date":      true,
		"posted_range_min": true,
		"posted_range_max": true,
		"equity":           true,
		"work_city":        true,
		"work_state":       true,
		"location":         true,
		"status":           true,
		"referral":         true,
	}

	if sortBy == "" || !validSortFields[sortBy] {
		sortBy = "applied_date"
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

	// Fetch roles
	records, err := h.app.FindRecordsByFilter(
		"roles",
		"",
		sortField,
		-1, // all records
		0,
	)
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return err
	}

	// Fetch all companies once to avoid N+1 queries
	companiesMap, err := util.FetchCompaniesMap(h.app)
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return err
	}

	// Convert records to Role structs with company names
	roles := make([]models.Role, len(records))
	for i, record := range records {
		role := recordToRole(record)
		// Look up company name from map
		if companyID := record.GetString("company"); companyID != "" {
			if companyName, ok := companiesMap[companyID]; ok {
				role.CompanyName = companyName
			}
		}
		roles[i] = role
	}

	// Sort in memory if needed
	if doInMemorySort {
		sortRolesByCompanyName(roles, order)
	}

	// Fetch companies for inline form dropdown
	companies, err := util.FetchCompaniesForDropdown(h.app)
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return err
	}

	return templates.RolesList(roles, sortBy, order, companies).Render(r.Context(), w)
}

func (h *RolesHandler) New(w http.ResponseWriter, r *http.Request) error {
	// Fetch companies for dropdown
	companies, err := util.FetchCompaniesForDropdown(h.app)
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return err
	}

	return templates.RoleFormNew(companies).Render(r.Context(), w)
}

func (h *RolesHandler) Create(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return err
	}

	collection, err := h.app.FindCollectionByNameOrId(util.CollectionRoles)
	if err != nil {
		http.Error(w, "Failed to find collection", http.StatusInternalServerError)
		return err
	}

	record := core.NewRecord(collection)
	record.Set("company", r.FormValue("company"))
	record.Set("name", r.FormValue("name"))
	record.Set("url", r.FormValue("url"))
	record.Set("description", r.FormValue("description"))
	record.Set("cover_letter", r.FormValue("cover_letter"))
	record.Set("application_location", r.FormValue("application_location"))
	record.Set("applied_date", r.FormValue("applied_date"))
	record.Set("closed_date", r.FormValue("closed_date"))

	// Parse numbers
	if minStr := r.FormValue("posted_range_min"); minStr != "" {
		if min, err := strconv.ParseInt(minStr, 10, 64); err == nil {
			record.Set("posted_range_min", min)
		}
	}
	if maxStr := r.FormValue("posted_range_max"); maxStr != "" {
		if max, err := strconv.ParseInt(maxStr, 10, 64); err == nil {
			record.Set("posted_range_max", max)
		}
	}

	record.Set("equity", r.FormValue("equity") == "on" || r.FormValue("equity") == "true")
	record.Set("work_city", r.FormValue("work_city"))
	record.Set("work_state", r.FormValue("work_state"))
	record.Set("location", r.FormValue("location"))
	record.Set("status", r.FormValue("status"))
	record.Set("discovery", r.FormValue("discovery"))
	record.Set("referral", r.FormValue("referral") == "on" || r.FormValue("referral") == "true")
	record.Set("notes", r.FormValue("notes"))

	if err := h.app.Save(record); err != nil {
		http.Error(w, "Failed to create role", http.StatusInternalServerError)
		return err
	}

	// If HTMX request, tell it to reload the page
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/roles")
		w.WriteHeader(http.StatusOK)
		return nil
	}

	http.Redirect(w, r, "/roles", http.StatusSeeOther)
	return nil
}

func (h *RolesHandler) Edit(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL path parameter
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return fmt.Errorf("missing id parameter")
	}

	record, err := h.app.FindRecordById(util.CollectionRoles, id)
	if err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return err
	}

	role := recordToRole(record)

	// Fetch company name for display
	if companyID := record.GetString("company"); companyID != "" {
		if companyRecord, err := h.app.FindRecordById(util.CollectionCompanies, companyID); err == nil {
			role.CompanyName = companyRecord.GetString("name")
		}
	}

	// Fetch companies for dropdown
	companies, err := util.FetchCompaniesForDropdown(h.app)
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return err
	}

	return templates.RoleFormEdit(role, companies).Render(r.Context(), w)
}

func (h *RolesHandler) Update(w http.ResponseWriter, r *http.Request) error {
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

	record, err := h.app.FindRecordById(util.CollectionRoles, id)
	if err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return err
	}

	record.Set("company", r.FormValue("company"))
	record.Set("name", r.FormValue("name"))
	record.Set("url", r.FormValue("url"))
	record.Set("description", r.FormValue("description"))
	record.Set("cover_letter", r.FormValue("cover_letter"))
	record.Set("application_location", r.FormValue("application_location"))
	record.Set("applied_date", r.FormValue("applied_date"))
	record.Set("closed_date", r.FormValue("closed_date"))

	// Parse numbers
	if minStr := r.FormValue("posted_range_min"); minStr != "" {
		if min, err := strconv.ParseInt(minStr, 10, 64); err == nil {
			record.Set("posted_range_min", min)
		}
	} else {
		record.Set("posted_range_min", nil)
	}
	if maxStr := r.FormValue("posted_range_max"); maxStr != "" {
		if max, err := strconv.ParseInt(maxStr, 10, 64); err == nil {
			record.Set("posted_range_max", max)
		}
	} else {
		record.Set("posted_range_max", nil)
	}

	record.Set("equity", r.FormValue("equity") == "on" || r.FormValue("equity") == "true")
	record.Set("work_city", r.FormValue("work_city"))
	record.Set("work_state", r.FormValue("work_state"))
	record.Set("location", r.FormValue("location"))
	record.Set("status", r.FormValue("status"))
	record.Set("discovery", r.FormValue("discovery"))
	record.Set("referral", r.FormValue("referral") == "on" || r.FormValue("referral") == "true")
	record.Set("notes", r.FormValue("notes"))

	if err := h.app.Save(record); err != nil {
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return err
	}

	http.Redirect(w, r, "/roles", http.StatusSeeOther)
	return nil
}

func (h *RolesHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL path parameter
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return fmt.Errorf("missing id parameter")
	}

	record, err := h.app.FindRecordById(util.CollectionRoles, id)
	if err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return err
	}

	if err := h.app.Delete(record); err != nil {
		http.Error(w, "Failed to delete role", http.StatusInternalServerError)
		return err
	}

	// If HTMX request, return empty response (row will be removed)
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	// Otherwise redirect
	http.Redirect(w, r, "/roles", http.StatusSeeOther)
	return nil
}
