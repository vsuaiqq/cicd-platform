package grpchandler

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/vsuaiqq/cicd/projects-service/internal/models"
	"github.com/vsuaiqq/cicd/projects-service/internal/repository"
	"github.com/vsuaiqq/cicd/projects-service/internal/service"
	pb "github.com/vsuaiqq/cicd/shared/proto/gen/projects"
)

type ProjectsServer struct {
	pb.UnimplementedProjectsServiceServer
	svc *service.ProjectService
}

func NewProjectsServer(svc *service.ProjectService) *ProjectsServer {
	return &ProjectsServer{svc: svc}
}

func (s *ProjectsServer) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (*pb.CreateProjectResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.Name == "" || req.RepositoryUrl == "" {
		return nil, status.Error(codes.InvalidArgument, "name and repository_url required")
	}

	p, err := s.svc.Create(ctx, req.UserId, req.Name, req.RepositoryUrl, req.DefaultBranch)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateProjectResponse{
		Id:            p.ID,
		Name:          p.Name,
		RepoUrl:       p.RepoURL,
		DefaultBranch: p.DefaultBranch,
		PublicKey:     p.PublicKey,
		WebhookSecret: p.WebhookSecret,
		WebhookUrl:    s.svc.WebhookURL(p.ID),
		Status:        p.Status,
		OwnerUserId:   p.UserID,
	}, nil
}

func (s *ProjectsServer) GetProject(ctx context.Context, req *pb.GetProjectRequest) (*pb.GetProjectResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}

	p, err := s.svc.GetByID(ctx, req.ProjectId, req.UserId)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetProjectResponse{
		Id:            p.ID,
		Name:          p.Name,
		RepoUrl:       p.RepoURL,
		DefaultBranch: p.DefaultBranch,
		PublicKey:     p.PublicKey,
		WebhookSecret: p.WebhookSecret,
		WebhookUrl:    s.svc.WebhookURL(p.ID),
		Status:        p.Status,
		OwnerUserId:   p.UserID,
	}, nil
}

func (s *ProjectsServer) ListProjects(ctx context.Context, req *pb.ListProjectsRequest) (*pb.ListProjectsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}

	list, err := s.svc.ListByUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	projects := make([]*pb.ProjectSummary, 0, len(list))
	for _, p := range list {
		projects = append(projects, &pb.ProjectSummary{
			Id:            p.ID,
			Name:          p.Name,
			RepoUrl:       p.RepoURL,
			DefaultBranch: p.DefaultBranch,
			Status:        p.Status,
		})
	}
	return &pb.ListProjectsResponse{Projects: projects}, nil
}

func (s *ProjectsServer) VerifyConnection(ctx context.Context, req *pb.VerifyConnectionRequest) (*pb.VerifyConnectionResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}

	err := s.svc.VerifyConnection(ctx, req.ProjectId, req.UserId)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return &pb.VerifyConnectionResponse{Success: false, Error: err.Error()}, nil
	}
	return &pb.VerifyConnectionResponse{Success: true, Status: "active"}, nil
}

func (s *ProjectsServer) DeleteProject(ctx context.Context, req *pb.DeleteProjectRequest) (*pb.DeleteProjectResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}

	err := s.svc.Delete(ctx, req.ProjectId, req.UserId)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteProjectResponse{}, nil
}

func (s *ProjectsServer) UpdateProject(ctx context.Context, req *pb.UpdateProjectRequest) (*pb.UpdateProjectResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name required")
	}

	p, err := s.svc.Update(ctx, req.ProjectId, req.UserId, req.Name, req.DefaultBranch)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateProjectResponse{
		Id:            p.ID,
		Name:          p.Name,
		RepoUrl:       p.RepoURL,
		DefaultBranch: p.DefaultBranch,
		PublicKey:     p.PublicKey,
		WebhookSecret: p.WebhookSecret,
		WebhookUrl:    s.svc.WebhookURL(p.ID),
		Status:        p.Status,
		OwnerUserId:   p.UserID,
	}, nil
}

func (s *ProjectsServer) GetEnvVars(ctx context.Context, req *pb.GetEnvVarsRequest) (*pb.GetEnvVarsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}

	vars, err := s.svc.GetEnvVars(ctx, req.ProjectId, req.UserId)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetEnvVarsResponse{Vars: toProtoEnvVars(vars)}, nil
}

