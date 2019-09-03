package main

import (
	"fmt"
	"strings"
)

type Verb string
const (
	RESERVE = Verb("RESERVE")
	BUY     = Verb("BUY")
	QUERY   = Verb("QUERY")
)
type Status string
const (
	OK       = Status("OK")
	FAIL     = Status("FAIL")
	FREE     = Status("FREE")
	SOLD     = Status("SOLD")
	RESERVED = Status("RESERVED")
)

var validResponses = []Status{OK, FAIL, FREE, SOLD}

type Command struct {
	verb  Verb
	seats []string
}

func (c Command) Serialize() string {
	return fmt.Sprintf("%s: %s", c.verb, strings.Join(c.seats, ","))
}

func AllocateSeats(seats ...string) Command {
	return Command{RESERVE, seats}
}

func BuySeats(seats ...string) Command {
	return Command{BUY, seats}
}

func QuerySeat(seat string) Command {
	return Command{QUERY, []string{seat}}
}

var verbFunc = map[Verb]func(...string) Command{
	BUY:     BuySeats,
	RESERVE: AllocateSeats,
}

func ParseCommand(message string) (Command, error) {
	verbAndPredicate := strings.Split(message, ": ")

	if len(verbAndPredicate) != 2 {
		return Command{}, fmt.Errorf("expected string [%s] to follow form [VERB: PREDICATE_1,PREDICATE_N]", message)
	}

	verb := Verb(strings.ReplaceAll(verbAndPredicate[0], ": ", ""))
	predicate := verbAndPredicate[1]

	predicateWithoutNoise := strings.Trim(strings.ReplaceAll(predicate, ",", ""), " ")
	if len(predicateWithoutNoise) == 0 {
		return Command{}, fmt.Errorf("expected string [%s] to follow form [VERB: PREDICATE_1,PREDICATE_N]", message)
	}

	seats := strings.Split(strings.Trim(predicate, " "), ",")

	commandConstructor := verbFunc[verb]
	if commandConstructor == nil {
		return Command{}, fmt.Errorf("cant find constructor for command [%s] in message [%s]", verb, message)
	}
	return commandConstructor(seats...), nil
}

func ParseResponse(response string) (Status, error) {
	for _, r := range validResponses {
		if r == Status(response) {
			return r, nil
		}
	}

	return "", fmt.Errorf("unexpected response [%s], should be one of %v", response, validResponses)
}
