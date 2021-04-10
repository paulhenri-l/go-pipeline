package pipeline

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/contracts"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewChain(t *testing.T) {
	var cnt1, cnt2 int64
	in := make(chan interface{})
	s1 := &counterStage{&cnt1}
	s2 := &counterStage{&cnt2}

	o := NewChain(
		context.Background(), in, []contracts.Stage{s1, s2},
	)

	in <- "hello"
	close(in)

	res := <-o
	_, ok := <-o

	assert.Equal(t, "hello", res)
	assert.False(t, ok)
	assert.Equal(t, int64(1), cnt1)
	assert.Equal(t, int64(1), cnt2)
}

func TestNewChain_Context(t *testing.T) {
	var cnt int64
	in := make(chan interface{})
	s1 := &counterStage{&cnt}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	o := NewChain(ctx, in, []contracts.Stage{s1})
	_, ok := <-o

	assert.False(t, ok)
	assert.Equal(t, int64(0), cnt)
}
