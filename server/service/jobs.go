package service

import (
	"context"
	"database/sql"
	"log"

	"render-box/shared/db/repo"
)

func CreateJob(db *sql.DB, data *repo.CreateJobParams) (*repo.Job, error) {
	r := repo.New(db)

	job, err := r.CreateJob(context.Background(), *data)
	if err != nil {
		log.Printf("Failed to create task: %v", err)
		return nil, err
	}

	return &job, nil
}

func GetJobs(db *sql.DB) (*[]repo.Job, error) {
	r := repo.New(db)

	tasks, err := r.GetJobs(context.Background())
	if err != nil {
		log.Printf("Failed to get tasks: %v", err)
		return nil, err
	}

	return &tasks, nil
}
