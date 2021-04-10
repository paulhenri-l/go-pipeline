package pipeline

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/contracts"
)

func NewChain(ctx context.Context, in <-chan interface{}, stages []contracts.Stage) <-chan interface{} {
	previousOut := in

	for _, stage := range stages {
		previousOut = stage.Start(ctx, previousOut)
	}

	return previousOut
}
