package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"render-box/server"
	"render-box/server/routes"
	"render-box/shared"
)

func main() {
	db, err := sql.Open("sqlite3", "render_box.db")
	if err != nil {
		log.Fatal("Failed to open database: %v", err)
	}
	defer db.Close()

	router := shared.NewMessageRouter()
	router.IncludeRouter(routes.InitJobRouter())
	router.IncludeRouter(routes.InitWorkerRouter())
	router.IncludeRouter(routes.InitTaskRouter())

	Port := "8000"
	server := server.NewServer(Port, db, router)
	server.Run()
}
