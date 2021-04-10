package chans

import "context"

func NewOrDone(ctx context.Context, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-in:
				if !ok {
					return
				}

				out<-i
			}
		}
	}()

	return out
}
