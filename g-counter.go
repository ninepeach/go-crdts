package crdt

import (
	"sync"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// EnableMutex is a constant that controls whether the mutex lock is used.
// Set to true to enable locking, false to disable locking.
const EnableMutex = true // Set this to false to disable locks

// GenerateUUID generates a unique identifier using timestamp and random bytes.
func GenerateUUID() string {
	// Get the current Unix timestamp in nanoseconds
	timestamp := time.Now().UnixNano()

	// Generate 8 random bytes
	randomBytes := make([]byte, 8)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic("failed to generate random bytes")
	}

	// Convert random bytes to hex string
	randomHex := hex.EncodeToString(randomBytes)

	// Format the UUID-like string
	return fmt.Sprintf("%d-%s", timestamp, randomHex)
}

// GCounter represents a G-counter in CRDT, which is a state-based grow-only counter
// that only supports increments.
type GCounter struct {
	// ident provides a unique identity to each replica.
	ident string

	// counter maps identity of each replica to their entry values
	// i.e., the counter value they individually have.
	counter map[string]int

	// Mutex for thread-safe access to the counter
	mu sync.Mutex
}

// NewGCounter returns a *GCounter by pre-assigning a unique identity to it.
func NewGCounter() *GCounter {
	return &GCounter{
		ident:   GenerateUUID(), // Unique ID for each machine/process
		counter: make(map[string]int),
	}
}

// Inc increments the GCounter by the value of 1 every time it is called.
func (g *GCounter) Inc() {
	g.IncVal(1)
}

// IncVal allows passing in an arbitrary delta to increment the current value of counter.
// Only positive values are accepted. If a negative value is provided, the implementation will panic.
func (g *GCounter) IncVal(incr int) {
	if incr < 0 {
		panic("cannot decrement a gcounter")
	}

	// Decide whether to use mutex lock based on EnableMutex constant
	if EnableMutex {
		g.mu.Lock()
		defer g.mu.Unlock()
	}

	// Increment the counter for the current machine's ID
	g.counter[g.ident] += incr
}

// Count returns the total count of this counter across all the present replicas.
func (g *GCounter) Count() (total int) {
	// Sum the values in the counter map
	for _, val := range g.counter {
		total += val
	}
	return
}

// Merge combines the counter values across multiple replicas.
// The property of idempotency is preserved here across multiple merges.
func (g *GCounter) Merge(c *GCounter) {
	// Decide whether to use mutex lock based on EnableMutex constant
	if EnableMutex {
		g.mu.Lock()
		defer g.mu.Unlock()
	}

	for ident, val := range c.counter {
		// Only keep the maximum value for each replica's counter
		if v, ok := g.counter[ident]; !ok || v < val {
			g.counter[ident] = val
		}
	}
}