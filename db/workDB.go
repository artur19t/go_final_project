package db

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"go-final/model"
)

func PostTaskDB(task model.Task) (string, bool, model.Task) {
	today := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = today
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err.Error(), true, task
	}

	back := ""
	task.Date, back = checkRepeat(task.Date, task.Repeat, today)
	if back != "" {
		return back, true, task
	}

	if task.Title == "" {
		return "не указан заголовок задачи", true, task
	}
	return "", false, task
}

// Запись задачи в БД
func CreateTask(task model.Task) (int64, error) {
	res, err := DB.Db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err != nil {
		return 1, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}

// Поиск задачи по ID
func FindTaskID(num int, idSearch string) (string, bool, model.TaskGet) {
	task := model.TaskGet{}

	row := DB.Db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", num))
	err := row.Scan(&num)
	if err != nil {
		return "не верно указан идентификатор", true, task
	}

	rows, err := DB.Db.Query(fmt.Sprintf("SELECT * FROM scheduler WHERE id = %s", idSearch))
	if err != nil {

		return "задача не найдена", true, task
	}
	defer rows.Close()
	for rows.Next() {
		id := 0
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			log.Println(err)
			return "error", true, task
		}
		task.ID = fmt.Sprint(id)
	}

	return "", false, task
}

// Изменение параметров задачи
func CorrectTaskDB(task model.TaskGet) (string, bool) {
	num, err := strconv.Atoi(task.ID)
	if err != nil {
		return "не верно указан идентификатор", true
	}
	row := DB.Db.QueryRow("SELECT id FROM scheduler WHERE id = :id", sql.Named("id", num))
	err = row.Scan(&num)
	if err != nil {
		return "не верно указан идентификатор", true
	}
	DB.Db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", num))
	return "", false
}

// Запрос задач
func GetTaskDB() (model.SliceTask, error) {
	tasks := model.SliceTask{}

	sliceTask := []model.TaskGet{}
	rows, err := DB.Db.Query("SELECT * FROM scheduler")
	if err != nil {
		return tasks, err
	}
	defer rows.Close()

	for rows.Next() {
		task := model.TaskGet{}
		id := 0
		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return tasks, err
		}
		task.ID = fmt.Sprint(id)
		sliceTask = append(sliceTask, task)
	}
	tasks.Tasks = sliceTask
	return tasks, nil
}

// Выполнение задачи
func DoneTaskDB(num int) (string, error) {
	task := model.TaskGet{}
	row := DB.Db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", num))
	err := row.Scan(&num, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return "не верно указан идентификатор", err
	}

	if task.Repeat == "" {
		_, err := DB.Db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", num))
		if err != nil {
			return err.Error(), err
		}
		return "", nil
	} else {
		date, _ := NextDate(time.Now(), task.Date, task.Repeat)
		if date == "-1" {
			return "ошибка переноса даты", err
		}
		DB.Db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
			sql.Named("date", date),
			sql.Named("title", task.Title),
			sql.Named("comment", task.Comment),
			sql.Named("repeat", task.Repeat),
			sql.Named("id", num))
		return "", nil
	}
}

// Удаление задачи
func DeleteTaskDB(num int) (string, error) {
	task := model.TaskGet{}

	row := DB.Db.QueryRow("SELECT * FROM scheduler WHERE id = :id", sql.Named("id", num))
	err := row.Scan(&num, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return "не верно указан идентификатор", err
	}
	_, err = DB.Db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", num))
	if err != nil {
		return err.Error(), err
	}
	return "", nil
}
