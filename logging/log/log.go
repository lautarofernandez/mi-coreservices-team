package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/coreservices-team/logging/metrics"
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

func (context logContext) Error(value interface{}, eventsAndTags ...interface{}) error {
	err := fmt.Errorf("%v", value)
	if Level <= ERROR {
		context.Log("error", fmt.Sprintf("%s", err), eventsAndTags...)
	}
	return err
}

func (context logContext) Critic(value interface{}, eventsAndTags ...interface{}) error {
	err := fmt.Errorf("%v", value)
	if Level <= CRITIC {
		context.Log("critic", fmt.Sprintf("%s", err), eventsAndTags...)
	}
	return err
}

func (context logContext) Errorf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	if Level <= ERROR {
		context.Log("error", fmt.Sprintf("%s", err))
	}
	return err
}

func (context logContext) Info(value interface{}, eventsAndTags ...interface{}) {
	if Level > INFO {
		return
	}
	context.Log("info", fmt.Sprintf("%v", value), eventsAndTags...)
}

func (context logContext) Debug(value interface{}, eventsAndTags ...interface{}) {
	if Level > DEBUG {
		return
	}
	context.Log("debug", fmt.Sprintf("%v", value), eventsAndTags...)
}

func (context logContext) Trace(value interface{}, eventsAndTags ...interface{}) {
	if Level > TRACE {
		return
	}
	context.Log("trace", fmt.Sprintf("%v", value), eventsAndTags...)
}

func (context logContext) Metric(value interface{}, eventsAndTags ...interface{}) {
	if Level > METRICS {
		return
	}
	context.Log("metric", fmt.Sprintf("%v", value), eventsAndTags...)
}

func (context logContext) Transaction(name string) logContext {
	if pushMetrics {
		return logContext{tags: context.tags, transaction: metrics.Trx(name)}
	}
	return context
}

func (context logContext) StartSegment(name string) *metrics.Segment {
	context.Metric(fmt.Sprintf("Segment \"%s\" started", name))
	if context.transaction != nil {
		return context.transaction.Segment(name)
	}
	return metrics.NullSegment()
}

func (context logContext) EndTransaction() {
	if context.transaction != nil {
		context.transaction.End()
	}
}

func (context logContext) Log(level string, message string, eventsAndTags ...interface{}) {
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

type logContext struct {
	transaction *metrics.Transaction
	tags        Tags
}

var defaultContext = logContext{tags: Tags{}}

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

func Transaction(name string) logContext {
	return defaultContext.Transaction(name)
}

func WithContext(tags Tags) logContext {
	return logContext{tags: tags}
}

func (context logContext) WithContext(tags Tags) logContext {
	return logContext{tags: context.tags.merge(tags)}
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
