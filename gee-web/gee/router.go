package gee

import (
	"fmt"
	"log"
	"net/http"
)

type router struct {
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

func (r *router) addRouter(method, pattern string, handler HandlerFunc) {
	log.Printf("Add Route %s => %s", method, pattern)
	key := method + "-" + pattern

	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path

	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.Writer.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(c.Writer, "404 Not Found: %s\n", c.Path)
	}
}