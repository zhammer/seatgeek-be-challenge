package main

import (
	"flag"
	"math/rand"
	"os"
)

func main() {
	consumerPort := flag.Int("consumer-port", 8099, "The port your server exposes to the CONSUMER")
	inventoryPort := flag.Int("inventory-port", 8099, "The port your server exposes to the INVENTORY")
	concurrencyLevel := flag.Int("concurrency", 10, "A positive value indicating how many concurrent clients to use")
	randomSeed := flag.Int64("seed", 42, "A positive value used to seed the random number generator")
	debugMode := flag.Bool("debug", true, "Prints some extra information and opens a HTTP server on port 8081")
	unluckiness := flag.Int("unluckiness", 5, "A % showing the probability of something bad happenning, like broken messages being sent or random disconnects")

	flag.Parse()
	rand.Seed(*randomSeed)

	logger := NewLogger(*debugMode)
	test := NewTester(*consumerPort, *inventoryPort, *concurrencyLevel, *unluckiness, logger)

	test.Start()
	test.Run()
	test.Finish()
	os.Exit(0)
}

type Test struct {
	consumerClient *Client
	//inventoryClient    int
	concurrencyLevel int
	unluckiness      int
	logger           *Logger
}

func (s *Test) Start() {
	err := s.consumerClient.Connect()
	if err != nil {
		os.Exit(1)
	}
}

func (s *Test) Run() {
	send, err := s.consumerClient.Send("boo")
	if err != nil {
		os.Exit(1)
	}
	println(send)
}

func (s *Test) Finish() {
	s.consumerClient.Disconnect()
}

func NewTester(consumerPort int, inventoryPort int, concurrencyLevel int, unluckiness int, logger *Logger) *Test {
	consumerClient := NewClient("inventory", inventoryPort, logger)
	return &Test{
		consumerClient,
		//inventoryPort,
		concurrencyLevel,
		unluckiness,
		logger,
	}
}
