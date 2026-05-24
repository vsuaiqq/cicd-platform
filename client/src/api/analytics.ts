import { apiRequest } from './client'

const PREFIX = '/api/v1/analytics'

export interface DailyPoint {
  date: string
  total: number
  success: number
  failed: number
}

export interface JobStat {
  job_name: string
  total_runs: number
  failure_rate: number
  avg_duration_sec: number
  avg_attempts: number
}

export interface DashboardData {
  total_runs: number
  success_count: number
  failed_count: number
  cancelled_count: number
  success_rate: number
  avg_duration_sec: number
  p50_duration_sec: number
  p95_duration_sec: number
  trend: DailyPoint[]
  top_failing_jobs: JobStat[]
  top_slow_jobs: JobStat[]
  flaky_jobs: JobStat[]
}

export type Period = '7d' | '30d' | '90d'

export function getDashboard(projectId: string, period: Period, token: string | null): Promise<DashboardData> {
  const params = new URLSearchParams({ project_id: projectId, period })
  return apiRequest<DashboardData>(`${PREFIX}/dashboard?${params}`, { token })
}
