package api

import (
	"bytes"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
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

	var buf bytes.Buffer
	for {
		select {
		case <-ticker.C:
			id, _ := uuid.NewRandom()
			str_id := id.String()
			ctx := PageData{Test: "Hello", Tasks: []Task{{Name: str_id}, {Name: str_id}}}

			err := c.Echo().Renderer.Render(&buf, "tasks", ctx, c)
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
