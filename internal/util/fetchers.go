package util

import (
	"reverse-ats/internal/models"

	"github.com/pocketbase/pocketbase"
)

// FetchCompaniesForDropdown fetches all companies and converts them to models.Company
// This is used by handlers that need to populate company dropdowns
func FetchCompaniesForDropdown(app *pocketbase.PocketBase) ([]models.Company, error) {
	// Fetch all companies sorted by name
	companyRecords, err := app.FindRecordsByFilter(CollectionCompanies, "", "name", -1, 0)
	if err != nil {
		return nil, err
	}

	companies := make([]models.Company, len(companyRecords))
	for i, record := range companyRecords {
		companies[i] = models.Company{
			ID:          record.Id,
			Name:        record.GetString("name"),
			Description: record.GetString("description"),
			Url:         record.GetString("url"),
			Linkedin:    record.GetString("linkedin"),
			HqCity:      record.GetString("hqCity"),
			HqState:     record.GetString("hqState"),
		}
	}

	return companies, nil
}

// FetchCompaniesMap fetches all companies and returns them as a map[id]name
// This is used to avoid N+1 queries when populating company names in lists
func FetchCompaniesMap(app *pocketbase.PocketBase) (map[string]string, error) {
	companyRecords, err := app.FindRecordsByFilter(CollectionCompanies, "", "", -1, 0)
	if err != nil {
		return nil, err
	}

	companiesMap := make(map[string]string, len(companyRecords))
	for _, record := range companyRecords {
		companiesMap[record.Id] = record.GetString("name")
	}

	return companiesMap, nil
}

// RoleInfo holds basic role information for map lookups
type RoleInfo struct {
	Name      string
	CompanyID string
}

// FetchRolesMap fetches all roles and returns them as a map[id]RoleInfo
// This is used to avoid N+1 queries when populating role names in lists
func FetchRolesMap(app *pocketbase.PocketBase) (map[string]RoleInfo, error) {
	roleRecords, err := app.FindRecordsByFilter(CollectionRoles, "", "", -1, 0)
	if err != nil {
		return nil, err
	}

	rolesMap := make(map[string]RoleInfo, len(roleRecords))
	for _, record := range roleRecords {
		rolesMap[record.Id] = RoleInfo{
			Name:      record.GetString("name"),
			CompanyID: record.GetString("company"),
		}
	}

	return rolesMap, nil
}
