// Package backoff provides a function for running a for loop with exponential
// backoff.
package backoff

import (
	"context"
	"math"
	"time"
)

// Exponential is for running a for loop with exponential backoff. The returned
// function should be used in the range expression, and controls the delay
// between each iteration.
//
// The default behavior is:
//   - multiplier = 2, use [WithMultiplier] to override
//   - min = 1s, use [WithMin] to override
//   - max = 1h, use [WithMax] to override
//   - run forever, use [WithTerminate] to override
func Exponential(ctx context.Context, opts ...func(*options)) func(func(int) bool) {
	o := options{
		time.Second,
		time.Hour,
		2,
		false,
		ctx,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o.rangeFunc
}

type options struct {
	min        time.Duration
	max        time.Duration
	multiplier float64
	terminate  bool
	ctx        context.Context
}

// after is [time.After] (except in tests, where it is overridden)
var after = time.After

func (o *options) rangeFunc(yield func(x int) bool) {
	s := float64(o.min)
	max := float64(o.max)
	for i := 0; yield(i); i++ {
		select {
		case <-after(time.Duration(s)):
		case <-o.ctx.Done():
			return
		}
		if s == max && o.terminate {
			yield(i + 1)
			return
		}
		s = math.Min(s*o.multiplier, max)
	}
}

// WithMin sets the duration of the sleep after the first iteration.
func WithMin(min time.Duration) func(*options) {
	return func(o *options) {
		o.min = min
	}
}

// WithMax sets maximum sleep duration.
func WithMax(max time.Duration) func(*options) {
	return func(o *options) {
		o.max = max
	}
}

// WithMultiplier sets the factor by which the sleep duration is incremented in
// every iteration.
func WithMultiplier(multiplier float64) func(*options) {
	return func(o *options) {
		o.multiplier = multiplier
	}
}

// WithTerminate terminates the iteration when the max sleep duration as
// configured by [WithMax] has been reached.
func WithTerminate(t bool) func(*options) {
	return func(o *options) {
		o.terminate = true
	}
}
