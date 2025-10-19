package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"reverse-ats/internal/db"
	"reverse-ats/internal/templates"
)

type CompaniesHandler struct {
	queries *db.Queries
}

func NewCompaniesHandler(queries *db.Queries) *CompaniesHandler {
	return &CompaniesHandler{queries: queries}
}

func (h *CompaniesHandler) List(w http.ResponseWriter, r *http.Request) {
	companies, err := h.queries.ListCompanies(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return
	}

	templates.CompaniesList(companies).Render(r.Context(), w)
}

func (h *CompaniesHandler) New(w http.ResponseWriter, r *http.Request) {
	templates.CompanyFormNew().Render(r.Context(), w)
}

func (h *CompaniesHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := nullString(r.FormValue("description"))
	url := nullString(r.FormValue("url"))
	hqCity := nullString(r.FormValue("hq_city"))
	hqState := nullString(r.FormValue("hq_state"))

	_, err := h.queries.CreateCompany(r.Context(), db.CreateCompanyParams{
		Name:        name,
		Description: description,
		Url:         url,
		HqCity:      hqCity,
		HqState:     hqState,
	})
	if err != nil {
		http.Error(w, "Failed to create company", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/companies", http.StatusSeeOther)
}

func (h *CompaniesHandler) Edit(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/companies/")
	idStr = strings.TrimSuffix(idStr, "/edit")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	company, err := h.queries.GetCompany(r.Context(), id)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return
	}

	templates.CompanyFormEdit(company).Render(r.Context(), w)
}

func (h *CompaniesHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/companies/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := nullString(r.FormValue("description"))
	url := nullString(r.FormValue("url"))
	hqCity := nullString(r.FormValue("hq_city"))
	hqState := nullString(r.FormValue("hq_state"))

	err = h.queries.UpdateCompany(r.Context(), db.UpdateCompanyParams{
		Name:        name,
		Description: description,
		Url:         url,
		HqCity:      hqCity,
		HqState:     hqState,
		CompanyID:   id,
	})
	if err != nil {
		http.Error(w, "Failed to update company", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/companies", http.StatusSeeOther)
}

func (h *CompaniesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/companies/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteCompany(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete company", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/companies", http.StatusSeeOther)
}

