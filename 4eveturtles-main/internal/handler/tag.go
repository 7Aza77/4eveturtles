package handler

import (
	"goevent/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	useCase usecase.TagUseCase
}

func NewTagHandler(useCase usecase.TagUseCase) *TagHandler {
	return &TagHandler{useCase: useCase}
}

type createTagInput struct {
	Name string `json:"name" binding:"required"`
}

type addTagsInput struct {
	TagIDs []int64 `json:"tag_ids" binding:"required"`
}

// @Summary Create Tag
// @Security ApiKeyAuth
// @Tags tags
// @Description Create a new event tag (e.g., Tech, Sports, Arts)
// @Accept json
// @Produce json
// @Param input body createTagInput true "tag name"
// @Success 200 {object} map[string]int64
// @Router /api/v1/tags [post]
func (h *TagHandler) create(c *gin.Context) {
	var input createTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.useCase.Create(c.Request.Context(), input.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary List Tags
// @Tags tags
// @Description Get list of all available tags
// @Produce json
// @Success 200 {array} entity.Tag
// @Router /api/v1/tags [get]
func (h *TagHandler) list(c *gin.Context) {
	tags, err := h.useCase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tags)
}

// @Summary Add Tags to Event
// @Security ApiKeyAuth
// @Tags tags
// @Description Assign tag IDs to a specific event
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param input body addTagsInput true "tag ids"
// @Success 200 {object} map[string]string
// @Router /api/v1/events/{id}/tags [post]
func (h *TagHandler) addTagsToEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	var input addTagsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.useCase.AddTagsToEvent(c.Request.Context(), eventId, input.TagIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Summary Remove All Tags from Event
// @Security ApiKeyAuth
// @Tags tags
// @Description Remove all tag associations from a specific event
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/events/{id}/tags [delete]
func (h *TagHandler) removeTagsFromEvent(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	if err := h.useCase.RemoveTagsFromEvent(c.Request.Context(), eventId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
