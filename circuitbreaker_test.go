package main

import (
	"testing"
)

func TestCircuitBreakerCalled(t *testing.T) {
	cb := newCircuitBreaker(breakLimit)
	_, err := cb.Execute()
	if err != nil {
		t.Log(err)
	}
	if cb.numCalled != 1 {
		t.Fatalf("Error in calling circuitbreaker, expected %d, found %d", 1, cb.numCalled)
	}
}

func TestHasReachedMaxCapacity_True(t *testing.T) {
	q := make(chan struct{}, 1)
	q <- struct{}{}
	cb := &circuitBreaker{queue: q}
	res := cb.hasReachedMaxCapacity()
	if res != true {
		t.Fatalf("got %t, want %t", res, true)
	}
}

func TestHasReachedMaxCapacity_False(t *testing.T) {
	q := make(chan struct{}, 1)
	cb := &circuitBreaker{queue: q}
	res := cb.hasReachedMaxCapacity()
	if res != false {
		t.Fatalf("got %t, want %t", res, true)
	}
}

func TestDoneExecuting(t *testing.T) {
	q := make(chan struct{}, 2)
	q <- struct{}{}
	cb := &circuitBreaker{queue: q}
	if len(cb.queue) != 1 {
		t.Fatalf("got %d, want %d", len(cb.queue), 1)
	}
	cb.doneExecuting()
	if len(cb.queue) != 2 {
		t.Fatalf("got %d, want %d", len(cb.queue), 2)
	}
}

func TestExecute_Error(t *testing.T) {
	q := make(chan struct{}, 1)
	cb := &circuitBreaker{numCalled: 0, queue: q}
	_, err := cb.Execute()
	if err == nil {
		t.Fatalf("error expected")
	}
}

func TestExecute_OK(t *testing.T) {
	q := make(chan struct{}, 1)
	cb := &circuitBreaker{numCalled: 0, queue: q}
	q <- struct{}{}
	_, err := cb.Execute()
	if err != nil {
		t.Fatalf("no error expected")
	}
}
