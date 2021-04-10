package chans

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFanOut(t *testing.T) {
	in := make(chan interface{})

	outs := NewFanOut(in, 2)
	in <- "hello"
	in <- "hello"
	close(in)

	r1 := <-outs[0]
	r2 := <-outs[1]
	_, ok1 := <-outs[0]
	_, ok2 := <-outs[1]

	assert.Equal(t, "hello", r1)
	assert.Equal(t, "hello", r2)
	assert.False(t, ok1)
	assert.False(t, ok2)
}

