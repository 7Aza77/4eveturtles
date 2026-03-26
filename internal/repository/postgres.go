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

	return db, nil
}
