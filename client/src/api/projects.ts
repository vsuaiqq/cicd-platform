import { apiRequest } from './client'

const PREFIX = '/api/v1/projects'

export interface CreateProjectRequest {
  name: string
  repository_url: string
  default_branch?: string
}

export interface UpdateProjectRequest {
  name: string
  default_branch: string
}


export interface Project {
  id: string
  user_id: string
  name: string
  repo_url: string
  default_branch: string
  public_key: string
  webhook_secret: string
  webhook_url: string
  status: string
}


export interface ProjectSummary {
  id: string
  name: string
  repo_url: string
  default_branch: string
  status: string
}

export interface VerifyProjectResponse {
  success: boolean
  status?: string
  error?: string
}

export interface EnvVar {
  key: string
  value: string
}

export function listProjects(token: string): Promise<ProjectSummary[]> {
  return apiRequest<ProjectSummary[]>(PREFIX, { token })
}

export function createProject(
  body: CreateProjectRequest,
  token: string
): Promise<Project> {
  return apiRequest<Project>(PREFIX, {
    method: 'POST',
    body: JSON.stringify(body),
    token,
  })
}

export function getProject(id: string, token: string): Promise<Project> {
  return apiRequest<Project>(`${PREFIX}/${encodeURIComponent(id)}`, { token })
}

export function verifyProject(id: string, token: string): Promise<VerifyProjectResponse> {
  return apiRequest<VerifyProjectResponse>(`${PREFIX}/${encodeURIComponent(id)}/verify`, {
    method: 'POST',
    token,
  })
}

export function deleteProject(id: string, token: string): Promise<void> {
  return apiRequest<void>(`${PREFIX}/${encodeURIComponent(id)}`, {
    method: 'DELETE',
    token,
  })
}

export function updateProject(
  id: string,
  body: UpdateProjectRequest,
  token: string,
): Promise<Project> {
  return apiRequest<Project>(`${PREFIX}/${encodeURIComponent(id)}`, {
    method: 'PATCH',
    body: JSON.stringify(body),
    token,
  })
}

export function getEnvVars(id: string, token: string): Promise<EnvVar[]> {
  return apiRequest<{ vars: EnvVar[] }>(`${PREFIX}/${encodeURIComponent(id)}/env`, { token })
    .then((r) => r.vars ?? [])
}

export function updateEnvVars(id: string, vars: EnvVar[], token: string): Promise<void> {
  return apiRequest<void>(`${PREFIX}/${encodeURIComponent(id)}/env`, {
    method: 'PUT',
    body: JSON.stringify({ vars }),
    token,
  })
}




export interface SecretMeta {
  key: string
  updated_at: string
}

export function listSecrets(id: string, token: string): Promise<SecretMeta[]> {
  return apiRequest<{ secrets: SecretMeta[] }>(`${PREFIX}/${encodeURIComponent(id)}/secrets`, { token })
    .then((r) => r.secrets ?? [])
}

export function setSecret(id: string, key: string, value: string, token: string): Promise<void> {
  return apiRequest<void>(`${PREFIX}/${encodeURIComponent(id)}/secrets/${encodeURIComponent(key)}`, {
    method: 'PUT',
    body: JSON.stringify({ value }),
    token,
  })
}

export function deleteSecret(id: string, key: string, token: string): Promise<void> {
  return apiRequest<void>(`${PREFIX}/${encodeURIComponent(id)}/secrets/${encodeURIComponent(key)}`, {
    method: 'DELETE',
    token,
  })
}

export function getPipelineYAML(id: string, token: string): Promise<string> {
  return apiRequest<{ yaml: string }>(`${PREFIX}/${encodeURIComponent(id)}/pipeline-yaml`, { token })
    .then((r) => r.yaml ?? '')
}

export function setPipelineYAML(id: string, yaml: string, token: string): Promise<void> {
  return apiRequest<void>(`${PREFIX}/${encodeURIComponent(id)}/pipeline-yaml`, {
    method: 'PUT',
    body: JSON.stringify({ yaml }),
    token,
  })
}



export interface ProjectMemberDTO {
  user_id:      string
  email:        string
  display_name: string
  role:         'owner' | 'editor' | 'viewer'
  invited_by:   string
  created_at:   string
}

export interface ListMembersResponse {
  members:        ProjectMemberDTO[]
  requester_role: 'owner' | 'editor' | 'viewer'
  owner_user_id:  string
}

export function listMembers(id: string, token: string): Promise<ListMembersResponse> {
  return apiRequest<ListMembersResponse>(`${PREFIX}/${encodeURIComponent(id)}/members`, { token })
    .then((r) => ({ members: r.members ?? [], requester_role: r.requester_role, owner_user_id: r.owner_user_id ?? '' }))
}

export function inviteMember(id: string, email: string, role: 'editor' | 'viewer', token: string): Promise<ProjectMemberDTO> {
  return apiRequest<ProjectMemberDTO>(`${PREFIX}/${encodeURIComponent(id)}/members`, {
    method: 'POST',
    body: JSON.stringify({ email, role }),
    token,
  })
}

export function updateMemberRole(id: string, userId: string, role: 'editor' | 'viewer', token: string): Promise<ProjectMemberDTO> {
  return apiRequest<ProjectMemberDTO>(`${PREFIX}/${encodeURIComponent(id)}/members/${encodeURIComponent(userId)}`, {
    method: 'PATCH',
    body: JSON.stringify({ role }),
    token,
  })
}

export function removeMember(id: string, userId: string, token: string): Promise<void> {
  return apiRequest<void>(`${PREFIX}/${encodeURIComponent(id)}/members/${encodeURIComponent(userId)}`, {
    method: 'DELETE',
    token,
  })
}
