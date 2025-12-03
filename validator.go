package ew11

import (
	"context"
	"time"
)

type Validatable[T comparable] interface {
	func(T) error
}

func Validate[T Validatable[E], E comparable](f T, status *E, want E, delay, timeout time.Duration) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		ticker := time.NewTicker(delay)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
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
	}()
}
