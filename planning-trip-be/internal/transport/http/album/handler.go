package album

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"planning-trip-be/internal/response"
	svc "planning-trip-be/internal/service/album"
)

func RegisterRoutes(router gin.IRouter, service svc.Service) {
	group := router.Group("/album")
	group.GET("", notImplemented(service, "list"))
	group.GET("/:id", notImplemented(service, "get"))
	group.POST("", notImplemented(service, "create"))
	group.PUT("/:id", notImplemented(service, "update"))
	group.DELETE("/:id", notImplemented(service, "delete"))
}

func notImplemented(_ svc.Service, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusNotImplemented,
			Message: action + " is not implemented yet",
		})
	}
}
