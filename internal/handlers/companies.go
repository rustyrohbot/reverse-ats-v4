package handlers

import (
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"reverse-ats/internal/models"
	"reverse-ats/internal/templates"
	"reverse-ats/internal/util"
)

type CompaniesHandler struct {
	app *pocketbase.PocketBase
}

func NewCompaniesHandler(app *pocketbase.PocketBase) *CompaniesHandler {
	return &CompaniesHandler{app: app}
}

func recordToCompany(record *core.Record) models.Company {
	return models.Company{
		ID:          record.Id,
		Name:        record.GetString("name"),
		Description: record.GetString("description"),
		Url:         record.GetString("url"),
		Linkedin:    record.GetString("linkedin"),
		HqCity:      record.GetString("hq_city"),
		HqState:     record.GetString("hq_state"),
		CreatedAt:   record.GetDateTime("created").String(),
		UpdatedAt:   record.GetDateTime("updated").String(),
	}
}

func (h *CompaniesHandler) List(w http.ResponseWriter, r *http.Request) error {
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Validate sort column
	validSortColumns := map[string]bool{
		"name":     true,
		"hq_city":  true,
		"hq_state": true,
	}

	if sortBy == "" || !validSortColumns[sortBy] {
		sortBy = "name"
	}

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Fetch companies from PocketBase
	sortField := sortBy
	if order == "desc" {
		sortField = "-" + sortBy
	}
	records, err := h.app.FindRecordsByFilter(
		"companies",
		"",
		sortField,
		-1, // all records
		0,
	)
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return err
	}

	// Convert records to Company structs
	companies := make([]models.Company, len(records))
	for i, record := range records {
		companies[i] = recordToCompany(record)
	}

	return templates.CompaniesList(companies, sortBy, order).Render(r.Context(), w)
}

func (h *CompaniesHandler) New(w http.ResponseWriter, r *http.Request) error {
	return templates.CompanyFormNew().Render(r.Context(), w)
}

func (h *CompaniesHandler) Create(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return err
	}

	collection, err := h.app.FindCollectionByNameOrId(util.CollectionCompanies)
	if err != nil {
		http.Error(w, "Failed to find collection", http.StatusInternalServerError)
		return err
	}

	record := core.NewRecord(collection)
	record.Set("name", r.FormValue("name"))
	record.Set("description", r.FormValue("description"))
	record.Set("url", r.FormValue("url"))
	record.Set("linkedin", r.FormValue("linkedin"))
	record.Set("hq_city", r.FormValue("hq_city"))
	record.Set("hq_state", r.FormValue("hq_state"))

	if err := h.app.Save(record); err != nil {
		http.Error(w, "Failed to create company", http.StatusInternalServerError)
		return err
	}

	// If HTMX request, return just the new row
	if r.Header.Get("HX-Request") == "true" {
		company := recordToCompany(record)
		return templates.CompanyRow(company).Render(r.Context(), w)
	}

	// Otherwise redirect
	http.Redirect(w, r, "/companies", http.StatusSeeOther)
	return nil
}

func (h *CompaniesHandler) Edit(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL path parameter
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return fmt.Errorf("missing id parameter")
	}

	record, err := h.app.FindRecordById(util.CollectionCompanies, id)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return err
	}

	company := recordToCompany(record)
	return templates.CompanyFormEdit(company).Render(r.Context(), w)
}

func (h *CompaniesHandler) Update(w http.ResponseWriter, r *http.Request) error {
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

	record, err := h.app.FindRecordById(util.CollectionCompanies, id)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return err
	}

	record.Set("name", r.FormValue("name"))
	record.Set("description", r.FormValue("description"))
	record.Set("url", r.FormValue("url"))
	record.Set("linkedin", r.FormValue("linkedin"))
	record.Set("hq_city", r.FormValue("hq_city"))
	record.Set("hq_state", r.FormValue("hq_state"))

	if err := h.app.Save(record); err != nil {
		http.Error(w, "Failed to update company", http.StatusInternalServerError)
		return err
	}

	http.Redirect(w, r, "/companies", http.StatusSeeOther)
	return nil
}

func (h *CompaniesHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL path parameter
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return fmt.Errorf("missing id parameter")
	}

	record, err := h.app.FindRecordById(util.CollectionCompanies, id)
	if err != nil {
		http.Error(w, "Company not found", http.StatusNotFound)
		return err
	}

	if err := h.app.Delete(record); err != nil {
		http.Error(w, "Failed to delete company", http.StatusInternalServerError)
		return err
	}

	// If HTMX request, return empty response (row will be removed)
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	// Otherwise redirect
	http.Redirect(w, r, "/companies", http.StatusSeeOther)
	return nil
}

