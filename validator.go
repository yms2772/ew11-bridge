package ew11

import "time"

func Validate[T comparable](f func(T) error, status *T, want T, delay, timeout time.Duration) {
	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-time.After(timeout):
			return
		case <-ticker.C:
			if err := f(want); err != nil {
				continue
			}
			if *status == want {
				return
			}
		}
	}
}
