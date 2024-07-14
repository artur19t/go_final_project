package db

import (
	"encoding/json"
	"go-final/model"
	"strconv"
	"time"
)

func checkRepeat(date, repeat, today string) (string, string) {
	if date < today && repeat == "" {
		date = today
	}
	if date < today && repeat != "" {
		date, _ = NextDate(time.Now(), date, repeat)
		if date == "-1" {
			return date, "не верно указано правило повторения"
		}
	}
	return date, ""
}

// Проверка корректности входящих данных
func CheckCorrectIn(task model.TaskGet) (string, bool) {
	today := time.Now().Format("20060102")
	if task.Date == "" {
		task.Date = today
	}

	_, err := time.Parse("20060102", task.Date)
	if err != nil {
		return err.Error(), true
	}
	if task.ID == "" {
		return "не указан идентификатор", true
	}
	back := ""
	task.Date, back = checkRepeat(task.Date, task.Repeat, today)
	if back != "" {
		return "не верно указано правило повторения", true
	}
	if task.Title == "" {
		return "не указан заголовок задачи", true
	}
	return "", false
}

// Проверка корректности введеного ID
func CheckID(idSearch string) (bool, []byte, int) {
	if idSearch == "" {
		err := model.StringValue{}
		err.Error = "не указан идентификатор"
		resp, _ := json.Marshal(err)
		return false, resp, 0
	}

	num, err := strconv.Atoi(idSearch)
	if err != nil {
		err := model.StringValue{}
		err.Error = "не верно указан идентификатор"
		resp, _ := json.Marshal(err)

		return false, resp, 0
	}
	return true, nil, num
}
