package service

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/vsuaiqq/cicd/projects-service/internal/models"
	"github.com/vsuaiqq/cicd/projects-service/internal/repository"
)

type ProjectService struct {
	repo           *repository.ProjectRepository
	encrypter      *Encrypter
	webhookBaseURL string
}

func NewProjectService(repo *repository.ProjectRepository, encrypter *Encrypter, webhookBaseURL string) *ProjectService {
	return &ProjectService{
		repo:           repo,
		encrypter:      encrypter,
		webhookBaseURL: webhookBaseURL,
	}
}

func (s *ProjectService) Create(ctx context.Context, userID, name, repoURL, defaultBranch string) (*models.Project, error) {
	if defaultBranch == "" {
		defaultBranch = "main"
	}

	privatePEM, publicKey, err := GenerateSSHKeyPair()
	if err != nil {
		return nil, fmt.Errorf("generate ssh key: %w", err)
	}
	webhookSecret, err := GenerateWebhookSecret()
	if err != nil {
		return nil, fmt.Errorf("generate webhook secret: %w", err)
	}
	encrypted, err := s.encrypter.Encrypt(privatePEM)
	if err != nil {
		return nil, fmt.Errorf("encrypt private key: %w", err)
	}

	p := &models.Project{
		ID:                  uuid.New().String(),
		UserID:              userID,
		Name:                name,
		RepoURL:             repoURL,
		DefaultBranch:       defaultBranch,
		PrivateKeyEncrypted: encrypted,
		PublicKey:           publicKey,
		WebhookSecret:       webhookSecret,
		Status:              "pending",
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProjectService) GetByID(ctx context.Context, id, userID string) (*models.Project, error) {
	return s.repo.GetByIDForUser(ctx, id, userID)
}

func (s *ProjectService) GetSSHKeyForProject(ctx context.Context, projectID string) (*models.Project, []byte, error) {
	p, err := s.repo.GetByID(ctx, projectID)
	if err != nil {
		return nil, nil, err
	}
	key, err := s.encrypter.Decrypt(p.PrivateKeyEncrypted)
	if err != nil {
		return nil, nil, fmt.Errorf("decrypt key: %w", err)
	}
	return p, key, nil
}

func (s *ProjectService) ListByUser(ctx context.Context, userID string) ([]*models.Project, error) {
	return s.repo.ListAccessible(ctx, userID)
}

func (s *ProjectService) WebhookURL(projectID string) string {
	if s.webhookBaseURL == "" {
		return ""
	}
	return strings.TrimRight(s.webhookBaseURL, "/") + "/" + projectID
}

func (s *ProjectService) Delete(ctx context.Context, id, userID string) error {
	return s.repo.DeleteByIDAndUser(ctx, id, userID)
}

func (s *ProjectService) Update(ctx context.Context, id, userID, name, defaultBranch string) (*models.Project, error) {
	return s.repo.Update(ctx, id, userID, name, defaultBranch)
}

func (s *ProjectService) GetEnvVars(ctx context.Context, id, userID string) ([]models.EnvVar, error) {
	if err := s.requireAtLeastViewer(ctx, id, userID); err != nil {
		return nil, err
	}
	return s.repo.GetEnvVars(ctx, id)
}

func (s *ProjectService) UpdateEnvVars(ctx context.Context, id, userID string, vars []models.EnvVar) error {
	if err := s.requireAtLeastEditor(ctx, id, userID); err != nil {
		return err
	}
	return s.repo.SetEnvVars(ctx, id, vars)
}

func (s *ProjectService) GetEnvVarsInternal(ctx context.Context, projectID string) ([]models.EnvVar, error) {
	return s.repo.GetEnvVars(ctx, projectID)
}

func (s *ProjectService) GetPipelineYAML(ctx context.Context, id, userID string) (string, error) {
	if err := s.requireAtLeastViewer(ctx, id, userID); err != nil {
		return "", err
	}
	return s.repo.GetPipelineYAML(ctx, id)
}

func (s *ProjectService) SetPipelineYAML(ctx context.Context, id, userID, yaml string) error {
	if err := s.requireAtLeastEditor(ctx, id, userID); err != nil {
		return err
	}

	return s.repo.SetPipelineYAMLByProjectID(ctx, id, yaml)
}

func (s *ProjectService) GetPipelineYAMLInternal(ctx context.Context, projectID string) (string, error) {
	return s.repo.GetPipelineYAML(ctx, projectID)
}

func (s *ProjectService) ListSecrets(ctx context.Context, projectID, userID string) ([]models.SecretMeta, error) {
	if err := s.requireAtLeastViewer(ctx, projectID, userID); err != nil {
		return nil, err
	}
	return s.repo.ListSecretKeys(ctx, projectID)
}

func (s *ProjectService) SetSecret(ctx context.Context, projectID, userID, key, value string) error {
	if err := s.requireAtLeastEditor(ctx, projectID, userID); err != nil {
		return err
	}
	encrypted, err := s.encrypter.Encrypt([]byte(value))
	if err != nil {
		return fmt.Errorf("encrypt secret: %w", err)
	}
	return s.repo.UpsertSecret(ctx, projectID, key, string(encrypted))
}

func (s *ProjectService) DeleteSecret(ctx context.Context, projectID, userID, key string) error {
	if err := s.requireAtLeastEditor(ctx, projectID, userID); err != nil {
		return err
	}
	return s.repo.DeleteSecret(ctx, projectID, key)
}

func (s *ProjectService) GetSecretsForPipeline(ctx context.Context, projectID string) (map[string]string, error) {
	encrypted, err := s.repo.GetSecretsEncrypted(ctx, projectID)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(encrypted))
	for k, v := range encrypted {
		plain, err := s.encrypter.Decrypt([]byte(v))
		if err != nil {
			return nil, fmt.Errorf("decrypt secret %q: %w", k, err)
		}
		out[k] = string(plain)
	}
	return out, nil
}

