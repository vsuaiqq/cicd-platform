import { useState } from 'react'
import { css } from '@emotion/react'
import { useI18n } from '../../i18n'

const terminal = css({
  background: 'var(--code-bg)',
  border: '1px solid var(--border-default)',
  borderRadius: 'var(--radius-md)',
  overflow: 'hidden',
  marginBottom: 10,
})

const terminalHeader = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: '8px 14px',
  borderBottom: '1px solid var(--border-subtle)',
  background: 'var(--bg-hover)',
})

const terminalDots = css({
  display: 'flex',
  gap: 6,
})

const terminalDot = css({
  width: 9,
  height: 9,
  borderRadius: '50%',
  background: 'var(--border-default)',
})

const copyBtn = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 5,
  padding: '3px 10px',
  borderRadius: 6,
  border: '1px solid var(--border-default)',
  background: 'transparent',
  color: 'var(--text-tertiary)',
  fontSize: '0.6875rem',
  fontWeight: 600,
  letterSpacing: '0.04em',
  textTransform: 'uppercase',
  cursor: 'pointer',
  transition: 'color 150ms, background 150ms, border-color 150ms',
  '&:hover': {
    color: 'var(--text-primary)',
    background: 'var(--bg-hover)',
    borderColor: 'var(--border-strong)',
  },
})

const copyBtnSuccess = css({
  color: 'var(--success)',
  borderColor: 'var(--success-glow)',
  background: 'var(--success-muted)',
  '&:hover': {
    color: 'var(--success)',
    background: 'var(--success-muted)',
  },
})

const pre = css({
  margin: 0,
  padding: '14px 16px',
  fontFamily: 'ui-monospace, "SF Mono", monospace',
  fontSize: '0.8125rem',
  lineHeight: 1.7,
  color: 'var(--code-text)',
  wordBreak: 'break-all',
  overflowX: 'auto',
  whiteSpace: 'pre-wrap',
})

function IconCopy() {
  return (
    <svg width="11" height="11" viewBox="0 0 14 14" fill="none" aria-hidden>
      <rect x="4" y="4" width="8" height="8" rx="1.5" stroke="currentColor" strokeWidth="1.2" />
      <path d="M2 10V2h8" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function IconCheck() {
  return (
    <svg width="11" height="11" viewBox="0 0 14 14" fill="none" aria-hidden>
      <path d="M2.5 7.5l3 3 6-6" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export default function CodeBlock({
  value,
  label,
}: {
  value: string
  label?: string
}) {
  const { t } = useI18n()
  const copyLabel = label ?? t.common.copy
  const [copied, setCopied] = useState(false)

  const copy = async () => {
    try {
      await navigator.clipboard.writeText(value)
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    } catch {

    }
  }

  return (
    <div css={terminal}>
      <div css={terminalHeader}>
        <div css={terminalDots}>
          <span css={terminalDot} />
          <span css={terminalDot} />
          <span css={terminalDot} />
        </div>
        <button
          type="button"
          css={[copyBtn, copied && copyBtnSuccess]}
          onClick={copy}
          aria-label={t.aria.copyToClipboard}
        >
          {copied ? <IconCheck /> : <IconCopy />}
          {copied ? t.common.copied : copyLabel}
        </button>
      </div>
      <pre css={pre}>{value}</pre>
    </div>
  )
}
