package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
type SetTask struct {
	ID int `json:"id"`
}

type SetError struct {
	Error error `json:"error"`
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	today := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = today
	}

	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if task.Date < today && task.Repeat == "" {
		task.Date = today
	}
	if task.Date < today && task.Repeat != "" {
		task.Date, _ = NextDate(time.Now(), task.Date, task.Repeat)
		if task.Date == "-1" {
			err1 := SetError{errors.New("не верно указано правило повторения")}
			resp, _ := json.Marshal(err1)
			w.Write(resp)
			//http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
	if task.Title == "" {
		err = errors.New("не указан заголовок задачи")
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dateToTable, _ := NextDate(time.Now(), task.Date, task.Repeat)
	if dateToTable == "-1" {
		err := errors.New("не верно указано правило повторения")
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		//fmt.Println(err)
		return
	}

	id, _ := res.LastInsertId()
	setId := SetTask{int(id)}
	resp, _ := json.Marshal(setId)
	w.Write(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

type FindTask struct {
	ID string `json:"id"`
}
type StringValue struct {
	Error string `json:"error"`
}

func fillTask(w http.ResponseWriter, r *http.Request) {
	idSearch := r.URL.Query().Get("id")
	if idSearch == "" {
		err := StringValue{}
		err.Error = "не указан идентификатор"
		resp, _ := json.Marshal(err)
		w.Write(resp)
		return
	}

	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	num, err := strconv.Atoi(idSearch)
	if err != nil {
		err := StringValue{}
		err.Error = "не верно указан идентификатор"
		resp, _ := json.Marshal(err)
		w.Write(resp)
		return
	}
	row := db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", num))
	err = row.Scan(&num)
	if err != nil {
		err1 := SetError{errors.New("не верно указан идентификатор")}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	task := TaskGet{}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM scheduler WHERE id = %s", idSearch))
	if err != nil {
		err1 := SetError{errors.New("задача не найдена")}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	defer rows.Close()
	for rows.Next() {
		id := 0
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			return
		}
		task.ID = fmt.Sprint(id)
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

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

	today := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = today
	}

	_, err = time.Parse("20060102", task.Date)
	if err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if task.ID == "" {
		err := StringValue{}
		err.Error = "не указан идентификатор"
		resp, _ := json.Marshal(err)
		w.Write(resp)
		return
	}
	if task.Date < today && task.Repeat == "" {
		task.Date = today
	}
	if task.Date < today && task.Repeat != "" {
		task.Date, _ = NextDate(time.Now(), task.Date, task.Repeat)
		if task.Date == "-1" {
			err1 := SetError{errors.New("не верно указано правило повторения")}
			resp, _ := json.Marshal(err1)
			w.Write(resp)
			//http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	}
	if task.Title == "" {
		err = errors.New("не указан заголовок задачи")
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dateToTable, _ := NextDate(time.Now(), task.Date, task.Repeat)
	if dateToTable == "-1" {
		err := errors.New("не верно указано правило повторения")
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		//http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	num, err := strconv.Atoi(task.ID)
	if err != nil {
		err := StringValue{}
		err.Error = "не верно указан идентификатор"
		resp, _ := json.Marshal(err)
		w.Write(resp)
		return
	}
	row := db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", num))
	err = row.Scan(&num)
	if err != nil {
		err1 := SetError{errors.New("не верно указан идентификатор")}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", num))

	space := Space{}
	resp, _ := json.Marshal(space)
	w.Write(resp)
}

type Space struct {
}

type TaskGet struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}
type SliceTask struct {
	Tasks []TaskGet `json:"tasks"`
}

func getTask(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		return
	}
	defer db.Close()
	tasks := SliceTask{}
	sliceTask := []TaskGet{}
	rows, err := db.Query("SELECT * FROM scheduler")
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		task := TaskGet{}
		id := 0
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return
		}
		task.ID = fmt.Sprint(id)
		sliceTask = append(sliceTask, task)
	}
	tasks.Tasks = sliceTask
	//fmt.Println(sliceTask)
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	//fmt.Println(resp)
	w.Write(resp)
}

func doneTask(w http.ResponseWriter, r *http.Request) {
	task := TaskGet{}
	idSearch := r.URL.Query().Get("id")
	if idSearch == "" {
		err := StringValue{}
		err.Error = "не указан идентификатор"
		resp, _ := json.Marshal(err)
		w.Write(resp)
		return
	}

	num, err := strconv.Atoi(idSearch)
	if err != nil {
		err := StringValue{}
		err.Error = "не верно указан идентификатор"
		resp, _ := json.Marshal(err)
		w.Write(resp)
		return
	}

	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", num))
	err = row.Scan(&num, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		err1 := SetError{errors.New("не верно указан идентификатор")}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}

	if task.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", num))
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		space := Space{}
		resp, _ := json.Marshal(space)
		w.Write(resp)
		return
	} else {
		date, _ := NextDate(time.Now(), task.Date, task.Repeat)
		if date == "-1" {
			err1 := SetError{errors.New("ошибка переноса даты")}
			resp, _ := json.Marshal(err1)
			w.Write(resp)
			return
		}
		db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
			sql.Named("date", date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat),
			sql.Named("id", num))

		space := Space{}
		resp, _ := json.Marshal(space)
		w.Write(resp)
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	task := TaskGet{}
	idSearch := r.URL.Query().Get("id")
	if idSearch == "" {
		err := StringValue{}
		err.Error = "не указан идентификатор"
		resp, _ := json.Marshal(err)
		w.Write(resp)
		return
	}

	num, err := strconv.Atoi(idSearch)
	if err != nil {
		err := StringValue{}
		err.Error = "не верно указан идентификатор"
		resp, _ := json.Marshal(err)
		w.Write(resp)
		return
	}

	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", num))
	err = row.Scan(&num, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		err1 := SetError{errors.New("не верно указан идентификатор")}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return
	}
	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", num))
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	space := Space{}
	resp, _ := json.Marshal(space)
	w.Write(resp)
}
