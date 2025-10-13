-- name: CreateUser :one
INSERT INTO users (id)
VALUES (?)
RETURNING *;

-- name: GetUser :one
SELECT * from users
WHERE id = ?;

-- name: CreateResume :one
INSERT INTO resumes (user_id, name, resume)
VALUES (?, ?, ?)
RETURNING *;

-- name: GetResume :one
SELECT * from resumes
WHERE user_id = ? AND name = ?;

-- name: GetResumes :many
SELECT * from resumes
WHERE user_id = ?;
