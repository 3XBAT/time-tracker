package storage

import (
	"errors"
	"github.com/3XBAT/time-tracker/internal/domain/models"
	"github.com/jmoiron/sqlx"
)

var (
	ErrUserNotFound = errors.New("users not found")
	ErrTaskNotFound = errors.New("tasks not found")
	ErrTaskEnded    = errors.New("task already finished")
	ErrUserExists   = errors.New("user already exists")
	ErrBadRequest   = errors.New("bad request")
)

type UserProvider interface {
	Users(params models.QueryParams) ([]models.User, error) //параметры нужны для фильтрации, если они пусты, то просто выводим все записи
	UserByID(id int) (models.User, error)
	Create(user models.User) (int, error)
	Update(user models.UpdateUserInput, id int) error
	Delete(id int) error
}

type TaskProvider interface {
	Create(task models.InputTaskCreate) (int, error)
	Update(task models.InputTaskUpdate) error
	Delete(task models.InputTaskDelete) error
	Tasks(task models.InputTask) ([]models.OutputTask, error)
}

type Storage struct {
	UserProvider
	TaskProvider
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		UserProvider: NewUserStorage(db),
		TaskProvider: NewTaskStorage(db),
	}
}
