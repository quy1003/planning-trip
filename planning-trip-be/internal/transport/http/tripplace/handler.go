package tripplace

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"planning-trip-be/internal/middleware"
	"planning-trip-be/internal/response"
	svc "planning-trip-be/internal/service/tripplace"
)

func RegisterRoutes(router gin.IRouter, service svc.Service, authMiddleware gin.HandlerFunc) {
	group := router.Group("/tripplace")
	group.GET("", authMiddleware, list(service))
	group.GET("/:id", authMiddleware, get(service))
	group.POST("", authMiddleware, create(service))
	group.PUT("/:id", authMiddleware, update(service))
	group.DELETE("/:id", authMiddleware, deleteTripPlace(service))
}

func list(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusUnauthorized, Message: "unauthorized"})
			return
		}
		tripID := c.Query("trip_id")
		if tripID == "" {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusBadRequest, Message: "trip_id is required"})
			return
		}

		var dayIndex *int
		if raw := c.Query("day_index"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil {
				response.WriteError(c.Writer, response.APIError{Status: http.StatusBadRequest, Message: "invalid day_index", Err: err})
				return
			}
			dayIndex = &parsed
		}

		rows, err := service.ListByTrip(c.Request.Context(), tripID, userID, dayIndex)
		if err != nil {
			writeServiceError(c, err, "cannot list trip places")
			return
		}
		response.Write(c.Writer, http.StatusOK, rows, "")
	}
}

func get(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusUnauthorized, Message: "unauthorized"})
			return
		}
		row, err := service.GetByID(c.Request.Context(), c.Param("id"), userID)
		if err != nil {
			writeServiceError(c, err, "cannot get trip place")
			return
		}
		response.Write(c.Writer, http.StatusOK, row, "")
	}
}

func create(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusUnauthorized, Message: "unauthorized"})
			return
		}

		var input svc.CreateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusBadRequest, Message: "invalid request body", Err: err})
			return
		}
		input.CreatedBy = userID

		row, err := service.Create(c.Request.Context(), input)
		if err != nil {
			writeServiceError(c, err, "cannot create trip place")
			return
		}
		response.Write(c.Writer, http.StatusCreated, row, "create trip place successfully")
	}
}

func update(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusUnauthorized, Message: "unauthorized"})
			return
		}
		var input svc.UpdateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusBadRequest, Message: "invalid request body", Err: err})
			return
		}
		row, err := service.Update(c.Request.Context(), c.Param("id"), userID, input)
		if err != nil {
			writeServiceError(c, err, "cannot update trip place")
			return
		}
		response.Write(c.Writer, http.StatusOK, row, "update trip place successfully")
	}
}

func deleteTripPlace(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusUnauthorized, Message: "unauthorized"})
			return
		}
		if err := service.Delete(c.Request.Context(), c.Param("id"), userID); err != nil {
			writeServiceError(c, err, "cannot delete trip place")
			return
		}
		response.Write(c.Writer, http.StatusOK, gin.H{}, "delete trip place successfully")
	}
}

func writeServiceError(c *gin.Context, err error, fallbackMessage string) {
	if errors.Is(err, svc.ErrInvalidInput) {
		response.WriteError(c.Writer, response.APIError{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}
	if errors.Is(err, svc.ErrForbidden) {
		response.WriteError(c.Writer, response.APIError{Status: http.StatusForbidden, Message: "forbidden"})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.WriteError(c.Writer, response.APIError{Status: http.StatusNotFound, Message: "trip place not found"})
		return
	}
	response.WriteError(c.Writer, response.APIError{Status: http.StatusInternalServerError, Message: fallbackMessage, Err: err})
}
