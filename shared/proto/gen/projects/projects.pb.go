package projects

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)

	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EnvVar struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         string                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *EnvVar) Reset() {
	*x = EnvVar{}
	mi := &file_projects_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EnvVar) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EnvVar) ProtoMessage() {}

func (x *EnvVar) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*EnvVar) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{0}
}

func (x *EnvVar) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *EnvVar) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type CreateProjectRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	RepositoryUrl string                 `protobuf:"bytes,3,opt,name=repository_url,json=repositoryUrl,proto3" json:"repository_url,omitempty"`
	DefaultBranch string                 `protobuf:"bytes,4,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateProjectRequest) Reset() {
	*x = CreateProjectRequest{}
	mi := &file_projects_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateProjectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateProjectRequest) ProtoMessage() {}

func (x *CreateProjectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*CreateProjectRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{1}
}

func (x *CreateProjectRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *CreateProjectRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateProjectRequest) GetRepositoryUrl() string {
	if x != nil {
		return x.RepositoryUrl
	}
	return ""
}

func (x *CreateProjectRequest) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

type CreateProjectResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	RepoUrl       string                 `protobuf:"bytes,3,opt,name=repo_url,json=repoUrl,proto3" json:"repo_url,omitempty"`
	DefaultBranch string                 `protobuf:"bytes,4,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	PublicKey     string                 `protobuf:"bytes,5,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	WebhookSecret string                 `protobuf:"bytes,6,opt,name=webhook_secret,json=webhookSecret,proto3" json:"webhook_secret,omitempty"`
	WebhookUrl    string                 `protobuf:"bytes,7,opt,name=webhook_url,json=webhookUrl,proto3" json:"webhook_url,omitempty"`
	Status        string                 `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	OwnerUserId   string                 `protobuf:"bytes,9,opt,name=owner_user_id,json=ownerUserId,proto3" json:"owner_user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateProjectResponse) Reset() {
	*x = CreateProjectResponse{}
	mi := &file_projects_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateProjectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateProjectResponse) ProtoMessage() {}

func (x *CreateProjectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*CreateProjectResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{2}
}

func (x *CreateProjectResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CreateProjectResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreateProjectResponse) GetRepoUrl() string {
	if x != nil {
		return x.RepoUrl
	}
	return ""
}

func (x *CreateProjectResponse) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

func (x *CreateProjectResponse) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

func (x *CreateProjectResponse) GetWebhookSecret() string {
	if x != nil {
		return x.WebhookSecret
	}
	return ""
}

func (x *CreateProjectResponse) GetWebhookUrl() string {
	if x != nil {
		return x.WebhookUrl
	}
	return ""
}

func (x *CreateProjectResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *CreateProjectResponse) GetOwnerUserId() string {
	if x != nil {
		return x.OwnerUserId
	}
	return ""
}

type GetProjectRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetProjectRequest) Reset() {
	*x = GetProjectRequest{}
	mi := &file_projects_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProjectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProjectRequest) ProtoMessage() {}

func (x *GetProjectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetProjectRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{3}
}

func (x *GetProjectRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetProjectRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type ListProjectsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListProjectsRequest) Reset() {
	*x = ListProjectsRequest{}
	mi := &file_projects_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListProjectsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProjectsRequest) ProtoMessage() {}

func (x *ListProjectsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListProjectsRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{4}
}

func (x *ListProjectsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type ListProjectsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Projects      []*ProjectSummary      `protobuf:"bytes,1,rep,name=projects,proto3" json:"projects,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListProjectsResponse) Reset() {
	*x = ListProjectsResponse{}
	mi := &file_projects_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListProjectsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProjectsResponse) ProtoMessage() {}

func (x *ListProjectsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListProjectsResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{5}
}

func (x *ListProjectsResponse) GetProjects() []*ProjectSummary {
	if x != nil {
		return x.Projects
	}
	return nil
}

type ProjectSummary struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	RepoUrl       string                 `protobuf:"bytes,3,opt,name=repo_url,json=repoUrl,proto3" json:"repo_url,omitempty"`
	DefaultBranch string                 `protobuf:"bytes,4,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	Status        string                 `protobuf:"bytes,5,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProjectSummary) Reset() {
	*x = ProjectSummary{}
	mi := &file_projects_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProjectSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectSummary) ProtoMessage() {}

func (x *ProjectSummary) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ProjectSummary) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{6}
}

func (x *ProjectSummary) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ProjectSummary) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProjectSummary) GetRepoUrl() string {
	if x != nil {
		return x.RepoUrl
	}
	return ""
}

func (x *ProjectSummary) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

func (x *ProjectSummary) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type GetProjectResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	RepoUrl       string                 `protobuf:"bytes,3,opt,name=repo_url,json=repoUrl,proto3" json:"repo_url,omitempty"`
	DefaultBranch string                 `protobuf:"bytes,4,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	PublicKey     string                 `protobuf:"bytes,5,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	WebhookSecret string                 `protobuf:"bytes,6,opt,name=webhook_secret,json=webhookSecret,proto3" json:"webhook_secret,omitempty"`
	WebhookUrl    string                 `protobuf:"bytes,7,opt,name=webhook_url,json=webhookUrl,proto3" json:"webhook_url,omitempty"`
	Status        string                 `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	OwnerUserId   string                 `protobuf:"bytes,9,opt,name=owner_user_id,json=ownerUserId,proto3" json:"owner_user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetProjectResponse) Reset() {
	*x = GetProjectResponse{}
	mi := &file_projects_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProjectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProjectResponse) ProtoMessage() {}

func (x *GetProjectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetProjectResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{7}
}

func (x *GetProjectResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetProjectResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetProjectResponse) GetRepoUrl() string {
	if x != nil {
		return x.RepoUrl
	}
	return ""
}

func (x *GetProjectResponse) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

func (x *GetProjectResponse) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

func (x *GetProjectResponse) GetWebhookSecret() string {
	if x != nil {
		return x.WebhookSecret
	}
	return ""
}

func (x *GetProjectResponse) GetWebhookUrl() string {
	if x != nil {
		return x.WebhookUrl
	}
	return ""
}

