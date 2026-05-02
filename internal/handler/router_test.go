package handler

import (
	"context"
	"goevent/internal/entity"
	"goevent/internal/usecase/mocks"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPingRoute(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil)
	router := h.InitRouter(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "pong")
}

func TestEventsListRoute(t *testing.T) {
	mockUseCase := new(mocks.EventUseCaseMock)
	h := NewHandler(nil, mockUseCase, nil, nil)
	router := h.InitRouter(nil)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("List", mock.Anything, mock.AnythingOfType("repository.EventFilter")).
			Return([]entity.Event{}, nil).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/events", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase.On("List", mock.Anything, mock.AnythingOfType("repository.EventFilter")).
			Return([]entity.Event{}, context.DeadlineExceeded).Once()

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/events", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockUseCase.AssertExpectations(t)
	})
}

func TestAuthSignUpRoute(t *testing.T) {
	mockAuth := new(mocks.AuthUseCaseMock)
	h := NewHandler(mockAuth, nil, nil, nil)
	router := h.InitRouter(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/sign-up", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateEventRequiresAuth(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil)
	router := h.InitRouter(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/events", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
