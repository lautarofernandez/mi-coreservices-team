package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/libs/go/errors"
)

// NoRouteHandler is a default handler that's usually used in conjunction with
// gins NoRoute method.
func NoRouteHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		errors.ReturnError(c, &errors.Error{
			Code:    errors.NotFoundApiError,
			Message: fmt.Sprintf("Resource %s not found.", c.Request.URL.Path),
			Values: map[string]string{
				"resource": c.Request.URL.Path,
			},
		})

		c.Abort()
	}
}

// BlockPublicTraffic returns 404 to the route when the request has the
// X-Public header set to true.
func BlockPublicTraffic() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header.Get("x-public") == "true" {
			c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  http.StatusNotFound,
				"message": "Not Found",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}

// HealthCheckHandler is a default handler that's used by Fury for checking if a
// given application instance is accepting requests.
func HealthCheckHandler(c *gin.Context) {
	c.String(http.StatusOK, "pong")
}
