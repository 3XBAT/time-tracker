package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/3XBAT/time-tracker/internal/domain/models"
	"github.com/jmoiron/sqlx"
)

type TaskStorage struct {
	db *sqlx.DB
}

func NewTaskStorage(db *sqlx.DB) *TaskStorage {
	return &TaskStorage{db: db}
}

func (s *TaskStorage) Create(input models.InputTaskCreate) (int, error) {
	const op = "storage.task.Create"

	if input.StartPeriod == nil {
		now := time.Now()
		input.StartPeriod = &now
	}
	query := fmt.Sprintf(`INSERT INTO tasks (user_id, name, start_time)
VALUES ($1, $2, $3) RETURNING id`)

	var id int

	err := s.db.QueryRow(query, input.UserID, input.Name, input.StartPeriod).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *TaskStorage) Update(input models.InputTaskUpdate) error {
	const op = "storage.task.Update"

	res, err := s.IsUpdate(input.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if res {
		return ErrTaskEnded
	}

	query := fmt.Sprintf(`UPDATE tasks SET end_time = $1 WHERE id = $2 AND user_id = $3 AND end_time IS NULL`)

	endTime := time.Now()

	_, err = s.db.Exec(query, endTime, input.Id, input.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err = s.IsUpdate(input.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if !res {
		return ErrBadRequest
	}

	return nil
}

func (s *TaskStorage) IsUpdate(taskID int) (bool, error) {
	const op = "storage.task.IsUpdate"

	query := fmt.Sprintf(`SELECT end_time FROM tasks WHERE id = $1`)
	var endTime sql.NullTime

	err := s.db.QueryRow(query, taskID).Scan(&endTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return endTime.Valid, nil
}

func (s *TaskStorage) Delete(input models.InputTaskDelete) error {
	const op = "storage.task.Delete"

	_, err := s.TaskById(input.TaskID)
	if err != nil {
		return err
	}

	query := fmt.Sprintf(`DELETE FROM tasks WHERE id = $1 AND user_id = $2`)
	_, err = s.db.Exec(query, input.TaskID, input.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *TaskStorage) Tasks(input models.InputTask) ([]models.OutputTask, error) {
	const op = "storage.task.Tasks"
	var tasks []models.OutputTask

	query, args := buildQueryForTasks(input)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.OutputTask
		var durationSeconds float64
		if err := rows.Scan(&task.Name, &durationSeconds); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		task.Duration = formatDuration(time.Duration(durationSeconds) * time.Second)
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during rows iteration: %w", err)
	}

	return tasks, nil
}

func (s *TaskStorage) TaskById(taskID int) (models.Task, error) {
	const op = "storage.task.TaskById"
	var task models.Task

	query := fmt.Sprintf(`SELECT * FROM tasks WHERE id = $1`)

	err := s.db.Get(&task, query, taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, ErrTaskNotFound
		}
		return models.Task{}, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil

}

func formatDuration(d time.Duration) string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	
	if days > 0 {
		return fmt.Sprintf("%dd %02dh %02dm", days, hours, minutes)
	}
	return fmt.Sprintf("%02dh %02dm", hours, minutes)
}

func buildQueryForTasks(input models.InputTask) (string, []interface{}) {
	baseQuery := `SELECT 
            name, 
            EXTRACT(EPOCH FROM (end_time - start_time)) AS duration
        FROM 
            tasks
        WHERE user_id = $1 AND 
            start_time IS NOT NULL AND 
            end_time IS NOT NULL `
	var conditions []string
	var args []interface{}
	args = append(args, input.UserID)

	if input.StartPeriod != nil {
		conditions = append(conditions, fmt.Sprintf("start_time >= $%d", len(args)+1))
		args = append(args, *input.StartPeriod)
	}
	if input.EndPeriod != nil {
		conditions = append(conditions, fmt.Sprintf("end_time <= $%d", len(args)+1))
		args = append(args, *input.EndPeriod)
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	baseQuery += ` ORDER BY 
            duration DESC;`

	return baseQuery, args
}
