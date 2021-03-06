//go:generate mockgen -package mocks -destination ../mocks/stage.go . Stage,MapStage,FlatMapStage

package contracts

import "context"

type Stage interface {
	Name() string
	Start(ctx context.Context, items <-chan interface{}) <-chan interface{}
}

type MapStage interface {
	Process(interface{}) (interface{}, error)
}

type FlatMapStage interface {
	Process(context.Context, interface{}) <-chan interface{}
}
