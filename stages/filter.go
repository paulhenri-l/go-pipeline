package stages

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/contracts"
)

type Filter struct {
	stage contracts.FilterStage
}

func NewFilter(stage contracts.FilterStage) *Filter {
	return &Filter{
		stage: stage,
	}
}

func (m *Filter) Name() string {
	return "Filter"
}

func (m *Filter) Start(ctx context.Context, items <-chan interface{}) <-chan interface{} {
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

				passes, err := m.stage.Process(item)
				if err != nil {
					// Should probably log or do something
					break
				}

				if passes {
					out <- item
				}
			}
		}
	}()

	return out
}
