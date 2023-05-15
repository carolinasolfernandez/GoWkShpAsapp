package main

import (
	"database/sql"
	"fmt"
	"time"
)

func main() {
	cb := newCircuitBreaker(breakLimit)
	concurrent := 10
	ch := make(chan string)
	for i := 0; i < concurrent; i++ {
		// time.Sleep(time.Millisecond * 1)
		go makeConcurrentRequest(cb, ch)
	}
	for i := 0; i < concurrent; i++ {
		fmt.Println(<-ch)
	}
}

// breakLimit capacidad maxima de la cola del circuit breaker
const breakLimit = 5

type sqlDB struct{}

// Query fakes query
func (db *sqlDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	time.Sleep(time.Millisecond * 10)
	return nil, nil
}

type circuitBreaker struct {
	sqlDB
	queue     chan struct{}
	numCalled int
}

// newCircuitBreaker crea un circuitbreaker con una cola limit de elementos
func newCircuitBreaker(limit int) *circuitBreaker {
	q := make(chan struct{}, limit)
	for i := 0; i < limit; i++ {
		q <- struct{}{}
	}
	s := sqlDB{}
	return &circuitBreaker{sqlDB: s, queue: q, numCalled: 0}
}

// makeConcurrentRequest simula requests concurrentes
func makeConcurrentRequest(cb *circuitBreaker, ch chan<- string) {
	start := time.Now()
	_, err := cb.Execute()
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2f segundos transcurridos. Error: %v", secs, err)
}

// Execute ejecuta queries sobre la db mientras no se exceda la capacidad de la cola del circuitbreaker
func (cb *circuitBreaker) Execute() (rows interface{}, err error) {
	cb.numCalled++
	// TODO
	r, e := cb.sqlDB.Query("SELECT FAKE FROM FAKE")
	// TODO
	return r, e
}

// hasReachedMaxCapacity utilizando channels determina si el circuitbreaker tiene capacidad para ejecutar una query
func (cb *circuitBreaker) hasReachedMaxCapacity() bool {
	// TODO
	return false
}

// doneExecuting utilizando channels, una vez terminada la ejecucion, que deberia pasar con la capacidad del circuitbreaker?
func (cb *circuitBreaker) doneExecuting() {
	// TODO
}
