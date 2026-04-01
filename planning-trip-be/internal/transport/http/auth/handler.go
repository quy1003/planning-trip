package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"planning-trip-be/internal/response"
	asvc "planning-trip-be/internal/service/auth"
)

func RegisterRoutes(router gin.IRouter, svc asvc.Service) {
	router.POST("/auth/register", register(svc))
	router.POST("/auth/login", login(svc))
}

func register(svc asvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input asvc.RegisterInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
				Err:     err,
			})
			return
		}

		result, err := svc.Register(c.Request.Context(), input)
		if err != nil {
			writeServiceError(c, err, "cannot register user")
			return
		}

		response.Write(c.Writer, http.StatusCreated, result, "register successfully")
	}
}

func login(svc asvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input asvc.LoginInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
				Err:     err,
			})
			return
		}

		result, err := svc.Login(c.Request.Context(), input)
		if err != nil {
			writeServiceError(c, err, "cannot login")
			return
		}

		response.Write(c.Writer, http.StatusOK, result, "login successfully")
	}
}

func writeServiceError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, asvc.ErrInvalidInput):
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	case errors.Is(err, asvc.ErrEmailAlreadyUsed):
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusConflict,
			Message: err.Error(),
		})
	case errors.Is(err, asvc.ErrInvalidCredentials):
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusUnauthorized,
			Message: err.Error(),
		})
	default:
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusInternalServerError,
			Message: fallbackMessage,
			Err:     err,
		})
	}
}
