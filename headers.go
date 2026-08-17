package main

import "strings"

type Headers map[string]string

func (h Headers) Get(name string) (string, bool) {
	value, ok := h[strings.ToLower(name)]
	return value, ok
}

func (h Headers) Set(name, value string) {
	h[strings.ToLower(name)] = value
}
