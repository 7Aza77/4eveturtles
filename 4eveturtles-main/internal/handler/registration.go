package handler

import (
	"goevent/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userId, _ := c.Get(userCtx)

	err = h.useCase.Register(c.Request.Context(), userId.(int64), eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "registered"})
}

func (h *RegistrationHandler) cancel(c *gin.Context) {
	eventId, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userId, _ := c.Get(userCtx)

	err := h.useCase.Cancel(c.Request.Context(), userId.(int64), eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "cancelled"})
}
