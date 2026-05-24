import { css } from '@emotion/react'
import { statusBg, statusColor } from '../../api/pipelines'
import { formatStatus } from '../../lib/formatStatus'
import { useI18n } from '../../i18n'

export default function RunStatusBadge({ status }: { status: string }) {
  const { t } = useI18n()
  const label = formatStatus(status, t.status)
  const dot = css({
    width: 6,
    height: 6,
    borderRadius: '50%',
    background: statusColor(status),
    flexShrink: 0,
    animation: status === 'running' ? 'pulse-dot 1.5s ease-in-out infinite' : 'none',
  })
  const badge = css({
    display: 'inline-flex',
    alignItems: 'center',
    gap: 6,
    padding: '4px 10px',
    borderRadius: 'var(--radius-full)',
    fontSize: '0.75rem',
    fontWeight: 600,
    letterSpacing: '0.01em',
    background: statusBg(status),
    color: statusColor(status),
    border: '1px solid var(--border-subtle)',
    whiteSpace: 'nowrap',
  })
  return (
    <span css={badge} role="status">
      <span css={dot} aria-hidden />
      {label}
    </span>
  )
}
