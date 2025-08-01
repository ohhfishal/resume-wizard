-- TODO: Remove this
-- name: GetNames :many
SELECT name from resumes
ORDER BY name;

-- name: GetResumes :many
SELECT * from resumes
ORDER BY name;

-- name: InsertResume :one
INSERT INTO resumes (name, body)
VALUES (?, ?)
RETURNING id;

-- name: GetApplications :many
SELECT * from applications
ORDER BY created_at;

-- name: InsertLog :one
INSERT INTO applications (resume_id, company, position)
VALUES (?, ?, ?)
RETURNING *;
