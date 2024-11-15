package main

import (
	"log"
	"net"

	"render-box/shared"
	"render-box/shared/db/repo"
)

func handleRead(conn *net.Conn) (*shared.Message, error) {
	header := make([]byte, 4)
	bodySize, err := shared.GetBodySize(conn, header)
	if err != nil {
		return nil, err
	}
	body, err := shared.ReadBody(conn, bodySize)
	if err != nil {
		return nil, err
	}

	log.Printf("MESSAGE: %v\n", body)
	return body, nil
}

func main() {
	listener := shared.NewTcpListener("8000")
	conn, err := listener.Run()
	if err != nil {
		log.Fatal(err)
	}
	task := repo.CreateJobParams{Name: "TestJob", Priority: int64(50), State: "waiting"}
	msg := shared.Message{Type: shared.MSGJobsCreate, Data: &task}
	msg.Send(conn)

	_, err = handleRead(conn)
	if err != nil {
		log.Printf("ERROR: could not read message: %v", err)
	}
}
