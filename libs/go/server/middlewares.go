package server

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/libs/go/errors"
	"github.com/mercadolibre/go-meli-toolkit/mlauth"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

// Auth Middleware
// Handles basic auth callerID, isAdmin request values retrieving
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawCallerID := mlauth.GetCaller(c.Request)
		isAdmin := mlauth.IsCallerAdmin(c.Request)
		callerID, err := strconv.ParseUint(rawCallerID, 10, 64)

		if err != nil {
			errors.ReturnError(c, &errors.Error{
				Code:    errors.BadRequestApiError,
				Cause:   "parsing header value",
				Message: "invalid caller.id",
				Values: map[string]string{
					"caller.id": rawCallerID,
				},
			})

			c.Abort()
			return
		}

		c.Set("callerID", callerID)
		c.Set("isAdmin", isAdmin)
		c.Next()
	}
}

// RenameNewRelicTransaction Middleware
// Rename Newrelic transaction to controller method name
func RenameNewRelicTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		trx := nrgin.Transaction(c)
		if trx != nil {
			splitURL := strings.Split(c.HandlerName(), "/")
			if len(splitURL) > 0 {
				trx.SetName(splitURL[len(splitURL)-1])
			}
		}

		c.Next()
	}
}
