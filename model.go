package main

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
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

type SetTask struct {
	ID int `json:"id"`
}

type SetError struct {
	Error error `json:"error"`
}

type FindTask struct {
	ID string `json:"id"`
}

type StringValue struct {
	Error string `json:"error"`
}