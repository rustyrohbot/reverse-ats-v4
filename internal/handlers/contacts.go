package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"reverse-ats/internal/db"
	"reverse-ats/internal/templates"
)

type ContactsHandler struct {
	queries *db.Queries
}

func NewContactsHandler(queries *db.Queries) *ContactsHandler {
	return &ContactsHandler{queries: queries}
}

func (h *ContactsHandler) List(w http.ResponseWriter, r *http.Request) {
	contacts, err := h.queries.ListContactsWithCompany(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch contacts", http.StatusInternalServerError)
		return
	}

	templates.ContactsList(contacts).Render(r.Context(), w)
}

func (h *ContactsHandler) New(w http.ResponseWriter, r *http.Request) {
	companies, err := h.queries.ListCompanies(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return
	}

	templates.ContactFormNew(companies).Render(r.Context(), w)
}

func (h *ContactsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	companyID, err := strconv.ParseInt(r.FormValue("company_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid company ID", http.StatusBadRequest)
		return
	}

	_, err = h.queries.CreateContact(r.Context(), db.CreateContactParams{
		CompanyID: companyID,
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Role:      nullString(r.FormValue("role")),
		Email:     nullString(r.FormValue("email")),
		Phone:     nullString(r.FormValue("phone")),
		Linkedin:  nullString(r.FormValue("linkedin")),
		Notes:     nullString(r.FormValue("notes")),
	})
	if err != nil {
		http.Error(w, "Failed to create contact", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (h *ContactsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/contacts/")
	idStr = strings.TrimSuffix(idStr, "/edit")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	contact, err := h.queries.GetContact(r.Context(), id)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	companies, err := h.queries.ListCompanies(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch companies", http.StatusInternalServerError)
		return
	}

	templates.ContactFormEdit(contact, companies).Render(r.Context(), w)
}

func (h *ContactsHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/contacts/")

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

	err = h.queries.UpdateContact(r.Context(), db.UpdateContactParams{
		CompanyID: companyID,
		FirstName: r.FormValue("first_name"),
		LastName:  r.FormValue("last_name"),
		Role:      nullString(r.FormValue("role")),
		Email:     nullString(r.FormValue("email")),
		Phone:     nullString(r.FormValue("phone")),
		Linkedin:  nullString(r.FormValue("linkedin")),
		Notes:     nullString(r.FormValue("notes")),
		ContactID: id,
	})
	if err != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

func (h *ContactsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/contacts/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteContact(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}
