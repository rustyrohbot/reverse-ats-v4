package models

// Role represents a job role/position
type Role struct {
	ID                  string
	CompanyID           string
	CompanyName         string // For display purposes
	Name                string
	Url                 string
	Description         string
	CoverLetter         string
	ApplicationLocation string
	AppliedDate         string
	ClosedDate          string
	PostedRangeMin      int64
	PostedRangeMax      int64
	Equity              bool
	WorkCity            string
	WorkState           string
	Location            string
	Status              string
	Discovery           string
	Referral            bool
	Notes               string
	CreatedAt           string
	UpdatedAt           string
}
