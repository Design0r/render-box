package api

import (
	"bytes"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"

	"render-box/shared"
	"render-box/shared/db/repo"
)

var upgrader = websocket.Upgrader{}

func WsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	listener := shared.NewTcpListener("8000")
	conn, err := listener.Run()
	taskMsg := shared.Message{Type: shared.MSGTasksAll, Data: nil}
	jobMsg := shared.Message{Type: shared.MSGJobsAll, Data: nil}

	var buf bytes.Buffer
	for {
		select {
		case <-ticker.C:
			taskMsg.Send(conn)
			response, err := shared.RecvMessage[[]repo.Task](conn)
			if err != nil {
				c.Logger().Error(err)
				return err
			}
			tasks := ((*response).Data).([]repo.Task)

			jobMsg.Send(conn)
			response, err = shared.RecvMessage[[]repo.Job](conn)
			if err != nil {
				c.Logger().Error(err)
				return err
			}
			jobs := ((*response).Data).([]repo.Job)
			ctx := PageData{Tasks: tasks, Jobs: jobs}

			err = c.Echo().Renderer.Render(&buf, "update", ctx, c)
			if err != nil {
				c.Logger().Error(err)
				return err
			}

			err = ws.WriteMessage(websocket.TextMessage, buf.Bytes())
			if err != nil {
				c.Logger().Error(err)
				return err
			}
		}
	}
}
