package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

type Repository struct {
	conn driver.Conn
}

func NewRepository(conn driver.Conn) *Repository {
	return &Repository{conn: conn}
}

func (r *Repository) Migrate(ctx context.Context) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS analytics_run_events (
			project_id   String,
			run_id       String,
			status       LowCardinality(String),
			branch       LowCardinality(String),
			duration_sec UInt32,
			created_at   DateTime
		)
		ENGINE = ReplacingMergeTree(created_at)
		PARTITION BY toYYYYMM(created_at)
		ORDER BY (project_id, run_id)
		TTL created_at + INTERVAL 2 YEAR DELETE
		SETTINGS index_granularity = 8192`,

		`CREATE TABLE IF NOT EXISTS analytics_job_events (
			run_id       String,
			job_id       String,
			project_id   String,
			job_name     LowCardinality(String),
			status       LowCardinality(String),
			duration_sec UInt32,
			attempt      UInt8,
			created_at   DateTime
		)
		ENGINE = ReplacingMergeTree(created_at)
		PARTITION BY toYYYYMM(created_at)
		ORDER BY (project_id, job_id)
		TTL created_at + INTERVAL 2 YEAR DELETE
		SETTINGS index_granularity = 8192`,
	}
	for _, stmt := range stmts {
		if err := r.conn.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("analytics migrate: %w", err)
		}
	}
	return nil
}

type RunEvent struct {
	ProjectID   string
	RunID       string
	Status      string
	Branch      string
	DurationSec uint32
	CreatedAt   time.Time
}

func (r *Repository) InsertRunEvent(ctx context.Context, e *RunEvent) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	return r.conn.Exec(ctx,
		`INSERT INTO analytics_run_events
		 (project_id, run_id, status, branch, duration_sec, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		e.ProjectID, e.RunID, e.Status, e.Branch, e.DurationSec, e.CreatedAt,
	)
}

type JobEvent struct {
	RunID       string
	JobID       string
	ProjectID   string
	JobName     string
	Status      string
	DurationSec uint32
	Attempt     uint8
	CreatedAt   time.Time
}

