package pipeline

import "sync"

type PipelineManager struct {
	pipelines map[string]*Pipeline
	mu        sync.RWMutex
}

func (m *PipelineManager) Create(config PipelineConfig) (string, error)
func (m *PipelineManager) Start(id string) error
func (m *PipelineManager) Stop(id string) error
func (m *PipelineManager) Delete(id string) error
func (m *PipelineManager) List() []PipelineInfo
