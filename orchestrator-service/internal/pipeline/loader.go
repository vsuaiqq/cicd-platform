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
		if len(job.Steps) == 0 {
			return fmt.Errorf("pipeline: job %q has no steps", name)
		}
	}
	return nil
}
