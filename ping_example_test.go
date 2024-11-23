package main

import (
	"context"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestPinger(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	r, w := io.Pipe()
	done := make(chan struct{})

	resetTimer := make(chan time.Duration, 1)
	resetTimer <- time.Second

	go func() {
		Pinger(ctx, w, resetTimer)
		close(done)
	}()

	receivePing := func(d time.Duration, r io.Reader) {
		if d >= 0 {
			t.Logf("Resetting timer %s", d)
			resetTimer <- d
		}

		now := time.Now()
		buf := make([]byte, 1024)
		n, err := r.Read(buf)

		if err != nil {
			t.Fatalf("Failed to read from pipe: %v", err)
		}

		t.Logf("Received %q (%s)", buf[:n], time.Since(now).Round(100*time.Millisecond))
	}

	testCases := []struct {
		name    string
		delay   time.Duration
		isFinal bool
	}{
		{"Initial Ping", 0, false},
		{"Short Delay Ping", 200 * time.Millisecond, false},
		{"Medium Delay Ping", 300 * time.Millisecond, false},
		{"Zero Delay Ping", 0, false},
		{"Negative Delay Ping 1", -1, true},
		{"Negative Delay Ping 2", -1, true},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Run %d: %s", i+1, tc.name), func(t *testing.T) {
			receivePing(tc.delay, r)
		})
	}

	cancel()
	<-done
}

