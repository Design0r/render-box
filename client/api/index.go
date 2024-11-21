package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"render-box/shared/db/repo"
)

type Task struct {
	Name string
}

type PageData struct {
	Tasks []repo.Task
	Jobs  []repo.Job
}

func HandleIndex(c echo.Context) error {
	ctx := PageData{Tasks: nil, Jobs: nil}
	c.Render(http.StatusOK, "index.html", ctx)

	return nil
}
