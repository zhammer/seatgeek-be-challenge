package main

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
)

type MockClient struct {
	NameToReturn            string
	ResponseToReturnAlways  string
	ListOfResponsesToReturn []string
	ErrorToReturn           error
	ListOfMessageReceived   []string
}

func (n *MockClient) Name() string {
	return n.NameToReturn
}

func (n *MockClient) Connect() error {
	return n.ErrorToReturn
}

func (n *MockClient) Disconnect() error {
	return n.ErrorToReturn
}

func (n *MockClient) Send(message string) (string, error) {
	n.ListOfMessageReceived = append(n.ListOfMessageReceived, message)
	responseToReturn := ""
	if n.ResponseToReturnAlways != "" {
		responseToReturn = n.ResponseToReturnAlways
	} else {
		responseToReturn = n.ListOfResponsesToReturn[0]
		n.ListOfResponsesToReturn = n.ListOfResponsesToReturn[1:]
	}
	return responseToReturn, n.ErrorToReturn
}

var logger = NewLogger(false)

func TestNewRepeatUntilAllOk(t *testing.T) {
	t.Run("strategy doest do anything if no seats in list", func(t *testing.T) {
		mockClient := &MockClient{
			ResponseToReturnAlways: "OK",
		}

		var seats []string
		verb := RESERVE

		repeater := newRepeatUntilAllOk(seats, verb, logger)

		success, err := repeater.Execute("name", mockClient)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if !success {
			t.Fatalf("failed when no seats left")
		}
	})

	t.Run("processes all seats in one go if all return ok", func(t *testing.T) {
		mockClient := &MockClient{
			ResponseToReturnAlways: "OK",
		}

		seats := []string{"A1", "B2", "C3"}
		verb := RESERVE

		repeater := newRepeatUntilAllOk(seats, verb, logger)

		finished, err := repeater.Execute("repeater", mockClient)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !finished {
			t.Fatalf("Strategy wasn't finished. Initial seats: %v Mock Client: %+v", seats, mockClient)
		}

		expectedMessagesReceived := []string{"RESERVE: A1", "RESERVE: B2", "RESERVE: C3"}
		sort.Strings(expectedMessagesReceived)
		sort.Strings(mockClient.ListOfMessageReceived)
		if !reflect.DeepEqual(mockClient.ListOfMessageReceived, expectedMessagesReceived) {
			t.Fatalf("Expecting client to receive messages %v, got %v", expectedMessagesReceived, mockClient.ListOfMessageReceived)
		}
	})

	t.Run("tries any failed ones in until done", func(t *testing.T) {
		mockClient := &MockClient{
			ListOfResponsesToReturn: []string{
				"OK", "FAIL", "OK",
				"FAIL",
				"OK",
			},
		}

		seats := []string{"A1", "B2", "C3"}
		verb := RESERVE

		repeater := newRepeatUntilAllOk(seats, verb, logger)

		//1st try
		finished, err := repeater.Execute("repeater", mockClient)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if finished {
			t.Fatalf("Strategy was finished on first iteration. Initial seats: %v Mock Client: %+v", seats, mockClient)
		}

		//2nd try
		finished, err = repeater.Execute("repeater", mockClient)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if finished {
			t.Fatalf("Strategy was finished on first iteration. Initial seats: %v Mock Client: %+v", seats, mockClient)
		}

		//3rd is a charm
		finished, err = repeater.Execute("repeater", mockClient)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !finished {
			t.Fatalf("Strategy wasn't finished. Initial seats: %v Mock Client: %+v", seats, mockClient)
		}

		expectedMessagesReceived := []string{
			"RESERVE: A1", "RESERVE: B2", "RESERVE: C3",
			"RESERVE: B2",
			"RESERVE: B2",
		}
		if !reflect.DeepEqual(mockClient.ListOfMessageReceived, expectedMessagesReceived) {
			t.Fatalf("Expecting client to receive messages %v, got %v", expectedMessagesReceived, mockClient.ListOfMessageReceived)
		}
	})

	t.Run("tries again after errors", func(t *testing.T) {
		mockClient := &MockClient{
			ResponseToReturnAlways: "OK",
		}

		seats := []string{"A1", "B2", "C3"}
		verb := RESERVE

		repeater := newRepeatUntilAllOk(seats, verb, logger)

		//1st try
		mockClient.ErrorToReturn = errors.New("expected")
		finished, err := repeater.Execute("repeater", mockClient)
		if err == nil {
			t.Fatalf("Expecting error, got nothing. Mock Client: %+v", mockClient)
		}

		if finished {
			t.Fatalf("Strategy was finished on first iteration. Initial seats: %v Mock Client: %+v", seats, mockClient)
		}

		//2nd try
		mockClient.ErrorToReturn = nil
		mockClient.ResponseToReturnAlways = "OK"
		finished, err = repeater.Execute("repeater", mockClient)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if !finished {
			t.Fatalf("Strategy wasn't finished. Initial seats: %v Mock Client: %+v", seats, mockClient)
		}
		expectedMessagesReceived := []string{
			"RESERVE: A1",
			"RESERVE: A1", "RESERVE: B2", "RESERVE: C3",
		}
		if !reflect.DeepEqual(mockClient.ListOfMessageReceived, expectedMessagesReceived) {
			t.Fatalf("Expecting client to receive messages %v, got %v", expectedMessagesReceived, mockClient.ListOfMessageReceived)
		}
	})
}

