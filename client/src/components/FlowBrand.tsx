
import { css } from '@emotion/react'

const logoMark = css({
  width: 30,
  height: 30,
  borderRadius: 8,
  background: 'linear-gradient(135deg, var(--accent) 0%, var(--indigo) 100%)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  flexShrink: 0,
  boxShadow: '0 2px 8px var(--accent-glow)',
  position: 'relative',
  overflow: 'hidden',
  '&::after': {
    content: '""',
    position: 'absolute',
    inset: 0,
    background: 'linear-gradient(135deg, var(--glass-highlight) 0%, transparent 60%)',
    borderRadius: 'inherit',
  },
})

const brandMarkCell = css({
  width: 40,
  height: 40,
  display: 'grid',
  placeItems: 'center',
  flexShrink: 0,
})

const brandRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  minWidth: 0,
})

const brandText = css({
  display: 'flex',
  alignItems: 'center',
  gap: 4,
  minWidth: 0,
  overflow: 'hidden',
  opacity: 1,
  transition: 'opacity var(--duration-sidebar) var(--ease-sidebar)',
  transitionDelay: '50ms',
  '.app-sidebar[data-collapsed] &': {
    opacity: 0,
    width: 0,
    overflow: 'hidden',
    position: 'absolute',
    pointerEvents: 'none',
    transitionDelay: '0ms',
  },
})

const logoName = css({
  fontWeight: 700,
  fontSize: '1.0625rem',
  letterSpacing: '-0.03em',
  color: 'var(--text-primary)',
  whiteSpace: 'nowrap',
})

const logoBadge = css({
  fontSize: '0.625rem',
  fontWeight: 700,
  letterSpacing: '0.08em',
  textTransform: 'uppercase',
  color: 'var(--accent)',
  background: 'var(--accent-muted)',
  border: '1px solid var(--border-accent)',
  borderRadius: 4,
  padding: '1px 6px',
  whiteSpace: 'nowrap',
  flexShrink: 0,
})

export function FlowIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden>
      <path d="M3 8h4M9 4l4 4-4 4" stroke="var(--text-on-accent)" strokeWidth="1.7" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export function FlowBrandLogo(_props: { collapsed?: boolean }) {
  return (
    <div css={brandRow}>
      <div css={brandMarkCell}>
        <div css={logoMark}>
          <FlowIcon />
        </div>
      </div>
      <div css={brandText}>
        <span css={logoName}>Flow</span>
        <span css={logoBadge}>CI/CD</span>
      </div>
    </div>
  )
}
