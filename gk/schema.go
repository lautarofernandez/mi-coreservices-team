package gk

import (
	"bytes"
	"io/ioutil"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/gk/jsonschema"
	"github.com/mercadolibre/coreservices-team/libs/go/errors"
)

// JSONSchema is a middleware that accepts a JSON schema name  that must
// be a valid JSON Schema (Draft #6) definition. It then uses this schema
// to validate the request body. It returns status 422 on failure.
func JSONSchema(schemaName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			errors.ReturnError(c, &errors.Error{
				Code:    errors.InternalServerApiError,
				Message: "Error reading JSON body from request",
				Cause:   err.Error(),
			})
			c.Abort()
			return
		}
		c.Request.Body.Close()

		buf := bytes.NewReader(body)
		if err := jsonschema.Validate(schemaName, buf); err != nil {
			values := map[string]string{}

			if verr, ok := err.(*jsonschema.ValidationError); ok {
				values = verr.ErrorsDescription()
			}

			errors.ReturnError(c, &errors.Error{
				Code:    errors.UnprocessableEntityApiError,
				Message: "Error validating body to JSON schema",
				Cause:   "Validation error",
				Values:  values,
			})
			c.Abort()
			return
		}

		// Rewind body buffer, encapsulate it in a NopCloser, and assign it to request body again.
		buf.Seek(0, 0)
		c.Request.Body = ioutil.NopCloser(buf)

		c.Next()
	}
}
