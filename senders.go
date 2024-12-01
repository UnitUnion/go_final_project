package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func sendResponse(w http.ResponseWriter, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	_, err := w.Write(data)
	if err != nil {
		log.Panic(err)
	}
}

func sendError(w http.ResponseWriter, errorMessage string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]string{"error": errorMessage})
	if err != nil {
		log.Panic(err)
	}
}
