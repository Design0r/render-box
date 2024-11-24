package server

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/mitchellh/mapstructure"

	"render-box/server/service"
	"render-box/shared"
	"render-box/shared/db/repo"
)

type Server struct {
	Addr     string
	Listener *net.Listener
	Db       *sql.DB
}

func NewServer(port string, db *sql.DB) *Server {
	return &Server{Addr: (":" + port), Listener: nil, Db: db}
}

func (self *Server) Run() {
	l, err := net.Listen("tcp", self.Addr)
	if err != nil {
		log.Fatalf("ERROR: Could not listen to port %s: %s\n", self.Addr, err)
	}
	self.Listener = &l
	log.Printf("TCP Server running on %s\n", self.Addr)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("ERROR: Could not accept connection: %s\n", err)
			continue
		}
		log.Printf("New connection with %s\n", conn.RemoteAddr().String())

		go handleConnection(&conn, self.Db)
	}
}

type ConnState struct {
	Worker *repo.Worker
	Task   *repo.Task
}

func handleConnection(conn *net.Conn, db *sql.DB) {
	c := *conn
	defer c.Close()

	state := &ConnState{}

	header := make([]byte, 4)
	for {
		bodySize, err := shared.GetBodySize(conn, header)
		if err != nil {
			break
		}
		body, err := shared.ReadBody(conn, bodySize)
		if err != nil {
			break
		}

		returnData, err := handleMessage(db, body, state)
		var returnMsg shared.Message
		if err != nil {
			returnMsg = shared.Message{Type: shared.MSGError, Data: nil}
		} else {
			returnMsg = shared.Message{Type: shared.MSGSuccess, Data: returnData}
		}
		log.Printf("%+v\n", returnMsg)

		if sendErr := returnMsg.Send(conn); sendErr != nil {
			log.Printf("ERROR: Could not send response to client: %v\n", sendErr)
			break
		}

	}

	if state.Worker != nil {
		service.UpdateWorkerState(db, "offline", state.Worker.ID)
	}
	if state.Task != nil {
		service.UpdateTaskState(db, "waiting", state.Task.ID)
	}

	log.Printf("Closed connection with %s\n", c.RemoteAddr().String())
}

func handleMessage(
	db *sql.DB,
	message *shared.Message,
	state *ConnState,
) (interface{}, error) {
	log.Printf("MESSAGE: %+v\n", message)
	switch message.Type {
	case shared.MSGTasksCreate:
		data := &repo.CreateTaskParams{}
		mapstructure.Decode(message.Data, data)
		task, err := service.CreateTask(db, data)
		if err != nil {
			return nil, err
		}
		return task, nil
	case shared.MSGTasksAll:
		tasks, err := service.GetTasks(db)
		if err != nil {
			return nil, err
		}
		return tasks, nil
	case shared.MSGJobsCreate:
		data := &repo.CreateJobParams{}
		mapstructure.Decode(message.Data, data)
		task, err := service.CreateJob(db, data)
		if err != nil {
			return nil, err
		}
		return task, nil
	case shared.MSGJobsAll:
		tasks, err := service.GetJobs(db)
		if err != nil {
			return nil, err
		}
		return tasks, nil
	case shared.MSGWorkerCreate:
		data := &repo.CreateWorkerParams{}
		mapstructure.Decode(message.Data, data)
		worker, err := service.CreateWorker(db, data)
		if err != nil {
			return nil, err
		}
		return worker, nil
	case shared.MSGWorkerAll:
		worker, err := service.GetWorkers(db)
		if err != nil {
			return nil, err
		}
		return worker, nil
	case shared.MSGWorkerRegister:
		name := (message.Data).(string)
		worker, err := service.RegisterWorker(db, name)
		if err != nil {
			return nil, err
		}
		state.Worker = worker
		return worker, err

	case shared.MSGTasksNext:
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

	default:
		return nil, fmt.Errorf("Invalid message type: %v", message.Type)
	}
}
