package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
)

type Client interface {
	Send(message string) (string, error)
}

type TcpClient struct {
	port   string
	conn   net.Conn
	logger *Logger
}

func (c *TcpClient) connect() error {
	c.logger.Debugf("Client connecting to port [%s]", c.port)
	conn, err := net.Dial("tcp", c.port)
	if err != nil {
		c.logger.Errorf("Client found error while connecting to port [%s]: %v", c.port, err)
		return err
	}
	c.conn = conn
	return nil
}

func (c *TcpClient) disconnect() error {
	err := c.conn.Close()
	if err != nil {
		c.logger.Errorf("Client found error while disconnecting from port [%s]: %v", c.port, err)
	}
	return err
}

func (c *TcpClient) Send(message string) (string, error) {
	c.logger.Debugf("Sending message [%s] to client on port [%s]", message, c.port)
	_, err := fmt.Fprintln(c.conn, message)
	if err == io.EOF {
		c.logger.Debugf("client on port [%s] closed connection", c.port)
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("client found error while writing to socket at [%s]: %v", c.port, err)
	}

	reader := bufio.NewReader(c.conn)
	c.logger.Debugf("Reading response from port [%s]", c.port)
	response, err := reader.ReadString('\n')
	if err == io.EOF {
		c.logger.Debugf("client on port [%s] closed connection", c.port)
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("client found error while reading socket at [%s]: %v", c.port, err)
	}
	responseMsg := strings.TrimRight(response, "\n")

	c.logger.Debugf("received message [%s] from client on port [%s]", responseMsg, c.port)

	return responseMsg, err
}

func NewTcpClient(port int, logger *Logger) (Client, error) {
	host := fmt.Sprintf("localhost:%d", port)
	conn, err := net.Dial("tcp", host)

	if err != nil {
		return nil, err
	}

	return &TcpClient{
		fmt.Sprintf(":%d", port),
		conn,
		logger,
	}, nil
}
