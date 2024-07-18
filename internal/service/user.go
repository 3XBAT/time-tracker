package service

import (
	"errors"
	"fmt"
	api2 "github.com/3XBAT/time-tracker/internal/api"
	"github.com/3XBAT/time-tracker/internal/config"
	"github.com/3XBAT/time-tracker/internal/domain/models"
	"github.com/3XBAT/time-tracker/internal/storage"
	"log/slog"
	"strings"
)

type UserService struct {
	storage storage.UserProvider
	log     *slog.Logger
}

func NewUserService(s storage.UserProvider, log *slog.Logger) *UserService {
	return &UserService{
		storage: s,
		log:     log,
	}
}

func (us *UserService) Users(params models.QueryParams) ([]models.User, error) {
	const op = "service.user.Users"
	log := us.log.With(slog.String("op", op))

	log.Debug("Received request with params", slog.Any("params", params))
	log.Info("attempting to get users")

	users, err := us.storage.Users(params)
	if err != nil {
		log.Warn(err.Error())
		return nil, err
	}

	log.Debug("Successfully retrieved users", slog.Any("users", users))
	log.Info("getting users was successful")
	return users, nil
}

func (us *UserService) UserById(id int) (models.User, error) {
	const op = "service.user.UserById"
	log := us.log.With(slog.String("op", op))

	log.Debug("Received request for user ID", slog.Int("id", id))
	log.Info("attempting to get user")

	user, err := us.storage.UserByID(id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn(err.Error())
			return models.User{}, err
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("Successfully retrieved user", slog.Any("user", user))
	log.Info("getting user was successful")

	return user, nil
}

func (us *UserService) Create(passportNumber string) (id int, err error) {
	const op = "service.user.CreateUser"
	log := us.log.With(slog.String("op", op))

	log.Debug("Received request to create user with passport number", slog.String("passportNumber", passportNumber))
	log.Info("attempting to get user info")

	psprtParams := strings.Split(passportNumber, " ")

	cfg := config.MustLoad()
	api := api2.NewApiClient(&cfg)

	user, err := api.UserInfo(psprtParams[0], psprtParams[1])
	if err != nil {
		log.Error(err.Error())
		return 0, err
	}
	if _, err := us.storage.UserByID(user.ID); err == nil {
		log.Error(err.Error())
		return 0, storage.ErrUserExists
	}

	id, err = us.storage.Create(*user)
	if err != nil {
		log.Error(err.Error())
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Debug("Successfully created user", slog.Int("id", id))
	log.Info("creating user was successful")

	return id, nil
}

func (us *UserService) Update(user models.UpdateUserInput, id int) error {
	const op = "service.user.UpdateUser"
	log := us.log.With(slog.String("op", op))

	log.Debug("Received request to update user", slog.Any("user", user), slog.Int("id", id))

	if _, err := us.storage.UserByID(id); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn(err.Error())
			return err
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("attempting to update user")

	err := us.storage.Update(user, id)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Debug("Successfully updated user", slog.Int("id", id))
	log.Info("updating user was successful")
	return nil
}

func (us *UserService) Delete(id int) error {
	const op = "service.user.Delete"
	log := us.log.With(slog.String("op", op))

	log.Debug("Received request to delete user", slog.Int("id", id))
	log.Info("attempting to delete user")

	if _, err := us.storage.UserByID(id); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info(err.Error())
			return err
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	err := us.storage.Delete(id)
	if err != nil {
		log.Error("error deleting user: " + err.Error())
		return err
	}

	log.Debug("Successfully deleted user", slog.Int("id", id))
	log.Info("deleting user was successful")
	return nil
}
