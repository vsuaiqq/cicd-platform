import { useEffect, useState } from 'react'
import { css } from '@emotion/react'
import { useI18n } from '../i18n'
import type { Theme } from '../lib/theme'
import type { Lang } from '../i18n'



const PREFS_KEY = 'cicd_notif_prefs'

export interface NotifPrefs {
  running:   boolean
  success:   boolean
  failed:    boolean
  cancelled: boolean
}

export function loadNotifPrefs(): NotifPrefs {
  try {
    const raw = localStorage.getItem(PREFS_KEY)
    if (raw) return { running: true, success: true, failed: true, cancelled: true, ...JSON.parse(raw) }
  } catch {  }
  return { running: true, success: true, failed: true, cancelled: true }
}

export function saveNotifPrefs(prefs: NotifPrefs): void {
  localStorage.setItem(PREFS_KEY, JSON.stringify(prefs))
}



const backdrop = css({
  position: 'fixed',
  inset: 0,
  background: 'var(--overlay-backdrop)',
  backdropFilter: 'blur(4px)',
  zIndex: 8000,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: 24,
  animation: 'modal-backdrop 0.15s ease both',
})

const dialog = css({
  width: '100%',
  maxWidth: 460,
  background: 'var(--bg-elevated)',
  border: '1px solid var(--border-default)',
  borderRadius: 'var(--radius-xl)',
  boxShadow: 'var(--shadow-overlay), 0 0 0 1px var(--overlay-inset) inset',
  overflow: 'hidden',
  animation: 'modal-slide 0.18s var(--ease-out) both',
})

const modalHeader = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: '16px 22px',
  borderBottom: '1px solid var(--border-subtle)',
})

const modalTitle = css({
  display: 'flex',
  alignItems: 'center',
  gap: 9,
  fontSize: '0.9375rem',
  fontWeight: 700,
  letterSpacing: '-0.02em',
  color: 'var(--text-primary)',
})

const closeBtn = css({
  width: 30,
  height: 30,
  borderRadius: 7,
  border: '1px solid var(--border-subtle)',
  background: 'transparent',
  color: 'var(--text-disabled)',
  cursor: 'pointer',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  transition: 'background 120ms, color 120ms',
  '&:hover': { background: 'var(--bg-hover)', color: 'var(--text-secondary)' },
})

const body = css({
  padding: '22px',
  display: 'flex',
  flexDirection: 'column',
  gap: 20,
  maxHeight: '70vh',
  overflowY: 'auto',
})

const section = css({
  display: 'flex',
  flexDirection: 'column',
  gap: 10,
})

const sectionLabel = css({
  fontSize: '0.6875rem',
  fontWeight: 700,
  textTransform: 'uppercase',
  letterSpacing: '0.09em',
  color: 'var(--text-disabled)',
})

const accountCard = css({
  display: 'flex',
  alignItems: 'center',
  gap: 12,
  padding: '12px 14px',
  borderRadius: 'var(--radius-lg)',
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
})

const avatar = css({
  width: 38,
  height: 38,
  borderRadius: 'var(--radius-full)',
  background: 'linear-gradient(135deg, var(--indigo), var(--accent))',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  fontSize: '0.875rem',
  fontWeight: 700,
  color: 'var(--text-on-accent)',
  flexShrink: 0,
  boxShadow: '0 0 0 2px var(--avatar-inset) inset',
})

const accountInfo = css({ flex: 1, minWidth: 0 })

const accountEmail = css({
  fontSize: '0.875rem',
  fontWeight: 600,
  color: 'var(--text-primary)',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

const accountRole = css({
  fontSize: '0.75rem',
  color: 'var(--text-disabled)',
  marginTop: 2,
  display: 'flex',
  alignItems: 'center',
  gap: 5,
})

const statusDot = css({
  width: 6,
  height: 6,
  borderRadius: '50%',
  background: 'var(--success)',
  boxShadow: '0 0 4px var(--success-glow)',
  flexShrink: 0,
})

const divider = css({
  height: 1,
  background: 'var(--border-subtle)',
})

const toggleRow = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: '10px 14px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
  gap: 12,
})

const toggleInfo = css({ flex: 1, minWidth: 0 })

const toggleLabel = css({
  fontSize: '0.875rem',
  fontWeight: 500,
  color: 'var(--text-primary)',
  marginBottom: 2,
})


const segmentRow = css({
  display: 'flex',
  gap: 6,
  padding: '4px',
  background: 'var(--bg-overlay)',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-subtle)',
})

