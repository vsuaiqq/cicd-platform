package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/vsuaiqq/cicd/projects-service/internal/models"
)

var ErrProjectNotFound = errors.New("project not found")

const projectColumns = `id, user_id, name, repo_url, default_branch, private_key_encrypted, public_key, webhook_secret, status, created_at, updated_at`

type rowScanner interface {
	Scan(dest ...any) error
}

func scanProject(s rowScanner) (*models.Project, error) {
	p := &models.Project{}
	err := s.Scan(
		&p.ID, &p.UserID, &p.Name, &p.RepoURL, &p.DefaultBranch,
		&p.PrivateKeyEncrypted, &p.PublicKey, &p.WebhookSecret,
		&p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return p, nil
}

type ProjectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, p *models.Project) error {
	const query = `
		INSERT INTO projects
			(id, user_id, name, repo_url, default_branch, private_key_encrypted, public_key, webhook_secret, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		p.ID, p.UserID, p.Name, p.RepoURL, p.DefaultBranch,
		p.PrivateKeyEncrypted, p.PublicKey, p.WebhookSecret, p.Status,
	).Scan(&p.CreatedAt, &p.UpdatedAt)
}

func (r *ProjectRepository) GetByID(ctx context.Context, id string) (*models.Project, error) {
	const query = `SELECT ` + projectColumns + ` FROM projects WHERE id = $1`
	p, err := scanProject(r.db.QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProjectNotFound
	}
	return p, err
}

func (r *ProjectRepository) GetByIDAndUser(ctx context.Context, id, userID string) (*models.Project, error) {
	const query = `SELECT ` + projectColumns + ` FROM projects WHERE id = $1 AND user_id = $2`
	p, err := scanProject(r.db.QueryRowContext(ctx, query, id, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProjectNotFound
	}
	return p, err
}

func (r *ProjectRepository) ListByUser(ctx context.Context, userID string) ([]*models.Project, error) {
	const query = `
		SELECT id, name, repo_url, default_branch, status
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Project
	for rows.Next() {
		p := &models.Project{UserID: userID}
		if err := rows.Scan(&p.ID, &p.Name, &p.RepoURL, &p.DefaultBranch, &p.Status); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *ProjectRepository) UpdateStatus(ctx context.Context, id, userID, status string) error {
	const query = `UPDATE projects SET status = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`
	res, err := r.db.ExecContext(ctx, query, status, id, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrProjectNotFound
	}
	return nil
}

func (r *ProjectRepository) DeleteByIDAndUser(ctx context.Context, id, userID string) error {
	const query = `DELETE FROM projects WHERE id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrProjectNotFound
	}
	return nil
}

func (r *ProjectRepository) Update(ctx context.Context, id, userID, name, defaultBranch string) (*models.Project, error) {
	const query = `
		UPDATE projects
		SET name = $1, default_branch = $2, updated_at = NOW()
		WHERE id = $3 AND user_id = $4
		RETURNING ` + projectColumns
	p, err := scanProject(r.db.QueryRowContext(ctx, query, name, defaultBranch, id, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProjectNotFound
	}
	return p, err
}

func (r *ProjectRepository) GetEnvVars(ctx context.Context, projectID string) ([]models.EnvVar, error) {
	const query = `SELECT key, value FROM project_env_vars WHERE project_id = $1 ORDER BY key`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vars []models.EnvVar
	for rows.Next() {
		var v models.EnvVar
		if err := rows.Scan(&v.Key, &v.Value); err != nil {
			return nil, err
		}
		vars = append(vars, v)
	}
	return vars, rows.Err()
}

func (r *ProjectRepository) ListSecretKeys(ctx context.Context, projectID string) ([]models.SecretMeta, error) {
	const query = `SELECT key, updated_at FROM project_secrets WHERE project_id = $1 ORDER BY key`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.SecretMeta
	for rows.Next() {
		var m models.SecretMeta
		if err := rows.Scan(&m.Key, &m.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *ProjectRepository) UpsertSecret(ctx context.Context, projectID, key, encryptedValue string) error {
	const query = `
		INSERT INTO project_secrets (project_id, key, value_encrypted, updated_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (project_id, key)
		DO UPDATE SET value_encrypted = EXCLUDED.value_encrypted, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, projectID, key, encryptedValue)
	return err
}

func (r *ProjectRepository) DeleteSecret(ctx context.Context, projectID, key string) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM project_secrets WHERE project_id = $1 AND key = $2`, projectID, key)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return errors.New("secret not found")
	}
	return nil
}

func (r *ProjectRepository) GetSecretsEncrypted(ctx context.Context, projectID string) (map[string]string, error) {
	const query = `SELECT key, value_encrypted FROM project_secrets WHERE project_id = $1`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]string)
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, rows.Err()
}

func (r *ProjectRepository) GetPipelineYAML(ctx context.Context, projectID string) (string, error) {
	var yaml string
	err := r.db.QueryRowContext(ctx,
		`SELECT pipeline_yaml_override FROM projects WHERE id = $1`, projectID,
	).Scan(&yaml)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrProjectNotFound
	}
	return yaml, err
}

func (r *ProjectRepository) SetPipelineYAML(ctx context.Context, projectID, userID, yaml string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE projects SET pipeline_yaml_override = $1, updated_at = NOW() WHERE id = $2 AND user_id = $3`,
		yaml, projectID, userID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrProjectNotFound
	}
	return nil
}

