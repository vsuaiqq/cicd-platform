import { useEffect, type ReactNode } from 'react'
import { css } from '@emotion/react'



const backdrop = css({
  position: 'fixed',
  inset: 0,
  background: 'var(--overlay-backdrop)',
  backdropFilter: 'blur(3px)',
  zIndex: 9000,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: 24,
  animation: 'modal-backdrop 0.18s ease both',
})

const dialog = css({
  width: '100%',
  maxWidth: 400,
  background: 'var(--bg-card)',
  border: '1px solid var(--border-default)',
  borderRadius: 'var(--radius-xl)',
  boxShadow: 'var(--shadow-lg), 0 0 0 1px var(--overlay-inset) inset',
  animation: 'modal-slide 0.2s var(--ease-out) both',
  overflow: 'hidden',
  position: 'relative',
  '&::before': {
    content: '""',
    position: 'absolute',
    top: 0,
    left: 0,
    right: 0,
    height: 1,
    background: 'linear-gradient(90deg, transparent 10%, var(--border-default) 50%, transparent 90%)',
  },
})

const dialogBody = css({
  padding: '28px 28px 24px',
})

const dialogIcon = css({
  width: 40,
  height: 40,
  borderRadius: 10,
  background: 'var(--danger-muted)',
  border: '1px solid var(--danger-glow)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: 'var(--danger)',
  marginBottom: 16,
})

const dialogTitle = css({
  fontSize: '1.0625rem',
  fontWeight: 700,
  letterSpacing: '-0.02em',
  color: 'var(--text-primary)',
  marginBottom: 8,
  lineHeight: 1.3,
})

const dialogDesc = css({
  fontSize: '0.875rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.6,
})

const dialogFooter = css({
  display: 'flex',
  gap: 8,
  justifyContent: 'flex-end',
  padding: '16px 20px',
  background: 'var(--bg-overlay)',
  borderTop: '1px solid var(--border-subtle)',
})

const btnCancel = css({
  padding: '8px 16px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'transparent',
  color: 'var(--text-secondary)',
  fontSize: '0.875rem',
  fontWeight: 500,
  cursor: 'pointer',
  transition: 'border-color 150ms, background 150ms',
  '&:hover': {
    borderColor: 'var(--border-strong)',
    background: 'var(--bg-hover)',
  },
})

const btnConfirm = css({
  padding: '8px 18px',
  borderRadius: 'var(--radius-md)',
  border: 'none',
  background: 'var(--danger)',
  color: 'var(--text-on-accent)',
  fontSize: '0.875rem',
  fontWeight: 600,
  cursor: 'pointer',
  transition: 'opacity 150ms, transform 150ms',
  '&:hover': { opacity: 0.88, transform: 'translateY(-1px)' },
  '&:active': { transform: 'translateY(0)', opacity: 1 },
  '&:disabled': { opacity: 0.5, cursor: 'not-allowed' },
})



interface ConfirmDialogProps {
  open: boolean
  title: string
  description: string | ReactNode
  confirmLabel?: string
  cancelLabel?: string
  loading?: boolean
  onConfirm: () => void
  onCancel: () => void
}

function DeleteIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 20 20" fill="none" aria-hidden>
      <path d="M4 6h12M8 6V4h4v2M6 6l.75 11h6.5L14 6" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export default function ConfirmDialog({
  open,
  title,
  description,
  confirmLabel = 'Delete',
  cancelLabel = 'Cancel',
  loading = false,
  onConfirm,
  onCancel,
}: ConfirmDialogProps) {
  useEffect(() => {
    if (!open) return
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onCancel()
      if (e.key === 'Enter') onConfirm()
    }
    window.addEventListener('keydown', handleKey)
    return () => window.removeEventListener('keydown', handleKey)
  }, [open, onCancel, onConfirm])

  if (!open) return null

  return (
    <div css={backdrop} onClick={(e) => e.target === e.currentTarget && onCancel()}>
      <div
        css={dialog}
        role="alertdialog"
        aria-modal="true"
        aria-labelledby="confirm-title"
        aria-describedby="confirm-desc"
      >
        <div css={dialogBody}>
          <div css={dialogIcon}>
            <DeleteIcon />
          </div>
          <div id="confirm-title" css={dialogTitle}>{title}</div>
          <p id="confirm-desc" css={dialogDesc}>{description}</p>
        </div>
        <div css={dialogFooter}>
          <button type="button" css={btnCancel} onClick={onCancel} disabled={loading}>
            {cancelLabel}
          </button>
          <button type="button" css={btnConfirm} onClick={onConfirm} disabled={loading}>
            {loading ? `${confirmLabel}…` : confirmLabel}
          </button>
        </div>
      </div>
    </div>
  )
}
