import { useEffect, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { timeAgo } from '../../api/pipelines'
import { useI18n } from '../../i18n'
import { SidebarItem } from './SidebarItem'
import { IconBell } from './icons'
import * as s from './styles'
import { css } from '@emotion/react'

export type RunNotifStatus = 'running' | 'success' | 'failed' | 'cancelled'

export interface RunNotification {
  id: string
  runId: string
  projectId: string
  status: RunNotifStatus
  receivedAt: number
  read: boolean
}

const STATUS_COLOR: Record<RunNotifStatus, string> = {
  running: 'var(--accent)',
  success: 'var(--success)',
  failed: 'var(--danger)',
  cancelled: 'var(--text-disabled)',
}

function StatusIcon({ status }: { status: RunNotifStatus }) {
  const c = STATUS_COLOR[status]
  if (status === 'success') {
    return (
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden>
        <circle cx="7" cy="7" r="6" stroke={c} strokeWidth="1.2" />
        <path d="M4 7l2.5 2.5 3.5-5" stroke={c} strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
      </svg>
    )
  }
  if (status === 'failed') {
    return (
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden>
        <circle cx="7" cy="7" r="6" stroke={c} strokeWidth="1.2" />
        <path d="M4.5 4.5l5 5M9.5 4.5l-5 5" stroke={c} strokeWidth="1.3" strokeLinecap="round" />
      </svg>
    )
  }
  if (status === 'cancelled') {
    return (
      <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden>
        <circle cx="7" cy="7" r="6" stroke={c} strokeWidth="1.2" />
        <path d="M5 9l2-2m0 0l2-2M7 7L5 5m2 2l2 2" stroke={c} strokeWidth="1.3" strokeLinecap="round" />
      </svg>
    )
  }
  return (
    <svg width="14" height="14" viewBox="0 0 14 14" fill="none" aria-hidden>
      <circle cx="7" cy="7" r="6" stroke={c} strokeWidth="1.2" />
      <path d="M5 7h4M7 5v4" stroke={c} strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  )
}

const panelHeader = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: '10px 12px 8px',
  borderBottom: '1px solid var(--border-subtle)',
})

const panelTitle = css({
  fontSize: '0.75rem',
  fontWeight: 700,
  color: 'var(--text-secondary)',
  letterSpacing: '0.04em',
  textTransform: 'uppercase',
})

const markAll = css({
  fontSize: '0.6875rem',
  color: 'var(--accent)',
  background: 'none',
  border: 'none',
  cursor: 'pointer',
  padding: '2px 4px',
  '&:hover': { opacity: 0.75 },
})

const list = css({
  overflowY: 'auto',
  maxHeight: 340,
})

const empty = css({
  padding: '24px 12px',
  textAlign: 'center',
  color: 'var(--text-disabled)',
  fontSize: '0.8125rem',
})

const row = css({
  display: 'flex',
  gap: 9,
  padding: '9px 12px',
  borderBottom: '1px solid var(--border-subtle)',
  cursor: 'pointer',
  '&:last-of-type': { borderBottom: 'none' },
  '&:hover': { background: 'var(--bg-hover)' },
})

const unread = css({
  width: 6,
  height: 6,
  borderRadius: 'var(--radius-full)',
  background: 'var(--accent)',
  flexShrink: 0,
  marginTop: 6,
})

interface Props {
  notifications: RunNotification[]
  panelOpen: boolean
  onToggle: () => void
  onMarkAllRead: (e: React.MouseEvent) => void
  onUpdate: (fn: (prev: RunNotification[]) => RunNotification[]) => void
}

export function NotificationBell({ notifications, panelOpen, onToggle, onMarkAllRead, onUpdate }: Props) {
  const navigate = useNavigate()
  const { t, lang } = useI18n()
  const hostRef = useRef<HTMLDivElement>(null)

  const unreadCount = notifications.filter((n) => !n.read).length

  const notifLabel = (status: RunNotifStatus) =>
    ({ running: t.notif.running, success: t.notif.success, failed: t.notif.failed, cancelled: t.notif.cancelled })[status]

  useEffect(() => {
    if (!panelOpen) return
    const onDown = (e: MouseEvent) => {
      if (hostRef.current && !hostRef.current.contains(e.target as Node)) onToggle()
    }
    document.addEventListener('mousedown', onDown)
    return () => document.removeEventListener('mousedown', onDown)
  }, [panelOpen, onToggle])

  const onItemClick = (n: RunNotification) => {
    onUpdate((prev) => prev.map((x) => (x.id === n.id ? { ...x, read: true } : x)))
    onToggle()
    if (n.projectId && n.runId) navigate(`/projects/${n.projectId}/runs/${n.runId}`)
  }

  return (
    <div css={s.notifHost} ref={hostRef}>
      {panelOpen && (
        <div css={s.notifPanel} role="dialog" aria-label={t.nav.notifications}>
          <div css={panelHeader}>
            <span css={panelTitle}>{t.nav.notifications}</span>
            {unreadCount > 0 && (
              <button type="button" css={markAll} onClick={onMarkAllRead}>{t.nav.markAllRead}</button>
            )}
          </div>
          <div css={list}>
            {notifications.length === 0 ? (
              <div css={empty}>{t.nav.noNotifications}</div>
            ) : (
              notifications.map((n) => (
                <div
                  key={n.id}
                  css={row}
                  role="button"
                  tabIndex={0}
                  onClick={() => onItemClick(n)}
                  onKeyDown={(e) => e.key === 'Enter' && onItemClick(n)}
                >
                  <StatusIcon status={n.status} />
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div style={{ fontSize: '0.8125rem', fontWeight: 500 }}>{notifLabel(n.status)}</div>
                    <div style={{ fontSize: '0.6875rem', color: 'var(--text-disabled)' }}>
                      {n.runId ? `#${n.runId.slice(0, 7)} · ` : ''}
                      {timeAgo(new Date(n.receivedAt).toISOString(), lang)}
                    </div>
                  </div>
                  {!n.read && <span css={unread} />}
                </div>
              ))
            )}
          </div>
        </div>
      )}
      <div style={{ position: 'relative' }}>
        <SidebarItem
          icon={<IconBell />}
          label={t.nav.notifications}
          onClick={onToggle}
          active={panelOpen}
          suffix={
            unreadCount > 0 ? (
              <span css={s.badge}>{unreadCount > 9 ? '9+' : unreadCount}</span>
            ) : undefined
          }
        />
        {unreadCount > 0 && <span css={s.badgeDot} aria-hidden />}
      </div>
    </div>
  )
}
