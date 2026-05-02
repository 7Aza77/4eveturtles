package handler

import (
	"goevent/internal/usecase"
	"goevent/pkg/lib/api/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	useCase usecase.AuthUseCase
}

func NewAuthHandler(useCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase: useCase}
}

type signUpInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) signUp(c *gin.Context) {
	var input signUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, response.ValidationError(errs))
			return
		}
		c.JSON(http.StatusBadRequest, response.Error("invalid input"))
		return
	}

	id, err := h.useCase.Register(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusConflict, response.Error("user with this email already exists"))
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status: response.StatusOk,
		Data:   gin.H{"id": id},
	})
}

type signInInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) signIn(c *gin.Context) {
	var input signInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, response.ValidationError(errs))
			return
		}
		c.JSON(http.StatusBadRequest, response.Error("invalid input"))
		return
	}

	token, err := h.useCase.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error("invalid credentials"))
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Status: response.StatusOk,
		Data:   gin.H{"token": token},
	})
}

func (h *AuthHandler) me(c *gin.Context) {
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

	user, err := h.useCase.Me(c.Request.Context(), userId)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error("user not found"))
		return
	}

	user.Password = ""

	c.JSON(http.StatusOK, response.Response{
		Status: response.StatusOk,
		Data:   user,
	})
}

func (h *AuthHandler) getParticipantsCount(c *gin.Context) {
	eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error("invalid event id"))
		return
	}
	_ = eventId
}
