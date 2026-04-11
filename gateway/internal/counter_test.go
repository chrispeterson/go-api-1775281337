package internal

import (
	"sync"
	"testing"
)

func TestNewCounter(t *testing.T) {
	c := NewCounter()
	if c == nil {
		t.Fatal("expected counter, got nil")
	}

	if c.Get() != 0 {
		t.Errorf("expected initial value 0, got %d", c.Get())
	}
}

func TestCounterGet(t *testing.T) {
	c := NewCounter()

	if got := c.Get(); got != 0 {
		t.Errorf("Get() = %d, want 0", got)
	}
}

func TestCounterSet(t *testing.T) {
	c := NewCounter()

	c.Set(42)
	if got := c.Get(); got != 42 {
		t.Errorf("after Set(42), Get() = %d, want 42", got)
	}

	c.Set(100)
	if got := c.Get(); got != 100 {
		t.Errorf("after Set(100), Get() = %d, want 100", got)
	}
}

func TestCounterIncrement(t *testing.T) {
	c := NewCounter()

	val := c.Increment()
	if val != 1 {
		t.Errorf("first Increment() = %d, want 1", val)
	}

	val = c.Increment()
	if val != 2 {
		t.Errorf("second Increment() = %d, want 2", val)
	}

	val = c.Increment()
	if val != 3 {
		t.Errorf("third Increment() = %d, want 3", val)
	}

	if got := c.Get(); got != 3 {
		t.Errorf("after 3 Increments, Get() = %d, want 3", got)
	}
}

func TestCounterThreadSafety(t *testing.T) {
	c := NewCounter()
	const numGoroutines = 100
	const incrementsPerGoroutine = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrementsPerGoroutine; j++ {
				c.Increment()
			}
		}()
	}

	wg.Wait()

	expected := int64(numGoroutines * incrementsPerGoroutine)
	if got := c.Get(); got != expected {
		t.Errorf("after concurrent increments, Get() = %d, want %d", got, expected)
	}
}

func TestCounterConcurrentGetSet(t *testing.T) {
	c := NewCounter()
	const numGoroutines = 50

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2)

	// Half goroutines increment
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			c.Increment()
		}()
	}

	// Half goroutines read
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			c.Get()
		}()
	}

	wg.Wait()

	if got := c.Get(); got != int64(numGoroutines) {
		t.Errorf("after concurrent reads/increments, Get() = %d, want %d", got, int64(numGoroutines))
	}
}
