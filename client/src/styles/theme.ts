

import { css } from '@emotion/react'

const ease = 'var(--ease-out)'
const fast = 'var(--duration-fast)'
const normal = 'var(--duration-normal)'


export const pageHeader = css({
  marginBottom: 32,
})

export const pageTitle = css({
  fontSize: '1.75rem',
  fontWeight: 700,
  letterSpacing: '-0.03em',
  color: 'var(--text-primary)',
  margin: 0,
  marginBottom: 6,
  lineHeight: 1.2,
})

export const pageDesc = css({
  fontSize: '0.9375rem',
  color: 'var(--text-secondary)',
  margin: 0,
  lineHeight: 1.65,
})

export const breadcrumb = css({
  marginBottom: 24,
  fontSize: '0.8125rem',
  color: 'var(--text-tertiary)',
  display: 'flex',
  alignItems: 'center',
  gap: 4,
  '& a': {
    color: 'var(--text-tertiary)',
    textDecoration: 'none',
    transition: `color ${fast} ${ease}`,
    '&:hover': { color: 'var(--text-primary)' },
  },
})


export const form = css({
  display: 'flex',
  flexDirection: 'column',
  gap: 20,
})

export const formLabel = css({
  display: 'flex',
  flexDirection: 'column',
  gap: 8,
  fontSize: '0.875rem',
  fontWeight: 500,
  color: 'var(--text-secondary)',
})

export const formInput = css({
  padding: '11px 14px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-overlay)',
  color: 'var(--text-primary)',
  fontSize: '0.9375rem',
  transition: `border-color ${fast} ${ease}, box-shadow ${fast} ${ease}`,
  '&:focus': {
    outline: 'none',
    borderColor: 'var(--accent)',
    boxShadow: '0 0 0 3px var(--accent-muted)',
  },
  '&::placeholder': {
    color: 'var(--text-disabled)',
  },
})


export const alertError = css({
  padding: '12px 16px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--danger-muted)',
  border: '1px solid var(--danger-glow)',
  color: 'var(--danger)',
  fontSize: '0.875rem',
  lineHeight: 1.5,
  display: 'flex',
  gap: 8,
  alignItems: 'flex-start',
})

export const alertSuccess = css({
  padding: '12px 16px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--success-muted)',
  border: '1px solid var(--success-glow)',
  color: 'var(--success)',
  fontSize: '0.875rem',
  lineHeight: 1.5,
  display: 'flex',
  gap: 8,
  alignItems: 'flex-start',
})


export const btnPrimary = css({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: 8,
  padding: '10px 20px',
  borderRadius: 'var(--radius-md)',
  border: 'none',
  background: 'var(--accent)',
  color: 'var(--text-on-accent)',
  fontSize: '0.9375rem',
  fontWeight: 600,
  letterSpacing: '-0.01em',
  textDecoration: 'none',
  cursor: 'pointer',
  transition: `background ${fast} ${ease}, box-shadow ${fast} ${ease}, transform ${fast} ${ease}`,
  '&:hover:not(:disabled)': {
    background: 'var(--accent-hover)',
    boxShadow: '0 4px 16px var(--accent-glow)',
    transform: 'translateY(-1px)',
  },
  '&:active:not(:disabled)': {
    transform: 'translateY(0)',
    boxShadow: 'none',
  },
  '&:disabled': {
    cursor: 'not-allowed',
    opacity: 0.55,
    transform: 'none',
  },
})

export const btnSecondary = css({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: 8,
  padding: '9px 16px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'transparent',
  color: 'var(--text-secondary)',
  fontSize: '0.875rem',
  fontWeight: 500,
  textDecoration: 'none',
  cursor: 'pointer',
  transition: `color ${fast} ${ease}, border-color ${fast} ${ease}, background ${fast} ${ease}`,
  '&:hover': {
    color: 'var(--text-primary)',
    borderColor: 'var(--border-strong)',
    background: 'var(--bg-hover)',
  },
})

export const btnDanger = css({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: 6,
  padding: '7px 14px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid transparent',
  background: 'var(--danger-muted)',
  color: 'var(--danger)',
  fontSize: '0.8125rem',
  fontWeight: 500,
  cursor: 'pointer',
  transition: `background ${fast} ${ease}, color ${fast} ${ease}, border-color ${fast} ${ease}`,
  '&:hover:not(:disabled)': {
    background: 'var(--danger)',
    color: 'var(--text-on-accent)',
    borderColor: 'var(--danger)',
  },
  '&:disabled': {
    cursor: 'not-allowed',
    opacity: 0.55,
  },
})

