package main

import (
	"fmt"
	"os"
	"reflect"
	"sync"
)

type Tester struct {
	consumerPort        int
	numSeats            int
	concurrency         int
	unluckiness         int
	blockingConsumers   []*Consumer
	backgroundConsumers []*Consumer
	expectedResults     map[string]Status
	logger              *Logger
}

func (t *Tester) finishTest(fail bool, format string, v ...interface{}) {
	var exitStatus int
	if fail {
		exitStatus = 1
		t.logger.InfoBannerf("❌ TEST FAILED: %v ❌", fmt.Sprintf(format, v...))
	} else {
		exitStatus = 0
		t.logger.InfoBannerf("✅ TEST SUCCESSFUL ✅")
	}

	os.Exit(exitStatus)
}

func (t *Tester) failE(err error) {
	t.failF("error: %v", err)
}

func (t *Tester) failF(format string, v ...interface{}) {
	t.finishTest(true, format, v...)
}

func (t *Tester) allKnownSeats() []string {
	allSeats := make([]string, 0, len(t.expectedResults))
	for seat, _ := range t.expectedResults {
		allSeats = append(allSeats, seat)
	}
	return allSeats
}

func (t *Tester) ensureServerHasExpectedState(expectedResults map[string]Status) {

	actualResults, err := QueryAllSeats(t.allKnownSeats(), t.newClient(), t.logger)
	if err != nil {
		t.failF("Error while querying state of all known seats: %v", err)
	}

	for _, s := range t.allKnownSeats() {
		expected := expectedResults[s]
		actual := actualResults[s]

		if expected == actual {
			t.logger.Debugf("%s - Expecting seat [%s] to have status [%s], got [%s]", "✅", s, expected, actual)
		} else {
			t.logger.Infof("%s - Expecting seat [%s] to have status [%s], got [%s]", "❌", s, expected, actual)
		}

	}

	if !reflect.DeepEqual(expectedResults, actualResults) {
		t.failF("Actual results different from expected!")
	}
}

func (t *Tester) newClient() Client {
	client, err := NewTcpClient(t.consumerPort, t.logger)
	if err != nil {
		t.failF("Error while connecting to server on port [%v]: %v", t.consumerPort, err)
	}

	return client
}

func (t *Tester) Run() {
	t.logger.Infof("Making sure server is clear")
	expectedState := map[string]Status{}
	for _, seat := range t.allKnownSeats() {
		expectedState[seat] = FREE
	}
	t.ensureServerHasExpectedState(expectedState)

	for _, c := range t.backgroundConsumers {
		t.logger.Infof("Starting consumer [%s]", c.name)
		go func() {
			for {
				c.Tick()
			}
		}()
	}

	wg := &sync.WaitGroup{}
	wg.Add(len(t.blockingConsumers))

	for _, c := range t.blockingConsumers {
		t.logger.Infof("Starting consumer [%s]", c.name)
		go func(tester *Tester, consumer *Consumer, barrier *sync.WaitGroup) {
			for {
				finished, err := consumer.Tick()

				if err != nil {
					tester.failF("[%v] found error, exiting test suite. error: %v", consumer.name, err)
				}

				if finished {
					tester.logger.Infof("[%s] finished its task", consumer.name)
					barrier.Done()
					break
				} else {
					tester.logger.Infof("[%s] has not finished its task yet", consumer.name)
				}
			}
		}(t, c, wg)
	}
	wg.Wait()
}

func (t *Tester) Finish() {
	t.logger.InfoBannerf("Finishing test")

	t.ensureServerHasExpectedState(t.expectedResults)

	t.finishTest(false, "")
}

func (t *Tester) Start() {
	t.logger.InfoBannerf("Starting test")
	numSeatsPerType := t.numSeats / 3
	numRepeatersPerType := t.concurrency / 3

	t.expectedResults = map[string]Status{}
	var seatsToRemain []string
	for i := 0; i < numSeatsPerType; i++ {
		seat := fmt.Sprintf("A%03d", i)
		t.expectedResults[seat] = FREE
		seatsToRemain = append(seatsToRemain, seat)
	}

	var seatsToBuy []string
	for i := 0; i < numSeatsPerType; i++ {
		seat := fmt.Sprintf("B%03d", i)
		t.expectedResults[seat] = SOLD
		seatsToBuy = append(seatsToBuy, seat)
	}

	var seatsToAllocateOnly []string
	for i := 0; i < numSeatsPerType; i++ {
		seat := fmt.Sprintf("C%03d", i)
		t.expectedResults[seat] = RESERVED
		seatsToAllocateOnly = append(seatsToAllocateOnly, seat)
	}

	allSeatsToAllocate := append(seatsToAllocateOnly, seatsToBuy...)

	var allocators []*Consumer
	for i, r := range NewManyRepeaters(numRepeatersPerType, allSeatsToAllocate, RESERVE, t.logger){
		name := fmt.Sprintf("allocator-%03d", i)
		allocators = append(allocators, NewConsumer(name, t.newClient(), r, t.logger))
	}

	var buyers []*Consumer
	for i, r := range NewManyRepeaters(numRepeatersPerType, seatsToBuy, BUY, t.logger){
		name := fmt.Sprintf("buyer-%03d", i)
		buyers = append(buyers, NewConsumer(name, t.newClient(), r, t.logger))
	}

	t.blockingConsumers = append(buyers, allocators...)

	brokenConsumer := NewConsumer("broken-consumer", t.newClient(), NewBrokenConsumer(t.failE, t.logger), t.logger)
	t.backgroundConsumers = []*Consumer{brokenConsumer}
}

func NewTester(consumerPort int, numSeats int, concurrencyLevel int, unluckiness int, logger *Logger) *Tester {

	return &Tester{
		consumerPort: consumerPort,
		concurrency:  concurrencyLevel,
		numSeats:     numSeats,
		logger:       logger,
	}
}
