-- name: ListRolesWithCompany :many
SELECT
    r.role_id,
    r.company_id,
    r.name as role_name,
    r.url,
    r.description,
    r.cover_letter,
    r.application_location,
    r.applied_date,
    r.closed_date,
    r.posted_range_min,
    r.posted_range_max,
    r.equity,
    r.work_city,
    r.work_state,
    r.location,
    r.status,
    r.discovery,
    r.referral,
    r.notes,
    r.created_at,
    r.updated_at,
    c.name as company_name
FROM roles r
INNER JOIN companies c ON r.company_id = c.company_id
ORDER BY r.applied_date DESC;

-- name: ListInterviewsWithRole :many
SELECT
    i.interview_id,
    i.role_id,
    i.date,
    i.start,
    i.end,
    i.notes,
    i.type,
    i.created_at,
    i.updated_at,
    r.name as role_name,
    c.company_id,
    c.name as company_name
FROM interviews i
INNER JOIN roles r ON i.role_id = r.role_id
INNER JOIN companies c ON r.company_id = c.company_id
ORDER BY i.date DESC, i.start DESC;

-- name: ListContactsWithCompany :many
SELECT
    ct.contact_id,
    ct.company_id,
    ct.first_name,
    ct.last_name,
    ct.role,
    ct.email,
    ct.phone,
    ct.linkedin,
    ct.notes,
    ct.created_at,
    ct.updated_at,
    c.name as company_name
FROM contacts ct
INNER JOIN companies c ON ct.company_id = c.company_id
ORDER BY c.name, ct.last_name, ct.first_name;

-- name: GetRoleWithCompany :one
SELECT
    r.role_id,
    r.name as role_name,
    r.url,
    r.description,
    r.status,
    r.applied_date,
    r.closed_date,
    r.posted_range_min,
    r.posted_range_max,
    r.work_city,
    r.work_state,
    r.location,
    r.discovery,
    r.referral,
    r.notes,
    c.company_id,
    c.name as company_name
FROM roles r
INNER JOIN companies c ON r.company_id = c.company_id
WHERE r.role_id = ?;

-- name: GetInterviewWithRoleAndCompany :one
SELECT
    i.interview_id,
    i.date,
    i.start,
    i.end,
    i.type,
    i.notes,
    r.role_id,
    r.name as role_name,
    c.company_id,
    c.name as company_name
FROM interviews i
INNER JOIN roles r ON i.role_id = r.role_id
INNER JOIN companies c ON r.company_id = c.company_id
WHERE i.interview_id = ?;