func (s *ProjectsServer) UpdateEnvVars(ctx context.Context, req *pb.UpdateEnvVarsRequest) (*pb.UpdateEnvVarsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}

	vars := fromProtoEnvVars(req.Vars)
	err := s.svc.UpdateEnvVars(ctx, req.ProjectId, req.UserId, vars)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateEnvVarsResponse{}, nil
}

func (s *ProjectsServer) ListSecrets(ctx context.Context, req *pb.ListSecretsRequest) (*pb.ListSecretsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}

	metas, err := s.svc.ListSecrets(ctx, req.ProjectId, req.UserId)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	out := make([]*pb.SecretMeta, len(metas))
	for i, m := range metas {
		out[i] = &pb.SecretMeta{
			Key:       m.Key,
			UpdatedAt: m.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}
	return &pb.ListSecretsResponse{Secrets: out}, nil
}

func (s *ProjectsServer) SetSecret(ctx context.Context, req *pb.SetSecretRequest) (*pb.SetSecretResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}
	if req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "key required")
	}
	if req.Value == "" {
		return nil, status.Error(codes.InvalidArgument, "value must not be empty")
	}

	err := s.svc.SetSecret(ctx, req.ProjectId, req.UserId, req.Key, req.Value)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.SetSecretResponse{}, nil
}

func (s *ProjectsServer) DeleteSecret(ctx context.Context, req *pb.DeleteSecretRequest) (*pb.DeleteSecretResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" || req.Key == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id and key required")
	}

	err := s.svc.DeleteSecret(ctx, req.ProjectId, req.UserId, req.Key)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.DeleteSecretResponse{}, nil
}

func (s *ProjectsServer) GetPipelineYAML(ctx context.Context, req *pb.GetPipelineYAMLRequest) (*pb.GetPipelineYAMLResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}

	yaml, err := s.svc.GetPipelineYAML(ctx, req.ProjectId, req.UserId)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetPipelineYAMLResponse{Yaml: yaml}, nil
}

func (s *ProjectsServer) SetPipelineYAML(ctx context.Context, req *pb.SetPipelineYAMLRequest) (*pb.SetPipelineYAMLResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}

	err := s.svc.SetPipelineYAML(ctx, req.ProjectId, req.UserId, req.Yaml)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.SetPipelineYAMLResponse{}, nil
}

func (s *ProjectsServer) Health(_ context.Context, _ *pb.HealthRequest) (*pb.HealthResponse, error) {
	return &pb.HealthResponse{Healthy: true, Status: "ok"}, nil
}

func (s *ProjectsServer) GetProjectInternal(ctx context.Context, req *pb.GetProjectInternalRequest) (*pb.GetProjectInternalResponse, error) {
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}

	p, sshKey, err := s.svc.GetSSHKeyForProject(ctx, req.ProjectId)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	envVars, err := s.svc.GetEnvVarsInternal(ctx, req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	pipelineYAML, err := s.svc.GetPipelineYAMLInternal(ctx, req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	secretsMap, err := s.svc.GetSecretsForPipeline(ctx, req.ProjectId)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &pb.GetProjectInternalResponse{
		Id:                   p.ID,
		RepoUrl:              p.RepoURL,
		WebhookSecret:        p.WebhookSecret,
		SshKey:               string(sshKey),
		EnvVars:              toProtoEnvVars(envVars),
		PipelineYamlOverride: pipelineYAML,
		Secrets:              toProtoSecrets(secretsMap),
	}, nil
}

func toProtoEnvVars(vars []models.EnvVar) []*pb.EnvVar {
	out := make([]*pb.EnvVar, len(vars))
	for i, v := range vars {
		out[i] = &pb.EnvVar{Key: v.Key, Value: v.Value}
	}
	return out
}

func toProtoSecrets(m map[string]string) []*pb.SecretKV {
	out := make([]*pb.SecretKV, 0, len(m))
	for k, v := range m {
		out = append(out, &pb.SecretKV{Key: k, Value: v})
	}
	return out
}

func fromProtoEnvVars(vars []*pb.EnvVar) []models.EnvVar {
	out := make([]models.EnvVar, 0, len(vars))
	for _, v := range vars {
		if v.Key != "" {
			out = append(out, models.EnvVar{Key: v.Key, Value: v.Value})
		}
	}
	return out
}

func (s *ProjectsServer) ListMembers(ctx context.Context, req *pb.ListMembersRequest) (*pb.ListMembersResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.Unauthenticated, "user_id required")
	}
	if req.ProjectId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id required")
	}
	members, requesterRole, ownerUserID, err := s.svc.ListMembers(ctx, req.ProjectId, req.UserId)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	out := make([]*pb.ProjectMember, 0, len(members))
	for _, m := range members {
		out = append(out, toProtoMember(m))
	}
	return &pb.ListMembersResponse{Members: out, RequesterRole: requesterRole, OwnerUserId: ownerUserID}, nil
}

