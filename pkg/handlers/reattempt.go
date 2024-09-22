package handlers

import (
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
}

// Retry mechanism
func reattempt(attempts int, sleep time.Duration, f func() error) error {
	if err := f(); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep))) //nolint:gosec
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return reattempt(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}