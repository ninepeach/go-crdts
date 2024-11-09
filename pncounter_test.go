package crdt

import (
	"testing"
)

func TestPNCounterIncrementDecrement(t *testing.T) {
	counter := NewPNCounter()
	counter.Increment("node1", 5)
	counter.Decrement("node1", 2)

	if counter.Value() != 3 {
		t.Errorf("Expected 3, got %d", counter.Value())
	}
}

func TestPNCounterMerge(t *testing.T) {
	counter1 := NewPNCounter()
	counter1.Increment("node1", 5)
	counter1.Decrement("node1", 2)

	counter2 := NewPNCounter()
	counter2.Increment("node2", 3)
	counter2.Decrement("node2", 1)

	counter1.Merge(counter2)

	expectedValue := 5
	if counter1.Value() != expectedValue {
		t.Errorf("Expected %d, got %d", expectedValue, counter1.Value())
	}
}
