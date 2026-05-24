package pipeline

import (
	"fmt"
	"time"

	"go.yaml.in/yaml/v3"
)

type Pipeline struct {
	Name string            `yaml:"name"`
	On   *Trigger          `yaml:"on"`
	Env  map[string]string `yaml:"env"`
	Jobs map[string]*Job   `yaml:"jobs"`
}

type Trigger struct {
	Push *PushTrigger `yaml:"push"`
}

type PushTrigger struct {
	Branches []string `yaml:"branches"`
}

type ApprovalRequired bool

func (a *ApprovalRequired) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	case "required", "true", "yes", "1":
		*a = true
	default:
		*a = false
	}
	return nil
}

type Duration time.Duration

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	dur, err := time.ParseDuration(value.Value)
	if err != nil {
		return fmt.Errorf("pipeline: invalid duration %q: %w", value.Value, err)
	}
	*d = Duration(dur)
	return nil
}

func (d Duration) Seconds() int { return int(time.Duration(d).Seconds()) }

type RetryConfig struct {
	Max int `yaml:"max,omitempty"`
}

func (r *RetryConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		return value.Decode(&r.Max)
	}
	type plain RetryConfig
	return value.Decode((*plain)(r))
}

type CacheConfig struct {
	Key   string   `yaml:"key"`
	Paths []string `yaml:"paths"`
}

type ArtifactDownload struct {
	Job string `yaml:"job"`
}

type ArtifactsConfig struct {
	Paths    []string           `yaml:"paths,omitempty"`
	Download []ArtifactDownload `yaml:"download,omitempty"`
}

type Job struct {
	Name  string            `yaml:"name"`
	Image string            `yaml:"image,omitempty"`
	Needs []string          `yaml:"needs"`
	Env   map[string]string `yaml:"env"`
	Steps []*Step           `yaml:"steps"`

	Approval ApprovalRequired `yaml:"approval,omitempty"`

	Timeout Duration `yaml:"timeout,omitempty"`

	Retry RetryConfig `yaml:"retry,omitempty"`

	Cache *CacheConfig `yaml:"cache,omitempty"`

	Artifacts *ArtifactsConfig `yaml:"artifacts,omitempty"`
}

type Step struct {
	Name            string `yaml:"name"`
	Run             string `yaml:"run,omitempty"`
	ContinueOnError bool   `yaml:"continue_on_error,omitempty"`

	Timeout Duration `yaml:"timeout,omitempty"`

	Retry int `yaml:"retry,omitempty"`
}

func (p *Pipeline) BranchAllowed(branch string) bool {
	if p.On == nil || p.On.Push == nil || len(p.On.Push.Branches) == 0 {
		return true
	}
	for _, b := range p.On.Push.Branches {
		if b == branch || b == "*" {
			return true
		}
	}
	return false
}

func (p *Pipeline) MergedEnv(job *Job, projectEnv map[string]string) map[string]string {
	size := len(p.Env) + len(job.Env) + len(projectEnv)
	merged := make(map[string]string, size)
	for k, v := range p.Env {
		merged[k] = v
	}
	for k, v := range job.Env {
		merged[k] = v
	}

	for k, v := range projectEnv {
		merged[k] = v
	}
	return merged
}
