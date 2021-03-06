package stages

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/paulhenri-l/go-pipeline/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMap(t *testing.T) {
	b := NewMap(newFakeMapStage(t))

	assert.IsType(t, &Map{}, b)
}

func TestMap_Name(t *testing.T) {
	b := NewMap(newFakeMapStage(t))

	assert.Equal(t, "Map", b.Name())
}

func TestMap_Start(t *testing.T) {
	b := NewMap(newFakeMapStage(t))
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	cancel()
	<-o
}

func TestMap_Start_ItemsAreMapped(t *testing.T) {
	s := newFakeMapStage(t)
	b := NewMap(s)

	s.EXPECT().Process(gomock.Any()).DoAndReturn(func(i interface{}) (interface{}, error) {
		return i.(int) * 2, nil
	}).Times(2)

	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	in <- 1
	i1 := <-o

	in <- 2
	i2 := <-o

	assert.Equal(t, 2, i1)
	assert.Equal(t, 4, i2)

	cancel()
	<-o
}

func TestMap_Start_ErrorsAreDiscarded(t *testing.T) {
	var cnt int
	s := newFakeMapStage(t)
	b := NewMap(s)

	s.EXPECT().Process(gomock.Any()).DoAndReturn(func(i interface{}) (interface{}, error) {
		cnt++
		if cnt == 1 {
			return i.(int) * 2, nil
		}

		return nil, errors.New("huh oh")
	}).Times(2)

	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	in <- 1
	i1 := <-o

	in <- 2
	cancel()
	i2, ok := <-o

	assert.Equal(t, 2, i1)
	assert.Nil(t, i2)
	assert.False(t, ok)

	<-o
}

func TestMap_Start_StopsWhenInClosed(t *testing.T) {
	b := NewMap(newFakeMapStage(t))
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	close(in)
	_, ok := <-o

	assert.False(t, ok)
	cancel()
}

func TestMap_Start_StopsWhenCancel(t *testing.T) {
	b := NewMap(newFakeMapStage(t))
	ctx, cancel := context.WithCancel(context.Background())
	in := make(chan interface{})

	o := b.Start(ctx, in)

	cancel()
	_, ok := <-o

	assert.False(t, ok)
}

func newFakeMapStage(t *testing.T) *mocks.MockMapStage {
	ctl := gomock.NewController(t)

	return mocks.NewMockMapStage(ctl)
}
