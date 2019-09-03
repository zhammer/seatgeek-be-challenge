package main

import (
	"fmt"
	"sync"
)

const (
	FREE         = "FREE"
	RESERVED     = "RESERVED"
	SOLD         = "SOLD"
)

type Inventory struct {
	seats map[Seat]string
	lock sync.Mutex
}

func (i *Inventory) Reserve(seat Seat) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	currentStatus := i.Get(seat)
	if currentStatus != FREE {
		return fmt.Errorf("seat [%s] can only be reserved if it is [%s], it is [%s]", seat, FREE, currentStatus)
	}

	i.seats[seat] = RESERVED
	return nil
}

func (i *Inventory) Buy(seat Seat) error {
	i.lock.Lock()
	defer i.lock.Unlock()
	currentStatus := i.Get(seat)
	if currentStatus != RESERVED {
		return fmt.Errorf("seat [%s] can only be bought if it is [%s], it is [%s]", seat, RESERVED, currentStatus)
	}

	i.seats[seat] = SOLD
	return nil
}

func (i *Inventory) Get(seat Seat) string {
	status := i.seats[seat]
	if status == "" {
		return FREE
	}
	return status
}

func NewInventory() *Inventory {
	return &Inventory{
		map[Seat]string{},
		sync.Mutex{},
	}
}
