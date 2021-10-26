package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	router *router
}

func New() *Engine {
	log.SetFlags(log.Lshortfile)

	return &Engine{
		router: newRouter(),
	}
}

func (engine *Engine) GET (pattern string, handler HandlerFunc) {
	engine.router.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST (pattern string, handler HandlerFunc)  {
	engine.router.addRoute("POST", pattern, handler)
}


func (engine *Engine) Run (addr string)  {
	http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP (w http.ResponseWriter, r *http.Request)  {
	engine.router.handle(NewContext(w, r))
}
