package main

import (
	"reflect"
	"testing"
)

func TestSerializeCommands(t *testing.T) {
	t.Run("Commands serialize as expected", func(t *testing.T) {
		expectations := map[string]Command{
			"RESERVE: A1,B11,C111": AllocateSeats("A1", "B11", "C111"),
			"BUY: A1,B11,C111":     BuySeats("A1", "B11", "C111"),
			"BUY: ":                BuySeats(),
			"RESERVE: ":            AllocateSeats(""),
			"QUERY: A1":            QuerySeat("A1"),
			"QUERY: C3421231":      QuerySeat("C3421231"),
		}

		for expectedString, command := range expectations {
			actualString := command.Serialize()
			if actualString != expectedString {
				t.Errorf("Expected command [%v] to serialize as [%s], got [%s]", command, expectedString, actualString)
			}
		}
	})
}

func TestParseCommands(t *testing.T) {
	t.Run("Valid commands are parsed as expected", func(t *testing.T) {
		expectations := map[string]Command{
			"RESERVE: A1,B11,C111": AllocateSeats("A1", "B11", "C111"),
			"BUY: A1,B11,C111":     BuySeats("A1", "B11", "C111"),
		}

		for actualString, expectedCommand := range expectations {
			actualCommand, err := ParseCommand(actualString)

			if err != nil {
				t.Errorf("Expected no error when parsing string [%s], got [%v]", actualString, err)
			}

			if !reflect.DeepEqual(actualCommand, actualCommand) {
				t.Errorf("Expected string [%s] to serialize as command [%+v], got [%+v]", actualString, expectedCommand, actualCommand)
			}
		}
	})
	t.Run("Invalid commands aren't parsed", func(t *testing.T) {
		invalidCommands := []string{
			"",
			"RESERVE",
			"BUY",
			"RELEASE",
			"OTHER",
			"RESERVE:",
			"BUY:",
			"RELEASE:",
			"OTHER:",
			"RESERVE: ,",
			"BUY: ,",
			"RELEASE: ,",
			"BUY: ,",
			"RELEASE:A1",
			"RELEASE: ,",
			": ,",
			": A1,A2",
			": A1,A1",
			"BUY: ",
			"RESERVE: ",
			"RELEASE: ,,",
			"OK",
			"FAIL",
		}

		for _, actualString := range invalidCommands {
			actualCommand, err := ParseCommand(actualString)
			if err == nil {
				t.Errorf("Expected parsing string [%s] to raise error, got command [%v]", actualString, actualCommand)
			}
		}
	})
}

func TestParseResponse(t *testing.T) {
	t.Run("parse valid responses as expected", func(t *testing.T) {
		expectations := map[string]Status{"OK": OK, "FAIL": FAIL, "FREE": FREE, "SOLD": SOLD}

		for response, expectedStatus := range expectations {
			actualStatus, err := ParseResponse(response)
			if err != nil {
				t.Errorf("Expected no error when parsing response [%s], got [%v]", expectedStatus, err)
			}

			if actualStatus != expectedStatus {
				t.Errorf("Expected parsing of response [%s] to yield [%s], but got [%s]", expectedStatus, expectedStatus, actualStatus)
			}
		}
	})

	t.Run("parse invalid responses as expected", func(t *testing.T) {
		expectations := []string{
			"",
			"banana",
			"GONE",
			"RELEASE",
			"FREE:",
			"BUY: A1,B11,C111",
		}

		for _, actualString := range expectations {
			actualBool, err := ParseResponse(actualString)
			if err == nil {
				t.Errorf("Expected error when parsing response [%s], got [%v]", actualString, actualBool)
			}
		}
	})
}
