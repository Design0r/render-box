package api

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"render-box/client/assets/templates"
	"render-box/client/schemas"
	"render-box/shared"
	"render-box/shared/db/repo"
)

var upgrader = websocket.Upgrader{}

func HandleClientActions(done chan struct{}, jobChan chan int64, ws *websocket.Conn) {
	defer close(done)
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			return
		}

		var msg map[string]interface{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("JSON Unmarshal error:", err)
			continue
		}

		log.Printf("Client Action: %+v", msg)

		if jobID, ok := msg["job_id"]; ok {
			id, err := strconv.Atoi(jobID.(string))
			if err != nil {
				continue
			}
			jobChan <- int64(id)
		}
	}
}

func fetchPageData(conn *net.Conn, jobId int64) (*schemas.PageData, error) {
	taskMsg := shared.Message{Type: shared.MSGTasksAll, Data: nil}
	taskMsg.Send(conn)
	response, err := shared.RecvMessage[[]repo.Task](conn)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	tasks := ((*response).Data).([]repo.Task)

	jobMsg := shared.Message{Type: shared.MSGJobsAll, Data: nil}
	jobMsg.Send(conn)
	response, err = shared.RecvMessage[[]repo.Job](conn)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	jobs := ((*response).Data).([]repo.Job)

	workerMsg := shared.Message{Type: shared.MSGWorkerAll, Data: nil}
	workerMsg.Send(conn)
	response, err = shared.RecvMessage[[]repo.Worker](conn)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	worker := ((*response).Data).([]repo.Worker)

	var activeTasks []repo.Task
	for _, t := range tasks {
		if t.JobID != jobId {
			continue
		}
		activeTasks = append(activeTasks, t)
	}

	return &schemas.PageData{Tasks: activeTasks, Jobs: jobs, Workers: worker}, nil
}

func updatePage(ws *websocket.Conn, conn *net.Conn, currentJobId int64) error {
	var buf bytes.Buffer
	ctx, err := fetchPageData(conn, currentJobId)
	if err != nil {
		return err
	}

	err = templates.Update(ctx).Render(context.Background(), &buf)
	if err != nil {
		return err
	}

	err = ws.WriteMessage(websocket.TextMessage, buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func WsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	var currentJobId int64

	done := make(chan struct{})
	jobChan := make(chan int64)
	go HandleClientActions(done, jobChan, ws)

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	listener := shared.NewTcpListener("8000")
	conn, err := listener.Run()
	defer (*conn).Close()
	if err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			updatePage(ws, conn, currentJobId)
		case <-done:
			return nil
		case id := <-jobChan:
			currentJobId = id
			updatePage(ws, conn, currentJobId)
		}
	}
}
