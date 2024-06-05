package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

func main() {
	CheckBase()
	dir := "web"
	http.Handle("/api/nextdate", http.HandlerFunc(nextDate))
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.Handle("/api/task", http.HandlerFunc(determinant))
	http.Handle("/api/tasks", http.HandlerFunc(getTasks))
	http.Handle("/api/task/done", http.HandlerFunc(doneTask))
	log.Printf("Запуск сервера на %d порту", 7540)
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}

func determinant(res http.ResponseWriter, req *http.Request) {
	method := fmt.Sprint(req.Method)
	if method == "POST" {
		postTask(res, req)
	}
	if method == "GET" {
		getTask(res, req)
	}
	if method == "PUT" {
		correctTask(res, req)
	}
	if method == "DELETE" {
		deleteTask(res, req)
	}
}

func nextDate(res http.ResponseWriter, req *http.Request) {
	now := req.URL.Query().Get("now")
	date := req.URL.Query().Get("date")
	repeat := req.URL.Query().Get("repeat")
	nowT, err := time.Parse("20060102", now)
	if err != nil {
		res.Write([]byte(fmt.Sprintf("%T", err)))
	}
	dateToRet, _ := NextDate(nowT, date, repeat)
	res.Write([]byte(dateToRet))
}
