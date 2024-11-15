package crdt

import (
	"sync"
	"testing"
	"time"
)

func TestGCounter_Inc(t *testing.T) {
	gcounter := NewGCounter()

	// Test a single increment
	gcounter.Inc()
	if gcounter.Count() != 1 {
		t.Errorf("Expected count to be 1, got %d", gcounter.Count())
	}

	// Test multiple increments
	gcounter.IncVal(5)
	if gcounter.Count() != 6 {
		t.Errorf("Expected count to be 6, got %d", gcounter.Count())
	}
}

func TestGCounter_Merge(t *testing.T) {
	gcounter1 := NewGCounter()
	gcounter2 := NewGCounter()

	// Increment counter in both replicas
	gcounter1.IncVal(3)
	gcounter2.IncVal(5)

	// Merge counters
	gcounter1.Merge(gcounter2)

	// After merging, the total should be the max of the two
	if gcounter1.Count() != 8 {
		t.Errorf("Expected merged count to be 8, got %d", gcounter1.Count())
	}
}

func TestGCounter_ConcurrentIncrement(t *testing.T) {
	gcounter := NewGCounter()
	var wg sync.WaitGroup

	// Number of goroutines
	numGoroutines := 100
	incrementAmount := 1

	// Using WaitGroup to wait for all goroutines to finish
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gcounter.IncVal(incrementAmount)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// The final count should be equal to numGoroutines * incrementAmount
	expected := numGoroutines * incrementAmount
	if gcounter.Count() != expected {
		t.Errorf("Expected count to be %d, got %d", expected, gcounter.Count())
	}
}

func TestGCounter_Performance(t *testing.T) {
	gcounter := NewGCounter()
	var wg sync.WaitGroup

	// Number of goroutines for performance test
	numGoroutines := 100000
	incrementAmount := 1

	// Start measuring time
	start := time.Now()

	// Increment counter concurrently
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gcounter.IncVal(incrementAmount)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Measure elapsed time
	duration := time.Since(start)

	// Calculate operations per second
	opsPerSecond := float64(numGoroutines) / duration.Seconds()

	// Print the result
	t.Logf("Processed %d increments in %v, which is %.2f ops/sec", numGoroutines, duration, opsPerSecond)

	// Check if performance is within expected time range (e.g., 1 second for 100,000 increments)
	if duration >= 2*time.Second {
		t.Errorf("Expected performance to be under 2 seconds, but got %v", duration)
	}

	// The final count should match the number of increments
	expected := numGoroutines * incrementAmount
	if gcounter.Count() != expected {
		t.Errorf("Expected count to be %d, got %d", expected, gcounter.Count())
	}
}