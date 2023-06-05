package service

import (
	"errors"
	"log"
	"todo_gorm/internal/repository"
	"todo_gorm/model"
)

type TodoTaskService struct {
	repo repository.TodoTask
}

func NewTodoTaskService(repo repository.TodoTask) *TodoTaskService {
	return &TodoTaskService{repo: repo}
}

func (s *TodoTaskService) CreateTask(userId int, task *model.Task) (int, error) {
	id, err := s.repo.CreateTask(userId, task)
	if err != nil {
		log.Println("failed to create task")
		return 0, err
	}
	return id, nil
}

func (s *TodoTaskService) GetAll(userId int) (model.Tasks, error) {
	tasks, err := s.repo.GetAll(userId)
	if err != nil {
		log.Println("failed to get list of tasks due to:", err.Error())
		return nil, err
	}

	return tasks, nil
}

func (s *TodoTaskService) GetTaskById(userId int, id int) (model.Task, error) {
	task, err := s.repo.GetTaskById(userId, id)
	if err != nil {
		log.Println("failed to get task due to:", err.Error())
		return model.Task{}, err
	}

	return task, nil
}

func (s *TodoTaskService) ValidateTask(task *model.Task) error {
	if task.Name == "" {
		return errors.New("update task has no values")
	}
	return nil
}

func (s *TodoTaskService) UpdateTask(userId int, id int, task *model.Task) error {
	err := s.repo.UpdateTask(userId, id, task)
	if err != nil {
		log.Println("failed to update task due to:", err.Error())
		return err
	}

	return nil
}

func (s *TodoTaskService) DeleteTask(userId int, id int) error {
	err := s.repo.DeleteTask(userId, id)
	if err != nil {
		log.Println("failed to delete task due to:", err.Error())
		return err
	}
	return nil
}
