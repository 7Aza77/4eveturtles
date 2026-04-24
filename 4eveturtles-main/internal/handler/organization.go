package handler

import (
	"goevent/internal/entity"
	"goevent/internal/usecase"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	useCase usecase.OrganizationUseCase
}

func NewOrganizationHandler(useCase usecase.OrganizationUseCase) *OrganizationHandler {
	return &OrganizationHandler{useCase: useCase}
}

type createOrgInput struct {
	Name          string `json:"name" binding:"required"`
	Description   string `json:"description"`
	UniversityID  string `json:"university_id"`
	GroupChatLink string `json:"group_chat_link"`
}

// @Summary Create Organization
// @Security ApiKeyAuth
// @Tags organizations
// @Description Create a new student organization
// @Accept json
// @Produce json
// @Param input body createOrgInput true "organization info"
// @Success 200 {object} map[string]int64
// @Router /api/v1/organizations [post]
func (h *OrganizationHandler) create(c *gin.Context) {
	var input createOrgInput
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

	org := entity.Organization{
		Name:          input.Name,
		Description:   input.Description,
		UniversityID:  input.UniversityID,
		GroupChatLink: input.GroupChatLink,
		OwnerID:       userId,
	}

	id, err := h.useCase.Create(c.Request.Context(), org)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

// @Summary List Organizations
// @Tags organizations
// @Description Get list of all student organizations
// @Produce json
// @Success 200 {array} entity.Organization
// @Router /api/v1/organizations [get]
func (h *OrganizationHandler) list(c *gin.Context) {
	orgs, err := h.useCase.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orgs)
}

// @Summary Get Organization by ID
// @Tags organizations
// @Description Get a specific organization by its ID
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} entity.Organization
// @Router /api/v1/organizations/{id} [get]
func (h *OrganizationHandler) getByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	org, err := h.useCase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organization not found"})
		return
	}

	c.JSON(http.StatusOK, org)
}

// @Summary Update Organization
// @Security ApiKeyAuth
// @Tags organizations
// @Description Update an existing organization
// @Accept json
// @Produce json
// @Param id path int true "Organization ID"
// @Param input body createOrgInput true "updated info"
// @Success 200 {object} map[string]string
// @Router /api/v1/organizations/{id} [put]
func (h *OrganizationHandler) update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input createOrgInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	org := entity.Organization{
		ID:            id,
		Name:          input.Name,
		Description:   input.Description,
		UniversityID:  input.UniversityID,
		GroupChatLink: input.GroupChatLink,
	}

	if err := h.useCase.Update(c.Request.Context(), org); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// @Summary Delete Organization
// @Security ApiKeyAuth
// @Tags organizations
// @Description Delete an organization by ID
// @Produce json
// @Param id path int true "Organization ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/organizations/{id} [delete]
func (h *OrganizationHandler) delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.useCase.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
