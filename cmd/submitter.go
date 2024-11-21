package main

import (
	"log"

	"render-box/shared"
	"render-box/shared/db/repo"
)

func main() {
	listener := shared.NewTcpListener("8000")
	conn, err := listener.Run()
	if err != nil {
		log.Fatal(err)
	}
	task := repo.CreateJobParams{Name: "TestJob", Priority: int64(50), State: "waiting"}
	msg := shared.Message{Type: shared.MSGJobsCreate, Data: &task}
	msg.Send(conn)

	response, err := shared.RecvMessage[any](conn)
	log.Println(response)
	if err != nil {
		log.Printf("ERROR: could not read message: %v", err)
	}
}
