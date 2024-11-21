package api

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

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

func fetchPageData(conn *net.Conn, jobId int64) (*PageData, error) {
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

	var activeTasks []repo.Task
	for _, t := range tasks {
		if t.JobID != jobId {
			continue
		}
		activeTasks = append(activeTasks, t)
	}

	return &PageData{Tasks: activeTasks, Jobs: jobs}, nil
}

func updatePage(c echo.Context, ws *websocket.Conn, conn *net.Conn, currentJobId int64) error {
	var buf bytes.Buffer
	ctx, err := fetchPageData(conn, currentJobId)
	if err != nil {
		return err
	}

	err = c.Echo().Renderer.Render(&buf, "update", ctx, c)
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

	for {
		select {
		case <-ticker.C:
			updatePage(c, ws, conn, currentJobId)
		case <-done:
			return nil
		case id := <-jobChan:
			currentJobId = id
			updatePage(c, ws, conn, currentJobId)
		}
	}
}
