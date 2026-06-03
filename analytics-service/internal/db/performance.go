package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type PerformanceMetricRecord struct {
	ProjectID  string
	RunID      string
	JobName    string
	Branch     string
	MetricName string
	Value      float64
	GatePassed bool
	CreatedAt  time.Time
}

type PerformanceGateResult struct {
	ProjectID       string
	RunID           string
	JobName         string
	Passed          bool
	ColdStart       bool
	BaselineSamples uint32
	Summary         string
	DetailsJSON     string
	CreatedAt       time.Time
}

func (r *Repository) migratePerformanceTables(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS performance_metrics (
			project_id   String,
			run_id       String,
			job_name     LowCardinality(String),
			branch       LowCardinality(String),
			metric_name  LowCardinality(String),
			metric_value Float64,
			gate_passed  UInt8 DEFAULT 0,
			created_at   DateTime
		)
		ENGINE = ReplacingMergeTree(created_at)
		PARTITION BY toYYYYMM(created_at)
		ORDER BY (project_id, run_id, job_name, metric_name)
		TTL created_at + INTERVAL 2 YEAR DELETE
		SETTINGS index_granularity = 8192`,

		`CREATE TABLE IF NOT EXISTS performance_gate_results (
			project_id        String,
			run_id            String,
			job_name          LowCardinality(String),
			passed            UInt8,
			cold_start        UInt8,
			baseline_samples  UInt32,
			summary           String,
			details_json      String,
			created_at        DateTime
		)
		ENGINE = ReplacingMergeTree(created_at)
		PARTITION BY toYYYYMM(created_at)
		ORDER BY (project_id, run_id, job_name)
		TTL created_at + INTERVAL 2 YEAR DELETE
		SETTINGS index_granularity = 8192`,
	}
	for _, stmt := range stmts {
		if err := r.conn.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("analytics migrate performance: %w", err)
		}
	}
	return nil
}

func (r *Repository) InsertPerformanceMetrics(ctx context.Context, records []*PerformanceMetricRecord) error {
	if len(records) == 0 {
		return nil
	}
	batch, err := r.conn.PrepareBatch(ctx,
		`INSERT INTO performance_metrics
		 (project_id, run_id, job_name, branch, metric_name, metric_value, gate_passed, created_at)`)
	if err != nil {
		return fmt.Errorf("analytics: prepare performance metrics batch: %w", err)
	}
	for _, rec := range records {
		createdAt := rec.CreatedAt
		if createdAt.IsZero() {
			createdAt = time.Now()
		}
		gatePassed := uint8(0)
		if rec.GatePassed {
			gatePassed = 1
		}
		if err := batch.Append(
			rec.ProjectID, rec.RunID, rec.JobName, rec.Branch,
			rec.MetricName, rec.Value, gatePassed, createdAt,
		); err != nil {
			return fmt.Errorf("analytics: append performance metric: %w", err)
		}
	}
	if err := batch.Send(); err != nil {
		return fmt.Errorf("analytics: send performance metrics batch: %w", err)
	}
	return nil
}

func (r *Repository) InsertPerformanceGateResult(ctx context.Context, result *PerformanceGateResult) error {
	if result.CreatedAt.IsZero() {
		result.CreatedAt = time.Now()
	}
	passed := uint8(0)
	if result.Passed {
		passed = 1
	}
	coldStart := uint8(0)
	if result.ColdStart {
		coldStart = 1
	}
	return r.conn.Exec(ctx,
		`INSERT INTO performance_gate_results
		 (project_id, run_id, job_name, passed, cold_start, baseline_samples, summary, details_json, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		result.ProjectID, result.RunID, result.JobName,
		passed, coldStart, result.BaselineSamples,
		result.Summary, result.DetailsJSON, result.CreatedAt,
	)
}