func (x *GetProjectResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *GetProjectResponse) GetOwnerUserId() string {
	if x != nil {
		return x.OwnerUserId
	}
	return ""
}

type VerifyConnectionRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VerifyConnectionRequest) Reset() {
	*x = VerifyConnectionRequest{}
	mi := &file_projects_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VerifyConnectionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyConnectionRequest) ProtoMessage() {}

func (x *VerifyConnectionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*VerifyConnectionRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{8}
}

func (x *VerifyConnectionRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *VerifyConnectionRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type VerifyConnectionResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	Error         string                 `protobuf:"bytes,3,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VerifyConnectionResponse) Reset() {
	*x = VerifyConnectionResponse{}
	mi := &file_projects_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VerifyConnectionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifyConnectionResponse) ProtoMessage() {}

func (x *VerifyConnectionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*VerifyConnectionResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{9}
}

func (x *VerifyConnectionResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *VerifyConnectionResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *VerifyConnectionResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

type DeleteProjectRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteProjectRequest) Reset() {
	*x = DeleteProjectRequest{}
	mi := &file_projects_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteProjectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteProjectRequest) ProtoMessage() {}

func (x *DeleteProjectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*DeleteProjectRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{10}
}

func (x *DeleteProjectRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *DeleteProjectRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type DeleteProjectResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteProjectResponse) Reset() {
	*x = DeleteProjectResponse{}
	mi := &file_projects_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteProjectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteProjectResponse) ProtoMessage() {}

func (x *DeleteProjectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*DeleteProjectResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{11}
}

type UpdateProjectRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Name          string                 `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	DefaultBranch string                 `protobuf:"bytes,4,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateProjectRequest) Reset() {
	*x = UpdateProjectRequest{}
	mi := &file_projects_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateProjectRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProjectRequest) ProtoMessage() {}

func (x *UpdateProjectRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*UpdateProjectRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{12}
}

func (x *UpdateProjectRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UpdateProjectRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *UpdateProjectRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateProjectRequest) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

type UpdateProjectResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	RepoUrl       string                 `protobuf:"bytes,3,opt,name=repo_url,json=repoUrl,proto3" json:"repo_url,omitempty"`
	DefaultBranch string                 `protobuf:"bytes,4,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	PublicKey     string                 `protobuf:"bytes,5,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	WebhookSecret string                 `protobuf:"bytes,6,opt,name=webhook_secret,json=webhookSecret,proto3" json:"webhook_secret,omitempty"`
	WebhookUrl    string                 `protobuf:"bytes,7,opt,name=webhook_url,json=webhookUrl,proto3" json:"webhook_url,omitempty"`
	Status        string                 `protobuf:"bytes,8,opt,name=status,proto3" json:"status,omitempty"`
	OwnerUserId   string                 `protobuf:"bytes,9,opt,name=owner_user_id,json=ownerUserId,proto3" json:"owner_user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateProjectResponse) Reset() {
	*x = UpdateProjectResponse{}
	mi := &file_projects_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateProjectResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateProjectResponse) ProtoMessage() {}

func (x *UpdateProjectResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*UpdateProjectResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{13}
}

func (x *UpdateProjectResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateProjectResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateProjectResponse) GetRepoUrl() string {
	if x != nil {
		return x.RepoUrl
	}
	return ""
}

func (x *UpdateProjectResponse) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

func (x *UpdateProjectResponse) GetPublicKey() string {
	if x != nil {
		return x.PublicKey
	}
	return ""
}

func (x *UpdateProjectResponse) GetWebhookSecret() string {
	if x != nil {
		return x.WebhookSecret
	}
	return ""
}

func (x *UpdateProjectResponse) GetWebhookUrl() string {
	if x != nil {
		return x.WebhookUrl
	}
	return ""
}

func (x *UpdateProjectResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *UpdateProjectResponse) GetOwnerUserId() string {
	if x != nil {
		return x.OwnerUserId
	}
	return ""
}

type GetEnvVarsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetEnvVarsRequest) Reset() {
	*x = GetEnvVarsRequest{}
	mi := &file_projects_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetEnvVarsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEnvVarsRequest) ProtoMessage() {}

func (x *GetEnvVarsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetEnvVarsRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{14}
}

func (x *GetEnvVarsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetEnvVarsRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type GetEnvVarsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Vars          []*EnvVar              `protobuf:"bytes,1,rep,name=vars,proto3" json:"vars,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetEnvVarsResponse) Reset() {
	*x = GetEnvVarsResponse{}
	mi := &file_projects_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetEnvVarsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetEnvVarsResponse) ProtoMessage() {}

func (x *GetEnvVarsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetEnvVarsResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{15}
}

func (x *GetEnvVarsResponse) GetVars() []*EnvVar {
	if x != nil {
		return x.Vars
	}
	return nil
}

type UpdateEnvVarsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Vars          []*EnvVar              `protobuf:"bytes,3,rep,name=vars,proto3" json:"vars,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateEnvVarsRequest) Reset() {
	*x = UpdateEnvVarsRequest{}
	mi := &file_projects_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateEnvVarsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateEnvVarsRequest) ProtoMessage() {}

func (x *UpdateEnvVarsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*UpdateEnvVarsRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{16}
}

func (x *UpdateEnvVarsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *UpdateEnvVarsRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *UpdateEnvVarsRequest) GetVars() []*EnvVar {
	if x != nil {
		return x.Vars
	}
	return nil
}

type UpdateEnvVarsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateEnvVarsResponse) Reset() {
	*x = UpdateEnvVarsResponse{}
	mi := &file_projects_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateEnvVarsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateEnvVarsResponse) ProtoMessage() {}

func (x *UpdateEnvVarsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*UpdateEnvVarsResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{17}
}

type GetPipelineYAMLRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPipelineYAMLRequest) Reset() {
	*x = GetPipelineYAMLRequest{}
	mi := &file_projects_proto_msgTypes[18]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPipelineYAMLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPipelineYAMLRequest) ProtoMessage() {}

func (x *GetPipelineYAMLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[18]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetPipelineYAMLRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{18}
}

func (x *GetPipelineYAMLRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *GetPipelineYAMLRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type GetPipelineYAMLResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`

	Yaml          string `protobuf:"bytes,1,opt,name=yaml,proto3" json:"yaml,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetPipelineYAMLResponse) Reset() {
	*x = GetPipelineYAMLResponse{}
	mi := &file_projects_proto_msgTypes[19]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetPipelineYAMLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetPipelineYAMLResponse) ProtoMessage() {}

func (x *GetPipelineYAMLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[19]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetPipelineYAMLResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{19}
}

func (x *GetPipelineYAMLResponse) GetYaml() string {
	if x != nil {
		return x.Yaml
	}
	return ""
}

type SetPipelineYAMLRequest struct {
	state     protoimpl.MessageState `protogen:"open.v1"`
	UserId    string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`

	Yaml          string `protobuf:"bytes,3,opt,name=yaml,proto3" json:"yaml,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetPipelineYAMLRequest) Reset() {
	*x = SetPipelineYAMLRequest{}
	mi := &file_projects_proto_msgTypes[20]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetPipelineYAMLRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetPipelineYAMLRequest) ProtoMessage() {}

func (x *SetPipelineYAMLRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[20]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SetPipelineYAMLRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{20}
}

func (x *SetPipelineYAMLRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *SetPipelineYAMLRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *SetPipelineYAMLRequest) GetYaml() string {
	if x != nil {
		return x.Yaml
	}
	return ""
}

type SetPipelineYAMLResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetPipelineYAMLResponse) Reset() {
	*x = SetPipelineYAMLResponse{}
	mi := &file_projects_proto_msgTypes[21]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetPipelineYAMLResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetPipelineYAMLResponse) ProtoMessage() {}

func (x *SetPipelineYAMLResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[21]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SetPipelineYAMLResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{21}
}

type SecretMeta struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	UpdatedAt     string                 `protobuf:"bytes,2,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SecretMeta) Reset() {
	*x = SecretMeta{}
	mi := &file_projects_proto_msgTypes[22]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SecretMeta) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretMeta) ProtoMessage() {}

func (x *SecretMeta) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[22]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SecretMeta) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{22}
}

func (x *SecretMeta) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SecretMeta) GetUpdatedAt() string {
	if x != nil {
		return x.UpdatedAt
	}
	return ""
}

type ListSecretsRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSecretsRequest) Reset() {
	*x = ListSecretsRequest{}
	mi := &file_projects_proto_msgTypes[23]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSecretsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSecretsRequest) ProtoMessage() {}

func (x *ListSecretsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[23]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListSecretsRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{23}
}

func (x *ListSecretsRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ListSecretsRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type ListSecretsResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Secrets       []*SecretMeta          `protobuf:"bytes,1,rep,name=secrets,proto3" json:"secrets,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListSecretsResponse) Reset() {
	*x = ListSecretsResponse{}
	mi := &file_projects_proto_msgTypes[24]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListSecretsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListSecretsResponse) ProtoMessage() {}

func (x *ListSecretsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[24]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListSecretsResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{24}
}

func (x *ListSecretsResponse) GetSecrets() []*SecretMeta {
	if x != nil {
		return x.Secrets
	}
	return nil
}

type SetSecretRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Key           string                 `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	Value         string                 `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetSecretRequest) Reset() {
	*x = SetSecretRequest{}
	mi := &file_projects_proto_msgTypes[25]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetSecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetSecretRequest) ProtoMessage() {}

func (x *SetSecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[25]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SetSecretRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{25}
}

func (x *SetSecretRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *SetSecretRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *SetSecretRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SetSecretRequest) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type SetSecretResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SetSecretResponse) Reset() {
	*x = SetSecretResponse{}
	mi := &file_projects_proto_msgTypes[26]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SetSecretResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SetSecretResponse) ProtoMessage() {}

func (x *SetSecretResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[26]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SetSecretResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{26}
}

type DeleteSecretRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	Key           string                 `protobuf:"bytes,3,opt,name=key,proto3" json:"key,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteSecretRequest) Reset() {
	*x = DeleteSecretRequest{}
	mi := &file_projects_proto_msgTypes[27]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteSecretRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSecretRequest) ProtoMessage() {}

func (x *DeleteSecretRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[27]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*DeleteSecretRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{27}
}

func (x *DeleteSecretRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *DeleteSecretRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *DeleteSecretRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type DeleteSecretResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteSecretResponse) Reset() {
	*x = DeleteSecretResponse{}
	mi := &file_projects_proto_msgTypes[28]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteSecretResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteSecretResponse) ProtoMessage() {}

func (x *DeleteSecretResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[28]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*DeleteSecretResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{28}
}

type SecretKV struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Key           string                 `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value         string                 `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SecretKV) Reset() {
	*x = SecretKV{}
	mi := &file_projects_proto_msgTypes[29]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SecretKV) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SecretKV) ProtoMessage() {}

func (x *SecretKV) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[29]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*SecretKV) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{29}
}

func (x *SecretKV) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *SecretKV) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

type HealthRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealthRequest) Reset() {
	*x = HealthRequest{}
	mi := &file_projects_proto_msgTypes[30]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthRequest) ProtoMessage() {}

func (x *HealthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[30]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*HealthRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{30}
}

type HealthResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Healthy       bool                   `protobuf:"varint,1,opt,name=healthy,proto3" json:"healthy,omitempty"`
	Status        string                 `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HealthResponse) Reset() {
	*x = HealthResponse{}
	mi := &file_projects_proto_msgTypes[31]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HealthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HealthResponse) ProtoMessage() {}

func (x *HealthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[31]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*HealthResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{31}
}

func (x *HealthResponse) GetHealthy() bool {
	if x != nil {
		return x.Healthy
	}
	return false
}

func (x *HealthResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

type ProjectMember struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Email         string                 `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	DisplayName   string                 `protobuf:"bytes,3,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	Role          string                 `protobuf:"bytes,4,opt,name=role,proto3" json:"role,omitempty"`
	InvitedBy     string                 `protobuf:"bytes,5,opt,name=invited_by,json=invitedBy,proto3" json:"invited_by,omitempty"`
	CreatedAt     string                 `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProjectMember) Reset() {
	*x = ProjectMember{}
	mi := &file_projects_proto_msgTypes[32]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProjectMember) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectMember) ProtoMessage() {}

