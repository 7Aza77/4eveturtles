package handler

import (
	"goevent/internal/metrics"
	"goevent/internal/usecase"
	"goevent/pkg/lib/api/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

type RegistrationHandler struct {
	useCase usecase.RegistrationUseCase
}

func NewRegistrationHandler(useCase usecase.RegistrationUseCase) *RegistrationHandler {
	return &RegistrationHandler{useCase: useCase}
}

func (h *RegistrationHandler) register(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid event id"))
		return
	}

	userIdRaw, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error("user not found"))
		return
	}

	userId, ok := userIdRaw.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error("invalid user id"))
		return
	}

	if err := h.useCase.Register(c.Request.Context(), userId, eventId); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	metrics.RegistrationsTotal.With(prometheus.Labels{"action": "register"}).Inc()
	c.JSON(http.StatusOK, response.OK())
}

func (h *RegistrationHandler) cancel(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid event id"))
		return
	}

	userIdRaw, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error("user not found"))
		return
	}

	userId, ok := userIdRaw.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error("invalid user id"))
		return
	}

	if err := h.useCase.Cancel(c.Request.Context(), userId, eventId); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	metrics.RegistrationsTotal.With(prometheus.Labels{"action": "cancel"}).Inc()
	c.JSON(http.StatusOK, response.OK())
}

func (h *RegistrationHandler) participants(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid event id"))
		return
	}

	ids, err := h.useCase.GetParticipants(c.Request.Context(), eventId)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status: response.StatusOk,
		Data:   gin.H{"event_id": eventId, "participants": ids, "count": len(ids)},
	})
}