func (r *Repository) InsertJobEvent(ctx context.Context, e *JobEvent) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	if e.Attempt < 1 {
		e.Attempt = 1
	}
	return r.conn.Exec(ctx,
		`INSERT INTO analytics_job_events
		 (run_id, job_id, project_id, job_name, status, duration_sec, attempt, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		e.RunID, e.JobID, e.ProjectID, e.JobName,
		e.Status, e.DurationSec, e.Attempt, e.CreatedAt,
	)
}

type RunStats struct {
	Total     int64
	Success   int64
	Failed    int64
	Cancelled int64
	AvgDur    float64
	P50Dur    float64
	P95Dur    float64
}

func (r *Repository) GetRunStats(ctx context.Context, projectID string, since time.Time) (*RunStats, error) {

	row := r.conn.QueryRow(ctx, `
		SELECT
			toInt64(count())                              AS total,
			toInt64(countIf(status = 'success'))          AS success,
			toInt64(countIf(status = 'failed'))           AS failed,
			toInt64(countIf(status = 'cancelled'))        AS cancelled,
			ifNull(avg(duration_sec), 0)                  AS avg_dur,
			ifNull(quantile(0.50)(duration_sec), 0)       AS p50,
			ifNull(quantile(0.95)(duration_sec), 0)       AS p95
		FROM analytics_run_events FINAL
		WHERE project_id = ? AND created_at >= ?`,
		projectID, since,
	)

	s := &RunStats{}
	if err := row.Scan(&s.Total, &s.Success, &s.Failed, &s.Cancelled, &s.AvgDur, &s.P50Dur, &s.P95Dur); err != nil {
		return nil, fmt.Errorf("analytics: get run stats: %w", err)
	}
	return s, nil
}

type DailyPoint struct {
	Date    string
	Total   int32
	Success int32
	Failed  int32
}

func (r *Repository) GetDailyTrend(ctx context.Context, projectID string, since time.Time) ([]*DailyPoint, error) {

	rows, err := r.conn.Query(ctx, `
		SELECT
			toDate(created_at)             AS day,
			toInt32(count())               AS total,
			toInt32(countIf(status = 'success')) AS success,
			toInt32(countIf(status = 'failed'))  AS failed
		FROM analytics_run_events FINAL
		WHERE project_id = ? AND created_at >= ?
		GROUP BY day
		ORDER BY day ASC
		WITH FILL FROM toDate(?) TO addDays(today(), 1) STEP 1`,
		projectID, since, since,
	)
	if err != nil {
		return nil, fmt.Errorf("analytics: get daily trend: %w", err)
	}
	defer rows.Close()

	var points []*DailyPoint
	for rows.Next() {
		var day time.Time
		p := &DailyPoint{}
		if err := rows.Scan(&day, &p.Total, &p.Success, &p.Failed); err != nil {
			return nil, fmt.Errorf("analytics: scan daily trend: %w", err)
		}
		p.Date = day.Format("2006-01-02")
		points = append(points, p)
	}
	return points, rows.Err()
}

type JobStat struct {
	JobName        string
	TotalRuns      int32
	FailureRate    float64
	AvgDurationSec float64
	AvgAttempts    float64
}

func (r *Repository) GetTopFailingJobs(ctx context.Context, projectID string, since time.Time, limit, minRuns int) ([]*JobStat, error) {
	return r.queryJobStats(ctx, `
		SELECT
			job_name,
			toInt32(count())                             AS total,
			avg(if(status = 'failed', 1.0, 0.0))         AS failure_rate,
			avg(duration_sec)                            AS avg_dur,
			avg(attempt)                                 AS avg_att
		FROM analytics_job_events FINAL
		WHERE project_id = ? AND created_at >= ? AND job_name != ''
		GROUP BY job_name
		HAVING total >= ?
		ORDER BY failure_rate DESC
		LIMIT ?`,
		projectID, since, minRuns, limit,
	)
}

func (r *Repository) GetTopSlowJobs(ctx context.Context, projectID string, since time.Time, limit, minRuns int) ([]*JobStat, error) {
	return r.queryJobStats(ctx, `
		SELECT
			job_name,
			toInt32(count())                             AS total,
			avg(if(status = 'failed', 1.0, 0.0))         AS failure_rate,
			avg(duration_sec)                            AS avg_dur,
			avg(attempt)                                 AS avg_att
		FROM analytics_job_events FINAL
		WHERE project_id = ? AND created_at >= ? AND job_name != ''
		GROUP BY job_name
		HAVING total >= ?
		ORDER BY avg_dur DESC
		LIMIT ?`,
		projectID, since, minRuns, limit,
	)
}

func (r *Repository) GetFlakyJobs(ctx context.Context, projectID string, since time.Time, limit, minRuns int) ([]*JobStat, error) {
	return r.queryJobStats(ctx, `
		SELECT
			job_name,
			toInt32(count())                             AS total,
			avg(if(status = 'failed', 1.0, 0.0))         AS failure_rate,
			avg(duration_sec)                            AS avg_dur,
			avg(attempt)                                 AS avg_att
		FROM analytics_job_events FINAL
		WHERE project_id = ? AND created_at >= ? AND job_name != ''
		GROUP BY job_name
		HAVING total >= ? AND avg_att > 1.0
		ORDER BY avg_att DESC
		LIMIT ?`,
		projectID, since, minRuns, limit,
	)
}

func (r *Repository) queryJobStats(ctx context.Context, query string, args ...any) ([]*JobStat, error) {
	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("analytics: job stats query: %w", err)
	}
	defer rows.Close()

	var stats []*JobStat
	for rows.Next() {
		s := &JobStat{}
		if err := rows.Scan(&s.JobName, &s.TotalRuns, &s.FailureRate, &s.AvgDurationSec, &s.AvgAttempts); err != nil {
			return nil, fmt.Errorf("analytics: scan job stats: %w", err)
		}
		stats = append(stats, s)
	}
	return stats, rows.Err()
}
