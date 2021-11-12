package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type GroupRouter struct {
	prefix string
	engine *engine
	middlewares []HandlerFunc
}

type engine struct {
	*GroupRouter
	router *router
	groups []*GroupRouter
}

func New() *engine {
	log.SetFlags(log.Lshortfile)

	engine := &engine{
		router: NewRouter(),
	}
	engine.GroupRouter = &GroupRouter{engine: engine}
	engine.groups = []*GroupRouter{engine.GroupRouter}

	return engine
}

func (group *GroupRouter) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *GroupRouter) Group(prefix string) *GroupRouter {
	engine := group.engine
	newGroup := &GroupRouter{
		engine: engine,
		prefix: group.prefix + prefix,
	}
	engine.groups = append(engine.groups, newGroup)

	return newGroup
}

func (group *GroupRouter) addRoute(method, pattern string, handler HandlerFunc) {
	group.engine.router.addRoute(method, group.prefix+pattern, handler)
}

func (group *GroupRouter) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *GroupRouter) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (e *engine) Run(addr string) {
	http.ListenAndServe(addr, e)
}

func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext(w, r)

	var middlewares []HandlerFunc
	for _, group := range e.groups {
		if strings.HasPrefix(c.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c.handlers = middlewares
	e.router.handler(c)
}
