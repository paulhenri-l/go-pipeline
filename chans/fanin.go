package chans

import "sync"

func NewFanIn(ins ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	wg := sync.WaitGroup{}

	for _, in := range ins {
		wg.Add(1)

		go func(in <-chan interface{}) {
			defer wg.Done()

			for i := range in {
				out <- i
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
