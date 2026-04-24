package handler

import (
	"context"
	"goevent/internal/entity"
	"goevent/internal/repository"
	"github.com/stretchr/testify/mock"
)

// Описываем структуру-заглушку, которая имитирует EventUseCase
type MockEventUseCase struct {
	mock.Mock
}

// Имитируем метод Create
func (m *MockEventUseCase) Create(ctx context.Context, event entity.Event) (int64, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(int64), args.Error(1)
}

// Имитируем метод GetByID
func (m *MockEventUseCase) GetByID(ctx context.Context, id int64) (entity.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Event), args.Error(1)
}

// Имитируем метод List
func (m *MockEventUseCase) List(ctx context.Context, f repository.EventFilter) ([]entity.Event, error) {
	args := m.Called(ctx, f)
	return args.Get(0).([]entity.Event), args.Error(1)
}

// Добавляем методы Update и Delete, чтобы удовлетворить интерфейсу
func (m *MockEventUseCase) Update(ctx context.Context, userID int64, role string, event entity.Event) error {
	return m.Called(ctx, userID, role, event).Error(0)
}

func (m *MockEventUseCase) Delete(ctx context.Context, userID int64, role string, id int64) error {
	return m.Called(ctx, userID, role, id).Error(0)
}