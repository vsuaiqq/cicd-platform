package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const (
	StatusPending          = "pending"
	StatusRunning          = "running"
	StatusSuccess          = "success"
	StatusFailed           = "failed"
	StatusCancelled        = "cancelled"
	StatusSkipped          = "skipped"
	StatusAwaitingApproval = "awaiting_approval"
)

type PipelineRun struct {
	ID           string
	ProjectID    string
	CommitSHA    string
	Branch       string
	RepoURL      string
	TriggerType  string
	Status       string
	CreatedAt    time.Time
	StartedAt    *time.Time
	FinishedAt   *time.Time
	PipelineYAML string
}

type PipelineJob struct {
	ID          string
	RunID       string
	Name        string
	DisplayName string
	Status      string
	Attempt     int
	StartedAt   *time.Time
	FinishedAt  *time.Time
}

type PipelineStep struct {
	ID         string
	JobID      string
	StepIndex  int
	Name       string
	Status     string
	LogOutput  string
	ExitCode   int
	StartedAt  *time.Time
	FinishedAt *time.Time
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateRun(ctx context.Context, projectID, commitSHA, branch, repoURL string) (*PipelineRun, error) {
	run := &PipelineRun{
		ID:          uuid.New().String(),
		ProjectID:   projectID,
		CommitSHA:   commitSHA,
		Branch:      branch,
		RepoURL:     repoURL,
		TriggerType: "push",
		Status:      StatusPending,
		CreatedAt:   time.Now(),
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO pipeline_runs (id, project_id, commit_sha, branch, repo_url, trigger_type, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		run.ID, run.ProjectID, run.CommitSHA, run.Branch, run.RepoURL,
		run.TriggerType, run.Status, run.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("db: create run: %w", err)
	}
	return run, nil
}

func (r *Repository) UpdateRunStatus(ctx context.Context, runID, status string) error {
	now := time.Now()
	var q string
	var args []any
	switch status {
	case StatusRunning:
		q = `UPDATE pipeline_runs SET status=$1, started_at=$2 WHERE id=$3`
		args = []any{status, now, runID}
	case StatusSuccess, StatusFailed, StatusCancelled:
		q = `UPDATE pipeline_runs SET status=$1, finished_at=$2 WHERE id=$3`
		args = []any{status, now, runID}
	default:
		q = `UPDATE pipeline_runs SET status=$1 WHERE id=$2`
		args = []any{status, runID}
	}
	_, err := r.db.ExecContext(ctx, q, args...)
	return err
}

func (r *Repository) GetRun(ctx context.Context, runID string) (*PipelineRun, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, project_id, commit_sha, branch, COALESCE(repo_url,''),
		        trigger_type, status, created_at, started_at, finished_at,
		        COALESCE(pipeline_yaml,'')
		 FROM pipeline_runs WHERE id=$1`, runID)
	return scanRun(row)
}

func (r *Repository) SetPipelineYAML(ctx context.Context, runID, yaml string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE pipeline_runs SET pipeline_yaml=$1 WHERE id=$2`, yaml, runID)
	return err
}

func (r *Repository) ListRunsByProject(ctx context.Context, projectID string, limit int) ([]*PipelineRun, error) {
	if limit <= 0 {
		limit = 20
	}
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, project_id, commit_sha, branch, COALESCE(repo_url,''),
		        trigger_type, status, created_at, started_at, finished_at
		 FROM pipeline_runs WHERE project_id=$1
		 ORDER BY created_at DESC LIMIT $2`, projectID, limit)
	if err != nil {
		return nil, fmt.Errorf("db: list runs: %w", err)
	}
	defer rows.Close()
	return scanRuns(rows)
}

func (r *Repository) CreateJob(ctx context.Context, runID, name, displayName string) (*PipelineJob, error) {
	job := &PipelineJob{
		ID:          uuid.New().String(),
		RunID:       runID,
		Name:        name,
		DisplayName: displayName,
		Status:      StatusPending,
		Attempt:     1,
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO pipeline_jobs (id, run_id, name, display_name, status, attempt)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		job.ID, job.RunID, job.Name, job.DisplayName, job.Status, job.Attempt,
	)
	if err != nil {
		return nil, fmt.Errorf("db: create job: %w", err)
	}
	return job, nil
}

func (r *Repository) UpdateJobAttempt(ctx context.Context, jobID string, attempt int) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE pipeline_jobs SET attempt=$1 WHERE id=$2`, attempt, jobID)
	return err
}

func (r *Repository) UpdateJobStatus(ctx context.Context, jobID, status string) error {
	now := time.Now()
	var q string
	var args []any
	switch status {
	case StatusRunning:
		q = `UPDATE pipeline_jobs SET status=$1, started_at=$2 WHERE id=$3`
		args = []any{status, now, jobID}
	case StatusSuccess, StatusFailed, StatusSkipped, StatusCancelled:
		q = `UPDATE pipeline_jobs SET status=$1, finished_at=$2 WHERE id=$3`
		args = []any{status, now, jobID}
	default:
		q = `UPDATE pipeline_jobs SET status=$1 WHERE id=$2`
		args = []any{status, jobID}
	}
	_, err := r.db.ExecContext(ctx, q, args...)
	return err
}