func (r *ProjectRepository) SetPipelineYAMLByProjectID(ctx context.Context, projectID, yaml string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE projects SET pipeline_yaml_override = $1, updated_at = NOW() WHERE id = $2`,
		yaml, projectID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrProjectNotFound
	}
	return nil
}

func (r *ProjectRepository) SetEnvVars(ctx context.Context, projectID string, vars []models.EnvVar) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM project_env_vars WHERE project_id = $1`, projectID); err != nil {
		return err
	}

	if len(vars) > 0 {
		const insert = `INSERT INTO project_env_vars (project_id, key, value) VALUES ($1, $2, $3)`
		for _, v := range vars {
			if _, err := tx.ExecContext(ctx, insert, projectID, v.Key, v.Value); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

var ErrMemberNotFound = errors.New("member not found")
var ErrMemberExists = errors.New("member already exists")
var ErrNotOwner = errors.New("only the project owner can perform this action")
var ErrCannotRemoveOwner = errors.New("cannot remove or change the role of the project owner")

func scanMember(s rowScanner) (*models.ProjectMember, error) {
	m := &models.ProjectMember{}
	return m, s.Scan(&m.ID, &m.ProjectID, &m.UserID, &m.Email, &m.DisplayName, &m.Role, &m.InvitedBy, &m.CreatedAt)
}

func (r *ProjectRepository) GetUserRoleAndOwner(ctx context.Context, projectID, userID string) (string, string, error) {
	var ownerID string
	err := r.db.QueryRowContext(ctx,
		`SELECT user_id FROM projects WHERE id = $1`, projectID,
	).Scan(&ownerID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", ErrProjectNotFound
	}
	if err != nil {
		return "", "", err
	}
	if ownerID == userID {
		return "owner", ownerID, nil
	}
	var role string
	err = r.db.QueryRowContext(ctx,
		`SELECT role FROM project_members WHERE project_id = $1 AND user_id = $2`, projectID, userID,
	).Scan(&role)
	if errors.Is(err, sql.ErrNoRows) {
		return "", "", ErrProjectNotFound
	}
	if err != nil {
		return "", "", err
	}
	return role, ownerID, nil
}

func (r *ProjectRepository) GetByIDForUser(ctx context.Context, id, userID string) (*models.Project, error) {
	const query = `
		SELECT ` + projectColumns + `
		FROM projects
		WHERE id = $1
		  AND (user_id = $2
		    OR EXISTS (SELECT 1 FROM project_members WHERE project_id = projects.id AND user_id = $2))
	`
	p, err := scanProject(r.db.QueryRowContext(ctx, query, id, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrProjectNotFound
	}
	return p, err
}

func (r *ProjectRepository) ListAccessible(ctx context.Context, userID string) ([]*models.Project, error) {
	const query = `
		SELECT id, name, repo_url, default_branch, status
		FROM projects
		WHERE user_id = $1
		   OR EXISTS (SELECT 1 FROM project_members WHERE project_id = projects.id AND user_id = $1)
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.Project
	for rows.Next() {
		p := &models.Project{}
		if err := rows.Scan(&p.ID, &p.Name, &p.RepoURL, &p.DefaultBranch, &p.Status); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *ProjectRepository) ListMembers(ctx context.Context, projectID string) ([]*models.ProjectMember, error) {
	const query = `
		SELECT id, project_id, user_id, email, display_name, role, invited_by, created_at
		FROM project_members
		WHERE project_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*models.ProjectMember
	for rows.Next() {
		m, err := scanMember(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *ProjectRepository) AddMember(ctx context.Context, m *models.ProjectMember) error {
	const query = `
		INSERT INTO project_members (id, project_id, user_id, email, display_name, role, invited_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING created_at
	`
	err := r.db.QueryRowContext(ctx, query,
		m.ID, m.ProjectID, m.UserID, m.Email, m.DisplayName, m.Role, m.InvitedBy,
	).Scan(&m.CreatedAt)
	if err != nil && strings.Contains(err.Error(), "duplicate key") {
		return ErrMemberExists
	}
	return err
}

func (r *ProjectRepository) UpdateMemberRole(ctx context.Context, projectID, userID, role string) (*models.ProjectMember, error) {
	const query = `
		UPDATE project_members
		SET role = $1
		WHERE project_id = $2 AND user_id = $3
		RETURNING id, project_id, user_id, email, display_name, role, invited_by, created_at
	`
	m, err := scanMember(r.db.QueryRowContext(ctx, query, role, projectID, userID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMemberNotFound
	}
	return m, err
}

func (r *ProjectRepository) RemoveMember(ctx context.Context, projectID, userID string) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM project_members WHERE project_id = $1 AND user_id = $2`, projectID, userID)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrMemberNotFound
	}
	return nil
}
