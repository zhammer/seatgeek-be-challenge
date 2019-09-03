package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

type Command string
type Seat string

func main() {
	logger := NewLogger(true)
	inventory := NewInventory()
	handler := newHandler(inventory)

	server := NewServer(8099, handler, logger)
	err := server.Start()
	if err != nil {
		os.Exit(1)
	}
}

func newHandler(inventory *Inventory) Handler {
	return func(conn net.Conn, logger *Logger) {
		defer func() {
			logger.Infof("Closing connection")
			conn.Close()
		}()

		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		for true {

			payload, err := rw.ReadString('\n')
			if err == io.EOF {
				return
			}
			if err != nil {
				logger.Errorf("%v", err)
				os.Exit(1)
			}
			line := strings.TrimSpace(payload)
			logger.Infof("Received message [%s]\n", line)

			var errorExecutingCommand error
			responseFromCommand := OK

			command, seat, err := ParseMessage(line)
			if err != nil {
				logger.Errorf("%v", err)
				errorExecutingCommand = err
			} else {
				logger.Infof("Executing command [%s] to seat[%s]", command, seat)
				switch command {
				case RESERVE:
					errorExecutingCommand = inventory.Reserve(seat)
				case BUY:
					errorExecutingCommand = inventory.Buy(seat)
				case QUERY:
					responseFromCommand = inventory.Get(seat)
				default:
					errorExecutingCommand = fmt.Errorf("unknown command [%s] in message [%s]", command, line)
				}
			}

			response := ""
			if errorExecutingCommand != nil {
				logger.Errorf("Error executing command: %v", errorExecutingCommand)
				response = FAIL
			} else {
				response = responseFromCommand
			}

			logger.Infof("Sending response [%s]", response)
			_, err = rw.WriteString(fmt.Sprintf("%s\n", response))
			if err != nil {
				logger.Errorf("%v", err)
				os.Exit(2)
			}
			rw.Flush()
		}

	}
}

func logPrefix() string {
	//from https://blog.sgmansfield.com/2015/12/goroutine-ids/
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return fmt.Sprintf("%04d", n)
}
