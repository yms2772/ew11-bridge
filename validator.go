package ew11

import (
	"context"
	"time"
)

type Validatable interface {
	func() error
}

func Validate[T Validatable, E comparable](f T, status *E, want E, delay, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(delay)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f(); err != nil {
				continue
			}
			if *status == want {
				return nil
			}
		}
	}
}