export const btnSecondaryLarge = css({
  padding: '10px 20px',
  fontSize: '0.9375rem',
})

export const btnGhost = css({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: 6,
  padding: '6px 10px',
  borderRadius: 'var(--radius-md)',
  border: 'none',
  background: 'transparent',
  color: 'var(--text-tertiary)',
  fontSize: '0.8125rem',
  fontWeight: 500,
  cursor: 'pointer',
  transition: `color ${fast} ${ease}, background ${fast} ${ease}`,
  '&:hover': {
    color: 'var(--text-primary)',
    background: 'var(--bg-hover)',
  },
})


export const card = css({
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-lg)',
  padding: 24,
  marginBottom: 16,
  boxShadow: 'var(--shadow-sm)',
  transition: `border-color ${fast} ${ease}`,
  '&:last-of-type': { marginBottom: 0 },
  '&:hover': {
    borderColor: 'var(--border-default)',
  },
})

export const cardTitle = css({
  fontSize: '1rem',
  fontWeight: 600,
  letterSpacing: '-0.01em',
  color: 'var(--text-primary)',
  marginBottom: 6,
  display: 'flex',
  alignItems: 'center',
  gap: 8,
})

export const cardDesc = css({
  fontSize: '0.875rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.65,
  marginBottom: 16,
})


export const projectCard = css({
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-lg)',
  padding: '20px 22px',
  display: 'flex',
  flexDirection: 'column',
  gap: 14,
  position: 'relative',
  overflow: 'hidden',
  transition: `border-color ${normal} ${ease}, box-shadow ${normal} ${ease}, transform ${fast} ${ease}`,
  animation: 'fade-in 0.3s var(--ease-out) both',
  '&:hover': {
    borderColor: 'var(--border-default)',
    boxShadow: 'var(--shadow-md)',
    transform: 'translateY(-1px)',
  },
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    height: '1px',
    background: 'linear-gradient(90deg, transparent, var(--accent-muted), transparent)',
    opacity: 0,
    transition: `opacity ${normal} ${ease}`,
  },
  '&:hover::before': {
    opacity: 1,
  },
})

export const projectCardName = css({
  fontSize: '0.9375rem',
  fontWeight: 600,
  letterSpacing: '-0.01em',
  color: 'var(--text-primary)',
  textDecoration: 'none',
  transition: `color ${fast} ${ease}`,
  '&:hover': { color: 'var(--accent)' },
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const projectCardMeta = css({
  display: 'flex',
  alignItems: 'center',
  gap: 6,
  fontSize: '0.8125rem',
  color: 'var(--text-tertiary)',
  fontFamily: 'ui-monospace, monospace',
  overflow: 'hidden',
  '& span': {
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
})

export const branchChip = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 5,
  padding: '3px 9px',
  borderRadius: 'var(--radius-full)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-overlay)',
  fontSize: '0.75rem',
  fontFamily: 'ui-monospace, monospace',
  color: 'var(--text-secondary)',
  flexShrink: 0,
})

export const projectCardFooter = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 8,
  marginTop: 'auto',
  paddingTop: 14,
  borderTop: '1px solid var(--border-subtle)',
})


export const statsRow = css({
  display: 'grid',
  gridTemplateColumns: 'repeat(3, 1fr)',
  gap: 14,
  marginBottom: 32,
})

export const statCard = css({
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-lg)',
  padding: '20px 22px',
  display: 'flex',
  flexDirection: 'column',
  gap: 6,
  animation: 'fade-in 0.25s var(--ease-out) both',
})

export const statLabel = css({
  fontSize: '0.75rem',
  fontWeight: 600,
  color: 'var(--text-tertiary)',
  textTransform: 'uppercase',
  letterSpacing: '0.07em',
})

export const statValue = css({
  fontSize: '2rem',
  fontWeight: 700,
  letterSpacing: '-0.04em',
  color: 'var(--text-primary)',
  lineHeight: 1,
})

export const statDesc = css({
  fontSize: '0.8125rem',
  color: 'var(--text-tertiary)',
})


export const statusBadge = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  padding: '4px 10px',
  borderRadius: 'var(--radius-full)',
  fontSize: '0.75rem',
  fontWeight: 600,
  letterSpacing: '0.01em',
})

export const statusActive = css({
  background: 'var(--success-muted)',
  color: 'var(--success)',
  border: '1px solid var(--success-glow)',
})

export const statusPending = css({
  background: 'var(--warning-muted)',
  color: 'var(--warning)',
  border: '1px solid var(--warning-glow)',
})

