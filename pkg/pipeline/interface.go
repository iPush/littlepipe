package pipeline

import (
	"time"

	"github.com/ipush/littlepipe/pkg/observability"
	"go.uber.org/zap"
)

type Closer interface {
	Close() error
}
type Source interface {
	Read() (*Message, error)
	BatchRead() ([]*Message, error)
	Closer
}

type Sink interface {
	Write(data *Message) error
	BatchWrite(data []*Message) error
	Closer
}

type Stage interface {
	Process(data *Message) (*Message, error)
	BatchProcess(data []*Message) ([]*Message, error)
	Closer
}

// ObservableStage 可观测的阶段
type ObservableStage struct {
	name    string
	stage   Stage
	logger  observability.Logger
	tracer  *observability.Tracer
	metrics *observability.Metrics
}

func NewObservableStage(name string, stage Stage, logger observability.Logger, tracer *observability.Tracer, metrics *observability.Metrics) *ObservableStage {
	return &ObservableStage{
		name:    name,
		stage:   stage,
		logger:  logger,
		tracer:  tracer,
		metrics: metrics,
	}
}

func (s *ObservableStage) Process(data *Message) (*Message, error) {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime)
		s.metrics.ProcessingDuration.WithLabelValues(s.name).Observe(duration.Seconds())
	}()

	s.metrics.MessagesInProgress.WithLabelValues(s.name).Inc()
	defer s.metrics.MessagesInProgress.WithLabelValues(s.name).Dec()

	s.logger.Info("processing message",
		zap.String("message_id", data.ID),
	)

	result, err := s.stage.Process(data)
	duration := time.Since(startTime)
	s.metrics.ProcessingDuration.WithLabelValues(s.name).Observe(float64(duration.Milliseconds()))

	if err != nil {
		s.metrics.ErrorsTotal.WithLabelValues(s.name, err.Error()).Inc()
		s.logger.Error("failed to process message",
			zap.String("message_id", data.ID),
			zap.Error(err),
			zap.Int("duration", int(duration.Milliseconds())))
	} else {
		s.metrics.MessagesTotal.WithLabelValues(s.name).Inc()
		s.logger.Info("processed message",
			zap.String("message_id", data.ID),
			zap.Int("duration", int(duration.Milliseconds())),
		)
	}

	return result, err
}