func (r *Repository) GetJobByID(ctx context.Context, jobID string) (*PipelineJob, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, run_id, name, display_name, status, COALESCE(attempt,1), started_at, finished_at
		 FROM pipeline_jobs WHERE id=$1`, jobID)
	return scanJob(row)
}

func (r *Repository) CancelRunNonTerminalJobs(ctx context.Context, runID string) (runningJobIDs []string, err error) {

	rows, err := r.db.QueryContext(ctx,
		`SELECT id FROM pipeline_jobs WHERE run_id=$1 AND status=$2`, runID, StatusRunning)
	if err != nil {
		return nil, fmt.Errorf("db: get running jobs: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		runningJobIDs = append(runningJobIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(ctx,
		`UPDATE pipeline_jobs
		 SET status=$1, finished_at=NOW()
		 WHERE run_id=$2
		   AND status IN ('pending','running','awaiting_approval')`,
		StatusCancelled, runID,
	)
	return runningJobIDs, err
}

func (r *Repository) GetJobsByRun(ctx context.Context, runID string) ([]*PipelineJob, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, run_id, name, display_name, status, COALESCE(attempt,1), started_at, finished_at
		 FROM pipeline_jobs WHERE run_id=$1 ORDER BY started_at ASC NULLS FIRST`, runID)
	if err != nil {
		return nil, fmt.Errorf("db: get jobs: %w", err)
	}
	defer rows.Close()
	return scanJobs(rows)
}

func (r *Repository) GetJobByRunAndName(ctx context.Context, runID, name string) (*PipelineJob, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, run_id, name, display_name, status, started_at, finished_at
		 FROM pipeline_jobs WHERE run_id=$1 AND name=$2`, runID, name)
	return scanJob(row)
}

func (r *Repository) UpsertStep(ctx context.Context, step *PipelineStep) error {
	if step.ID == "" {
		step.ID = uuid.New().String()
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO pipeline_steps (id, job_id, step_index, name, status, log_output, exit_code, started_at, finished_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		 ON CONFLICT (job_id, step_index) DO UPDATE SET
		   id=$1, name=$4, status=$5, log_output=$6, exit_code=$7, started_at=$8, finished_at=$9`,
		step.ID, step.JobID, step.StepIndex, step.Name,
		step.Status, step.LogOutput, step.ExitCode, step.StartedAt, step.FinishedAt,
	)
	return err
}

func (r *Repository) DeleteStepsByJob(ctx context.Context, jobID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM pipeline_steps WHERE job_id=$1`, jobID)
	return err
}

func (r *Repository) GetStepsByJob(ctx context.Context, jobID string) ([]*PipelineStep, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, job_id, step_index, name, status, log_output, exit_code, started_at, finished_at
		 FROM pipeline_steps WHERE job_id=$1 ORDER BY step_index`, jobID)
	if err != nil {
		return nil, fmt.Errorf("db: get steps: %w", err)
	}
	defer rows.Close()
	return scanSteps(rows)
}

func scanRun(row *sql.Row) (*PipelineRun, error) {
	var r PipelineRun
	err := row.Scan(&r.ID, &r.ProjectID, &r.CommitSHA, &r.Branch, &r.RepoURL,
		&r.TriggerType, &r.Status, &r.CreatedAt, &r.StartedAt, &r.FinishedAt, &r.PipelineYAML)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &r, err
}

func scanRuns(rows *sql.Rows) ([]*PipelineRun, error) {
	var out []*PipelineRun
	for rows.Next() {
		var r PipelineRun
		if err := rows.Scan(&r.ID, &r.ProjectID, &r.CommitSHA, &r.Branch, &r.RepoURL,
			&r.TriggerType, &r.Status, &r.CreatedAt, &r.StartedAt, &r.FinishedAt); err != nil {
			return nil, err
		}
		out = append(out, &r)
	}
	return out, rows.Err()
}

func scanJob(row *sql.Row) (*PipelineJob, error) {
	var j PipelineJob
	err := row.Scan(&j.ID, &j.RunID, &j.Name, &j.DisplayName, &j.Status, &j.Attempt, &j.StartedAt, &j.FinishedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &j, err
}

func scanJobs(rows *sql.Rows) ([]*PipelineJob, error) {
	var out []*PipelineJob
	for rows.Next() {
		var j PipelineJob
		if err := rows.Scan(&j.ID, &j.RunID, &j.Name, &j.DisplayName, &j.Status, &j.Attempt, &j.StartedAt, &j.FinishedAt); err != nil {
			return nil, err
		}
		out = append(out, &j)
	}
	return out, rows.Err()
}

func scanSteps(rows *sql.Rows) ([]*PipelineStep, error) {
	var out []*PipelineStep
	for rows.Next() {
		var s PipelineStep
		if err := rows.Scan(&s.ID, &s.JobID, &s.StepIndex, &s.Name, &s.Status, &s.LogOutput, &s.ExitCode, &s.StartedAt, &s.FinishedAt); err != nil {
			return nil, err
		}
		out = append(out, &s)
	}
	return out, rows.Err()
}
