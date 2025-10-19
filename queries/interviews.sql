-- name: GetInterview :one
SELECT * FROM interviews
WHERE interview_id = ?;

-- name: ListInterviews :many
SELECT * FROM interviews
ORDER BY date DESC, start DESC;

-- name: ListInterviewsByRole :many
SELECT * FROM interviews
WHERE role_id = ?
ORDER BY date DESC, start DESC;

-- name: CreateInterview :one
INSERT INTO interviews (role_id, date, start, end, notes, type)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateInterview :exec
UPDATE interviews
SET role_id = ?, date = ?, start = ?, end = ?, notes = ?, type = ?, updated_at = CURRENT_TIMESTAMP
WHERE interview_id = ?;

-- name: DeleteInterview :exec
DELETE FROM interviews
WHERE interview_id = ?;

-- name: LinkInterviewContact :exec
INSERT INTO interviews_contacts (interview_id, contact_id)
VALUES (?, ?)
ON CONFLICT(interview_id, contact_id) DO NOTHING;

-- name: UnlinkInterviewContact :exec
DELETE FROM interviews_contacts
WHERE interview_id = ? AND contact_id = ?;

-- name: GetInterviewContacts :many
SELECT c.* FROM contacts c
INNER JOIN interviews_contacts ic ON c.contact_id = ic.contact_id
WHERE ic.interview_id = ?;
