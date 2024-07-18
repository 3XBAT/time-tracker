package service

import (
	"github.com/3XBAT/time-tracker/internal/domain/models"
	"github.com/3XBAT/time-tracker/internal/storage"
	"log/slog"
)

type UserProvider interface {
	Users(params models.QueryParams) ([]models.User, error) //параметры нужны для фильтрации, если они пусты, то просто выводим все записи
	Create(passportNumber string) (int, error)
	Update(user models.UpdateUserInput, id int) error
	Delete(id int) error
	UserById(id int) (models.User, error)
}

type TaskProvider interface {
	Create(input models.InputTaskCreate) (int, error)
	Update(task models.InputTaskUpdate) error
	Delete(taskDeleteRequest models.InputTaskDelete) error
	Tasks(task models.InputTask) ([]models.OutputTask, error)
}

type Service struct {
	UserProvider
	TaskProvider
}

func NewService(log *slog.Logger, s *storage.Storage) *Service {
	return &Service{
		UserProvider: NewUserService(s.UserProvider, log),
		TaskProvider: NewTaskService(s.TaskProvider, log),
	}
}
