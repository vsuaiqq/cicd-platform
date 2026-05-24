import { css } from '@emotion/react'

export const backdrop = css({
  position: 'fixed',
  inset: 0,
  zIndex: 8000,
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'center',
  padding: 'min(10vh, 88px) 20px 20px',
  background: 'var(--overlay-backdrop)',
  backdropFilter: 'blur(10px)',
  WebkitBackdropFilter: 'blur(10px)',
  animation: 'modal-backdrop 0.18s var(--ease-out) both',
})

export const dialog = css({
  position: 'relative',
  width: '100%',
  maxWidth: 600,
  background: 'var(--bg-elevated)',
  border: '1px solid var(--border-default)',
  borderRadius: 14,
  boxShadow: 'var(--shadow-overlay), 0 0 0 1px var(--overlay-inset) inset, 0 20px 64px rgba(0,0,0,0.4)',
  overflow: 'hidden',
  display: 'flex',
  flexDirection: 'column',
  maxHeight: 'min(480px, 68vh)',
  animation: 'modal-slide 0.2s var(--ease-out) both',
})

export const searchRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  padding: '0 12px 0 14px',
  borderBottom: '1px solid var(--border-subtle)',
  flexShrink: 0,
})

export const searchIcon = css({
  display: 'flex',
  color: 'var(--text-tertiary)',
  flexShrink: 0,
})

export const searchInput = css({
  flex: 1,
  minWidth: 0,
  padding: '15px 8px',
  border: 'none',
  outline: 'none',
  background: 'transparent',
  color: 'var(--text-primary)',
  fontSize: '1rem',
  fontWeight: 400,
  lineHeight: 1.35,
  caretColor: 'var(--text-primary)',
  boxShadow: 'none',
  WebkitAppearance: 'none',
  '&::placeholder': { color: 'var(--text-disabled)' },
  '&:focus': {
    outline: 'none',
    boxShadow: 'none',
  },
  '&:focus-visible': {
    outline: 'none',
    boxShadow: 'none',
  },
})

export const kbd = css({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  minWidth: 22,
  height: 22,
  padding: '0 6px',
  borderRadius: 6,
  border: '1px solid var(--border-default)',
  background: 'var(--bg-base)',
  fontSize: '0.6875rem',
  fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
  fontWeight: 500,
  color: 'var(--text-disabled)',
  flexShrink: 0,
  letterSpacing: '0.02em',
})

export const spinner = css({
  width: 16,
  height: 16,
  borderRadius: '50%',
  border: '2px solid var(--border-subtle)',
  borderTopColor: 'var(--text-tertiary)',
  flexShrink: 0,
  animation: 'spin 0.65s linear infinite',
})

export const list = css({
  flex: 1,
  minHeight: 100,
  overflowY: 'auto',
  padding: '6px 0 8px',
  '&::-webkit-scrollbar': { width: 6 },
  '&::-webkit-scrollbar-thumb': { background: 'var(--border-default)', borderRadius: 6 },
})

export const groupHead = css({
  padding: '10px 18px 4px',
  fontSize: '0.6875rem',
  fontWeight: 600,
  color: 'var(--text-disabled)',
  letterSpacing: '0.07em',
  textTransform: 'uppercase',
})

export const row = css({
  position: 'relative',
  display: 'flex',
  alignItems: 'center',
  gap: 12,
  margin: '1px 6px',
  padding: '8px 12px 8px 14px',
  borderRadius: 'var(--radius-md)',
  cursor: 'pointer',
  userSelect: 'none',
  transition: 'background 0.1s var(--ease-out)',
  '&[data-active="true"]': {
    background: 'var(--bg-active)',
    '&::before': {
      content: '""',
      position: 'absolute',
      left: 4,
      top: 6,
      bottom: 6,
      width: 2,
      borderRadius: 2,
      background: 'var(--text-tertiary)',
    },
  },
  '&:hover': {
    background: 'var(--bg-hover)',
  },
})

export const rowIcon = (tone: 'accent' | 'muted' | 'default') =>
  css({
    width: 34,
    height: 34,
    borderRadius: 9,
    border: '1px solid var(--border-subtle)',
    background: 'var(--bg-overlay)',
    display: 'grid',
    placeItems: 'center',
    flexShrink: 0,
    color:
      tone === 'accent'
        ? 'var(--text-secondary)'
        : tone === 'muted'
          ? 'var(--text-disabled)'
          : 'var(--text-tertiary)',
  })

export const rowBody = css({ flex: 1, minWidth: 0 })

export const rowTitle = css({
  display: 'block',
  fontSize: '0.9375rem',
  fontWeight: 500,
  color: 'var(--text-primary)',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const rowSub = css({
  display: 'block',
  marginTop: 2,
  fontSize: '0.75rem',
  color: 'var(--text-disabled)',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
  fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
})

export const rowMark = css({
  background: 'var(--bg-active)',
  color: 'var(--text-primary)',
  borderRadius: 3,
  padding: '0 2px',
  fontWeight: 600,
})

export const rowEnter = css({
  flexShrink: 0,
  color: 'var(--text-disabled)',
  opacity: 0,
  transition: 'opacity 0.1s',
  '[data-active="true"] &': { opacity: 1 },
})

export const emptyWrap = css({
  padding: '36px 24px 32px',
  textAlign: 'center',
})

export const emptyIcon = css({
  width: 44,
  height: 44,
  margin: '0 auto 12px',
  borderRadius: 12,
  border: '1px solid var(--border-subtle)',
  background: 'var(--bg-overlay)',
  display: 'grid',
  placeItems: 'center',
  color: 'var(--text-disabled)',
})

export const emptyTitle = css({
  fontSize: '0.9375rem',
  fontWeight: 600,
  color: 'var(--text-secondary)',
  marginBottom: 6,
})

export const emptyHint = css({
  fontSize: '0.8125rem',
  color: 'var(--text-disabled)',
  lineHeight: 1.45,
})

export const footer = css({
  display: 'flex',
  alignItems: 'center',
  flexWrap: 'wrap',
  gap: '12px 16px',
  padding: '9px 14px',
  borderTop: '1px solid var(--border-subtle)',
  background: 'var(--bg-overlay)',
  flexShrink: 0,
})

export const footerHint = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  fontSize: '0.6875rem',
  color: 'var(--text-disabled)',
})

export const footerMeta = css({
  marginLeft: 'auto',
  fontSize: '0.6875rem',
  color: 'var(--text-disabled)',
})

export const skeletonRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 12,
  margin: '1px 6px',
  padding: '8px 12px 8px 14px',
})

export const skeletonIcon = css({
  width: 34,
  height: 34,
  borderRadius: 9,
  background: 'var(--bg-overlay)',
  border: '1px solid var(--border-subtle)',
  flexShrink: 0,
  animation: 'pulse-subtle 1.2s ease-in-out infinite',
})

export const skeletonLines = css({
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  gap: 6,
})

export const skeletonLine = css({
  height: 10,
  borderRadius: 4,
  background: 'var(--bg-overlay)',
  animation: 'pulse-subtle 1.2s ease-in-out infinite',
})
