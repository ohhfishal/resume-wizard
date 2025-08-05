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

-- name: InsertLog :one
INSERT INTO applications (resume_id, company, position)
VALUES (?, ?, ?)
RETURNING *;
