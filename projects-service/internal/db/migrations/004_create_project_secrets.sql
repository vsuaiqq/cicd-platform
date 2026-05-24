CREATE TABLE IF NOT EXISTS project_secrets (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id       VARCHAR(36) NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    key              TEXT NOT NULL,
    value_encrypted  TEXT NOT NULL,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (project_id, key)
);

CREATE INDEX IF NOT EXISTS idx_project_secrets_project_id ON project_secrets(project_id);
