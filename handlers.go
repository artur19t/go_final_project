package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"go-final/db"
	"go-final/model"
)

// Хендлер добавления задачи
func postTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err1 := model.SetError{Error: err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		err1 := model.SetError{Error: err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	strErr, haveErr, task := db.PostTaskDB(task)
	if haveErr {
		err1 := model.SetError{Error: errors.New(strErr)}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	num, err := db.CreateTask(task)
	if num == 0 && err != nil {
		err1 := model.SetError{Error: err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	setId := model.SetTask{ID: int(num)}
	resp, _ := json.Marshal(setId)
	w.Write(resp)
	w.Header().Set("Content-Type", "application/json")
}

// Хендлер для отправки данных о задаче
func getTask(w http.ResponseWriter, r *http.Request) {
	idSearch := r.URL.Query().Get("id")
	next, resp, num := db.CheckID(idSearch)
	if !next {
		w.Write(resp)
		return
	}

	errStr, haveErr, task := db.FindTaskID(num, idSearch)
	if haveErr {
		err1 := model.SetError{Error: errors.New(errStr)}
		resp, _ = json.Marshal(err1)
		w.Write(resp)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		err1 := model.SetError{Error: errors.New(err.Error())}
		resp, _ = json.Marshal(err1)
		w.Write(resp)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Хендлер изменения параметров задачи
func correctTask(w http.ResponseWriter, r *http.Request) {
	var task model.TaskGet
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err1 := model.SetError{Error: err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		err1 := model.SetError{Error: err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}

	strErr, haveErr := db.CheckCorrectIn(task)
	if haveErr {
		err1 := model.SetError{Error: errors.New(strErr)}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}

	strErr, haveErr = db.CorrectTaskDB(task)
	if haveErr {
		err1 := model.SetError{Error: errors.New(strErr)}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	if !haveErr {
		space := model.Space{}
		resp, _ := json.Marshal(space)
		w.Write(resp)
	}
}

// Хендлер запроса задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.GetTaskDB()
	if err != nil {
		return
	}

	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Хендлер выполнения задачи
func doneTask(w http.ResponseWriter, r *http.Request) {
	idSearch := r.URL.Query().Get("id")
	next, resp, num := db.CheckID(idSearch)
	if !next {
		w.Write(resp)
		return
	}

	errText, err := db.DoneTaskDB(num)
	if err != nil {
		err1 := model.SetError{Error: errors.New(errText)}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	space := model.Space{}
	resp, _ = json.Marshal(space)
	w.Write(resp)
}

// Хендлер удаления задачи
func deleteTask(w http.ResponseWriter, r *http.Request) {
	idSearch := r.URL.Query().Get("id")
	next, resp, num := db.CheckID(idSearch)
	if !next {
		w.Write(resp)
		return
	}
	errText, err := db.DeleteTaskDB(num)
	if err != nil {
		err1 := model.SetError{Error: errors.New(errText)}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	space := model.Space{}
	resp, _ = json.Marshal(space)
	w.Write(resp)
}
