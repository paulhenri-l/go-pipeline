package stages

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/chans"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	b := NewBuilder(newSimpleBluePrint())

	assert.IsType(t, &Builder{}, b)
}

func TestBuilder_Name(t *testing.T) {
	b := NewBuilder(newSimpleBluePrint())

	assert.Equal(t, "Builder", b.Name())
}

func TestBuilder_Start(t *testing.T) {
	b := NewBuilder(newSimpleBluePrint())
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	cancel()
	<-o
}

func TestBuilder_Start_ItemsArePassedToTheBluePrint(t *testing.T) {
	b := NewBuilder(newSimpleBluePrint())
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	in <- "hello"
	i := <-o

	assert.Equal(t, "hello", i)
	cancel()
	<-o
}

func TestBuilder_Start_StopsWhenInClosed(t *testing.T) {
	b := NewBuilder(newSimpleBluePrint())
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)
	close(in)
	_, ok := <-o

	assert.False(t, ok)
	cancel()
}

func TestBuilder_Start_StopsWhenCancel(t *testing.T) {
	b := NewBuilder(newSimpleBluePrint())
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)
	cancel()
	_, ok := <-o

	assert.False(t, ok)
	cancel()
}

func newSimpleBluePrint() BlueprintFunc {
	return func(ctx context.Context, in <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})

		go func() {
			defer close(out)

			for i := range chans.NewOrDone(ctx, in) {
				out <- i
			}
		}()

		return out
	}
}
