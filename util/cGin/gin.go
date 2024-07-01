package cGin

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/ars0915/gogolook-exercise/util/cError"
	"github.com/ars0915/gogolook-exercise/util/log"
	"github.com/ars0915/gogolook-exercise/util/paging"
)

type meta struct {
	*paging.Paginator
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// func (p meta) MarshalJSON() ([]byte, error) {
// 	type Alias meta
// 	s := struct {
// 		*Alias
// 	}{
// 		Alias: (*Alias)(&p),
// 	}

// 	if p.Paginator != nil {
// 		snakeCaseJSON := jsoniter.Config{}.Froze()
// 		bytes, err := snakeCaseJSON.Marshal(&s)
// 		return bytes, err
// 	}

// 	return json.Marshal(&s)
// }

type Wrap struct {
	Meta meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Context struct {
	*gin.Context
	wrap        Wrap
	err         error
	code        int
	customError *CustomError
}

type DetailFields map[string]interface{}

var prefixCode *int

func SetResponseCodePrefix(prefix int) {
	prefixCode = &prefix
}

func NewContext(c *gin.Context) *Context {
	return &Context{
		Context: c,
		wrap:    Wrap{},
	}
}

type HandlerFunc func(*Context)

func (cGinFun HandlerFunc) GinFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := NewContext(c)
		cGinFun(ctx)
	}
}

func (c *Context) GetPaginator() paging.Paginator {
	var (
		page  = paging.DefaultPage
		limit = paging.DefaultLimit
		err   error
	)
	pageStr, pageOK := c.GetQuery(paging.PageKeyName)
	if pageOK {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page <= 0 {
			page = paging.DefaultPage
		}
	}

	limitStr, limitOK := c.GetQuery(paging.LimitKeyName)
	if limitOK {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit == 0 || limit > paging.DefaultMaxLimit {
			limit = paging.DefaultLimit
		}
	}

	offset := (page - 1) * limit
	return paging.Paginator{
		Limit:  limit,
		Page:   page,
		Offset: offset,
	}
}

// Response serializes the given struct as JSON into the response body.
func (c *Context) Response(httpCode int, msg string) {
	c.beforeResponse()
	if c.customError != nil {
		c.wrap.Meta.Code = c.customError.Code
		c.wrap.Meta.Message = c.customError.Message
		httpCode = c.customError.HTTPCode
	} else {
		if c.code == 0 {
			c.WithCode(httpCode)
		}
		c.wrap.Meta.Code = c.code
		c.wrap.Meta.Message = msg
	}

	c.logError(httpCode)

	c.JSON(httpCode, c.wrap)
	c.Abort()
}

// WithPaginator set paginator
func (c *Context) WithPaginator(page paging.Paginator) *Context {
	c.wrap.Meta.Paginator = &page
	return c
}

// WithData set response data
func (c *Context) WithData(data interface{}) *Context {
	c.wrap.Data = data
	return c
}

// WithError set error
func (c *Context) WithError(err error) *Context {
	c.err = err
	return c
}

// WithCode set code
func (c *Context) WithCode(code int) *Context {
	if code >= 1000 {
		logrus.Errorf("Invalid Error Code %d", code)
	}
	if prefixCode != nil {
		code = (*prefixCode * 1000) + code
	}
	c.code = code
	return c
}

func (c *Context) logError(httpCode int) {
	msg := c.wrap.Meta.Message
	_, file, line, _ := runtime.Caller(2)
	ginLog := logrus.WithFields(log.Fields{
		"httpSource":     fmt.Sprintf("%s:%d", file, line),
		"responseStatus": httpCode,
		"requestMethod":  c.Request.Method,
		"requestURL":     c.Request.URL.String(),
		"requestHeader":  c.Request.Header,
		"err":            c.err,
	})

	logLevel := logrus.InfoLevel
	if httpCode >= http.StatusInternalServerError {
		logLevel = logrus.ErrorLevel
	} else if httpCode >= http.StatusBadRequest {
		logLevel = logrus.WarnLevel
	}

	ginLog.Log(logLevel, msg)
}

func (c *Context) beforeResponse() *Context {
	err := c.err
	if err != nil {
		err = cError.Unwrap(err)
	}
	switch cErr := err.(type) {
	case CustomError:
		c.customError = &cErr
		c.code = cErr.Code
	}
	return c
}

// // Done invoke parent Done method
// func (c *Context) Done() <-chan struct{} {
// 	return c.Context.Done()
// }

// // Err invoke parent Err method
// func (c *Context) Err() error {
// 	return c.Context.Err()
// }

// Value invoke parent Value method
func (c *Context) Value(key interface{}) interface{} {
	if c.Request != nil {
		ctx := c.Request.Context()
		if v := ctx.Value(key); v != nil {
			return v
		}
	}
	return c.Context.Value(key)
}
