package pipeline

import (
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

func LoadBytes(data []byte) (*Pipeline, error) {
	var p Pipeline
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("pipeline: yaml unmarshal: %w", err)
	}
	if err := validate(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func LoadFile(path string) (*Pipeline, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("pipeline: read file: %w", err)
	}
	return LoadBytes(data)
}

func validate(p *Pipeline) error {
	if len(p.Jobs) == 0 {
		return fmt.Errorf("pipeline: at least one job is required")
	}
	for name, job := range p.Jobs {
		if job == nil {
			return fmt.Errorf("pipeline: job %q is nil", name)
		}
		for _, need := range job.Needs {
			if _, ok := p.Jobs[need]; !ok {
				return fmt.Errorf("pipeline: job %q needs unknown job %q", name, need)
			}
		}
		if job.IsPerformanceGate() {
			if err := validatePerformanceGate(name, job, p.Jobs); err != nil {
				return err
			}
			continue
		}
		if len(job.Steps) == 0 {
			return fmt.Errorf("pipeline: job %q has no steps", name)
		}
	}
	return nil
}

func validatePerformanceGate(name string, job *Job, allJobs map[string]*Job) error {
	pg := job.PerformanceGate
	if pg.SourceJob == "" {
		return fmt.Errorf("pipeline: job %q performance_gate.source_job is required", name)
	}
	src, ok := allJobs[pg.SourceJob]
	if !ok {
		return fmt.Errorf("pipeline: job %q performance_gate.source_job %q not found", name, pg.SourceJob)
	}
	if src.IsPerformanceGate() {
		return fmt.Errorf("pipeline: job %q cannot use performance gate job %q as source", name, pg.SourceJob)
	}
	hasNeed := false
	for _, need := range job.Needs {
		if need == pg.SourceJob {
			hasNeed = true
			break
		}
	}
	if !hasNeed {
		return fmt.Errorf("pipeline: job %q must list performance_gate.source_job %q in needs", name, pg.SourceJob)
	}
	return nil
}
