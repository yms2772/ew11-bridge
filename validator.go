package ew11

import (
	"context"
	"time"
)

type Validatable[T comparable] interface {
	func(T) error
}

func Validate[T Validatable[E], E comparable](f T, status *E, want E, delay, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if err := f(want); err != nil {
				continue
			}

			time.Sleep(delay)
			if *status == want {
				return nil
			}
		}
	}
}
