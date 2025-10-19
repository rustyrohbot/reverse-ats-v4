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

5. (Optional) Import sample CSV data:
```bash
go run cmd/import/main.go
```
Note: The importer looks for CSV files in the `./data` directory.

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
├── data/                # Your CSV files for import (gitignored)
└── sample_data/         # Sample CSV files for reference
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

The import tool allows you to bulk-load data from CSV files into the database. This is useful for:
- Initial data migration from spreadsheets
- Testing with sample data
- Backup and restore operations

### Quick Start

1. **Place your CSV files** in the `./data` directory with these exact names:
   - `reverse-ats - Companies.csv`
   - `reverse-ats - Roles.csv`
   - `reverse-ats - Interviews.csv`
   - `reverse-ats - Contacts.csv`
   - `reverse-ats - InterviewsContacts.csv` (optional)

2. **Run the importer**:
   ```bash
   go run cmd/import/main.go
   ```

The importer will automatically:
- Read all CSV files from `./data`
- Import them into `./data.db`
- Display progress and results

### CSV File Formats

Each CSV file must have specific columns. Here are the required formats:

#### Companies (`reverse-ats - Companies.csv`)
```csv
name,description,url,hq_city,hq_state
TechCorp,Leading technology company,https://techcorp.com,San Francisco,CA
DataSys,Data analytics platform,https://datasys.com,Seattle,WA
```

**Columns:**
- `name` (required) - Company name
- `description` (optional) - Company description
- `url` (optional) - Company website
- `hq_city` (optional) - Headquarters city
- `hq_state` (optional) - Headquarters state

#### Roles (`reverse-ats - Roles.csv`)
```csv
company_id,name,url,description,cover_letter,application_location,applied_date,closed_date,posted_range_min,posted_range_max,equity,work_city,work_state,location,status,discovery,referral,notes
1,Senior Backend Engineer,https://techcorp.com/jobs/123,Build scalable systems,Dear Hiring Manager...,https://apply.techcorp.com,2025-01-15,,150,200,0.05-0.15%,San Francisco,CA,HYBRID,APPLIED,LinkedIn,,Great team
```

**Columns:**
- `company_id` (required) - References company (must exist in Companies)
- `name` (required) - Job title/role name
- `url` (optional) - Job posting URL
- `description` (optional) - Job description
- `cover_letter` (optional) - Your cover letter text
- `application_location` (optional) - Where you applied (e.g., "LinkedIn", "Company Website")
- `applied_date` (optional) - Date applied (YYYY-MM-DD format)
- `closed_date` (optional) - Date position closed (YYYY-MM-DD format)
- `posted_range_min` (optional) - Minimum salary in thousands (e.g., 150 for $150k)
- `posted_range_max` (optional) - Maximum salary in thousands (e.g., 200 for $200k)
- `equity` (optional) - Equity range (e.g., "0.05-0.15%")
- `work_city` (optional) - Work location city
- `work_state` (optional) - Work location state
- `location` (optional) - REMOTE, HYBRID, or ON_SITE
- `status` (optional) - RESEARCHING, APPLIED, INTERVIEWING, OFFERED, ACCEPTED, REJECTED, WITHDRAWN
- `discovery` (optional) - How you found the role
- `referral` (optional) - Who referred you
- `notes` (optional) - Additional notes

#### Interviews (`reverse-ats - Interviews.csv`)
```csv
role_id,date,start,end,type,notes
1,2025-01-22,10:00,10:30,RECRUITER,Initial phone screen with Sarah
3,2025-01-25,14:00,15:00,TECH_SCREEN,System design discussion
```

**Columns:**
- `role_id` (required) - References role (must exist in Roles)
- `date` (required) - Interview date (YYYY-MM-DD format)
- `start` (required) - Start time (HH:MM format, 24-hour)
- `end` (required) - End time (HH:MM format, 24-hour)
- `type` (required) - RECRUITER, TECH_SCREEN, MANAGER, LOOP, or MISC
- `notes` (optional) - Interview notes

#### Contacts (`reverse-ats - Contacts.csv`)
```csv
company_id,first_name,last_name,email,phone,linkedin,notes
1,Sarah,Johnson,sarah@techcorp.com,415-555-0123,https://linkedin.com/in/sarahjohnson,Recruiter - very responsive
```

**Columns:**
- `company_id` (required) - References company (must exist in Companies)
- `first_name` (required) - Contact's first name
- `last_name` (required) - Contact's last name
- `email` (optional) - Email address
- `phone` (optional) - Phone number
- `linkedin` (optional) - LinkedIn profile URL
- `notes` (optional) - Notes about the contact

#### InterviewsContacts (`reverse-ats - InterviewsContacts.csv`)
```csv
interview_id,contact_id
1,1
2,1
```

**Columns:**
- `interview_id` (required) - References interview
- `contact_id` (required) - References contact

### Sample Data

Sample CSV files are included in `./sample_data` directory for reference:
- `reverse-ats - Companies.csv`
- `reverse-ats - Roles.csv`
- `reverse-ats - Interviews.csv`
- `reverse-ats - Contacts.csv`
- `reverse-ats - InterviewsContacts.csv`

You can examine these files to understand the expected format. To use the sample data:
```bash
# Copy sample files to data directory
cp sample_data/* data/

# Run the importer
go run cmd/import/main.go
```

### Creating Your Own CSV Files

**Option 1: Export from Google Sheets / Excel**
1. Create sheets with the column headers listed above
2. Fill in your data
3. Export each sheet as CSV with the exact filenames shown
4. Place them in the `./data` directory

**Option 2: Export from Existing Spreadsheet**
If you're already tracking job applications in a spreadsheet:
1. Add/rename columns to match the required format above
2. For `company_id` and `role_id`, create a simple numbering system (1, 2, 3...)
3. Ensure dates are in YYYY-MM-DD format
4. Export as CSV with the correct filenames

**Option 3: Start Fresh**
1. Use the application's web interface to add data manually
2. No CSV import needed!

### Import Tips

- **Import Order**: The importer automatically handles the correct order (Companies → Roles → Interviews/Contacts)
- **IDs**: Company and role IDs in your CSV are matched during import
- **Empty Values**: Leave cells blank for optional fields - don't use "N/A" or "null"
- **Dates**: Always use YYYY-MM-DD format (e.g., 2025-01-15)
- **Times**: Use 24-hour format HH:MM (e.g., 14:30 for 2:30 PM)
- **Re-importing**: The importer appends data - delete `./data.db` first if you want to start fresh

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

### Roles
- **Create**: Click "Add role" button, select company, fill form, submit
- **Read**: View all roles at `/roles`
- **Update**: Click "Edit" on any role row, modify form, submit
- **Delete**: Click "Delete" button with confirmation

### Interviews
- **Create**: Click "Add interview" button, select role, fill form, submit
- **Read**: View all interviews at `/interviews`
- **Update**: Click "Edit" on any interview row, modify form, submit
- **Delete**: Click "Delete" button with confirmation

### Contacts
- **Create**: Click "Add contact" button, select company, fill form, submit
- **Read**: View all contacts at `/contacts`
- **Update**: Click "Edit" on any contact row, modify form, submit
- **Delete**: Click "Delete" button with confirmation

## License

This project is for personal use.
