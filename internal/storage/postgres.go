package storage

import (
	"fmt"
	"github.com/3XBAT/time-tracker/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg config.Config) (*sqlx.DB, error) {
	const op = "storage.storage.NewStorage"

	db, err := sqlx.Open("postgres",
		fmt.Sprintf("port=%s user=%s host=%s dbname=%s password=%s sslmode=%s",
			cfg.DB.Port, cfg.DB.Username, cfg.DB.Host, cfg.DB.DBName, cfg.DB.Password, cfg.DB.SSLMode),
	)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping() // тут возникает ошибка failed to connect to database
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return db, nil
}
