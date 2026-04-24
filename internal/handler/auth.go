package handler

import (
	"goevent/internal/usecase"
	"goevent/pkg/lib/api/response"
	"net/http"

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
		c.JSON(http.StatusInternalServerError, response.Error("failed to create user"))
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
