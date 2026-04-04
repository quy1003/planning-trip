package upload

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"planning-trip-be/internal/response"
	usvc "planning-trip-be/internal/service/upload"
)

func RegisterRoutes(router gin.IRouter, service usvc.Service, authMiddleware gin.HandlerFunc) {
	group := router.Group("/upload")
	group.POST("/image", authMiddleware, uploadImage(service))
}

func uploadImage(service usvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "missing file",
				Err:     err,
			})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusBadRequest,
				Message: "invalid file",
				Err:     err,
			})
			return
		}
		defer file.Close()

		result, err := service.UploadImage(c.Request.Context(), file, fileHeader.Filename)
		if err != nil {
			writeServiceError(c, err)
			return
		}

		response.Write(c.Writer, http.StatusOK, result, "upload image successfully")
	}
}

func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usvc.ErrInvalidInput):
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
	case errors.Is(err, usvc.ErrNotConfigured):
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	case errors.Is(err, usvc.ErrUploadFailed):
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusBadGateway,
			Message: err.Error(),
		})
	default:
		response.WriteError(c.Writer, response.APIError{
			Status:  http.StatusInternalServerError,
			Message: "cannot upload image",
			Err:     err,
		})
	}
}
