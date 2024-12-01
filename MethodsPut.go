package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func TaskMethodPut(w http.ResponseWriter, r *http.Request) {
	var Now = time.Now()
	var task Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		sendError(w, "json error", 500)
		return
	}

	err = IDChecker(task.ID)
	if err != nil {
		sendError(w, err.Error(), 400)
		return
	}

	_, err = GetTaskByID(task.ID)
	if err != nil {
		sendError(w, err.Error(), 404)
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
			sendError(w, err.Error(), 400)
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
		if taskDate.Before(time.Now()) {
			if task.Repeat != "" {
				task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					sendError(w, err.Error(), 400)
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
		sendError(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal("")
	if err != nil {
		sendError(w, "Ошибка json", 500)
		return
	}
	sendResponse(w, out)
}
