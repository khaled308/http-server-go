package main

import (
	"bytes"
	"fmt"
	"strconv"
)

var statusText = map[int]string{
	200: "OK", 201: "Created", 204: "No Content",
	301: "Moved Permanently", 302: "Found", 304: "Not Modified",
	400: "Bad Request", 401: "Unauthorized", 403: "Forbidden",
	404: "Not Found", 405: "Method Not Allowed", 500: "Internal Server Error",
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func (response Response) Bytes() ([]byte, error) {
	if response.StatusCode == 0 {
		response.StatusCode = 200
	}

	status, ok := statusText[response.StatusCode]
	if !ok {
		return nil, fmt.Errorf("unsupported status code: %d", response.StatusCode)
	}

	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}

	if _, ok := response.Headers["Content-Type"]; !ok {
		response.Headers["Content-Type"] = "text/plain"
	}

	if _, ok := response.Headers["Content-Length"]; !ok {
		response.Headers["Content-Length"] = strconv.Itoa(len(response.Body))
	}

	if _, ok := response.Headers["Connection"]; !ok {
		response.Headers["Connection"] = "close"
	}

	var res bytes.Buffer

	fmt.Fprintf(
		&res,
		"HTTP/1.1 %d %s\r\n",
		response.StatusCode,
		status,
	)

	for key, value := range response.Headers {
		fmt.Fprintf(&res, "%s: %s\r\n", key, value)
	}

	res.WriteString("\r\n")
	res.Write(response.Body)

	return res.Bytes(), nil
}
