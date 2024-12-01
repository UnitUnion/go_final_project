package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func TaskMethodPost(w http.ResponseWriter, r *http.Request) {
	var Now = time.Now()
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		sendError(w, "Json err", 500)
		return
	}
	// Ловим пустой заг
	if task.Title == "" {
		sendError(w, "Не указан заголовок задачи", 400)
		return
	}
	// Проверяем репит через NextDate
	if task.Repeat != "" {
		_, err := NextDate(Now, Now.Format(DateFormat), task.Repeat)
		if err != nil {
			sendError(w, "Ошибка правила повторения задачи", 400)
			return

		}
	}
	// Ловим не корректную дату
	if task.Date != "" {
		taskDate, err := time.Parse(DateFormat, task.Date)
		if err != nil {
			sendError(w, "Неверная дата", 400)
			return
		}
		taskDate = taskDate.Truncate(24 * time.Hour)
		if taskDate.Before(time.Now().Truncate(24 * time.Hour)) {
			if task.Repeat != "" {

				task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					sendError(w, "Неверная дата", 400)
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

	ID, err := AddTask(task)
	result := map[string]string{"id": fmt.Sprintf("%v", ID)}
	out, err := json.Marshal(result)
	if err != nil {
		sendError(w, "Ошибка json", 500)
		return
	}
	sendResponse(w, out)
}

func TaskDoneMethodPost(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	err := IDChecker(id)
	if err != nil {
		sendError(w, err.Error(), 400)
		return
	}

	var task Task
	task, err = GetTaskByID(id)
	if err != nil {
		sendError(w, err.Error(), 404)
		return
	}
	switch {
	case task.Repeat != "":
		now := time.Now()
		newDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			sendError(w, err.Error(), 400)
			return
		}
		task.Date = newDate
		err = UpdateTaskByID(task)
		if err != nil {
			sendError(w, err.Error(), 400)
			return
		}
		out, err := json.Marshal("{}")
		if err != nil {
			sendError(w, "Ошибка json", 500)
			return
		}
		sendResponse(w, out)
	default:
		err := DeleteTaskByID(id)
		if err != nil {
			sendError(w, err.Error(), 500)
			return
		}
		out, err := json.Marshal("{}")
		if err != nil {
			sendError(w, "Ошибка json", 500)
			return
		}
		sendResponse(w, out)

	}

}
