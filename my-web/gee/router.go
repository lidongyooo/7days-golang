package gee

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	roots map[string]*node
	handlers map[string]HandlerFunc
}

func NewRouter() *router {
	return &router{
		roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	parts := strings.Split(pattern, "/")

	values := make([]string, 0)
	for _, part := range parts {
		if part != "" {
			values = append(values, part)
			if part[0] == '*' {
				break
			}
		}
	}
	return values
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	values := parsePattern(pattern)

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = NewNode()
	}
	r.roots[method].Insert(pattern, values, 0)

	key := method + "-" + pattern
	r.handlers[key] = handler
}

func (r *router) getRoute(method, pattern string) (*node, map[string]string) {
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	searchValues := parsePattern(pattern)
	n := root.Search(searchValues, 0)
	params := make(map[string]string)
	if n != nil {
		values := parsePattern(n.pattern)
		for index, value := range values{
			if value[0] == ':' {
				params[value[1:]] = searchValues[index]
			} else if value[0] == '*' {
				params[value[1:]] = strings.Join(searchValues[index:], "/")
			}
		}

		return n, params
	}

	return nil, nil
}

func (r *router) handler(c *Context)  {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil{
		c.Params = params
		key := c.Method + "-" + n.pattern

		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.Status(http.StatusNotFound)
			fmt.Fprintf(c.Writer, "Not Found.\n")
		})
	}

	c.Next()
}