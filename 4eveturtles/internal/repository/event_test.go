package repository

import (
	"context"
	"goevent/internal/entity"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock" // Мок для имитации базы данных
	"github.com/jmoiron/sqlx"       // Твой драйвер
	"github.com/stretchr/testify/assert"
)

// Тест №1: Проверяем получение события по ID
func TestEventPostgres_GetByID(t *testing.T) {
	// 1. Инициализируем мок базы данных
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Оборачиваем в sqlx.DB, так как это используется в твоем EventPostgres
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewEventPostgres(sqlxDB)

	expectedEvent := entity.Event{
		ID:       1,
		Title:    "Концерт в КБТУ",
		Location: "Двор",
		Date:     time.Now(),
	}

	// Настраиваем ожидание SELECT запроса
	rows := sqlmock.NewRows([]string{"id", "title", "description", "date", "location", "max_participants", "creator_id"}).
		AddRow(expectedEvent.ID, expectedEvent.Title, "", expectedEvent.Date, expectedEvent.Location, 100, 1)
	
	mock.ExpectQuery("SELECT \\* FROM events WHERE id = \\$1").
		WithArgs(int64(1)).
		WillReturnRows(rows)

	// Выполняем
	res, err := repo.GetByID(context.Background(), 1)

	// Проверяем результат
	assert.NoError(t, err)
	assert.Equal(t, expectedEvent.Title, res.Title)
}

// Тест №2: Проверяем создание нового события (Критично для QA)
func TestEventPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := NewEventPostgres(sqlxDB)

	now := time.Now()
	event := entity.Event{
		Title:           "Хакатон 2026",
		Description:     "Главное событие для разработчиков",
		Date:            now,
		Location:        "KBTU",
		MaxParticipants: 50,
		CreatorID:       123,
	}

	// Настраиваем ожидание INSERT запроса. 
	// Мы ждем, что база вернет нам ID = 100 после вставки.
	mock.ExpectQuery("INSERT INTO events").
		WithArgs(event.Title, event.Description, event.Date, event.Location, event.MaxParticipants, event.CreatorID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(100)))

	// Выполняем
	id, err := repo.Create(context.Background(), event)

	// Проверяем
	assert.NoError(t, err)
	assert.Equal(t, int64(100), id, "Должен вернуться ID созданного события")
}