const segmentBtn = (active: boolean) => css({
  flex: 1,
  padding: '6px 12px',
  borderRadius: 6,
  border: 'none',
  background: active ? 'var(--bg-elevated)' : 'transparent',
  color: active ? 'var(--text-primary)' : 'var(--text-disabled)',
  fontSize: '0.8125rem',
  fontWeight: active ? 600 : 400,
  cursor: 'pointer',
  transition: 'background 150ms, color 150ms',
  boxShadow: active ? 'var(--shadow-sm)' : 'none',
  '&:hover': { color: active ? 'var(--text-primary)' : 'var(--text-secondary)' },
})

function Toggle({ checked, onChange }: { checked: boolean; onChange: (v: boolean) => void }) {
  return (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      onClick={() => onChange(!checked)}
      css={css({
        width: 38,
        height: 22,
        borderRadius: 11,
        cursor: 'pointer',
        transition: 'background 200ms, box-shadow 200ms',
        position: 'relative',
        flexShrink: 0,
        outline: 'none',
        background: checked ? 'var(--accent)' : 'var(--bg-overlay)',
        boxShadow: checked ? '0 0 8px var(--accent-glow)' : 'none',
        border: checked ? 'none' : '1px solid var(--border-default)',
        '&:focus-visible': { outline: '2px solid var(--accent)', outlineOffset: 2 },
      })}
    >
      <span css={css({
        position: 'absolute',
        top: 2,
        left: checked ? 18 : 2,
        width: 18,
        height: 18,
        borderRadius: '50%',
        background: 'var(--toggle-thumb)',
        transition: 'left 200ms',
        boxShadow: '0 1px 3px var(--toggle-thumb-shadow)',
      })} />
    </button>
  )
}

const notifColorDot = (color: string) => css({
  width: 8,
  height: 8,
  borderRadius: '50%',
  background: color,
  flexShrink: 0,
})

const footer = css({
  padding: '14px 22px',
  borderTop: '1px solid var(--border-subtle)',
  display: 'flex',
  justifyContent: 'flex-end',
  gap: 8,
})

const btnSecondary = css({
  padding: '8px 16px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'transparent',
  color: 'var(--text-secondary)',
  fontSize: '0.875rem',
  fontWeight: 500,
  cursor: 'pointer',
  transition: 'border-color 150ms, background 150ms',
  '&:hover': { background: 'var(--bg-hover)', borderColor: 'var(--border-strong)' },
})

const btnPrimary = css({
  padding: '8px 18px',
  borderRadius: 'var(--radius-md)',
  border: 'none',
  background: 'var(--accent)',
  color: 'var(--text-on-accent)',
  fontSize: '0.875rem',
  fontWeight: 600,
  cursor: 'pointer',
  transition: 'opacity 150ms',
  '&:hover': { opacity: 0.88 },
})

const savedPill = css({
  fontSize: '0.75rem',
  fontWeight: 600,
  color: 'var(--success)',
  display: 'flex',
  alignItems: 'center',
  gap: 5,
  marginRight: 'auto',
  animation: 'fade-in 0.2s ease both',
})

const prefRow = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 16,
})

const prefRowLabel = css({
  fontSize: '0.875rem',
  fontWeight: 500,
  color: 'var(--text-secondary)',
})



interface SettingsModalProps {
  open: boolean
  onClose: () => void
  userId: string | null
  onPrefsChange: (prefs: NotifPrefs) => void
  theme?: Theme
  onThemeChange?: (theme: Theme) => void
  lang?: Lang
  onLangChange?: (lang: Lang) => void
}

function getInitials(id: string | null | undefined): string {
  if (!id) return '?'
  const clean = id.replace(/[^a-zA-Z0-9@.]/g, '')
  const atIdx = clean.indexOf('@')
  const local = atIdx > 0 ? clean.slice(0, atIdx) : clean
  return local.slice(0, 2).toUpperCase()
}

