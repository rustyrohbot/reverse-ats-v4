package importer

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// ID mapping tables to handle conversion from old integer IDs to PocketBase string IDs
type IDMappings struct {
	Companies map[string]string // old ID -> new ID
	Roles     map[string]string
	Contacts  map[string]string
	Interviews map[string]string
}

func NewIDMappings() *IDMappings {
	return &IDMappings{
		Companies: make(map[string]string),
		Roles:     make(map[string]string),
		Contacts:  make(map[string]string),
		Interviews: make(map[string]string),
	}
}

// emptyToNull converts empty string or "NULL" to empty string for PocketBase
func emptyToNull(s string) string {
	if s == "NULL" {
		return ""
	}
	return s
}

// parseBool converts string to boolean
func parseBool(s string) bool {
	if s == "" || s == "NULL" {
		return false
	}
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "true" || s == "yes" || s == "1" || s == "t" || s == "y"
}

// parseInt64 parses int64 from string, handling currency format
func parseInt64(s string) int64 {
	if s == "" || s == "NULL" {
		return 0
	}

	// Remove currency symbols, commas, and decimal portions
	s = strings.TrimPrefix(s, "$")
	s = strings.ReplaceAll(s, ",", "")

	// If there's a decimal point, take only the integer part
	if idx := strings.Index(s, "."); idx != -1 {
		s = s[:idx]
	}

	var result int64
	fmt.Sscanf(s, "%d", &result)
	return result
}

// parseDate converts date string to ISO format (YYYY-MM-DD)
// Handles both "April 14, 2025" and "2025-04-14" formats
func parseDate(s string) string {
	if s == "" || s == "NULL" {
		return ""
	}

	// Try ISO format first (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.Format("2006-01-02")
	}

	// Try text format (January 2, 2006)
	if t, err := time.Parse("January 2, 2006", s); err == nil {
		return t.Format("2006-01-02")
	}

	// Return as-is if can't parse (will likely fail validation)
	return s
}

// ImportCompanies imports companies from CSV file
func ImportCompanies(app *pocketbase.PocketBase, filepath string, mappings *IDMappings) error {
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

	collection, err := app.FindCollectionByNameOrId("companies")
	if err != nil {
		return fmt.Errorf("failed to find companies collection: %w", err)
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
		oldID := record[0]

		pbRecord := core.NewRecord(collection)
		pbRecord.Set("name", record[1])
		pbRecord.Set("description", emptyToNull(record[2]))
		pbRecord.Set("url", emptyToNull(record[3]))
		pbRecord.Set("linkedin", emptyToNull(record[4]))
		pbRecord.Set("hq_city", emptyToNull(record[5]))
		pbRecord.Set("hq_state", emptyToNull(record[6]))

		if err := app.Save(pbRecord); err != nil {
			return fmt.Errorf("failed to insert company %s: %w", record[1], err)
		}

		// Store ID mapping
		mappings.Companies[oldID] = pbRecord.Id
		count++
	}

	fmt.Printf("Imported %d companies\n", count)
	return nil
}

// ImportContacts imports contacts from CSV file
func ImportContacts(app *pocketbase.PocketBase, filepath string, mappings *IDMappings) error {
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

	collection, err := app.FindCollectionByNameOrId("contacts")
	if err != nil {
		return fmt.Errorf("failed to find contacts collection: %w", err)
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
		oldID := record[0]
		oldCompanyID := record[1]

		// Map old company ID to new PocketBase ID
		newCompanyID, ok := mappings.Companies[oldCompanyID]
		if !ok {
			return fmt.Errorf("company ID %s not found in mapping for contact %s %s", oldCompanyID, record[2], record[3])
		}

		pbRecord := core.NewRecord(collection)
		pbRecord.Set("company", newCompanyID)
		pbRecord.Set("first_name", record[2])
		pbRecord.Set("last_name", record[3])
		pbRecord.Set("role", emptyToNull(record[4]))
		pbRecord.Set("email", emptyToNull(record[5]))
		pbRecord.Set("phone", emptyToNull(record[6]))
		pbRecord.Set("linkedin", emptyToNull(record[7]))
		pbRecord.Set("notes", emptyToNull(record[8]))

		if err := app.Save(pbRecord); err != nil {
			return fmt.Errorf("failed to insert contact %s %s: %w", record[2], record[3], err)
		}

		// Store ID mapping
		mappings.Contacts[oldID] = pbRecord.Id
		count++
	}

	fmt.Printf("Imported %d contacts\n", count)
	return nil
}

