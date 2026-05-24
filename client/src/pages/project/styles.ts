import { css } from '@emotion/react'

export const pageWrap = css({
  animation: 'fade-in 0.25s var(--ease-out) both',
})

export const hero = css({
  position: 'relative',
  marginBottom: 24,
  padding: '22px 24px',
  borderRadius: 'var(--radius-lg)',
  border: '1px solid var(--border-subtle)',
  background: 'linear-gradient(165deg, var(--bg-card) 0%, var(--bg-overlay) 100%)',
  overflow: 'hidden',
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    height: 1,
    background: 'linear-gradient(90deg, transparent, var(--border-strong), transparent)',
    opacity: 0.6,
  },
})

export const heroTop = css({
  display: 'flex',
  alignItems: 'center',
  gap: 16,
  flexWrap: 'wrap',
  width: '100%',
})

export const projectAvatar = css({
  width: 48,
  height: 48,
  borderRadius: 12,
  flexShrink: 0,
  display: 'grid',
  placeItems: 'center',
  fontSize: '1.125rem',
  fontWeight: 700,
  letterSpacing: '-0.02em',
  color: 'var(--text-primary)',
  background: 'var(--bg-overlay)',
  border: '1px solid var(--border-default)',
})

export const heroMain = css({ flex: 1, minWidth: 0 })

export const heroTitleRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  flexWrap: 'wrap',
  marginBottom: 10,
})

export const heroTitle = css({
  margin: 0,
  fontSize: '1.5rem',
  fontWeight: 700,
  letterSpacing: '-0.03em',
  color: 'var(--text-primary)',
  lineHeight: 1.2,
})

export const heroMeta = css({
  display: 'flex',
  alignItems: 'center',
  flexWrap: 'wrap',
  gap: '8px 12px',
})

export const repoLink = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  maxWidth: 'min(100%, 480px)',
  fontSize: '0.8125rem',
  fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
  color: 'var(--text-tertiary)',
  textDecoration: 'none',
  transition: 'color 0.15s',
  '&:hover': { color: 'var(--text-secondary)' },
})

export const repoText = css({
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const roleBadge = css({
  display: 'inline-flex',
  alignItems: 'center',
  padding: '3px 8px',
  borderRadius: 'var(--radius-full)',
  fontSize: '0.6875rem',
  fontWeight: 600,
  letterSpacing: '0.04em',
  textTransform: 'uppercase',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-overlay)',
  color: 'var(--text-tertiary)',
})

export const heroActions = css({
  display: 'flex',
  alignItems: 'center',
  flexWrap: 'wrap',
  gap: 8,
  flexShrink: 0,
  marginLeft: 'auto',
  '@media (max-width: 720px)': {
    marginLeft: 0,
    width: '100%',
    paddingTop: 4,
  },
})

export const projectNotice = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
  flexWrap: 'wrap',
  marginBottom: 20,
  padding: '12px 16px',
  borderRadius: 'var(--radius-lg)',
  border: '1px solid color-mix(in srgb, var(--warning) 35%, var(--border-subtle))',
  background: 'var(--warning-muted)',
  fontSize: '0.8125rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.45,
})

export const projectNoticeLink = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  flexShrink: 0,
  padding: '7px 12px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--bg-elevated)',
  border: '1px solid var(--border-default)',
  color: 'var(--text-primary)',
  fontSize: '0.8125rem',
  fontWeight: 600,
  textDecoration: 'none',
  transition: 'background 0.15s, border-color 0.15s',
  '&:hover': {
    background: 'var(--bg-hover)',
    borderColor: 'var(--border-strong)',
  },
  '& svg': { display: 'block' },
})

export const btnGhost = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  padding: '8px 14px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-overlay)',
  color: 'var(--text-secondary)',
  fontSize: '0.8125rem',
  fontWeight: 500,
  fontFamily: 'inherit',
  textDecoration: 'none',
  cursor: 'pointer',
  transition: 'background 0.15s, border-color 0.15s, color 0.15s, box-shadow 0.15s',
  '&:hover': {
    background: 'var(--bg-hover)',
    color: 'var(--text-primary)',
    borderColor: 'var(--border-strong)',
  },
  '& svg': { display: 'block', flexShrink: 0 },
  '&:focus-visible': {
    outline: '2px solid var(--border-strong)',
    outlineOffset: 2,
  },
})

export const btnDangerGhost = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  padding: '8px 12px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid transparent',
  background: 'transparent',
  color: 'var(--text-disabled)',
  fontSize: '0.8125rem',
  fontWeight: 500,
  cursor: 'pointer',
  transition: 'color 0.15s, background 0.15s, border-color 0.15s',
  '&:hover:not(:disabled)': {
    color: 'var(--danger)',
    background: 'var(--danger-muted)',
    borderColor: 'color-mix(in srgb, var(--danger) 25%, transparent)',
  },
  '&:disabled': { opacity: 0.5, cursor: 'not-allowed' },
})

