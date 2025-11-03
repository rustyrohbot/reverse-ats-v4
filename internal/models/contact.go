package models

// Contact represents a person at a company
type Contact struct {
	ID          string
	CompanyID   string
	CompanyName string // For display purposes
	FirstName   string
	LastName    string
	Role        string
	Email       string
	Phone       string
	Linkedin    string
	Notes       string
	CreatedAt   string
	UpdatedAt   string
}
