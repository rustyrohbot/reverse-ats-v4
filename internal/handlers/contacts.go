package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"

	"reverse-ats/internal/models"
	"reverse-ats/internal/templates"
	"reverse-ats/internal/util"
)

type ContactsHandler struct {
	app *pocketbase.PocketBase
}

func NewContactsHandler(app *pocketbase.PocketBase) *ContactsHandler {
	return &ContactsHandler{app: app}
}

func recordToContact(record *core.Record) models.Contact {
	contact := models.Contact{
		ID:        record.Id,
		CompanyID: record.GetString("company"),
		FirstName: record.GetString("first_name"),
		LastName:  record.GetString("last_name"),
		Role:      record.GetString("role"),
		Email:     record.GetString("email"),
		Phone:     record.GetString("phone"),
		Linkedin:  record.GetString("linkedin"),
		Notes:     record.GetString("notes"),
		CreatedAt: record.GetDateTime("created").String(),
		UpdatedAt: record.GetDateTime("updated").String(),
	}

	// Get company name from expanded relation
	if companyRecord := record.ExpandedOne("company"); companyRecord != nil {
		contact.CompanyName = companyRecord.GetString("name")
	}

	return contact
}

func sortContactsByCompanyName(contacts []models.Contact, order string) {
	sort.Slice(contacts, func(i, j int) bool {
		cmpResult := strings.Compare(
			strings.ToLower(contacts[i].CompanyName),
			strings.ToLower(contacts[j].CompanyName),
		)
		if order == "desc" {
			return cmpResult > 0
		}
		return cmpResult < 0
	})
}

func (h *ContactsHandler) List(w http.ResponseWriter, r *http.Request) error {
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Validate sort field
	validSortFields := map[string]bool{
		"company_name": true,
		"first_name":   true,
		"last_name":    true,
	}

	if sortBy == "" || !validSortFields[sortBy] {
		sortBy = "first_name"
	}

	if order != "asc" && order != "desc" {
		order = "asc"
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

	// Fetch contacts
	records, err := h.app.FindRecordsByFilter(
		"contacts",
		"",
		sortField,
		-1, // all records
		0,
	)
	if err != nil {
		http.Error(w, "Failed to fetch contacts", http.StatusInternalServerError)
		return err
	}

	// Fetch all companies once to avoid N+1 queries
	companiesMap, err := util.FetchCompaniesMap(h.app)
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return err
	}

	// Convert records to Contact structs with company names
	contacts := make([]models.Contact, len(records))
	for i, record := range records {
		contact := recordToContact(record)
		// Look up company name from map
		if companyID := record.GetString("company"); companyID != "" {
			if companyName, ok := companiesMap[companyID]; ok {
				contact.CompanyName = companyName
			}
		}
		contacts[i] = contact
	}

	// Sort in memory if needed
	if doInMemorySort {
		sortContactsByCompanyName(contacts, order)
	}

	// Fetch companies for inline form dropdown
	companies, err := util.FetchCompaniesForDropdown(h.app)
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return err
	}

	return templates.ContactsList(contacts, sortBy, order, companies).Render(r.Context(), w)
}

func (h *ContactsHandler) New(w http.ResponseWriter, r *http.Request) error {
	// Fetch companies for dropdown
	companies, err := util.FetchCompaniesForDropdown(h.app)
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return err
	}

	return templates.ContactFormNew(companies).Render(r.Context(), w)
}

func (h *ContactsHandler) Create(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return err
	}

	collection, err := h.app.FindCollectionByNameOrId(util.CollectionContacts)
	if err != nil {
		http.Error(w, "Failed to find collection", http.StatusInternalServerError)
		return err
	}

	record := core.NewRecord(collection)
	record.Set("company", r.FormValue("company"))
	record.Set("first_name", r.FormValue("first_name"))
	record.Set("last_name", r.FormValue("last_name"))
	record.Set("role", r.FormValue("role"))
	record.Set("email", r.FormValue("email"))
	record.Set("phone", r.FormValue("phone"))
	record.Set("linkedin", r.FormValue("linkedin"))
	record.Set("notes", r.FormValue("notes"))

	if err := h.app.Save(record); err != nil {
		http.Error(w, "Failed to create contact", http.StatusInternalServerError)
		return err
	}

	// If HTMX request, return just the new row
	if r.Header.Get("HX-Request") == "true" {
		contact := recordToContact(record)
		// Fetch company name for display
		if companyID := record.GetString("company"); companyID != "" {
			if companyRecord, err := h.app.FindRecordById(util.CollectionCompanies, companyID); err == nil {
				contact.CompanyName = companyRecord.GetString("name")
			}
		}
		return templates.ContactRow(contact).Render(r.Context(), w)
	}

	// Otherwise redirect
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	return nil
}

func (h *ContactsHandler) Edit(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL path parameter
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return fmt.Errorf("missing id parameter")
	}

	record, err := h.app.FindRecordById(util.CollectionContacts, id)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return err
	}

	contact := recordToContact(record)

	// Fetch company name for display
	if companyID := record.GetString("company"); companyID != "" {
		if companyRecord, err := h.app.FindRecordById(util.CollectionCompanies, companyID); err == nil {
			contact.CompanyName = companyRecord.GetString("name")
		}
	}

	// Fetch companies for dropdown
	companies, err := util.FetchCompaniesForDropdown(h.app)
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return err
	}

	return templates.ContactFormEdit(contact, companies).Render(r.Context(), w)
}

func (h *ContactsHandler) Update(w http.ResponseWriter, r *http.Request) error {
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

	record, err := h.app.FindRecordById(util.CollectionContacts, id)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return err
	}

	record.Set("company", r.FormValue("company"))
	record.Set("first_name", r.FormValue("first_name"))
	record.Set("last_name", r.FormValue("last_name"))
	record.Set("role", r.FormValue("role"))
	record.Set("email", r.FormValue("email"))
	record.Set("phone", r.FormValue("phone"))
	record.Set("linkedin", r.FormValue("linkedin"))
	record.Set("notes", r.FormValue("notes"))

	if err := h.app.Save(record); err != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return err
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	return nil
}

func (h *ContactsHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	// Extract ID from URL path parameter
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return fmt.Errorf("missing id parameter")
	}

	record, err := h.app.FindRecordById(util.CollectionContacts, id)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return err
	}

	if err := h.app.Delete(record); err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return err
	}

	// If HTMX request, return empty response (row will be removed)
	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return nil
	}

	// Otherwise redirect
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	return nil
}
