package generators

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewChan(t *testing.T) {
	in := make(chan interface{})
	c := NewChan(in)

	assert.IsType(t, &Chan{}, c)
}

func TestChan_Start(t *testing.T) {
	in := make(chan interface{})
	c := NewChan(in)
	ctx, cancel := context.WithCancel(context.Background())

	out := c.Start(ctx)

	in <- 1
	received := <-out
	cancel()

	<-out
	assert.Equal(t, 1, received)
}

func TestChan_Start_StopsOnClose(t *testing.T) {
	in := make(chan interface{})
	c := NewChan(in)

	out := c.Start(context.Background())
	close(in)

	_, ok := <-out
	assert.False(t, ok)
}

func TestChan_Start_StopsOnCancel(t *testing.T) {
	in := make(chan interface{})
	c := NewChan(in)
	ctx, cancel := context.WithCancel(context.Background())

	out := c.Start(ctx)
	cancel()

	_, ok := <-out
	assert.False(t, ok)
}