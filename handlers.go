package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// TaskHandler Handler задачи. Методы: POST GET PUT DELETE
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		var task Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			sendError(w, "Json err")
			return
		}
		// Ловим пустой заг
		if task.Title == "" {
			sendError(w, "Не указан заголовок задачи")
			return
		}
		// Проверяем репит через NextDate
		if task.Repeat != "" {
			_, err := NextDate(Now, Now.Format(DateFormat), task.Repeat)
			if err != nil {
				sendError(w, "Ошибка правила повторения задачи")
				return

			}
		}
		// Ловим не корректную дату
		if task.Date != "" {
			taskDate, err := time.Parse(DateFormat, task.Date)
			if err != nil {
				sendError(w, "Неверная дата")
				return
			}
			taskDate = taskDate.Truncate(24 * time.Hour)
			if taskDate.Before(time.Now().Truncate(24 * time.Hour)) {
				if task.Repeat != "" {

					task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
					if err != nil {
						sendError(w, "Неверная дата")
						return
					}
				} else {
					task.Date = time.Now().Format(DateFormat)
				}
			}
		}
		// Ловим пустую дату
		if task.Date == "" {
			task.Date = time.Now().Format(DateFormat)
		}

		ID := AddTask(task)
		response := map[string]string{"id": fmt.Sprintf("%v", ID)}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			sendError(w, "Ошибка json")
		}

	case r.Method == "GET":
		id := r.FormValue("id")
		err := IDChecker(id)
		if err != nil {
			sendError(w, err.Error())
			return
		}
		task, err := GetTaskByID(id)
		if err != nil {
			sendError(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err = json.NewEncoder(w).Encode(task)
		if err != nil {
			sendError(w, "Ошибка json")
		}

	case r.Method == "PUT":
		var task Task

		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			sendError(w, "json error")
			return
		}

		err = IDChecker(task.ID)
		if err != nil {
			sendError(w, err.Error())
			return
		}

		_, err = GetTaskByID(task.ID)
		if err != nil {
			sendError(w, err.Error())
			return
		}

		// Ловим пустой заг
		if task.Title == "" {
			sendError(w, "Не указан заголовок задачи")
			return
		}
		// Проверяем репит через NextDate
		if task.Repeat != "" {
			_, err := NextDate(Now, Now.Format(DateFormat), task.Repeat)
			if err != nil {
				sendError(w, err.Error())
				return

			}
		}
		// Ловим не корректную дату
		if task.Date != "" {
			taskDate, err := time.Parse(DateFormat, task.Date)
			if err != nil {
				sendError(w, "Неверная дата")
				return
			}
			if taskDate.Before(time.Now()) {
				if task.Repeat != "" {
					task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
					if err != nil {
						sendError(w, err.Error())
						return
					}
				} else {
					task.Date = time.Now().Format(DateFormat)
				}
			}
		}
		// Ловим пустую дату
		if task.Date == "" {
			task.Date = time.Now().Format(DateFormat)
		}
		err = UpdateTaskByID(task)
		if err != nil {
			sendError(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{}"))

	case r.Method == "DELETE":
		id := r.FormValue("id")
		err := IDChecker(id)
		if err != nil {
			sendError(w, err.Error())
			return
		}

		_, err = GetTaskByID(id)
		if err != nil {
			sendError(w, err.Error())
			return
		}
		err = DeleteTaskByID(id)
		if err != nil {
			sendError(w, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte("{}"))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

// TasksHandler Handler нескольких задач. Методы: GET
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		resp, err := GetTasks()
		if err != nil {
			sendError(w, "Нет задач")
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			sendError(w, "Ошибка json")
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}

}

// TaskDoneHandler Handler выполненой задачи. Методы: POST
func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "POST":
		id := r.FormValue("id")

		err := IDChecker(id)
		if err != nil {
			sendError(w, err.Error())
			return
		}

		var task Task
		task, err = GetTaskByID(id)
		if err != nil {
			sendError(w, err.Error())
			return
		}
		switch {
		case task.Repeat != "":
			now := time.Now()
			newDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				sendError(w, err.Error())
				return
			}
			task.Date = newDate
			err = UpdateTaskByID(task)
			if err != nil {
				sendError(w, err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{}"))

		default:
			err := DeleteTaskByID(id)
			if err != nil {
				sendError(w, err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{}"))

		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}

}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == "GET":
		var res string
		n := r.FormValue("now")
		d := r.FormValue("date")
		repeat := r.FormValue("repeat")
		if n == "" {
			n = time.Now().Format(DateFormat)
		}
		_, err := time.Parse(DateFormat, d)
		if err != nil {
			res = ""
		}
		now, err := time.Parse(DateFormat, n)
		if err != nil {
			res = ""
		}
		res, err = NextDate(now, d, repeat)
		if err != nil {
			res = ""
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write([]byte(res))

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func sendError(w http.ResponseWriter, errorMessage string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})

}
