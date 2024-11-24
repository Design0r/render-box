package main

import (
	"log"
	"net"

	"render-box/shared"
	"render-box/shared/db/repo"
)

func createJob(conn *net.Conn) (*repo.Job, error) {
	job := repo.CreateJobParams{Name: "TestJob", Priority: int64(50), State: "waiting"}
	msg := shared.Message{Type: shared.MSGJobsCreate, Data: &job}
	msg.Send(conn)

	response, err := shared.RecvMessage[repo.Job](conn)
	log.Println(response)
	if err != nil {
		log.Printf("ERROR: could not read message: %v", err)
		return nil, err
	}
	jobResposne := (response.Data).(repo.Job)

	return &jobResposne, nil
}

func createTask(conn *net.Conn, jobId int64) error {
	log.Println("Creating task with jobId:", jobId)
	task := repo.CreateTaskParams{
		Priority: 10,
		Data:     "",
		State:    "waiting",
		JobID:    jobId,
	}
	msg := shared.Message{Type: shared.MSGTasksCreate, Data: &task}
	msg.Send(conn)

	response, err := shared.RecvMessage[interface{}](conn)
	log.Printf("Task: %+v", response)
	if err != nil {
		log.Printf("ERROR: could not read message: %v", err)
		return err
	}

	return nil
}

func main() {
	listener := shared.NewTcpListener("8000")
	conn, err := listener.Run()
	if err != nil {
		log.Fatal(err)
	}
	job, err := createJob(conn)
	if err != nil {
		return
	}

	for range 5 {
		createTask(conn, job.ID)
	}
}
