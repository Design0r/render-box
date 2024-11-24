package service

import (
	"context"
	"database/sql"
	"log"

	"render-box/shared/db/repo"
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

func GetNextTask(db *sql.DB) (*repo.Task, error) {
	r := repo.New(db)

	task, err := r.GetNextTask(context.Background())
	if err != nil {
		log.Printf("Failed to get next task: %v", err)
		return nil, err
	}

	data := repo.UpdateTaskStateParams{State: "progress", ID: task.ID}
	task, err = r.UpdateTaskState(context.Background(), data)
	if err != nil {
		log.Printf("Failed to get next task: %v", err)
		return nil, err
	}

	return &task, nil
}

func UpdateTaskState(db *sql.DB, state string, taskId int64) (*repo.Task, error) {
	r := repo.New(db)

	data := repo.UpdateTaskStateParams{State: state, ID: taskId}
	task, err := r.UpdateTaskState(context.Background(), data)
	if err != nil {
		log.Printf("Failed to update tast state: %v", err)
		return nil, err
	}

	return &task, nil
}
