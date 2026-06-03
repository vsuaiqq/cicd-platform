import { createApi } from '@reduxjs/toolkit/query/react'
import type { DashboardData, Period, PerformanceGateResult } from '../../api/analytics'
import type { PipelineRun } from '../../api/pipelines'
import type {
  CreateProjectRequest,
  EnvVar,
  ListMembersResponse,
  Project,
  ProjectMemberDTO,
  ProjectSummary,
  SecretMeta,
  UpdateProjectRequest,
  VerifyProjectResponse,
} from '../../api/projects'
import { apiBaseQuery } from './baseQuery'

const projects = '/api/v1/projects'
const project = (id: string) => `${projects}/${encodeURIComponent(id)}`
const pipeline = '/api/v1/pipeline'

export const apiSlice = createApi({
  reducerPath: 'api',
  baseQuery: apiBaseQuery,
  tagTypes: ['ProjectList', 'Project', 'Runs', 'Run', 'EnvVars', 'Secrets', 'Members', 'PipelineYaml', 'Dashboard', 'PerformanceGate'],

  keepUnusedDataFor: 300,

  refetchOnMountOrArgChange: 60,
  refetchOnFocus: true,
  refetchOnReconnect: true,
  endpoints: (build) => ({

    listProjects: build.query<ProjectSummary[], void>({
      query: () => projects,
      providesTags: [{ type: 'ProjectList', id: 'LIST' }],
    }),

    getProject: build.query<Project, string>({
      query: (id) => project(id),
      providesTags: (_r, _e, id) => [{ type: 'Project', id }],
    }),

    createProject: build.mutation<Project, CreateProjectRequest>({
      query: (body) => ({ url: projects, method: 'POST', body }),
      invalidatesTags: [{ type: 'ProjectList', id: 'LIST' }],
    }),

    updateProject: build.mutation<Project, { id: string; body: UpdateProjectRequest }>({
      query: ({ id, body }) => ({ url: project(id), method: 'PATCH', body }),
      invalidatesTags: (_r, _e, { id }) => [
        { type: 'Project', id },
        { type: 'ProjectList', id: 'LIST' },
      ],
    }),

    deleteProject: build.mutation<void, string>({
      query: (id) => ({ url: project(id), method: 'DELETE' }),
      invalidatesTags: (_r, _e, id) => [
        { type: 'Project', id },
        { type: 'ProjectList', id: 'LIST' },
        { type: 'Runs', id },
      ],
    }),

    verifyProject: build.mutation<VerifyProjectResponse, string>({
      query: (id) => ({ url: `${project(id)}/verify`, method: 'POST' }),
      invalidatesTags: (_r, _e, id) => [{ type: 'Project', id }],
    }),


    getEnvVars: build.query<EnvVar[], string>({
      query: (id) => `${project(id)}/env`,
      transformResponse: (r: { vars?: EnvVar[] }) => r.vars ?? [],
      providesTags: (_r, _e, id) => [{ type: 'EnvVars', id }],
    }),

    updateEnvVars: build.mutation<void, { id: string; vars: EnvVar[] }>({
      query: ({ id, vars }) => ({
        url: `${project(id)}/env`,
        method: 'PUT',
        body: { vars },
      }),
      invalidatesTags: (_r, _e, { id }) => [{ type: 'EnvVars', id }],
    }),

    listSecrets: build.query<SecretMeta[], string>({
      query: (id) => `${project(id)}/secrets`,
      transformResponse: (r: { secrets?: SecretMeta[] }) => r.secrets ?? [],
      providesTags: (_r, _e, id) => [{ type: 'Secrets', id }],
    }),

    setSecret: build.mutation<void, { projectId: string; key: string; value: string }>({
      query: ({ projectId, key, value }) => ({
        url: `${project(projectId)}/secrets/${encodeURIComponent(key)}`,
        method: 'PUT',
        body: { value },
      }),
      invalidatesTags: (_r, _e, { projectId }) => [{ type: 'Secrets', id: projectId }],
    }),

    deleteSecret: build.mutation<void, { projectId: string; key: string }>({
      query: ({ projectId, key }) => ({
        url: `${project(projectId)}/secrets/${encodeURIComponent(key)}`,
        method: 'DELETE',
      }),
      invalidatesTags: (_r, _e, { projectId }) => [{ type: 'Secrets', id: projectId }],
    }),

    getPipelineYaml: build.query<string, string>({
      query: (id) => `${project(id)}/pipeline-yaml`,
      transformResponse: (r: { yaml?: string }) => r.yaml ?? '',
      providesTags: (_r, _e, id) => [{ type: 'PipelineYaml', id }],
    }),

    setPipelineYaml: build.mutation<void, { id: string; yaml: string }>({
      query: ({ id, yaml }) => ({
        url: `${project(id)}/pipeline-yaml`,
        method: 'PUT',
        body: { yaml },
      }),
      invalidatesTags: (_r, _e, { id }) => [{ type: 'PipelineYaml', id }],
    }),


    listMembers: build.query<ListMembersResponse, string>({
      query: (id) => `${project(id)}/members`,
      transformResponse: (r: ListMembersResponse) => ({
        members: r.members ?? [],
        requester_role: r.requester_role,
        owner_user_id: r.owner_user_id ?? '',
      }),
      providesTags: (_r, _e, id) => [{ type: 'Members', id }],
    }),

    inviteMember: build.mutation<
      ProjectMemberDTO,
      { projectId: string; email: string; role: 'editor' | 'viewer' }
    >({
      query: ({ projectId, email, role }) => ({
        url: `${project(projectId)}/members`,
        method: 'POST',
        body: { email, role },
      }),
      invalidatesTags: (_r, _e, { projectId }) => [{ type: 'Members', id: projectId }],
    }),

    updateMemberRole: build.mutation<
      ProjectMemberDTO,
      { projectId: string; userId: string; role: 'editor' | 'viewer' }
    >({
      query: ({ projectId, userId, role }) => ({
        url: `${project(projectId)}/members/${encodeURIComponent(userId)}`,
        method: 'PATCH',
        body: { role },
      }),
      invalidatesTags: (_r, _e, { projectId }) => [{ type: 'Members', id: projectId }],
    }),

    removeMember: build.mutation<void, { projectId: string; userId: string }>({
      query: ({ projectId, userId }) => ({
        url: `${project(projectId)}/members/${encodeURIComponent(userId)}`,
        method: 'DELETE',
      }),
      invalidatesTags: (_r, _e, { projectId }) => [{ type: 'Members', id: projectId }],
    }),


    listRuns: build.query<PipelineRun[], string>({
      query: (projectId) => {
        const qs = new URLSearchParams({ project_id: projectId })
        return `${pipeline}/runs?${qs}`
      },
      providesTags: (result, _e, projectId) =>
        result
          ? [
              ...result.map((r) => ({ type: 'Run' as const, id: r.id })),
              { type: 'Runs', id: projectId },
            ]
          : [{ type: 'Runs', id: projectId }],
    }),

    getRun: build.query<PipelineRun, string>({
      query: (runId) => `${pipeline}/runs/${encodeURIComponent(runId)}`,
      providesTags: (_r, _e, runId) => [{ type: 'Run', id: runId }],
    }),

    cancelRun: build.mutation<void, string>({
      query: (runId) => ({
        url: `${pipeline}/runs/${encodeURIComponent(runId)}/cancel`,
        method: 'POST',
      }),
      invalidatesTags: (_r, _e, runId) => [{ type: 'Run', id: runId }],
    }),

    approveJob: build.mutation<void, { runId: string; jobId: string }>({
      query: ({ runId, jobId }) => ({
        url: `${pipeline}/runs/${encodeURIComponent(runId)}/jobs/${encodeURIComponent(jobId)}/approve`,
        method: 'POST',
      }),
      invalidatesTags: (_r, _e, { runId }) => [{ type: 'Run', id: runId }],
    }),

    rejectJob: build.mutation<void, { runId: string; jobId: string }>({
      query: ({ runId, jobId }) => ({
        url: `${pipeline}/runs/${encodeURIComponent(runId)}/jobs/${encodeURIComponent(jobId)}/reject`,
        method: 'POST',
      }),
      invalidatesTags: (_r, _e, { runId }) => [{ type: 'Run', id: runId }],
    }),


    getDashboard: build.query<DashboardData, { projectId: string; period: Period }>({
      query: ({ projectId, period }) => {
        const params = new URLSearchParams({ project_id: projectId, period })
        return `/api/v1/analytics/dashboard?${params}`
      },
      transformResponse: (d: DashboardData) => ({
        ...d,
        trend: d.trend ?? [],
        top_failing_jobs: d.top_failing_jobs ?? [],
        top_slow_jobs: d.top_slow_jobs ?? [],
        flaky_jobs: d.flaky_jobs ?? [],
      }),
      providesTags: (_r, _e, { projectId, period }) => [
        { type: 'Dashboard', id: `${projectId}-${period}` },
      ],
    }),

    getPerformanceGateResult: build.query<
      PerformanceGateResult,
      { runId: string; jobName: string }
    >({
      query: ({ runId, jobName }) => {
        const params = new URLSearchParams({ run_id: runId, job_name: jobName })
        return `/api/v1/analytics/performance-gate?${params}`
      },
      transformResponse: (r: PerformanceGateResult) => ({
        ...r,
        metrics: r.metrics ?? [],
      }),
      providesTags: (_r, _e, { runId, jobName }) => [
        { type: 'PerformanceGate', id: `${runId}-${jobName}` },
      ],
    }),
  }),
})

export const {
  useListProjectsQuery,
  useGetProjectQuery,
  useCreateProjectMutation,
  useUpdateProjectMutation,
  useDeleteProjectMutation,
  useVerifyProjectMutation,
  useGetEnvVarsQuery,
  useUpdateEnvVarsMutation,
  useListSecretsQuery,
  useSetSecretMutation,
  useDeleteSecretMutation,
  useGetPipelineYamlQuery,
  useSetPipelineYamlMutation,
  useListMembersQuery,
  useInviteMemberMutation,
  useUpdateMemberRoleMutation,
  useRemoveMemberMutation,
  useListRunsQuery,
  useGetRunQuery,
  useCancelRunMutation,
  useApproveJobMutation,
  useRejectJobMutation,
  useGetDashboardQuery,
  useGetPerformanceGateResultQuery,
} = apiSlice
