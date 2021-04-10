package chans

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTee(t *testing.T) {
	in := make(chan interface{})

	outs := NewTee(in, 3)
	in <- "hello"
	close(in)

	r1 := <-outs[0]
	r2 := <-outs[1]
	r3 := <-outs[2]

	_, ok1 := <-outs[0]
	_, ok2 := <-outs[1]
	_, ok3 := <-outs[2]

	assert.Equal(t, "hello", r1)
	assert.Equal(t, "hello", r2)
	assert.Equal(t, "hello", r3)

	assert.False(t, ok1)
	assert.False(t, ok2)
	assert.False(t, ok3)
}
