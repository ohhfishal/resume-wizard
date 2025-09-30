-- name: InsertResume :one
INSERT INTO resumes (user_id, name, resume)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetResume :one
SELECT * from resumes
WHERE user_id = ? AND id = ?;

-- name: GetResumes :many
SELECT * from resumes
WHERE user_id = ?
ORDER BY created_at; -- Last used??

-- name: InsertApplicationActivity :one
INSERT INTO application_history (
    user_id,
    resume_id,
    company,
    position,
    description,
    status
) VALUES (
    ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetApplicationHistory :many
SELECT * from application_history
WHERE user_id = ?;
