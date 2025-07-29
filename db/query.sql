-- name: GetNames :many
SELECT name from resumes
ORDER BY name;

-- name: GetRows :many
SELECT * from resumes
ORDER BY name;
