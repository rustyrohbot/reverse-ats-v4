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
cd reverse-ats-v4
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
Note: The importer looks for CSV files in the `./import` directory.

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
│   ├── import/          # CSV import CLI utility
│   └── export/          # CSV export CLI utility
├── internal/
│   ├── handlers/        # HTTP request handlers
│   ├── db/              # sqlc generated database code
│   ├── database/        # Database connection setup
│   ├── importer/        # CSV import logic (shared)
│   ├── exporter/        # CSV export logic (shared)
│   ├── util/            # Shared utilities (date formatting, etc.)
│   └── templates/       # templ template files
├── migrations/          # goose SQL migration files
├── queries/             # SQL queries for sqlc
├── static/              # CSS, JS, and static assets
├── import/              # CSV files for import (gitignored)
├── export/              # Exported CSV files (gitignored)
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

- **Data Import/Export** - Flexible data management
  - Web-based CSV import via drag-and-drop modal
  - CLI-based batch import from directory
  - CLI-based export to CSV files
  - Consistent NULL value handling across import/export

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

# Import CSV data from ./import directory
go run cmd/import/main.go

# Export all data to ./export directory
go run cmd/export/main.go
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

## Data Import and Export

The application provides two methods for importing and exporting data:

### Web-Based Import

Import CSV files directly through the web interface:

1. Click the **Import** button in the top navigation
2. Upload CSV files (one or more of the following):
   - Companies
   - Roles
   - Contacts
   - Interviews
   - InterviewsContacts
3. Files are validated for size (max 10MB) and type (.csv only)
4. Data is imported with proper foreign key ordering

**Security features:**
- File type validation (CSV only)
- File size limits (10MB per file)
- Secure temporary file handling
- Automatic cleanup after processing

### CLI-Based Import

The CLI import tool allows you to bulk-load data from CSV files into the database. This is useful for:
- Initial data migration from spreadsheets
- Testing with sample data
- Backup and restore operations
- Batch imports without using the web interface

#### CLI Import Quick Start

1. **Place your CSV files** in the `./import` directory with these exact names:
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
- Read all CSV files from `./import`
- Import them into `./data.db`
- Display progress and results

### CSV File Formats

Each CSV file must have specific columns. Here are the required formats:

#### Companies (`reverse-ats - Companies.csv`)
```csv
companyID,name,description,url,linkedin,hqCity,hqState
1,TechCorp,Leading technology company,https://techcorp.com,NULL,San Francisco,CA
2,DataSys,NULL,https://datasys.com,NULL,Seattle,WA
```

**Columns:**
- `companyID` (optional for import) - Company ID (auto-generated if not provided)
- `name` (required) - Company name
- `description` (optional) - Company description (leave empty or use "NULL")
- `url` (optional) - Company website (leave empty or use "NULL")
- `linkedin` (optional) - Company LinkedIn URL (leave empty or use "NULL")
- `hqCity` (optional) - Headquarters city (leave empty or use "NULL")
- `hqState` (optional) - Headquarters state (leave empty or use "NULL")

#### Roles (`reverse-ats - Roles.csv`)
```csv
company_id,name,url,description,cover_letter,application_location,applied_date,closed_date,posted_range_min,posted_range_max,equity,work_city,work_state,location,status,discovery,referral,notes
1,Senior Backend Engineer,https://techcorp.com/jobs/123,Build scalable systems,NULL,LinkedIn,2025-01-15,NULL,150,200,0.05-0.15%,San Francisco,CA,HYBRID,APPLIED,LinkedIn,NULL,Great team
```

**Columns:**
- `company_id` (required) - References company (must exist in Companies)
- `name` (required) - Job title/role name
- `url` (optional) - Job posting URL (use "NULL" if empty)
- `description` (optional) - Job description (use "NULL" if empty)
- `cover_letter` (optional) - Your cover letter text (use "NULL" if empty)
- `application_location` (optional) - Where you applied (use "NULL" if empty)
- `applied_date` (optional) - Date applied in YYYY-MM-DD format (use "NULL" if empty)
- `closed_date` (optional) - Date position closed in YYYY-MM-DD format (use "NULL" if empty)
- `posted_range_min` (optional) - Minimum salary in thousands (use "NULL" if empty)
- `posted_range_max` (optional) - Maximum salary in thousands (use "NULL" if empty)
- `equity` (optional) - Equity range like "0.05-0.15%" (use "NULL" if empty)
- `work_city` (optional) - Work location city (use "NULL" if empty)
- `work_state` (optional) - Work location state (use "NULL" if empty)
- `location` (optional) - REMOTE, HYBRID, or ON_SITE (use "NULL" if empty)
- `status` (optional) - RESEARCHING, APPLIED, INTERVIEWING, OFFERED, ACCEPTED, REJECTED, WITHDRAWN (use "NULL" if empty)
- `discovery` (optional) - How you found the role (use "NULL" if empty)
- `referral` (optional) - Who referred you (use "NULL" if empty)
- `notes` (optional) - Additional notes (use "NULL" if empty)

