package words

import (
	"fmt"
	"strings"
)

// UppercaseStage 将文本转换为大写
type UppercaseStage struct{}

func NewUppercaseStage() *UppercaseStage {
	return &UppercaseStage{}
}

func (s *UppercaseStage) Process(data interface{}) (interface{}, error) {
	if str, ok := data.(string); ok {
		return strings.ToUpper(str), nil
	}
	return nil, fmt.Errorf("UppercaseStage: unsupported data type: %T", data)
}

// FilterEmptyStage 过滤空行
type FilterEmptyStage struct{}

func NewFilterEmptyStage() *FilterEmptyStage {
	return &FilterEmptyStage{}
}

func (s *FilterEmptyStage) Process(data interface{}) (interface{}, error) {
	if str, ok := data.(string); ok {
		if strings.TrimSpace(str) == "" {
			return nil, nil // 返回 nil 表示该数据应被过滤掉
		}
		return str, nil
	}
	return nil, fmt.Errorf("FilterEmptyStage: unsupported data type: %T", data)
}
