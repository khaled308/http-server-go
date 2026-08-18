package main

import (
	"fmt"
	"net"
	"os"
)

func startServer() {
	setupRoutes()
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
		response := Response{
			StatusCode: 400,
			Body:       []byte(err.Error()),
		}

		data, _ := response.Bytes()

		conn.Write(data)
		return
	}

	if handler, ok := router.Match(request.Method, request.Path); ok {
		data, _ := handler(request).Bytes()
		conn.Write(data)
	} else {
		response := Response{
			StatusCode: 404,
			Body:       []byte("Not Found"),
		}

		data, _ := response.Bytes()

		conn.Write(data)
	}
}
