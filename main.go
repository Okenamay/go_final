package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"go_final/database"
	"go_final/handler"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	r := chi.NewRouter()
	database.DBInit()

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	handler := handler.NewHandler(db)

	fmt.Println("Starting server at port 7540")

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.Get("/api/task", handler.GetTask)
	r.Post("/api/task", handler.AddTask)
	r.Put("/api/task", handler.EditTask)
	r.Delete("/api/task", handler.DeleteTask)
	r.Get("/api/tasks", handler.GetAllTasks)
	r.Get("/api/nextdate", handler.GetNextDate)
	r.Post("/api/task/done", handler.SetTaskDone)

	err = http.ListenAndServe(":7540", r)
	if err != nil {
		fmt.Printf("Server failed to start: %s\n", err)
	}

}
