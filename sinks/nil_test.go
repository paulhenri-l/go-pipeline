package sinks

import (
	"context"
	"testing"
)

func TestNil_Start(t *testing.T) {
	n := &Nil{}
	in := make(chan interface{})

	done := n.Start(context.Background(), in)

	in<-1
	in<-1
	in<-1
	close(in)

	<-done
}

func TestNil_Start_CtxCancel(t *testing.T) {
	n := &Nil{}
	in := make(chan interface{})
	ctx, cancel := context.WithCancel(context.Background())

	done := n.Start(ctx, in)
	cancel()

	<-done
}
