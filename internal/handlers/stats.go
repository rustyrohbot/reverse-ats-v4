package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"

	"reverse-ats/internal/templates"
	"reverse-ats/internal/util"
)

type StatsHandler struct {
	app *pocketbase.PocketBase
}

func NewStatsHandler(app *pocketbase.PocketBase) *StatsHandler {
	return &StatsHandler{app: app}
}

func (h *StatsHandler) Show(w http.ResponseWriter, r *http.Request) error {
	dateRange := r.URL.Query().Get("range")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	// Default to last 30 days
	if dateRange == "" && startDate == "" {
		dateRange = "30"
	}

	// Calculate date boundaries
	var startDateFilter, endDateFilter time.Time
	now := time.Now()
	var filterDates bool

	if dateRange == "custom" && startDate != "" && endDate != "" {
		var err error
		startDateFilter, err = time.Parse("2006-01-02", startDate)
		if err == nil {
			endDateFilter, err = time.Parse("2006-01-02", endDate)
			if err == nil {
				filterDates = true
			}
		}
	} else if dateRange == "all" {
		// No date filter
		filterDates = false
	} else {
		// Calculate based on range
		var days int
		switch dateRange {
		case "7":
			days = 7
		case "90":
			days = 90
		case "180":
			days = 180
		case "365":
			days = 365
		default:
			days = 30
			dateRange = "30"
		}
		startDateFilter = now.AddDate(0, 0, -days)
		endDateFilter = now
		filterDates = true
	}

	// Build SQL date filter
	var dateClause string
	var whereClause string
	if filterDates {
		startStr := startDateFilter.Format("2006-01-02")
		endStr := endDateFilter.Format("2006-01-02")

		// PocketBase stores dates in ISO format, so we can compare directly
		dateClause = " AND applied_date BETWEEN '" + startStr + "' AND '" + endStr + "'"
		whereClause = " WHERE applied_date IS NOT NULL AND applied_date != '' AND applied_date BETWEEN '" + startStr + "' AND '" + endStr + "'"
	}

	stats := templates.StatsData{
		DateRange: dateRange,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Get database connection
	db := h.app.DB()

	// Query: First and Last Application Dates
	var firstDateQuery, lastDateQuery string
	if whereClause != "" {
		firstDateQuery = "SELECT applied_date FROM roles " + whereClause +
			" ORDER BY applied_date ASC LIMIT 1"
		lastDateQuery = "SELECT applied_date FROM roles " + whereClause +
			" ORDER BY applied_date DESC LIMIT 1"
	} else {
		firstDateQuery = "SELECT applied_date FROM roles WHERE applied_date IS NOT NULL AND applied_date != '' " +
			"ORDER BY applied_date ASC LIMIT 1"
		lastDateQuery = "SELECT applied_date FROM roles WHERE applied_date IS NOT NULL AND applied_date != '' " +
			"ORDER BY applied_date DESC LIMIT 1"
	}
	var firstDate, lastDate sql.NullString
	db.NewQuery(firstDateQuery).Row(&firstDate)
	db.NewQuery(lastDateQuery).Row(&lastDate)
	if firstDate.Valid && firstDate.String != "" {
		stats.FirstApplicationDate = util.FormatDateToText(firstDate.String)
	}
	if lastDate.Valid && lastDate.String != "" {
		stats.LastApplicationDate = util.FormatDateToText(lastDate.String)
	}

	// Query: Roles Applied (count of roles with applied_date in range)
	rolesQuery := "SELECT COUNT(*) FROM roles WHERE applied_date IS NOT NULL AND applied_date != ''" + dateClause
	db.NewQuery(rolesQuery).Row(&stats.RolesApplied)

	// Query: Offers Received
	var offersQuery string
	if whereClause != "" {
		offersQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'OFFER'"
	} else {
		offersQuery = "SELECT COUNT(*) FROM roles WHERE status = 'OFFER'"
	}
	db.NewQuery(offersQuery).Row(&stats.OffersReceived)

	// Query: Rejections
	var rejectionsQuery string
	if whereClause != "" {
		rejectionsQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'REJECTED'"
	} else {
		rejectionsQuery = "SELECT COUNT(*) FROM roles WHERE status = 'REJECTED'"
	}
	db.NewQuery(rejectionsQuery).Row(&stats.Rejections)

	// Query: Interviewing
	var interviewingQuery string
	if whereClause != "" {
		interviewingQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'INTERVIEWING'"
	} else {
		interviewingQuery = "SELECT COUNT(*) FROM roles WHERE status = 'INTERVIEWING'"
	}
	db.NewQuery(interviewingQuery).Row(&stats.Interviewing)

	// Query: Ghosted
	var ghostedQuery string
	if whereClause != "" {
		ghostedQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'GHOSTED'"
	} else {
		ghostedQuery = "SELECT COUNT(*) FROM roles WHERE status = 'GHOSTED'"
	}
	db.NewQuery(ghostedQuery).Row(&stats.Ghosted)

	// Query: Freeze
	var freezeQuery string
	if whereClause != "" {
		freezeQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'FREEZE'"
	} else {
		freezeQuery = "SELECT COUNT(*) FROM roles WHERE status = 'FREEZE'"
	}
	db.NewQuery(freezeQuery).Row(&stats.Freeze)

	// Query: Withdrew
	var withdrewQuery string
	if whereClause != "" {
		withdrewQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'WITHDREW'"
	} else {
		withdrewQuery = "SELECT COUNT(*) FROM roles WHERE status = 'WITHDREW'"
	}
	db.NewQuery(withdrewQuery).Row(&stats.Withdrew)

	// Query: Average Posted Min
	var avgMinQuery string
	if whereClause != "" {
		avgMinQuery = "SELECT AVG(CAST(posted_range_min AS REAL)) FROM roles" + whereClause + " AND posted_range_min IS NOT NULL AND posted_range_min != 0"
	} else {
		avgMinQuery = "SELECT AVG(CAST(posted_range_min AS REAL)) FROM roles WHERE posted_range_min IS NOT NULL AND posted_range_min != 0"
	}
	var avgMin sql.NullFloat64
	db.NewQuery(avgMinQuery).Row(&avgMin)
	if avgMin.Valid {
		stats.AvgPostedMin = avgMin.Float64
	}

	// Query: Average Posted Max
	var avgMaxQuery string
	if whereClause != "" {
		avgMaxQuery = "SELECT AVG(CAST(posted_range_max AS REAL)) FROM roles" + whereClause + " AND posted_range_max IS NOT NULL AND posted_range_max != 0"
	} else {
		avgMaxQuery = "SELECT AVG(CAST(posted_range_max AS REAL)) FROM roles WHERE posted_range_max IS NOT NULL AND posted_range_max != 0"
	}
	var avgMax sql.NullFloat64
	db.NewQuery(avgMaxQuery).Row(&avgMax)
	if avgMax.Valid {
		stats.AvgPostedMax = avgMax.Float64
	}

	// Query: Absolute Posted Min
	var absMinQuery string
	if whereClause != "" {
		absMinQuery = "SELECT MIN(posted_range_min) FROM roles" + whereClause + " AND posted_range_min IS NOT NULL AND posted_range_min != 0"
	} else {
		absMinQuery = "SELECT MIN(posted_range_min) FROM roles WHERE posted_range_min IS NOT NULL AND posted_range_min != 0"
	}
	var absMin sql.NullInt64
	db.NewQuery(absMinQuery).Row(&absMin)
	if absMin.Valid {
		stats.AbsPostedMin = absMin.Int64
	}

	// Query: Absolute Posted Max
	var absMaxQuery string
	if whereClause != "" {
		absMaxQuery = "SELECT MAX(posted_range_max) FROM roles" + whereClause + " AND posted_range_max IS NOT NULL AND posted_range_max != 0"
	} else {
		absMaxQuery = "SELECT MAX(posted_range_max) FROM roles WHERE posted_range_max IS NOT NULL AND posted_range_max != 0"
	}
	var absMax sql.NullInt64
	db.NewQuery(absMaxQuery).Row(&absMax)
	if absMax.Valid {
		stats.AbsPostedMax = absMax.Int64
	}

	// Query: Remote Roles
	var remoteQuery string
	if whereClause != "" {
		remoteQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND location = 'REMOTE'"
	} else {
		remoteQuery = "SELECT COUNT(*) FROM roles WHERE location = 'REMOTE'"
	}
	db.NewQuery(remoteQuery).Row(&stats.RemoteRoles)

	// Query: Hybrid Roles
	var hybridQuery string
	if whereClause != "" {
		hybridQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND location = 'HYBRID'"
	} else {
		hybridQuery = "SELECT COUNT(*) FROM roles WHERE location = 'HYBRID'"
	}
	db.NewQuery(hybridQuery).Row(&stats.HybridRoles)

	// Query: Onsite Roles
	var onsiteQuery string
	if whereClause != "" {
		onsiteQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND location = 'ONSITE'"
	} else {
		onsiteQuery = "SELECT COUNT(*) FROM roles WHERE location = 'ONSITE'"
	}
	db.NewQuery(onsiteQuery).Row(&stats.OnsiteRoles)

	// For interviews, apply the same date filter (using 'date' column)
	var interviewWhereClause string
	if filterDates {
		startStr := startDateFilter.Format("2006-01-02")
		endStr := endDateFilter.Format("2006-01-02")
		interviewWhereClause = " WHERE date BETWEEN '" + startStr + "' AND '" + endStr + "'"
	}

	// Query: Total Interviews
	var totalInterviewsQuery string
	if interviewWhereClause != "" {
		totalInterviewsQuery = "SELECT COUNT(*) FROM interviews" + interviewWhereClause
	} else {
		totalInterviewsQuery = "SELECT COUNT(*) FROM interviews"
	}
	db.NewQuery(totalInterviewsQuery).Row(&stats.TotalInterviews)

	// Query: Recruiter Interviews
	var recruiterQuery string
	if interviewWhereClause != "" {
		recruiterQuery = "SELECT COUNT(*) FROM interviews" + interviewWhereClause + " AND type = 'RECRUITER'"
	} else {
		recruiterQuery = "SELECT COUNT(*) FROM interviews WHERE type = 'RECRUITER'"
	}
	db.NewQuery(recruiterQuery).Row(&stats.RecruiterInterviews)

	// Query: Manager Interviews
	var managerQuery string
	if interviewWhereClause != "" {
		managerQuery = "SELECT COUNT(*) FROM interviews" + interviewWhereClause + " AND type = 'MANAGER'"
	} else {
		managerQuery = "SELECT COUNT(*) FROM interviews WHERE type = 'MANAGER'"
	}
	db.NewQuery(managerQuery).Row(&stats.ManagerInterviews)

	// Query: Loop Interviews
	var loopQuery string
	if interviewWhereClause != "" {
		loopQuery = "SELECT COUNT(*) FROM interviews" + interviewWhereClause + " AND type = 'LOOP'"
	} else {
		loopQuery = "SELECT COUNT(*) FROM interviews WHERE type = 'LOOP'"
	}
	db.NewQuery(loopQuery).Row(&stats.LoopInterviews)

	// Query: Tech Screen Interviews
	var techScreenQuery string
	if interviewWhereClause != "" {
		techScreenQuery = "SELECT COUNT(*) FROM interviews" + interviewWhereClause + " AND type = 'TECH_SCREEN'"
	} else {
		techScreenQuery = "SELECT COUNT(*) FROM interviews WHERE type = 'TECH_SCREEN'"
	}
	db.NewQuery(techScreenQuery).Row(&stats.TechScreenInterviews)

	return templates.Stats(stats).Render(r.Context(), w)
}
