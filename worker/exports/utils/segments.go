package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent/v3/integrations/nrgin"
	"github.com/newrelic/go-agent/v3/newrelic"
)

//StartSegment start a segment for new relic
func StartSegment(c *gin.Context, name string) *newrelic.Segment {
	return newrelic.StartSegment(nrgin.Transaction(c), name)
}

//StartExternalSegment start a external segment for new relic
func StartExternalSegment(c *gin.Context, url string) *newrelic.ExternalSegment {
	return &newrelic.ExternalSegment{
		URL:       url,
		StartTime: newrelic.StartSegmentNow(nrgin.Transaction(c)),
	}
}
