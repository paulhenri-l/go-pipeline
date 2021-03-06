package stages

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/contracts"
)

type FlatMap struct {
	stage contracts.FlatMapStage
}

func NewFlatMap(stage contracts.FlatMapStage) *FlatMap {
	return &FlatMap{
		stage: stage,
	}
}

func (m *FlatMap) Name() string {
	return "FlatMap"
}

func (m *FlatMap) Start(ctx context.Context, items <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return

			case item, ok := <-items:
				if !ok {
					return
				}

				for i := range m.stage.Process(ctx, item) {
					out <- i
				}
			}
		}
	}()

	return out
}
