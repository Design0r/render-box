package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Task struct {
	Name string
}

type PageData struct {
	Test  string
	Tasks []Task
}

func HandleIndex(c echo.Context) error {
	tasks := []Task{{Name: "Render 01"}, {Name: "Render 02"}}
	ctx := PageData{Test: "Hello World", Tasks: tasks}
	c.Render(http.StatusOK, "index.html", ctx)

	return nil
}
