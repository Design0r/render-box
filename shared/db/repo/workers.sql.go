// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: workers.sql

package repo

import (
	"context"
)

const createWorker = `-- name: CreateWorker :one
INSERT INTO workers (name, state, metadata, task_id)
VALUES (?, ?, ?, ?)
RETURNING id, name, state, metadata, created_at, edited_at, task_id
`

type CreateWorkerParams struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	Metadata string `json:"metadata"`
	TaskID   *int64 `json:"task_id"`
}

func (q *Queries) CreateWorker(ctx context.Context, arg CreateWorkerParams) (Worker, error) {
	row := q.db.QueryRowContext(ctx, createWorker,
		arg.Name,
		arg.State,
		arg.Metadata,
		arg.TaskID,
	)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.State,
		&i.Metadata,
		&i.CreatedAt,
		&i.EditedAt,
		&i.TaskID,
	)
	return i, err
}

const getWorkerByName = `-- name: GetWorkerByName :one
SELECT id, name, state, metadata, created_at, edited_at, task_id FROM workers
WHERE name = ?
`

func (q *Queries) GetWorkerByName(ctx context.Context, name string) (Worker, error) {
	row := q.db.QueryRowContext(ctx, getWorkerByName, name)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.State,
		&i.Metadata,
		&i.CreatedAt,
		&i.EditedAt,
		&i.TaskID,
	)
	return i, err
}

const getWorkers = `-- name: GetWorkers :many
SELECT id, name, state, metadata, created_at, edited_at, task_id FROM workers
`

func (q *Queries) GetWorkers(ctx context.Context) ([]Worker, error) {
	rows, err := q.db.QueryContext(ctx, getWorkers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Worker
	for rows.Next() {
		var i Worker
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.State,
			&i.Metadata,
			&i.CreatedAt,
			&i.EditedAt,
			&i.TaskID,
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

const updateWorkerState = `-- name: UpdateWorkerState :one
UPDATE workers
SET state = ?, edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, name, state, metadata, created_at, edited_at, task_id
`

type UpdateWorkerStateParams struct {
	State string `json:"state"`
	ID    int64  `json:"id"`
}

func (q *Queries) UpdateWorkerState(ctx context.Context, arg UpdateWorkerStateParams) (Worker, error) {
	row := q.db.QueryRowContext(ctx, updateWorkerState, arg.State, arg.ID)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.State,
		&i.Metadata,
		&i.CreatedAt,
		&i.EditedAt,
		&i.TaskID,
	)
	return i, err
}

const updateWorkerTask = `-- name: UpdateWorkerTask :one
UPDATE workers
SET task_id = ?, edited_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING id, name, state, metadata, created_at, edited_at, task_id
`

type UpdateWorkerTaskParams struct {
	TaskID *int64 `json:"task_id"`
	ID     int64  `json:"id"`
}

func (q *Queries) UpdateWorkerTask(ctx context.Context, arg UpdateWorkerTaskParams) (Worker, error) {
	row := q.db.QueryRowContext(ctx, updateWorkerTask, arg.TaskID, arg.ID)
	var i Worker
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.State,
		&i.Metadata,
		&i.CreatedAt,
		&i.EditedAt,
		&i.TaskID,
	)
	return i, err
}
