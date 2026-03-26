package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Драйвер для Postgres
)

func NewPostgresDB(host, port, user, password, dbname string) (*sqlx.DB, error) {
	// Строка подключения (DSN)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Создаем таблицы, если они не существуют (для упрощения защиты)
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'student',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		date TIMESTAMP NOT NULL,
		location VARCHAR(255),
		max_participants INTEGER DEFAULT 0,
		creator_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS event_registrations (
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		event_id INTEGER REFERENCES events(id) ON DELETE CASCADE,
		registered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (user_id, event_id)
	);`
	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}
