package internal

import (
	"sync"
)

// Counter provides thread-safe access to a count value.
type Counter struct {
	mu    sync.Mutex
	value int64
}

// NewCounter creates a new Counter initialized to 0.
func NewCounter() *Counter {
	return &Counter{
		value: 0,
	}
}

// Get returns the current count value.
func (c *Counter) Get() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.value
}

// Increment increments the counter by 1 and returns the new value.
func (c *Counter) Increment() int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
	return c.value
}

// Set sets the counter to a specific value.
func (c *Counter) Set(value int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = value
}
