package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ipush/littlepipe/pkg/observability"
	"go.uber.org/zap"
)

// AsyncStageWrapper 将普通 Stage 包装成支持异步和并发处理的形式
type AsyncStageWrapper struct {
	stage      Stage
	workers    int
	bufferSize int
	logger     observability.Logger
	metrics    *observability.Metrics
}

func NewAsyncStageWrapper(stage Stage, workers, bufferSize int, logger observability.Logger, metrics *observability.Metrics) *AsyncStageWrapper {
	return &AsyncStageWrapper{
		stage:      stage,
		workers:    workers,
		bufferSize: bufferSize,
		logger:     logger,
		metrics:    metrics,
	}
}

// ProcessAsync 异步处理消息
func (w *AsyncStageWrapper) ProcessAsync(ctx context.Context) (chan<- *Message, <-chan *Message, <-chan error) {
	input := make(chan *Message, w.bufferSize)
	output := make(chan *Message, w.bufferSize)
	errChan := make(chan error, 1)

	// 启动工作池
	var wg sync.WaitGroup
	for i := 0; i < w.workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			w.worker(ctx, workerID, input, output, errChan)
		}(i)
	}

	// 监控协程，等待所有工作完成后关闭通道
	go func() {
		wg.Wait()
		close(output)
		close(errChan)
	}()

	return input, output, errChan
}

func (w *AsyncStageWrapper) worker(ctx context.Context, workerID int, input <-chan *Message, output chan<- *Message, errChan chan<- error) {
	for {
		select {
		case msg, ok := <-input:
			if !ok {
				return
			}

			startTime := time.Now()
			result, err := w.stage.Process(msg)
			duration := time.Since(startTime)

			// 记录指标
			w.metrics.ProcessingDuration.WithLabelValues(fmt.Sprintf("worker_%d", workerID)).Observe(duration.Seconds())

			if err != nil {
				w.metrics.ErrorsTotal.WithLabelValues(fmt.Sprintf("worker_%d", workerID), err.Error()).Inc()
				w.logger.Error("worker processing failed",
					zap.Int("worker_id", workerID),
					zap.String("message_id", msg.ID),
					zap.Error(err))

				// 发送错误但继续处理
				select {
				case errChan <- err:
				case <-ctx.Done():
					return
				}
				continue
			}

			// 发送处理结果
			select {
			case output <- result:
				w.metrics.MessagesTotal.WithLabelValues(fmt.Sprintf("worker_%d", workerID)).Inc()
			case <-ctx.Done():
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// ProcessBatch 批量异步处理
func (w *AsyncStageWrapper) ProcessBatch(ctx context.Context, batch []*Message) ([]*Message, error) {
	if len(batch) == 0 {
		return nil, nil
	}

	input, output, errChan := w.ProcessAsync(ctx)
	results := make([]*Message, 0, len(batch))
	errors := make([]error, 0)

	// 启动发送协程
	go func() {
		for _, msg := range batch {
			select {
			case input <- msg:
			case <-ctx.Done():
				close(input)
				return
			}
		}
		close(input)
	}()

	// 收集结果
	for i := 0; i < len(batch); i++ {
		select {
		case result := <-output:
			if result != nil {
				results = append(results, result)
			}
		case err := <-errChan:
			errors = append(errors, err)
		case <-ctx.Done():
			return results, ctx.Err()
		}
	}

	if len(errors) > 0 {
		// 可以根据需要决定如何处理错误
		return results, fmt.Errorf("batch processing had %d errors, first error: %v", len(errors), errors[0])
	}

	return results, nil
}
