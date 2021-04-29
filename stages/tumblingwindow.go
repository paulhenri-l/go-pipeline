package stages

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/contracts"
	"github.com/paulhenri-l/go-pipeline/repo"
	"github.com/pkg/errors"
	"math"
	"time"
)

type TumblingWindow struct {
	stage            contracts.TumblingWindowStage
	windowStore      *repo.WindowStore
	lateWindowsCount int
	windowSize       time.Duration
}

func NewTumblingWindow(s contracts.TumblingWindowStage, windowSize time.Duration, lateWindowsCount int) *TumblingWindow {
	return &TumblingWindow{
		stage:            s,
		windowStore:      repo.NewWindowStore(),
		windowSize:       windowSize,
		lateWindowsCount: lateWindowsCount,
	}
}

func (t *TumblingWindow) Name() string {
	return "Tumbling window"
}

func (t *TumblingWindow) Start(ctx context.Context, items <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	emmitInterval := time.NewTicker(time.Second * 1)

	go func() {
		defer emmitInterval.Stop()
		defer close(out)

		for {
			select {
			case <-ctx.Done():
				t.flushAll(out)
				return

			case <-emmitInterval.C:
				t.emmitOldWindows(out)
				emmitInterval.Reset(time.Second * 1)

			case i, ok := <-items:
				if !ok {
					t.flushAll(out)
					return
				}

				t.process(i)
			}
		}
	}()

	return out
}

func (t *TumblingWindow) flushAll(out chan<- interface{}) {
	for windowDob, points := range t.windowStore.PopAll() {
		for _, p := range t.stage.ToDataPoints(windowDob, points) {
			out <- p
		}
	}
}

func (t *TumblingWindow) emmitOldWindows(out chan<- interface{}) {
	maxWindow := time.Now().Add(t.windowSize * -1 * time.Duration(t.lateWindowsCount)).Unix()

	for windowDob, points := range t.windowStore.PopOlderThan(maxWindow) {
		for _, p := range t.stage.ToDataPoints(windowDob, points) {
			out <- p
		}
	}
}

func (t *TumblingWindow) process(item interface{}) {
	event, ok := item.(contracts.Timestamped)
	if !ok {
		t.stage.HandleError(
			errors.Errorf("item does not implement Timestamped interface %+v", item),
		)
		return
	}

	bin := int64(math.Floor(float64(event.GetTimestamp())/t.windowSize.Seconds()) * t.windowSize.Seconds())
	w := t.windowStore.Get(bin)
	t.stage.Process(w, item)
}
