package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func TaskMethodGet(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	err := IDChecker(id)
	if err != nil {
		sendError(w, err.Error(), 500)
		return
	}
	task, err := GetTaskByID(id)
	if err != nil {
		sendError(w, err.Error(), 404)
		return
	}

	out, err := json.Marshal(task)
	if err != nil {
		sendError(w, "Ошибка json", 500)
		return
	}
	sendResponse(w, out)
}

func TasksMethodGet(w http.ResponseWriter) {
	tasks, err := GetTasks()
	if err != nil {
		sendError(w, "Нет задач", 404)
		return
	}
	result := map[string][]Task{"tasks": tasks}

	out, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err.Error())
	}

	sendResponse(w, out)
}

func NextDateMethodGet(w http.ResponseWriter, r *http.Request) {
	var result string
	n := r.FormValue("now")
	d := r.FormValue("date")
	repeat := r.FormValue("repeat")
	if n == "" {
		n = time.Now().Format(DateFormat)
	}
	_, err := time.Parse(DateFormat, d)
	if err != nil {
		sendError(w, err.Error(), 400)
		return
	}
	now, err := time.Parse(DateFormat, n)
	if err != nil {
		sendError(w, err.Error(), 500)
		return
	}
	result, err = NextDate(now, d, repeat)
	if err != nil {
		sendError(w, err.Error(), 400)
		return
	}

	out, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err.Error(), 500)
	}

	sendResponse(w, out)

}
