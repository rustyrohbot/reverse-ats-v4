package util

import "time"

// FormatDateToText converts ISO format dates to text format for display.
// "2025-10-01" → "October 1, 2025"
// "2025-10-20 00:00:00.000Z" → "October 20, 2025"
// "October 1, 2025" → "October 1, 2025" (unchanged)
func FormatDateToText(dateStr string) string {
	if dateStr == "" {
		return ""
	}

	// Try multiple date/datetime formats
	formats := []string{
		"2006-01-02",                        // ISO date: 2025-10-20
		"2006-01-02 15:04:05.999999999Z07:00", // ISO datetime with nanoseconds: 2025-10-20 00:00:00.000Z
		"2006-01-02 15:04:05.999999999Z",    // ISO datetime with nanoseconds: 2025-10-20 00:00:00.000Z
		"2006-01-02T15:04:05.999999999Z",    // ISO8601 datetime: 2025-10-20T00:00:00.000Z
		"2006-01-02 15:04:05",               // datetime without timezone: 2025-10-20 00:00:00
		time.RFC3339,                        // RFC3339 format
		time.RFC3339Nano,                    // RFC3339 with nanoseconds
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			// Successfully parsed, convert to text format
			return t.Format("January 2, 2006")
		}
	}

	// Not a recognized format, return as-is (assume it's already in text format)
	return dateStr
}

// FormatTimeTo12Hour converts 24-hour time format to 12-hour format for display.
// "14:30" → "2:30 PM"
// "09:00" → "9:00 AM"
// "12:00 PM" → "12:00 PM" (unchanged if already 12-hour)
func FormatTimeTo12Hour(timeStr string) string {
	if timeStr == "" {
		return ""
	}

	// Try to parse as 24-hour format
	formats := []string{
		"15:04",    // 24-hour HH:MM
		"3:04 PM",  // Already 12-hour
		"03:04 PM", // Already 12-hour with leading zero
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			// Return in 12-hour format without leading zero
			return t.Format("3:04 PM")
		}
	}

	// If all parsing fails, return original
	return timeStr
}