func (s *ProjectService) requireAtLeastViewer(ctx context.Context, projectID, userID string) error {
	_, _, err := s.repo.GetUserRoleAndOwner(ctx, projectID, userID)
	return err
}

func (s *ProjectService) requireAtLeastEditor(ctx context.Context, projectID, userID string) error {
	role, _, err := s.repo.GetUserRoleAndOwner(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if role == "viewer" {
		return repository.ErrProjectNotFound
	}
	return nil
}

func (s *ProjectService) requireOwner(ctx context.Context, projectID, userID string) error {
	role, _, err := s.repo.GetUserRoleAndOwner(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if role != "owner" {
		return repository.ErrNotOwner
	}
	return nil
}

func (s *ProjectService) ListMembers(ctx context.Context, projectID, userID string) ([]*models.ProjectMember, string, string, error) {
	role, ownerUserID, err := s.repo.GetUserRoleAndOwner(ctx, projectID, userID)
	if err != nil {
		return nil, "", "", err
	}
	members, err := s.repo.ListMembers(ctx, projectID)
	if err != nil {
		return nil, "", "", err
	}
	return members, role, ownerUserID, nil
}

func (s *ProjectService) InviteMember(ctx context.Context, projectID, ownerUserID, inviteeUserID, inviteeEmail, inviteeDisplayName, role string) (*models.ProjectMember, error) {
	if err := s.requireOwner(ctx, projectID, ownerUserID); err != nil {
		return nil, err
	}
	if role != "editor" && role != "viewer" {
		role = "viewer"
	}

	if inviteeUserID == ownerUserID {
		return nil, repository.ErrMemberExists
	}
	m := &models.ProjectMember{
		ID:          uuid.New().String(),
		ProjectID:   projectID,
		UserID:      inviteeUserID,
		Email:       inviteeEmail,
		DisplayName: inviteeDisplayName,
		Role:        role,
		InvitedBy:   ownerUserID,
	}
	if err := s.repo.AddMember(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *ProjectService) UpdateMemberRole(ctx context.Context, projectID, ownerUserID, targetUserID, role string) (*models.ProjectMember, error) {
	if err := s.requireOwner(ctx, projectID, ownerUserID); err != nil {
		return nil, err
	}
	if targetUserID == ownerUserID {
		return nil, repository.ErrCannotRemoveOwner
	}
	if role != "editor" && role != "viewer" {
		role = "viewer"
	}
	return s.repo.UpdateMemberRole(ctx, projectID, targetUserID, role)
}

func (s *ProjectService) RemoveMember(ctx context.Context, projectID, ownerUserID, targetUserID string) error {
	if err := s.requireOwner(ctx, projectID, ownerUserID); err != nil {
		return err
	}
	if targetUserID == ownerUserID {
		return repository.ErrCannotRemoveOwner
	}
	return s.repo.RemoveMember(ctx, projectID, targetUserID)
}

func (s *ProjectService) VerifyConnection(ctx context.Context, id, userID string) error {
	p, err := s.repo.GetByIDForUser(ctx, id, userID)
	if err != nil {
		return err
	}

	privatePEM, err := s.encrypter.Decrypt(p.PrivateKeyEncrypted)
	if err != nil {
		return fmt.Errorf("decrypt key: %w", err)
	}

	dir, err := os.MkdirTemp("", "cicd-verify-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	keyPath := filepath.Join(dir, "key")
	if err := os.WriteFile(keyPath, privatePEM, 0600); err != nil {
		return fmt.Errorf("write key file: %w", err)
	}

	sshCmd := fmt.Sprintf(
		"ssh -i '%s' -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o BatchMode=yes",
		keyPath,
	)
	cmd := exec.CommandContext(ctx, "git", "ls-remote", p.RepoURL)

	cmd.Env = []string{
		"HOME=" + os.TempDir(),
		"PATH=" + os.Getenv("PATH"),
		"GIT_SSH_COMMAND=" + sshCmd,
		"GIT_TERMINAL_PROMPT=0",
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git ls-remote failed: %w (output: %s)", err, string(out))
	}

	return s.repo.UpdateStatus(ctx, id, userID, "active")
}
