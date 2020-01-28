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
	var body string
	for scanner.Scan() {
		l := scanner.Text()
		body += l + "\n"
	}
	res := fmt.Sprintf("HTTP/1.1 200 OK\nContent-Length: %d\n\n%s", len(body), body)
	_, err := conn.Write([]byte(res))
	if err != nil {
		log.Fatal(err)
	}
}
