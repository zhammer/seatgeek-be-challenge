package main

import (
	"fmt"
	"net"
	"testing"
)

func respondWith(t *testing.T, server net.Listener, responseCode string) {
	for {
		conn, err := server.Accept()
		if err != nil {
			t.Fatalf("Error reading socket: %v", err)
		}
		fmt.Fprintln(conn, responseCode)
	}
}

func TestTcpClient(t *testing.T) {
	t.Run("Sends many messages to server socket", func(t *testing.T) {
		goodPort := 8081
		expectedReturn := "OK"
		goodServer, err := net.Listen("tcp", fmt.Sprintf(":%d", goodPort))

		if err != nil {
			t.Fatalf("Error opening test server: %v", err)
		}

		go respondWith(t, goodServer, expectedReturn)

		client, err := NewTcpClient(goodPort, NewLogger(false))
		if err != nil {
			t.Fatalf("Error connecting to server: %v", err)
		}

		responseCode, err := client.Send("A")

		if err != nil {
			t.Errorf("Error sending message to server: %v", err)
		}

		if responseCode != expectedReturn {
			t.Errorf("Expected responseCode to be 1, got %v", responseCode)
		}

	})
}
