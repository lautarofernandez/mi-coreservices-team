package metrics

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/logging/format"
	"github.com/mercadolibre/go-meli-toolkit/gingonic/mlhandlers"
	"github.com/mercadolibre/go-meli-toolkit/godog"
	"github.com/newrelic/go-agent/v3/newrelic"
)

type Metric struct {
	metricType string
	Name       string
	Value      float64
}

type Metrics struct {
	Values []Metric
}

type Transaction struct {
	nrTrx *newrelic.Transaction
}

var NewRelicApp *newrelic.Application

const (
	FULL     = "F"
	SIMPLE   = "S"
	COMPOUND = "C"
	ERROR    = "E" // Sends error to NewRelic
)

var namePrefix = ""

func UsePrefix(prefix string) {
	namePrefix = prefix
}

// Returns a metric of type "full"
func (metrics Metrics) Full(name string, value float64) Metrics {
	return Metrics{append(metrics.Values, Metric{FULL, name, value})}
}

// Returns a metric of type "simple"
func (metrics Metrics) Simple(name string, value float64) Metrics {
	return Metrics{append(metrics.Values, Metric{SIMPLE, name, value})}
}

// Returns a metric of type "compound"
func (metrics Metrics) Compound(name string, value float64) Metrics {
	return Metrics{append(metrics.Values, Metric{COMPOUND, name, value})}
}

// Returns a metric of type "simple" with a value of 1
func (metrics Metrics) Counter(name string) Metrics {
	return Metrics{append(metrics.Values, Metric{SIMPLE, name, float64(1)})}
}

// Returns a metric of type "simple" with a value of 1
func (metrics Metrics) Error(name string) Metrics {
	return Metrics{append(metrics.Values, Metric{ERROR, name, float64(1)})}
}

// Returns a metric of type "full"
func Full(name string, value float64) Metrics {
	return Metrics{[]Metric{{FULL, name, value}}}
}

// Returns a metric of type "simple"
func Simple(name string, value float64) Metrics {
	return Metrics{[]Metric{{SIMPLE, name, value}}}
}

// Returns a metric of type "error"
func Error(name string) Metrics {
	return Metrics{[]Metric{{ERROR, name, float64(1)}}}
}

// Returns a metric of type "compund"

// Returns a metric of type "simple" with a value of 1
func Counter(name string) Metrics {
	return Metrics{[]Metric{{SIMPLE, name, float64(1)}}}
}

// Pushes a metric
func PushMetric(metric Metric, trx *Transaction, tags ...string) error {
	name := namePrefix + "." + metric.Name
	switch metric.metricType {
	case FULL:
		godog.RecordFullMetric(name, metric.Value, tags...)
	case SIMPLE:
		godog.RecordSimpleMetric(name, metric.Value, tags...)
	case COMPOUND:
		godog.RecordCompoundMetric(name, metric.Value, tags...)
	case ERROR:
		if trx != nil {
			fmt.Println("Sending error")
			trx.NoticeError(name)
		}
		godog.RecordSimpleMetric(name, float64(1), tags...)
	default:
		return fmt.Errorf("Unkown metric type: %s", metric.metricType)
	}
	return nil
}

func GingonicHandlers() []gin.HandlerFunc {
	return []gin.HandlerFunc{mlhandlers.Datadog(), NewRelic()}
}

func InitNewRelic(debug bool, environment string, appName string, appKey string) error {
	fmt.Println(environment)

	var configName = newrelic.ConfigAppName(fmt.Sprintf("%s.%s", environment, appName))
	var configLicense = newrelic.ConfigLicense(appKey)

	if app, err := newrelic.NewApplication(
		configName,
		configLicense,
		func(config *newrelic.Config) {
			if debug {
				newrelic.ConfigDebugLogger(os.Stdout)
			}
		}); err != nil {
		return fmt.Errorf("Could not create newrelic agent: %s", err)
	} else {
		NewRelicApp = app
	}

	return nil
}

// Helpers

func MinutesSince(t time.Time) float64 {
	return t.Sub(time.Now()).Minutes()
}

func ElapsedMilliseconds(t time.Time) float64 {
	return format.Milliseconds(time.Since(t))
}

func Trx(id string) *Transaction {
	nrTrx := NewRelicApp.StartTransaction(id)
	return &Transaction{nrTrx}
}

func TrxWithTransaction(nrTrx *newrelic.Transaction) *Transaction {
	return &Transaction{nrTrx: nrTrx}
}

func (trx *Transaction) Segment(name string) *Segment {
	return &Segment{newrelic.StartSegment(trx.nrTrx, name)}
}

func (trx *Transaction) NoticeError(name string) {
	if trx.nrTrx != nil {
		trx.nrTrx.NoticeError(errors.New(name))
	}
}

func (trx *Transaction) End() {
	if trx.nrTrx != nil {
		trx.nrTrx.End()
	}
}

type Segment struct {
	nrSeg *newrelic.Segment
}

func NullSegment() *Segment {
	return &Segment{}
}

func (seg *Segment) End() {
	seg.nrSeg.End()
}

// Middleware to use with New Relic
func NewRelic() gin.HandlerFunc {
	return func(c *gin.Context) {
		txn := NewRelicApp.StartTransaction(c.Request.URL.String())
		txn.SetWebResponse(c.Writer)
		txn.SetWebRequestHTTP(c.Request)
		defer txn.End()
		c.Set("NR_TXN", txn)
		c.Next()
	}
}
