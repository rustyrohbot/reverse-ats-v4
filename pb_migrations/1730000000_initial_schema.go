package pb_migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Create companies collection
		companies := core.NewBaseCollection("companies")

		nameField := &core.TextField{Name: "name", Required: true}
		nameField.Max = 500

		descField := &core.TextField{Name: "description"}
		descField.Max = 50000

		cityField := &core.TextField{Name: "hq_city"}
		cityField.Max = 200

		stateField := &core.TextField{Name: "hq_state"}
		stateField.Max = 100

		companies.Fields.Add(
			nameField,
			descField,
			&core.URLField{Name: "url"},
			&core.URLField{Name: "linkedin"},
			cityField,
			stateField,
		)
		if err := app.Save(companies); err != nil {
			return err
		}

		// Create contacts collection
		contacts := core.NewBaseCollection("contacts")

		firstNameField := &core.TextField{Name: "first_name", Required: true}
		firstNameField.Max = 200

		lastNameField := &core.TextField{Name: "last_name", Required: true}
		lastNameField.Max = 200

		roleField := &core.TextField{Name: "role"}
		roleField.Max = 500

		phoneField := &core.TextField{Name: "phone"}
		phoneField.Max = 50

		contactNotesField := &core.TextField{Name: "notes"}
		contactNotesField.Max = 50000

		contacts.Fields.Add(
			&core.RelationField{
				Name:          "company",
				Required:      true,
				CollectionId:  companies.Id,
				CascadeDelete: true,
				MaxSelect:     1,
			},
			firstNameField,
			lastNameField,
			roleField,
			&core.EmailField{Name: "email"},
			phoneField,
			&core.URLField{Name: "linkedin"},
			contactNotesField,
		)
		if err := app.Save(contacts); err != nil {
			return err
		}

		// Create roles collection
		roles := core.NewBaseCollection("roles")

		roleNameField := &core.TextField{Name: "name", Required: true}
		roleNameField.Max = 500

		roleDescField := &core.TextField{Name: "description"}
		roleDescField.Max = 50000

		coverLetterField := &core.TextField{Name: "cover_letter"}
		coverLetterField.Max = 50000

		appLocationField := &core.TextField{Name: "application_location"}
		appLocationField.Max = 1000

		workCityField := &core.TextField{Name: "work_city"}
		workCityField.Max = 200

		workStateField := &core.TextField{Name: "work_state"}
		workStateField.Max = 100

		locationField := &core.TextField{Name: "location"}
		locationField.Max = 100

		statusField := &core.TextField{Name: "status"}
		statusField.Max = 100

		discoveryField := &core.TextField{Name: "discovery"}
		discoveryField.Max = 500

		roleNotesField := &core.TextField{Name: "notes"}
		roleNotesField.Max = 50000

		roles.Fields.Add(
			&core.RelationField{
				Name:          "company",
				Required:      true,
				CollectionId:  companies.Id,
				CascadeDelete: true,
				MaxSelect:     1,
			},
			roleNameField,
			&core.URLField{Name: "url"},
			roleDescField,
			coverLetterField,
			appLocationField,
			&core.DateField{Name: "applied_date"},
			&core.DateField{Name: "closed_date"},
			&core.NumberField{Name: "posted_range_min"},
			&core.NumberField{Name: "posted_range_max"},
			&core.BoolField{Name: "equity"},
			workCityField,
			workStateField,
			locationField,
			statusField,
			discoveryField,
			&core.BoolField{Name: "referral"},
			roleNotesField,
		)
		if err := app.Save(roles); err != nil {
			return err
		}

		// Create interviews collection
		interviews := core.NewBaseCollection("interviews")

		startField := &core.TextField{Name: "start", Required: true}
		startField.Max = 20

		endField := &core.TextField{Name: "end", Required: true}
		endField.Max = 20

		interviewNotesField := &core.TextField{Name: "notes"}
		interviewNotesField.Max = 50000

		interviews.Fields.Add(
			&core.RelationField{
				Name:          "role",
				Required:      true,
				CollectionId:  roles.Id,
				CascadeDelete: true,
				MaxSelect:     1,
			},
			&core.DateField{Name: "date", Required: true},
			startField,
			endField,
			interviewNotesField,
			&core.SelectField{
				Name:      "type",
				Required:  true,
				MaxSelect: 1,
				Values:    []string{"RECRUITER", "LOOP", "TECH_SCREEN", "MANAGER", "MISC"},
			},
			&core.RelationField{
				Name:          "contacts",
				CollectionId:  contacts.Id,
				CascadeDelete: false,
				// MaxSelect omitted for unlimited (many-to-many)
			},
		)
		if err := app.Save(interviews); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		// Down migration - drop all collections
		collections := []string{"interviews", "roles", "contacts", "companies"}
		for _, name := range collections {
			collection, err := app.FindCollectionByNameOrId(name)
			if err == nil {
				if err := app.Delete(collection); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
