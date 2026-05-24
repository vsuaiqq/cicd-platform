import { Outlet, useNavigate } from 'react-router-dom'
import { css } from '@emotion/react'
import { useCallback, useEffect, useRef, useState } from 'react'
import { useAppDispatch, useAppSelector } from '../store'
import { apiSlice } from '../store/api/apiSlice'
import { logout } from '../store/authSlice'
import { connectGlobalWS, type WSEvent } from '../api/pipelines'
import { useToast } from './ui'
import CommandPalette from './CommandPalette'
import HelpDocsModal from './HelpDocsModal'
import SettingsModal, { loadNotifPrefs, type NotifPrefs } from './SettingsModal'
import {
  AppSidebar,
  loadSidebarCollapsed,
  saveSidebarCollapsed,
  type RunNotification,
  type RunNotifStatus,
} from './sidebar'
import { useI18n } from '../i18n'
import { useTheme } from '../lib/themeContext'

const shell = css({
  display: 'flex',
  minHeight: '100vh',
  background: 'var(--bg-base)',
})

const main = css({
  flex: 1,
  minWidth: 0,
  padding: '36px 48px 64px',
  maxWidth: 1160,
})

const MAX_RECONNECT_DELAY = 30_000
const MAX_RETRIES = 10

function getInitials(label: string): string {
  const parts = label.trim().split(/\s+/).filter(Boolean)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  if (parts.length === 1 && parts[0].length >= 2) return parts[0].slice(0, 2).toUpperCase()
  return label.slice(0, 2).toUpperCase() || '?'
}

