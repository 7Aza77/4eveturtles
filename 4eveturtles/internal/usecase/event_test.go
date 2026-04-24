package usecase

import (
	"context"
	"errors"
	"goevent/internal/entity"
	"goevent/internal/repository"
	"goevent/internal/repository/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEvent_GetByID(t *testing.T) {
	mockRepo := new(mocks.EventMock)
	useCase := NewEvent(mockRepo, nil) // Redis is nil for now

	t.Run("success", func(t *testing.T) {
		expectedEvent := entity.Event{ID: 1, Title: "Test Event"}
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedEvent, nil).Once()

		event, err := useCase.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedEvent.Title, event.Title)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByID", mock.Anything, int64(2)).Return(entity.Event{}, errors.New("not found")).Once()

		_, err := useCase.GetByID(context.Background(), 2)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestEvent_List(t *testing.T) {
	mockRepo := new(mocks.EventMock)
	useCase := NewEvent(mockRepo, nil)

	t.Run("success", func(t *testing.T) {
		expectedEvents := []entity.Event{{ID: 1, Title: "E1"}, {ID: 2, Title: "E2"}}
		filter := repository.EventFilter{Limit: 10, Offset: 0}
		mockRepo.On("List", mock.Anything, filter).Return(expectedEvents, nil).Once()

		events, err := useCase.List(context.Background(), filter)

		assert.NoError(t, err)
		assert.Len(t, events, 2)
		mockRepo.AssertExpectations(t)
	})
}
