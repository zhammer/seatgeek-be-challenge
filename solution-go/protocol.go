package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	RESERVE = "RESERVE"
	BUY     = "BUY"
	QUERY   = "QUERY"
	OK      = "OK"
	FAIL    = "FAIL"
)

func ParseMessage(line string) (Command, Seat, error) {
	matches, err := regexp.MatchString("^\\w+: \\w+$", line)
	if err != nil {
		panic("Could not compile regular expression")
	}

	if !matches {
		return "", "", fmt.Errorf("invalid message [%s]", line)
	}

	split := strings.Split(line, ": ")

	command := Command(strings.TrimSpace(split[0]))
	seat := Seat(strings.TrimSpace(split[1]))

	if command != RESERVE &&
		command != BUY &&
		command != QUERY {
		return "", "", fmt.Errorf("invalid command [%s] in message [%s]", command, line)
	}

	return command, seat, nil
}
