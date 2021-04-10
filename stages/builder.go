package stages

import "context"

type Builder struct {
	blueprint BlueprintFunc
}

type BlueprintFunc func(ctx context.Context, in <-chan interface{}) <-chan interface{}

func NewBuilder(bp BlueprintFunc) *Builder {
	return &Builder{
		blueprint: bp,
	}
}

func (m *Builder) Name() string {
	return "Builder"
}

func (m *Builder) Start(ctx context.Context, items <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		o := m.blueprint(ctx, items)

		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-o:
				if !ok {
					return
				}

				out <- i
			}
		}
	}()

	return out
}
