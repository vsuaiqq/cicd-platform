package client

import (
	"context"
	"time"

	"google.golang.org/grpc"

	pb "github.com/vsuaiqq/cicd/shared/proto/gen/projects"
)

type ProjectsClient struct {
	conn    *grpc.ClientConn
	client  pb.ProjectsServiceClient
	timeout time.Duration
}

func NewProjectsClient(address string, timeout time.Duration) (*ProjectsClient, error) {
	conn, err := newGRPCConn(address)
	if err != nil {
		return nil, err
	}
	return &ProjectsClient{
		conn:    conn,
		client:  pb.NewProjectsServiceClient(conn),
		timeout: timeout,
	}, nil
}

func (c *ProjectsClient) Close() error {
	return c.conn.Close()
}

func (c *ProjectsClient) CreateProject(ctx context.Context, userID, name, repositoryURL, defaultBranch string) (*pb.CreateProjectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.CreateProject(ctx, &pb.CreateProjectRequest{
		UserId:        userID,
		Name:          name,
		RepositoryUrl: repositoryURL,
		DefaultBranch: defaultBranch,
	})
}

func (c *ProjectsClient) GetProject(ctx context.Context, userID, projectID string) (*pb.GetProjectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.GetProject(ctx, &pb.GetProjectRequest{
		UserId:    userID,
		ProjectId: projectID,
	})
}

func (c *ProjectsClient) ListProjects(ctx context.Context, userID string) (*pb.ListProjectsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.ListProjects(ctx, &pb.ListProjectsRequest{UserId: userID})
}

func (c *ProjectsClient) VerifyConnection(ctx context.Context, userID, projectID string) (*pb.VerifyConnectionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.VerifyConnection(ctx, &pb.VerifyConnectionRequest{
		UserId:    userID,
		ProjectId: projectID,
	})
}

func (c *ProjectsClient) DeleteProject(ctx context.Context, userID, projectID string) (*pb.DeleteProjectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.DeleteProject(ctx, &pb.DeleteProjectRequest{
		UserId:    userID,
		ProjectId: projectID,
	})
}

func (c *ProjectsClient) UpdateProject(ctx context.Context, userID, projectID, name, defaultBranch string) (*pb.UpdateProjectResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.UpdateProject(ctx, &pb.UpdateProjectRequest{
		UserId:        userID,
		ProjectId:     projectID,
		Name:          name,
		DefaultBranch: defaultBranch,
	})
}

func (c *ProjectsClient) GetEnvVars(ctx context.Context, userID, projectID string) (*pb.GetEnvVarsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.GetEnvVars(ctx, &pb.GetEnvVarsRequest{
		UserId:    userID,
		ProjectId: projectID,
	})
}

func (c *ProjectsClient) UpdateEnvVars(ctx context.Context, userID, projectID string, vars []*pb.EnvVar) (*pb.UpdateEnvVarsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.UpdateEnvVars(ctx, &pb.UpdateEnvVarsRequest{
		UserId:    userID,
		ProjectId: projectID,
		Vars:      vars,
	})
}

func (c *ProjectsClient) ListSecrets(ctx context.Context, userID, projectID string) (*pb.ListSecretsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.ListSecrets(ctx, &pb.ListSecretsRequest{
		UserId:    userID,
		ProjectId: projectID,
	})
}

func (c *ProjectsClient) SetSecret(ctx context.Context, userID, projectID, key, value string) (*pb.SetSecretResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.SetSecret(ctx, &pb.SetSecretRequest{
		UserId:    userID,
		ProjectId: projectID,
		Key:       key,
		Value:     value,
	})
}

func (c *ProjectsClient) DeleteSecret(ctx context.Context, userID, projectID, key string) (*pb.DeleteSecretResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.DeleteSecret(ctx, &pb.DeleteSecretRequest{
		UserId:    userID,
		ProjectId: projectID,
		Key:       key,
	})
}

func (c *ProjectsClient) GetPipelineYAML(ctx context.Context, userID, projectID string) (*pb.GetPipelineYAMLResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.GetPipelineYAML(ctx, &pb.GetPipelineYAMLRequest{
		UserId:    userID,
		ProjectId: projectID,
	})
}

func (c *ProjectsClient) SetPipelineYAML(ctx context.Context, userID, projectID, yaml string) (*pb.SetPipelineYAMLResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.SetPipelineYAML(ctx, &pb.SetPipelineYAMLRequest{
		UserId:    userID,
		ProjectId: projectID,
		Yaml:      yaml,
	})
}

func (c *ProjectsClient) Health(ctx context.Context) (*pb.HealthResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.Health(ctx, &pb.HealthRequest{})
}

func (c *ProjectsClient) ListMembers(ctx context.Context, userID, projectID string) (*pb.ListMembersResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.ListMembers(ctx, &pb.ListMembersRequest{
		UserId:    userID,
		ProjectId: projectID,
	})
}

func (c *ProjectsClient) InviteMember(ctx context.Context, ownerUserID, projectID, inviteeUserID, inviteeEmail, inviteeDisplayName, role string) (*pb.InviteMemberResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.InviteMember(ctx, &pb.InviteMemberRequest{
		OwnerUserId:        ownerUserID,
		ProjectId:          projectID,
		InviteeUserId:      inviteeUserID,
		InviteeEmail:       inviteeEmail,
		InviteeDisplayName: inviteeDisplayName,
		Role:               role,
	})
}

func (c *ProjectsClient) UpdateMemberRole(ctx context.Context, ownerUserID, projectID, targetUserID, role string) (*pb.UpdateMemberRoleResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.UpdateMemberRole(ctx, &pb.UpdateMemberRoleRequest{
		OwnerUserId:  ownerUserID,
		ProjectId:    projectID,
		TargetUserId: targetUserID,
		Role:         role,
	})
}

func (c *ProjectsClient) RemoveMember(ctx context.Context, ownerUserID, projectID, targetUserID string) (*pb.RemoveMemberResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	return c.client.RemoveMember(ctx, &pb.RemoveMemberRequest{
		OwnerUserId:  ownerUserID,
		ProjectId:    projectID,
		TargetUserId: targetUserID,
	})
}
