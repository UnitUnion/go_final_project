package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	taskDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", err
	}
	switch {
	case repeat == "":
		err := errors.New("Ошибка правила повторения задачи")
		return "", err

	case repeat == "y":
		for {
			taskDate = taskDate.AddDate(1, 0, 0)
			if taskDate.After(now) {
				return taskDate.Format(DateFormat), nil
			}
		}

	case strings.HasPrefix(repeat, "d "):
		for {
			daysStr := strings.TrimPrefix(repeat, "d ")
			days, err := strconv.Atoi(daysStr)
			if err != nil || days < 1 || days > 400 {
				err := errors.New("Ошибка правила повторения задачи")
				return "", err
			}
			taskDate = taskDate.AddDate(0, 0, days)
			if taskDate.After(now) {
				return taskDate.Format(DateFormat), nil
			}

		}

	default:
		err := errors.New("Ошибка правила повторения задачи")
		return "", err
	}

}

func IDChecker(id string) error {
	if id == "" {
		err := errors.New("ID пустой")
		return err
	}
	_, err := strconv.Atoi(id)
	if err != nil {
		err := errors.New("ID не целое число")
		return err
	}
	return nil
}
