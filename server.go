package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
)

func startServer() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer listener.Close()

	fmt.Println("Listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	request, err := readRequest(conn)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Printf("%+v\n", request)
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
