package pipeline

import (
	"context"
	"sync/atomic"
	"time"
)

type numGen struct {
	max int
}

func (n *numGen) Start(ctx context.Context) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for i := 0; i < n.max; i++ {
			out <- i
		}

		<-ctx.Done()
	}()

	return out
}

type counterStage struct{
	counter *int64
}

func (c counterStage) Name() string {
	return "counterstage"
}

func (c *counterStage) Start(ctx context.Context, events <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return

			case e, ok := <-events:
				if !ok {
					return
				}
				atomic.AddInt64(c.counter, 1)
				out <- e
			}
		}
	}()

	return out
}

type waitStage struct {
	wait time.Duration
}

func (ws *waitStage) Name() string {
	return "waitstage"
}

func (ws *waitStage) Start(ctx context.Context, events <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return

			case e, ok := <-events:
				if !ok {
					return
				}
				<-time.NewTimer(ws.wait).C
				out <- e
			}
		}
	}()

	return out
}
