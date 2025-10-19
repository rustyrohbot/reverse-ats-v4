# Reverse ATS

A web application for tracking job applications, interviews, and contacts during your job search. Built to replace a spreadsheet-based workflow with a modern, user-friendly interface.

## Tech Stack

- **Backend**: Go with standard library HTTP server
- **Database**: SQLite with [goose](https://github.com/pressly/goose) migrations
- **Type-safe SQL**: [sqlc](https://sqlc.dev/) for generating Go code from SQL queries
- **Frontend**: [HTMX](https://htmx.org/) for dynamic interactions
- **Templates**: [templ](https://templ.guide/) for type-safe HTML templates
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)

## Getting Started

### Prerequisites

- Go 1.23 or later
- Node.js (for Tailwind CSS)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd reverse-ats-v3
```

2. Install Go dependencies:
```bash
go mod download
```

3. Install development tools:
```bash
make install-tools
```

4. Run database migrations:
```bash
make migrate-up
```

5. (Optional) Import existing CSV data:
```bash
go run cmd/import/main.go
```

### Development

Start the development server:
```bash
make run
```

The application will be available at `http://localhost:8080`

## Project Structure

```
├── cmd/
│   ├── server/          # Main application entry point
│   └── import/          # CSV import utility
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── db/              # sqlc generated database code
│   ├── database/        # Database connection setup
│   ├── importer/        # CSV import logic
│   └── templates/       # templ template files
├── migrations/          # goose SQL migration files
├── queries/             # SQL queries for sqlc
├── static/              # CSS, JS, and static assets
└── *.csv               # Initial data export files
```

## Features

- **Company Management** - Track organizations with full CRUD operations
  - Add, edit, and delete companies
  - Store company description, URL, and headquarters location

- **Role Tracking** - Monitor job applications across companies
  - View all applications in a comprehensive table
  - Track salary ranges, equity, location (Remote/Hybrid/On-site)
  - Record application status, dates, and cover letters

- **Interview Scheduling** - Organize interview sessions
  - Link interviews to specific roles
  - Track interview type, date, time, and notes

- **Contact Management** - Maintain recruiter and hiring manager information
  - Associate contacts with companies
  - Store email, phone, LinkedIn, and role information

- **Responsive UI** - Modern interface with Tailwind CSS
  - Full-width tables with proper gridlines
  - HTMX-powered interactions without page reloads
  - Scrollable text columns for descriptions and notes

## Database Schema

The application manages five main entities:

- **Companies** - Organizations you're applying to
- **Roles** - Job positions at companies
- **Interviews** - Interview sessions for specific roles
- **Contacts** - Recruiters and hiring managers at companies
- **InterviewsContacts** - Junction table linking interviews to contacts

See [CLAUDE.md](./CLAUDE.md) for detailed schema information.

## Common Commands

### Database

```bash
# Run migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check migration status
make migrate-status

# Import CSV data
go run cmd/import/main.go
```

### Code Generation

```bash
# Generate all code (sqlc + templ)
make generate

# Generate only sqlc code
sqlc generate

# Generate only templ templates
templ generate
```

### Development

```bash
# Run the application
make run

# Build the application
make build

# Run tests
go test ./...

# Clean generated files
make clean
```

## CSV Import

The import tool reads CSV files from the current directory and populates the database. The expected CSV files are:

- `reverse-ats - Companies.csv`
- `reverse-ats - Contacts.csv`
- `reverse-ats - Roles.csv`
- `reverse-ats - Interviews.csv`
- `reverse-ats - InterviewsContacts.csv`

Usage:
```bash
# Import from current directory
go run cmd/import/main.go

# Import from specific directory
go run cmd/import/main.go -dir /path/to/csvs

# Use custom database path
go run cmd/import/main.go -db ./custom.db
```

## Development Workflow

1. **Make schema changes**: Edit migration files in `migrations/`
2. **Update queries**: Modify SQL files in `queries/`
3. **Regenerate code**: Run `make generate`
4. **Create templates**: Add `.templ` files in `internal/templates/`
5. **Build handlers**: Implement HTTP handlers in `internal/handlers/`
6. **Test**: Run `go test ./...`

## CRUD Operations

Currently implemented:

### Companies
- **Create**: Click "Add company" button, fill form, submit
- **Read**: View all companies at `/companies`
- **Update**: Click "Edit" on any company row, modify form, submit
- **Delete**: Click "Delete" button with confirmation

### Roles, Interviews, Contacts
- Read operations available
- Create, Update, Delete operations: Coming soon

## License

This project is for personal use.
