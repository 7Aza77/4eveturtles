package mocks

import (
	"context"
	"goevent/internal/entity"
	"goevent/internal/repository"

	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

func (m *AuthUseCaseMock) Register(ctx context.Context, email, password string) (int64, error) {
	args := m.Called(ctx, email, password)
	return args.Get(0).(int64), args.Error(1)
}

func (m *AuthUseCaseMock) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

type EventUseCaseMock struct {
	mock.Mock
}

func (m *EventUseCaseMock) Create(ctx context.Context, event entity.Event) (int64, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(int64), args.Error(1)
}

func (m *EventUseCaseMock) GetByID(ctx context.Context, id int64) (entity.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(entity.Event), args.Error(1)
}

func (m *EventUseCaseMock) List(ctx context.Context, filter repository.EventFilter) ([]entity.Event, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]entity.Event), args.Error(1)
}

func (m *EventUseCaseMock) Update(ctx context.Context, userId int64, role string, event entity.Event) error {
	args := m.Called(ctx, userId, role, event)
	return args.Error(0)
}

func (m *EventUseCaseMock) Delete(ctx context.Context, userId int64, role string, eventId int64) error {
	args := m.Called(ctx, userId, role, eventId)
	return args.Error(0)
}

type RegistrationUseCaseMock struct {
	mock.Mock
}

func (m *RegistrationUseCaseMock) Register(ctx context.Context, userId, eventId int64) error {
	args := m.Called(ctx, userId, eventId)
	return args.Error(0)
}

func (m *RegistrationUseCaseMock) Cancel(ctx context.Context, userId, eventId int64) error {
	args := m.Called(ctx, userId, eventId)
	return args.Error(0)
}
