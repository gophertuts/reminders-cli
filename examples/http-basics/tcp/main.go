package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	listener, _ := net.Listen("tcp", ":8080")
	defer listener.Close()
	for {
		conn, _ := listener.Accept()
		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()
	_ = conn.SetReadDeadline(time.Now().Add(time.Second))
	scanner := bufio.NewScanner(conn)
	var req string
	for scanner.Scan() {
		l := scanner.Text()
		req += l + "\n"
	}
	res := "HTTP/1.1 200 OK\r\nContent-Length: 0\r\n\r\n"
	_, err := conn.Write([]byte(res))
	if err != nil {
		log.Fatal(err)
	}
	if req != "" {
		fmt.Printf("THE REQUEST:\n%s", req)
	}
}
