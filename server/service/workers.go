package service

import (
	"context"
	"database/sql"
	"log"

	"render-box/shared/db/repo"
)

func CreateWorker(db *sql.DB, data *repo.CreateWorkerParams) (*repo.Worker, error) {
	r := repo.New(db)

	worker, err := r.CreateWorker(context.Background(), *data)
	if err != nil {
		log.Printf("Failed to create worker: %v", err)
		return nil, err
	}

	return &worker, nil
}

func GetWorkers(db *sql.DB) (*[]repo.Worker, error) {
	r := repo.New(db)

	worker, err := r.GetWorkers(context.Background())
	if err != nil {
		log.Printf("Failed to get workers: %v", err)
		return nil, err
	}

	return &worker, nil
}
