package mocks

import (
	"context"
	"goevent/internal/entity"
	"goevent/internal/repository"

	"github.com/stretchr/testify/mock"
)

type EventMock struct {
	mock.Mock
}

func (m *EventMock) Create(ctx context.Context, event entity.Event) (int64, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(int64), args.Error(1)
}

func (m *EventMock) GetByID(ctx context.Context, id int64) (entity.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Event), args.Error(1)
}

func (m *EventMock) List(ctx context.Context, filter repository.EventFilter) ([]entity.Event, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entity.Event), args.Error(1)
}

func (m *EventMock) Update(ctx context.Context, event entity.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *EventMock) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
