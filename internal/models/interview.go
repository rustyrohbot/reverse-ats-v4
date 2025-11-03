package models

// Interview represents an interview for a role
type Interview struct {
	ID          string
	RoleID      string
	RoleName    string   // For display purposes
	CompanyID   string   // For display purposes
	CompanyName string   // For display purposes
	Date        string
	Start       string
	End         string
	Notes       string
	Type        string
	ContactIDs  []string // Multiple contacts can be associated
	Contacts    []Contact // Expanded contacts for display
	CreatedAt   string
	UpdatedAt   string
}
