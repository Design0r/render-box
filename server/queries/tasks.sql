-- name: GetNextTask :one
SELECT *
FROM tasks t
JOIN jobs j ON t.job_id = j.id
WHERE j.priority = (
    SELECT MAX(priority)
    FROM jobs
    WHERE state IN ('progress', 'waiting')
)
AND j.state IN ('progress', 'waiting')
AND t.priority = (
    SELECT MAX(priority)
    FROM tasks
    WHERE state = 'waiting' AND job_id = j.id
)
AND t.state = 'waiting'
ORDER BY t.created_at ASC
LIMIT 1;

-- name: UpdateTaskState :exec
UPDATE tasks
SET state = 'progress', edited_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: CreateTask :one
INSERT INTO tasks (priority, data, state, job_id)
VALUES (?, ?, ?, ?)
RETURNING *;

-- name: GetTasks :many
SELECT * from tasks;
