package routes

import (
	"database/sql"

	"render-box/server/service"
	"render-box/shared"
	"render-box/shared/db/repo"
)

func InitJobRouter() *shared.MessageRouter {
	router := shared.NewMessageRouter()
	router.Register("jobs.create", CreateJob)
	router.Register("jobs.all", AllJobs)

	return router
}

func CreateJob(
	db *sql.DB,
	message *shared.Message,
	state *shared.ConnState,
) (any, error) {
	data, err := shared.UnmarshallBody[repo.CreateJobParams](message.Data)
	if err != nil {
		return nil, err
	}
	task, err := service.CreateJob(db, data)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func AllJobs(
	db *sql.DB,
	message *shared.Message,
	state *shared.ConnState,
) (any, error) {
	tasks, err := service.GetJobs(db)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
