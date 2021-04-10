package chans

func NewFanOut(in <-chan interface{}, count int) []<-chan interface{} {
	outs := make([]<-chan interface{}, count)

	for i := 0; i < count; i++ {
		out := make(chan interface{})

		go func(out chan interface{}) {
			defer close(out)

			for i := range in {
				out <- i
			}
		}(out)

		outs[i] = out
	}

	return outs
}
