package api

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"render-box/shared"
	"render-box/shared/db/repo"
)

type PageData struct {
	Tasks   []repo.Task
	Jobs    []repo.Job
	Workers []repo.Worker
}

func HandleIndex(c echo.Context) error {
	listener := shared.NewTcpListener("8000")
	conn, err := listener.Run()
	defer (*conn).Close()
	if err != nil {
		return err
	}
	ctx, err := fetchPageData(conn, int64(0))
	if err != nil {
		return err
	}
	c.Render(http.StatusOK, "index.html", ctx)

	return nil
}