export default function SettingsModal({
  open,
  onClose,
  userId,
  onPrefsChange,
  theme = 'dark',
  onThemeChange,
  lang: currentLang = 'en',
  onLangChange,
}: SettingsModalProps) {
  const { t } = useI18n()
  const [prefs, setPrefs] = useState<NotifPrefs>(loadNotifPrefs)
  const [saved, setSaved] = useState(false)

  useEffect(() => {
    if (open) {
      setPrefs(loadNotifPrefs())
      setSaved(false)
    }
  }, [open])

  useEffect(() => {
    if (!open) return
    const handler = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [open, onClose])

  const handleSave = () => {
    saveNotifPrefs(prefs)
    onPrefsChange(prefs)
    setSaved(true)
    setTimeout(() => setSaved(false), 2500)
  }

  const setOnePref = (key: keyof NotifPrefs, val: boolean) => {
    setPrefs(p => ({ ...p, [key]: val }))
    setSaved(false)
  }

  if (!open) return null

  const notifItems: { key: keyof NotifPrefs; label: string; color: string }[] = [
    { key: 'running',   label: t.settings.notifRunning,   color: 'var(--accent)' },
    { key: 'success',   label: t.settings.notifSuccess,   color: 'var(--success)' },
    { key: 'failed',    label: t.settings.notifFailed,    color: 'var(--danger)' },
    { key: 'cancelled', label: t.settings.notifCancelled, color: 'var(--text-disabled)' },
  ]

  return (
    <div css={backdrop} onMouseDown={(e) => e.target === e.currentTarget && onClose()}>
      <div css={dialog} role="dialog" aria-modal aria-label={t.settings.title}>
        <div css={modalHeader}>
          <div css={modalTitle}>
            <svg width="15" height="15" viewBox="0 0 16 16" fill="none" aria-hidden>
              <circle cx="8" cy="8" r="2.2" stroke="currentColor" strokeWidth="1.2" />
              <path d="M8 1.5v1.3M8 13.2v1.3M1.5 8h1.3M13.2 8h1.3M3.5 3.5l.9.9M11.6 11.6l.9.9M3.5 12.5l.9-.9M11.6 4.4l.9-.9" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" />
            </svg>
            {t.settings.title}
          </div>
          <button type="button" css={closeBtn} onClick={onClose} aria-label={t.settings.close}>
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
              <path d="M2 2l8 8M10 2l-8 8" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
            </svg>
          </button>
        </div>

        <div css={body}>

          <div css={section}>
            <span css={sectionLabel}>{t.settings.account}</span>
            <div css={accountCard}>
              <div css={avatar}>{getInitials(userId)}</div>
              <div css={accountInfo}>
                <div css={accountEmail}>{userId ?? '—'}</div>
                <div css={accountRole}>
                  <span css={statusDot} />
                  {t.settings.statusActive}
                </div>
              </div>
            </div>
          </div>

          <div css={divider} />


          <div css={section}>
            <span css={sectionLabel}>{t.prefs.appearance}</span>
            <div css={prefRow}>
              <span css={prefRowLabel}>{t.prefs.appearance}</span>
              <div css={segmentRow}>
                <button
                  type="button"
                  css={segmentBtn(theme === 'dark')}
                  onClick={() => onThemeChange?.('dark')}
                >
                  {t.prefs.dark}
                </button>
                <button
                  type="button"
                  css={segmentBtn(theme === 'light')}
                  onClick={() => onThemeChange?.('light')}
                >
                  {t.prefs.light}
                </button>
              </div>
            </div>
            <div css={prefRow}>
              <span css={prefRowLabel}>{t.prefs.language}</span>
              <div css={segmentRow}>
                <button
                  type="button"
                  css={segmentBtn(currentLang === 'en')}
                  onClick={() => onLangChange?.('en')}
                >
                  {t.prefs.langEn}
                </button>
                <button
                  type="button"
                  css={segmentBtn(currentLang === 'ru')}
                  onClick={() => onLangChange?.('ru')}
                >
                  {t.prefs.langRu}
                </button>
              </div>
            </div>
          </div>

          <div css={divider} />


          <div css={section}>
            <span css={sectionLabel}>{t.settings.notifications}</span>
            {notifItems.map(({ key, label, color }) => (
              <div key={key} css={toggleRow}>
                <span css={notifColorDot(color)} />
                <div css={toggleInfo}>
                  <div css={toggleLabel}>{label}</div>
                </div>
                <Toggle checked={prefs[key]} onChange={v => setOnePref(key, v)} />
              </div>
            ))}
          </div>
        </div>

        <div css={footer}>
          {saved && (
            <span css={savedPill}>
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                <path d="M2 6.5l3 3 5-6" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
              {t.common.save}
            </span>
          )}
          <button type="button" css={btnSecondary} onClick={onClose}>{t.settings.cancel}</button>
          <button type="button" css={btnPrimary} onClick={handleSave}>{t.settings.save}</button>
        </div>
      </div>
    </div>
  )
}
