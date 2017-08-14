MercadoPago CoreServices Logger
===

MPCS logger library for golang. 

Instalation
---

Use go get to fetch the package:

$ go get github.com/mercadolibre/coreservices-team/libs/go/logger

Usage
---

Import and use it. 

```
import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/mercadolibre/coreservices-team/libs/go/logger"
)

func Ping(c *gin.Context) {
	logger := logger.LoggerWithName(c, "Ping")

	logger.Info("PingLog", logger.Attrs{
		"atts_1": "value 1",
		"...":    "value ...",
		"atts_n": "value n",
	})

	c.String(http.StatusOK, "pong")
}
```

Methods
---

`LoggerWithName(c *gin.Context, name string) *Logger`

Returns a logger pointer with gin context information. The name is used for more information. 

`(l *Logger) LogWithLevel(level string, event string, attrs ...Attrs) *Logger`

Basic log method. Level `INFO`, `DEBUG`, `WARN` or `ERROR`, and an event are required. 
Optional attrs of type `Attrs - map[string]interface{}` can be passed.

Direct methods for each level are provided. 

`(l *Logger) Debug(event string, attrs ...Attrs) *Logger`

`(l *Logger) Error(event string, attrs ...Attrs) *Logger`

`(l *Logger) Warning(event string, attrs ...Attrs) *Logger`

`(l *Logger) Info(event string, attrs ...Attrs) *Logger`

 

Changelog
---

0.0.1 - 2017-08-14 

- Initial commit. 

Author
---

Fernando Russ (`fernando.russ@mercadolibre.com)