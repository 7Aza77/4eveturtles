package handler

import (
	"goevent/internal/entity"
	"goevent/internal/repository"
	"goevent/internal/usecase"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	useCase usecase.EventUseCase
}

func NewEventHandler(useCase usecase.EventUseCase) *EventHandler {
	return &EventHandler{useCase: useCase}
}

type createEventInput struct {
	Title           string  `json:"title" binding:"required"`
	Description     string  `json:"description"`
	Date            string  `json:"date" binding:"required"` // RFC3339
	Location        string  `json:"location"`
	MaxParticipants int     `json:"max_participants"`
	OrganizationID  *int64  `json:"organization_id"`
	GroupChatLink   string  `json:"group_chat_link"`
}

// @Summary Create Event
// @Security ApiKeyAuth
// @Tags events
// @Description create a new event
// @Accept json
// @Produce json
// @Param input body createEventInput true "event info"
// @Success 200 {integer} integer 1
// @Router /api/v1/events [post]
func (h *EventHandler) create(c *gin.Context) {
	var input createEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdRaw, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	userId, ok := userIdRaw.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	eventDate, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date must be in RFC3339 format"})
		return
	}

	event := entity.Event{
		Title:           input.Title,
		Description:     input.Description,
		Date:            eventDate,
		Location:        input.Location,
		MaxParticipants: input.MaxParticipants,
		CreatorID:       userId,
		OrganizationID:  input.OrganizationID,
		GroupChatLink:   input.GroupChatLink,
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
// @Accept json
// @Produce json
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
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input createEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIdRaw, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	userId, ok := userIdRaw.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	eventDate, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date must be in RFC3339 format"})
		return
	}

	event := entity.Event{
		ID:              id,
		Title:           input.Title,
		Description:     input.Description,
		Date:            eventDate,
		Location:        input.Location,
		MaxParticipants: input.MaxParticipants,
		OrganizationID:  input.OrganizationID,
		GroupChatLink:   input.GroupChatLink,
	}

	if err := h.useCase.Update(c.Request.Context(), userId, event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (h *EventHandler) delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userIdRaw, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	userId, ok := userIdRaw.(int64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), userId, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
