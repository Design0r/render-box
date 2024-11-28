package routes

import (
	"database/sql"
	"log"

	"render-box/server/service"
	"render-box/shared"
	"render-box/shared/db/repo"
)

func InitTaskRouter() *shared.MessageRouter {
	router := shared.NewMessageRouter()
	router.Register(string(shared.MSGTasksCreate), CreateTask)
	router.Register(string(shared.MSGTasksAll), AllTasks)
	router.Register(string(shared.MSGTasksNext), NextTask)
	router.Register(string(shared.MSGTasksComplete), CompleteTask)

	return router
}

func CreateTask(
	db *sql.DB,
	message *shared.Message,
	state *shared.ConnState,
) (interface{}, error) {
	data, err := shared.UnmarshallBody[repo.CreateTaskParams](message.Data)
	if err != nil {
		return nil, err
	}
	task, err := service.CreateTask(db, data)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func AllTasks(
	db *sql.DB,
	message *shared.Message,
	state *shared.ConnState,
) (interface{}, error) {
	tasks, err := service.GetTasks(db)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func NextTask(
	db *sql.DB,
	message *shared.Message,
	state *shared.ConnState,
) (interface{}, error) {
	workerState := "working"
	var workerTaskId *int64

	task, terr := service.GetNextTask(db)
	if terr != nil {
		log.Println("No waiting tasks...")
		workerState = "waiting"
	} else {
		workerTaskId = &task.ID
		err := service.UpdateJobState(db, "progress", task.JobID)
		if err != nil {
			return nil, err
		}
	}

	_, err := service.UpdateWorkerState(db, workerState, state.Worker.ID)
	if err != nil {
		return nil, err
	}

	worker, err := service.UpdateWorkerTask(db, state.Worker.ID, workerTaskId)
	if err != nil {
		return nil, err
	}

	state.Worker = worker
	state.Task = task

	if terr != nil {
		return nil, terr
	}

	return task, nil
}

func CompleteTask(
	db *sql.DB,
	message *shared.Message,
	state *shared.ConnState,
) (interface{}, error) {
	task, err := service.UpdateTaskState(db, "completed", state.Task.ID)
	if err != nil {
		return nil, err
	}

	_, err = service.UpdateCompletedJob(db, task.ID)
	if err != nil {
		return nil, err
	}

	state.Task = task
	return task, nil
}
