package stages

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/contracts"
	"sync"
)

type StageFactory func() contracts.Stage

type Parallel struct {
	factory     StageFactory
	parallelism int
}

func NewParallel(factory StageFactory, parallelism int) *Parallel {
	return &Parallel{
		factory: factory,
		parallelism: parallelism,
	}
}

func (b *Parallel) Name() string {
	return "Parallel"
}

func (b *Parallel) Start(ctx context.Context, items <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	wg := sync.WaitGroup{}

	for i := 0; i < b.parallelism; i++ {
		localOut := b.factory().Start(ctx, items)
		wg.Add(1)

		go func() {
			defer wg.Done()
			for item := range localOut {
				out <- item
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
