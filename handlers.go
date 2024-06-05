package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Хендлер добавления задачи
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var task1 TaskGet
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	strErr := checkCorrectIn(task1, w)
	if strErr == "error" {
		return
	}
	strErr = writeTask(task, w)
	if strErr == "error" {
		return
	}
}

// Хендлер для отправки данных о задаче
func getTask(w http.ResponseWriter, r *http.Request) {
	idSearch := r.URL.Query().Get("id")
	next, resp, num := checkID(idSearch)
	if !next {
		w.Write(resp)
		return
	}

	errStr := findTaskID(num, idSearch, w)
	if errStr == "error" {
		return
	}
}

// Хендлер изменения параметров задачи
func correctTask(w http.ResponseWriter, r *http.Request) {
	var task TaskGet
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}

	errStr := checkCorrectIn(task, w)
	if errStr == "error" {
		return
	}

	resp, err := correctTaskDB(task)
	if resp != nil {
		w.Write(resp)
	}
	if err != nil {
		return
	}
}

// Хендлер запроса задач
func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := getTaskDB()
	if err != nil {
		if resp != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Хендлер выполнения задачи
func doneTask(w http.ResponseWriter, r *http.Request) {

	idSearch := r.URL.Query().Get("id")
	next, resp, num := checkID(idSearch)
	if !next {
		w.Write(resp)
		return
	}
	resp, err := doneTaskDB(num)
	w.Write(resp)
	if err != nil {
		return
	}
}

// Хендлер удаления задачи
func deleteTask(w http.ResponseWriter, r *http.Request) {
	idSearch := r.URL.Query().Get("id")
	next, resp, num := checkID(idSearch)
	if !next {
		w.Write(resp)
		return
	}
	resp, err := deleteTaskDB(num)
	if err != nil {
		w.Write(resp)
	}
	space := Space{}
	resp, _ = json.Marshal(space)
	w.Write(resp)
}
