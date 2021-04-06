package stages

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/paulhenri-l/go-pipeline/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFilter(t *testing.T) {
	b := NewFilter(newFakeFilterStage(t))

	assert.IsType(t, &Filter{}, b)
}

func TestFilter_Name(t *testing.T) {
	b := NewFilter(newFakeFilterStage(t))

	assert.Equal(t, "Filter", b.Name())
}

func TestFilter_Start(t *testing.T) {
	b := NewFilter(newFakeFilterStage(t))
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	cancel()
	<-o
}

func TestFilter_Start_ItemsAreFiltered(t *testing.T) {
	s := newFakeFilterStage(t)
	b := NewFilter(s)

	s.EXPECT().Process(gomock.Any()).DoAndReturn(func(i interface{}) (interface{}, error) {
		return false, nil
	})

	s.EXPECT().Process(gomock.Any()).DoAndReturn(func(i interface{}) (interface{}, error) {
		return true, nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	in <- "filtered"
	in <- "not_filtered"

	i := <-o

	assert.Equal(t, "not_filtered", i)

	cancel()
	<-o
}

func TestFilter_Start_ErrorsAreDiscarded(t *testing.T) {
	var cnt int
	s := newFakeFilterStage(t)
	b := NewFilter(s)

	s.EXPECT().Process(gomock.Any()).DoAndReturn(func(i interface{}) (interface{}, error) {
		cnt++
		if cnt == 1 {
			return true, nil
		}

		return nil, errors.New("huh oh")
	}).Times(2)

	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	in <- "no_error"
	i1 := <-o

	in <- "error"
	cancel()
	i2, ok := <-o

	assert.Equal(t, "no_error", i1)
	assert.Nil(t, i2)
	assert.False(t, ok)

	<-o
}

func TestFilter_Start_StopsWhenInClosed(t *testing.T) {
	b := NewFilter(newFakeFilterStage(t))
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	close(in)
	_, ok := <-o

	assert.False(t, ok)
	cancel()
}

func TestFilter_Start_StopsWhenCancel(t *testing.T) {
	b := NewFilter(newFakeFilterStage(t))
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	cancel()
	_, ok := <-o

	assert.False(t, ok)
}

func newFakeFilterStage(t *testing.T) *mocks.MockFilterStage {
	ctl := gomock.NewController(t)

	return mocks.NewMockFilterStage(ctl)
}
