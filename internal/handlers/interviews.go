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

type InterviewsHandler struct {
	queries *db.Queries
	dbConn  *sql.DB
}

func NewInterviewsHandler(queries *db.Queries, dbConn *sql.DB) *InterviewsHandler {
	return &InterviewsHandler{queries: queries, dbConn: dbConn}
}

func (h *InterviewsHandler) List(w http.ResponseWriter, r *http.Request) {
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	// Map display names to actual column/field names
	sortColumnMap := map[string]string{
		"company_name": "c.name",
		"date":         "i.date",
	}

	sortCol, ok := sortColumnMap[sortBy]
	if !ok {
		sortBy = "date"
		sortCol = "i.date"
	}

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	query := fmt.Sprintf(`
		SELECT
			i.interview_id,
			i.role_id,
			i.date,
			i.start,
			i.end,
			i.notes,
			i.type,
			i.created_at,
			i.updated_at,
			r.name as role_name,
			c.company_id,
			c.name as company_name
		FROM interviews i
		INNER JOIN roles r ON i.role_id = r.role_id
		INNER JOIN companies c ON r.company_id = c.company_id
		ORDER BY %s %s, i.start %s`, sortCol, strings.ToUpper(order), strings.ToUpper(order))

	rows, err := h.dbConn.QueryContext(r.Context(), query)
	if err != nil {
		http.Error(w, "Failed to fetch interviews", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var interviews []db.ListInterviewsWithRoleRow
	for rows.Next() {
		var interview db.ListInterviewsWithRoleRow
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
			&interview.RoleName,
			&interview.CompanyID,
			&interview.CompanyName,
		)
		if err != nil {
			http.Error(w, "Failed to scan interviews", http.StatusInternalServerError)
			return
		}
		interviews = append(interviews, interview)
	}

	templates.InterviewsList(interviews, sortBy, order).Render(r.Context(), w)
}

func (h *InterviewsHandler) New(w http.ResponseWriter, r *http.Request) {
	roles, err := h.queries.ListRolesWithCompany(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}

	templates.InterviewFormNew(roles).Render(r.Context(), w)
}

func (h *InterviewsHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	roleID, err := strconv.ParseInt(r.FormValue("role_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	_, err = h.queries.CreateInterview(r.Context(), db.CreateInterviewParams{
		RoleID: roleID,
		Date:   r.FormValue("date"),
		Start:  r.FormValue("start"),
		End:    r.FormValue("end"),
		Notes:  nullString(r.FormValue("notes")),
		Type:   r.FormValue("type"),
	})
	if err != nil {
		http.Error(w, "Failed to create interview", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/interviews", http.StatusSeeOther)
}

func (h *InterviewsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/interviews/")
	idStr = strings.TrimSuffix(idStr, "/edit")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	interview, err := h.queries.GetInterview(r.Context(), id)
	if err != nil {
		http.Error(w, "Interview not found", http.StatusNotFound)
		return
	}

	roles, err := h.queries.ListRolesWithCompany(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch roles", http.StatusInternalServerError)
		return
	}

	templates.InterviewFormEdit(interview, roles).Render(r.Context(), w)
}

func (h *InterviewsHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/interviews/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	roleID, err := strconv.ParseInt(r.FormValue("role_id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = h.queries.UpdateInterview(r.Context(), db.UpdateInterviewParams{
		RoleID:      roleID,
		Date:        r.FormValue("date"),
		Start:       r.FormValue("start"),
		End:         r.FormValue("end"),
		Notes:       nullString(r.FormValue("notes")),
		Type:        r.FormValue("type"),
		InterviewID: id,
	})
	if err != nil {
		http.Error(w, "Failed to update interview", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/interviews", http.StatusSeeOther)
}

func (h *InterviewsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/interviews/")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = h.queries.DeleteInterview(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete interview", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/interviews", http.StatusSeeOther)
}
