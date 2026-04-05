package place

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"planning-trip-be/internal/middleware"
	"planning-trip-be/internal/response"
	svc "planning-trip-be/internal/service/place"
)

func RegisterRoutes(router gin.IRouter, service svc.Service, authMiddleware gin.HandlerFunc) {
	group := router.Group("/place")
	group.GET("", authMiddleware, list(service))
	group.GET("/:id", authMiddleware, get(service))
	group.POST("", authMiddleware, create(service))
	group.PUT("/:id", authMiddleware, update(service))
	group.DELETE("/:id", authMiddleware, deletePlace(service))
}

func list(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 20
		if raw := c.Query("limit"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil {
				response.WriteError(c.Writer, response.APIError{Status: http.StatusBadRequest, Message: "invalid limit", Err: err})
				return
			}
			limit = parsed
		}

		rows, err := service.List(c.Request.Context(), c.Query("q"), limit)
		if err != nil {
			writeServiceError(c, err, "cannot list places")
			return
		}
		response.Write(c.Writer, http.StatusOK, rows, "")
	}
}

func get(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		row, err := service.GetByID(c.Request.Context(), c.Param("id"))
		if err != nil {
			writeServiceError(c, err, "cannot get place")
			return
		}
		response.Write(c.Writer, http.StatusOK, row, "")
	}
}

func create(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input svc.CreateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusBadRequest, Message: "invalid request body", Err: err})
			return
		}

		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusUnauthorized, Message: "unauthorized"})
			return
		}
		input.CreatedBy = &userID

		row, err := service.Create(c.Request.Context(), input)
		if err != nil {
			writeServiceError(c, err, "cannot create place")
			return
		}
		response.Write(c.Writer, http.StatusCreated, row, "create place successfully")
	}
}

func update(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input svc.UpdateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{Status: http.StatusBadRequest, Message: "invalid request body", Err: err})
			return
		}
		row, err := service.Update(c.Request.Context(), c.Param("id"), input)
		if err != nil {
			writeServiceError(c, err, "cannot update place")
			return
		}
		response.Write(c.Writer, http.StatusOK, row, "update place successfully")
	}
}

func deletePlace(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := service.Delete(c.Request.Context(), c.Param("id")); err != nil {
			writeServiceError(c, err, "cannot delete place")
			return
		}
		response.Write(c.Writer, http.StatusOK, gin.H{}, "delete place successfully")
	}
}

func writeServiceError(c *gin.Context, err error, fallbackMessage string) {
	if errors.Is(err, svc.ErrInvalidInput) {
		response.WriteError(c.Writer, response.APIError{Status: http.StatusBadRequest, Message: err.Error()})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.WriteError(c.Writer, response.APIError{Status: http.StatusNotFound, Message: "place not found"})
		return
	}
	response.WriteError(c.Writer, response.APIError{Status: http.StatusInternalServerError, Message: fallbackMessage, Err: err})
}
