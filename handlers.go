package main

import (
	"net/http"
)

// TaskHandler Handler задачи. Методы: POST GET PUT DELETE
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		TaskMethodPost(w, r)

	case r.Method == "GET":
		TaskMethodGet(w, r)

	case r.Method == "PUT":
		TaskMethodPut(w, r)

	case r.Method == "DELETE":
		TaskMethodDelete(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

// TasksHandler Handler нескольких задач. Методы: GET
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		TasksMethodGet(w)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}

}

// TaskDoneHandler Handler выполненой задачи. Методы: POST
func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		TaskDoneMethodPost(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}

}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		NextDateMethodGet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}
