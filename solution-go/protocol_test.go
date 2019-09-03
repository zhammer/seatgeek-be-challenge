package main

import "testing"

func TestParseMessage(t *testing.T) {
	t.Run("Parses valid messages", func(t *testing.T) {
		expectations := map[string][]string{
			"BUY: B0":         {BUY, "B0"},
			"RESERVE: A2342A": {RESERVE, "A2342A"},
			"QUERY: 987423d":  {QUERY, "987423d"},
		}

		for message, expectedOutput := range expectations {
			command, seat, err := ParseMessage(message)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if command != Command(expectedOutput[0]) || seat != Seat(expectedOutput[1]) {
				t.Errorf("Expected message [%s] to parse into %v, got [%s][%s]", message, expectedOutput, command, seat)
			}
		}
	})

	t.Run("Rejects invalid messages", func(t *testing.T) {
		invalidMessages := []string{
			"BUY:B0",
			"BUY: B0:",
			"B:UY: B0:",
			"RESERVE: ",
			": 987423d",
			"987423d",
			"APRICOT: 987423d",
		}

		for _, invalidMessage := range invalidMessages {
			command, seat, err := ParseMessage(invalidMessage)
			if err == nil {
				t.Errorf("No error for invalid message [%v], got: [%s][%s]", invalidMessage, command, seat)
			}
		}
	})
}
