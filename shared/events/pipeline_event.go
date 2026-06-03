package events

type PipelineJobTask struct {
	RunID     string            `json:"run_id"`
	ProjectID string            `json:"project_id,omitempty"`
	JobID     string            `json:"job_id"`
	JobName   string            `json:"job_name"`
	Branch    string            `json:"branch,omitempty"`
	Image     string            `json:"image,omitempty"`
	RepoURL   string            `json:"repo_url"`
	CommitSHA string            `json:"commit_sha"`
	SSHKey    string            `json:"ssh_key"`
	Env       map[string]string `json:"env,omitempty"`

	Secrets map[string]string `json:"secrets,omitempty"`

	TimeoutSeconds int `json:"timeout_seconds,omitempty"`
	RetryMax       int `json:"retry_max,omitempty"`
	AttemptNumber  int `json:"attempt_number,omitempty"`

	Cache []CacheEntry `json:"cache,omitempty"`

	ArtifactsUpload []string `json:"artifacts_upload,omitempty"`

	ArtifactsDownload []ArtifactSource `json:"artifacts_download,omitempty"`

	Steps []PipelineStep `json:"steps"`
}

type CacheEntry struct {
	Key   string   `json:"key"`
	Paths []string `json:"paths"`
}

type ArtifactSource struct {
	JobName string `json:"job_name"`
	JobID   string `json:"job_id"`
}

type PipelineStep struct {
	Index           int    `json:"index"`
	Name            string `json:"name"`
	Run             string `json:"run"`
	ContinueOnError bool   `json:"continue_on_error,omitempty"`
	TimeoutSeconds  int    `json:"timeout_seconds,omitempty"`
	RetryMax        int    `json:"retry_max,omitempty"`
}

type PerformanceMetricValue struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type PipelineJobResult struct {
	RunID         string       `json:"run_id"`
	JobID         string       `json:"job_id"`
	JobName       string       `json:"job_name,omitempty"`
	ProjectID     string       `json:"project_id,omitempty"`
	Branch        string       `json:"branch,omitempty"`
	Status        string       `json:"status"`
	AttemptNumber int          `json:"attempt_number,omitempty"`
	Steps         []StepResult `json:"steps"`
	StartedAt     int64        `json:"started_at"`
	FinishedAt    int64        `json:"finished_at"`

	PerformanceMetrics []PerformanceMetricValue `json:"performance_metrics,omitempty"`
}

type CancelJobEvent struct {
	RunID string `json:"run_id"`
	JobID string `json:"job_id,omitempty"`
}

type RunCompletedEvent struct {
	RunID       string `json:"run_id"`
	ProjectID   string `json:"project_id"`
	Status      string `json:"status"`
	Branch      string `json:"branch"`
	CommitSHA   string `json:"commit_sha"`
	DurationSec int    `json:"duration_sec"`
	StartedAt   int64  `json:"started_at"`
	FinishedAt  int64  `json:"finished_at"`
}

type StepResult struct {
	Index      int    `json:"index"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	Output     string `json:"output"`
	ExitCode   int    `json:"exit_code"`
	StartedAt  int64  `json:"started_at"`
	FinishedAt int64  `json:"finished_at"`
}
