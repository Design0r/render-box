package main

import (
	"log"
	"time"

	"github.com/google/uuid"

	"render-box/shared"
	"render-box/shared/db/repo"
)

func main() {
	listener := shared.NewTcpListener("8000")
	conn, err := listener.Run()
	if err != nil {
		log.Fatal(err)
	}
	defer (*conn).Close()

	// name, err := os.Hostname()
	randName, err := uuid.NewRandom()
	name := randName.String()
	if err != nil {
		log.Fatal(err)
	}

	msg := shared.Message{Type: shared.MSGWorkerRegister, Data: name}
	msg.Send(conn)

	response, err := shared.RecvMessage[any](conn)
	log.Println(response)
	if err != nil {
		log.Printf("ERROR: could not read message: %v", err)
	}

	msg = shared.Message{Type: shared.MSGTasksNext, Data: nil}
	complete := shared.Message{Type: shared.MSGTasksComplete, Data: nil}
	for {
		msg.Send(conn)
		response, err := shared.RecvMessage[repo.Task](conn)
		if err != nil {
			log.Printf("Failed getting task: %v", err)
			break
		}
		if response.Type == shared.MSGError {
			log.Printf("No Task: %v", err)
			time.Sleep(time.Second * 2)
			continue
		}

		task := (response.Data).(repo.Task)
		log.Printf("Got Task: %+v", task)
		time.Sleep(time.Second * 5)

		complete.Send(conn)
		response, err = shared.RecvMessage[repo.Task](conn)
		if err != nil {
			log.Printf("Failed completing task: %v", err)
			break
		}
	}
}
