CREATE TABLE IF NOT EXISTS pipeline_runs (
    id            VARCHAR(36) PRIMARY KEY,
    project_id    VARCHAR(36) NOT NULL,
    commit_sha    VARCHAR(40) NOT NULL,
    branch        VARCHAR(255) NOT NULL,
    trigger_type  VARCHAR(32) NOT NULL DEFAULT 'push',
    status        VARCHAR(32) NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at    TIMESTAMPTZ,
    finished_at   TIMESTAMPTZ,
    pipeline_yaml TEXT,
    repo_url      VARCHAR(2048)
);
ALTER TABLE pipeline_runs ADD COLUMN IF NOT EXISTS pipeline_yaml TEXT;
ALTER TABLE pipeline_runs ADD COLUMN IF NOT EXISTS repo_url      VARCHAR(2048);
CREATE INDEX IF NOT EXISTS idx_pipeline_runs_project_id ON pipeline_runs(project_id);
CREATE INDEX IF NOT EXISTS idx_pipeline_runs_status     ON pipeline_runs(status);

CREATE TABLE IF NOT EXISTS pipeline_jobs (
    id           VARCHAR(36) PRIMARY KEY,
    run_id       VARCHAR(36) NOT NULL REFERENCES pipeline_runs(id) ON DELETE CASCADE,
    name         VARCHAR(255) NOT NULL,
    display_name VARCHAR(255),
    status       VARCHAR(32) NOT NULL DEFAULT 'pending',
    attempt      INT NOT NULL DEFAULT 1,
    started_at   TIMESTAMPTZ,
    finished_at  TIMESTAMPTZ
);
ALTER TABLE pipeline_jobs ADD COLUMN IF NOT EXISTS attempt INT NOT NULL DEFAULT 1;
CREATE INDEX IF NOT EXISTS idx_pipeline_jobs_run_id ON pipeline_jobs(run_id);

CREATE TABLE IF NOT EXISTS pipeline_steps (
    id           VARCHAR(36) PRIMARY KEY,
    job_id       VARCHAR(36) NOT NULL REFERENCES pipeline_jobs(id) ON DELETE CASCADE,
    step_index   INT NOT NULL,
    name         VARCHAR(255) NOT NULL,
    status       VARCHAR(32) NOT NULL DEFAULT 'pending',
    log_output   TEXT,
    exit_code    INT,
    started_at   TIMESTAMPTZ,
    finished_at  TIMESTAMPTZ,
    UNIQUE (job_id, step_index)
);
DO $$ BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint
        WHERE conname = 'pipeline_steps_job_id_step_index_key'
    ) THEN
        ALTER TABLE pipeline_steps ADD CONSTRAINT pipeline_steps_job_id_step_index_key
            UNIQUE (job_id, step_index);
    END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_pipeline_steps_job_id ON pipeline_steps(job_id);
