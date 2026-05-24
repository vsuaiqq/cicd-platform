ALTER TABLE projects
    ADD COLUMN IF NOT EXISTS pipeline_yaml_override TEXT NOT NULL DEFAULT '';
