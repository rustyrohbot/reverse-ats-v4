-- +goose Up
-- +goose StatementBegin
CREATE TABLE companies (
    company_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    url TEXT,
    linkedin TEXT,
    hq_city TEXT,
    hq_state TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE contacts (
    contact_id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    role TEXT,
    email TEXT,
    phone TEXT,
    linkedin TEXT,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(company_id) ON DELETE CASCADE
);

CREATE TABLE roles (
    role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    company_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    url TEXT,
    description TEXT,
    cover_letter TEXT,
    application_location TEXT,
    applied_date TEXT,
    closed_date TEXT,
    posted_range_min INTEGER,
    posted_range_max INTEGER,
    equity BOOLEAN,
    work_city TEXT,
    work_state TEXT,
    location TEXT,
    status TEXT,
    discovery TEXT,
    referral BOOLEAN,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(company_id) ON DELETE CASCADE
);

CREATE TABLE interviews (
    interview_id INTEGER PRIMARY KEY AUTOINCREMENT,
    role_id INTEGER NOT NULL,
    date TEXT NOT NULL,
    start TEXT NOT NULL,
    end TEXT NOT NULL,
    notes TEXT,
    type TEXT NOT NULL CHECK(type IN ('RECRUITER', 'LOOP', 'TECH_SCREEN', 'MANAGER', 'MISC')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE
);

CREATE TABLE interviews_contacts (
    interviews_contact_id INTEGER PRIMARY KEY AUTOINCREMENT,
    interview_id INTEGER NOT NULL,
    contact_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (interview_id) REFERENCES interviews(interview_id) ON DELETE CASCADE,
    FOREIGN KEY (contact_id) REFERENCES contacts(contact_id) ON DELETE CASCADE,
    UNIQUE(interview_id, contact_id)
);

-- Create indexes for better query performance
CREATE INDEX idx_contacts_company_id ON contacts(company_id);
CREATE INDEX idx_roles_company_id ON roles(company_id);
CREATE INDEX idx_interviews_role_id ON interviews(role_id);
CREATE INDEX idx_interviews_contacts_interview_id ON interviews_contacts(interview_id);
CREATE INDEX idx_interviews_contacts_contact_id ON interviews_contacts(contact_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_interviews_contacts_contact_id;
DROP INDEX IF EXISTS idx_interviews_contacts_interview_id;
DROP INDEX IF EXISTS idx_interviews_role_id;
DROP INDEX IF EXISTS idx_roles_company_id;
DROP INDEX IF EXISTS idx_contacts_company_id;
DROP TABLE IF EXISTS interviews_contacts;
DROP TABLE IF EXISTS interviews;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS companies;
-- +goose StatementEnd
