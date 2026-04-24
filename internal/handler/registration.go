package handler

import (
	"goevent/internal/usecase"
	"goevent/pkg/lib/api/response"
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

// @Summary Register for Event
// @Security ApiKeyAuth
// @Description register user for an event
// @Tags events
// @ID register-event
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} response.Response
// @Failure 400,401 {object} response.Response
// @Router /api/v1/events/{id}/register [post]
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

	c.JSON(http.StatusOK, response.OK())
}

// @Summary Cancel Event Registration
// @Security ApiKeyAuth
// @Description unregister user from an event
// @Tags events
// @ID cancel-event
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} response.Response
// @Failure 400,401 {object} response.Response
// @Router /api/v1/events/{id}/unregister [delete]
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

	c.JSON(http.StatusOK, response.OK())
}
