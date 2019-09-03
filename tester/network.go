package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	name   string
	port   string
	logger *Logger
	conn   net.Conn
}

func (c *Client) Connect() error {
	c.logger.Infof("Client [%s] connecting to port [%s]", c.name, c.port)
	conn, err := net.Dial("tcp", c.port)
	if err != nil {
		c.logger.Errorf("Client [%s] found error while connecting to port [%s]: %v", c.name, c.port, err)
		return err
	}
	c.conn = conn
	return nil
}

func (c *Client) Disconnect() error {
	err := c.conn.Close()
	if err != nil {
		c.logger.Errorf("Client [%s] found error while disconnecting from port [%s]: %v", c.name, c.port, err)
	}
	return err
}

func (c *Client) Send(message string) (string, error) {
	c.logger.Debugf("Sending message [%s] to client [%s] on port [%s]", message, c.name, c.port)
	_, err := fmt.Fprintln(c.conn, message)
	if err != nil {
		return "", fmt.Errorf("client [%s] found error while writing to socket at [%s]: %v", c.name, c.port, err)
	}

	if err != nil {
		return "", fmt.Errorf("client [%s] found error while flushing socket at [%s]: %v", c.name, c.port, err)
	}

	reader := bufio.NewReader(c.conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("client [%s] found error while reading socket at [%s]: %v", c.name, c.port, err)
	}
	responseMsg := strings.TrimRight(response, "\n")

	c.logger.Debugf("received message [%s] from client [%s] on port [%s]", responseMsg, c.name, c.port)
	return responseMsg, nil
}

func NewClient(name string, port int, logger *Logger) *Client {
	return &Client{
		name:   name,
		port:   fmt.Sprintf(":%d", port),
		logger: logger,
	}
}
