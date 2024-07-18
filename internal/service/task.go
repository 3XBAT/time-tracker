package service

import (
	"errors"
	"github.com/3XBAT/time-tracker/internal/domain/models"
	"github.com/3XBAT/time-tracker/internal/storage"
	"log/slog"
)

type TaskService struct {
	storage storage.TaskProvider
	log     *slog.Logger
}

func NewTaskService(s storage.TaskProvider, log *slog.Logger) *TaskService {
	return &TaskService{
		storage: s,
		log:     log,
	}
}

func (ts *TaskService) Create(input models.InputTaskCreate) (int, error) {
	const op = "service.task.Create"

	log := ts.log.With(slog.String("op", op))

	log.Debug("received request to create task", slog.Any("input", input))
	log.Info("starting create task")

	id, err := ts.storage.Create(input)
	if err != nil {
		log.Warn("failed creating task:", err)
		return id, err
	}

	log.Debug("successfully created task", slog.Int("id", id))
	log.Info("task created")

	return id, nil
}

func (ts *TaskService) Update(task models.InputTaskUpdate) error {
	const op = "service.task.Update"

	log := ts.log.With(slog.String("op", op))

	log.Debug("received request to update task", slog.Any("task", task))
	log.Info("trying to update task")

	err := ts.storage.Update(task)
	if err != nil {
		if errors.Is(err, storage.ErrTaskEnded) {
			log.Warn(err.Error())
			return err
		}
		if errors.Is(err, storage.ErrTaskNotFound) {
			log.Warn(err.Error())
			return err
		}
		log.Warn("failed updating task:", err)
		return err
	}

	log.Debug("successfully updated task", slog.Any("task", task))
	log.Info("task updated successfully")

	return err
}

func (ts *TaskService) Delete(task models.InputTaskDelete) error {
	const op = "service.task.Delete"

	log := ts.log.With(slog.String("op", op))

	log.Debug("Received request to delete task", slog.Any("task", task))
	log.Info("trying to delete task")

	err := ts.storage.Delete(task)
	if err != nil {
		log.Warn("failed deleting task:", err.Error())
	}

	log.Debug("successfully deleted task", slog.Any("task", task))
	log.Info("task deleted")

	return err
}

func (ts *TaskService) Tasks(input models.InputTask) ([]models.OutputTask, error) {
	const op = "service.task.Tasks"

	log := ts.log.With(slog.String("op", op))

	log.Debug("received request to get tasks", slog.Any("input", input))
	log.Info("trying to get tasks")

	tasks, err := ts.storage.Tasks(input)
	if err != nil {
		log.Warn("failed getting tasks:", err)
		return nil, err
	}

	log.Debug("successfully retrieved tasks", slog.Int("count", len(tasks)))
	log.Info("getting tasks successfully")
	return tasks, nil
}
