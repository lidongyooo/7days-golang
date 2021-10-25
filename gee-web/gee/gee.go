package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

func (engine *Engine) AddRoute (method, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern

	engine.router[key] = handler
}

func (engine *Engine) Get (pattern string, handler HandlerFunc) {
	engine.AddRoute("Get", pattern, handler)
}

func (engine *Engine) Post (pattern string, handler HandlerFunc)  {
	engine.AddRoute("Post", pattern, handler)
}


func (engine *Engine) Run (addr string)  {
	http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP (w http.ResponseWriter, r *http.Request)  {
	key := r.Method + "-" + r.URL.Path

	if handler, ok := engine.router[key]; ok {
		handler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Not Found: %s\n", r.URL)
	}
}
