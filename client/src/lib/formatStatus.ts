import type { Translations } from '../i18n'

type StatusT = Translations['status']

const STATUS_MAP: Record<string, keyof StatusT> = {
  running: 'running',
  success: 'success',
  failed: 'failed',
  cancelled: 'cancelled',
  pending: 'pending',
  skipped: 'skipped',
  awaiting_approval: 'awaitingApproval',
  live: 'live',
  active: 'active',
  error: 'error',
}

export function formatStatus(status: string, statusT: StatusT): string {
  const key = STATUS_MAP[status.toLowerCase()]
  return key ? statusT[key] : status
}