func (x *ProjectMember) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[32]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ProjectMember) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{32}
}

func (x *ProjectMember) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ProjectMember) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *ProjectMember) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *ProjectMember) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

func (x *ProjectMember) GetInvitedBy() string {
	if x != nil {
		return x.InvitedBy
	}
	return ""
}

func (x *ProjectMember) GetCreatedAt() string {
	if x != nil {
		return x.CreatedAt
	}
	return ""
}

type ListMembersRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListMembersRequest) Reset() {
	*x = ListMembersRequest{}
	mi := &file_projects_proto_msgTypes[33]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListMembersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMembersRequest) ProtoMessage() {}

func (x *ListMembersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[33]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListMembersRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{33}
}

func (x *ListMembersRequest) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ListMembersRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type ListMembersResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Members       []*ProjectMember       `protobuf:"bytes,1,rep,name=members,proto3" json:"members,omitempty"`
	RequesterRole string                 `protobuf:"bytes,2,opt,name=requester_role,json=requesterRole,proto3" json:"requester_role,omitempty"`
	OwnerUserId   string                 `protobuf:"bytes,3,opt,name=owner_user_id,json=ownerUserId,proto3" json:"owner_user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListMembersResponse) Reset() {
	*x = ListMembersResponse{}
	mi := &file_projects_proto_msgTypes[34]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListMembersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMembersResponse) ProtoMessage() {}

func (x *ListMembersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[34]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListMembersResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{34}
}

func (x *ListMembersResponse) GetMembers() []*ProjectMember {
	if x != nil {
		return x.Members
	}
	return nil
}

func (x *ListMembersResponse) GetRequesterRole() string {
	if x != nil {
		return x.RequesterRole
	}
	return ""
}

func (x *ListMembersResponse) GetOwnerUserId() string {
	if x != nil {
		return x.OwnerUserId
	}
	return ""
}

type InviteMemberRequest struct {
	state              protoimpl.MessageState `protogen:"open.v1"`
	OwnerUserId        string                 `protobuf:"bytes,1,opt,name=owner_user_id,json=ownerUserId,proto3" json:"owner_user_id,omitempty"`
	ProjectId          string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	InviteeUserId      string                 `protobuf:"bytes,3,opt,name=invitee_user_id,json=inviteeUserId,proto3" json:"invitee_user_id,omitempty"`
	InviteeEmail       string                 `protobuf:"bytes,4,opt,name=invitee_email,json=inviteeEmail,proto3" json:"invitee_email,omitempty"`
	InviteeDisplayName string                 `protobuf:"bytes,5,opt,name=invitee_display_name,json=inviteeDisplayName,proto3" json:"invitee_display_name,omitempty"`
	Role               string                 `protobuf:"bytes,6,opt,name=role,proto3" json:"role,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *InviteMemberRequest) Reset() {
	*x = InviteMemberRequest{}
	mi := &file_projects_proto_msgTypes[35]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InviteMemberRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InviteMemberRequest) ProtoMessage() {}

func (x *InviteMemberRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[35]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*InviteMemberRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{35}
}

func (x *InviteMemberRequest) GetOwnerUserId() string {
	if x != nil {
		return x.OwnerUserId
	}
	return ""
}

func (x *InviteMemberRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *InviteMemberRequest) GetInviteeUserId() string {
	if x != nil {
		return x.InviteeUserId
	}
	return ""
}

func (x *InviteMemberRequest) GetInviteeEmail() string {
	if x != nil {
		return x.InviteeEmail
	}
	return ""
}

func (x *InviteMemberRequest) GetInviteeDisplayName() string {
	if x != nil {
		return x.InviteeDisplayName
	}
	return ""
}

func (x *InviteMemberRequest) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

type InviteMemberResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Member        *ProjectMember         `protobuf:"bytes,1,opt,name=member,proto3" json:"member,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *InviteMemberResponse) Reset() {
	*x = InviteMemberResponse{}
	mi := &file_projects_proto_msgTypes[36]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *InviteMemberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InviteMemberResponse) ProtoMessage() {}

func (x *InviteMemberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[36]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*InviteMemberResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{36}
}

func (x *InviteMemberResponse) GetMember() *ProjectMember {
	if x != nil {
		return x.Member
	}
	return nil
}

type UpdateMemberRoleRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OwnerUserId   string                 `protobuf:"bytes,1,opt,name=owner_user_id,json=ownerUserId,proto3" json:"owner_user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	TargetUserId  string                 `protobuf:"bytes,3,opt,name=target_user_id,json=targetUserId,proto3" json:"target_user_id,omitempty"`
	Role          string                 `protobuf:"bytes,4,opt,name=role,proto3" json:"role,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateMemberRoleRequest) Reset() {
	*x = UpdateMemberRoleRequest{}
	mi := &file_projects_proto_msgTypes[37]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateMemberRoleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMemberRoleRequest) ProtoMessage() {}

func (x *UpdateMemberRoleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[37]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*UpdateMemberRoleRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{37}
}

func (x *UpdateMemberRoleRequest) GetOwnerUserId() string {
	if x != nil {
		return x.OwnerUserId
	}
	return ""
}

func (x *UpdateMemberRoleRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *UpdateMemberRoleRequest) GetTargetUserId() string {
	if x != nil {
		return x.TargetUserId
	}
	return ""
}

func (x *UpdateMemberRoleRequest) GetRole() string {
	if x != nil {
		return x.Role
	}
	return ""
}

type UpdateMemberRoleResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Member        *ProjectMember         `protobuf:"bytes,1,opt,name=member,proto3" json:"member,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateMemberRoleResponse) Reset() {
	*x = UpdateMemberRoleResponse{}
	mi := &file_projects_proto_msgTypes[38]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateMemberRoleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMemberRoleResponse) ProtoMessage() {}

