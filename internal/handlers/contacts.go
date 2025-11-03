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

type ContactsHandler struct {
	queries *db.Queries
	dbConn  *sql.DB
}

func NewContactsHandler(queries *db.Queries, dbConn *sql.DB) *ContactsHandler {
	return &ContactsHandler{queries: queries, dbConn: dbConn}
}

func (h *ContactsHandler) List(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Map display names to actual column/field names
	sortColumnMap := map[string]string{
		"company_name": "c.name",
		"first_name":   "ct.first_name",
		"last_name":    "ct.last_name",
	}

	sortCol, ok := sortColumnMap[sortBy]
	if !ok {
		sortBy = "company_name"
		sortCol = "c.name"
	}

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	query := fmt.Sprintf(`
		SELECT
			ct.contact_id,
			ct.company_id,
			ct.first_name,
			ct.last_name,
			ct.role,
			ct.email,
			ct.phone,
			ct.linkedin,
			ct.notes,
			ct.created_at,
			ct.updated_at,
			c.name as company_name
		FROM contacts ct
		INNER JOIN companies c ON ct.company_id = c.company_id
		ORDER BY %s %s, ct.last_name ASC, ct.first_name ASC`, sortCol, strings.ToUpper(order))

	rows, err := h.dbConn.QueryContext(r.Context(), query)
	if err != nil {
		http.Error(w, "Failed to fetch contacts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var contacts []db.ListContactsWithCompanyRow
	for rows.Next() {
		var contact db.ListContactsWithCompanyRow
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
			&contact.CompanyName,
		)
		if err != nil {
			http.Error(w, "Failed to scan contacts", http.StatusInternalServerError)
			return
		}
		contacts = append(contacts, contact)
	}

	templates.ContactsList(contacts, sortBy, order).Render(r.Context(), w)
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