export default function AppLayout() {
  const dispatch = useAppDispatch()
  const invalidateRuns = useCallback(
    (runId: string, projectId?: string) => {
      const tags: { type: 'Run' | 'Runs'; id: string }[] = [{ type: 'Run', id: runId }]
      if (projectId) tags.push({ type: 'Runs', id: projectId })
      dispatch(apiSlice.util.invalidateTags(tags))
    },
    [dispatch],
  )
  const navigate = useNavigate()
  const toast = useToast()
  const { userId, email, username, accessToken } = useAppSelector((s) => s.auth)
  const [sidebarCollapsed, setSidebarCollapsed] = useState(loadSidebarCollapsed)
  const { t, lang, setLang } = useI18n()
  const { theme, setTheme, toggleTheme } = useTheme()

  const notifLabel = (status: RunNotifStatus): string =>
    ({
      running: t.notif.running,
      success: t.notif.success,
      failed: t.notif.failed,
      cancelled: t.notif.cancelled,
    })[status]

  const [searchOpen, setSearchOpen] = useState(false)
  const [helpOpen, setHelpOpen] = useState(false)
  const [settingsOpen, setSettingsOpen] = useState(false)
  const [notifPrefs, setNotifPrefs] = useState<NotifPrefs>(loadNotifPrefs)
  const [notifications, setNotifications] = useState<RunNotification[]>([])
  const [panelOpen, setPanelOpen] = useState(false)

  const profileDisplayName = username?.trim() || email?.trim() || userId || '—'
  const profileSubline = email?.trim() || userId || ''

  const toggleSidebar = useCallback(() => {
    setSidebarCollapsed((prev) => {
      const next = !prev
      saveSidebarCollapsed(next)
      if (next) setPanelOpen(false)
      return next
    })
  }, [])

  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
        e.preventDefault()
        setSearchOpen((v) => !v)
      }
    }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [])

  const wsRef = useRef<WebSocket | null>(null)
  const retryCountRef = useRef(0)
  const retryTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined)
  const tokenRef = useRef(accessToken)
  const mountEpochRef = useRef(0)
  const handleRunEventRef = useRef<(event: WSEvent) => void>(() => {})
  const notifPrefsRef = useRef(notifPrefs)

  useEffect(() => { tokenRef.current = accessToken }, [accessToken])
  useEffect(() => { notifPrefsRef.current = notifPrefs }, [notifPrefs])

  const handleRunEvent = useCallback(
    (event: WSEvent) => {
      if (event.type !== 'run_updated') return
      const status = event.status as RunNotifStatus
      if (!(['running', 'success', 'failed', 'cancelled'] as RunNotifStatus[]).includes(status)) return

      const notif: RunNotification = {
        id: `${event.run_id}-${status}`,
        runId: event.run_id,
        projectId: event.project_id ?? '',
        status,
        receivedAt: Date.now(),
        read: false,
      }

      setNotifications((prev) => {
        const filtered = prev.filter((n) => n.id !== notif.id)
        return [notif, ...filtered].slice(0, 50)
      })

      if (notifPrefsRef.current[status]) {
        const variant =
          status === 'running' ? 'info' : status === 'success' ? 'success' : status === 'failed' ? 'error' : 'warning'
        toast.addToast(notifLabel(status), variant, 6000)
      }

      if (['success', 'failed', 'cancelled'].includes(status)) {
        invalidateRuns(event.run_id, event.project_id)
      }
    },
    [toast, t, invalidateRuns],
  )

  useEffect(() => { handleRunEventRef.current = handleRunEvent }, [handleRunEvent])

  const connect = useCallback(() => {
    if (!tokenRef.current) return
    const epoch = mountEpochRef.current
    const ws = connectGlobalWS({
      token: tokenRef.current,
      onEvent: (ev) => handleRunEventRef.current(ev),
      onOpen: () => { retryCountRef.current = 0 },
      onClose: () => {
        wsRef.current = null
        if (epoch !== mountEpochRef.current) return
        if (!tokenRef.current) return
        if (retryCountRef.current >= MAX_RETRIES) return
        const delay = Math.min(1000 * 2 ** retryCountRef.current, MAX_RECONNECT_DELAY)
        retryCountRef.current++
        retryTimerRef.current = setTimeout(() => {
          if (epoch === mountEpochRef.current) connect()
        }, delay)
      },
    })
    wsRef.current = ws
  }, [])

  useEffect(() => {
    if (!accessToken) {
      mountEpochRef.current++
      clearTimeout(retryTimerRef.current)
      wsRef.current?.close()
      wsRef.current = null
      return
    }
    mountEpochRef.current++
    retryCountRef.current = 0
    connect()
    return () => {
      mountEpochRef.current++
      clearTimeout(retryTimerRef.current)
      wsRef.current?.close()
      wsRef.current = null
    }
  }, [accessToken, connect])

  const markAllRead = (e: React.MouseEvent) => {
    e.stopPropagation()
    setNotifications((prev) => prev.map((n) => ({ ...n, read: true })))
  }

  const handleLogout = () => {
    dispatch(logout())
    navigate('/login')
  }

  return (
    <div css={shell}>
      <AppSidebar
        collapsed={sidebarCollapsed}
        onToggleCollapse={toggleSidebar}
        profileDisplayName={profileDisplayName}
        profileSubline={profileSubline}
        profileInitials={getInitials(profileDisplayName)}
        connectedTitle={t.aria.connected}
        theme={theme}
        onToggleTheme={toggleTheme}
        lang={lang}
        onToggleLang={() => setLang(lang === 'en' ? 'ru' : 'en')}
        onOpenSearch={() => setSearchOpen(true)}
        onOpenSettings={() => setSettingsOpen(true)}
        onOpenHelp={() => setHelpOpen(true)}
        onLogout={handleLogout}
        notifications={notifications}
        notifPanelOpen={panelOpen}
        onNotifToggle={() => setPanelOpen((v) => !v)}
        onNotifMarkAllRead={markAllRead}
        onNotifUpdate={setNotifications}
      />

      <main css={main}>
        <Outlet />
      </main>

      <CommandPalette
        open={searchOpen}
        onClose={() => setSearchOpen(false)}
        onOpenSettings={() => {
          setSearchOpen(false)
          setSettingsOpen(true)
        }}
      />
      <HelpDocsModal open={helpOpen} onClose={() => setHelpOpen(false)} />
      <SettingsModal
        open={settingsOpen}
        onClose={() => setSettingsOpen(false)}
        userId={userId}
        onPrefsChange={setNotifPrefs}
        theme={theme}
        onThemeChange={setTheme}
        lang={lang}
        onLangChange={setLang}
      />
    </div>
  )
}
