package logger

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

const (
	StatusDebug   = "DEBUG"
	StatusError   = "ERROR"
	StatusInfo    = "INFO"
	StatusWarning = "WARN"

	writeQueueChannelSize = 1000
)

var (
	defaultLogWriter io.Writer = os.Stdout
)

type Attrs map[string]interface{}

type Logger struct {
	Attributes Attrs
	Writer     io.Writer
}

var loggerChannel chan Logger

func init() {
	loggerChannel = channelWriter()
}

func Log(item Logger) {

	var logLine string

	for k, v := range item.Attributes {
		logLine += fmt.Sprintf("[%s:%+v]", k, v)
	}

	fmt.Fprintln(item.Writer, logLine)
}

func channelWriter() chan Logger {
	out := make(chan Logger, writeQueueChannelSize)
	go func() {
		for item := range out {
			Log(item)
		}
	}()
	return out
}

func LoggerWithName(c *gin.Context, name string) *Logger {
	reqID, ok := c.Get("RequestId")
	logger := &Logger{
		Attributes: map[string]interface{}{
			"source": name,
		},
		Writer: defaultLogWriter,
	}
	if ok {
		logger.Attributes["request_id"] = reqID.(string)
	}
	return logger
}

func (l *Logger) LogWithLevel(level string, event string, attrs ...Attrs) *Logger {

	item := Logger{
		Attributes: make(map[string]interface{}, 0),
		Writer:     defaultLogWriter,
	}

	// user supplied attributes
	for _, ts := range attrs {
		for k, v := range ts {
			item.Attributes[k] = v
		}
	}
	// default attributes
	for k, v := range l.Attributes {
		item.Attributes[k] = v
	}
	// base attriutes
	item.Attributes["level"] = level
	item.Attributes["event"] = event

	loggerChannel <- item

	return l
}

func (l *Logger) Debug(event string, attrs ...Attrs) *Logger {
	return l.LogWithLevel(StatusDebug, event, attrs...)
}

func (l *Logger) Error(event string, attrs ...Attrs) *Logger {
	return l.LogWithLevel(StatusError, event, attrs...)
}

func (l *Logger) Warning(event string, attrs ...Attrs) *Logger {
	return l.LogWithLevel(StatusWarning, event, attrs...)
}

func (l *Logger) Info(event string, attrs ...Attrs) *Logger {
	return l.LogWithLevel(StatusWarning, event, attrs...)
}