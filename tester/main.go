package main

import (
	"flag"
	"math/rand"
	"os"
)

func main() {
	consumerPort := flag.Int("consumer-port", 8099, "The port your server exposes to the CONSUMER")
	numSeats := flag.Int("seats", 500000, "A positive value indicating how many concurrent clients to use")
	concurrencyLevel := flag.Int("concurrency", 150, "A positive value indicating how many concurrent clients to use")
	randomSeed := flag.Int64("seed", 42, "A positive value used to seed the random number generator")
	debugMode := flag.Bool("debug", false, "Prints some extra information and opens a HTTP server on port 8081")
	unluckiness := flag.Int("unluckiness", 5, "A % showing the probability of something bad happenning, like broken messages being sent or random disconnects")

	flag.Parse()
	rand.Seed(*randomSeed)

	logger := NewLogger(*debugMode)
	test := NewTester(*consumerPort, *numSeats, *concurrencyLevel, *unluckiness, logger)

	test.Start()
	test.Run()
	test.Finish()
	os.Exit(0)
}
