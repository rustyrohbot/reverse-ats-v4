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
