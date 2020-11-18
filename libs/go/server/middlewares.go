package server

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/libs/go/errors"
	"github.com/mercadolibre/go-meli-toolkit/mlauth"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
)

const (
	apiRulesTestHeader = "X-Gordik-Mode"
	apiRulesTestValue  = "endpoint-test"
)

// MeliAPIRules middleware encapsulates all endpoints of Gordik and is used specifically
// for returning 200 when a specific header is set. The idea is to use this header
// whenever we want to test the existence of an endpoint and not its functionality.
func MeliAPIRules() gin.HandlerFunc {
	return func(c *gin.Context) {
		v := c.GetHeader(apiRulesTestHeader)
		if v == apiRulesTestValue {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

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
