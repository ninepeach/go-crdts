package crdt

// PNCounter represents a state-based PN-Counter. It is
// implemented as sets of two G-Counters, one that tracks
// increments while the other decrements.
type PNCounter struct {
	pCounter *GCounter // Counter for increments
	nCounter *GCounter // Counter for decrements
}

// NewPNCounter returns a new *PNCounter with both its
// G-Counters initialized.
func NewPNCounter() *PNCounter {
	return &PNCounter{
		pCounter: NewGCounter(), // Initialize the positive GCounter
		nCounter: NewGCounter(), // Initialize the negative GCounter
	}
}

// Inc monotonically increments the current value of the
// PN-Counter by one.
func (pn *PNCounter) Inc() {
	pn.IncVal(1)
}

// IncVal increments the current value of the PN-Counter
// by the delta incr that is provided. The value of delta
// has to be >= 0. If the value of delta is < 0, then this
// implementation panics.
func (pn *PNCounter) IncVal(incr int) {
	if incr < 0 {
		panic("cannot increment with a negative value")
	}
	pn.pCounter.IncVal(incr) // Increment the positive counter
}

// Dec monotonically decrements the current value of the
// PN-Counter by one.
func (pn *PNCounter) Dec() {
	pn.DecVal(1)
}

// DecVal adds a decrement to the current value of
// PN-Counter by the value of delta decr. Similar to
// IncVal, the value of decr cannot be less than 0.
func (pn *PNCounter) DecVal(decr int) {
	if decr < 0 {
		panic("cannot decrement with a negative value")
	}
	pn.nCounter.IncVal(decr) // Increment the negative counter (effectively decrement)
}

// Count returns the current value of the counter. It
// subtracts the value of negative G-Counter from the
// positive grow-only counter and returns the result.
// Because this counter can grow in either direction,
// negative integers as results are possible.
func (pn *PNCounter) Count() int {
	return pn.pCounter.Count() - pn.nCounter.Count() // Positive counter minus negative counter
}

// Merge combines both the PN-Counters and saves the result
// in the invoking counter. Respective G-Counters are merged
// i.e. +ve with +ve and -ve with -ve, but no computation
// is actually performed for the count, just the merging of states.
func (pn *PNCounter) Merge(pnpn *PNCounter) {
	// Merge positive counters with positive counters and negative with negative
	pn.pCounter.Merge(pnpn.pCounter)
	pn.nCounter.Merge(pnpn.nCounter)
}
