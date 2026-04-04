package trip

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"planning-trip-be/internal/middleware"
	"planning-trip-be/internal/response"
	svc "planning-trip-be/internal/service/trip"
)

func RegisterRoutes(router gin.IRouter, service svc.Service, authMiddleware gin.HandlerFunc) {
	group := router.Group("/trip")
	group.GET("", authMiddleware, list(service))
	group.GET("/:id", get(service))
	group.POST("", authMiddleware, create(service))
	group.PUT("/:id", authMiddleware, notImplemented(service, "update"))
	group.DELETE("/:id", authMiddleware, notImplemented(service, "delete"))
}

func list(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		trips, err := service.ListByUser(c.Request.Context(), userID)
		if err != nil {
			writeServiceError(c, err, "cannot list trips")
			return
		}
		response.Write(c.Writer, http.StatusOK, trips, "")
	}
}

func create(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input svc.CreateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
				Err:     err,
			})
			return
		}

		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}
		input.CreatorID = userID

		trip, err := service.Create(c.Request.Context(), input)
		if err != nil {
			writeServiceError(c, err, "cannot create trip")
			return
		}
		response.Write(c.Writer, http.StatusCreated, trip, "create trip successfully")
	}
}

func get(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		tripID := c.Param("id")
		trip, err := service.GetByID(c.Request.Context(), tripID)
		if err != nil {
			writeServiceError(c, err, "cannot get trip")
			return
		}
		response.Write(c.Writer, http.StatusOK, trip, "")
	}
}

func writeServiceError(c *gin.Context, err error, fallbackMessage string) {
	if errors.Is(err, svc.ErrInvalidInput) {
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusNotFound,
			Message: "trip not found",
		})
		return
	}

	response.WriteError(c.Writer, response.APIError{
		Status:  http.StatusInternalServerError,
		Message: fallbackMessage,
		Err:     err,
	})
}

func notImplemented(_ svc.Service, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusNotImplemented,
			Message: action + " is not implemented yet",
		})
	}
}
