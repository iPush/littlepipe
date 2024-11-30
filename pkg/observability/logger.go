package observability

import "go.uber.org/zap"

type Logger interface {
	Info(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	With(fields ...Field) Logger
}
type Field = zap.Field

func NewLogger(opts ...zap.Option) Logger {
	logger, _ := zap.NewProduction()
	return &zapLogger{logger: logger}
}

type zapLogger struct {
	logger *zap.Logger
}

func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{logger: l.logger.With(fields...)}
}
