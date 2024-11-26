package littlepipe

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

type Config struct {
	BufferSize  int
	Concurrency int
	RetryCount  int
	RetryDelay  time.Duration
}

type LittlePipe struct {
	source  Source
	stages  []Stage
	sink    Sink
	errChan chan error
	config  Config
	ctx     context.Context
	cancel  context.CancelFunc
}

func NewLittlePipe(config Config) *LittlePipe {
	ctx, cancel := context.WithCancel(context.Background())
	return &LittlePipe{
		config:  config,
		ctx:     ctx,
		cancel:  cancel,
		errChan: make(chan error),
	}
}

func (p *LittlePipe) SetSource(source Source) *LittlePipe {
	p.source = source
	return p
}

func (p *LittlePipe) AddStage(stage Stage) *LittlePipe {
	p.stages = append(p.stages, stage)
	return p
}

func (p *LittlePipe) SetSink(sink Sink) *LittlePipe {
	p.sink = sink
	return p
}

func (p *LittlePipe) Run() error {
	if p.source == nil || p.sink == nil {
		return fmt.Errorf("source and sink are required")
	}

	channels := make([]chan any, len(p.stages)+1)
	for i := range channels {
		channels[i] = make(chan any, p.config.BufferSize)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(p.stages)+2)

	// start Source
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(channels[0])

		for {
			data, err := p.source.Read()
			if err != nil {
				if err == io.EOF {
					return
				}
				errChan <- err
				return
			}
			select {
			case channels[0] <- data:
			case <-p.ctx.Done():
				return
			}
		}
	}()

	// start Stages
	for i, stage := range p.stages {
		wg.Add(1)
		go func(index int, stage Stage) {
			defer wg.Done()
			defer close(channels[index+1])

			for data := range channels[index] {
				result, err := stage.Process(data)
				if err != nil {
					errChan <- fmt.Errorf("stage %d: %w", index, err)
					return
				}
				select {
				case channels[index+1] <- result:
				case <-p.ctx.Done():
				}
			}
		}(i, stage)
	}

	// start Sink
	wg.Add(1)
	go func() {
		defer wg.Done()
		for data := range channels[len(channels)-1] {
			if err := p.sink.Write(data); err != nil {
				errChan <- fmt.Errorf("sink: %w", err)
				return
			}
		}
	}()

	// wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()

	select {
	case err := <-errChan:
		p.cancel()
		return err
	case <-p.ctx.Done():
		return p.ctx.Err()
	case <-waitChanClosed(errChan):
		return nil
	}
}

func waitChanClosed(ch <-chan error) chan struct{} {
	done := make(chan struct{})
	go func() {
		for range ch {
			// 消耗所有错误
		}
		close(done)
	}()
	return done
}

type Source interface {
	Read() (interface{}, error)
}

type Stage interface {
	Process(data any) (any, error)
}

type Sink interface {
	Write(data any) error
}
