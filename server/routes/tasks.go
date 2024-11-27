package routes

import (
	"database/sql"

	"render-box/server/service"
	"render-box/shared"
	"render-box/shared/db/repo"
)

func InitTaskRouter() *shared.MessageRouter {
	router := shared.NewMessageRouter()
	router.Register(string(shared.MSGTasksCreate), CreateTask)
	router.Register(string(shared.MSGTasksAll), AllTasks)
	router.Register(string(shared.MSGTasksNext), NextTask)

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
	task, err := service.GetNextTask(db)
	if err != nil {
		return nil, err
	}

	_, err = service.UpdateWorkerState(db, "working", state.Worker.ID)
	if err != nil {
		return nil, err
	}

	worker, err := service.UpdateWorkerTask(db, state.Worker.ID, task.ID)
	if err != nil {
		return nil, err
	}

	state.Worker = worker
	state.Task = task
	return task, err
}
