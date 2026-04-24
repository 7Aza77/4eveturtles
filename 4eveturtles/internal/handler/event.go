package handler

import (
	"goevent/internal/entity"
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

	c.JSON(http.StatusOK, response.Response{
		Status: response.StatusOk,
		Data:   gin.H{"id": id},
	})
}

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

	event := entity.Event{
		ID:              id,
		Title:           input.Title,
		Description:     input.Description,
		Date:            eventDate,
		Location:        input.Location,
		MaxParticipants: input.MaxParticipants,
	}

	if err := h.useCase.Update(c.Request.Context(), userId, event); err != nil {
		c.JSON(http.StatusForbidden, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK())
}

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

	if err := h.useCase.Delete(c.Request.Context(), userId, id); err != nil {
		c.JSON(http.StatusForbidden, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.OK())
}
