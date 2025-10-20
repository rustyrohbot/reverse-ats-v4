package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"reverse-ats/internal/db"
	"reverse-ats/internal/templates"
)

type StatsHandler struct {
	queries *db.Queries
	dbConn  *sql.DB
}

func NewStatsHandler(queries *db.Queries, dbConn *sql.DB) *StatsHandler {
	return &StatsHandler{queries: queries, dbConn: dbConn}
}

func (h *StatsHandler) Show(w http.ResponseWriter, r *http.Request) {
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

	// Build SQL date filter using SQLite's date parsing
	var dateClause string
	var whereClause string
	if filterDates {
		// Convert text dates like "April 8, 2025" to ISO format for comparison
		startStr := startDateFilter.Format("2006-01-02")
		endStr := endDateFilter.Format("2006-01-02")

		// SQL to convert "April 8, 2025" → "2025-04-08"
		dateConversion := "substr(applied_date, -4) || '-' || " +
			"CASE substr(applied_date, 1, instr(applied_date, ' ')-1) " +
			"WHEN 'January' THEN '01' WHEN 'February' THEN '02' WHEN 'March' THEN '03' " +
			"WHEN 'April' THEN '04' WHEN 'May' THEN '05' WHEN 'June' THEN '06' " +
			"WHEN 'July' THEN '07' WHEN 'August' THEN '08' WHEN 'September' THEN '09' " +
			"WHEN 'October' THEN '10' WHEN 'November' THEN '11' WHEN 'December' THEN '12' END || '-' || " +
			"printf('%02d', CAST(replace(substr(applied_date, instr(applied_date, ' ')+1, instr(substr(applied_date, instr(applied_date, ' ')+1), ',') - 1), ' ', '') AS INTEGER))"

		dateClause = " AND (" + dateConversion + ") BETWEEN '" + startStr + "' AND '" + endStr + "'"
		whereClause = " WHERE applied_date IS NOT NULL AND (" + dateConversion + ") BETWEEN '" + startStr + "' AND '" + endStr + "'"
	}

	stats := templates.StatsData{
		DateRange: dateRange,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Query: Roles Applied (count of roles with applied_date in range)
	rolesQuery := "SELECT COUNT(*) FROM roles WHERE applied_date IS NOT NULL" + dateClause
	h.dbConn.QueryRow(rolesQuery).Scan(&stats.RolesApplied)

	// Query: Offers Received
	var offersQuery string
	if whereClause != "" {
		offersQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'OFFER'"
	} else {
		offersQuery = "SELECT COUNT(*) FROM roles WHERE status = 'OFFER'"
	}
	h.dbConn.QueryRow(offersQuery).Scan(&stats.OffersReceived)

	// Query: Rejections
	var rejectionsQuery string
	if whereClause != "" {
		rejectionsQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'REJECTED'"
	} else {
		rejectionsQuery = "SELECT COUNT(*) FROM roles WHERE status = 'REJECTED'"
	}
	h.dbConn.QueryRow(rejectionsQuery).Scan(&stats.Rejections)

	// Query: Ghosted
	var ghostedQuery string
	if whereClause != "" {
		ghostedQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'GHOSTED'"
	} else {
		ghostedQuery = "SELECT COUNT(*) FROM roles WHERE status = 'GHOSTED'"
	}
	h.dbConn.QueryRow(ghostedQuery).Scan(&stats.Ghosted)

	// Query: Freeze
	var freezeQuery string
	if whereClause != "" {
		freezeQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'FREEZE'"
	} else {
		freezeQuery = "SELECT COUNT(*) FROM roles WHERE status = 'FREEZE'"
	}
	h.dbConn.QueryRow(freezeQuery).Scan(&stats.Freeze)

	// Query: Withdrew
	var withdrewQuery string
	if whereClause != "" {
		withdrewQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND status = 'WITHDREW'"
	} else {
		withdrewQuery = "SELECT COUNT(*) FROM roles WHERE status = 'WITHDREW'"
	}
	h.dbConn.QueryRow(withdrewQuery).Scan(&stats.Withdrew)

	// Query: Average Posted Min
	var avgMinQuery string
	if whereClause != "" {
		avgMinQuery = "SELECT AVG(CAST(posted_range_min AS REAL)) FROM roles" + whereClause + " AND posted_range_min IS NOT NULL"
	} else {
		avgMinQuery = "SELECT AVG(CAST(posted_range_min AS REAL)) FROM roles WHERE posted_range_min IS NOT NULL"
	}
	var avgMin sql.NullFloat64
	h.dbConn.QueryRow(avgMinQuery).Scan(&avgMin)
	if avgMin.Valid {
		stats.AvgPostedMin = avgMin.Float64
	}

	// Query: Average Posted Max
	var avgMaxQuery string
	if whereClause != "" {
		avgMaxQuery = "SELECT AVG(CAST(posted_range_max AS REAL)) FROM roles" + whereClause + " AND posted_range_max IS NOT NULL"
	} else {
		avgMaxQuery = "SELECT AVG(CAST(posted_range_max AS REAL)) FROM roles WHERE posted_range_max IS NOT NULL"
	}
	var avgMax sql.NullFloat64
	h.dbConn.QueryRow(avgMaxQuery).Scan(&avgMax)
	if avgMax.Valid {
		stats.AvgPostedMax = avgMax.Float64
	}

	// Query: Absolute Posted Min
	var absMinQuery string
	if whereClause != "" {
		absMinQuery = "SELECT MIN(posted_range_min) FROM roles" + whereClause + " AND posted_range_min IS NOT NULL"
	} else {
		absMinQuery = "SELECT MIN(posted_range_min) FROM roles WHERE posted_range_min IS NOT NULL"
	}
	var absMin sql.NullInt64
	h.dbConn.QueryRow(absMinQuery).Scan(&absMin)
	if absMin.Valid {
		stats.AbsPostedMin = absMin.Int64
	}

	// Query: Absolute Posted Max
	var absMaxQuery string
	if whereClause != "" {
		absMaxQuery = "SELECT MAX(posted_range_max) FROM roles" + whereClause + " AND posted_range_max IS NOT NULL"
	} else {
		absMaxQuery = "SELECT MAX(posted_range_max) FROM roles WHERE posted_range_max IS NOT NULL"
	}
	var absMax sql.NullInt64
	h.dbConn.QueryRow(absMaxQuery).Scan(&absMax)
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
	h.dbConn.QueryRow(remoteQuery).Scan(&stats.RemoteRoles)

	// Query: Hybrid Roles
	var hybridQuery string
	if whereClause != "" {
		hybridQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND location = 'HYBRID'"
	} else {
		hybridQuery = "SELECT COUNT(*) FROM roles WHERE location = 'HYBRID'"
	}
	h.dbConn.QueryRow(hybridQuery).Scan(&stats.HybridRoles)

	// Query: Onsite Roles
	var onsiteQuery string
	if whereClause != "" {
		onsiteQuery = "SELECT COUNT(*) FROM roles" + whereClause + " AND location = 'ONSITE'"
	} else {
		onsiteQuery = "SELECT COUNT(*) FROM roles WHERE location = 'ONSITE'"
	}
	h.dbConn.QueryRow(onsiteQuery).Scan(&stats.OnsiteRoles)

	// For interviews, apply the same date filter (using 'date' column)
	var interviewWhereClause string
	if filterDates {
		startStr := startDateFilter.Format("2006-01-02")
		endStr := endDateFilter.Format("2006-01-02")

		// SQL to convert interview date field "April 8, 2025" → "2025-04-08"
		interviewDateConversion := "substr(date, -4) || '-' || " +
			"CASE substr(date, 1, instr(date, ' ')-1) " +
			"WHEN 'January' THEN '01' WHEN 'February' THEN '02' WHEN 'March' THEN '03' " +
			"WHEN 'April' THEN '04' WHEN 'May' THEN '05' WHEN 'June' THEN '06' " +
			"WHEN 'July' THEN '07' WHEN 'August' THEN '08' WHEN 'September' THEN '09' " +
			"WHEN 'October' THEN '10' WHEN 'November' THEN '11' WHEN 'December' THEN '12' END || '-' || " +
			"printf('%02d', CAST(replace(substr(date, instr(date, ' ')+1, instr(substr(date, instr(date, ' ')+1), ',') - 1), ' ', '') AS INTEGER))"

		interviewWhereClause = " WHERE (" + interviewDateConversion + ") BETWEEN '" + startStr + "' AND '" + endStr + "'"
	}

	// Query: Total Interviews
	var totalInterviewsQuery string
	if interviewWhereClause != "" {
		totalInterviewsQuery = "SELECT COUNT(*) FROM interviews" + interviewWhereClause
	} else {
		totalInterviewsQuery = "SELECT COUNT(*) FROM interviews"
	}
	h.dbConn.QueryRow(totalInterviewsQuery).Scan(&stats.TotalInterviews)

	// Query: Recruiter Interviews
	var recruiterQuery string
	if interviewWhereClause != "" {
		recruiterQuery = "SELECT COUNT(*) FROM interviews" + interviewWhereClause + " AND type = 'RECRUITER'"
	} else {
		recruiterQuery = "SELECT COUNT(*) FROM interviews WHERE type = 'RECRUITER'"
	}
	h.dbConn.QueryRow(recruiterQuery).Scan(&stats.RecruiterInterviews)

	// Query: Manager Interviews
	var managerQuery string
	if interviewWhereClause != "" {
		managerQuery = "SELECT COUNT(*) FROM interviews" + interviewWhereClause + " AND type = 'MANAGER'"
	} else {
		managerQuery = "SELECT COUNT(*) FROM interviews WHERE type = 'MANAGER'"
	}
	h.dbConn.QueryRow(managerQuery).Scan(&stats.ManagerInterviews)

	// Query: Loop Interviews
	var loopQuery string
	if interviewWhereClause != "" {
		loopQuery = "SELECT COUNT(*) FROM interviews" + interviewWhereClause + " AND type = 'LOOP'"
	} else {
		loopQuery = "SELECT COUNT(*) FROM interviews WHERE type = 'LOOP'"
	}
	h.dbConn.QueryRow(loopQuery).Scan(&stats.LoopInterviews)

	// Query: Tech Screen Interviews
	var techScreenQuery string
	if interviewWhereClause != "" {
		techScreenQuery = "SELECT COUNT(*) FROM interviews" + interviewWhereClause + " AND type = 'TECH_SCREEN'"
	} else {
		techScreenQuery = "SELECT COUNT(*) FROM interviews WHERE type = 'TECH_SCREEN'"
	}
	h.dbConn.QueryRow(techScreenQuery).Scan(&stats.TechScreenInterviews)

	templates.Stats(stats).Render(r.Context(), w)
}
