package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"render-box/client/api"
	"render-box/client/assets"
)

func main() {
	e := echo.New()
	e.Renderer = assets.NewTemplate()
	apiV1 := e.Group("")


	server := api.NewServer(e)
	server.InitMiddleware()
	server.InitRoutes(apiV1)

	log.Fatal(server.Start(":8080"))
}
