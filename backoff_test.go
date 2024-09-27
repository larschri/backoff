package backoff

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

type mockSleep time.Duration

func (m *mockSleep) after(d time.Duration) <-chan time.Time {
	*m = mockSleep(d)
	return time.NewTimer(0).C
}

func durations(s string) []time.Duration {
	var t []time.Duration
	for _, s := range strings.Split(s, ";") {
		d, err := time.ParseDuration(s)
		if err != nil {
			panic(err)
		}
		t = append(t, d)
	}
	return t
}

var tests = map[string]struct {
	opts   []func(*options)
	sleeps []time.Duration
}{
	"default": {
		opts:   []func(*options){},
		sleeps: durations("0s;1s;2s;4s;8s;16s;32s;1m4s;2m8s;4m16s;8m32s;17m4s;34m8s;1h0m0s"),
	},
	"multiplier5": {
		opts:   []func(*options){WithMultiplier(5)},
		sleeps: durations("0s;1s;5s;25s;2m5s;10m25s;52m5s;1h0m0s"),
	},
	"range": {
		opts:   []func(*options){WithMin(3 * time.Second), WithMax(12 * time.Second)},
		sleeps: durations("0s;3s;6s;12s"),
	},
}

func TestExponential(t *testing.T) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var m mockSleep
			after = m.after
			for i := range Exponential(context.Background(), test.opts...) {
				if i > 100 {
					break
				}
				idx := len(test.sleeps) - 1
				if i < idx {
					idx = i
				}
				if time.Duration(m) != test.sleeps[idx] {
					t.Errorf("%d: expected %v, got %v", i, test.sleeps[idx], time.Duration(m))
				}
			}
		})
	}
}

func TestWithTerminate(t *testing.T) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var m mockSleep
			after = m.after
			opts := append(test.opts, WithTerminate(true))
			n := 0
			for _ = range Exponential(context.Background(), opts...) {
				n++
			}
			if n != len(test.sleeps) {
				t.Errorf("expected %v iterations, had %v", len(test.sleeps), n)
			}
		})
	}
}

func ExampleExponential() {
	for range Exponential(context.Background(),
		WithMultiplier(2),    // configure multiplier - 2 is default
		WithMin(time.Second), // configure initial wait time - 1 second is default
		WithMax(time.Hour),   // configure max wait time - 1 hour is default
		WithTerminate(false), // configure the loop to run forever
	) {
		if _, err := os.Hostname(); err == nil {
			fmt.Println("yay, hostname")
			break
		}
	}
	// Output: yay, hostname
}
