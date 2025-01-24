package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"final/database"
	"final/handler"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

func main() {

	r := chi.NewRouter()
	database.DBInit()

	db, err := sql.Open("sqlite3", "scheduler.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := database.NewContainer(db)
	service := handler.NewTaskService(store)

	fmt.Println("Starting server at port 7540")

	r.Handle("/*", http.FileServer(http.Dir("./web")))
	r.HandleFunc("/api/task/done", service.DoneTask)
	r.HandleFunc("/api/task", service.Task)
	r.HandleFunc("/api/nextdate", handler.NextDeadLine)
	r.HandleFunc("/api/tasks", service.GetTasks)

	err = http.ListenAndServe(":7540", r)
	if err != nil {
		panic(err)
	}

}
