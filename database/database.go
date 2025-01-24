package database

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"final/daterules"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile = "scheduler.db"
	limit  = 50
)

func DBInit() *sql.DB {
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
		log.Println("Database creation success!")
	}

	return db
}

func AddEntry(db *sql.DB, task daterules.Task) (string, error) {
	AddEntry := `INSERT INTO scheduler (date, title, comment, repeat) 
	VALUES (:date, :title, :comment, :repeat)`
	result, err := db.Exec(AddEntry,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return "", err
	}
	idb, err := result.LastInsertId()
	if err != nil {
		return "", err
	}

	return strconv.Itoa(int(idb)), nil
}

func DeleteEntry(db *sql.DB, id string) error {
	DeleteEntry := `DELETE FROM scheduler WHERE id = ?`
	result, err := db.Exec(DeleteEntry, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("wrong row id")
	}

	return nil
}

func EditEntry(db *sql.DB, task daterules.Task) error {
	EditEntry := `UPDATE scheduler 
	SET date = ?, title = ?, comment = ?, repeat = ? 
	WHERE id = ?;
	`
	result, err := db.Exec(EditEntry,
		task.Date,
		task.Title,
		task.Comment,
		task.Repeat,
		task.ID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("wrong row id")
	}

	return nil
}

func GetEntry(db *sql.DB, id string) (daterules.Task, error) {
	GetEntry := `SELECT id, date, title, comment, repeat 
	FROM scheduler WHERE id = ?`
	row := db.QueryRow(GetEntry, id)
	var entry daterules.Task
	err := row.Scan(&entry.ID, &entry.Date, &entry.Title, &entry.Comment, &entry.Repeat)
	if err != nil {
		return daterules.Task{}, err
	}

	return entry, nil
}

func GetAllEntries(db *sql.DB) ([]daterules.Task, error) {
	GetAllEntries := `SELECT id, date, title, comment, repeat 
	FROM scheduler 
	WHERE date >= strftime('%Y %m %d', 'now') 
	ORDER BY date ASC 
	LIMIT ?;
	`
	rows, err := db.Query(GetAllEntries, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := []daterules.Task{}
	for rows.Next() {
		var entry daterules.Task
		err := rows.Scan(&entry.ID, &entry.Date, &entry.Title, &entry.Comment, &entry.Repeat)
		if err != nil {
			return nil, errors.New("database error")
		}
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
