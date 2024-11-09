package crdt

import (
	"testing"
)

func TestGCounterIncrement(t *testing.T) {
	counter := NewGCounter()
	counter.Increment("node1", 5)
	if counter.Value() != 5 {
		t.Errorf("Expected 5, got %d", counter.Value())
	}
}

func TestGCounterMerge(t *testing.T) {
	counter1 := NewGCounter()
	counter1.Increment("node1", 5)

	counter2 := NewGCounter()
	counter2.Increment("node1", 3)
	counter2.Increment("node2", 4)

	counter1.Merge(counter2)

	if counter1.Value() != 9 {
		t.Errorf("Expected 9, got %d", counter1.Value())
	}
}
