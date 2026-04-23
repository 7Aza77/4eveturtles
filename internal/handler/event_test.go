package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestEventHandler_Create(t *testing.T) {
	// Переводим Gin в тестовый режим, чтобы не спамить в консоль лишними логами
	gin.SetMode(gin.TestMode)

	t.Run("Invalid Date Format", func(t *testing.T) {
		// 1. Инициализируем мок бизнес-логики и сам хендлер
		mockUC := new(MockEventUseCase)
		h := NewEventHandler(mockUC)

		// 2. Настраиваем тестовый роутер
		r := gin.New()

		// ВАЖНО: Добавляем анонимный Middleware. 
		// Он имитирует работу Auth-Middleware, добавляя ID пользователя в контекст.
		// Без этого хендлер сразу вернет 401 Unauthorized.
		r.Use(func(c *gin.Context) {
			// Ключ "userId" должен совпадать с тем, что используется в handler/event.go (userCtx)
			c.Set(userCtx, int64(1)) 
			c.Next()
		})

		r.POST("/events", h.create)

		// 3. Подготавливаем "плохие" входные данные
		// Мы специально передаем дату в неверном формате (без времени и часового пояса), 
		// чтобы сработала ошибка валидации в хендлере.
		badInput := map[string]interface{}{
			"title": "Событие с плохой датой",
			"date":  "2025-06-01", // Правильный формат: RFC3339 (например, 2025-06-01T15:00:00Z)
		}
		
		body, _ := json.Marshal(badInput)
		
		// Создаем HTTP-запрос
		req, _ := http.NewRequest("POST", "/events", bytes.NewBuffer(body))
		// Recorder будет записывать ответ сервера
		w := httptest.NewRecorder()

		// 4. Запускаем выполнение запроса
		r.ServeHTTP(w, req)

		// 5. ПРОВЕРКА (Assertion)
		// Теперь, когда мы "авторизованы", хендлер дойдет до парсинга даты, 
		// увидит ошибку и вернет статус 400.
		assert.Equal(t, http.StatusBadRequest, w.Code, "Должен вернуться статус 400, так как формат даты неверный")
	})
}