func (x *UpdateMemberRoleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[38]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*UpdateMemberRoleResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{38}
}

func (x *UpdateMemberRoleResponse) GetMember() *ProjectMember {
	if x != nil {
		return x.Member
	}
	return nil
}

type RemoveMemberRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	OwnerUserId   string                 `protobuf:"bytes,1,opt,name=owner_user_id,json=ownerUserId,proto3" json:"owner_user_id,omitempty"`
	ProjectId     string                 `protobuf:"bytes,2,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	TargetUserId  string                 `protobuf:"bytes,3,opt,name=target_user_id,json=targetUserId,proto3" json:"target_user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveMemberRequest) Reset() {
	*x = RemoveMemberRequest{}
	mi := &file_projects_proto_msgTypes[39]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveMemberRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveMemberRequest) ProtoMessage() {}

func (x *RemoveMemberRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[39]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*RemoveMemberRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{39}
}

func (x *RemoveMemberRequest) GetOwnerUserId() string {
	if x != nil {
		return x.OwnerUserId
	}
	return ""
}

func (x *RemoveMemberRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

func (x *RemoveMemberRequest) GetTargetUserId() string {
	if x != nil {
		return x.TargetUserId
	}
	return ""
}

type RemoveMemberResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RemoveMemberResponse) Reset() {
	*x = RemoveMemberResponse{}
	mi := &file_projects_proto_msgTypes[40]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RemoveMemberResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveMemberResponse) ProtoMessage() {}

func (x *RemoveMemberResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[40]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*RemoveMemberResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{40}
}

type GetProjectInternalRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	ProjectId     string                 `protobuf:"bytes,1,opt,name=project_id,json=projectId,proto3" json:"project_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetProjectInternalRequest) Reset() {
	*x = GetProjectInternalRequest{}
	mi := &file_projects_proto_msgTypes[41]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProjectInternalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProjectInternalRequest) ProtoMessage() {}

func (x *GetProjectInternalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[41]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetProjectInternalRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{41}
}

func (x *GetProjectInternalRequest) GetProjectId() string {
	if x != nil {
		return x.ProjectId
	}
	return ""
}

type GetProjectInternalResponse struct {
	state                protoimpl.MessageState `protogen:"open.v1"`
	Id                   string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	RepoUrl              string                 `protobuf:"bytes,2,opt,name=repo_url,json=repoUrl,proto3" json:"repo_url,omitempty"`
	WebhookSecret        string                 `protobuf:"bytes,3,opt,name=webhook_secret,json=webhookSecret,proto3" json:"webhook_secret,omitempty"`
	SshKey               string                 `protobuf:"bytes,4,opt,name=ssh_key,json=sshKey,proto3" json:"ssh_key,omitempty"`
	EnvVars              []*EnvVar              `protobuf:"bytes,5,rep,name=env_vars,json=envVars,proto3" json:"env_vars,omitempty"`
	PipelineYamlOverride string                 `protobuf:"bytes,6,opt,name=pipeline_yaml_override,json=pipelineYamlOverride,proto3" json:"pipeline_yaml_override,omitempty"`
	Secrets              []*SecretKV            `protobuf:"bytes,7,rep,name=secrets,proto3" json:"secrets,omitempty"`
	DefaultBranch        string                 `protobuf:"bytes,8,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *GetProjectInternalResponse) Reset() {
	*x = GetProjectInternalResponse{}
	mi := &file_projects_proto_msgTypes[42]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProjectInternalResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProjectInternalResponse) ProtoMessage() {}

func (x *GetProjectInternalResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[42]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*GetProjectInternalResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{42}
}

func (x *GetProjectInternalResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetProjectInternalResponse) GetRepoUrl() string {
	if x != nil {
		return x.RepoUrl
	}
	return ""
}

func (x *GetProjectInternalResponse) GetWebhookSecret() string {
	if x != nil {
		return x.WebhookSecret
	}
	return ""
}

func (x *GetProjectInternalResponse) GetSshKey() string {
	if x != nil {
		return x.SshKey
	}
	return ""
}

func (x *GetProjectInternalResponse) GetEnvVars() []*EnvVar {
	if x != nil {
		return x.EnvVars
	}
	return nil
}

func (x *GetProjectInternalResponse) GetPipelineYamlOverride() string {
	if x != nil {
		return x.PipelineYamlOverride
	}
	return ""
}

func (x *GetProjectInternalResponse) GetSecrets() []*SecretKV {
	if x != nil {
		return x.Secrets
	}
	return nil
}

func (x *GetProjectInternalResponse) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

type ListProjectsInternalRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListProjectsInternalRequest) Reset() {
	*x = ListProjectsInternalRequest{}
	mi := &file_projects_proto_msgTypes[43]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListProjectsInternalRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProjectsInternalRequest) ProtoMessage() {}

func (x *ListProjectsInternalRequest) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[43]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListProjectsInternalRequest) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{43}
}

type ProjectInternalSummary struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	DefaultBranch string                 `protobuf:"bytes,2,opt,name=default_branch,json=defaultBranch,proto3" json:"default_branch,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProjectInternalSummary) Reset() {
	*x = ProjectInternalSummary{}
	mi := &file_projects_proto_msgTypes[44]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProjectInternalSummary) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectInternalSummary) ProtoMessage() {}

func (x *ProjectInternalSummary) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[44]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ProjectInternalSummary) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{44}
}

func (x *ProjectInternalSummary) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ProjectInternalSummary) GetDefaultBranch() string {
	if x != nil {
		return x.DefaultBranch
	}
	return ""
}

type ListProjectsInternalResponse struct {
	state         protoimpl.MessageState    `protogen:"open.v1"`
	Projects      []*ProjectInternalSummary `protobuf:"bytes,1,rep,name=projects,proto3" json:"projects,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListProjectsInternalResponse) Reset() {
	*x = ListProjectsInternalResponse{}
	mi := &file_projects_proto_msgTypes[45]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListProjectsInternalResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListProjectsInternalResponse) ProtoMessage() {}