// ImportRoles imports roles from CSV file
func ImportRoles(app *pocketbase.PocketBase, filepath string, mappings *IDMappings) error {
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

	collection, err := app.FindCollectionByNameOrId("roles")
	if err != nil {
		return fmt.Errorf("failed to find roles collection: %w", err)
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
		if record[1] == "" || record[1] == "NULL" {
			continue
		}

		oldID := record[0]
		oldCompanyID := record[1]

		// Map old company ID to new PocketBase ID
		newCompanyID, ok := mappings.Companies[oldCompanyID]
		if !ok {
			return fmt.Errorf("company ID %s not found in mapping for role %s", oldCompanyID, record[2])
		}

		pbRecord := core.NewRecord(collection)
		pbRecord.Set("company", newCompanyID)
		pbRecord.Set("name", record[2])
		pbRecord.Set("url", emptyToNull(record[3]))
		pbRecord.Set("description", emptyToNull(record[4]))
		pbRecord.Set("cover_letter", emptyToNull(record[5]))
		pbRecord.Set("application_location", emptyToNull(record[6]))
		pbRecord.Set("applied_date", parseDate(record[7]))
		pbRecord.Set("closed_date", parseDate(record[8]))
		pbRecord.Set("posted_range_min", parseInt64(record[9]))
		pbRecord.Set("posted_range_max", parseInt64(record[10]))
		pbRecord.Set("equity", parseBool(record[11]))
		pbRecord.Set("work_city", emptyToNull(record[12]))
		pbRecord.Set("work_state", emptyToNull(record[13]))
		pbRecord.Set("location", emptyToNull(record[14]))
		pbRecord.Set("status", emptyToNull(record[15]))
		pbRecord.Set("discovery", emptyToNull(record[16]))
		pbRecord.Set("referral", parseBool(record[17]))
		pbRecord.Set("notes", emptyToNull(record[18]))

		if err := app.Save(pbRecord); err != nil {
			return fmt.Errorf("failed to insert role %s: %w", record[2], err)
		}

		// Store ID mapping
		mappings.Roles[oldID] = pbRecord.Id
		count++
	}

	fmt.Printf("Imported %d roles\n", count)
	return nil
}

// ImportInterviews imports interviews from CSV file
func ImportInterviews(app *pocketbase.PocketBase, filepath string, mappings *IDMappings) error {
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

	collection, err := app.FindCollectionByNameOrId("interviews")
	if err != nil {
		return fmt.Errorf("failed to find interviews collection: %w", err)
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
		oldID := record[0]
		oldRoleID := record[1]

		// Map old role ID to new PocketBase ID
		newRoleID, ok := mappings.Roles[oldRoleID]
		if !ok {
			return fmt.Errorf("role ID %s not found in mapping for interview on %s", oldRoleID, record[2])
		}

		pbRecord := core.NewRecord(collection)
		pbRecord.Set("role", newRoleID)
		pbRecord.Set("date", parseDate(record[2]))
		pbRecord.Set("start", record[3])
		pbRecord.Set("end", record[4])
		pbRecord.Set("notes", emptyToNull(record[5]))
		pbRecord.Set("type", record[6])

		if err := app.Save(pbRecord); err != nil {
			return fmt.Errorf("failed to insert interview on %s: %w", record[2], err)
		}

		// Store ID mapping
		mappings.Interviews[oldID] = pbRecord.Id
		count++
	}

	fmt.Printf("Imported %d interviews\n", count)
	return nil
}

// ImportInterviewsContacts imports interview-contact relationships from CSV file
func ImportInterviewsContacts(app *pocketbase.PocketBase, filepath string, mappings *IDMappings) error {
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
		oldInterviewID := record[1]
		oldContactID := record[2]

		// Map old IDs to new PocketBase IDs
		newInterviewID, ok := mappings.Interviews[oldInterviewID]
		if !ok {
			return fmt.Errorf("interview ID %s not found in mapping", oldInterviewID)
		}

		newContactID, ok := mappings.Contacts[oldContactID]
		if !ok {
			return fmt.Errorf("contact ID %s not found in mapping", oldContactID)
		}

		// Fetch the interview record and update its contacts field
		interview, err := app.FindRecordById("interviews", newInterviewID)
		if err != nil {
			return fmt.Errorf("failed to find interview %s: %w", newInterviewID, err)
		}

		// Get existing contacts (if any)
		existingContacts := interview.GetStringSlice("contacts")
		// Add new contact if not already present
		found := false
		for _, c := range existingContacts {
			if c == newContactID {
				found = true
				break
			}
		}
		if !found {
			existingContacts = append(existingContacts, newContactID)
			interview.Set("contacts", existingContacts)
			if err := app.Save(interview); err != nil {
				return fmt.Errorf("failed to update interview %s with contact %s: %w", newInterviewID, newContactID, err)
			}
		}

		count++
	}

	fmt.Printf("Imported %d interview-contact links\n", count)
	return nil
}

// ImportStep represents a single import operation
type ImportStep struct {
	Name     string
	Filename string
	Filepath string
	Fn       func(*pocketbase.PocketBase, string, *IDMappings) error
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
func ImportFromSteps(app *pocketbase.PocketBase, steps []ImportStep, skipMissing bool) []error {
	var errors []error
	mappings := NewIDMappings()

	for _, step := range steps {
		if step.Filepath == "" {
			if skipMissing {
				continue
			}
			err := fmt.Errorf("%s: no file provided", step.Name)
			fmt.Printf("ERROR: %v\n", err)
			errors = append(errors, err)
			continue
		}

		// Check if file exists
		if _, err := os.Stat(step.Filepath); os.IsNotExist(err) {
			if skipMissing {
				continue
			}
			err := fmt.Errorf("%s: file not found: %s", step.Name, step.Filepath)
			fmt.Printf("ERROR: %v\n", err)
			errors = append(errors, err)
			continue
		}

		fmt.Printf("Importing %s from %s...\n", step.Name, step.Filepath)
		if err := step.Fn(app, step.Filepath, mappings); err != nil {
			wrappedErr := fmt.Errorf("%s: %w", step.Name, err)
			fmt.Printf("ERROR: %v\n", wrappedErr)
			errors = append(errors, wrappedErr)
			continue
		}
	}

	return errors
}

func ImportAll(app *pocketbase.PocketBase, dir string) error {
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
	errors := ImportFromSteps(app, steps, false)
	if len(errors) > 0 {
		return errors[0]
	}

	fmt.Println("\n✅ All data imported successfully!")
	return nil
}
