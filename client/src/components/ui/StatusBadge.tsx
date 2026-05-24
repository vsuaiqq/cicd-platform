import { useI18n } from '../../i18n'
import { formatStatus } from '../../lib/formatStatus'
import { statusBadge, statusActive, statusPending, statusError } from '../../styles/theme'
import { css } from '@emotion/react'

type Status = 'active' | 'pending' | 'failed' | 'error' | string

const dot = css({
  width: 6,
  height: 6,
  borderRadius: '50%',
  flexShrink: 0,
})

const dotActive = css({
  background: 'var(--success)',
  animation: 'pulse-dot 2s ease-in-out infinite',
})

const dotPending = css({
  background: 'var(--warning)',
})

const dotError = css({
  background: 'var(--danger)',
})

export default function StatusBadge({ status }: { status: Status }) {
  const { t } = useI18n()
  const s = status?.toLowerCase() ?? ''
  const isActive = s === 'active'
  const isPending = s === 'pending'
  const isError = s === 'failed' || s === 'error'

  const badgeStyle = isActive ? statusActive : isPending ? statusPending : isError ? statusError : statusPending
  const dotStyle = isActive ? dotActive : isPending ? dotPending : isError ? dotError : dotPending
  const label = formatStatus(status, t.status)

  return (
    <span
      css={[statusBadge, badgeStyle]}
      role="status"
      aria-label={`${t.settings.statusLabel}: ${label}`}
    >
      <span css={[dot, dotStyle]} />
      {label}
    </span>
  )
}
