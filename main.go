package main

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"net/http"
	"os"
	"time"
)

const DateFormat = "20060102"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

var db *sql.DB
var err error
var Now = time.Now()

func main() {
	db, err = initDB()
	if err != nil {
		log.Panic(err)
	}
	ServicePort, ok := os.LookupEnv("ServicePort")
	if !ok {
		ServicePort = "7540"
	}
	addr := ":" + ServicePort
	fmt.Println("Service up on port:", addr)

	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/task", TaskHandler)
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", TaskDoneHandler)
	http.HandleFunc("/api/nextdate", NextDateHandler)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}
