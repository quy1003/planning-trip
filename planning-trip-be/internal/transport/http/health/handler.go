package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"planning-trip-be/internal/response"
	hsvc "planning-trip-be/internal/service/health"
)

func Handler(svc hsvc.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := svc.Check(c.Request.Context())
		if err != nil {
			response.WriteError(c.Writer, response.APIError{
				Status:  http.StatusInternalServerError,
				Message: "health check failed",
				Err:     err,
			})
			return
		}

		response.Write(c.Writer, http.StatusOK, res, "")
	}
}
