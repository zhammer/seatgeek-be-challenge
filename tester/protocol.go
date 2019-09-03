package main

import (
	"fmt"
	"strings"
)

type Command struct {
	verb  string
	seats []string
}

const (
	RESERVE = "RESERVE"
	BUY     = "BUY"
	RELEASE = "RELEASE"
)

func (c Command) Serialize() string {
	return fmt.Sprintf("%s: %s", c.verb, strings.Join(c.seats, ","))
}

func ReserveSeats(seats ...string) Command {
	return Command{RESERVE, seats}
}

func BuySeats(seats ...string) Command {
	return Command{BUY, seats}
}

func ReleaseSeats(seats ...string) Command {
	return Command{RELEASE, seats}
}

var string2Func = map[string]func(...string) Command{
	RELEASE: ReleaseSeats,
	BUY:     BuySeats,
	RESERVE: ReserveSeats,
}

func ParseCommand(message string) (Command, error) {
	verbAndPredicate := strings.Split(message, ": ")

	if len(verbAndPredicate) != 2 {
		return Command{}, fmt.Errorf("expected string [%s] to follow form [VERB: PREDICATE_1,PREDICATE_N]", message)
	}

	verb := strings.ReplaceAll(verbAndPredicate[0], ": ", "")
	predicate := verbAndPredicate[1]

	predicateWithoutNoise := strings.Trim(strings.ReplaceAll(predicate, ",", ""), " ")
	if len(predicateWithoutNoise) == 0 {
		return Command{}, fmt.Errorf("expected string [%s] to follow form [VERB: PREDICATE_1,PREDICATE_N]", message)
	}

	seats := strings.Split(strings.Trim(predicate, " "), ",")


	commandConstructor := string2Func[verb]
	if commandConstructor == nil {
		return Command{}, fmt.Errorf("cant find constructor for command [%s] in message [%s]", verb, message)
	}
	return commandConstructor(seats...), nil
}
