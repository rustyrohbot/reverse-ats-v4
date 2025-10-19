package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"reverse-ats/internal/db"
	"reverse-ats/internal/templates"
)

type RolesHandler struct {
	queries *db.Queries
}

func NewRolesHandler(queries *db.Queries) *RolesHandler {
	return &RolesHandler{queries: queries}
}

func (h *RolesHandler) List(w http.ResponseWriter, r *http.Request) {
	roles, err := h.queries.ListRolesWithCompany(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}

	templates.RolesList(roles).Render(r.Context(), w)
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
		Equity:              nullString(r.FormValue("equity")),
		WorkCity:            nullString(r.FormValue("work_city")),
		WorkState:           nullString(r.FormValue("work_state")),
		Location:            nullString(r.FormValue("location")),
		Status:              nullString(r.FormValue("status")),
		Discovery:           nullString(r.FormValue("discovery")),
		Referral:            nullString(r.FormValue("referral")),
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
		Equity:              nullString(r.FormValue("equity")),
		WorkCity:            nullString(r.FormValue("work_city")),
		WorkState:           nullString(r.FormValue("work_state")),
		Location:            nullString(r.FormValue("location")),
		Status:              nullString(r.FormValue("status")),
		Discovery:           nullString(r.FormValue("discovery")),
		Referral:            nullString(r.FormValue("referral")),
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
