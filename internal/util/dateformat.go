package util

import "time"

// FormatDateToText converts ISO format dates to text format for display.
// "2025-10-01" → "October 1, 2025"
// "October 1, 2025" → "October 1, 2025" (unchanged)
func FormatDateToText(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// Try to parse as ISO format first
	t, err := time.Parse("2006-01-02", dateStr)
	if err == nil {
		// Successfully parsed as ISO, convert to text format
		return t.Format("January 2, 2006")
	}

	// Not ISO format, return as-is (assume it's already in text format)
	return dateStr
}
