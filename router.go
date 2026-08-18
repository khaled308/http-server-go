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
	path = strings.TrimRight(path, "/")

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

func (r *Router) extractPathParams(method, path string) (map[string]string, string) {
	params := make(map[string]string)

	for url, _ := range r.routes[method] {
		if !strings.Contains(url, ":") {
			continue
		}

		p1 := strings.Split(url, "/")
		p2 := strings.Split(path, "/")

		if len(p1) != len(p2) {
			continue
		}

		i := 0

		for i < len(p1) {
			if p1[i] == p2[i] {
				i++
			} else {
				if !strings.HasPrefix(p1[i], ":") {
					break
				}

				params[p1[i][1:]] = p2[i]
				i++
			}
		}

		if i >= len(p1) {
			return params, url
		}
	}

	return nil, ""
}

func (r *Router) Match(method, path string) (Handler, bool) {
	method = strings.ToUpper(method)

	path = strings.TrimRight(path, "/")
	idx := strings.IndexRune(path, '?')

	if idx != -1 {
		path = path[:idx]
	}

	if handler, ok := r.routes[method][path]; ok {
		return handler, true
	}

	if params, url := router.extractPathParams(method, path); params != nil {
		return r.routes[method][url], true
	}

	return nil, false
}
