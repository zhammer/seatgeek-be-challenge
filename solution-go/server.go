package main

import (
	"fmt"
	"net"
)

type Handler func(conn net.Conn, logger *Logger)

type Server struct {
	logger  *Logger
	port    int
	handler Handler
}

func (s *Server) Start() error {

	address := fmt.Sprintf(":%d", s.port)
	s.logger.Infof("starting server at [%s]", address)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		s.logger.Errorf("error while opening socket: %v", err)
		return err
	}
	for {
		s.logger.Infof("ready to accept connections")
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Errorf("error accepting connection: %v", err)
			return err
		}
		s.logger.Debugf("conn Accepted")
		go s.handler(conn, s.logger)
	}
}

func NewServer(port int, handler Handler, logger *Logger) *Server {
	return &Server{logger, port, handler}
}
