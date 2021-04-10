package chans

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewOrDone(t *testing.T) {
	in := make(chan interface{})

	out := NewOrDone(context.Background(), in)

	in<-"hello"
	close(in)

	r := <-out
	_, ok := <-out

	assert.Equal(t, "hello", r)
	assert.False(t, ok)
}

func TestNewOrDone_Context_Done(t *testing.T) {
	in := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())

	out := NewOrDone(ctx, in)

	in<-"hello"
	cancel()

	r := <-out
	_, ok := <-out

	assert.Equal(t, "hello", r)
	assert.False(t, ok)
}
