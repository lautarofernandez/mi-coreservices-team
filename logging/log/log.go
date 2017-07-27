package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/logging/metrics"
	newrelic "github.com/newrelic/go-agent"
)

const (
	TRACE   = -1
	DEBUG   = 0
	INFO    = 1
	METRICS = 1
	WARN    = 2
	ERROR   = 4
	CRITIC  = 4
	NONE    = 100
)

var levelNames = map[string]int{
	"TRACE":  TRACE,
	"DEBUG":  DEBUG,
	"INFO":   INFO,
	"WARN":   WARN,
	"ERROR":  ERROR,
	"CRITIC": CRITIC,
	"NONE":   NONE}

var Level = NONE
var pushMetrics = false

func SetLevel(l int) {
	Level = l
}

func SetLevelByName(name string) {
	level, ok := levelNames[name]
	if !ok {
		panic(fmt.Sprintf("Invalid log level: %s", name))
	}
	SetLevel(level)
}

func SetLevelFromEnv() bool {
	level := os.Getenv("LOG_LEVEL")
	if level != "" {
		SetLevelByName(strings.ToUpper(level))
		return true
	}
	return false
}

func (context LogContext) Error(value interface{}, eventsAndTags ...interface{}) error {
	err := fmt.Errorf("%v", value)
	if Level <= ERROR {
		context.Log("error", fmt.Sprintf("%s", err), eventsAndTags...)
	}
	return err
}

func (context LogContext) Critic(value interface{}, eventsAndTags ...interface{}) error {
	err := fmt.Errorf("%v", value)
	if Level <= CRITIC {
		context.Log("critic", fmt.Sprintf("%s", err), eventsAndTags...)
	}
	return err
}

func (context LogContext) Errorf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	if Level <= ERROR {
		context.Log("error", fmt.Sprintf("%s", err))
	}
	return err
}

func (context LogContext) Info(value interface{}, eventsAndTags ...interface{}) {
	if Level > INFO {
		return
	}
	context.Log("info", fmt.Sprintf("%v", value), eventsAndTags...)
}

func (context LogContext) Debug(value interface{}, eventsAndTags ...interface{}) {
	if Level > DEBUG {
		return
	}
	context.Log("debug", fmt.Sprintf("%v", value), eventsAndTags...)
}

func (context LogContext) Trace(value interface{}, eventsAndTags ...interface{}) {
	if Level > TRACE {
		return
	}
	context.Log("trace", fmt.Sprintf("%v", value), eventsAndTags...)
}

func (context LogContext) Metric(value interface{}, eventsAndTags ...interface{}) {
	if Level > METRICS {
		return
	}
	context.Log("metric", fmt.Sprintf("%v", value), eventsAndTags...)
}

func (context LogContext) Transaction(name string) LogContext {
	if pushMetrics {
		return LogContext{tags: context.tags, transaction: metrics.Trx(name)}
	}
	return context
}

func (context LogContext) WithTransaction(trxName string, ginContext *gin.Context) LogContext {
	if pushMetrics {
		var transactionTrx newrelic.Transaction
		transaction, castBool := ginContext.Get("NR_TXN")
		if castBool {
			transactionTrx = transaction.(newrelic.Transaction)
			transactionTrx.SetName(trxName)
		}
		return LogContext{tags: context.tags, transaction: metrics.TrxWithTransaction(transactionTrx)}
	}
	return context
}

func (context LogContext) StartSegment(name string) *metrics.Segment {
	context.Metric(fmt.Sprintf("Segment \"%s\" started", name))
	if context.transaction != nil {
		return context.transaction.Segment(name)
	}
	return metrics.NullSegment()
}

func (context LogContext) EndTransaction() {
	if context.transaction != nil {
		context.transaction.End()
	}
}

func (context LogContext) Log(level string, message string, eventsAndTags ...interface{}) {
	var tags = Tags{}
	var metric metrics.Metrics // TODO: merge multiple metrics
	if len(eventsAndTags) > 0 {
		for _, eventOrTag := range eventsAndTags {
			if event, ok := eventOrTag.(string); ok {
				tags = tags.merge(Tags{"event": event})
			} else if extraTags, ok := eventOrTag.(Tags); ok {
				tags = tags.merge(extraTags)
			} else {
				if m, ok := eventOrTag.(metrics.Metrics); ok {
					metric = m
					for _, value := range m.Values {
						tags = tags.merge(Tags{value.Name: value.Value})
					}
				} else {
					panic(fmt.Sprintf("Argument must be of type Tags, Metrics or string: %v", eventOrTag))
				}
			}
		}
	}
	Log(context.tags.merge(Tags{"level": level, "message": message}).merge(tags))
	if pushMetrics {
		for _, m := range metric.Values {
			if err := metrics.PushMetric(m, context.transaction, tags.asMetricTags()...); err != nil {
				context.Errorf("Error pushing metric: %s", err)
			}
		}
	}
}

type Tags map[string]interface{}

func Log(attrs Tags) {
	var line string
	for k, v := range attrs {
		line += fmt.Sprintf(`[%s:"%+v"]`, k, v)
	}
	fmt.Println(line)
}

func (tags Tags) merge(other Tags) Tags {
	merged := make(Tags, len(tags)+len(other))
	for k, v := range tags {
		merged[k] = v
	}
	for k, v := range other {
		merged[k] = v
	}
	return merged
}

type LogContext struct {
	transaction *metrics.Transaction
	tags        Tags
}

var defaultContext = LogContext{tags: Tags{}}

func Error(value interface{}, eventsAndTags ...interface{}) error {
	return defaultContext.Error(value, eventsAndTags...)
}

func Errorf(format string, a ...interface{}) error {
	return defaultContext.Errorf(format, a)
}

func Info(value interface{}, eventsAndTags ...interface{}) {
	defaultContext.Info(value, eventsAndTags...)
}

func Debug(value interface{}, eventsAndTags ...interface{}) {
	defaultContext.Debug(value, eventsAndTags...)
}

func Trace(value interface{}, eventsAndTags ...interface{}) {
	defaultContext.Trace(value, eventsAndTags...)
}

func Critic(value interface{}, eventsAndTags ...interface{}) {
	defaultContext.Critic(value, eventsAndTags...)
}

func Metric(value interface{}, eventsAndTags ...interface{}) {
	defaultContext.Metric(value, eventsAndTags...)
}

func Transaction(name string) LogContext {
	return defaultContext.Transaction(name)
}

func WithContext(tags Tags) LogContext {
	return LogContext{tags: tags}
}

func WithEventContext(eventName string) LogContext {
	return LogContext{tags: Tags{"event": eventName}}
}

func (context LogContext) WithContext(tags Tags) LogContext {
	return LogContext{tags: context.tags.merge(tags)}
}

func PushMetrics(prefix string) {
	pushMetrics = true
	metrics.UsePrefix(prefix)
}

func init() {
	SetLevelFromEnv()
}

func (tags Tags) asMetricTags() []string {
	res := make([]string, 0, len(tags))
	for k, v := range tags {
		res = append(res, fmt.Sprintf("%v:%v", k, v))
	}
	return res
}
