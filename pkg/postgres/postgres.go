package postgres

import (
	"errors"
	"fmt"
	"movie-service/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func New(cfg config.Postgres) (*sqlx.DB, error) {
	const op = "postgres.New"

	// Form connection string
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	// Connect to postgres
	conn, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to postgres: %w", op, err)
	}

	// Ping db in order to check connection because successful connection creation
	// guarantees nothing
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("%s: failed to ping postgres: %w", op, err)
	}

	// Migrations
	m, err := migrate.New(
		"file://migrations",
		connString,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create migrate instance: %w", op, err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("%s: failed to perform db migrations: %w", op, err)
	}

	return conn, nil
}
