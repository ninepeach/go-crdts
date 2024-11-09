package crdt

import "sync"

// GCounter is a grow-only counter for distributed systems.
type GCounter struct {
	counts map[string]int
	mu     sync.Mutex
}

// NewGCounter initializes a new GCounter.
func NewGCounter() *GCounter {
	return &GCounter{counts: make(map[string]int)}
}

// Increment increases the counter for the given node ID.
func (g *GCounter) Increment(nodeID string, value int) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.counts[nodeID] += value
}

// Value returns the total count by summing values from all nodes.
func (g *GCounter) Value() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	total := 0
	for _, v := range g.counts {
		total += v
	}
	return total
}

// Merge updates the counter with another GCounter by taking the max count per node.
func (g *GCounter) Merge(other *GCounter) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for nodeID, count := range other.counts {
		if count > g.counts[nodeID] {
			g.counts[nodeID] = count
		}
	}
}
