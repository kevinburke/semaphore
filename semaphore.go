package semaphore

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Semaphore is an implementation of semaphore.
type Semaphore struct {
	permits int
	avail   int
	channel chan struct{}
	aMutex  *sync.RWMutex
	rMutex  *sync.Mutex
}

// New creates a new Semaphore with specified number of concurrent workers.
func New(n int) *Semaphore {
	if n < 1 {
		panic("Invalid number of permits. Less than 1")
	}

	// fill channel buffer
	channel := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		channel <- struct{}{}
	}

	return &Semaphore{
		n,
		n,
		channel,
		&sync.RWMutex{},
		&sync.Mutex{},
	}
}

// Acquire blocks until a worker becomes available.
func (s *Semaphore) Acquire() {
	s.aMutex.Lock()
	defer s.aMutex.Unlock()

	<-s.channel
	s.avail--
}

// AcquireMany is similar to Acquire() but for many workers.
// An error is returned if n is greater number of workers in the semaphore.
func (s *Semaphore) AcquireMany(n int) error {
	if n > s.permits {
		return errors.New("Too many requested permits")
	}
	s.aMutex.Lock()
	defer s.aMutex.Unlock()

	s.avail -= n
	for ; n > 0; n-- {
		<-s.channel
	}
	s.avail += n
	return nil
}

// AcquireWithin is similar to AcquireMany() but cancels if duration elapses before getting the permits.
// Returns true if successful and false if timeout occurs.
func (s *Semaphore) AcquireContext(n int, ctx context.Context) bool {
	go func() {
		time.Sleep(d)
		timeout <- true
	}()
	go func() {
		s.AcquireMany(n)
		timeout <- false
		if <-cancel {
			s.ReleaseMany(n)
		}
	}()
	if <-timeout {
		cancel <- true
		return false
	}
	cancel <- false
	return true
}

// Release releases one worker.
func (s *Semaphore) Release() {
	s.rMutex.Lock()
	defer s.rMutex.Unlock()

	s.channel <- struct{}{}
	s.avail++
}

// ReleaseMany releases n permits.
func (s *Semaphore) ReleaseMany(n int) {
	if n > s.permits {
		panic("Too many requested releases")
	}
	for ; n > 0; n-- {
		s.Release()
	}
}

// AvailablePermits gives number of available unacquired permits.
func (s *Semaphore) AvailablePermits() int {
	s.aMutex.RLock()
	defer s.aMutex.RUnlock()

	if s.avail < 0 {
		return 0
	}
	return s.avail
}

// DrainPermits acquires all available permits and return the number of permits acquired.
func (s *Semaphore) DrainPermits() int {
	n := s.AvailablePermits()
	if n > 0 {
		s.AcquireMany(n)
	}
	return n
}
