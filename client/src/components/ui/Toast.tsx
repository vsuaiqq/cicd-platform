import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from 'react'
import { css } from '@emotion/react'
import { useI18n } from '../../i18n'



export type ToastVariant = 'success' | 'error' | 'warning' | 'info'

export interface Toast {
  id: string
  message: string
  variant: ToastVariant
  duration?: number
}

interface ToastContextValue {
  addToast: (message: string, variant?: ToastVariant, duration?: number) => void
  success: (message: string, duration?: number) => void
  error: (message: string, duration?: number) => void
  warning: (message: string, duration?: number) => void
  info: (message: string, duration?: number) => void
}



const ToastContext = createContext<ToastContextValue | null>(null)

export function useToast(): ToastContextValue {
  const ctx = useContext(ToastContext)
  if (!ctx) throw new Error('useToast must be used inside <ToastProvider>')
  return ctx
}



const container = css({
  position: 'fixed',
  bottom: 24,
  right: 24,
  zIndex: 9999,
  display: 'flex',
  flexDirection: 'column',
  gap: 8,
  maxWidth: 380,
  pointerEvents: 'none',
})

const toastBase = css({
  display: 'flex',
  alignItems: 'flex-start',
  gap: 10,
  padding: '12px 14px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid',
  backdropFilter: 'blur(8px)',
  boxShadow: 'var(--shadow-md), 0 0 0 1px var(--overlay-inset) inset',
  animation: 'toast-in 0.22s var(--ease-out) both',
  pointerEvents: 'all',
  minWidth: 260,
  cursor: 'default',
  transition: 'opacity 0.2s, transform 0.2s',
})

const toastVariants: Record<ToastVariant, ReturnType<typeof css>> = {
  success: css({
    background: 'var(--toast-success-bg)',
    borderColor: 'var(--success-glow)',
    color: 'var(--success)',
  }),
  error: css({
    background: 'var(--toast-error-bg)',
    borderColor: 'var(--danger-glow)',
    color: 'var(--danger)',
  }),
  warning: css({
    background: 'var(--toast-warning-bg)',
    borderColor: 'var(--warning-glow)',
    color: 'var(--warning)',
  }),
  info: css({
    background: 'var(--toast-info-bg)',
    borderColor: 'var(--accent-glow)',
    color: 'var(--accent)',
  }),
}

const toastMessage = css({
  flex: 1,
  fontSize: '0.875rem',
  fontWeight: 500,
  lineHeight: 1.45,
  color: 'var(--text-primary)',
})

const dismissBtn = css({
  flexShrink: 0,
  background: 'none',
  border: 'none',
  cursor: 'pointer',
  color: 'var(--text-disabled)',
  padding: 2,
  lineHeight: 1,
  transition: 'color 120ms',
  borderRadius: 4,
  marginTop: 1,
  '&:hover': { color: 'var(--text-secondary)' },
})

const progressTrack = css({
  position: 'absolute',
  bottom: 0,
  left: 0,
  right: 0,
  height: 2,
  background: 'var(--progress-track)',
  borderRadius: '0 0 var(--radius-md) var(--radius-md)',
  overflow: 'hidden',
})



function ToastIcon({ variant }: { variant: ToastVariant }) {
  if (variant === 'success') return (
    <svg width="14" height="14" viewBox="0 0 14 14" fill="none" style={{ flexShrink: 0, marginTop: 2 }}>
      <circle cx="7" cy="7" r="6" stroke="currentColor" strokeWidth="1.2" />
      <path d="M4 7l2.5 2.5 3.5-5" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
  if (variant === 'error') return (
    <svg width="14" height="14" viewBox="0 0 14 14" fill="none" style={{ flexShrink: 0, marginTop: 2 }}>
      <circle cx="7" cy="7" r="6" stroke="currentColor" strokeWidth="1.2" />
      <path d="M4.5 4.5l5 5M9.5 4.5l-5 5" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  )
  if (variant === 'warning') return (
    <svg width="14" height="14" viewBox="0 0 14 14" fill="none" style={{ flexShrink: 0, marginTop: 2 }}>
      <path d="M7 1L13 12H1L7 1z" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M7 5v3M7 10v.5" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  )
  return (
    <svg width="14" height="14" viewBox="0 0 14 14" fill="none" style={{ flexShrink: 0, marginTop: 2 }}>
      <circle cx="7" cy="7" r="6" stroke="currentColor" strokeWidth="1.2" />
      <path d="M7 6v4M7 4.5v-.5" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  )
}



function ProgressBar({ duration, variant }: { duration: number; variant: ToastVariant }) {
  const color = variant === 'success' ? 'var(--success)' : variant === 'error' ? 'var(--danger)' : variant === 'warning' ? 'var(--warning)' : 'var(--accent)'
  return (
    <div css={progressTrack}>
      <div
        style={{
          height: '100%',
          background: color,
          opacity: 0.45,
          animation: `toast-progress ${duration}ms linear forwards`,
          transformOrigin: 'left',
        }}
      />
    </div>
  )
}



function ToastItem({ toast, onDismiss }: { toast: Toast; onDismiss: (id: string) => void }) {
  const { t } = useI18n()
  const duration = toast.duration ?? 4000

  useEffect(() => {
    const timer = setTimeout(() => onDismiss(toast.id), duration)
    return () => clearTimeout(timer)
  }, [toast.id, duration, onDismiss])

  return (
    <div css={[toastBase, toastVariants[toast.variant]]} role="alert" aria-live="assertive" style={{ position: 'relative', overflow: 'hidden' }}>
      <ToastIcon variant={toast.variant} />
      <span css={toastMessage}>{toast.message}</span>
      <button
        css={dismissBtn}
        type="button"
        onClick={() => onDismiss(toast.id)}
        aria-label={t.aria.dismissNotification}
      >
        <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
          <path d="M2 2l8 8M10 2l-8 8" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
        </svg>
      </button>
      <ProgressBar duration={duration} variant={toast.variant} />
    </div>
  )
}



export function ToastProvider({ children }: { children: ReactNode }) {
  const { t } = useI18n()
  const [toasts, setToasts] = useState<Toast[]>([])
  const counterRef = useRef(0)

  const addToast = useCallback((message: string, variant: ToastVariant = 'info', duration = 4000) => {
    const id = `toast-${++counterRef.current}`
    setToasts((prev) => [...prev, { id, message, variant, duration }])
  }, [])

  const dismiss = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id))
  }, [])

  const value: ToastContextValue = {
    addToast,
    success: (msg, dur) => addToast(msg, 'success', dur),
    error: (msg, dur) => addToast(msg, 'error', dur),
    warning: (msg, dur) => addToast(msg, 'warning', dur),
    info: (msg, dur) => addToast(msg, 'info', dur),
  }

  return (
    <ToastContext.Provider value={value}>
      {children}
      <div css={container} aria-label={t.aria.notifications}>
        {toasts.map((t) => (
          <ToastItem key={t.id} toast={t} onDismiss={dismiss} />
        ))}
      </div>
    </ToastContext.Provider>
  )
}
