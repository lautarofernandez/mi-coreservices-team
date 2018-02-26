package gk

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/libs/go/logger"
	"github.com/mercadolibre/go-meli-toolkit/mlauth"
	"github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"github.com/satori/go.uuid"
)

// Measurable is the interface of the exposed methods used for measuring
// code execution time and reporting errors.
type Measurable interface {
	StartSegment(name string) newrelic.Segment
	StartExternalSegment(url string) newrelic.ExternalSegment
	NoticeError(err error)
}

// Caller is the type that contains the information inside a request that
// represents the user that generated it.
type Caller struct {
	ID      uint64
	IsAdmin bool
}

// Context contains all the resources we use during a given request
type Context struct {
	ClientID  string
	Caller    Caller
	RequestID string
	Log       *logger.Logger

	nrTransaction newrelic.Transaction
}

// HandlerFunc defines the signature of our http handlers
type HandlerFunc func(*gin.Context, *Context)

// Handler receives a MeliHandlerFunc and allows it to be called from inside gin
// where a gin.HandlerFunc is expected.
func Handler(f HandlerFunc) gin.HandlerFunc {
	// Get caller function name so that we can rename newrelic transaction
	callerName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	return func(c *gin.Context) {
		rawCallerID := mlauth.GetCaller(c.Request)
		clientID := mlauth.GetClientId(c.Request)

		// If we can't parse callerID then it remains 0
		callerID, _ := strconv.ParseUint(rawCallerID, 10, 64)

		reqID := c.GetString("RequestId")

		context := &Context{
			Caller: Caller{
				ID:      callerID,
				IsAdmin: mlauth.IsCallerAdmin(c.Request),
			},
			ClientID:  clientID,
			RequestID: reqID,
			Log: &logger.Logger{
				Attributes: logger.Attrs{"request_id": reqID},
			},
			nrTransaction: nrgin.Transaction(c),
		}

		// Rename NewRelic transaction name to the name of the function that's being
		// wrapped by our context.
		if context.nrTransaction != nil {
			splitURL := strings.Split(callerName, "/")
			if len(splitURL) > 0 {
				context.nrTransaction.SetName(splitURL[len(splitURL)-1])
			}
		}

		f(c, context)
	}
}

// StartSegment makes it easy to instrument segments.
// After starting a segment do `defer segment.End()`
func (c *Context) StartSegment(name string) newrelic.Segment {
	return newrelic.StartSegment(c.nrTransaction, name)
}

// StartExternalSegment makes it easy to instrument segments that call external services.
func (c *Context) StartExternalSegment(url string) newrelic.ExternalSegment {
	return newrelic.ExternalSegment{
		URL:       url,
		StartTime: newrelic.StartSegmentNow(c.nrTransaction),
	}
}

// NoticeError records an error.  The first five errors per transaction are recorded.
func (c *Context) NoticeError(err error) {
	if c.nrTransaction == nil {
		return
	}

	c.nrTransaction.NoticeError(err)
}

// DatastoreSegment records a segment pertaining an operation with a datastore
func (c *Context) DatastoreSegment(db newrelic.DatastoreProduct, collection string, operation DBOperation) newrelic.DatastoreSegment {
	return newrelic.DatastoreSegment{
		StartTime:  newrelic.StartSegmentNow(c.nrTransaction),
		Product:    db,
		Collection: collection,
		Operation:  string(operation),
	}
}

// CreateTestContext returns a MPCS Context ready to use for testing purposes. The
// context is only populated with a functioning logger and a valid request id.
// If more information is required, then the user should add it in its end.
func CreateTestContext() *Context {
	reqID := uuid.NewV4().String()

	return &Context{
		RequestID: reqID,
		Log: &logger.Logger{
			Attributes: logger.Attrs{"request_id": reqID},
		},
	}
}
