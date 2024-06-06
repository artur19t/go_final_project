package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func postTaskDB(task Task, w http.ResponseWriter) string {
	today := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = today
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return "error"
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
			return "error"
		}

	}
	if task.Title == "" {
		err = errors.New("не указан заголовок задачи")
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return "error"
	}
	dateToTable, _ := NextDate(time.Now(), task.Date, task.Repeat)
	if dateToTable == "-1" {
		err := errors.New("не верно указано правило повторения")
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return "error"
	}
	resp := writeTask(task, w)
	if resp[0] == 'e' {
		return "error"
	}
	return ""
}

// Запись задачи в БД
func writeTask(task Task, w http.ResponseWriter) []byte {
	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return []byte{'e'}
	}
	defer db.Close()

	res, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		//fmt.Println(err)
		return []byte{'e'}
	}

	id, _ := res.LastInsertId()
	setId := SetTask{int(id)}
	resp, _ := json.Marshal(setId)
	w.Write(resp)
	w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusCreated)
	return resp
}

// Проверка корректности входящих данных
func checkCorrectIn(task TaskGet, w http.ResponseWriter) string {
	today := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = today
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return "error"
	}
	if task.ID == "" {
		err := StringValue{}
		err.Error = "не указан идентификатор"
		resp, _ := json.Marshal(err)
		w.Write(resp)
		return "error"
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
			return "error"
		}

	}
	if task.Title == "" {
		err = errors.New("не указан заголовок задачи")
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return "error"
	}
	dateToTable, _ := NextDate(time.Now(), task.Date, task.Repeat)
	if dateToTable == "-1" {
		err := errors.New("не верно указано правило повторения")
		err1 := SetError{err}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return "error"
	}
	return ""
}

// Поиск задачи по ID
func findTaskID(num int, idSearch string, w http.ResponseWriter) string {
	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		log.Println(err)
		return "error"
	}
	defer db.Close()
	row := db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", num))
	err = row.Scan(&num)
	if err != nil {
		err1 := SetError{errors.New("не верно указан идентификатор")}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return "error"
	}
	task := TaskGet{}
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM scheduler WHERE id = %s", idSearch))
	if err != nil {
		err1 := SetError{errors.New("задача не найдена")}
		resp, _ := json.Marshal(err1)
		w.Write(resp)
		return "error"
	}
	defer rows.Close()
	for rows.Next() {
		id := 0
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			return "error"
		}
		task.ID = fmt.Sprint(id)
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return "error"
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
	return ""
}

// Изменение параметров задачи
func correctTaskDB(task TaskGet) ([]byte, error) {

	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer db.Close()

	num, err := strconv.Atoi(task.ID)
	if err != nil {
		err1 := StringValue{}
		err1.Error = "не верно указан идентификатор"
		resp, _ := json.Marshal(err1)
		return resp, err
	}
	row := db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", num))
	err = row.Scan(&num)
	if err != nil {
		err1 := SetError{errors.New("не верно указан идентификатор")}
		resp, _ := json.Marshal(err1)
		return resp, err
	}
	db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", num))

	space := Space{}
	resp, _ := json.Marshal(space)
	return resp, nil
}

// Проверка корректности введеного ID
func checkID(idSearch string) (bool, []byte, int) {
	if idSearch == "" {
		err := StringValue{}
		err.Error = "не указан идентификатор"
		resp, _ := json.Marshal(err)
		return false, resp, 0
	}

	num, err := strconv.Atoi(idSearch)
	if err != nil {
		err := StringValue{}
		err.Error = "не верно указан идентификатор"
		resp, _ := json.Marshal(err)

		return false, resp, 0
	}
	return true, nil, num
}

// Запрос задач
func getTaskDB() ([]byte, error) {
	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		return []byte{}, err
	}
	defer db.Close()
	tasks := SliceTask{}
	sliceTask := []TaskGet{}
	rows, err := db.Query("SELECT * FROM scheduler")
	if err != nil {
		return []byte{}, err
	}
	defer rows.Close()
	for rows.Next() {
		task := TaskGet{}
		id := 0
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return []byte{}, err
		}
		task.ID = fmt.Sprint(id)
		sliceTask = append(sliceTask, task)
	}
	tasks.Tasks = sliceTask
	//fmt.Println(sliceTask)
	resp, err := json.Marshal(tasks)
	if err != nil {
		return []byte{'e'}, err
	}
	return resp, nil
}

// Выполнение задачи
func doneTaskDB(num int) ([]byte, error) {
	task := TaskGet{}
	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", num))
	err = row.Scan(&num, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		err1 := SetError{errors.New("не верно указан идентификатор")}
		resp, _ := json.Marshal(err1)
		return resp, err
	}

	if task.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", num))
		if err != nil {
			return []byte(err.Error()), err
		}
		space := Space{}
		resp, _ := json.Marshal(space)
		return resp, nil
	} else {
		date, _ := NextDate(time.Now(), task.Date, task.Repeat)
		if date == "-1" {
			err1 := SetError{errors.New("ошибка переноса даты")}
			resp, _ := json.Marshal(err1)
			return resp, err
		}
		db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
			sql.Named("date", date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat),
			sql.Named("id", num))

		space := Space{}
		resp, _ := json.Marshal(space)
		return resp, nil
	}
}

// Удаление задачи
func deleteTaskDB(num int) ([]byte, error) {
	task := TaskGet{}
	db, err := sql.Open("sqlite", "base/scheduler.db")
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", num))
	err = row.Scan(&num, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		err1 := SetError{errors.New("не верно указан идентификатор")}
		resp, _ := json.Marshal(err1)
		return resp, err
	}
	_, err = db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", num))
	if err != nil {
		return []byte(err.Error()), err
	}
	return []byte{}, nil
}
