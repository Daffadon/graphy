package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MigratePostgres(pool *pgxpool.Pool, s *slog.Logger) {
	sqlDB, err := sql.Open("postgres", pool.Config().ConnString())
	if err != nil {
		s.Error(fmt.Sprintf("failed to open sql.DB for migration: %v", err))
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		s.Error(fmt.Sprintf("failed to create migration driver: %v", err))
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/pkg/database/migrations/postgresql",
		"postgres", driver)
	if err != nil {
		s.Error(fmt.Sprintf("failed to create migrate instance: %v", err))
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		s.Error(fmt.Sprintf("migration failed: %v", err))
	}
}

func ClosePostgres(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
