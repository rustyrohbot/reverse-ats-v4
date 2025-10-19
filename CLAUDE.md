# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A web application for tracking job applications, interviews, and contacts during a job search - replacing a spreadsheet-based workflow with a user-friendly CRUD interface.

**Tech Stack:**
- **Backend**: Go (standard library + `net/http`)
- **Frontend**: HTMX for interactivity
- **Templating**: templ (type-safe Go templates)
- **Styling**: Tailwind CSS
- **Database**: SQLite
- **Migrations**: goose (github.com/pressly/goose)
- **SQL Code Generation**: sqlc (github.com/sqlc-dev/sqlc)

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
- CSV files use the string "NULL" for missing values
- Date fields are formatted as "April 14, 2025" (text format)
- Time fields use 12-hour format with AM/PM

## Development Commands

### Database Setup
```bash
# Install goose CLI (optional, can use programmatically)
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir migrations sqlite3 ./data.db up

# Create a new migration
goose -dir migrations create migration_name sql

# Reset database
goose -dir migrations sqlite3 ./data.db down
# or
rm -f data.db && goose -dir migrations sqlite3 ./data.db up

# Check migration status
goose -dir migrations sqlite3 ./data.db status
```

### Running the Application
```bash
# Run development server with hot reload (if using air)
air

# Run directly
go run cmd/server/main.go
```

### Building
```bash
# Generate all code (run this after changing SQL queries or templates)
make generate
# or manually:
sqlc generate          # Generate Go code from SQL queries
templ generate         # Generate Go code from templ templates

# Build CSS with Tailwind
npx tailwindcss -i ./static/input.css -o ./static/output.css --watch

# Build application
go build -o bin/server cmd/server/main.go
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
│   └── server/          # Main application entry point
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── db/              # sqlc generated database code
│   ├── database/        # Database connection setup
│   └── templates/       # templ template files
├── migrations/          # goose SQL migration files
├── queries/             # SQL queries for sqlc
├── static/              # CSS, JS, and static assets
├── sqlc.yaml            # sqlc configuration
└── *.csv               # Initial data files
```

### HTMX Patterns
- Use `hx-get`, `hx-post`, `hx-put`, `hx-delete` for CRUD operations
- Return partial HTML fragments for targeted updates
- Use `hx-target` and `hx-swap` for DOM manipulation
- Leverage `hx-trigger` for interactive events

### Database Access with sqlc
- Write SQL queries in `queries/*.sql` files
- Use sqlc annotations for type-safe query generation:
  ```sql
  -- name: GetCompany :one
  SELECT * FROM companies WHERE company_id = ?;

  -- name: ListCompanies :many
  SELECT * FROM companies ORDER BY name;

  -- name: CreateCompany :execresult
  INSERT INTO companies (name, url, hq_city, hq_state)
  VALUES (?, ?, ?, ?);
  ```
- Run `sqlc generate` to create type-safe Go code in `internal/db/`
- Use generated `Queries` struct in handlers:
  ```go
  queries := db.New(dbConn)
  company, err := queries.GetCompany(ctx, companyID)
  ```
- Handle NULL values using `sql.Null*` types or custom nullable types

### Template Rendering
- Use templ for type-safe template generation
- Run `templ generate` after modifying `.templ` files
- Keep template logic minimal; put business logic in handlers/models

### Styling
- Use Tailwind utility classes
- Run Tailwind CLI in watch mode during development
- Consider using Tailwind forms plugin for better form styling
