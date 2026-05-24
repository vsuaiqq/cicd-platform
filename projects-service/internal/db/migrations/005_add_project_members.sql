CREATE TABLE IF NOT EXISTS project_members (
  id           VARCHAR(36)  NOT NULL PRIMARY KEY,
  project_id   VARCHAR(36)  NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
  user_id      VARCHAR(36)  NOT NULL,
  email        VARCHAR(255) NOT NULL,
  display_name VARCHAR(255) NOT NULL DEFAULT '',
  role         VARCHAR(32)  NOT NULL DEFAULT 'viewer',
  invited_by   VARCHAR(36)  NOT NULL,
  created_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
  UNIQUE (project_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_project_members_project_id ON project_members(project_id);
CREATE INDEX IF NOT EXISTS idx_project_members_user_id    ON project_members(user_id);
