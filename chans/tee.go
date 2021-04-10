package chans

func NewTee(in <-chan interface{}, count int) []<-chan interface{} {
	outs := make([]chan interface{}, count)
	roOuts := make([]<-chan interface{}, count)

	for i := 0; i < count; i++ {
		out := make(chan interface{})
		outs[i] = out
		roOuts[i] = out
	}

	go func() {
		defer func() {
			for _, out := range outs {
				close(out)
			}
		}()

		for i := range in {
			for _, out := range outs {
				// Not that if a reader blocks it will block the whole tee
				// This means that in a pipeline the pipeline will be as slow as
				// The slowest tee reader
				out <- i
			}
		}
	}()

	return roOuts
}
