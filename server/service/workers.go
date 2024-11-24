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

func UpdateWorkerState(db *sql.DB, state string, workerId int64) (*repo.Worker, error) {
	r := repo.New(db)

	data := repo.UpdateWorkerStateParams{State: state, ID: workerId}
	worker, err := r.UpdateWorkerState(context.Background(), data)
	if err != nil {
		log.Printf("Failed to update worker state: %v", err)
		return nil, err
	}

	return &worker, err
}

func UpdateWorkerTask(db *sql.DB, workerId int64, taskId int64) (*repo.Worker, error) {
	r := repo.New(db)

	data := repo.UpdateWorkerTaskParams{ID: workerId, TaskID: &taskId}
	worker, err := r.UpdateWorkerTask(context.Background(), data)
	if err != nil {
		log.Printf("Failed to update worker task: %v", err)
		return nil, err
	}

	return &worker, err
}

func GetWorkerByName(db *sql.DB, name string) (*repo.Worker, error) {
	r := repo.New(db)

	worker, err := r.GetWorkerByName(context.Background(), name)
	if err != nil {
		log.Printf("Failed to get worker: %v", err)
		return nil, err
	}

	return &worker, nil
}

func RegisterWorker(db *sql.DB, name string) (*repo.Worker, error) {
	worker, err := GetWorkerByName(db, name)
	if err != nil {
		data := &repo.CreateWorkerParams{Name: name, State: "waiting", Metadata: "", TaskID: nil}
		worker, err := CreateWorker(db, data)
		if err != nil {
			return nil, err
		}
		return worker, err
	}

	worker, err = UpdateWorkerState(db, "waiting", worker.ID)
	if err != nil {
		return nil, err
	}

	return worker, err
}
