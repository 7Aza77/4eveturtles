package repository

import (
	"fmt"
	"log"
	"time"

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
			log.Printf("Попытка %d: ошибка открытия БД: %s", i+1, err.Error())
		} else {
			err = db.Ping()
			if err == nil {
				log.Println("Успешное подключение к БД!")
				break
			}
			log.Printf("Попытка %d: БД еще не готова: %s", i+1, err.Error())
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД после 10 попыток: %w", err)
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
