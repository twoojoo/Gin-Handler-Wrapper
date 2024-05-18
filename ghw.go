package ghw

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiErrBodyParser func(err APIError) gin.H

type APIErrorOption struct {
	kind  string
	value any
}

func OptionPreventAbort() APIErrorOption {
	return APIErrorOption{kind: "preventAbort", value: true}
}

func OptionBodyParser(p func(err ApiErrBodyParser) gin.H) APIErrorOption {
	return APIErrorOption{kind: "bodyParser", value: p}
}

type APIError struct {
	StatusCode   int    `json:"statusCode"`
	Message      string `json:"message"`
	preventAbort bool
	bodyParser   ApiErrBodyParser
}

func (e APIError) Error() string {
	return e.Message
}

func (e APIError) PreventAbort() APIError {
	e.preventAbort = true
	return e
}

func (e APIError) WithBodyParser(p ApiErrBodyParser) APIError {
	e.bodyParser = p
	return e
}

func appendOptions(e APIError, opts []APIErrorOption) APIError {
	if opts == nil {
		return e
	}

	if len(opts) == 0 {
		return e
	}

	for _, opt := range opts {
		switch opt.kind {
		case "preventAbort":
			e.preventAbort = opt.value.(bool)
		case "bodyParser":
			e.bodyParser = opt.value.(ApiErrBodyParser)
		}
	}

	return e
}

func NewAPIError(statusCode int, message string, opts ...APIErrorOption) APIError {
	e := APIError{
		StatusCode: statusCode,
		Message:    message,
	}

	return appendOptions(e, opts)
}

func ErrBadRequest(message string, opts ...APIErrorOption) APIError {
	return NewAPIError(http.StatusBadRequest, message, opts...)
}

func ErrNotFound(message string, opts ...APIErrorOption) APIError {
	return NewAPIError(http.StatusNotFound, message, opts...)
}

func ErrInternalServerError(message string, opts ...APIErrorOption) APIError {
	return NewAPIError(http.StatusInternalServerError, message, opts...)
}

func ErrUnauthorized(message string, opts ...APIErrorOption) APIError {
	return NewAPIError(http.StatusUnauthorized, message, opts...)
}

func ErrForbidden(message string, opts ...APIErrorOption) APIError {
	return NewAPIError(http.StatusForbidden, message, opts...)
}

func ErrConflict(message string, opts ...APIErrorOption) APIError {
	return NewAPIError(http.StatusConflict, message, opts...)
}

func ErrUnprocessableEntity(message string, opts ...APIErrorOption) APIError {
	return NewAPIError(http.StatusUnprocessableEntity, message, opts...)
}

func ErrTooManyRequests(message string, opts ...APIErrorOption) APIError {
	return NewAPIError(http.StatusTooManyRequests, message, opts...)
}

func Wrap(f func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := f(c)
		if err != nil {
			if apiErr, ok := err.(APIError); ok {
				if !apiErr.preventAbort {
					c.Abort()
				}

				if apiErr.bodyParser != nil {
					c.JSON(apiErr.StatusCode, apiErr.bodyParser(apiErr))
					return
				}

				c.JSON(apiErr.StatusCode, apiErr)
				return
			}

			c.Abort()
			c.JSON(500, APIError{Message: err.Error()})
			return
		}
	}
}
