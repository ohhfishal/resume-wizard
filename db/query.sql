-- name: CreateUser :one
INSERT INTO users (id)
VALUES (?)
RETURNING *;

-- name: GetUser :one
SELECT * from users
WHERE id = ?;

