package projects

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9

const (
	ProjectsService_CreateProject_FullMethodName        = "/projects.ProjectsService/CreateProject"
	ProjectsService_GetProject_FullMethodName           = "/projects.ProjectsService/GetProject"
	ProjectsService_ListProjects_FullMethodName         = "/projects.ProjectsService/ListProjects"
	ProjectsService_VerifyConnection_FullMethodName     = "/projects.ProjectsService/VerifyConnection"
	ProjectsService_DeleteProject_FullMethodName        = "/projects.ProjectsService/DeleteProject"
	ProjectsService_UpdateProject_FullMethodName        = "/projects.ProjectsService/UpdateProject"
	ProjectsService_Health_FullMethodName               = "/projects.ProjectsService/Health"
	ProjectsService_GetEnvVars_FullMethodName           = "/projects.ProjectsService/GetEnvVars"
	ProjectsService_UpdateEnvVars_FullMethodName        = "/projects.ProjectsService/UpdateEnvVars"
	ProjectsService_GetPipelineYAML_FullMethodName      = "/projects.ProjectsService/GetPipelineYAML"
	ProjectsService_SetPipelineYAML_FullMethodName      = "/projects.ProjectsService/SetPipelineYAML"
	ProjectsService_ListSecrets_FullMethodName          = "/projects.ProjectsService/ListSecrets"
	ProjectsService_SetSecret_FullMethodName            = "/projects.ProjectsService/SetSecret"
	ProjectsService_DeleteSecret_FullMethodName         = "/projects.ProjectsService/DeleteSecret"
	ProjectsService_GetProjectInternal_FullMethodName   = "/projects.ProjectsService/GetProjectInternal"
	ProjectsService_ListProjectsInternal_FullMethodName = "/projects.ProjectsService/ListProjectsInternal"
	ProjectsService_ListMembers_FullMethodName          = "/projects.ProjectsService/ListMembers"
	ProjectsService_InviteMember_FullMethodName         = "/projects.ProjectsService/InviteMember"
	ProjectsService_UpdateMemberRole_FullMethodName     = "/projects.ProjectsService/UpdateMemberRole"
	ProjectsService_RemoveMember_FullMethodName         = "/projects.ProjectsService/RemoveMember"
)

type ProjectsServiceClient interface {
	CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*CreateProjectResponse, error)
	GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*GetProjectResponse, error)
	ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error)
	VerifyConnection(ctx context.Context, in *VerifyConnectionRequest, opts ...grpc.CallOption) (*VerifyConnectionResponse, error)
	DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*DeleteProjectResponse, error)
	UpdateProject(ctx context.Context, in *UpdateProjectRequest, opts ...grpc.CallOption) (*UpdateProjectResponse, error)
	Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error)
	GetEnvVars(ctx context.Context, in *GetEnvVarsRequest, opts ...grpc.CallOption) (*GetEnvVarsResponse, error)
	UpdateEnvVars(ctx context.Context, in *UpdateEnvVarsRequest, opts ...grpc.CallOption) (*UpdateEnvVarsResponse, error)
	GetPipelineYAML(ctx context.Context, in *GetPipelineYAMLRequest, opts ...grpc.CallOption) (*GetPipelineYAMLResponse, error)
	SetPipelineYAML(ctx context.Context, in *SetPipelineYAMLRequest, opts ...grpc.CallOption) (*SetPipelineYAMLResponse, error)

	ListSecrets(ctx context.Context, in *ListSecretsRequest, opts ...grpc.CallOption) (*ListSecretsResponse, error)
	SetSecret(ctx context.Context, in *SetSecretRequest, opts ...grpc.CallOption) (*SetSecretResponse, error)
	DeleteSecret(ctx context.Context, in *DeleteSecretRequest, opts ...grpc.CallOption) (*DeleteSecretResponse, error)

	GetProjectInternal(ctx context.Context, in *GetProjectInternalRequest, opts ...grpc.CallOption) (*GetProjectInternalResponse, error)
	ListProjectsInternal(ctx context.Context, in *ListProjectsInternalRequest, opts ...grpc.CallOption) (*ListProjectsInternalResponse, error)

	ListMembers(ctx context.Context, in *ListMembersRequest, opts ...grpc.CallOption) (*ListMembersResponse, error)
	InviteMember(ctx context.Context, in *InviteMemberRequest, opts ...grpc.CallOption) (*InviteMemberResponse, error)
	UpdateMemberRole(ctx context.Context, in *UpdateMemberRoleRequest, opts ...grpc.CallOption) (*UpdateMemberRoleResponse, error)
	RemoveMember(ctx context.Context, in *RemoveMemberRequest, opts ...grpc.CallOption) (*RemoveMemberResponse, error)
}

type projectsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProjectsServiceClient(cc grpc.ClientConnInterface) ProjectsServiceClient {
	return &projectsServiceClient{cc}
}

func (c *projectsServiceClient) CreateProject(ctx context.Context, in *CreateProjectRequest, opts ...grpc.CallOption) (*CreateProjectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateProjectResponse)
	err := c.cc.Invoke(ctx, ProjectsService_CreateProject_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) GetProject(ctx context.Context, in *GetProjectRequest, opts ...grpc.CallOption) (*GetProjectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProjectResponse)
	err := c.cc.Invoke(ctx, ProjectsService_GetProject_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) ListProjects(ctx context.Context, in *ListProjectsRequest, opts ...grpc.CallOption) (*ListProjectsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListProjectsResponse)
	err := c.cc.Invoke(ctx, ProjectsService_ListProjects_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) VerifyConnection(ctx context.Context, in *VerifyConnectionRequest, opts ...grpc.CallOption) (*VerifyConnectionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VerifyConnectionResponse)
	err := c.cc.Invoke(ctx, ProjectsService_VerifyConnection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) DeleteProject(ctx context.Context, in *DeleteProjectRequest, opts ...grpc.CallOption) (*DeleteProjectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteProjectResponse)
	err := c.cc.Invoke(ctx, ProjectsService_DeleteProject_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) UpdateProject(ctx context.Context, in *UpdateProjectRequest, opts ...grpc.CallOption) (*UpdateProjectResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateProjectResponse)
	err := c.cc.Invoke(ctx, ProjectsService_UpdateProject_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, ProjectsService_Health_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) GetEnvVars(ctx context.Context, in *GetEnvVarsRequest, opts ...grpc.CallOption) (*GetEnvVarsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetEnvVarsResponse)
	err := c.cc.Invoke(ctx, ProjectsService_GetEnvVars_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) UpdateEnvVars(ctx context.Context, in *UpdateEnvVarsRequest, opts ...grpc.CallOption) (*UpdateEnvVarsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateEnvVarsResponse)
	err := c.cc.Invoke(ctx, ProjectsService_UpdateEnvVars_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) GetPipelineYAML(ctx context.Context, in *GetPipelineYAMLRequest, opts ...grpc.CallOption) (*GetPipelineYAMLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPipelineYAMLResponse)
	err := c.cc.Invoke(ctx, ProjectsService_GetPipelineYAML_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) SetPipelineYAML(ctx context.Context, in *SetPipelineYAMLRequest, opts ...grpc.CallOption) (*SetPipelineYAMLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetPipelineYAMLResponse)
	err := c.cc.Invoke(ctx, ProjectsService_SetPipelineYAML_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) ListSecrets(ctx context.Context, in *ListSecretsRequest, opts ...grpc.CallOption) (*ListSecretsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListSecretsResponse)
	err := c.cc.Invoke(ctx, ProjectsService_ListSecrets_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) SetSecret(ctx context.Context, in *SetSecretRequest, opts ...grpc.CallOption) (*SetSecretResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetSecretResponse)
	err := c.cc.Invoke(ctx, ProjectsService_SetSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) DeleteSecret(ctx context.Context, in *DeleteSecretRequest, opts ...grpc.CallOption) (*DeleteSecretResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteSecretResponse)
	err := c.cc.Invoke(ctx, ProjectsService_DeleteSecret_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) GetProjectInternal(ctx context.Context, in *GetProjectInternalRequest, opts ...grpc.CallOption) (*GetProjectInternalResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProjectInternalResponse)
	err := c.cc.Invoke(ctx, ProjectsService_GetProjectInternal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) ListProjectsInternal(ctx context.Context, in *ListProjectsInternalRequest, opts ...grpc.CallOption) (*ListProjectsInternalResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListProjectsInternalResponse)
	err := c.cc.Invoke(ctx, ProjectsService_ListProjectsInternal_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) ListMembers(ctx context.Context, in *ListMembersRequest, opts ...grpc.CallOption) (*ListMembersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListMembersResponse)
	err := c.cc.Invoke(ctx, ProjectsService_ListMembers_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) InviteMember(ctx context.Context, in *InviteMemberRequest, opts ...grpc.CallOption) (*InviteMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(InviteMemberResponse)
	err := c.cc.Invoke(ctx, ProjectsService_InviteMember_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) UpdateMemberRole(ctx context.Context, in *UpdateMemberRoleRequest, opts ...grpc.CallOption) (*UpdateMemberRoleResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateMemberRoleResponse)
	err := c.cc.Invoke(ctx, ProjectsService_UpdateMemberRole_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectsServiceClient) RemoveMember(ctx context.Context, in *RemoveMemberRequest, opts ...grpc.CallOption) (*RemoveMemberResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RemoveMemberResponse)
	err := c.cc.Invoke(ctx, ProjectsService_RemoveMember_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type ProjectsServiceServer interface {
	CreateProject(context.Context, *CreateProjectRequest) (*CreateProjectResponse, error)
	GetProject(context.Context, *GetProjectRequest) (*GetProjectResponse, error)
	ListProjects(context.Context, *ListProjectsRequest) (*ListProjectsResponse, error)
	VerifyConnection(context.Context, *VerifyConnectionRequest) (*VerifyConnectionResponse, error)
	DeleteProject(context.Context, *DeleteProjectRequest) (*DeleteProjectResponse, error)
	UpdateProject(context.Context, *UpdateProjectRequest) (*UpdateProjectResponse, error)
	Health(context.Context, *HealthRequest) (*HealthResponse, error)
	GetEnvVars(context.Context, *GetEnvVarsRequest) (*GetEnvVarsResponse, error)
	UpdateEnvVars(context.Context, *UpdateEnvVarsRequest) (*UpdateEnvVarsResponse, error)
	GetPipelineYAML(context.Context, *GetPipelineYAMLRequest) (*GetPipelineYAMLResponse, error)
	SetPipelineYAML(context.Context, *SetPipelineYAMLRequest) (*SetPipelineYAMLResponse, error)

	ListSecrets(context.Context, *ListSecretsRequest) (*ListSecretsResponse, error)
	SetSecret(context.Context, *SetSecretRequest) (*SetSecretResponse, error)
	DeleteSecret(context.Context, *DeleteSecretRequest) (*DeleteSecretResponse, error)

	GetProjectInternal(context.Context, *GetProjectInternalRequest) (*GetProjectInternalResponse, error)
	ListProjectsInternal(context.Context, *ListProjectsInternalRequest) (*ListProjectsInternalResponse, error)

	ListMembers(context.Context, *ListMembersRequest) (*ListMembersResponse, error)
	InviteMember(context.Context, *InviteMemberRequest) (*InviteMemberResponse, error)
	UpdateMemberRole(context.Context, *UpdateMemberRoleRequest) (*UpdateMemberRoleResponse, error)
	RemoveMember(context.Context, *RemoveMemberRequest) (*RemoveMemberResponse, error)
	mustEmbedUnimplementedProjectsServiceServer()
}

type UnimplementedProjectsServiceServer struct{}

func (UnimplementedProjectsServiceServer) CreateProject(context.Context, *CreateProjectRequest) (*CreateProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateProject not implemented")
}
func (UnimplementedProjectsServiceServer) GetProject(context.Context, *GetProjectRequest) (*GetProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProject not implemented")
}
func (UnimplementedProjectsServiceServer) ListProjects(context.Context, *ListProjectsRequest) (*ListProjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjects not implemented")
}
func (UnimplementedProjectsServiceServer) VerifyConnection(context.Context, *VerifyConnectionRequest) (*VerifyConnectionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyConnection not implemented")
}
func (UnimplementedProjectsServiceServer) DeleteProject(context.Context, *DeleteProjectRequest) (*DeleteProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteProject not implemented")
}
func (UnimplementedProjectsServiceServer) UpdateProject(context.Context, *UpdateProjectRequest) (*UpdateProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateProject not implemented")
}
func (UnimplementedProjectsServiceServer) Health(context.Context, *HealthRequest) (*HealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Health not implemented")
}
func (UnimplementedProjectsServiceServer) GetEnvVars(context.Context, *GetEnvVarsRequest) (*GetEnvVarsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEnvVars not implemented")
}
func (UnimplementedProjectsServiceServer) UpdateEnvVars(context.Context, *UpdateEnvVarsRequest) (*UpdateEnvVarsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEnvVars not implemented")
}
func (UnimplementedProjectsServiceServer) GetPipelineYAML(context.Context, *GetPipelineYAMLRequest) (*GetPipelineYAMLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPipelineYAML not implemented")
}
func (UnimplementedProjectsServiceServer) SetPipelineYAML(context.Context, *SetPipelineYAMLRequest) (*SetPipelineYAMLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetPipelineYAML not implemented")
}
func (UnimplementedProjectsServiceServer) ListSecrets(context.Context, *ListSecretsRequest) (*ListSecretsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListSecrets not implemented")
}
func (UnimplementedProjectsServiceServer) SetSecret(context.Context, *SetSecretRequest) (*SetSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetSecret not implemented")
}
func (UnimplementedProjectsServiceServer) DeleteSecret(context.Context, *DeleteSecretRequest) (*DeleteSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteSecret not implemented")
}
func (UnimplementedProjectsServiceServer) GetProjectInternal(context.Context, *GetProjectInternalRequest) (*GetProjectInternalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProjectInternal not implemented")
}
func (UnimplementedProjectsServiceServer) ListProjectsInternal(context.Context, *ListProjectsInternalRequest) (*ListProjectsInternalResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListProjectsInternal not implemented")
}
func (UnimplementedProjectsServiceServer) ListMembers(context.Context, *ListMembersRequest) (*ListMembersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListMembers not implemented")
}
func (UnimplementedProjectsServiceServer) InviteMember(context.Context, *InviteMemberRequest) (*InviteMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteMember not implemented")
}
func (UnimplementedProjectsServiceServer) UpdateMemberRole(context.Context, *UpdateMemberRoleRequest) (*UpdateMemberRoleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMemberRole not implemented")
}
func (UnimplementedProjectsServiceServer) RemoveMember(context.Context, *RemoveMemberRequest) (*RemoveMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveMember not implemented")
}
func (UnimplementedProjectsServiceServer) mustEmbedUnimplementedProjectsServiceServer() {}
func (UnimplementedProjectsServiceServer) testEmbeddedByValue()                         {}

type UnsafeProjectsServiceServer interface {
	mustEmbedUnimplementedProjectsServiceServer()
}

func RegisterProjectsServiceServer(s grpc.ServiceRegistrar, srv ProjectsServiceServer) {

	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProjectsService_ServiceDesc, srv)
}

func _ProjectsService_CreateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).CreateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_CreateProject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).CreateProject(ctx, req.(*CreateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_GetProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).GetProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_GetProject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).GetProject(ctx, req.(*GetProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_ListProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProjectsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).ListProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_ListProjects_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).ListProjects(ctx, req.(*ListProjectsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_VerifyConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyConnectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).VerifyConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_VerifyConnection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).VerifyConnection(ctx, req.(*VerifyConnectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_DeleteProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).DeleteProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_DeleteProject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).DeleteProject(ctx, req.(*DeleteProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_UpdateProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateProjectRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).UpdateProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_UpdateProject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).UpdateProject(ctx, req.(*UpdateProjectRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_Health_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).Health(ctx, req.(*HealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_GetEnvVars_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetEnvVarsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).GetEnvVars(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_GetEnvVars_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).GetEnvVars(ctx, req.(*GetEnvVarsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_UpdateEnvVars_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEnvVarsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).UpdateEnvVars(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_UpdateEnvVars_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).UpdateEnvVars(ctx, req.(*UpdateEnvVarsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_GetPipelineYAML_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPipelineYAMLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).GetPipelineYAML(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_GetPipelineYAML_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).GetPipelineYAML(ctx, req.(*GetPipelineYAMLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_SetPipelineYAML_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetPipelineYAMLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).SetPipelineYAML(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_SetPipelineYAML_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).SetPipelineYAML(ctx, req.(*SetPipelineYAMLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_ListSecrets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListSecretsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).ListSecrets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_ListSecrets_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).ListSecrets(ctx, req.(*ListSecretsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_SetSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).SetSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_SetSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).SetSecret(ctx, req.(*SetSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_DeleteSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).DeleteSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_DeleteSecret_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).DeleteSecret(ctx, req.(*DeleteSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_GetProjectInternal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProjectInternalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).GetProjectInternal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_GetProjectInternal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).GetProjectInternal(ctx, req.(*GetProjectInternalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_ListProjectsInternal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProjectsInternalRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).ListProjectsInternal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_ListProjectsInternal_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).ListProjectsInternal(ctx, req.(*ListProjectsInternalRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_ListMembers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMembersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).ListMembers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_ListMembers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).ListMembers(ctx, req.(*ListMembersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_InviteMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InviteMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).InviteMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_InviteMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).InviteMember(ctx, req.(*InviteMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_UpdateMemberRole_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMemberRoleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).UpdateMemberRole(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_UpdateMemberRole_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).UpdateMemberRole(ctx, req.(*UpdateMemberRoleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectsService_RemoveMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectsServiceServer).RemoveMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectsService_RemoveMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectsServiceServer).RemoveMember(ctx, req.(*RemoveMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var ProjectsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "projects.ProjectsService",
	HandlerType: (*ProjectsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateProject",
			Handler:    _ProjectsService_CreateProject_Handler,
		},
		{
			MethodName: "GetProject",
			Handler:    _ProjectsService_GetProject_Handler,
		},
		{
			MethodName: "ListProjects",
			Handler:    _ProjectsService_ListProjects_Handler,
		},
		{
			MethodName: "VerifyConnection",
			Handler:    _ProjectsService_VerifyConnection_Handler,
		},
		{
			MethodName: "DeleteProject",
			Handler:    _ProjectsService_DeleteProject_Handler,
		},
		{
			MethodName: "UpdateProject",
			Handler:    _ProjectsService_UpdateProject_Handler,
		},
		{
			MethodName: "Health",
			Handler:    _ProjectsService_Health_Handler,
		},
		{
			MethodName: "GetEnvVars",
			Handler:    _ProjectsService_GetEnvVars_Handler,
		},
		{
			MethodName: "UpdateEnvVars",
			Handler:    _ProjectsService_UpdateEnvVars_Handler,
		},
		{
			MethodName: "GetPipelineYAML",
			Handler:    _ProjectsService_GetPipelineYAML_Handler,
		},
		{
			MethodName: "SetPipelineYAML",
			Handler:    _ProjectsService_SetPipelineYAML_Handler,
		},
		{
			MethodName: "ListSecrets",
			Handler:    _ProjectsService_ListSecrets_Handler,
		},
		{
			MethodName: "SetSecret",
			Handler:    _ProjectsService_SetSecret_Handler,
		},
		{
			MethodName: "DeleteSecret",
			Handler:    _ProjectsService_DeleteSecret_Handler,
		},
		{
			MethodName: "GetProjectInternal",
			Handler:    _ProjectsService_GetProjectInternal_Handler,
		},
		{
			MethodName: "ListProjectsInternal",
			Handler:    _ProjectsService_ListProjectsInternal_Handler,
		},
		{
			MethodName: "ListMembers",
			Handler:    _ProjectsService_ListMembers_Handler,
		},
		{
			MethodName: "InviteMember",
			Handler:    _ProjectsService_InviteMember_Handler,
		},
		{
			MethodName: "UpdateMemberRole",
			Handler:    _ProjectsService_UpdateMemberRole_Handler,
		},
		{
			MethodName: "RemoveMember",
			Handler:    _ProjectsService_RemoveMember_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "projects.proto",
}
