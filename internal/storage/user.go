package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/3XBAT/time-tracker/internal/domain/models"
	"github.com/jmoiron/sqlx"
	"strings"
)

const (
	ID             = "id"
	Name           = "name"
	Surname        = "surname"
	PassportNumber = "passport_number"
	Address        = "address"
)

type UserStorage struct {
	db *sqlx.DB
}

func NewUserStorage(db *sqlx.DB) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) Users(params models.QueryParams) ([]models.User, error) {
	const op = "storage.Users"
	var users []models.User

	query, args := buildQuery(params)

	err := s.db.Select(&users, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return users, nil
}

func (s *UserStorage) UserByID(id int) (models.User, error) {
	const op = "storage.UserByID"
	var user models.User

	query := fmt.Sprintf("SELECT * FROM users WHERE id = $1")

	err := s.db.Get(&user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *UserStorage) Create(user models.User) (int, error) {
	const op = "storage.CreateUser"

	var id int

	query := fmt.Sprintf(`INSERT INTO users (name, surname, patronymic, passport_number, address)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`)

	row := s.db.QueryRow(query, user.Name, user.Surname, user.Patronymic, user.PassportNumber, user.Address)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *UserStorage) Update(input models.UpdateUserInput, id int) error {
	const op = "storage.UpdateUser"

	setValues := make([]string, 0)
	args := make([]interface{}, 0)

	if input.Address != nil {
		setValues = append(setValues, fmt.Sprintf("addr=$%d", len(setValues)+1))
		args = append(args, *input.Address)
	}

	if input.PassportNumber != nil {
		setValues = append(setValues, fmt.Sprintf("passport_number=$%d", len(setValues)+1))
		args = append(args, *input.PassportNumber)
	}

	setQuery := strings.Join(setValues, ", ")
	fmt.Println(setQuery)
	query := fmt.Sprintf(`UPDATE users SET %s WHERE id = $%d`, setQuery, len(setValues)+1)
	args = append(args, id)

	fmt.Println(query, args)

	_, err := s.db.Exec(query, args...)

	return err
}

func (s *UserStorage) Delete(userID int) error {
	const op = "storage.Delete"

	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	deleteTasksQuery := fmt.Sprintf(`DELETE FROM tasks WHERE user_id = $1`)
	_, err = tx.Exec(deleteTasksQuery, userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	deleteUserQuery := fmt.Sprintf(`DELETE FROM users WHERE id = $1`)
	_, err = tx.Exec(deleteUserQuery, userID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit()
}

func buildQuery(params models.QueryParams) (string, []interface{}) {
	baseQuery := "SELECT * FROM users WHERE 1=1"
	var args []interface{}
	var conditions []string

	if params.ID != "" {
		conditions, args = appendCondition(conditions, args, ID, params.ID)
	}
	if params.Name != "" {
		conditions, args = appendCondition(conditions, args, Name, params.Name)
	}
	if params.PassportNumber != "" {
		conditions, args = appendCondition(conditions, args, PassportNumber, params.PassportNumber)
	}
	if params.Address != "" {
		conditions, args = appendCondition(conditions, args, Address, params.Address)
	}
	if params.Surname != "" {
		conditions, args = appendCondition(conditions, args, Surname, params.Surname)
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)+1, len(args)+2)
	args = append(args, params.Limit, params.Offset)

	return baseQuery, args
}

func appendCondition(conditions []string, args []interface{}, field, value string) ([]string, []interface{}) {
	if strings.HasPrefix(value, ">=") || strings.HasPrefix(value, "<=") || strings.HasPrefix(value, ">") || strings.HasPrefix(value, "<") {
		operator := value[:2]
		if operator[1] != '=' {
			operator = value[:1]
			conditions = append(conditions, fmt.Sprintf("%s %s $%d", field, operator, len(args)+1))
			args = append(args, value[1:])
		} else {
			conditions = append(conditions, fmt.Sprintf("%s %s $%d", field, operator, len(args)+1))
			args = append(args, value[2:])
		}
	} else {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", field, len(args)+1))
		args = append(args, value)
	}
	return conditions, args
}
