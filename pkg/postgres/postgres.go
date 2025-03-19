package postgres

import (
	"fmt"
	"order-service/internal/config"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func New(cfg config.Postgres) (*sqlx.DB, error) {
	const op = "postgres.New"

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return conn, nil
}
