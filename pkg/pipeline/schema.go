package pipeline

import (
	"fmt"
	"time"
)

type Record struct {
	Schema    *Schema
	Data      map[string]Value
	Timestamp time.Time
	Version   int
}

type Schema struct {
	Fields     []Field
	PrimaryKey string
	Version    int
	Metadata   map[string]string
}

type Field struct {
	Name     string
	Type     FieldType
	Required bool
}

type Value struct {
	Type  FieldType
	Value any
}

type FieldType int

const (
	TypeUnknown FieldType = iota
	TypeString
	TypeInt64
	TypeFloat64
	TypeBoolean
	TypeList
	TypeDict
	TypeDecimal
	TypeJSON
)

func (r *Record) GetValue(name string) (Value, bool) {
	value, ok := r.Data[name]
	return value, ok
}

func (s *Schema) Validate(record *Record) error {
	for _, field := range s.Fields {
		value, ok := record.Data[field.Name]
		if !ok && field.Required {
			return fmt.Errorf("field %s is required but not present", field.Name)
		}

		if ok && value.Type != field.Type {
			return fmt.Errorf("field %s type mismatch, want %s, get %s", field.Name,
				field.Type, value.Type)
		}
	}
	return nil
}
