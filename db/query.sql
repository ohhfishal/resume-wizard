-- name: GetNames :many
SELECT name from resumes
ORDER BY name;

-- name: GetRows :many
SELECT * from resumes
ORDER BY name;

-- name: InsertResume :one
INSERT INTO resumes (name, body)
VALUES (?, ?)
RETURNING id;

-- name: InsertLog :one
INSERT INTO applications (resume_id, company, position)
VALUES (?, ?, ?)
RETURNING *;
