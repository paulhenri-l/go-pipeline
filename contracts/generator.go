//go:generate mockgen -package mocks -destination ../mocks/generator.go . Generator

package contracts

import "context"

type Generator interface {
	Start(ctx context.Context, done <-chan struct{}) <-chan interface{}
}
