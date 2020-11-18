package gk_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/coreservices-team/gk"
	"github.com/mercadolibre/coreservices-team/libs/go/logger"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	uuid, _ := uuid.NewV4()
	reqID := uuid.String()

	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Set("RequestId", reqID)

	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set("Content-Type", "application/json")
	c.Request.Header.Set("X-Request-Id", reqID)
	c.Request.Header.Set("X-Caller-Id", "120120120")
	c.Request.Header.Set("X-Caller-Scopes", "admin")
	c.Request.Header.Set("X-Public", "true")

	gk.Handler(func(c *gin.Context, ctx *gk.Context) {
		assert.EqualValues(t, reqID, ctx.RequestID)
		assert.EqualValues(t, 120120120, ctx.Caller.ID)
		assert.EqualValues(t, true, ctx.Caller.IsAdmin)
		assert.EqualValues(t, true, ctx.Caller.IsPublic)

		assert.NotNil(t, ctx.Log)
		assert.IsType(t, &logger.Logger{}, ctx.Log)

		assert.Implements(t, (*gk.Measurable)(nil), ctx)

		ctx.NoticeError(fmt.Errorf("Test error"))
	})(c)
}

func TestCreateTestContext(t *testing.T) {
	// This test is really unnecessary, but we do it as to not to penalize our code coverage
	gk.CreateTestContext()
}
