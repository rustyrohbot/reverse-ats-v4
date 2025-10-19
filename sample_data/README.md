# Sample Data

This directory contains sample CSV files that demonstrate the expected format for importing data into the Reverse ATS application.

## Files

- `reverse-ats - Companies.csv` - 8 sample companies
- `reverse-ats - Roles.csv` - 8 sample job roles with various statuses
- `reverse-ats - Interviews.csv` - 7 sample interviews
- `reverse-ats - Contacts.csv` - 9 sample contacts
- `reverse-ats - InterviewsContacts.csv` - 2 sample interview-contact relationships

## Using Sample Data

To test the application with sample data:

1. Copy the sample files to the `data/` directory:
   ```bash
   cp sample_data/* data/
   ```

2. Run the importer:
   ```bash
   go run cmd/import/main.go
   ```

## Creating Your Own Data

Use these files as templates for creating your own CSV files. The column headers must match exactly, but you can modify the data rows to match your job search information.

For detailed documentation on the CSV format and column descriptions, see the main [README.md](../README.md#csv-import).
