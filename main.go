package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	// "github.com/Okenamay/final/database"

	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

func dbInit() *sql.DB {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if install {
		query := ` 
		CREATE TABLE IF NOT EXISTS scheduler (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            date TEXT NOT NULL,
            title TEXT NOT NULL,
            comment TEXT,
            repeat TEXT NOT NULL CHECK(length(repeat) <= 128)
        );
        CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
		`
		_, err = db.Exec(query)
		if err != nil {
			panic(err)
		}
		log.Println("Database creation successful!")
	}
	return db
}

func main() {

	dbInit()

	r := chi.NewRouter()
	fmt.Println("Starting server at port 7540")

	r.Handle("/*", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(":7540", r)
	if err != nil {
		panic(err)
	}

}
