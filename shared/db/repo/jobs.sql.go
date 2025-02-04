// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: jobs.sql

package repo

import (
	"context"
)

const createJob = `-- name: CreateJob :one
INSERT INTO jobs (name, priority, state)
VALUES (?, ?, ?)
RETURNING id, name, priority, state, created_at, edited_at
`

type CreateJobParams struct {
	Name     string `json:"name"`
	Priority int64  `json:"priority"`
	State    string `json:"state"`
}

func (q *Queries) CreateJob(ctx context.Context, arg CreateJobParams) (Job, error) {
	row := q.db.QueryRowContext(ctx, createJob, arg.Name, arg.Priority, arg.State)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Priority,
		&i.State,
		&i.CreatedAt,
		&i.EditedAt,
	)
	return i, err
}

const getJobs = `-- name: GetJobs :many
SELECT id, name, priority, state, created_at, edited_at FROM jobs
`

func (q *Queries) GetJobs(ctx context.Context) ([]Job, error) {
	rows, err := q.db.QueryContext(ctx, getJobs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Job
	for rows.Next() {
		var i Job
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Priority,
			&i.State,
			&i.CreatedAt,
			&i.EditedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const restoreJobState = `-- name: RestoreJobState :exec
UPDATE jobs
SET state = 'waiting', edited_at = CURRENT_TIMESTAMP
WHERE jobs.id = ?
AND (SELECT COUNT(*) FROM tasks t 
      WHERE t.job_id = jobs.id 
      AND t.state in ('progress', 'waiting')
    ) > 0
`

func (q *Queries) RestoreJobState(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, restoreJobState, id)
	return err
}

const updateCompletedJob = `-- name: UpdateCompletedJob :one
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
RETURNING id, name, priority, state, created_at, edited_at
`

func (q *Queries) UpdateCompletedJob(ctx context.Context, id int64) (Job, error) {
	row := q.db.QueryRowContext(ctx, updateCompletedJob, id)
	var i Job
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Priority,
		&i.State,
		&i.CreatedAt,
		&i.EditedAt,
	)
	return i, err
}

const updateJobState = `-- name: UpdateJobState :exec
UPDATE jobs
SET state = ?, edited_at = CURRENT_TIMESTAMP
WHERE jobs.id = ?
`

type UpdateJobStateParams struct {
	State string `json:"state"`
	ID    int64  `json:"id"`
}

func (q *Queries) UpdateJobState(ctx context.Context, arg UpdateJobStateParams) error {
	_, err := q.db.ExecContext(ctx, updateJobState, arg.State, arg.ID)
	return err
}
