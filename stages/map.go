package stages

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/contracts"
)

type Map struct {
	stage contracts.MapStage
}

func NewMap(stage contracts.MapStage) *Map {
	return &Map{
		stage: stage,
	}
}

func (m *Map) Name() string {
	return "Map"
}

func (m *Map) Start(ctx context.Context, items <-chan interface{}) <-chan interface{} {
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

				i, err := m.stage.Process(item)
				if err != nil {
					// Should probably log or do something
					break
				}

				out <- i
			}
		}
	}()

	return out
}