func (x *ListProjectsInternalResponse) ProtoReflect() protoreflect.Message {
	mi := &file_projects_proto_msgTypes[45]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListProjectsInternalResponse) Descriptor() ([]byte, []int) {
	return file_projects_proto_rawDescGZIP(), []int{45}
}

func (x *ListProjectsInternalResponse) GetProjects() []*ProjectInternalSummary {
	if x != nil {
		return x.Projects
	}
	return nil
}

var File_projects_proto protoreflect.FileDescriptor

const file_projects_proto_rawDesc = "" +
	"\n" +
	"\x0eprojects.proto\x12\bprojects\"0\n" +
	"\x06EnvVar\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value\"\x91\x01\n" +
	"\x14CreateProjectRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12%\n" +
	"\x0erepository_url\x18\x03 \x01(\tR\rrepositoryUrl\x12%\n" +
	"\x0edefault_branch\x18\x04 \x01(\tR\rdefaultBranch\"\xa0\x02\n" +
	"\x15CreateProjectResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x19\n" +
	"\brepo_url\x18\x03 \x01(\tR\arepoUrl\x12%\n" +
	"\x0edefault_branch\x18\x04 \x01(\tR\rdefaultBranch\x12\x1d\n" +
	"\n" +
	"public_key\x18\x05 \x01(\tR\tpublicKey\x12%\n" +
	"\x0ewebhook_secret\x18\x06 \x01(\tR\rwebhookSecret\x12\x1f\n" +
	"\vwebhook_url\x18\a \x01(\tR\n" +
	"webhookUrl\x12\x16\n" +
	"\x06status\x18\b \x01(\tR\x06status\x12\"\n" +
	"\rowner_user_id\x18\t \x01(\tR\vownerUserId\"K\n" +
	"\x11GetProjectRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\".\n" +
	"\x13ListProjectsRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\"L\n" +
	"\x14ListProjectsResponse\x124\n" +
	"\bprojects\x18\x01 \x03(\v2\x18.projects.ProjectSummaryR\bprojects\"\x8e\x01\n" +
	"\x0eProjectSummary\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x19\n" +
	"\brepo_url\x18\x03 \x01(\tR\arepoUrl\x12%\n" +
	"\x0edefault_branch\x18\x04 \x01(\tR\rdefaultBranch\x12\x16\n" +
	"\x06status\x18\x05 \x01(\tR\x06status\"\x9d\x02\n" +
	"\x12GetProjectResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x19\n" +
	"\brepo_url\x18\x03 \x01(\tR\arepoUrl\x12%\n" +
	"\x0edefault_branch\x18\x04 \x01(\tR\rdefaultBranch\x12\x1d\n" +
	"\n" +
	"public_key\x18\x05 \x01(\tR\tpublicKey\x12%\n" +
	"\x0ewebhook_secret\x18\x06 \x01(\tR\rwebhookSecret\x12\x1f\n" +
	"\vwebhook_url\x18\a \x01(\tR\n" +
	"webhookUrl\x12\x16\n" +
	"\x06status\x18\b \x01(\tR\x06status\x12\"\n" +
	"\rowner_user_id\x18\t \x01(\tR\vownerUserId\"Q\n" +
	"\x17VerifyConnectionRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\"b\n" +
	"\x18VerifyConnectionResponse\x12\x18\n" +
	"\asuccess\x18\x01 \x01(\bR\asuccess\x12\x16\n" +
	"\x06status\x18\x02 \x01(\tR\x06status\x12\x14\n" +
	"\x05error\x18\x03 \x01(\tR\x05error\"N\n" +
	"\x14DeleteProjectRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\"\x17\n" +
	"\x15DeleteProjectResponse\"\x89\x01\n" +
	"\x14UpdateProjectRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\x12\x12\n" +
	"\x04name\x18\x03 \x01(\tR\x04name\x12%\n" +
	"\x0edefault_branch\x18\x04 \x01(\tR\rdefaultBranch\"\xa0\x02\n" +
	"\x15UpdateProjectResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12\x19\n" +
	"\brepo_url\x18\x03 \x01(\tR\arepoUrl\x12%\n" +
	"\x0edefault_branch\x18\x04 \x01(\tR\rdefaultBranch\x12\x1d\n" +
	"\n" +
	"public_key\x18\x05 \x01(\tR\tpublicKey\x12%\n" +
	"\x0ewebhook_secret\x18\x06 \x01(\tR\rwebhookSecret\x12\x1f\n" +
	"\vwebhook_url\x18\a \x01(\tR\n" +
	"webhookUrl\x12\x16\n" +
	"\x06status\x18\b \x01(\tR\x06status\x12\"\n" +
	"\rowner_user_id\x18\t \x01(\tR\vownerUserId\"K\n" +
	"\x11GetEnvVarsRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\":\n" +
	"\x12GetEnvVarsResponse\x12$\n" +
	"\x04vars\x18\x01 \x03(\v2\x10.projects.EnvVarR\x04vars\"t\n" +
	"\x14UpdateEnvVarsRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\x12$\n" +
	"\x04vars\x18\x03 \x03(\v2\x10.projects.EnvVarR\x04vars\"\x17\n" +
	"\x15UpdateEnvVarsResponse\"P\n" +
	"\x16GetPipelineYAMLRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\"-\n" +
	"\x17GetPipelineYAMLResponse\x12\x12\n" +
	"\x04yaml\x18\x01 \x01(\tR\x04yaml\"d\n" +
	"\x16SetPipelineYAMLRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\x12\x12\n" +
	"\x04yaml\x18\x03 \x01(\tR\x04yaml\"\x19\n" +
	"\x17SetPipelineYAMLResponse\"=\n" +
	"\n" +
	"SecretMeta\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x1d\n" +
	"\n" +
	"updated_at\x18\x02 \x01(\tR\tupdatedAt\"L\n" +
	"\x12ListSecretsRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\"E\n" +
	"\x13ListSecretsResponse\x12.\n" +
	"\asecrets\x18\x01 \x03(\v2\x14.projects.SecretMetaR\asecrets\"r\n" +
	"\x10SetSecretRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\x12\x10\n" +
	"\x03key\x18\x03 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x04 \x01(\tR\x05value\"\x13\n" +
	"\x11SetSecretResponse\"_\n" +
	"\x13DeleteSecretRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\x12\x10\n" +
	"\x03key\x18\x03 \x01(\tR\x03key\"\x16\n" +
	"\x14DeleteSecretResponse\"2\n" +
	"\bSecretKV\x12\x10\n" +
	"\x03key\x18\x01 \x01(\tR\x03key\x12\x14\n" +
	"\x05value\x18\x02 \x01(\tR\x05value\"\x0f\n" +
	"\rHealthRequest\"B\n" +
	"\x0eHealthResponse\x12\x18\n" +
	"\ahealthy\x18\x01 \x01(\bR\ahealthy\x12\x16\n" +
	"\x06status\x18\x02 \x01(\tR\x06status\"\xb3\x01\n" +
	"\rProjectMember\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x14\n" +
	"\x05email\x18\x02 \x01(\tR\x05email\x12!\n" +
	"\fdisplay_name\x18\x03 \x01(\tR\vdisplayName\x12\x12\n" +
	"\x04role\x18\x04 \x01(\tR\x04role\x12\x1d\n" +
	"\n" +
	"invited_by\x18\x05 \x01(\tR\tinvitedBy\x12\x1d\n" +
	"\n" +
	"created_at\x18\x06 \x01(\tR\tcreatedAt\"L\n" +
	"\x12ListMembersRequest\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\"\x93\x01\n" +
	"\x13ListMembersResponse\x121\n" +
	"\amembers\x18\x01 \x03(\v2\x17.projects.ProjectMemberR\amembers\x12%\n" +
	"\x0erequester_role\x18\x02 \x01(\tR\rrequesterRole\x12\"\n" +
	"\rowner_user_id\x18\x03 \x01(\tR\vownerUserId\"\xeb\x01\n" +
	"\x13InviteMemberRequest\x12\"\n" +
	"\rowner_user_id\x18\x01 \x01(\tR\vownerUserId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\x12&\n" +
	"\x0finvitee_user_id\x18\x03 \x01(\tR\rinviteeUserId\x12#\n" +
	"\rinvitee_email\x18\x04 \x01(\tR\finviteeEmail\x120\n" +
	"\x14invitee_display_name\x18\x05 \x01(\tR\x12inviteeDisplayName\x12\x12\n" +
	"\x04role\x18\x06 \x01(\tR\x04role\"G\n" +
	"\x14InviteMemberResponse\x12/\n" +
	"\x06member\x18\x01 \x01(\v2\x17.projects.ProjectMemberR\x06member\"\x96\x01\n" +
	"\x17UpdateMemberRoleRequest\x12\"\n" +
	"\rowner_user_id\x18\x01 \x01(\tR\vownerUserId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\x12$\n" +
	"\x0etarget_user_id\x18\x03 \x01(\tR\ftargetUserId\x12\x12\n" +
	"\x04role\x18\x04 \x01(\tR\x04role\"K\n" +
	"\x18UpdateMemberRoleResponse\x12/\n" +
	"\x06member\x18\x01 \x01(\v2\x17.projects.ProjectMemberR\x06member\"~\n" +
	"\x13RemoveMemberRequest\x12\"\n" +
	"\rowner_user_id\x18\x01 \x01(\tR\vownerUserId\x12\x1d\n" +
	"\n" +
	"project_id\x18\x02 \x01(\tR\tprojectId\x12$\n" +
	"\x0etarget_user_id\x18\x03 \x01(\tR\ftargetUserId\"\x16\n" +
	"\x14RemoveMemberResponse\":\n" +
	"\x19GetProjectInternalRequest\x12\x1d\n" +
	"\n" +
	"project_id\x18\x01 \x01(\tR\tprojectId\"\xbf\x02\n" +
	"\x1aGetProjectInternalResponse\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x19\n" +
	"\brepo_url\x18\x02 \x01(\tR\arepoUrl\x12%\n" +
	"\x0ewebhook_secret\x18\x03 \x01(\tR\rwebhookSecret\x12\x17\n" +
	"\assh_key\x18\x04 \x01(\tR\x06sshKey\x12+\n" +
	"\benv_vars\x18\x05 \x03(\v2\x10.projects.EnvVarR\aenvVars\x124\n" +
	"\x16pipeline_yaml_override\x18\x06 \x01(\tR\x14pipelineYamlOverride\x12,\n" +
	"\asecrets\x18\a \x03(\v2\x12.projects.SecretKVR\asecrets\x12%\n" +
	"\x0edefault_branch\x18\b \x01(\tR\rdefaultBranch\"\x1d\n" +
	"\x1bListProjectsInternalRequest\"O\n" +
	"\x16ProjectInternalSummary\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12%\n" +
	"\x0edefault_branch\x18\x02 \x01(\tR\rdefaultBranch\"\\\n" +
	"\x1cListProjectsInternalResponse\x12<\n" +
	"\bprojects\x18\x01 \x03(\v2 .projects.ProjectInternalSummaryR\bprojects2\xf0\f\n" +
	"\x0fProjectsService\x12P\n" +
	"\rCreateProject\x12\x1e.projects.CreateProjectRequest\x1a\x1f.projects.CreateProjectResponse\x12G\n" +
	"\n" +
	"GetProject\x12\x1b.projects.GetProjectRequest\x1a\x1c.projects.GetProjectResponse\x12M\n" +
	"\fListProjects\x12\x1d.projects.ListProjectsRequest\x1a\x1e.projects.ListProjectsResponse\x12Y\n" +
	"\x10VerifyConnection\x12!.projects.VerifyConnectionRequest\x1a\".projects.VerifyConnectionResponse\x12P\n" +
	"\rDeleteProject\x12\x1e.projects.DeleteProjectRequest\x1a\x1f.projects.DeleteProjectResponse\x12P\n" +
	"\rUpdateProject\x12\x1e.projects.UpdateProjectRequest\x1a\x1f.projects.UpdateProjectResponse\x12;\n" +
	"\x06Health\x12\x17.projects.HealthRequest\x1a\x18.projects.HealthResponse\x12G\n" +
	"\n" +
	"GetEnvVars\x12\x1b.projects.GetEnvVarsRequest\x1a\x1c.projects.GetEnvVarsResponse\x12P\n" +
	"\rUpdateEnvVars\x12\x1e.projects.UpdateEnvVarsRequest\x1a\x1f.projects.UpdateEnvVarsResponse\x12V\n" +
	"\x0fGetPipelineYAML\x12 .projects.GetPipelineYAMLRequest\x1a!.projects.GetPipelineYAMLResponse\x12V\n" +
	"\x0fSetPipelineYAML\x12 .projects.SetPipelineYAMLRequest\x1a!.projects.SetPipelineYAMLResponse\x12J\n" +
	"\vListSecrets\x12\x1c.projects.ListSecretsRequest\x1a\x1d.projects.ListSecretsResponse\x12D\n" +
	"\tSetSecret\x12\x1a.projects.SetSecretRequest\x1a\x1b.projects.SetSecretResponse\x12M\n" +
	"\fDeleteSecret\x12\x1d.projects.DeleteSecretRequest\x1a\x1e.projects.DeleteSecretResponse\x12_\n" +
	"\x12GetProjectInternal\x12#.projects.GetProjectInternalRequest\x1a$.projects.GetProjectInternalResponse\x12e\n" +
	"\x14ListProjectsInternal\x12%.projects.ListProjectsInternalRequest\x1a&.projects.ListProjectsInternalResponse\x12J\n" +
	"\vListMembers\x12\x1c.projects.ListMembersRequest\x1a\x1d.projects.ListMembersResponse\x12M\n" +
	"\fInviteMember\x12\x1d.projects.InviteMemberRequest\x1a\x1e.projects.InviteMemberResponse\x12Y\n" +
	"\x10UpdateMemberRole\x12!.projects.UpdateMemberRoleRequest\x1a\".projects.UpdateMemberRoleResponse\x12M\n" +
	"\fRemoveMember\x12\x1d.projects.RemoveMemberRequest\x1a\x1e.projects.RemoveMemberResponseB3Z1github.com/vsuaiqq/cicd/shared/proto/gen/projectsb\x06proto3"

