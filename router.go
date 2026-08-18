package main

import "strings"

type Handler func(Request) Response

type Router struct {
	routes map[string]map[string]Handler
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]map[string]Handler),
	}
}

func (r *Router) Handle(method, path string, handler Handler) {
	method = strings.ToUpper(method)

	if r.routes[method] == nil {
		r.routes[method] = make(map[string]Handler)
	}

	r.routes[method][path] = handler
}

func (r *Router) GET(path string, handler Handler) {
	r.Handle("GET", path, handler)
}

func (r *Router) POST(path string, handler Handler) {
	r.Handle("POST", path, handler)
}

func (r *Router) PUT(path string, handler Handler) {
	r.Handle("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler Handler) {
	r.Handle("DELETE", path, handler)
}

func (r *Router) Match(method, path string) (Handler, bool) {
	method = strings.ToUpper(method)
	if handler, ok := r.routes[method][path]; ok {
		return handler, true
	}

	return nil, false
}
