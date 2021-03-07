package generators

import "context"

type Nil struct{}

func NewNil() *Nil {
	return &Nil{}
}

func (n *Nil) Start(ctx context.Context) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		<-ctx.Done()
	}()

	return out
}

