package repository

import (
	"gorm.io/gorm"
	"todo_gorm/model"
)

type Authorization interface {
	CreateUser(user *model.User) (int, error)
	GetUser(email string) (model.User, error)
	IsEmailUsed(email string) bool
}

type TodoTask interface {
	CreateTask(userId int, task *model.Task) (int, error)
	GetAll(userId int) (model.Tasks, error)
	GetTaskById(userId int, id int) (model.Task, error)
	UpdateTask(userId int, id int, task *model.Task) error
	DeleteTask(userId int, id int) error
}

type Repository struct {
	Authorization
	TodoTask
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoTask:      NewTodoTaskPostgres(db),
	}
}
