package crdt

// PNCounter is a Positive-Negative counter using two GCounters for increment and decrement.
type PNCounter struct {
	inc, dec *GCounter
}

// NewPNCounter initializes a new PNCounter.
func NewPNCounter() *PNCounter {
	return &PNCounter{
		inc: NewGCounter(),
		dec: NewGCounter(),
	}
}

// Increment increases the positive counter for the node.
func (pn *PNCounter) Increment(nodeID string, value int) {
	pn.inc.Increment(nodeID, value)
}

// Decrement increases the negative counter for the node.
func (pn *PNCounter) Decrement(nodeID string, value int) {
	pn.dec.Increment(nodeID, value)
}

// Value returns the net result (inc - dec).
func (pn *PNCounter) Value() int {
	return pn.inc.Value() - pn.dec.Value()
}

// Merge merges another PNCounter by merging both inc and dec counters.
func (pn *PNCounter) Merge(other *PNCounter) {
	pn.inc.Merge(other.inc)
	pn.dec.Merge(other.dec)
}
