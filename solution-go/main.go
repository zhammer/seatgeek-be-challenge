package main

import (
	"bufio"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":8099")
	if err != nil {
		// handle error
	}
	for {
		println("loop")
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()
	line, err := rw.ReadString('\n')
	if err != nil && err != io.EOF {
		println(err.Error())
		os.Exit(1)
	}
	message := line
	println(1, message, 1)
	response := strings.Repeat(message, 10)
	println(2, response, 2)
	_, err = rw.WriteString(response + "\n")
	if err != nil {
		println(err.Error())
		os.Exit(2)
	}
	rw.Flush()
}
