package chans

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFanIn(t *testing.T) {
	var results = make([]interface{}, 3)
	in1 := make(chan interface{}, 2)
	in2 := make(chan interface{}, 2)

	out := NewFanIn(in1, in2)
	in1<-"hello1"
	in2<-"hello2"
	close(in1)
	in2<-"still running"
	close(in2)

	results[0] = <-out
	results[1] = <-out
	results[2] = <-out
	_, ok := <-out

	assert.Contains(t, results, "hello1")
	assert.Contains(t, results, "hello2")
	assert.Contains(t, results, "still running")
	assert.False(t, ok)
}
