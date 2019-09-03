package main

import "testing"

func TestNewInventory(t *testing.T) {
	t.Run("Seats are free when no other action was performed", func(t *testing.T) {
		seats := []Seat{"A1", "b4", "qe33"}
		inventory := NewInventory()

		expectAllSeatsToHaveStatus(t, inventory, seats, FREE)
	})

	t.Run("Free seats can be reserved", func(t *testing.T) {
		seatsToReserve := []Seat{"A1", "b4", "A321"}
		seatsToRemainFree := []Seat{"ABC", "ORCH1", "Z2"}
		inventory := NewInventory()

		for _, seat := range seatsToReserve {
			err := inventory.Reserve(seat)
			if err != nil {
				t.Errorf("Unexpected error when reserving seat [%s]: %v", seat, err)
			}
		}

		expectAllSeatsToHaveStatus(t, inventory, seatsToRemainFree, FREE)
		expectAllSeatsToHaveStatus(t, inventory, seatsToReserve, RESERVED)
	})

	t.Run("Reserved seats can be bought", func(t *testing.T) {
		seatsToReserveOnly := []Seat{"A1", "A4", "A321"}
		seatsToRemainFree := []Seat{"ABC", "ORCH1", "Z2"}
		seatsToBuy := []Seat{"LL", "67231", "32y789dwund"}

		inventory := NewInventory()
		allSeatsToReserve := append(seatsToReserveOnly, seatsToBuy...)
		for _, seat := range allSeatsToReserve {
			err := inventory.Reserve(seat)
			if err != nil {
				t.Errorf("Unexpected error when reserving seat [%s]: %v", seat, err)
			}
		}

		for _, seat := range seatsToBuy {
			err := inventory.Buy(seat)
			if err != nil {
				t.Errorf("Unexpected error when buying seat [%s]: %v", seat, err)
			}
		}

		expectAllSeatsToHaveStatus(t, inventory, seatsToRemainFree, FREE)
		expectAllSeatsToHaveStatus(t, inventory, seatsToReserveOnly, RESERVED)
		expectAllSeatsToHaveStatus(t, inventory, seatsToBuy, SOLD)
	})

	t.Run("Seats cannot be bought unless reserved", func(t *testing.T) {
		seatsToReserveOnly := []Seat{"A1", "A4", "A321"}
		seatsToRemainFree := []Seat{"ABC", "ORCH1", "Z2"}
		seatsThatDontExist := []Seat{"XYZ", "XCC", "XCC"}

		inventory := NewInventory()

		for _, seat := range seatsToReserveOnly {
			err := inventory.Reserve(seat)
			if err != nil {
				t.Errorf("Unexpected error when reserving seat [%s]: %v", seat, err)
			}
		}

		allSeats := append(seatsToRemainFree, seatsThatDontExist...)
		for _, seat := range allSeats {
			err := inventory.Buy(seat)
			if err == nil {
				currentStatus := inventory.Get(seat)
				t.Fatalf("Expecting error when buying seat [%s], got nothing. Seart currently marked as [%s]", seat, currentStatus)
			}
		}

		expectAllSeatsToHaveStatus(t, inventory, seatsToRemainFree, FREE)
		expectAllSeatsToHaveStatus(t, inventory, seatsToReserveOnly, RESERVED)
		expectAllSeatsToHaveStatus(t, inventory, seatsThatDontExist, FREE)
	})
}

func expectAllSeatsToHaveStatus(t *testing.T, inventory *Inventory, seats []Seat, desiredStatus string) {
	for _, seat := range seats {
		seatStatus := inventory.Get(seat)
		if seatStatus != desiredStatus {
			t.Errorf("Seat [%s] expected to be [%s], got [%s]", seat, desiredStatus, seatStatus)
		}
	}
}
