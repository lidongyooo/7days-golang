package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Request *http.Request
	Params map[string]string

	//request info
	Method string
	Path string

	handlers []HandlerFunc
	index int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Request: r,
		Method: r.Method,
		Params: make(map[string]string),
		Path: r.URL.Path,
		index: -1,
	}
}

func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) SetHeader(key ,value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) Status(code int)  {
	c.Writer.WriteHeader(code)
}

func (c *Context) Fail(data interface{}) {
	c.index = len(c.handlers)
	c.JSON(http.StatusInternalServerError, data)
}

func (c *Context) JSON(statusCode int, data interface{}) {
	c.SetHeader("Content-Type", "application/json; utf-8")
	c.Status(statusCode)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(data); err != nil {
		c.Status(http.StatusInternalServerError)
		fmt.Fprintf(c.Writer, "Server Error.")
	}
}

func (c *Context) String(statusCode int, data string) {
	c.SetHeader("Content-Type", "text/plain; utf-8")
	c.Status(statusCode)
	fmt.Fprintf(c.Writer, data)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) Next() {
	c.index++

	len := len(c.handlers)
	for ; c.index < len; c.index++ {
		c.handlers[c.index](c)
	}
}
