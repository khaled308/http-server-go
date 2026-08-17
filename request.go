package main

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}

var methods = map[string]bool{
	"GET":     true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"HEAD":    true,
	"OPTIONS": true,
	"PATCH":   true,
}

func isVersion(version string) bool {
	return version == "HTTP/1.0" || version == "HTTP/1.1"
}

func parseRequest(request []byte) (Request, error) {
	headerEnd := bytes.Index(request, []byte("\r\n\r\n"))
	if headerEnd == -1 {
		return Request{}, fmt.Errorf("incomplete request headers")
	}

	headerSection := string(request[:headerEnd])
	body := request[headerEnd+4:]

	lines := strings.Split(headerSection, "\r\n")

	// Request line
	parts := strings.Split(lines[0], " ")

	if len(parts) != 3 {
		return Request{}, fmt.Errorf("invalid request line")
	}

	method := parts[0]
	path := parts[1]
	version := parts[2]

	if !methods[method] {
		return Request{}, fmt.Errorf("unsupported method: %s", method)
	}

	if !isVersion(version) {
		return Request{}, fmt.Errorf("unsupported HTTP version: %s", version)
	}

	if !strings.HasPrefix(path, "/") {
		return Request{}, fmt.Errorf("invalid path: %s", path)
	}

	headers := make(map[string]string)

	for _, line := range lines[1:] {
		idx := strings.IndexByte(line, ':')
		if idx == -1 {
			return Request{}, fmt.Errorf("malformed header: %q", line)
		}

		name := strings.ToLower(strings.TrimSpace(line[:idx]))
		value := strings.TrimSpace(line[idx+1:])

		if name == "" {
			return Request{}, fmt.Errorf("empty header name")
		}

		headers[name] = value
	}

	if value, ok := headers["content-length"]; ok {
		contentLength, err := strconv.Atoi(value)
		if err != nil || contentLength < 0 {
			return Request{}, fmt.Errorf("invalid Content-Length")
		}

		if len(body) != contentLength {
			return Request{}, fmt.Errorf(
				"body length mismatch: expected %d, got %d",
				contentLength,
				len(body),
			)
		}
	}

	return Request{
		Method:  method,
		Path:    path,
		Version: version,
		Headers: headers,
		Body:    body,
	}, nil
}

func readRequest(conn net.Conn) (Request, error) {
	var rawRequest []byte
	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return Request{}, err
		}

		rawRequest = append(rawRequest, buffer[:n]...)

		headerEnd := bytes.Index(rawRequest, []byte("\r\n\r\n"))
		if headerEnd == -1 {
			continue
		}

		request, err := parseRequest(rawRequest)

		if err == nil {
			return request, nil
		}

		if len(request.Body) == 0 {
			return request, err
		}
	}
}
