package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"render-box/server"
)

func main() {
	db, err := sql.Open("sqlite3", "render_box.db")
	if err != nil {
		log.Fatal("Failed to open database: %v", err)
	}
	defer db.Close()

	Port := "8000"
	server := server.NewServer(Port, db)
	server.Run()
}
