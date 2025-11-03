package importer

import (
	"context"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"reverse-ats/internal/db"
)

// nullString converts "NULL" string to sql.NullString
func nullString(s string) sql.NullString {
	if s == "" || s == "NULL" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// NullInt64 converts string to sql.NullInt64
// Handles currency format like "$130,047.00"
func NullInt64(s string) sql.NullInt64 {
	if s == "" || s == "NULL" {
		return sql.NullInt64{Valid: false}
	}

	// Remove currency symbols, commas, and decimal portions
	s = strings.TrimPrefix(s, "$")
	s = strings.ReplaceAll(s, ",", "")

	// If there's a decimal point, take only the integer part
	if idx := strings.Index(s, "."); idx != -1 {
		s = s[:idx]
	}

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: i, Valid: true}
}

// nullBool converts string to sql.NullBool
func nullBool(s string) sql.NullBool {
	if s == "" || s == "NULL" {
		return sql.NullBool{Valid: false}
	}

	// Handle various boolean representations
	s = strings.ToLower(strings.TrimSpace(s))
	switch s {
	case "true", "yes", "1", "t", "y":
		return sql.NullBool{Bool: true, Valid: true}
	case "false", "no", "0", "f", "n":
		return sql.NullBool{Bool: false, Valid: true}
	default:
		return sql.NullBool{Valid: false}
	}
}

// ImportCompanies imports companies from CSV file
func ImportCompanies(queries *db.Queries, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open companies CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		// Skip if not enough fields
		if len(record) < 7 {
			continue
		}

		// CSV columns: companyID,name,description,url,linkedin,hqCity,hqState
		_, err = queries.CreateCompany(context.Background(), db.CreateCompanyParams{
			Name:        record[1],
			Description: nullString(record[2]),
			Url:         nullString(record[3]),
			Linkedin:    nullString(record[4]),
			HqCity:      nullString(record[5]),
			HqState:     nullString(record[6]),
		})
		if err != nil {
			return fmt.Errorf("failed to insert company %s: %w", record[1], err)
		}
		count++
	}

	fmt.Printf("Imported %d companies\n", count)
	return nil
}

// ImportContacts imports contacts from CSV file
func ImportContacts(queries *db.Queries, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open contacts CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		// Skip if not enough fields
		if len(record) < 9 {
			continue
		}

		// CSV columns: contactID,companyID,firstName,lastName,role,email,phone,linkedin,notes
		companyID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid company ID %s: %w", record[1], err)
		}

		_, err = queries.CreateContact(context.Background(), db.CreateContactParams{
			CompanyID: companyID,
			FirstName: record[2],
			LastName:  record[3],
			Role:      nullString(record[4]),
			Email:     nullString(record[5]),
			Phone:     nullString(record[6]),
			Linkedin:  nullString(record[7]),
			Notes:     nullString(record[8]),
		})
		if err != nil {
			return fmt.Errorf("failed to insert contact %s %s: %w", record[2], record[3], err)
		}
		count++
	}

	fmt.Printf("Imported %d contacts\n", count)
	return nil
}

// ImportRoles imports roles from CSV file
func ImportRoles(queries *db.Queries, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open roles CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		// Skip if not enough fields
		if len(record) < 19 {
			continue
		}

		// CSV columns: roleID,companyID,name,url,description,coverLetter,applicationLocation,
		// appliedDate,closedDate,postedRangeMin,postedRangeMax,equity,workCity,workState,
		// location,status,discovery,referral,notes

		// Skip rows with empty company ID (malformed CSV)
		if record[1] == "" {
			continue
		}

		companyID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid company ID %s: %w", record[1], err)
		}

		_, err = queries.CreateRole(context.Background(), db.CreateRoleParams{
			CompanyID:           companyID,
			Name:                record[2],
			Url:                 nullString(record[3]),
			Description:         nullString(record[4]),
			CoverLetter:         nullString(record[5]),
			ApplicationLocation: nullString(record[6]),
			AppliedDate:         nullString(record[7]),
			ClosedDate:          nullString(record[8]),
			PostedRangeMin:      NullInt64(record[9]),
			PostedRangeMax:      NullInt64(record[10]),
			Equity:              nullBool(record[11]),
			WorkCity:            nullString(record[12]),
			WorkState:           nullString(record[13]),
			Location:            nullString(record[14]),
			Status:              nullString(record[15]),
			Discovery:           nullString(record[16]),
			Referral:            nullBool(record[17]),
			Notes:               nullString(record[18]),
		})
		if err != nil {
			return fmt.Errorf("failed to insert role %s: %w", record[2], err)
		}
		count++
	}

	fmt.Printf("Imported %d roles\n", count)
	return nil
}

