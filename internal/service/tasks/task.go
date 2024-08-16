package tasks

import (
	"fmt"
	"regexp"
	"time"

	"cactus3d/go_final_project/internal/models"
	"cactus3d/go_final_project/internal/nextdate"
)

type Service struct {
	store StoreProvider
}

type StoreProvider interface {
	CreateTask(task *models.Task) (int, error)
	GetTaskById(id string) (*models.Task, error)
	GetTasks() ([]models.Task, error)
	GetTasksBySearch(search string) ([]models.Task, error)
	GetTasksByDate(date string) ([]models.Task, error)
	UpdateTask(*models.Task) (int64, error)
	DeleteTaskById(id string) (int64, error)
}

func New(store StoreProvider) *Service {
	return &Service{store: store}
}

func (s *Service) Add(date, title, comment, repeat string) (int, error) {
	_, err := time.Parse(nextdate.DateFormat, date)
	if err != nil {
		return 0, err
	}
	n := time.Now().Format(nextdate.DateFormat)

	if date < n {
		if repeat != "" {
			date, err = nextdate.NextDate(time.Now(), date, repeat)
			if err != nil {
				return 0, err
			}
		} else {
			date = time.Now().Format(nextdate.DateFormat)
		}
	}

	task := &models.Task{
		Title:   title,
		Date:    date,
		Comment: comment,
		Repeat:  repeat,
	}

	return s.store.CreateTask(task)
}

func (s *Service) GetAll(search string) ([]models.Task, error) {

	if search != "" {
		_, err := regexp.Compile(`[0-3][0-9]\.[0-1][0-9]\.20[0-9][0-9]`)
		if err != nil {
			return nil, err
		}
		matched, _ := regexp.MatchString(`[0-3][0-9]\.[0-1][0-9]\.20[0-9][0-9]`, search)
		if matched {
			date, err := time.Parse("02.01.2006", search)
			if err != nil {
				return nil, err
			}
			search = date.Format(nextdate.DateFormat)
			return s.store.GetTasksByDate(search)
		}

		search = "%" + search + "%"
		return s.store.GetTasksBySearch(search)
	}

	return s.store.GetTasks()
}

func (s *Service) Get(id string) (*models.Task, error) {
	task, err := s.store.GetTaskById(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, fmt.Errorf("not found")
	}

	return task, nil
}

func (s *Service) Update(task *models.Task) error {
	count, err := s.store.UpdateTask(task)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("no rows with such id")
	}
	return nil
}

func (s *Service) Delete(id string) error {
	count, err := s.store.DeleteTaskById(id)
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("nothing to delete")
	}

	return nil
}

func (s *Service) Done(id string) error {
	task, err := s.Get(id)
	if err != nil {
		return err
	}

	if task.Repeat == "" {
		err = s.Delete(task.Id)
		if err != nil {
			return err
		}
		return nil
	}

	date, err := nextdate.NextDate(time.Now(), task.Date, task.Repeat)
	if err != nil {
		return err
	}

	task.Date = date
	err = s.Update(task)
	if err != nil {
		return err
	}

	return nil
}
