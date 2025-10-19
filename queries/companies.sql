-- name: GetCompany :one
SELECT * FROM companies
WHERE company_id = ?;

-- name: ListCompanies :many
SELECT * FROM companies
ORDER BY name;

-- name: CreateCompany :one
INSERT INTO companies (name, description, url, hq_city, hq_state)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateCompany :exec
UPDATE companies
SET name = ?, description = ?, url = ?, hq_city = ?, hq_state = ?, updated_at = CURRENT_TIMESTAMP
WHERE company_id = ?;

-- name: DeleteCompany :exec
DELETE FROM companies
WHERE company_id = ?;
