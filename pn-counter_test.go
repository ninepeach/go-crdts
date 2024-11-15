package crdt

import (
	"testing"
	"sync"
	"time"
)

func TestPNCounter_Inc(t *testing.T) {
	pn := NewPNCounter()
	pn.Inc()

	if pn.Count() != 1 {
		t.Errorf("expected count 1, got %d", pn.Count())
	}
}

func TestPNCounter_Dec(t *testing.T) {
	pn := NewPNCounter()
	pn.Inc()
	pn.Dec()

	if pn.Count() != 0 {
		t.Errorf("expected count 0, got %d", pn.Count())
	}
}

func TestPNCounter_IncVal(t *testing.T) {
	pn := NewPNCounter()
	pn.IncVal(5)

	if pn.Count() != 5 {
		t.Errorf("expected count 5, got %d", pn.Count())
	}
}

func TestPNCounter_DecVal(t *testing.T) {
	pn := NewPNCounter()
	pn.IncVal(5)
	pn.DecVal(2)

	if pn.Count() != 3 {
		t.Errorf("expected count 3, got %d", pn.Count())
	}
}

func TestPNCounter_Merge(t *testing.T) {
	pn1 := NewPNCounter()
	pn2 := NewPNCounter()

	pn1.IncVal(5)
	pn2.IncVal(3)
	pn2.DecVal(1)

	pn1.Merge(pn2)

	if pn1.Count() != 7 {
		t.Errorf("expected count 7 after merge, got %d", pn1.Count())
	}
}

func TestPNCounter_ConcurrentIncrement(t *testing.T) {
	pn := NewPNCounter()

	var wg sync.WaitGroup
	const numRoutines = 1000

	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			pn.Inc()
		}()
	}
	wg.Wait()

	if pn.Count() != numRoutines {
		t.Errorf("expected count %d, got %d", numRoutines, pn.Count())
	}
}

func TestPNCounter_ConcurrentDecrement(t *testing.T) {
	pn := NewPNCounter()
	pn.IncVal(1000)

	var wg sync.WaitGroup
	const numRoutines = 1000

	wg.Add(numRoutines)
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer wg.Done()
			pn.Dec()
		}()
	}
	wg.Wait()

	if pn.Count() != 0 {
		t.Errorf("expected count 0, got %d", pn.Count())
	}
}

func TestPNCounter_Performance(t *testing.T) {
	pn := NewPNCounter()
	const numOps = 1000000

	start := time.Now()
	for i := 0; i < numOps; i++ {
		pn.Inc()
	}
	duration := time.Since(start)
	t.Logf("Increment %d operations took %v", numOps, duration)

	start = time.Now()
	for i := 0; i < numOps; i++ {
		pn.Dec()
	}
	duration = time.Since(start)
	t.Logf("Decrement %d operations took %v", numOps, duration)
}
