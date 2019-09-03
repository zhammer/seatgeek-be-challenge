package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
)

type Strategy interface {
	Execute(name string, c Client) (bool, error)
}

type ActualStrategy struct {
	leftToProcess map[string]bool
	lock          sync.Mutex
	verb          Verb
	logger        *Logger
}

func (s *ActualStrategy) Execute(name string, c Client) (bool, error) {
	everythingWasOk := true
	s.lock.Lock()
	defer s.lock.Unlock()

	seats := sortSet(s.leftToProcess)

	for _, seat := range seats {
		command := Command{s.verb, []string{seat}}
		message := command.Serialize()

		response, err := c.Send(message)

		if err != nil {
			s.logger.Errorf("[%s] RECEIVED ERROR SENDING MESSAGE [%v], error: %v", name, message, err)
			return false, err
		}

		status := Status(response)
		if status != OK {
			s.logger.Infof("[%s] FAILED SENDING MESSAGE [%s], response: %s", name, message, status)
		} else {
			s.logger.Debugf("[%s] apple: %v", name, s.leftToProcess)
			delete(s.leftToProcess, seat)
			s.logger.Infof("[%s] SUCCEEDED SENDING MESSAGE [%s], response: %s", name, message, status)
			s.logger.Debugf("[%s] has left to process: %v", name, s.leftToProcess)
		}
		everythingWasOk = everythingWasOk && (status == OK)
	}

	return everythingWasOk, nil
}

type Consumer struct {
	name     string
	client   Client
	strategy Strategy
	logger   *Logger
}

func (c *Consumer) Tick() (bool, error) {
	return c.strategy.Execute(c.name, c.client)
}

func sortSet(leftToProcess map[string]bool) []string {
	//Necessary to make this deterministic
	var seats []string
	for seat, _ := range leftToProcess {
		seats = append(seats, seat)
	}
	sort.Strings(seats)
	return seats
}

func newRepeatUntilAllOk(initialSeats []string, verb Verb, l *Logger) Strategy {
	leftToProcess := map[string]bool{}
	for _, seat := range initialSeats {
		leftToProcess[seat] = true
	}

	return &ActualStrategy{
		leftToProcess: leftToProcess,
		verb:          verb,
		logger:        l,
	}
}

func NewManyRepeaters(minNumRepeaters int, initialSeats []string, verb Verb, logger *Logger) []Strategy {
	minSeatsPerRepeater := len(initialSeats) / minNumRepeaters

	workingList := make([]string, len(initialSeats))
	copy(workingList, initialSeats)

	var repeaters []Strategy
	for i := 0; i < minNumRepeaters; i++ {
		howManySeatsToAllocate := minSeatsPerRepeater
		if minSeatsPerRepeater > len(workingList) {
			howManySeatsToAllocate = len(workingList)
		}

		allocatedSeats := workingList[0:howManySeatsToAllocate]
		workingList = workingList[len(allocatedSeats):]
		repeater := newRepeatUntilAllOk(allocatedSeats, verb, logger)
		repeaters = append(repeaters, repeater)
		logger.Debugf("Repeater %s#%d allocated to seats: %v", verb, i, allocatedSeats)
	}
	if len(workingList) != 0 {
		extraRepeater := newRepeatUntilAllOk(workingList, verb, logger)
		repeaters = append(repeaters, extraRepeater)
	}

	logger.Infof("Created [%d] repeaters for verb [%s]", len(repeaters), verb)

	return repeaters
}

type BrokenStrategy struct {
	exit func(error)
}

var possibleBrokenMessages = []string{
	"🍏", "🍎", " 🍐", "🍊 🍋", "🍌", "🍉", "🍇", "🍓", "🍈", "🍒", "🍑", "🍍", "🥭", "🥥", "🥝", "🍅", "🍆", "🥑", "🥦",
	"🥖", "🥨", "🥯", "🧀", "🥚", "🍳", "🥞 🥓", "🥩", "🍗", "🍖", "🌭", "🍔", "🍟", "🍕", "🥪", "🥙", "🌮", "🌯", "🥗",
	"BOUGHT: Z12", "QUERY Z121", "", "Z1",
}

func (b *BrokenStrategy) Execute(name string, c Client) (bool, error) {
	message := possibleBrokenMessages[rand.Intn(len(possibleBrokenMessages))]
	response, err := c.Send(message)

	if err != nil {
		err = fmt.Errorf("[%s] RECEIVED ERROR SENDING MESSAGE [%s], error: %v", name, message, err)
		b.exit(err)
		return false, err
	}

	status := Status(response)
	if status == OK {
		err = fmt.Errorf("[%s] EXPECTED FAILED BUT SUCCEEDED SENDING MESSAGE [%s], response: %s", name, message, response)
		b.exit(err)
		return false, err
	}

	return true, nil
}

func NewBrokenConsumer(exit func(error), l *Logger) Strategy {
	return &BrokenStrategy{
		exit,
	}
}

func QueryAllSeats(seatsToQuery []string, c Client, l *Logger) (map[string]Status, error) {
	queryResponse := map[string]Status{}
	name := "q"
	for _, seat := range seatsToQuery {
		message := QuerySeat(seat).Serialize()

		response, err := c.Send(message)
		status := Status(response)

		if err != nil {
			l.Errorf("[%s] RECEIVED ERROR SENDING MESSAGE [%s], error: %v", name, message, err)
			return nil, err
		}

		l.Debugf("[%s] RESPONSE FOR MESSAGE [%s] WAS [%s]", name, message, response)
		queryResponse[seat] = status
	}
	return queryResponse, nil
}

func NewConsumer(name string, client Client, strategy Strategy, logger *Logger) *Consumer {
	return &Consumer{
		name,
		client,
		strategy,
		logger,
	}
}
