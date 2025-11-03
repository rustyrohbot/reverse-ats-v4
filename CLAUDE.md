# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A web application for tracking job applications, interviews, and contacts during a job search - replacing a spreadsheet-based workflow with a user-friendly CRUD interface.

**Tech Stack:**
- **Backend**: Go with PocketBase v0.31.0 (as Go package)
- **Frontend**: HTMX for interactivity
- **Templating**: templ (type-safe Go templates)
- **Styling**: Tailwind CSS
- **Database**: PocketBase (SQLite with built-in type-safe API)

## Initial Data Schema

The CSV files in this repository represent the initial data export from a spreadsheet. These define the database schema:

### Tables and Relationships

**Companies** (`reverse-ats - Companies.csv`)
- PK: `companyID`
- Columns: name, description, url, hqCity, hqState

**Roles** (`reverse-ats - Roles.csv`)
- PK: `roleID`
- FK: `companyID` → Companies
- Columns: name, url, description, coverLetter, applicationLocation, appliedDate, closedDate, postedRangeMin, postedRangeMax, equity, workCity, workState, location, status, discovery, referral, notes

**Interviews** (`reverse-ats - Interviews.csv`)
- PK: `interviewID`
- FK: `roleID` → Roles
- Columns: date, start, end, notes, type
- Types: RECRUITER, LOOP, TECH_SCREEN, MANAGER, MISC

**Contacts** (`reverse-ats - Contacts.csv`)
- PK: `contactID`
- FK: `companyID` → Companies
- Columns: firstName, lastName, role, email, phone, linkedin, notes

**InterviewsContacts** (`reverse-ats - InterviewsContacts.csv`)
- Junction table (many-to-many)
- FK: `interviewId` → Interviews
- FK: `contactId` → Contacts

### Data Notes
- CSV files use the string "NULL" (all caps) for missing values
- Date fields in CSV files must be in YYYY-MM-DD format (e.g., "2025-01-15")
- Date fields are displayed in the UI as text format (e.g., "January 15, 2025")
- Time fields use 24-hour HH:MM format (e.g., "14:30")
- Import/export maintains consistent NULL formatting for round-trip compatibility

## Development Commands

### Database Setup
```bash
# Reset database (delete data directory)
rm -rf pb_data

# PocketBase auto-migrates on startup using pb_migrations/
# Create new migrations by editing pb_migrations/ files

# Access PocketBase admin UI
# Navigate to http://localhost:5627/_/
```

### Running the Application
```bash
# Run directly
go run cmd/server/main.go

# Or using make
make run
```

### Building
```bash
# Generate templ templates (run this after changing .templ files)
make generate
# or manually:
templ generate

# Build CSS with Tailwind
npx tailwindcss -i ./static/input.css -o ./static/output.css --watch
# or using npm scripts:
npm run watch:css

# Build application
go build -o bin/server cmd/server/main.go
# or using make:
make build
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/handlers
```

## Architecture Guidelines

### Project Structure
```
├── cmd/
│   ├── server/          # Main application entry point
│   ├── import/          # CLI import utility
│   └── export/          # CLI export utility
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Domain models
│   ├── importer/        # Shared CSV import logic
│   ├── exporter/        # Shared CSV export logic
│   ├── util/            # Shared utilities (date formatting, etc.)
│   └── templates/       # templ template files
├── pb_migrations/       # PocketBase schema migrations
├── pb_data/             # PocketBase data directory (gitignored)
├── static/              # CSS, JS, and static assets
├── import/              # CSV files for import (gitignored)
├── export/              # Exported CSV files (gitignored)
└── sample_data/         # Sample CSV files for reference
```

### HTMX Patterns
- Use `hx-get`, `hx-post`, `hx-put`, `hx-delete` for CRUD operations
- Return partial HTML fragments for targeted updates
- Use `hx-target` and `hx-swap` for DOM manipulation
- Leverage `hx-trigger` for interactive events

### Database Access with PocketBase
- PocketBase provides a type-safe API for database operations
- Collections are defined in `pb_migrations/` files
- Access records via PocketBase app instance:
  ```go
  // Find collection
  collection, err := app.FindCollectionByNameOrId("companies")

  // Create record
  record := core.NewRecord(collection)
  record.Set("name", "Company Name")
  record.Set("url", "https://example.com")
  err := app.Save(record)

  // Find record
  record, err := app.FindRecordById("companies", recordId)

  // Find all records with filter
  records, err := app.FindRecordsByFilter("companies", "name ~ 'Tech'", "-created", 100)

  // Update record
  record.Set("name", "Updated Name")
  err := app.Save(record)

  // Delete record
  err := app.Delete(record)
  ```
- PocketBase handles NULL values automatically with empty strings
- Use relation fields for foreign keys (automatic cascade delete support)

### Template Rendering
- Use templ for type-safe template generation
- Run `templ generate` after modifying `.templ` files
- Keep template logic minimal; put business logic in handlers/models

### Styling
- Use Tailwind utility classes
- Run Tailwind CLI in watch mode during development
- Consider using Tailwind forms plugin for better form styling
- **IMPORTANT**: After modifying `.templ` files, always rebuild Tailwind CSS:
  ```bash
  npx tailwindcss -i ./static/input.css -o ./static/output.css --minify
  ```

### Code Organization and Shared Logic
- Shared logic between CLI and web handlers should be extracted to separate packages
- Example: `internal/importer` and `internal/exporter` contain shared CSV processing logic
- Utility functions (like date formatting) belong in `internal/util`
- This ensures consistency and reduces code duplication

### Import/Export Architecture
- **Import Flow**:
  - Web handler (`internal/handlers/import.go`) accepts file uploads
  - CLI (`cmd/import/main.go`) reads from `./import` directory
  - Both use shared logic from `internal/importer/importer.go`
  - Import order is enforced: Companies → Roles → Contacts → Interviews → InterviewsContacts

- **Export Flow**:
  - Web handler (`internal/handlers/export.go`) creates a zip file
  - CLI (`cmd/export/main.go`) writes to `./export` directory
  - Both use shared logic from `internal/exporter/exporter.go`
  - NULL values are consistently written as "NULL" string

### Date Handling
- PocketBase DateField stores dates in ISO format (YYYY-MM-DD)
- CSV files must use ISO format (YYYY-MM-DD)
- The importer converts text format dates like "April 14, 2025" to ISO format automatically
- UI displays dates in text format using `internal/util/dateformat.go`
- PocketBase handles date filtering and sorting natively
