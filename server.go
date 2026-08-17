package main

import (
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

	_, err := readRequest(conn)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	response := Response{
		StatusCode: 200,
		Body:       []byte("Hello, World!"),
	}

	data, _ := response.Bytes()

	conn.Write(data)
}
