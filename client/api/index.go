package api

import (
	"context"

	"github.com/labstack/echo/v4"

	"render-box/client/assets/templates"
	"render-box/shared"
)

func HandleIndex(c echo.Context) error {
	listener := shared.NewTcpListener("8000")
	conn, err := listener.Run()
	defer (*conn).Close()
	if err != nil {
		return err
	}
	ctx, err := fetchPageData(conn, int64(1))
	if err != nil {
		return err
	}

	templates.Index(ctx).Render(context.Background(), c.Response().Writer)

	return nil
}
