package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"go_final/database"
	"go_final/daterules"

	_ "github.com/mattn/go-sqlite3"
)

type Id struct {
	Id string `json:"id"`
}

type Err struct {
	Error string `json:"error"`
}

func (err Err) Bytes() []byte {
	data, _ := json.Marshal(err)
	return data
}

type TasksRes struct {
	Tasks []daterules.Task `json:"tasks"`
}

var (
	TimeFormat string = daterules.TimeFormat
)

type Handler struct {
	DB *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		DB: db,
	}
}

func taskFromReq(bytes []byte) (daterules.Task, error) {
	var t daterules.Task
	err := json.Unmarshal(bytes, &t)
	if err != nil {
		return daterules.Task{}, err
	}

	if t.Title == "" {
		return daterules.Task{}, errors.New("no title specified")
	}

	y, m, d := time.Now().Date()
	nowDate := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	var taskDate time.Time
	if t.Date == "" || t.Date == "today" {
		taskDate = nowDate
	} else {
		taskDate, err = time.Parse(TimeFormat, t.Date)
		if err != nil {
			return daterules.Task{}, err
		}

		for taskDate.Before(nowDate) {
			taskDate, err = daterules.NextTime(nowDate, "", t.Repeat)
			if err != nil {
				return daterules.Task{}, err
			}
		}
	}

	t.Date = taskDate.Format(daterules.TimeFormat)

	return t, nil
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.Write(Err{Error: "no id specified"}.Bytes())
		return
	}

	task, err := database.GetEntry(h.DB, id)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	response, err := json.Marshal(task)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *Handler) AddTask(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	t, err := taskFromReq(bytes)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	id, err := database.AddEntry(h.DB, t)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	res := Id{
		Id: id,
	}

	response, err := json.Marshal(res)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (h *Handler) EditTask(w http.ResponseWriter, r *http.Request) {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	t, err := taskFromReq(bytes)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	err = database.EditEntry(h.DB, t)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.Write(Err{Error: "task id is not set"}.Bytes())
		return
	}

	err := database.DeleteEntry(h.DB, id)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

func (h *Handler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := database.GetAllEntries(h.DB)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	res := TasksRes{
		Tasks: tasks,
	}

	resp, err := json.Marshal(res)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *Handler) GetNextDate(w http.ResponseWriter, r *http.Request) {
	now := r.URL.Query().Get("now")
	if now == "" {
		w.Write([]byte{})
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		w.Write([]byte{})
		return
	}

	repeat := r.URL.Query().Get("repeat")
	if repeat == "" {
		w.Write([]byte{})
		return
	}

	taskDate, err := time.Parse(TimeFormat, date)
	if err != nil {
		w.Write([]byte{})
		return
	}

	nowDate, err := time.Parse(TimeFormat, now)
	if err != nil {
		w.Write([]byte{})
		return
	}

	resultT, err := daterules.NextTime(taskDate, "", repeat)
	if err != nil {
		w.Write([]byte{})
		return
	}

	for taskDate.Before(nowDate) {
		taskDate, err = daterules.NextTime(taskDate, "", repeat)
		if err != nil {
			w.Write([]byte{})
			return
		}

		resultT = taskDate
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resultT.Format(TimeFormat)))
}

func (h *Handler) SetTaskDone(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.Write(Err{Error: "task id is not set"}.Bytes())
		return
	}

	task, err := database.GetEntry(h.DB, id)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	if task.Repeat == "" {
		err := database.DeleteEntry(h.DB, id)
		if err != nil {
			w.Write(Err{Error: err.Error()}.Bytes())
			return
		}

		w.Write([]byte("{}"))
		return
	}

	task.Date, err = task.NextDay()
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}
	err = database.EditEntry(h.DB, task)
	if err != nil {
		w.Write(Err{Error: err.Error()}.Bytes())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}
