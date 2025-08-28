-- name: InsertBase :one
INSERT INTO base_resumes (user_id, name, resume)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetBaseResume :one
SELECT * from base_resumes
WHERE user_id = ? AND id = ?;

-- name: GetBaseResumes :many
SELECT * from base_resumes
WHERE user_id = ?
ORDER BY created_at; -- Last used??

-- name: CreateSession :one
-- Create a session of a user working on an application
INSERT INTO sessions (
  uuid,
  base_resume_id,
  user_id,
  company,
  position,
  description,
  resume -- TODO: Think this can be removed
) VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: GetSession :one
-- Get a session (user space)
SELECT * FROM sessions
WHERE uuid = ? AND user_id = ? AND deleted_at IS NULL;

-- name: AddResumeToSession :exec
UPDATE sessions 
SET resume = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE uuid = ? AND user_id = ? AND deleted_at IS NULL;

-- name: SoftDeleteSession :exec
UPDATE sessions 
SET deleted_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE uuid = ? AND user_id = ? AND deleted_at IS NULL;

-- name: CreateApplication :one
INSERT INTO applications (
    user_id,
    base_resume_id,
    company,
    position,
    description,
    resume,
    status
) VALUES (
    ?, ?, ?, ?, ?, ?, 'pending'
) RETURNING *;

-- name: GetApplications :many
SELECT * from applications
WHERE user_id = ? AND deleted_at IS NULL;

-- name: UpdateApplication :one
UPDATE applications 
SET 
    applied_at = ?,
    status = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = ? AND id = ? AND deleted_at IS NULL
RETURNING *;
