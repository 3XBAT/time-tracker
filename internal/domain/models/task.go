package models

import "time"

type Task struct {
	Id        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Name      string    `json:"name" db:"name"`
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
}

type InputTaskUpdate struct {
	Id     int `json:"id"`
	UserID int `json:"user_id"`
}

type InputTask struct {
	UserID      int        `json:"user_id" db:"user_id"`
	StartPeriod *time.Time `json:"start_time" db:"start_time"`
	EndPeriod   *time.Time `json:"end_time" db:"end_time"`
}

type InputTaskCreate struct {
	UserID      int        `json:"user_id"`
	Name        string     `json:"name"`
	StartPeriod *time.Time `json:"start_time"`
}

type OutputTask struct {
	Name     string `json:"name" `
	Duration string `json:"duration"`
}

type InputTaskDelete struct {
	UserID int `json:"user_id"`
	TaskID int `json:"task_id"`
}
