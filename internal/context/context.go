package context

import (
	"fmt"

	"github.com/albertojnk/stonks/internal/common"
	"github.com/albertojnk/stonks/internal/core/domains"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
)

const (
	environmentdev  = "dev"
	environmentprod = "prod"
)

const ()

//Context .
type Context struct {
	LoggedAuth domains.Auth
	CID        string `json:"cid"`
	HTTPPrefix string
	Logger     *log.Entry
	Version    string
	Localizer  *i18n.Localizer
	Lang       LangType
	HostURL    string `json:"host_url"`
}

type LangType string

const (
	LangUS = "en-US"
	LangBR = "pt-BR"
)

// Result is a struct of result.
type Result struct {
	State            ResultState
	Error            error
	ErrorCode        int
	ErrorDescription string
}

// ResultState represent actual state of result.
type ResultState int

const (
	// ResultStateSuccess state is when the result is ok.
	ResultStateSuccess ResultState = 0

	// ResultStateError state is when the result has validation error.
	ResultStateError ResultState = 1

	// ResultStateUnexpected state is when the result has unexpected error.
	ResultStateUnexpected ResultState = 2

	// ResultStateNotFound when a register was not found
	ResultStateNotFound ResultState = 3

	// ResultStateUnauthorized unauthorized
	ResultStateUnauthorized ResultState = 4
)

// ResultError base struct for application errors.
type ResultError struct {
	Code        int    `json:"error_code"`
	Description string `json:"error_description"`
}

//New create a new context instance
func New() *Context {
	return &Context{
		CID:     common.GenerateUUID(),
		Version: common.GetEnv("API_VERSION", "0.0.0"),
	}
}

//WithLogger add logger to context.
func (c *Context) WithLogger() *Context {
	contextLogger := log.WithFields(log.Fields{
		"cid": c.CID,
	})

	if common.GetEnv("environment", "dev") == environmentprod {
		log.SetReportCaller(true)
		log.SetFormatter(&log.JSONFormatter{})
	}

	c.Logger = contextLogger
	return c
}

// ResultSuccess build a common struct for success.
func (c *Context) ResultSuccess() Result {
	result := Result{
		State: ResultStateSuccess,
	}

	return result
}

// ResultError build a common struct for validation error result.
func (c *Context) ResultError(errorCode int, description string, params ...interface{}) Result {

	if params != nil && len(params) > 0 {
		description = fmt.Sprintf(description, params...)
	}

	result := Result{
		State:            ResultStateError,
		ErrorCode:        errorCode,
		ErrorDescription: description,
	}

	c.Logger.Warningf("Code: %v - Message: %v\n", errorCode, description)

	return result
}

// ResultNotFound when any record was not found
func (c *Context) ResultNotFound(errorCode int, description string, params ...interface{}) Result {

	if params != nil && len(params) > 0 {
		description = fmt.Sprintf(description, params...)
	}

	result := Result{
		State:            ResultStateNotFound,
		ErrorCode:        errorCode,
		ErrorDescription: description,
	}

	c.Logger.Warningf("Code: %v - Message: %v\n", errorCode, description)

	return result
}

// ResultUnauthorized Unauthorized
func (c *Context) ResultUnauthorized() Result {
	result := Result{
		State: ResultStateUnauthorized,
	}

	return result
}

// Unexpected build a struct for unexpected error result.
func (c *Context) Unexpected(err error, data ...interface{}) Result {
	errorReturn := 999999
	description := "We don't have idea what's going on. But probably will be fixed in the next version."

	result := Result{
		State:            ResultStateUnexpected,
		Error:            err,
		ErrorCode:        errorReturn,
		ErrorDescription: description,
	}

	c.Logger.Error(description)

	return result
}

// JSON convert Result to HTTP Response.
func (r Result) JSON(c *gin.Context, body interface{}) {
	if r.State == ResultStateSuccess {
		if body == nil {
			c.String(200, "")
		} else {
			c.JSON(200, body)
		}
	} else {
		errorReturn := ResultError{
			Description: r.ErrorDescription,
			Code:        r.ErrorCode,
		}

		if r.State == ResultStateError {
			c.JSON(400, errorReturn)
		} else if r.State == ResultStateNotFound {
			c.JSON(404, errorReturn)
		} else if r.State == ResultStateUnauthorized {
			c.String(401, "")
		} else {
			c.JSON(500, errorReturn)
		}
	}
}
