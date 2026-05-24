CREATE TABLE IF NOT EXISTS project_env_vars (
    project_id VARCHAR(36) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    key        VARCHAR(255) NOT NULL,
    value      TEXT         NOT NULL DEFAULT '',
    PRIMARY KEY (project_id, key)
);

CREATE INDEX IF NOT EXISTS idx_project_env_vars_project_id ON project_env_vars(project_id);
