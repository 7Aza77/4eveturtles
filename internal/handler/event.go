package handler

import (
	"goevent/internal/entity"
	"goevent/internal/metrics"
	"goevent/internal/repository"
	"goevent/internal/usecase"
	"goevent/pkg/lib/api/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type EventHandler struct {
	useCase usecase.EventUseCase
}

func NewEventHandler(useCase usecase.EventUseCase) *EventHandler {
	return &EventHandler{useCase: useCase}
}

type createEventInput struct {
	Title           string `json:"title"            binding:"required"`
	Description     string `json:"description"`
	Date            string `json:"date"             binding:"required"`
	Location        string `json:"location"`
	MaxParticipants int    `json:"max_participants"  binding:"min=0"`
}

// @Summary Create Event
// @Security ApiKeyAuth
// @Description create a new event
// @Tags events
// @Accept json
// @Produce json
// @Param input body createEventInput true "event info"
// @Success 200 {object} response.Response
// @Failure 400,401 {object} response.Response
// @Router /api/v1/events/ [post]
func (h *EventHandler) create(c *gin.Context) {
	var input createEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, response.ValidationError(errs))
			return
		}
		c.JSON(http.StatusBadRequest, response.Error("invalid input body"))
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

	eventDate, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("date must be in RFC3339 format, e.g. 2025-06-01T15:00:00Z"))
		return
	}

	event := entity.Event{
		Title:           input.Title,
		Description:     input.Description,
		Date:            eventDate,
		Location:        input.Location,
		MaxParticipants: input.MaxParticipants,
		CreatorID:       userId,
	}

	id, err := h.useCase.Create(c.Request.Context(), event)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	metrics.EventsCreatedTotal.Inc()
	c.JSON(http.StatusOK, response.Response{
		Status: response.StatusOk,
		Data:   gin.H{"id": id},
	})
}

// @Summary List Events
// @Description get list of events with optional filtering and pagination
// @Tags events
// @Accept json
// @Produce json
// @Param limit query int false "limit" default(10)
// @Param offset query int false "offset" default(0)
// @Param sort_by query string false "sort field" Enums(id, title, date, location)
// @Param order query string false "order" Enums(asc, desc)
// @Param title query string false "filter by title"
// @Param location query string false "filter by location"
// @Param from_date query string false "filter from date (RFC3339)"
// @Param to_date query string false "filter to date (RFC3339)"
// @Success 200 {object} response.Response
// @Router /api/v1/events/ [get]
func (h *EventHandler) list(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	sortBy := c.DefaultQuery("sort_by", "id")
	order := c.DefaultQuery("order", "asc")
	title := c.Query("title")
	location := c.Query("location")
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")

	events, err := h.useCase.List(c.Request.Context(), repository.EventFilter{
		Limit:    limit,
		Offset:   offset,
		SortBy:   sortBy,
		Order:    order,
		Title:    title,
		Location: location,
		FromDate: fromDate,
		ToDate:   toDate,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error("failed to retrieve events"))
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status: response.StatusOk,
		Data:   events,
	})
}

// @Summary Get Event By ID
// @Description get a single event by id
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/events/{id} [get]
func (h *EventHandler) getByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	event, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("event not found"))
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status: response.StatusOk,
		Data:   event,
	})
}

// @Summary Update Event
// @Security ApiKeyAuth
// @Description update an event (creator only)
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param input body createEventInput true "event info"
// @Success 200 {object} response.Response
// @Failure 400,403 {object} response.Response
// @Router /api/v1/events/{id} [put]
func (h *EventHandler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
		return
	}

	var input createEventInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, response.ValidationError(errs))
			return
		}
		c.JSON(http.StatusBadRequest, response.Error("invalid input body"))
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

	eventDate, err := time.Parse(time.RFC3339, input.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("date must be in RFC3339 format"))
		return
	}

	roleRaw, ok := c.Get(roleCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error("role not found"))
		return
	}
	role, ok := roleRaw.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error("invalid role"))
		return
	}

	event := entity.Event{
		ID:              id,
		Title:           input.Title,
		Description:     input.Description,
		Date:            eventDate,
		Location:        input.Location,
		MaxParticipants: input.MaxParticipants,
	}

	if err := h.useCase.Update(c.Request.Context(), userId, role, event); err != nil {
		c.JSON(http.StatusForbidden, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK())
}

// @Summary Delete Event
// @Security ApiKeyAuth
// @Description delete an event (creator only)
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/events/{id} [delete]
func (h *EventHandler) delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid id"))
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

	roleRaw, ok := c.Get(roleCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error("role not found"))
		return
	}
	role, ok := roleRaw.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Error("invalid role"))
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), userId, role, id); err != nil {
		c.JSON(http.StatusForbidden, response.Error(err.Error()))
		return
	}

	metrics.EventsDeletedTotal.Inc()
	c.JSON(http.StatusOK, response.OK())
}