// ImportInterviews imports interviews from CSV file
func ImportInterviews(queries *db.Queries, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open interviews CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		// Skip if not enough fields
		if len(record) < 7 {
			continue
		}

		// CSV columns: interviewID,roleID,date,start,end,notes,type
		roleID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid role ID %s: %w", record[1], err)
		}

		_, err = queries.CreateInterview(context.Background(), db.CreateInterviewParams{
			RoleID: roleID,
			Date:   record[2],
			Start:  record[3],
			End:    record[4],
			Notes:  nullString(record[5]),
			Type:   record[6],
		})
		if err != nil {
			return fmt.Errorf("failed to insert interview on %s: %w", record[2], err)
		}
		count++
	}

	fmt.Printf("Imported %d interviews\n", count)
	return nil
}

// ImportInterviewsContacts imports interview-contact relationships from CSV file
func ImportInterviewsContacts(queries *db.Queries, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open interviews-contacts CSV: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Skip header
	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read CSV record: %w", err)
		}

		// Skip if not enough fields
		if len(record) < 3 {
			continue
		}

		// CSV columns: interviewsContactId,interviewId,contactId
		interviewID, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid interview ID %s: %w", record[1], err)
		}

		contactID, err := strconv.ParseInt(record[2], 10, 64)
		if err != nil {
			return fmt.Errorf("invalid contact ID %s: %w", record[2], err)
		}

		err = queries.LinkInterviewContact(context.Background(), db.LinkInterviewContactParams{
			InterviewID: interviewID,
			ContactID:   contactID,
		})
		if err != nil {
			return fmt.Errorf("failed to link interview %d with contact %d: %w", interviewID, contactID, err)
		}
		count++
	}

	fmt.Printf("Imported %d interview-contact links\n", count)
	return nil
}

// ImportAll imports all CSV files from the specified directory
// ImportStep represents a single import operation
type ImportStep struct {
	Name     string
	Filename string
	Filepath string
	Fn       func(*db.Queries, string) error
}

// GetImportSteps returns the import steps in the correct order (respecting foreign keys)
func GetImportSteps() []ImportStep {
	return []ImportStep{
		{"companies", "reverse-ats - Companies.csv", "", ImportCompanies},
		{"roles", "reverse-ats - Roles.csv", "", ImportRoles},
		{"contacts", "reverse-ats - Contacts.csv", "", ImportContacts},
		{"interviews", "reverse-ats - Interviews.csv", "", ImportInterviews},
		{"interviews-contacts", "reverse-ats - InterviewsContacts.csv", "", ImportInterviewsContacts},
	}
}

// ImportFromSteps imports data from a list of steps, skipping missing files
// Returns a list of errors encountered (doesn't stop on first error)
func ImportFromSteps(queries *db.Queries, steps []ImportStep, skipMissing bool) []error {
	var errors []error

	for _, step := range steps {
		if step.Filepath == "" {
			if skipMissing {
				continue
			}
			errors = append(errors, fmt.Errorf("%s: no file provided", step.Name))
			continue
		}

		// Check if file exists
		if _, err := os.Stat(step.Filepath); os.IsNotExist(err) {
			if skipMissing {
				continue
			}
			errors = append(errors, fmt.Errorf("%s: file not found: %s", step.Name, step.Filepath))
			continue
		}

		fmt.Printf("Importing %s from %s...\n", step.Name, step.Filepath)
		if err := step.Fn(queries, step.Filepath); err != nil {
			errors = append(errors, fmt.Errorf("%s: %w", step.Name, err))
			continue
		}
	}

	return errors
}

func ImportAll(queries *db.Queries, dir string) error {
	steps := GetImportSteps()

	// Set filepaths for all steps
	for i := range steps {
		filepath := dir + "/" + steps[i].Filename
		// Check if file exists
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			// Try without space in filename
			filepath = strings.ReplaceAll(filepath, " ", "")
			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				return fmt.Errorf("file not found: %s", steps[i].Filename)
			}
		}
		steps[i].Filepath = filepath
	}

	// Import with error collection (but fail on first error for CLI compatibility)
	errors := ImportFromSteps(queries, steps, false)
	if len(errors) > 0 {
		return errors[0]
	}

	fmt.Println("\n✅ All data imported successfully!")
	return nil
}
