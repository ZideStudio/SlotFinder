package auth

import (
	"time"
)

// time.NewTicker can't be intercepted directly, so cleanRefreshTokens takes
// a ticker built from this seam, letting tests swap in a fake clock instead
// of waiting on real wall-clock time.
type clock interface {
	NewTicker(d time.Duration) ticker
}

type ticker interface {
	C() <-chan time.Time
	Stop()
}

type realClock struct{}

func (realClock) NewTicker(d time.Duration) ticker {
	return &realTicker{t: time.NewTicker(d)}
}

type realTicker struct{ t *time.Ticker }

func (r *realTicker) C() <-chan time.Time { return r.t.C }
func (r *realTicker) Stop()               { r.t.Stop() }
