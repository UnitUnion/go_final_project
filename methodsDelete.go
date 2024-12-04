package main

import (
	"net/http"
)

func TaskMethodDelete(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	err := IDChecker(id)
	if err != nil {
		sendError(w, err.Error(), 400)
		return
	}

	_, err = GetTaskByID(id)
	if err != nil {
		sendError(w, err.Error(), 404)
		return
	}
	err = DeleteTaskByID(id)
	if err != nil {
		sendError(w, err.Error(), 500)
		return
	}

	out := []byte("{}")
	sendResponse(w, out)
}
