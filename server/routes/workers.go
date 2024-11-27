package routes

import (
	"database/sql"

	"render-box/server/service"
	"render-box/shared"
	"render-box/shared/db/repo"
)

func InitWorkerRouter() *shared.MessageRouter {
	router := shared.NewMessageRouter()
	router.Register(string(shared.MSGWorkerCreate), CreateWorker)
	router.Register(string(shared.MSGWorkerAll), AllWorkers)
	router.Register(string(shared.MSGWorkerRegister), RegisterWorker)

	return router
}

func CreateWorker(
	db *sql.DB,
	message *shared.Message,
	state *shared.ConnState,
) (interface{}, error) {
	data, err := shared.UnmarshallBody[repo.CreateWorkerParams](message.Data)
	if err != nil {
		return nil, err
	}
	worker, err := service.CreateWorker(db, data)
	if err != nil {
		return nil, err
	}
	return worker, nil
}

func AllWorkers(
	db *sql.DB,
	message *shared.Message,
	state *shared.ConnState,
) (interface{}, error) {
	worker, err := service.GetWorkers(db)
	if err != nil {
		return nil, err
	}
	return worker, nil
}

func RegisterWorker(
	db *sql.DB,
	message *shared.Message,
	state *shared.ConnState,
) (interface{}, error) {
	name := (message.Data).(string)
	worker, err := service.RegisterWorker(db, name)
	if err != nil {
		return nil, err
	}
	state.Worker = worker
	return worker, err
}
