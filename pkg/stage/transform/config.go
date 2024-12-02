package transform

type TransformRule struct {
	Target string `json:"target"`
	Expr   string `json:"expr"`

	Required bool `json:"required"`
}

type TransformConfig struct {
	Rules []TransformRule `json:"rules"`
}