export const statusError = css({
  background: 'var(--danger-muted)',
  color: 'var(--danger)',
  border: '1px solid var(--danger-glow)',
})


export const tableWrap = css({
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-lg)',
  overflow: 'hidden',
  boxShadow: 'var(--shadow-sm)',
})

export const table = css({
  width: '100%',
  borderCollapse: 'collapse',
})

export const tableTh = css({
  textAlign: 'left',
  padding: '12px 20px',
  fontSize: '0.6875rem',
  fontWeight: 600,
  color: 'var(--text-tertiary)',
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  background: 'var(--bg-overlay)',
  borderBottom: '1px solid var(--border-subtle)',
  '&:last-of-type': { textAlign: 'right' },
})

export const tableTd = css({
  padding: '15px 20px',
  borderBottom: '1px solid var(--border-subtle)',
  fontSize: '0.9375rem',
  color: 'var(--text-primary)',
  verticalAlign: 'middle',
  '&:last-of-type': { textAlign: 'right' },
})

export const tableTr = css({
  transition: `background ${fast} ${ease}`,
  '&:last-child td': { borderBottom: 'none' },
  '&:hover': { background: 'var(--bg-hover)' },
})


export const emptyState = css({
  padding: '64px 32px',
  textAlign: 'center',
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-lg)',
  animation: 'fade-in 0.3s var(--ease-out) both',
})

export const emptyStateIcon = css({
  width: 60,
  height: 60,
  margin: '0 auto 20px',
  borderRadius: 'var(--radius-lg)',
  background: 'linear-gradient(135deg, var(--bg-overlay), var(--bg-elevated))',
  border: '1px solid var(--border-default)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: 'var(--text-tertiary)',
  fontSize: '1.5rem',
  boxShadow: '0 0 0 6px var(--bg-card), 0 0 0 7px var(--border-subtle)',
})

export const emptyTitle = css({
  fontSize: '1.125rem',
  fontWeight: 600,
  letterSpacing: '-0.02em',
  color: 'var(--text-primary)',
  marginBottom: 8,
})

export const emptyDesc = css({
  fontSize: '0.9375rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.65,
  marginBottom: 28,
  maxWidth: 380,
  marginLeft: 'auto',
  marginRight: 'auto',
})


export const codeBlock = css({
  padding: '14px 18px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--code-bg)',
  border: '1px solid var(--border-default)',
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.8125rem',
  lineHeight: 1.7,
  wordBreak: 'break-all',
  marginBottom: 0,
  color: 'var(--code-text)',
})

export const stepCard = [
  card,
  css({
    '&:last-of-type': { marginBottom: 0 },
  }),
]

export const stepTitle = css({
  fontSize: '0.9375rem',
  fontWeight: 600,
  color: 'var(--text-primary)',
  marginBottom: 4,
  display: 'flex',
  alignItems: 'center',
  gap: 10,
})

export const stepSub = css({
  fontSize: '0.8125rem',
  color: 'var(--text-tertiary)',
  marginBottom: 16,
})

export const instructions = css({
  fontSize: '0.875rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.7,
  marginBottom: 12,
  '& ol': { margin: '0 0 8px 20px', padding: 0 },
  '& li': { marginBottom: 4 },
})


export const toolbar = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  marginBottom: 20,
  flexWrap: 'wrap',
  gap: 12,
})

export const sectionTitle = css({
  fontSize: '0.75rem',
  fontWeight: 600,
  color: 'var(--text-tertiary)',
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
})

export const metaMono = css({
  fontSize: '0.8125rem',
  color: 'var(--text-tertiary)',
  fontFamily: 'ui-monospace, monospace',
})

export const linkPrimary = css({
  color: 'var(--text-primary)',
  fontWeight: 600,
  textDecoration: 'none',
  transition: `color ${fast} ${ease}`,
  '&:hover': { color: 'var(--accent)' },
})

export const loadingText = css({
  fontSize: '0.9375rem',
  color: 'var(--text-tertiary)',
})

export const loadingBlock = css({
  padding: '40px 20px',
  fontSize: '0.9375rem',
  color: 'var(--text-tertiary)',
  textAlign: 'center',
})


export const sectionDivider = css({
  height: 1,
  background: 'var(--border-subtle)',
  margin: '32px 0',
})


export const chip = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 5,
  padding: '3px 10px',
  borderRadius: 'var(--radius-full)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-overlay)',
  fontSize: '0.75rem',
  fontWeight: 500,
  color: 'var(--text-secondary)',
})
