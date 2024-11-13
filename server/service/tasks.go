package service

import (
	"context"
	"database/sql"
	"log"

	"render-box/server/db/repo"
)

func CreateTask(db *sql.DB, data *repo.CreateTaskParams) (*repo.Task, error) {
	r := repo.New(db)

	task, err := r.CreateTask(context.Background(), *data)
	if err != nil {
		log.Printf("Failed to create task: %v", err)
		return nil, err
	}

	return &task, nil
}

func GetTasks(db *sql.DB) (*[]repo.Task, error) {
	r := repo.New(db)

	tasks, err := r.GetTasks(context.Background())
	if err != nil {
		log.Printf("Failed to get tasks: %v", err)
		return nil, err
	}

	return &tasks, nil
}
