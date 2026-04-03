package repository

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер для Postgres
)

func NewPostgresDB(host, port, user, password, dbname string) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var db *sqlx.DB
	var err error

	// Повторные попытки подключения (для Docker Compose)
	for i := 0; i < 10; i++ {
		db, err = sqlx.Open("postgres", dsn)
		if err != nil {
			slog.Warn("failed to open db", "attempt", i+1, "error", err)
		} else {
			err = db.Ping()
			if err == nil {
				slog.Debug("database connected")
				break
			}
			slog.Warn("database not ready", "attempt", i+1, "error", err)
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД после 10 попыток: %w", err)
	}

	// Запуск миграций
	if err := runMigrations(host, port, user, password, dbname); err != nil {
		return nil, fmt.Errorf("ошибка запуска миграций: %w", err)
	}

	return db, nil
}

func runMigrations(host, port, user, password, dbname string) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname)

	m, err := migrate.New("file://migrations", dsn)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	slog.Info("migrations applied successfully")
	return nil
}
