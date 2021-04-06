//go:generate mockgen -package mocks -destination ../mocks/stage.go . Stage,MapStage,FilterStage,FlatMapStage

package contracts

import "context"

type Stage interface {
	Name() string
	Start(ctx context.Context, items <-chan interface{}) <-chan interface{}
}

type MapStage interface {
	Process(interface{}) (interface{}, error)
}

type FilterStage interface {
	Process(interface{}) (bool, error)
}

type FlatMapStage interface {
	Process(context.Context, interface{}) <-chan interface{}
}
