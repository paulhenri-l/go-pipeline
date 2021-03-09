package stages

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/contracts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewParallel(t *testing.T) {
	_, f := newWaitStageFactory()
	b := NewParallel(f, 2)

	assert.IsType(t, &Parallel{}, b)
}

func TestParallel_Name(t *testing.T) {
	_, f := newWaitStageFactory()
	b := NewParallel(f, 2)

	assert.Equal(t, "Parallel", b.Name())
}

func TestParallel_Start(t *testing.T) {
	_, f := newWaitStageFactory()
	b := NewParallel(f, 2)
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	cancel()
	<-o
}

func TestParallel_Start_StopsWhenInClosed(t *testing.T) {
	_, f := newWaitStageFactory()
	b := NewParallel(f, 2)
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	close(in)
	_, ok := <-o

	assert.False(t, ok)
	cancel()
}

func TestParallel_Start_StopsWhenCancel(t *testing.T) {
	_, f := newWaitStageFactory()
	b := NewParallel(f, 2)
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	cancel()
	_, ok := <-o

	assert.False(t, ok)
}

func TestParallel_Start_ParallelizeStages(t *testing.T) {
	wait := make(chan struct{})
	s1 := &waitStage{wait: wait}
	s2 := &waitStage{wait: wait}
	cnt := 0
	b := NewParallel(func() contracts.Stage {
		cnt++
		if cnt == 1 {
			return s1
		} else if cnt == 2 {
			return s2
		} else {
			t.Error("called more than two times")
			return s2
		}
	}, 2)

	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	// Since parallelism is two I should be able to push two without blocking
	// The wait chan is here to make sure that while pushing events the workers
	// Do not become available again.
	in <- "item"
	in <- "item"
	close(wait)

	assert.Equal(t, "item", <-o)
	assert.Equal(t, "item", <-o)
	assert.Equal(t, 1, s1.receivedEvents)
	assert.Equal(t, 1, s2.receivedEvents)
	cancel()
	<-o
}

func newWaitStageFactory() (chan struct{}, StageFactory) {
	wait := make(chan struct{})

	return wait, func() contracts.Stage {
		return &waitStage{
			wait: wait,
		}
	}
}

type waitStage struct {
	wait           chan struct{}
	receivedEvents int
}

func (ws *waitStage) Name() string {
	return "wait stage"
}

func (ws *waitStage) Start(ctx context.Context, items <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				return
			case i, ok := <-items:
				if !ok {
					return
				}

				ws.receivedEvents++

				<-ws.wait

				out <- i
			}
		}
	}()

	return out
}
