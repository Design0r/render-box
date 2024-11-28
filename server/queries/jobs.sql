-- name: CreateJob :one
INSERT INTO jobs (name, priority, state)
VALUES (?, ?, ?)
RETURNING *;


-- name: GetJobs :many
SELECT * FROM jobs;

-- name: UpdateCompletedJob :one
UPDATE jobs
SET state = 'completed', edited_at = CURRENT_TIMESTAMP
WHERE id = (
    SELECT job_id FROM tasks t WHERE t.id = ?
)
AND NOT EXISTS (
    SELECT 1 FROM tasks
    WHERE job_id = jobs.id
      AND state = 'waiting'
)
RETURNING *;

-- name: UpdateJobState :exec
UPDATE jobs
SET state = ?, edited_at = CURRENT_TIMESTAMP
WHERE jobs.id = ?;

-- name: RestoreJobState :exec
UPDATE jobs
SET state = 'waiting', edited_at = CURRENT_TIMESTAMP
WHERE jobs.id = ?
AND (SELECT COUNT(*) FROM tasks t 
      WHERE t.job_id = jobs.id 
      AND t.state in ('progress', 'waiting')
    ) > 0;