export const tabShell = css({
  marginBottom: 24,
  maxWidth: '100%',
})

export const tabList = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 4,
  padding: 4,
  borderRadius: 'var(--radius-lg)',
  border: '1px solid var(--border-subtle)',
  background: 'var(--bg-overlay)',
  maxWidth: '100%',
  boxSizing: 'border-box',
  '@media (max-width: 640px)': {
    display: 'flex',
    width: '100%',
    overflowX: 'auto',
    WebkitOverflowScrolling: 'touch',
    scrollbarWidth: 'none',
    scrollPaddingInline: 4,
    '&::-webkit-scrollbar': { display: 'none' },
  },
})

export const tabPanel = css({
  animation: 'fade-in 0.22s var(--ease-out) both',
  outline: 'none',
})

export const tabBtn = (active: boolean) =>
  css({
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'center',
    gap: 7,
    padding: '8px 14px',
    borderRadius: 'var(--radius-md)',
    border: 'none',
    background: active ? 'var(--bg-elevated)' : 'transparent',
    boxShadow: active ? '0 1px 2px rgba(0,0,0,0.12), 0 0 0 1px var(--border-subtle)' : 'none',
    color: active ? 'var(--text-primary)' : 'var(--text-tertiary)',
    fontSize: '0.8125rem',
    fontWeight: active ? 600 : 500,
    cursor: 'pointer',
    fontFamily: 'inherit',
    flexShrink: 0,
    whiteSpace: 'nowrap',
    transition: 'background 0.15s, color 0.15s, box-shadow 0.15s',
    '&:hover': { color: 'var(--text-primary)' },
    '& svg': { display: 'block', flexShrink: 0 },
    '&:focus-visible': {
      outline: '2px solid var(--border-strong)',
      outlineOffset: 2,
    },
  })

export const tabCount = css({
  minWidth: 18,
  height: 18,
  padding: '0 5px',
  borderRadius: 'var(--radius-full)',
  fontSize: '0.6875rem',
  fontWeight: 600,
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: 'var(--bg-active)',
  color: 'var(--text-secondary)',
})

export const panel = css({
  borderRadius: 'var(--radius-lg)',
  border: '1px solid var(--border-subtle)',
  background: 'var(--bg-card)',
  overflow: 'hidden',
})

export const panelHeader = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
  padding: '14px 18px',
  borderBottom: '1px solid var(--border-subtle)',
  background: 'var(--bg-overlay)',
  flexWrap: 'wrap',
  '@media (max-width: 520px)': {
    flexDirection: 'column',
    alignItems: 'stretch',
  },
})

export const panelHeaderActions = css({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  flexShrink: 0,
  marginLeft: 'auto',
  '@media (max-width: 520px)': {
    marginLeft: 0,
    width: '100%',
    justifyContent: 'space-between',
  },
})

export const panelTitle = css({
  margin: 0,
  fontSize: '0.875rem',
  fontWeight: 600,
  color: 'var(--text-primary)',
  display: 'flex',
  alignItems: 'center',
  gap: 8,
  minWidth: 0,
  '& svg': { display: 'block', flexShrink: 0 },
})

export const panelMeta = css({
  fontSize: '0.75rem',
  fontWeight: 500,
  color: 'var(--text-disabled)',
  fontVariantNumeric: 'tabular-nums',
})

export const runsTableHead = css({
  display: 'grid',
  gridTemplateColumns: '120px minmax(0, 1fr) 88px 100px 28px',
  gap: 12,
  padding: '10px 18px',
  borderBottom: '1px solid var(--border-subtle)',
  fontSize: '0.6875rem',
  fontWeight: 600,
  letterSpacing: '0.06em',
  textTransform: 'uppercase',
  color: 'var(--text-disabled)',
  '@media (max-width: 640px)': {
    display: 'none',
  },
})

export const runRow = css({
  display: 'grid',
  gridTemplateColumns: '120px minmax(0, 1fr) 88px 100px 28px',
  gap: 12,
  alignItems: 'center',
  padding: '12px 18px',
  borderBottom: '1px solid var(--border-subtle)',
  textDecoration: 'none',
  color: 'inherit',
  transition: 'background 0.12s',
  '&:last-child': { borderBottom: 'none' },
  '&:hover': { background: 'var(--bg-hover)' },
  '&:focus-visible': {
    outline: '2px solid var(--border-strong)',
    outlineOffset: -2,
    background: 'var(--bg-hover)',
  },
  '@media (max-width: 640px)': {
    gridTemplateColumns: '1fr auto',
    gridTemplateRows: 'auto auto',
    gap: '8px 12px',
  },
})

