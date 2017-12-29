package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

//StartSegment start a segment for new relic
func StartSegment(c *gin.Context, name string) newrelic.Segment {
	return newrelic.StartSegment(nrgin.Transaction(c), name)
}

//StartExternalSegment start a external segment for new relic
func StartExternalSegment(c *gin.Context, url string) newrelic.ExternalSegment {
	nrSearchSegment := newrelic.ExternalSegment{URL: url}
	nrSearchSegment.StartTime = newrelic.StartSegmentNow(nrgin.Transaction(c))
	return nrSearchSegment
}
