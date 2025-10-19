package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"reverse-ats/internal/db"
	"reverse-ats/internal/templates"
)

type InterviewsHandler struct {
	queries *db.Queries
}

func NewInterviewsHandler(queries *db.Queries) *InterviewsHandler {
	return &InterviewsHandler{queries: queries}
}

func (h *InterviewsHandler) List(w http.ResponseWriter, r *http.Request) {
	interviews, err := h.queries.ListInterviewsWithRole(r.Context())
	if err != nil {
		http.Error(w, "Failed to fetch interviews", http.StatusInternalServerError)
		return
	}

	templates.InterviewsList(interviews).Render(r.Context(), w)
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
