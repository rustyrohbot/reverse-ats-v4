package handlers

import (
	"database/sql"
	"strconv"
)

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func nullInt64(s string) sql.NullInt64 {
	if s == "" {
		return sql.NullInt64{Valid: false}
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: val, Valid: true}
}

func nullBool(s string) sql.NullBool {
	// Checkbox is only present in form value if checked, value will be "true"
	if s == "true" {
		return sql.NullBool{Bool: true, Valid: true}
	}
	// If empty or any other value, treat as false but valid
	if s == "" {
		return sql.NullBool{Bool: false, Valid: false}
	}
	return sql.NullBool{Bool: false, Valid: true}
}
