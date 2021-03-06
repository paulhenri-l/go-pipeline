package stages

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/paulhenri-l/go-pipeline/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFlatMap(t *testing.T) {
	b := NewFlatMap(newFakeFlatMapStage(t))

	assert.IsType(t, &FlatMap{}, b)
}

func TestFlatMap_Name(t *testing.T) {
	b := NewFlatMap(newFakeFlatMapStage(t))

	assert.Equal(t, "FlatMap", b.Name())
}

func TestFlatMap_Start(t *testing.T) {
	b := NewFlatMap(newFakeFlatMapStage(t))
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	cancel()
	<-o
}

func TestFlatMap_Start_ItemsAreFlatMapped(t *testing.T) {
	s := newFakeFlatMapStage(t)
	b := NewFlatMap(s)

	s.EXPECT().Process(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, i interface{}) <-chan interface{} {
		out := make(chan interface{}, 2)
		out<-i
		out<-i
		close(out)
		return out
	}).Times(2)

	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	in <- 1
	i1a := <-o
	i1b := <-o

	in <- 2
	i2a := <-o
	i2b := <-o

	assert.Equal(t, 1, i1a)
	assert.Equal(t, 1, i1b)
	assert.Equal(t, 2, i2a)
	assert.Equal(t, 2, i2b)

	cancel()
	<-o
}

func TestFlatMap_Start_StopsWhenInClosed(t *testing.T) {
	b := NewFlatMap(newFakeFlatMapStage(t))
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	close(in)
	_, ok := <-o

	assert.False(t, ok)
	cancel()
}

func TestFlatMap_Start_StopsWhenCancel(t *testing.T) {
	b := NewFlatMap(newFakeFlatMapStage(t))
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	cancel()
	_, ok := <-o

	assert.False(t, ok)
}

func newFakeFlatMapStage(t *testing.T) *mocks.MockFlatMapStage {
	ctl := gomock.NewController(t)

	return mocks.NewMockFlatMapStage(ctl)
}