func (r *Repository) GetBaselineMetricValues(
	ctx context.Context,
	projectID, sourceJob, branch, excludeRunID string,
	metricNames []string,
	since time.Time,
	limit int,
) (map[string][]float64, error) {
	if limit <= 0 {
		limit = 50
	}
	out := make(map[string][]float64, len(metricNames))
	for _, name := range metricNames {
		rows, err := r.conn.Query(ctx, `
			SELECT metric_value
			FROM performance_metrics FINAL
			WHERE project_id = ?
			  AND job_name = ?
			  AND branch = ?
			  AND metric_name = ?
			  AND gate_passed = 1
			  AND run_id != ?
			  AND created_at >= ?
			ORDER BY created_at DESC
			LIMIT ?`,
			projectID, sourceJob, branch, name, excludeRunID, since, limit,
		)
		if err != nil {
			return nil, fmt.Errorf("analytics: baseline query for %s: %w", name, err)
		}
		var values []float64
		for rows.Next() {
			var v float64
			if err := rows.Scan(&v); err != nil {
				rows.Close()
				return nil, fmt.Errorf("analytics: scan baseline: %w", err)
			}
			values = append(values, v)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return nil, err
		}
		// Reverse to chronological order (oldest first).
		for i, j := 0, len(values)-1; i < j; i, j = i+1, j-1 {
			values[i], values[j] = values[j], values[i]
		}
		out[name] = values
	}
	return out, nil
}

func (r *Repository) MarkRunMetricsGatePassed(
	ctx context.Context,
	projectID, runID, jobName string,
	passed bool,
) error {
	rows, err := r.conn.Query(ctx, `
		SELECT branch, metric_name, metric_value, created_at
		FROM performance_metrics FINAL
		WHERE project_id = ? AND run_id = ? AND job_name = ?`,
		projectID, runID, jobName,
	)
	if err != nil {
		return fmt.Errorf("analytics: load run metrics for gate mark: %w", err)
	}
	defer rows.Close()

	gatePassed := uint8(0)
	if passed {
		gatePassed = 1
	}

	var records []struct {
		branch, metricName string
		value              float64
		createdAt          time.Time
	}
	for rows.Next() {
		var rec struct {
			branch, metricName string
			value              float64
			createdAt          time.Time
		}
		if err := rows.Scan(&rec.branch, &rec.metricName, &rec.value, &rec.createdAt); err != nil {
			return fmt.Errorf("analytics: scan run metric: %w", err)
		}
		records = append(records, rec)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(records) == 0 {
		return nil
	}

	now := time.Now()
	batch, err := r.conn.PrepareBatch(ctx,
		`INSERT INTO performance_metrics
		 (project_id, run_id, job_name, branch, metric_name, metric_value, gate_passed, created_at)`)
	if err != nil {
		return err
	}
	for _, rec := range records {
		if err := batch.Append(
			projectID, runID, jobName, rec.branch, rec.metricName, rec.value, gatePassed, now,
		); err != nil {
			return err
		}
	}
	return batch.Send()
}

func (r *Repository) GetRunPerformanceMetrics(
	ctx context.Context,
	projectID, runID, jobName string,
) (map[string]float64, error) {
	rows, err := r.conn.Query(ctx, `
		SELECT metric_name, metric_value
		FROM performance_metrics FINAL
		WHERE project_id = ? AND run_id = ? AND job_name = ?`,
		projectID, runID, jobName,
	)
	if err != nil {
		return nil, fmt.Errorf("analytics: get run metrics: %w", err)
	}
	defer rows.Close()

	out := make(map[string]float64)
	for rows.Next() {
		var name string
		var value float64
		if err := rows.Scan(&name, &value); err != nil {
			return nil, fmt.Errorf("analytics: scan run metric: %w", err)
		}
		out[name] = value
	}
	return out, rows.Err()
}

func (r *Repository) GetPerformanceGateResult(
	ctx context.Context,
	runID, jobName string,
) (*PerformanceGateResult, error) {
	row := r.conn.QueryRow(ctx, `
		SELECT project_id, run_id, job_name, passed, cold_start, baseline_samples, summary, details_json, created_at
		FROM performance_gate_results FINAL
		WHERE run_id = ? AND job_name = ?
		ORDER BY created_at DESC
		LIMIT 1`,
		runID, jobName,
	)

	res := &PerformanceGateResult{}
	var passed, coldStart uint8
	if err := row.Scan(
		&res.ProjectID, &res.RunID, &res.JobName,
		&passed, &coldStart, &res.BaselineSamples,
		&res.Summary, &res.DetailsJSON, &res.CreatedAt,
	); err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, nil
		}
		return nil, fmt.Errorf("analytics: get gate result: %w", err)
	}
	res.Passed = passed == 1
	res.ColdStart = coldStart == 1
	return res, nil
}

func EncodeGateDetails(details any) (string, error) {
	data, err := json.Marshal(details)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
