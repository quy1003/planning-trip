package trip

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"planning-trip-be/internal/middleware"
	"planning-trip-be/internal/response"
	svc "planning-trip-be/internal/service/trip"
)

func RegisterRoutes(router gin.IRouter, service svc.Service, authMiddleware gin.HandlerFunc) {
	group := router.Group("/trip")
	group.GET("", authMiddleware, list(service))
	group.GET("/:id", authMiddleware, getDetail(service))
	group.POST("", authMiddleware, create(service))
	group.PATCH("/:id/status", authMiddleware, updateStatus(service))
	group.GET("/:id/members", authMiddleware, listMembers(service))
	group.GET("/:id/schedule-items", authMiddleware, listScheduleItems(service))
	group.POST("/:id/schedule-items", authMiddleware, createScheduleItem(service))
	group.PUT("/:id/schedule-items/:item_id", authMiddleware, updateScheduleItem(service))
	group.DELETE("/:id/schedule-items/:item_id", authMiddleware, deleteScheduleItem(service))
	group.PATCH("/:id/schedule-items/reorder", authMiddleware, reorderScheduleItems(service))
	group.GET("/:id/album-preview", authMiddleware, listAlbumPreview(service))
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

func getDetail(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		tripID := c.Param("id")
		trip, err := service.GetDetail(c.Request.Context(), tripID, userID)
		if err != nil {
			writeServiceError(c, err, "cannot get trip detail")
			return
		}
		response.Write(c.Writer, http.StatusOK, trip, "")
	}
}

func updateStatus(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		var input svc.UpdateStatusInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
				Err:     err,
			})
			return
		}

		tripID := c.Param("id")
		trip, err := service.UpdateStatus(c.Request.Context(), tripID, userID, input)
		if err != nil {
			writeServiceError(c, err, "cannot update trip status")
			return
		}

		response.Write(c.Writer, http.StatusOK, trip, "update trip status successfully")
	}
}

func listMembers(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		tripID := c.Param("id")
		members, err := service.ListMembers(c.Request.Context(), tripID, userID)
		if err != nil {
			writeServiceError(c, err, "cannot list trip members")
			return
		}

		response.Write(c.Writer, http.StatusOK, members, "")
	}
}

func listScheduleItems(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		var dayIndex *int
		if raw := c.Query("day_index"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil {
				response.WriteError(c.Writer, response.APIError{
					Status:  http.StatusBadRequest,
					Message: "invalid day_index",
					Err:     err,
				})
				return
			}
			dayIndex = &parsed
		}

		tripID := c.Param("id")
		items, err := service.ListScheduleItems(c.Request.Context(), tripID, userID, dayIndex)
		if err != nil {
			writeServiceError(c, err, "cannot list schedule items")
			return
		}

		response.Write(c.Writer, http.StatusOK, items, "")
	}
}

func createScheduleItem(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		var input svc.CreateScheduleItemInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
				Err:     err,
			})
			return
		}

		tripID := c.Param("id")
		item, err := service.CreateScheduleItem(c.Request.Context(), tripID, userID, input)
		if err != nil {
			writeServiceError(c, err, "cannot create schedule item")
			return
		}

		response.Write(c.Writer, http.StatusCreated, item, "create schedule item successfully")
	}
}

func updateScheduleItem(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		var input svc.UpdateScheduleItemInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
				Err:     err,
			})
			return
		}

		tripID := c.Param("id")
		itemID := c.Param("item_id")
		item, err := service.UpdateScheduleItem(c.Request.Context(), tripID, itemID, userID, input)
		if err != nil {
			writeServiceError(c, err, "cannot update schedule item")
			return
		}

		response.Write(c.Writer, http.StatusOK, item, "update schedule item successfully")
	}
}

func deleteScheduleItem(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		tripID := c.Param("id")
		itemID := c.Param("item_id")
		if err := service.DeleteScheduleItem(c.Request.Context(), tripID, itemID, userID); err != nil {
			writeServiceError(c, err, "cannot delete schedule item")
			return
		}

		response.Write(c.Writer, http.StatusOK, gin.H{}, "delete schedule item successfully")
	}
}

func reorderScheduleItems(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		var input svc.ReorderScheduleInput
		if err := c.ShouldBindJSON(&input); err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid request body",
				Err:     err,
			})
			return
		}

		tripID := c.Param("id")
		if err := service.ReorderSchedule(c.Request.Context(), tripID, userID, input); err != nil {
			writeServiceError(c, err, "cannot reorder schedule items")
			return
		}

		response.Write(c.Writer, http.StatusOK, gin.H{}, "reorder schedule items successfully")
	}
}

func listAlbumPreview(service svc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := middleware.AuthUserID(c)
		if !ok {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusUnauthorized,
				Message: "unauthorized",
			})
			return
		}

		limit := 12
		if raw := c.Query("limit"); raw != "" {
			parsed, err := strconv.Atoi(raw)
			if err != nil {
				response.WriteError(c.Writer, response.APIError{
					Status:  http.StatusBadRequest,
					Message: "invalid limit",
					Err:     err,
				})
				return
			}
			limit = parsed
		}

		tripID := c.Param("id")
		photos, err := service.ListAlbumPreview(c.Request.Context(), tripID, userID, limit)
		if err != nil {
			writeServiceError(c, err, "cannot list album preview")
			return
		}

		response.Write(c.Writer, http.StatusOK, photos, "")
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

	if errors.Is(err, svc.ErrForbidden) {
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusForbidden,
			Message: "forbidden",
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
