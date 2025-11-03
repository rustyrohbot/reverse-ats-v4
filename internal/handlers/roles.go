package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"reverse-ats/internal/db"
	"reverse-ats/internal/templates"
)

type RolesHandler struct {
	queries *db.Queries
	dbConn  *sql.DB
}

func NewRolesHandler(queries *db.Queries, dbConn *sql.DB) *RolesHandler {
	return &RolesHandler{queries: queries, dbConn: dbConn}
}

func (h *RolesHandler) List(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Map display names to actual column/field names
	sortColumnMap := map[string]string{
		"company_name":      "c.name",
		"applied_date":      "r.applied_date",
		"closed_date":       "r.closed_date",
		"posted_range_min":  "r.posted_range_min",
		"posted_range_max":  "r.posted_range_max",
		"equity":            "r.equity",
		"work_city":         "r.work_city",
		"work_state":        "r.work_state",
		"location":          "r.location",
		"status":            "r.status",
		"referral":          "r.referral",
	}

	sortCol, ok := sortColumnMap[sortBy]
	if !ok {
		sortBy = "applied_date"
		sortCol = "r.applied_date"
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	query := fmt.Sprintf(`
		SELECT
			r.role_id,
			r.company_id,
			r.name as role_name,
			r.url,
			r.description,
			r.cover_letter,
			r.application_location,
			r.applied_date,
			r.closed_date,
			r.posted_range_min,
			r.posted_range_max,
			r.equity,
			r.work_city,
			r.work_state,
			r.location,
			r.status,
			r.discovery,
			r.referral,
			r.notes,
			r.created_at,
			r.updated_at,
			c.name as company_name
		FROM roles r
		INNER JOIN companies c ON r.company_id = c.company_id
		ORDER BY %s %s`, sortCol, strings.ToUpper(order))

	rows, err := h.dbConn.QueryContext(r.Context(), query)
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var roles []db.ListRolesWithCompanyRow
	for rows.Next() {
		var role db.ListRolesWithCompanyRow
		err := rows.Scan(
			&role.RoleID,
			&role.CompanyID,
			&role.RoleName,
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
			&role.CompanyName,
		)
		if err != nil {
			http.Error(w, "Failed to scan roles", http.StatusInternalServerError)
			return
		}
		roles = append(roles, role)
	}

	templates.RolesList(roles, sortBy, order).Render(r.Context(), w)
}

func (h *RolesHandler) New(w http.ResponseWriter, r *http.Request) {
	companies, err := h.queries.ListCompanies(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return
	}

	templates.RoleFormNew(companies).Render(r.Context(), w)
}

func (h *RolesHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	companyID, err := strconv.ParseInt(r.FormValue("company_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	_, err = h.queries.CreateRole(r.Context(), db.CreateRoleParams{
		CompanyID:           companyID,
		Name:                r.FormValue("name"),
		Url:                 nullString(r.FormValue("url")),
		Description:         nullString(r.FormValue("description")),
		CoverLetter:         nullString(r.FormValue("cover_letter")),
		ApplicationLocation: nullString(r.FormValue("application_location")),
		AppliedDate:         nullString(r.FormValue("applied_date")),
		ClosedDate:          nullString(r.FormValue("closed_date")),
		PostedRangeMin:      nullInt64(r.FormValue("posted_range_min")),
		PostedRangeMax:      nullInt64(r.FormValue("posted_range_max")),
		Equity:              nullBool(r.FormValue("equity")),
		WorkCity:            nullString(r.FormValue("work_city")),
		WorkState:           nullString(r.FormValue("work_state")),
		Location:            nullString(r.FormValue("location")),
		Status:              nullString(r.FormValue("status")),
		Discovery:           nullString(r.FormValue("discovery")),
		Referral:            nullBool(r.FormValue("referral")),
		Notes:               nullString(r.FormValue("notes")),
	})
	if err != nil {
		http.Error(w, "Failed to create role", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/roles", http.StatusSeeOther)
}

func (h *RolesHandler) Edit(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/roles/")
	idStr = strings.TrimSuffix(idStr, "/edit")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	role, err := h.queries.GetRole(r.Context(), id)
	if err != nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	companies, err := h.queries.ListCompanies(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return
	}

	templates.RoleFormEdit(role, companies).Render(r.Context(), w)
}

func (h *RolesHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/roles/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	companyID, err := strconv.ParseInt(r.FormValue("company_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	err = h.queries.UpdateRole(r.Context(), db.UpdateRoleParams{
		CompanyID:           companyID,
		Name:                r.FormValue("name"),
		Url:                 nullString(r.FormValue("url")),
		Description:         nullString(r.FormValue("description")),
		CoverLetter:         nullString(r.FormValue("cover_letter")),
		ApplicationLocation: nullString(r.FormValue("application_location")),
		AppliedDate:         nullString(r.FormValue("applied_date")),
		ClosedDate:          nullString(r.FormValue("closed_date")),
		PostedRangeMin:      nullInt64(r.FormValue("posted_range_min")),
		PostedRangeMax:      nullInt64(r.FormValue("posted_range_max")),
		Equity:              nullBool(r.FormValue("equity")),
		WorkCity:            nullString(r.FormValue("work_city")),
		WorkState:           nullString(r.FormValue("work_state")),
		Location:            nullString(r.FormValue("location")),
		Status:              nullString(r.FormValue("status")),
		Discovery:           nullString(r.FormValue("discovery")),
		Referral:            nullBool(r.FormValue("referral")),
		Notes:               nullString(r.FormValue("notes")),
		RoleID:              id,
	})
	if err != nil {
		http.Error(w, "Failed to update role", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/roles", http.StatusSeeOther)
}

func (h *RolesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/roles/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteRole(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete role", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/roles", http.StatusSeeOther)
}
