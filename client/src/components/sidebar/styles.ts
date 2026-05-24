import { css } from '@emotion/react'

const dur = 'var(--duration-sidebar)'
const ease = 'var(--ease-sidebar)'
const fast = 'var(--duration-fast)'


export const sidebarRoot = css({
  '--sb-pad-x': '12px',
  '--sb-item-h': '40px',
  '--sb-icon-col': '40px',
  flexShrink: 0,
  position: 'sticky',
  top: 0,
  height: '100vh',
  width: 'var(--sidebar-width-expanded)',
  display: 'flex',
  flexDirection: 'column',
  background: 'var(--bg-elevated)',
  borderRight: '1px solid var(--border-subtle)',
  padding: '12px var(--sb-pad-x)',
  overflowX: 'hidden',
  overflowY: 'auto',
  transition: `width ${dur} ${ease}, padding ${dur} ${ease}`,
  '&[data-collapsed]': {
    '--sb-pad-x': '12px',
    width: 'var(--sidebar-width-collapsed)',
  },
})

export const sidebarHeader = css({
  display: 'flex',
  alignItems: 'center',
  gap: 8,
  minHeight: 44,
  marginBottom: 8,
  flexShrink: 0,
  '.app-sidebar[data-collapsed] &': {
    flexDirection: 'column',
    gap: 10,
    paddingBottom: 4,
  },
})

export const brandLink = css({
  flex: 1,
  minWidth: 0,
  display: 'flex',
  alignItems: 'center',
  textDecoration: 'none',
  '.app-sidebar[data-collapsed] &': {
    flex: 'none',
    justifyContent: 'center',
    width: '100%',
  },
})

export const collapseButton = css({
  flexShrink: 0,
  width: 32,
  height: 32,
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-subtle)',
  background: 'var(--bg-overlay)',
  color: 'var(--text-tertiary)',
  display: 'grid',
  placeItems: 'center',
  cursor: 'pointer',
  appearance: 'none',
  WebkitAppearance: 'none',
  transition: `background ${fast} var(--ease-out), color ${fast} var(--ease-out), border-color ${fast} var(--ease-out)`,
  '&:hover': {
    background: 'var(--bg-hover)',
    color: 'var(--text-primary)',
    borderColor: 'var(--border-default)',
  },
  '.app-sidebar[data-collapsed] &': {
    width: 36,
    height: 36,
  },
})

export const sidebarScroll = css({
  flex: 1,
  minHeight: 0,
  display: 'flex',
  flexDirection: 'column',
  gap: 2,
  margin: '0 calc(-1 * var(--sb-pad-x))',
  padding: '0 var(--sb-pad-x)',
})

export const sectionTitle = css({
  fontSize: '0.6875rem',
  fontWeight: 600,
  color: 'var(--text-disabled)',
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  padding: '14px 10px 6px',
  margin: 0,
  whiteSpace: 'nowrap',
  overflow: 'hidden',
  transition: `opacity ${dur} ${ease}, max-height ${dur} ${ease}, padding ${dur} ${ease}`,
  maxHeight: 32,
  '.app-sidebar[data-collapsed] &': {
    opacity: 0,
    maxHeight: 0,
    paddingTop: 0,
    paddingBottom: 0,
    pointerEvents: 'none',
  },
})

export const divider = css({
  height: 1,
  background: 'var(--border-subtle)',
  margin: '6px 10px',
  flexShrink: 0,
  '.app-sidebar[data-collapsed] &': {
    margin: '6px 4px',
  },
})


const itemBase = css({
  display: 'grid',
  gridTemplateColumns: 'var(--sb-icon-col) 1fr auto',
  alignItems: 'center',
  minHeight: 'var(--sb-item-h)',
  width: '100%',
  padding: 0,
  margin: 0,
  border: 'none',
  borderRadius: 'var(--radius-md)',
  background: 'transparent',
  color: 'var(--text-tertiary)',
  fontSize: '0.875rem',
  fontWeight: 500,
  fontFamily: 'inherit',
  lineHeight: 1.3,
  textDecoration: 'none',
  textAlign: 'left',
  cursor: 'pointer',
  appearance: 'none',
  WebkitAppearance: 'none',
  position: 'relative',
  transition: `color ${fast} var(--ease-out), background ${fast} var(--ease-out)`,
  '&:hover': {
    color: 'var(--text-secondary)',
    background: 'var(--bg-hover)',
  },
  '&.active': {
    color: 'var(--text-primary)',
    background: 'var(--bg-active)',
    '& [data-slot="icon"]': { color: 'var(--accent)' },
    '&::before': {
      content: '""',
      position: 'absolute',
      left: 0,
      top: '22%',
      bottom: '22%',
      width: 2,
      borderRadius: 2,
      background: 'var(--accent)',
    },
  },
  '.app-sidebar[data-collapsed] &': {
    gridTemplateColumns: '1fr',
    justifyItems: 'center',
    width: 'var(--sb-item-h)',
    marginInline: 'auto',
  },
})

export const itemLink = css(itemBase)
export const itemButton = css(itemBase)

export const itemIcon = css({
  width: 'var(--sb-icon-col)',
  height: 'var(--sb-item-h)',
  display: 'grid',
  placeItems: 'center',
  flexShrink: 0,
  color: 'currentColor',
  transition: `color ${fast} var(--ease-out)`,
  '.app-sidebar[data-collapsed] &': {
    width: 'var(--sb-item-h)',
  },
})

