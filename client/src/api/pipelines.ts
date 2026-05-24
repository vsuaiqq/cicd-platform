import { apiRequest, apiBase } from './client'

const PREFIX = '/api/v1/pipeline'

export type RunStatus  = 'pending' | 'running' | 'success' | 'failed' | 'cancelled'
export type JobStatus  = 'pending' | 'running' | 'success' | 'failed' | 'skipped' | 'cancelled' | 'awaiting_approval'
export type StepStatus = 'pending' | 'running' | 'success' | 'failed' | 'cancelled'

export interface PipelineStep {
  index:       number
  name:        string
  status:      StepStatus
  log_output:  string
  exit_code:   number
  started_at:  string | null
  finished_at: string | null
}

export interface PipelineJob {
  id:           string
  name:         string
  display_name: string
  status:       JobStatus
  attempt:      number
  started_at:   string | null
  finished_at:  string | null
  steps:        PipelineStep[]
}

export interface PipelineRun {
  id:            string
  project_id:    string
  commit_sha:    string
  branch:        string
  trigger_type:  string
  status:        RunStatus
  created_at:    string
  started_at:    string | null
  finished_at:   string | null
  pipeline_yaml?: string
  jobs?:         PipelineJob[]
}

export function listRuns(projectId: string, token: string): Promise<PipelineRun[]> {
  const qs = new URLSearchParams({ project_id: projectId })
  return apiRequest<PipelineRun[]>(`${PREFIX}/runs?${qs}`, { token })
}

export function getRun(runId: string, token: string): Promise<PipelineRun> {
  return apiRequest<PipelineRun>(`${PREFIX}/runs/${encodeURIComponent(runId)}`, { token })
}

export function cancelRun(runId: string, token: string): Promise<void> {
  return apiRequest<void>(`${PREFIX}/runs/${encodeURIComponent(runId)}/cancel`, {
    method: 'POST',
    token,
  })
}

export function approveJob(runId: string, jobId: string, token: string): Promise<void> {
  return apiRequest<void>(
    `${PREFIX}/runs/${encodeURIComponent(runId)}/jobs/${encodeURIComponent(jobId)}/approve`,
    { method: 'POST', token },
  )
}

export function rejectJob(runId: string, jobId: string, token: string): Promise<void> {
  return apiRequest<void>(
    `${PREFIX}/runs/${encodeURIComponent(runId)}/jobs/${encodeURIComponent(jobId)}/reject`,
    { method: 'POST', token },
  )
}



export type WSEventType =
  | 'run_updated'
  | 'job_updated'
  | 'job_awaiting_approval'



  | 'heartbeat'




export interface StepEvent {
  index:       number
  name:        string
  status:      string
  exit_code:   number
  log_output:  string
  started_at:  number
  finished_at: number
}

export interface WSEvent {
  type:       WSEventType
  run_id:     string

  project_id?: string
  status?:    string
  job_name?:  string
  job_id?:    string
  steps?:     StepEvent[]



  server_time_ms?:    number

  run_finished_at_ms?: number

  job_started_at_ms?:  number

  job_finished_at_ms?: number
}

export interface ConnectRunWSOptions {
  token:    string
  onEvent:  (event: WSEvent) => void

  onOpen?:  () => void
  onClose?: () => void
  onError?: () => void
}


export function connectRunWS(runId: string, options: ConnectRunWSOptions): WebSocket {
  let wsBase: string
  if (apiBase) {
    wsBase = apiBase.replace(/^http(s?):\/\//, 'ws$1://')
  } else {
    const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    wsBase = `${proto}//${window.location.host}`
  }

  const url = `${wsBase}/api/v1/pipeline/ws/runs/${encodeURIComponent(runId)}?token=${encodeURIComponent(options.token)}`
  const ws  = new WebSocket(url)

  ws.onopen    = () => options.onOpen?.()
  ws.onmessage = (e) => {
    try {
      options.onEvent(JSON.parse(e.data) as WSEvent)
    } catch {

    }
  }
  ws.onclose = () => options.onClose?.()
  ws.onerror = () => options.onError?.()

  return ws
}

export interface ConnectGlobalWSOptions {
  token:    string
  onEvent:  (event: WSEvent) => void
  onOpen?:  () => void
  onClose?: () => void
  onError?: () => void
}


export function connectGlobalWS(options: ConnectGlobalWSOptions): WebSocket {
  let wsBase: string
  if (apiBase) {
    wsBase = apiBase.replace(/^http(s?):\/\//, 'ws$1://')
  } else {
    const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
    wsBase = `${proto}//${window.location.host}`
  }

  const url = `${wsBase}/api/v1/pipeline/ws/events?token=${encodeURIComponent(options.token)}`
  const ws  = new WebSocket(url)

  ws.onopen    = () => options.onOpen?.()
  ws.onmessage = (e) => {
    try {
      options.onEvent(JSON.parse(e.data) as WSEvent)
    } catch {

    }
  }
  ws.onclose = () => options.onClose?.()
  ws.onerror = () => options.onError?.()

  return ws
}



export type RunStatusEvent = WSEvent



export function statusColor(status: string): string {
  switch (status) {
    case 'success':           return 'var(--success)'
    case 'failed':            return 'var(--danger)'
    case 'running':           return 'var(--accent)'
    case 'skipped':
    case 'cancelled':         return 'var(--text-disabled)'
    case 'awaiting_approval': return 'var(--warning)'
    default:                  return 'var(--text-tertiary)'
  }
}

export function statusBg(status: string): string {
  switch (status) {
    case 'success':           return 'var(--success-muted)'
    case 'failed':            return 'var(--danger-muted)'
    case 'running':           return 'var(--accent-muted)'
    case 'awaiting_approval': return 'var(--warning-muted)'
    default:                  return 'var(--bg-overlay)'
  }
}




export function durationMs(
  run: Pick<PipelineRun, 'started_at' | 'finished_at'>,
  nowMs?: number,
): number | null {
  if (!run.started_at) return null
  const end = run.finished_at ? new Date(run.finished_at) : new Date(nowMs ?? Date.now())
  return end.getTime() - new Date(run.started_at).getTime()
}

export function formatDuration(ms: number): string {
  const s = Math.floor(ms / 1000)
  if (s < 60) return `${s}s`
  const m   = Math.floor(s / 60)
  const rem = s % 60
  return `${m}m ${rem}s`
}

export function timeAgo(iso: string, lang: 'en' | 'ru' = 'en'): string {
  const diff = Date.now() - new Date(iso).getTime()
  const s    = Math.floor(diff / 1000)
  if (lang === 'ru') {
    if (s < 60)  return 'только что'
    const m = Math.floor(s / 60)
    if (m < 60)  return `${m} мин. назад`
    const h = Math.floor(m / 60)
    if (h < 24)  return `${h} ч. назад`
    return `${Math.floor(h / 24)} д. назад`
  }
  if (s < 60)  return `${s}s ago`
  const m = Math.floor(s / 60)
  if (m < 60)  return `${m}m ago`
  const h = Math.floor(m / 60)
  if (h < 24)  return `${h}h ago`
  return `${Math.floor(h / 24)}d ago`
}
