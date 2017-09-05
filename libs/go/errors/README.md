MercadoPago CoreServices Errors
===

MPCS errors library for golang. 

Instalation
---

Use go get to fetch the package:

$ go get github.com/mercadolibre/coreservices-team/libs/go/errors

Usage
---

Import and use it. 

```
package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/libs/go/errors"
	"net/http"
)

// Ping func
// health status
func Ping(c *gin.Context) {
	if true { // must return an error
		errors.ReturnError(c, &errors.Error{
			Code: errors.AuthorizationApiError,
			Message: "Ping not authorized",
		})
	}

	c.String(http.StatusOK, "pong")
}
``` 

Methods
---

`ReturnError(c *gin.Context, err *Error)`

Returns an error to Gin Gonic.


Changelog
---

0.0.1 - 2017-08-14 

- Initial commit. 

Author
---

Core Services Team (`coreservices@mercadolibre.com`)
Slack (`https://mercadopago-team.slack.com/messages/C45S2LB5K`)
