//go:generate mockgen -package mocks -destination ../mocks/stage.go . Stage,MapStage,FilterStage,FlatMapStage,TumblingWindowStage

package contracts

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/repo"
)

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

type TumblingWindowStage interface {
	Process(window *repo.Window, item interface{})
	ToDataPoints(window int64, rawPoints *repo.Window) []interface{}
	HandleError(err error)
}
