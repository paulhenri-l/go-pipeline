package pipeline

import (
	"context"
	"github.com/paulhenri-l/go-pipeline/contracts"
	"github.com/paulhenri-l/go-pipeline/stages"
	"github.com/pkg/errors"
	"sync"
)

type Pipeline struct {
	mtx             *sync.Mutex
	generator       contracts.Generator
	stages          []contracts.Stage
	sink            contracts.Sink
	done            <-chan struct{}
	generatorCancel context.CancelFunc
	stagesCancel    context.CancelFunc
	sinkCancel      context.CancelFunc
	started         bool
	shuttingDown    bool
	stopped         bool
}

func New(gen contracts.Generator, stages []contracts.Stage, sink contracts.Sink) *Pipeline {
	return &Pipeline{
		mtx:       &sync.Mutex{},
		generator: gen,
		stages:    stages,
		sink:      sink,
		started:   false,
	}
}

func NewWithBluePrint(gen contracts.Generator, bp stages.BlueprintFunc, sink contracts.Sink) *Pipeline {
	return New(gen, []contracts.Stage{stages.NewBuilder(bp)}, sink)
}

func (p *Pipeline) Start() error {
	if p.shuttingDown {
		return errors.New("pipeline stopped")
	}
	if p.started {
		return errors.New("pipeline already started")
	}

	p.started = true
	genOut := p.startGenerator()
	stagesOut := p.startStages(genOut)
	p.done = p.startSink(stagesOut)

	go func() {
		<-p.done
		p.mtx.Lock()
		p.stopped = true
		p.mtx.Unlock()
	}()

	return nil
}

func (p *Pipeline) Stop() {
	_ = p.StopWithContext(context.Background())
}

func (p *Pipeline) StopWithContext(ctx context.Context) error {
	var err error
	if p.shuttingDown {
		return nil
	}

	p.shuttingDown = true
	p.generatorCancel()

	select {
	case <-ctx.Done():
		p.stagesCancel()
		p.sinkCancel()
		err = errors.WithStack(ctx.Err())
	case <-p.done:
	}

	p.mtx.Lock()
	p.stopped = true
	p.mtx.Unlock()

	return err
}

func (p *Pipeline) Started() bool {
	return p.started
}

func (p *Pipeline) Stopped() bool {
	return p.stopped
}

func (p *Pipeline) startGenerator() <-chan interface{} {
	ctx, cancel := context.WithCancel(context.Background())
	p.generatorCancel = cancel

	return p.generator.Start(ctx)
}

func (p *Pipeline) startStages(in <-chan interface{}) <-chan interface{} {
	lastOutput := in
	ctx, cancel := context.WithCancel(context.Background())
	p.stagesCancel = cancel

	for _, stage := range p.stages {
		lastOutput = stage.Start(ctx, lastOutput)
	}

	return lastOutput
}

func (p *Pipeline) startSink(in <-chan interface{}) <-chan struct{} {
	ctx, cancel := context.WithCancel(context.Background())
	p.sinkCancel = cancel

	return p.sink.Start(ctx, in)
}
