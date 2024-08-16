package models

import (
	"fmt"
	"strconv"
	"time"

	"cactus3d/go_final_project/internal/nextdate"
)

type Task struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Date    string `json:"date"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func (task *Task) Validate() (bool, error) {
	if task.Id == "" {
		return false, fmt.Errorf("не указан id")
	}
	if id, err := strconv.Atoi(task.Id); err != nil || id < 0 {
		return false, fmt.Errorf("id должен быть положительным числом")
	}

	if task.Title == "" {
		return false, fmt.Errorf("не указан заголовок задачи")
	}

	if task.Date == "" {
		task.Date = time.Now().Format(nextdate.DateFormat)
	} else {
		_, err := time.Parse(nextdate.DateFormat, task.Date)
		if err != nil {
			return false, err
		}
	}

	if task.Repeat != "" {
		_, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}
