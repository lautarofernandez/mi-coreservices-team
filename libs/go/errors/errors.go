package errors

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	"github.com/newrelic/go-agent/_integrations/nrgin/v1"
	"net/http"
)

type ErrorCode struct {
	Literal   string
	Status    int
	Alertable bool
}

type Error struct {
	Code    ErrorCode         `json:"code,omitempty"`
	Cause   string            `json:"cause,omitempty"`
	Message string            `json:"message,omitempty"`
	Values  map[string]string `json:"values,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s - %s: %s", e.Code.Literal, e.Cause, e.Message)
}

func (e *Error) FullError() error {
	var values string
	if e.Values != nil {
		content, err := json.Marshal(e.Values)
		if err == nil {
			values = fmt.Sprintf(" VALUES: %s", content)
		}
	}
	return fmt.Errorf("%s%s", e.Error(), values)
}

func (e *Error) MarshalJSON() ([]byte, error) {
	// serialize CODE as
	return json.Marshal(&struct {
		Error   string            `json:"error,omitempty"`
		Cause   string            `json:"cause,omitempty"`
		Message string            `json:"message,omitempty"`
		Values  map[string]string `json:"values,omitempty"`
	}{
		Error:   e.Code.Literal,
		Cause:   e.Cause,
		Message: e.Message,
		Values:  e.Values,
	})
}

var (
	BadRequestApiError = ErrorCode{
		Status:    http.StatusBadRequest,
		Literal:   "BadRequestApiError",
		Alertable: false,
	}

	NotFoundApiError = ErrorCode{
		Status:    http.StatusNotFound,
		Literal:   "NotFoundApiError",
		Alertable: false,
	}

	AuthorizationApiError = ErrorCode{
		Status:    http.StatusUnauthorized,
		Literal:   "AuthorizationApiError",
		Alertable: false,
	}

	InternalServerApiError = ErrorCode{
		Status:    http.StatusInternalServerError,
		Literal:   "InternalServerApiError",
		Alertable: true,
	}

	BadGatewayApiError = ErrorCode{
		Status:    http.StatusBadGateway,
		Literal:   "BadGatewayApiError",
		Alertable: true,
	}

	UnsupportedIndexPageSizeApiError = ErrorCode{
		Status:    http.StatusRequestEntityTooLarge,
		Literal:   "UnsupportedIndexPageSizeApiError",
		Alertable: true,
	}

)

func ReturnError(c *gin.Context, err *Error) {

	requestID, hasRequestID := c.Get("RequestId")
	if hasRequestID {
		c.Header("X-Request-Id", requestID.(string))
	}

	c.JSON(err.Code.Status, err)

	if err.Code.Alertable {
		transaction := nrgin.Transaction(c)
		if transaction != nil {
			transaction.AddAttribute(newrelic.AttributeResponseCode, err.Code.Status)
			if hasRequestID {
				transaction.AddAttribute("error.request_id", requestID.(string))
			}
			for k, v := range err.Values {
				transaction.AddAttribute(fmt.Sprintf("error.%s", k), v)
			}
			transaction.NoticeError(err)
		}
	}
}

func Wrap(c *gin.Context, code ErrorCode, err error) *Error {
	return nil
}
