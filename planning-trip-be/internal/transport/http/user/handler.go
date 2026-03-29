package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"planning-trip-be/internal/response"
	usvc "planning-trip-be/internal/service/user"

	"gorm.io/gorm"
)

func RegisterRoutes(router gin.IRouter, svc usvc.Service) {
	router.GET("/users", listUsers(svc))
	router.GET("/users/:id", getUser(svc))
	router.POST("/users", createUser(svc))
	router.PUT("/users/:id", updateUser(svc))
	router.DELETE("/users/:id", deleteUser(svc))
}

func listUsers(svc usvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := svc.List(c.Request.Context())
		if err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusInternalServerError,
				Message: "cannot list users",
				Err:     err,
			})
			return
		}
		response.Write(c.Writer, http.StatusOK, users, "")
	}
}

func getUser(svc usvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		user, err := svc.GetByID(c.Request.Context(), userID)
		if err != nil {
			writeServiceError(c, err, "cannot get user")
			return
		}
		response.Write(c.Writer, http.StatusOK, user, "")
	}
}

func createUser(svc usvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input usvc.CreateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
				Err:     err,
			})
			return
		}

		user, err := svc.Create(c.Request.Context(), input)
		if err != nil {
			writeServiceError(c, err, "cannot create user")
			return
		}
		response.Write(c.Writer, http.StatusCreated, user, "")
	}
}

func updateUser(svc usvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		var input usvc.UpdateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
				Err:     err,
			})
			return
		}

		user, err := svc.Update(c.Request.Context(), userID, input)
		if err != nil {
			writeServiceError(c, err, "cannot update user")
			return
		}
		response.Write(c.Writer, http.StatusOK, user, "")
	}
}

func deleteUser(svc usvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		if err := svc.Delete(c.Request.Context(), userID); err != nil {
			writeServiceError(c, err, "cannot delete user")
			return
		}
		response.Write(c.Writer, http.StatusOK, gin.H{"deleted": true}, "")
	}
}

func writeServiceError(c *gin.Context, err error, fallbackMessage string) {
	if errors.Is(err, usvc.ErrInvalidInput) {
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusNotFound,
			Message: "user not found",
		})
		return
	}

	response.WriteError(c.Writer, response.APIError{
		Status:  http.StatusInternalServerError,
		Message: fallbackMessage,
		Err:     err,
	})
}