func TestNewManyRepeaters(t *testing.T) {
	t.Run("evenly yet randomly splits a list amongst repeaters", func(t *testing.T) {
		var list []string
		for i := 0; i < 30; i++ {
			list = append(list, fmt.Sprintf("seat#%03d", i))
		}

		minNumRepeatersExpected := 3
		repeaters := NewManyRepeaters(minNumRepeatersExpected, list, BUY, logger)
		if len(repeaters) < minNumRepeatersExpected {
			t.Fatalf("Expeting max [%d] repeaters, got [%d]", minNumRepeatersExpected, len(repeaters))
		}

		var allReceivedMessages []string
		for i, repeater := range repeaters {
			client := &MockClient{ResponseToReturnAlways: "OK"}
			repeaterName := fmt.Sprintf("repeater#%d", i)
			finished, err := repeater.Execute(repeaterName, client)
			if err != nil {
				t.Fatalf("Unexpected error for repeater [%s]: %v", repeaterName, err)
			}

			if !finished {
				t.Fatalf("Expected repeater [%s] to finish, it didnt", repeaterName)
			}

			if len(client.ListOfMessageReceived) == 0 {
				t.Fatalf("Expected repeater [%s] to send messages to client, it didnt. Client: %+v", repeaterName, client)
			}

			if len(client.ListOfMessageReceived) == len(list) {
				t.Fatalf("Repeater [%s] seems to have been allocated the whole seat list", repeaterName)
			}
			allReceivedMessages = append(allReceivedMessages, client.ListOfMessageReceived...)
		}

		var allExpectedMessages []string
		for _, l := range list {
			allExpectedMessages = append(allExpectedMessages, fmt.Sprintf("BUY: %s", l))
		}

		sort.Strings(allExpectedMessages)
		sort.Strings(allReceivedMessages)
		if !reflect.DeepEqual(allReceivedMessages, allExpectedMessages) {
			t.Fatalf("Expected messages to be sent: %v, got %v", allExpectedMessages, allReceivedMessages)
		}
	})
}

func TestNewBrokenConsumer(t *testing.T) {
	t.Run("broken consumer sends invalid messages", func(t *testing.T) {
		mockClient := &MockClient{
			ResponseToReturnAlways: "FAIL",
		}

		var exitError error
		exitFn := func(e error) {
			exitError = e
		}

		broken := NewBrokenConsumer(exitFn, logger)

		success, err := broken.Execute("name", mockClient)

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if exitError != nil {
			t.Fatalf("Unexpected exit error: %v", err)
		}
		if !success {
			t.Fatalf("failed to send broken message to client [%+v]", mockClient)
		}

		validVerbs := []string{"RESERVE", "BUY", "QUERY"}

		for _, verb := range validVerbs {
			messageReceived := mockClient.ListOfMessageReceived[0]
			if strings.Contains(messageReceived, verb) {
				t.Errorf("message [%s] was supposed to be invalid, but contained valid verb [%s]", messageReceived, verb)
			}
		}
	})

	t.Run("returns error if error buying seats", func(t *testing.T) {
		mockClient := &MockClient{
			ResponseToReturnAlways: "OK",
			ErrorToReturn:          errors.New("expected error"),
		}

		var exitError error
		exitFn := func(e error) {
			exitError = e
		}

		broken := NewBrokenConsumer(exitFn, logger)
		success, err := broken.Execute("name", mockClient)

		if err == nil {
			t.Fatalf("expeting error, got nothing")
		}
		if exitError == nil {
			t.Fatalf("expecting exit error, got nothing")
		}

		if success {
			t.Fatalf("expecting failure, got success")
		}
	})
}

func TestQueryAllSeats(t *testing.T) {
	t.Run("queries all seats and reports results", func(t *testing.T) {
		expectedSeatStatuses := map[string]Status{
			"A1": OK,
			"B2": FAIL,
			"C3": OK,
		}

		var seats []string
		var responses []string
		for seat, response := range expectedSeatStatuses {
			seats = append(seats, seat)
			responses = append(responses, string(response))
		}

		mockClient := &MockClient{
			ListOfResponsesToReturn: responses,
		}

		results, err := QueryAllSeats(seats, mockClient, logger)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if !reflect.DeepEqual(results, expectedSeatStatuses) {
			t.Fatalf("Expected resutls to be %v, got %v", expectedSeatStatuses, results)
		}
	})
}

func TestEqualWhenSorted(t *testing.T) {
	t.Run("Compares slices as expected", func(t *testing.T) {
		type expectation struct {
			listA  []string
			listB  []string
			result bool
		}
		expectations := []expectation{
			{
				[]string{"a", "b", "c", ""},
				[]string{"a", "b", "c", ""},
				true,
			},
			{
				[]string{"c", "b", "a"},
				[]string{"c", "a", "b"},
				true,
			},
			{
				[]string{"a", "b", "c"},
				[]string{"a", "b"},
				false,
			},
			{
				[]string{"a", "b", "c"},
				[]string{},
				false,
			},

			{
				[]string{},
				[]string{},
				true,
			},
		}

		for _, e := range expectations {
			sort.Strings(e.listA)
			sort.Strings(e.listB)
			result := reflect.DeepEqual(e.listA, e.listB)

			if result != e.result {
				t.Fatalf("Expected comparing %v to %v to be [%t], got [%t]", e.listA, e.listB, e.result, result)
			}
		}
	})
}
