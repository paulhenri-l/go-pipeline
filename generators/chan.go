package generators

import "context"

type Chan struct {
	in <-chan interface{}
}

func NewChan(in <-chan interface{}) *Chan {
	return &Chan{
		in: in,
	}
}

func (c *Chan) Start(ctx context.Context) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return

			case item, ok := <- c.in:
				if !ok {
					return
				}

				out<-item
			}
		}
	}()

	return out
}
