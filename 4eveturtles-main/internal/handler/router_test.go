package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil, nil, nil)
	router := h.InitRouter(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "pong"}`, w.Body.String())
}

func TestEventsRoute(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil, nil, nil)
	router := h.InitRouter(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/events/", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status": "Список мероприятий скоро будет здесь"}`, w.Body.String())
}
