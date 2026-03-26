package handler

import (
	"goevent/internal/entity"
	"goevent/internal/repository"
	"goevent/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	useCase usecase.EventUseCase
}

func NewEventHandler(useCase usecase.EventUseCase) *EventHandler {
	return &EventHandler{useCase: useCase}
}

type createEventInput struct {
	Title           string `json:"title" binding:"required"`
	Description     string `json:"description"`
	Date            string `json:"date" binding:"required"` // В формате RFC3339
	Location        string `json:"location"`
	MaxParticipants int    `json:"max_participants"`
}

// @Summary Create Event
// @Security ApiKeyAuth
// @Tags events
// @Description create a new event
// @Accept  json
// @Produce  json
// @Param input body createEventInput true "event info"
// @Success 200 {integer} integer 1
// @Router /api/v1/events [post]
func (h *EventHandler) create(c *gin.Context) {
	var input createEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Парсинг даты (упрощенно)
	// eventDate, _ := time.Parse(time.RFC3339, input.Date)

	event := entity.Event{
		Title:           input.Title,
		Description:     input.Description,
		Location:        input.Location,
		MaxParticipants: input.MaxParticipants,
		CreatorID:       userId.(int64),
	}

	id, err := h.useCase.Create(c.Request.Context(), event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary List Events
// @Tags events
// @Description get list of events
// @Accept  json
// @Produce  json
// @Param limit query int false "limit" default(10)
// @Param offset query int false "offset" default(0)
// @Success 200 {array} entity.Event
// @Router /api/v1/events [get]
func (h *EventHandler) list(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort_by", "id")
	order := c.DefaultQuery("order", "asc")

	events, err := h.useCase.List(c.Request.Context(), repository.EventFilter{
		Limit:  limit,
		Offset: offset,
		SortBy: sortBy,
		Order:  order,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (h *EventHandler) getByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	event, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func (h *EventHandler) update(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	var input createEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, _ := c.Get(userCtx)

	event := entity.Event{
		ID:              id,
		Title:           input.Title,
		Description:     input.Description,
		Location:        input.Location,
		MaxParticipants: input.MaxParticipants,
	}

	err := h.useCase.Update(c.Request.Context(), userId.(int64), event)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *EventHandler) delete(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	userId, _ := c.Get(userCtx)

	err := h.useCase.Delete(c.Request.Context(), userId.(int64), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
