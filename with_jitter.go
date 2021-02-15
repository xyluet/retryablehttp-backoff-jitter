package backoff

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

var rander = rand.New(rand.NewSource(time.Now().UnixNano()))

type withJitter struct {
	retryablehttp.Backoff
	maxDuration time.Duration
}

func (j *withJitter) backoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	sleep := j.Backoff(min, max, attemptNum, resp)
	jitter := time.Duration(rander.Int63n(int64(j.maxDuration)))
	sleepWithJitter := sleep + jitter
	if sleepWithJitter > max {
		return max
	}
	return sleepWithJitter
}

// WithJitterOption configures backoff with jitter.
type WithJitterOption func(*withJitter)

// WithJitterMaxDuration sets the jitter max duration.
func WithJitterMaxDuration(dur time.Duration) WithJitterOption {
	return func(j *withJitter) { j.maxDuration = dur }
}

// WithJitter decorates the retryablehttp Backoff and adds jitter.
func WithJitter(backoff retryablehttp.Backoff, options ...WithJitterOption) retryablehttp.Backoff {
	j := &withJitter{
		Backoff:     backoff,
		maxDuration: time.Second,
	}
	for _, opt := range options {
		opt(j)
	}
	return j.backoff
}
