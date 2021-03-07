package generators

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNil(t *testing.T) {
	c := NewNil()

	assert.IsType(t, &Nil{}, c)
}

func TestNil_Start_StopsOnCancel(t *testing.T) {
	c := NewNil()
	ctx, cancel := context.WithCancel(context.Background())

	out := c.Start(ctx)
	cancel()

	_, ok := <-out
	assert.False(t, ok)
}