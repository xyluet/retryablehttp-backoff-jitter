package backoff_test

import (
	"net/http"
	"testing"
	"time"

	backoff "github.com/xyluet/retryablehttp-backoff-jitter"
)

func durationInBetween(t *testing.T, in, min, max time.Duration) {
	if in < min || in > max {
		t.Fatalf("%s is not between %s and %s", in, min, max)
	}
}

func TestWithJitter(t *testing.T) {
	t.Run("attemp num 0", func(t *testing.T) {
		backoffWithJitter := backoff.WithJitter(
			func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
				return 5 * time.Second
			},
		)
		sleep := backoffWithJitter(time.Second, 30*time.Second, 0, nil)
		durationInBetween(t, sleep, 5*time.Second, 6*time.Second)
	})

	t.Run("MUST less or equal than max", func(t *testing.T) {
		backoffWithJitter := backoff.WithJitter(
			func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
				return 30 * time.Second
			},
		)
		sleep := backoffWithJitter(time.Second, 30*time.Second, 5, nil)
		durationInBetween(t, sleep, 30*time.Second, 30*time.Second)
	})

	t.Run("config max jitter duration", func(t *testing.T) {
		backoffWithJitter := backoff.WithJitter(
			func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
				return 1 * time.Second
			},
			backoff.WithJitterMaxDuration(100*time.Millisecond),
		)
		sleep := backoffWithJitter(time.Second, 30*time.Second, 0, nil)
		durationInBetween(t, sleep, 1*time.Second, 1100*time.Millisecond)
	})
}
