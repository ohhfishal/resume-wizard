-- TODO: Remove this
-- name: GetNames :many
SELECT name from resumes
ORDER BY name;

-- name: GetResumes :many
SELECT * from resumes
ORDER BY name;

-- name: GetResumeByID :one
SELECT * from resumes
WHERE ID = ?;

-- name: InsertResume :one
INSERT INTO resumes (name, body)
VALUES (?, ?)
RETURNING id;

-- name: UpdateResume :exec
UPDATE resumes
SET 
  body = ?
WHERE 
  id = ?;

-- name: GetApplications :many
SELECT * from applications
ORDER BY created_at;

-- name: UpdateApplication :exec
UPDATE applications
SET
  status = ?
WHERE position = ? AND company = ?;

-- name: InsertLog :one
INSERT INTO applications (resume_id, company, position)
VALUES (?, ?, ?)
RETURNING *;

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
INSERT INTO applications_v2 (
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

-- name: GetApplicationsV2 :many
SELECT * from applications_v2
WHERE user_id = ? AND deleted_at IS NULL;

