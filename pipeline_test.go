package pipeline

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/chans"
	"github.com/paulhenri-l/go-pipeline/contracts"
	"github.com/paulhenri-l/go-pipeline/sinks"
	"github.com/paulhenri-l/go-pipeline/stages"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewPipeline(t *testing.T) {
	p := New(
		&numGen{max: 1}, []contracts.Stage{}, &sinks.Nil{},
	)

	assert.IsType(t, &Pipeline{}, p)
}

func TestNewWithSingleStage(t *testing.T) {
	p := NewWithBluePrint(&numGen{max: 1}, newForwardBluePrint(), &sinks.Nil{})

	assert.IsType(t, &Pipeline{}, p)
}

func TestPipeline_Start(t *testing.T) {
	p := NewWithBluePrint(&numGen{max: 1}, newForwardBluePrint(), &sinks.Nil{})

	err := p.Start()

	assert.NoError(t, err)
	assert.Equal(t, true, p.Started())
}

func TestPipeline_Start_Twice(t *testing.T) {
	p := NewWithBluePrint(&numGen{max: 1}, newForwardBluePrint(), &sinks.Nil{})
	err1 := p.Start()
	err2 := p.Start()

	assert.Nil(t, err1)
	assert.EqualError(t, err2, "pipeline already started")
}

func TestPipeline_Start_AfterStop(t *testing.T) {
	p := NewWithBluePrint(&numGen{max: 1}, newForwardBluePrint(), &sinks.Nil{})
	_ = p.Start()
	p.Stop()

	err := p.Start()

	assert.EqualError(t, err, "pipeline stopped")
}

func TestPipeline_Stop(t *testing.T) {
	p := NewWithBluePrint(&numGen{max: 1}, newForwardBluePrint(), &sinks.Nil{})
	_ = p.Start()
	p.Stop()

	assert.True(t, p.Stopped())
}

func TestPipeline_Stop_DrainsPipeline(t *testing.T) {
	var cnt int64
	p := New(&numGen{max: 2}, []contracts.Stage{
		&waitStage{wait: 1 * time.Second},
		&counterStage{&cnt},
	}, &sinks.Nil{})

	_ = p.Start()
	_ = p.StopWithContext(context.Background())

	assert.Equal(t, int64(2), cnt)
}

func TestPipeline_Stop_Twice(t *testing.T) {
	p := NewWithBluePrint(&numGen{max: 1}, newForwardBluePrint(), &sinks.Nil{})

	_ = p.Start()
	e1 := p.StopWithContext(context.Background())
	e2 := p.StopWithContext(context.Background())

	assert.NoError(t, e1)
	assert.NoError(t, e2)
}

func TestPipeline_StopWithContext(t *testing.T) {
	var cnt int64
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	cancel()

	p := New(&numGen{max: 2}, []contracts.Stage{
		&waitStage{wait: 10 * time.Second},
		&counterStage{&cnt},
	}, &sinks.Nil{})

	_ = p.Start()
	err := p.StopWithContext(ctx)

	assert.Equal(t, int64(0), cnt)
	assert.Error(t, err)
}

func newForwardBluePrint() stages.BlueprintFunc {
	return func(ctx context.Context, in <-chan interface{}) <-chan interface{} {
		out := make(chan interface{})

		go func() {
			defer close(out)

			for i := range chans.NewOrDone(ctx, in) {
				out <- i
			}
		}()

		return out
	}
}
