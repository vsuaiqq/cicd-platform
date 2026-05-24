import type { Translations } from '../../i18n'

export function formatTriggerType(trigger: string, t: Translations['project']): string {
  const key = trigger.trim().toLowerCase()
  const map: Record<string, string | undefined> = {
    push: t.triggerPush,
    manual: t.triggerManual,
    webhook: t.triggerWebhook,
    schedule: t.triggerSchedule,
  }
  return map[key] ?? trigger
}
export function sortRunsByNewest<T extends { created_at: string }>(runs: T[]): T[] {
  return [...runs].sort(
    (a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
  )
}

