package pipeline

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/contracts"
	"github.com/paulhenri-l/go-pipeline/sinks"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var counter int64 = 0

func TestNewPipeline(t *testing.T) {
	p := NewPipeline(
		&numGen{max: 1}, []contracts.Stage{}, &sinks.Nil{},
	)

	assert.IsType(t, &Pipeline{}, p)
}

func TestPipeline_Start(t *testing.T) {
	p := NewPipeline(
		&numGen{max: 1}, []contracts.Stage{}, &sinks.Nil{},
	)

	err := p.Start()

	assert.NoError(t, err)
	assert.Equal(t, true, p.Started())
}

func TestPipeline_Start_Twice(t *testing.T) {
	p := NewPipeline(
		&numGen{max: 1}, []contracts.Stage{}, &sinks.Nil{},
	)
	err1 := p.Start()
	err2 := p.Start()

	assert.Nil(t, err1)
	assert.EqualError(t, err2, "pipeline already started")
}

func TestPipeline_Start_AfterStop(t *testing.T) {
	p := NewPipeline(
		&numGen{max: 1}, []contracts.Stage{}, &sinks.Nil{},
	)
	_ = p.Start()
	p.Stop()

	err := p.Start()

	assert.EqualError(t, err, "pipeline stopped")
}

func TestPipeline_Stop(t *testing.T) {
	p := NewPipeline(
		&numGen{max: 1}, []contracts.Stage{}, &sinks.Nil{},
	)
	_ = p.Start()
	p.Stop()

	assert.True(t, p.Stopped())
}

func TestPipeline_Stop_DrainsPipeline(t *testing.T) {
	p := NewPipeline(&numGen{max: 2}, []contracts.Stage{
		&waitStage{wait: 1 * time.Second},
		&counterStage{},
	}, &sinks.Nil{})

	_ = p.Start()
	_ = p.StopWithContext(context.Background())

	assert.Equal(t, int64(2), counter)
}

func TestPipeline_Stop_Twice(t *testing.T) {
	p := NewPipeline(
		&numGen{max: 1}, []contracts.Stage{}, &sinks.Nil{},
	)

	_ = p.Start()
	e1 := p.StopWithContext(context.Background())
	e2 := p.StopWithContext(context.Background())

	assert.NoError(t, e1)
	assert.NoError(t, e2)
}

func TestPipeline_StopWithContext(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	cancel()

	p := NewPipeline(&numGen{max: 2}, []contracts.Stage{
		&waitStage{wait: 10 * time.Second},
		&counterStage{},
	}, &sinks.Nil{})

	_ = p.Start()
	err := p.StopWithContext(ctx)

	assert.Equal(t, int64(0), counter)
	assert.Error(t, err)
}

type numGen struct {
	max int
}

func (n *numGen) Start(ctx context.Context) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)

		for i := 0; i < n.max; i++ {
			out <- i
		}

		<-ctx.Done()
	}()

	return out
}

type counterStage struct{}

func (c counterStage) Name() string {
	return "counterstage"
}

func (c *counterStage) Start(ctx context.Context, events <-chan interface{}) <-chan interface{} {
	counter = 0
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return

			case e, ok := <-events:
				if !ok {
					return
				}
				atomic.AddInt64(&counter, 1)
				out <- e
			}
		}
	}()

	return out
}

type waitStage struct {
	wait time.Duration
}

func (ws *waitStage) Name() string {
	return "waitstage"
}

func (ws *waitStage) Start(ctx context.Context, events <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return

			case e, ok := <-events:
				if !ok {
					return
				}
				<-time.NewTimer(ws.wait).C
				out <- e
			}
		}
	}()

	return out
}
