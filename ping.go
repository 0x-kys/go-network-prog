package main

import (
	"context"
	"io"
	"time"
)

const defaultPingInterval = 30 * time.Second

func Pinger(ctx context.Context, w io.Writer, reset <-chan time.Duration) {
	var interval time.Duration

	// explain this select
	// also, why select?
	select {
	case <-ctx.Done():
		return
	case interval = <-reset:
	default:
	}

	if interval <= 0 {
		interval = defaultPingInterval
	}

	timer := time.NewTimer(interval)
	// explain this function
	defer func() {
		if !timer.Stop() {
			<-timer.C
		}
	}() // explain

	// explain this for, select & all the cases
	for {
		select {
		case <-ctx.Done():
			return
		case newInterval := <-reset:
			if !timer.Stop() {
				<-timer.C
			}

			if newInterval > 0 {
				interval = newInterval
			}
		case <-timer.C:
			if _, err := w.Write([]byte("ping")); err != nil {
				// track and act on consecutive timeout here
				return
			}
		}

		_ = timer.Reset(interval)
	}
}
