package main

import (
	"log"
	"os"
	"time"

	"render-box/shared"
	"render-box/shared/db/repo"
)

func main() {
	listener := shared.NewTcpListener("8000")
	conn, err := listener.Run()
	if err != nil {
		log.Fatal(err)
	}
	name, err := os.Hostname()
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
	for {
		msg.Send(conn)
		response, err := shared.RecvMessage[repo.Task](conn)
		if err != nil {
			log.Printf("Failed getting task: %v", err)
			break
		}

		task := (response.Data).(repo.Task)
		log.Printf("Got Task: %+v", task)
		time.Sleep(time.Second * 10)
	}
}
