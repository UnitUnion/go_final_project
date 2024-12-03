package main

import (
	"database/sql"
	"errors"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
)

// Создали / подключили БД
func initDB() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	_, err = os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, err
	}

	if install {
		createTableQuery := `
        CREATE TABLE scheduler (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		date INTEGER,
		title TEXT,
		comment TEXT,
		repeat TEXT(128)
		);`
		_, err = db.Exec(createTableQuery)
		if err != nil {
			log.Panic(err)
		}
		log.Println("Таблица scheduler создана успешно.")
	}

	return db, nil

}

// Добавили таску в БД
func AddTask(t Task) (int64, error) {

	result, err := db.Exec("INSERT INTO scheduler ( date, title, comment, repeat ) VALUES (:date, :title, :comment, :repeat)",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat))
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Получили таску по ID
func GetTaskByID(id string) (Task, error) {
	var tasklist Task

	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id).Scan(&tasklist.ID, &tasklist.Date, &tasklist.Title, &tasklist.Comment, &tasklist.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			err := errors.New("Нет задачи")
			return tasklist, err
		}
		err := errors.New("Ошибка базы данных")
		return tasklist, err
	}

	return tasklist, nil

}

// Обновили таску по ID
func UpdateTaskByID(t Task) error {
	r, err := db.Exec("UPDATE scheduler SET date=:date, title=:title, comment=:comment, repeat=:repeat  WHERE id=:id",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.ID))
	if err != nil {
		err := errors.New("Ошибка базы данных")
		return err
	}
	affected, err := r.RowsAffected()
	if err != nil || affected == 0 {
		err := errors.New("Нет задач")
		return err
	}
	return nil
}

// Удалили таску по ID
func DeleteTaskByID(id string) error {
	r, err := db.Exec("DELETE FROM scheduler WHERE id=?", id)
	if err != nil {
		err := errors.New("Ошибка базы данных")
		log.Panic(err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		err := errors.New("Ошибка подсчета строк")
		return err
	}
	if affected == 0 {
		err := errors.New("Задача не найдена")
		return err
	}

	return nil
}

// Получили список ближайших задач
func GetTasks() ([]Task, error) {
	var tasks []Task
	result, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date")
	if err != nil {
		err := errors.New("Ошибка базы данных")
		return tasks, err
	}

	for result.Next() {
		var tasklist Task
		err := result.Scan(&tasklist.ID, &tasklist.Date, &tasklist.Title, &tasklist.Comment, &tasklist.Repeat)
		if err != nil {
			err := errors.New("Задачи не найдены")
			return tasks, err
		}
		tasks = append(tasks, tasklist)

	}

	return tasks, nil
}
