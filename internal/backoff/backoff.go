// Package backoff implement the backoff strategy for WeCross-Go-SDK
package backoff

import (
	"time"

	"github.com/WeBankBlockchain/WeCross-Go-SDK/internal/wecrossrand"
)

// Stratgy defines the methodology for backing off after a connection
// failure.
type Stratgy interface {
	// Backoff returns the amount of time to wait before the next retry given
	// the number of consecutive failures.
	Backoff(retries int) time.Duration
}

// DefaultExponential is an exponential backoff implementation using the
// default values.
var DefaultExponential = Exponential{Config: DefaultConfig}

// Exponential implements exponential backoff algorithm as defined.
type Exponential struct {
	// Config contains all options to configure the backoff algorithm.
	Config Config
}

// Backoff returns the amount of time to wait before the next retry given the
// number of retries.
func (bc Exponential) Backoff(retries int) time.Duration {
	if retries == 0 {
		return bc.Config.BaseDelay
	}
	backoff, max := float64(bc.Config.BaseDelay), float64(bc.Config.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= bc.Config.Multiplier
		retries--
	}
	if backoff > max {
		backoff = max
	}
	// Randomize backoff delays so that if other requests start at
	// the same time, they won't operate in lockstep.
	backoff *= 1 + bc.Config.Jitter*(wecrossrand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}

// The following structs are kept in internal until the WeCross-Go-SDK project decides whether or not to
// allow alternative Config.
//
// Config defines the configuration options for backoff.
type Config struct {
	// BaseDelay is the amount of time to backoff after the first failure.
	BaseDelay time.Duration
	// Multiplier is the factor with which to multiply backoffs after a
	// failed retry. Should ideally be greater than 1.
	Multiplier float64
	// Jitter is the factor whith which backoffs are randomized.
	Jitter float64
	// MaxDelay is the upper bound of backoff delay.
	MaxDelay time.Duration
}

// DefaultConfig is a backoff configuration with the default values.
var DefaultConfig = Config{
	BaseDelay:  1.0 * time.Second,
	Multiplier: 1.6,
	Jitter:     0.2,
	MaxDelay:   120 * time.Second,
}
