-- name: GetContact :one
SELECT * FROM contacts
WHERE contact_id = ?;

-- name: ListContacts :many
SELECT * FROM contacts
ORDER BY last_name, first_name;

-- name: ListContactsByCompany :many
SELECT * FROM contacts
WHERE company_id = ?
ORDER BY last_name, first_name;

-- name: CreateContact :one
INSERT INTO contacts (company_id, first_name, last_name, role, email, phone, linkedin, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateContact :exec
UPDATE contacts
SET company_id = ?, first_name = ?, last_name = ?, role = ?, email = ?, phone = ?, linkedin = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
WHERE contact_id = ?;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE contact_id = ?;
