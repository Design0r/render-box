-- name: CreateJob :one
INSERT INTO jobs (name, priority, state)
VALUES (?, ?, ?)
RETURNING *;


-- name: GetJobs :many
SELECT * FROM jobs
