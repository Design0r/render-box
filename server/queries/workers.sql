-- name: CreateWorker :one
INSERT INTO workers (name, state, metadata, task_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetWorkers :many
SELECT * FROM workers;

-- name: GetWorkerByName :one
SELECT * FROM workers
WHERE name = ?;

-- name: UpdateWorkerState :one
UPDATE workers
SET state = ?, edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateWorkerTask :one
UPDATE workers
SET task_id = ?, edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