func (s *ProjectsServer) InviteMember(ctx context.Context, req *pb.InviteMemberRequest) (*pb.InviteMemberResponse, error) {
	if req.OwnerUserId == "" {
		return nil, status.Error(codes.Unauthenticated, "owner_user_id required")
	}
	if req.ProjectId == "" || req.InviteeUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id and invitee_user_id required")
	}
	m, err := s.svc.InviteMember(ctx, req.ProjectId, req.OwnerUserId, req.InviteeUserId, req.InviteeEmail, req.InviteeDisplayName, req.Role)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if errors.Is(err, repository.ErrNotOwner) {
		return nil, status.Error(codes.PermissionDenied, "only the project owner can invite members")
	}
	if errors.Is(err, repository.ErrMemberExists) {
		return nil, status.Error(codes.AlreadyExists, "user is already a member of this project")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.InviteMemberResponse{Member: toProtoMember(m)}, nil
}

func (s *ProjectsServer) UpdateMemberRole(ctx context.Context, req *pb.UpdateMemberRoleRequest) (*pb.UpdateMemberRoleResponse, error) {
	if req.OwnerUserId == "" {
		return nil, status.Error(codes.Unauthenticated, "owner_user_id required")
	}
	if req.ProjectId == "" || req.TargetUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id and target_user_id required")
	}
	m, err := s.svc.UpdateMemberRole(ctx, req.ProjectId, req.OwnerUserId, req.TargetUserId, req.Role)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if errors.Is(err, repository.ErrNotOwner) {
		return nil, status.Error(codes.PermissionDenied, "only the project owner can update member roles")
	}
	if errors.Is(err, repository.ErrCannotRemoveOwner) {
		return nil, status.Error(codes.InvalidArgument, "cannot change the owner's own role")
	}
	if errors.Is(err, repository.ErrMemberNotFound) {
		return nil, status.Error(codes.NotFound, "member not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateMemberRoleResponse{Member: toProtoMember(m)}, nil
}

func (s *ProjectsServer) RemoveMember(ctx context.Context, req *pb.RemoveMemberRequest) (*pb.RemoveMemberResponse, error) {
	if req.OwnerUserId == "" {
		return nil, status.Error(codes.Unauthenticated, "owner_user_id required")
	}
	if req.ProjectId == "" || req.TargetUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "project_id and target_user_id required")
	}
	err := s.svc.RemoveMember(ctx, req.ProjectId, req.OwnerUserId, req.TargetUserId)
	if errors.Is(err, repository.ErrProjectNotFound) {
		return nil, status.Error(codes.NotFound, "project not found")
	}
	if errors.Is(err, repository.ErrNotOwner) {
		return nil, status.Error(codes.PermissionDenied, "only the project owner can remove members")
	}
	if errors.Is(err, repository.ErrCannotRemoveOwner) {
		return nil, status.Error(codes.InvalidArgument, "cannot remove the project owner")
	}
	if errors.Is(err, repository.ErrMemberNotFound) {
		return nil, status.Error(codes.NotFound, "member not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.RemoveMemberResponse{}, nil
}

func toProtoMember(m *models.ProjectMember) *pb.ProjectMember {
	return &pb.ProjectMember{
		UserId:      m.UserID,
		Email:       m.Email,
		DisplayName: m.DisplayName,
		Role:        m.Role,
		InvitedBy:   m.InvitedBy,
		CreatedAt:   m.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}
