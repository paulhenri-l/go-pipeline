package sinks

import (
	"context"
)

type Nil struct{}

func (n *Nil) Start(ctx context.Context, in <-chan interface{}) <-chan struct{}{
	done := make(chan struct{})

	go func() {
		defer close(done)

		for {
			select {
			case <-ctx.Done():
				return

			case _, ok := <-in:
				if !ok {
					return
				}
			}
		}
	}()

	return done
}