export const emptyPanelBody = css({
  padding: '12px 8px 28px',
})

export const runBranch = css({
  fontSize: '0.875rem',
  fontWeight: 500,
  color: 'var(--text-primary)',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const runCommit = css({
  fontSize: '0.75rem',
  color: 'var(--text-tertiary)',
  fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace',
})

export const runTime = css({
  fontSize: '0.75rem',
  color: 'var(--text-disabled)',
  textAlign: 'right',
  whiteSpace: 'nowrap',
  '@media (max-width: 640px)': {
    gridColumn: '2',
    gridRow: '1 / 3',
    alignSelf: 'center',
  },
})

export const runDuration = css({
  fontSize: '0.75rem',
  color: 'var(--text-tertiary)',
  fontVariantNumeric: 'tabular-nums',
  textAlign: 'right',
})

export const settingsLayout = css({
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 240px) minmax(0, 1fr)',
  gap: 20,
  alignItems: 'start',
  minWidth: 0,
  '@media (max-width: 900px)': {
    gridTemplateColumns: '1fr',
    gap: 16,
  },
})

export const settingsNav = css({
  display: 'flex',
  flexDirection: 'column',
  gap: 2,
  padding: 8,
  borderRadius: 'var(--radius-lg)',
  border: '1px solid var(--border-subtle)',
  background: 'var(--bg-overlay)',
  position: 'sticky',
  top: 16,
  minWidth: 0,
  maxWidth: '100%',
  boxSizing: 'border-box',
  overflow: 'hidden',
  '@media (max-width: 900px)': {
    position: 'static',
    flexDirection: 'row',
    flexWrap: 'nowrap',
    gap: 6,
    padding: '6px 8px',
    overflowX: 'auto',
    overflowY: 'hidden',
    WebkitOverflowScrolling: 'touch',
    scrollbarWidth: 'none',
    scrollPaddingInline: 8,
    '&::-webkit-scrollbar': { display: 'none' },
  },
})

export const settingsNavIcon = css({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  width: 18,
  height: 18,
  flexShrink: 0,
  color: 'inherit',
  opacity: 0.9,
  '& svg': { display: 'block' },
})

export const settingsNavLabel = css({
  minWidth: 0,
  lineHeight: 1.3,
  textAlign: 'left',
})

export const settingsNavBtn = (active: boolean) =>
  css({
    display: 'inline-flex',
    alignItems: 'center',
    justifyContent: 'flex-start',
    gap: 10,
    width: '100%',
    minWidth: 0,
    maxWidth: '100%',
    boxSizing: 'border-box',
    padding: '10px 12px',
    borderRadius: 'var(--radius-md)',
    border: 'none',
    background: active ? 'var(--bg-elevated)' : 'transparent',
    color: active ? 'var(--text-primary)' : 'var(--text-tertiary)',
    fontSize: '0.8125rem',
    fontWeight: active ? 600 : 500,
    cursor: 'pointer',
    fontFamily: 'inherit',
    whiteSpace: 'nowrap',
    transition: 'background 0.12s, color 0.12s, box-shadow 0.12s',
    boxShadow: active ? '0 0 0 1px var(--border-subtle)' : 'none',
    '&:hover': {
      color: 'var(--text-primary)',
      background: active ? 'var(--bg-elevated)' : 'var(--bg-hover)',
    },
    '&:focus-visible': {
      outline: '2px solid var(--border-strong)',
      outlineOffset: 2,
    },
    '@media (max-width: 900px)': {
      width: 'auto',
      flexShrink: 0,
      padding: '8px 14px',
    },
  })


export const periodBar = tabList

export const periodBtn = tabBtn

export const analyticsPeriodWrap = css({
  marginBottom: 24,
})

export const settingsPanel = css({
  borderRadius: 'var(--radius-lg)',
  border: '1px solid var(--border-subtle)',
  background: 'var(--bg-card)',
  overflow: 'hidden',
})

export const settingsPanelHead = css({
  padding: '16px 20px',
  borderBottom: '1px solid var(--border-subtle)',
  background: 'var(--bg-overlay)',
})

export const settingsPanelTitle = css({
  margin: 0,
  fontSize: '0.9375rem',
  fontWeight: 600,
  color: 'var(--text-primary)',
  lineHeight: 1.35,
  letterSpacing: '-0.01em',
})

export const settingsPanelDesc = css({
  margin: '6px 0 0',
  fontSize: '0.8125rem',
  lineHeight: 1.5,
  color: 'var(--text-tertiary)',
})

export const settingsPanelBody = css({
  padding: '20px',
})

export const loadingBlock = css({
  padding: '32px 20px',
  textAlign: 'center',
  fontSize: '0.875rem',
  color: 'var(--text-tertiary)',
})
