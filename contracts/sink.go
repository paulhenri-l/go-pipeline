//go:generate mockgen -package mocks -destination ../mocks/sink.go . Sink

package contracts

import "context"

type Sink interface {
	Start(ctx context.Context, in <-chan interface{}) <-chan struct{}
}