#### Interviews (`reverse-ats - Interviews.csv`)
```csv
role_id,date,start,end,type,notes
1,2025-01-22,10:00,10:30,RECRUITER,Initial phone screen with Sarah
3,2025-01-25,14:00,15:00,TECH_SCREEN,NULL
```

**Columns:**
- `role_id` (required) - References role (must exist in Roles)
- `date` (required) - Interview date in YYYY-MM-DD format
- `start` (required) - Start time in HH:MM format (24-hour)
- `end` (required) - End time in HH:MM format (24-hour)
- `type` (required) - RECRUITER, TECH_SCREEN, MANAGER, LOOP, or MISC
- `notes` (optional) - Interview notes (use "NULL" if empty)

#### Contacts (`reverse-ats - Contacts.csv`)
```csv
company_id,first_name,last_name,email,phone,linkedin,notes
1,Sarah,Johnson,sarah@techcorp.com,415-555-0123,https://linkedin.com/in/sarahjohnson,Recruiter - very responsive
2,Michael,Chen,NULL,NULL,https://linkedin.com/in/michaelchen,NULL
```

**Columns:**
- `company_id` (required) - References company (must exist in Companies)
- `first_name` (required) - Contact's first name
- `last_name` (required) - Contact's last name
- `email` (optional) - Email address (use "NULL" if empty)
- `phone` (optional) - Phone number (use "NULL" if empty)
- `linkedin` (optional) - LinkedIn profile URL (use "NULL" if empty)
- `notes` (optional) - Notes about the contact (use "NULL" if empty)

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
# Create import directory and copy sample files
mkdir -p import
cp sample_data/* import/

# Run the importer
go run cmd/import/main.go
```

### Creating Your Own CSV Files

**Option 1: Export from Google Sheets / Excel**
1. Create sheets with the column headers listed above
2. Fill in your data
3. Export each sheet as CSV with the exact filenames shown
4. Place them in the `./import` directory

**Option 2: Export from Existing Spreadsheet**
If you're already tracking job applications in a spreadsheet:
1. Add/rename columns to match the required format above
2. For `company_id` and `role_id`, create a simple numbering system (1, 2, 3...)
3. Ensure dates are in YYYY-MM-DD format
4. Export as CSV with the correct filenames

**Option 3: Start Fresh**
1. Use the application's web interface to add data manually
2. No CSV import needed!


### CLI-Based Export

Export all data to CSV files for backup or external analysis:

```bash
# Export all data to ./export directory
go run cmd/export/main.go
```

The export tool will:
- Create the `./export` directory if it doesn't exist
- Export all five tables to separate CSV files
- Use consistent NULL formatting ("NULL" string for missing values)
- Name files with the standard naming convention

**Exported files:**
- `reverse-ats - Companies.csv`
- `reverse-ats - Roles.csv`
- `reverse-ats - Contacts.csv`
- `reverse-ats - Interviews.csv`
- `reverse-ats - InterviewsContacts.csv`

You can also export via the web interface by clicking the **Export** button, which downloads a zip file containing all CSV files.

### Import/Export Tips

- **Import Order**: The importer automatically handles the correct order (Companies → Roles → Contacts → Interviews → InterviewsContacts)
- **IDs**: Company and role IDs in your CSV are matched during import
- **NULL Values**: Use the string "NULL" (all caps) for missing/empty values in CSV files. Both import and export use this convention for consistency.
- **Dates**: Both import and export accept/preserve multiple formats:
  - ISO format: YYYY-MM-DD (e.g., "2025-01-15")
  - Text format: "Month Day, Year" (e.g., "January 15, 2025" or "April 14, 2025")
  - Data is stored as-is in the database
  - UI displays dates in text format regardless of storage format
- **Times**: Both import and export accept/preserve multiple formats:
  - 24-hour format: HH:MM (e.g., "14:30")
  - 12-hour format: HH:MM AM/PM (e.g., "2:30 PM")
  - Data is stored as-is in the database
- **Re-importing**: The importer appends data - delete `./data.db` first if you want to start fresh
- **Round-trip Compatibility**: Files exported via the CLI can be directly re-imported without modification

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
