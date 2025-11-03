-- name: GetRole :one
SELECT * FROM roles
WHERE role_id = ?;

-- name: ListRoles :many
SELECT * FROM roles
ORDER BY applied_date DESC;

-- name: ListRolesByCompany :many
SELECT * FROM roles
WHERE company_id = ?
ORDER BY applied_date DESC;

-- name: CreateRole :one
INSERT INTO roles (
    company_id, name, url, description, cover_letter, application_location,
    applied_date, closed_date, posted_range_min, posted_range_max, equity,
    work_city, work_state, location, status, discovery, referral, notes
)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateRole :exec
UPDATE roles
SET company_id = ?, name = ?, url = ?, description = ?, cover_letter = ?, application_location = ?,
    applied_date = ?, closed_date = ?, posted_range_min = ?, posted_range_max = ?, equity = ?,
    work_city = ?, work_state = ?, location = ?, status = ?, discovery = ?, referral = ?, notes = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE role_id = ?;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE role_id = ?;