var (
	file_projects_proto_rawDescOnce sync.Once
	file_projects_proto_rawDescData []byte
)

func file_projects_proto_rawDescGZIP() []byte {
	file_projects_proto_rawDescOnce.Do(func() {
		file_projects_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_projects_proto_rawDesc), len(file_projects_proto_rawDesc)))
	})
	return file_projects_proto_rawDescData
}

var file_projects_proto_msgTypes = make([]protoimpl.MessageInfo, 46)
var file_projects_proto_goTypes = []any{
	(*EnvVar)(nil),
	(*CreateProjectRequest)(nil),
	(*CreateProjectResponse)(nil),
	(*GetProjectRequest)(nil),
	(*ListProjectsRequest)(nil),
	(*ListProjectsResponse)(nil),
	(*ProjectSummary)(nil),
	(*GetProjectResponse)(nil),
	(*VerifyConnectionRequest)(nil),
	(*VerifyConnectionResponse)(nil),
	(*DeleteProjectRequest)(nil),
	(*DeleteProjectResponse)(nil),
	(*UpdateProjectRequest)(nil),
	(*UpdateProjectResponse)(nil),
	(*GetEnvVarsRequest)(nil),
	(*GetEnvVarsResponse)(nil),
	(*UpdateEnvVarsRequest)(nil),
	(*UpdateEnvVarsResponse)(nil),
	(*GetPipelineYAMLRequest)(nil),
	(*GetPipelineYAMLResponse)(nil),
	(*SetPipelineYAMLRequest)(nil),
	(*SetPipelineYAMLResponse)(nil),
	(*SecretMeta)(nil),
	(*ListSecretsRequest)(nil),
	(*ListSecretsResponse)(nil),
	(*SetSecretRequest)(nil),
	(*SetSecretResponse)(nil),
	(*DeleteSecretRequest)(nil),
	(*DeleteSecretResponse)(nil),
	(*SecretKV)(nil),
	(*HealthRequest)(nil),
	(*HealthResponse)(nil),
	(*ProjectMember)(nil),
	(*ListMembersRequest)(nil),
	(*ListMembersResponse)(nil),
	(*InviteMemberRequest)(nil),
	(*InviteMemberResponse)(nil),
	(*UpdateMemberRoleRequest)(nil),
	(*UpdateMemberRoleResponse)(nil),
	(*RemoveMemberRequest)(nil),
	(*RemoveMemberResponse)(nil),
	(*GetProjectInternalRequest)(nil),
	(*GetProjectInternalResponse)(nil),
	(*ListProjectsInternalRequest)(nil),
	(*ProjectInternalSummary)(nil),
	(*ListProjectsInternalResponse)(nil),
}
var file_projects_proto_depIdxs = []int32{
	6,
	0,
	0,
	22,
	32,
	32,
	32,
	0,
	29,
	44,
	1,
	3,
	4,
	8,
	10,
	12,
	30,
	14,
	16,
	18,
	20,
	23,
	25,
	27,
	41,
	43,
	33,
	35,
	37,
	39,
	2,
	7,
	5,
	9,
	11,
	13,
	31,
	15,
	17,
	19,
	21,
	24,
	26,
	28,
	42,
	45,
	34,
	36,
	38,
	40,
	30,
	10,
	10,
	10,
	0,
}

func init() { file_projects_proto_init() }
func file_projects_proto_init() {
	if File_projects_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_projects_proto_rawDesc), len(file_projects_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   46,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_projects_proto_goTypes,
		DependencyIndexes: file_projects_proto_depIdxs,
		MessageInfos:      file_projects_proto_msgTypes,
	}.Build()
	File_projects_proto = out.File
	file_projects_proto_goTypes = nil
	file_projects_proto_depIdxs = nil
}
