-- name: CreateWorker :one
INSERT INTO workers (name, state, metadata, task_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetWorkers :many
SELECT * FROM workers;
