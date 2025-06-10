package executor

import "sync"

// Pool limits the number of concurrent goroutines
type Pool struct {
	sem chan struct{}
	wg  sync.WaitGroup
}

// NewPool creates a pool with size limit
func NewPool(max int) *Pool {
	return &Pool{sem: make(chan struct{}, max)}
}

// Go runs the function in a goroutine while respecting concurrency limit
func (p *Pool) Go(fn func()) {
	p.wg.Add(1)
	p.sem <- struct{}{}
	go func() {
		defer func() {
			<-p.sem
			p.wg.Done()
		}()
		fn()
	}()
}

// Wait waits for all goroutines to finish
func (p *Pool) Wait() { p.wg.Wait() }
