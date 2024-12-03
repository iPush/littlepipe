package pipeline

import (
	"encoding/json"
	"time"
)

// ComponentType 组件类型
type ComponentType string

const (
	// Source 类型
	SourceTypeMySQL ComponentType = "mysql"
	SourceTypeKafka ComponentType = "kafka"

	// Stage 类型
	StageTypeTransform ComponentType = "transform"
	StageTypeFilter    ComponentType = "filter"

	// Sink 类型
	SinkTypeKafka ComponentType = "kafka"
	SinkTypeES    ComponentType = "elasticsearch"
)

// Status Pipeline 状态
type Status string

const (
	StatusCreated Status = "created"
	StatusRunning Status = "running"
	StatusStopped Status = "stopped"
	StatusError   Status = "error"
)

// PipelineConfig 完整的 Pipeline 配置
type PipelineConfig struct {
	// 基础信息
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 组件配置
	Source SourceConfig  `json:"source"`
	Stages []StageConfig `json:"stages"`
	Sink   SinkConfig    `json:"sink"`

	// 运行时配置
	BufferSize  int           `json:"buffer_size"`
	Concurrency int           `json:"concurrency"`
	RetryCount  int           `json:"retry_count"`
	RetryDelay  time.Duration `json:"retry_delay"`

	// 监控配置
	EnableMetrics bool `json:"enable_metrics"`
	EnableTracing bool `json:"enable_tracing"`
	EnableLogging bool `json:"enable_logging"`
}

// SourceConfig Source 配置
type SourceConfig struct {
	Type   ComponentType   `json:"type"`
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"` // 具体配置使用 RawMessage
}

// MySQL Source 具体配置
type MySQLSourceConfig struct {
	DSN         string        `json:"dsn"`
	Query       string        `json:"query"`
	BatchSize   int           `json:"batch_size"`
	PollTimeout time.Duration `json:"poll_timeout"`
	// 增量配置
	IncrementColumn string      `json:"increment_column"`
	LastValue       interface{} `json:"last_value"`
}

// StageConfig Stage 配置
type StageConfig struct {
	Type   ComponentType   `json:"type"`
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

// Transform Stage 具体配置
type TransformConfig struct {
	Rules []TransformRule `json:"rules"`
}

type TransformRule struct {
	Target   string             `json:"target"`
	Expr     string             `json:"expr"`
	Type     pipeline.FieldType `json:"type"`
	Required bool               `json:"required"`
}

// SinkConfig Sink 配置
type SinkConfig struct {
	Type   ComponentType   `json:"type"`
	Name   string          `json:"name"`
	Config json.RawMessage `json:"config"`
}

// Kafka Sink 具体配置
type KafkaSinkConfig struct {
	Brokers []string `json:"brokers"`
	Topic   string   `json:"topic"`
	// 生产者配置
	BatchSize    int           `json:"batch_size"`
	BatchTimeout time.Duration `json:"batch_timeout"`
	RetryMax     int           `json:"retry_max"`
	RequiredAcks int           `json:"required_acks"`
}
