import type { ReactNode } from 'react'
import { css } from '@emotion/react'
import { emptyState, emptyStateIcon, emptyTitle, emptyDesc } from '../../styles/theme'

const iconGlow = css({
  '&::after': {
    content: '""',
    position: 'absolute',
    inset: -1,
    borderRadius: 'inherit',
    background: 'radial-gradient(circle at 50% 0%, var(--accent-muted), transparent 70%)',
    pointerEvents: 'none',
  },
  position: 'relative',
})

function DefaultIcon() {
  return (
    <svg width="26" height="26" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" aria-hidden>
      <rect x="3" y="3" width="7" height="7" rx="1" />
      <rect x="14" y="3" width="7" height="7" rx="1" />
      <rect x="3" y="14" width="7" height="7" rx="1" />
      <path d="M17.5 14v7M14 17.5h7" />
    </svg>
  )
}

export default function EmptyState({
  title,
  description,
  action,
  icon,
}: {
  title: string
  description: string
  action: ReactNode
  icon?: ReactNode
}) {
  return (
    <div css={emptyState}>
      <div css={[emptyStateIcon, iconGlow]}>
        {icon ?? <DefaultIcon />}
      </div>
      <div css={emptyTitle}>{title}</div>
      <p css={emptyDesc}>{description}</p>
      {action}
    </div>
  )
}
