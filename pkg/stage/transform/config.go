package transform

import "github.com/ipush/littlepipe/pkg/pipeline"

type TransformRule struct {
	Target   string             `json:"target"`
	Expr     string             `json:"expr"`
	Type     pipeline.FieldType `json:"type"`
	Required bool               `json:"required"`
}

type TransformConfig struct {
	Rules []TransformRule `json:"rules"`
}
