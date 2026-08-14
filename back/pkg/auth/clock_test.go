package auth

import (
	"sync"
	"time"
)

// fakeClock lets tests advance time deterministically to fire tickers
// created against the real production interval.
type fakeClock struct {
	mu      sync.Mutex
	now     time.Time
	tickers []*fakeTicker
}

func newFakeClock() *fakeClock {
	return &fakeClock{now: time.Now()}
}

func (f *fakeClock) NewTicker(d time.Duration) ticker {
	f.mu.Lock()
	defer f.mu.Unlock()

	t := &fakeTicker{c: make(chan time.Time, 1), interval: d, next: f.now.Add(d)}
	f.tickers = append(f.tickers, t)
	return t
}

// Advance moves the fake clock forward by d, firing every ticker whose next
// tick falls within the new time.
func (f *fakeClock) Advance(d time.Duration) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.now = f.now.Add(d)
	for _, t := range f.tickers {
		t.fireDue(f.now)
	}
}

type fakeTicker struct {
	mu       sync.Mutex
	c        chan time.Time
	interval time.Duration
	next     time.Time
	stopped  bool
}

func (t *fakeTicker) fireDue(now time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.stopped {
		return
	}
	for !t.next.After(now) {
		select {
		case t.c <- t.next:
		default:
		}
		t.next = t.next.Add(t.interval)
	}
}

func (t *fakeTicker) C() <-chan time.Time { return t.c }

func (t *fakeTicker) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.stopped = true
}