export const itemLabel = css({
  minWidth: 0,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
  paddingRight: 4,
  opacity: 1,
  transition: `opacity ${dur} ${ease}`,
  transitionDelay: '50ms',
  '.app-sidebar[data-collapsed] &': {
    opacity: 0,
    width: 0,
    padding: 0,
    overflow: 'hidden',
    position: 'absolute',
    pointerEvents: 'none',
    transitionDelay: '0ms',
  },
})

export const itemSuffix = css({
  flexShrink: 0,
  paddingRight: 8,
  opacity: 1,
  transition: `opacity ${dur} ${ease}`,
  transitionDelay: '50ms',
  '.app-sidebar[data-collapsed] &': {
    opacity: 0,
    width: 0,
    padding: 0,
    overflow: 'hidden',
    position: 'absolute',
    pointerEvents: 'none',
    transitionDelay: '0ms',
  },
})

export const kbHint = css({
  display: 'inline-flex',
  alignItems: 'center',
  padding: '2px 6px',
  borderRadius: 4,
  border: '1px solid var(--border-default)',
  background: 'var(--bg-overlay)',
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.625rem',
  fontWeight: 500,
  color: 'var(--text-disabled)',
  letterSpacing: '0.02em',
})

export const badge = css({
  minWidth: 18,
  height: 18,
  padding: '0 5px',
  borderRadius: 'var(--radius-full)',
  background: 'var(--accent)',
  color: 'var(--text-on-accent)',
  fontSize: '0.625rem',
  fontWeight: 700,
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  boxShadow: '0 0 8px var(--accent-glow)',
})

export const badgeDot = css({
  position: 'absolute',
  top: 8,
  right: 8,
  width: 7,
  height: 7,
  borderRadius: 'var(--radius-full)',
  background: 'var(--accent)',
  boxShadow: '0 0 6px var(--accent-glow)',
  pointerEvents: 'none',
  '.app-sidebar:not([data-collapsed]) &': {
    display: 'none',
  },
})

export const sidebarFooter = css({
  flexShrink: 0,
  display: 'flex',
  flexDirection: 'column',
  gap: 2,
  marginTop: 4,
  paddingTop: 10,
  borderTop: '1px solid var(--border-subtle)',
})

export const prefsRow = css({
  display: 'grid',
  gridTemplateColumns: '1fr 1fr',
  gap: 6,
  padding: '4px 0',
  '.app-sidebar[data-collapsed] &': {
    gridTemplateColumns: '1fr',
    gap: 4,
  },
})

export const prefChip = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  gap: 6,
  minHeight: 34,
  padding: '0 8px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-subtle)',
  background: 'var(--bg-overlay)',
  color: 'var(--text-disabled)',
  fontSize: '0.75rem',
  fontWeight: 500,
  fontFamily: 'inherit',
  appearance: 'none',
  cursor: 'pointer',
  transition: `color ${fast} var(--ease-out), background ${fast} var(--ease-out), border-color ${fast} var(--ease-out)`,
  '&:hover': {
    color: 'var(--text-secondary)',
    borderColor: 'var(--border-default)',
    background: 'var(--bg-hover)',
  },
  '.app-sidebar[data-collapsed] & span[data-pref-label]': {
    display: 'none',
  },
})

export const profileMeta = css([
  itemLabel,
  {
    display: 'flex',
    flexDirection: 'column',
    gap: 1,
    lineHeight: 1.25,
    paddingRight: 4,
  },
])

export const profileName = css({
  fontSize: '0.8125rem',
  fontWeight: 600,
  color: 'var(--text-primary)',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const profileSub = css({
  fontSize: '0.6875rem',
  color: 'var(--text-disabled)',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

export const avatar = css({
  width: 28,
  height: 28,
  borderRadius: 'var(--radius-full)',
  background: 'linear-gradient(135deg, var(--indigo), var(--accent))',
  display: 'grid',
  placeItems: 'center',
  fontSize: '0.6875rem',
  fontWeight: 700,
  color: 'var(--text-on-accent)',
  position: 'relative',
  boxShadow: '0 0 0 1px var(--avatar-inset) inset',
})

export const avatarStatus = css({
  position: 'absolute',
  right: -1,
  bottom: -1,
  width: 8,
  height: 8,
  borderRadius: 'var(--radius-full)',
  background: 'var(--success)',
  border: '2px solid var(--bg-elevated)',
})

export const logoutItem = css(itemBase, {
  color: 'var(--text-disabled)',
  '&:hover': {
    background: 'var(--danger-muted)',
    color: 'var(--danger)',
  },
})

export const notifHost = css({
  position: 'relative',
  width: '100%',
})

export const notifPanel = css({
  position: 'absolute',
  bottom: 'calc(100% + 8px)',
  left: 0,
  right: 0,
  zIndex: 500,
  background: 'var(--bg-elevated)',
  border: '1px solid var(--border-default)',
  borderRadius: 'var(--radius-lg)',
  boxShadow: 'var(--shadow-overlay)',
  overflow: 'hidden',
  maxHeight: 400,
  display: 'flex',
  flexDirection: 'column',
  animation: 'slide-up 0.18s var(--ease-out) both',
  '.app-sidebar[data-collapsed] &': {
    position: 'fixed',
    left: 'calc(var(--sidebar-width-collapsed) + 8px)',
    right: 'auto',
    bottom: 72,
    width: 300,
  },
})
