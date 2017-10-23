package server

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/libs/go/errors"
	"github.com/mercadolibre/go-meli-toolkit/mlauth"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

// Auth Middleware. It checks that either the caller id or an admin scope is present
// in the request. If neither is present, it fails with 400 Bad Request.
// If prerequisites are met, then the found values are added to Gins context.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawCallerID := mlauth.GetCaller(c.Request)
		isAdmin := mlauth.IsCallerAdmin(c.Request)

		callerID, err := strconv.ParseUint(rawCallerID, 10, 64)

		// If request is not from an admin, and we failed parsing caller ID, fail
		if !isAdmin && err != nil {
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
