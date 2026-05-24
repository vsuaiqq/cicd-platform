import { useEffect } from 'react'
import { Link } from 'react-router-dom'
import { css } from '@emotion/react'
import { useAppSelector } from '../store'
import { useI18n } from '../i18n'
import { pageTitle, pageDesc, btnSecondary } from '../styles/theme'

const pageWrap = css({
  animation: 'fade-in 0.3s var(--ease-out) both',
  maxWidth: 560,
})

const heroCard = css({
  display: 'flex',
  alignItems: 'center',
  gap: 20,
  padding: '24px 28px',
  borderRadius: 'var(--radius-lg)',
  border: '1px solid var(--border-subtle)',
  background: 'var(--bg-elevated)',
  marginBottom: 24,
})

const avatar = css({
  width: 64,
  height: 64,
  borderRadius: 'var(--radius-full)',
  background: 'linear-gradient(135deg, var(--indigo), var(--accent))',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  fontSize: '1.25rem',
  fontWeight: 700,
  color: 'var(--text-on-accent)',
  flexShrink: 0,
  boxShadow: '0 0 0 1px var(--avatar-inset) inset',
})

const heroText = css({ minWidth: 0 })

const displayNameStyle = css({
  fontSize: '1.375rem',
  fontWeight: 700,
  letterSpacing: '-0.02em',
  color: 'var(--text-primary)',
  lineHeight: 1.25,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

const heroSub = css({
  fontSize: '0.875rem',
  color: 'var(--text-tertiary)',
  marginTop: 4,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

const fieldList = css({
  display: 'flex',
  flexDirection: 'column',
  gap: 0,
  borderRadius: 'var(--radius-lg)',
  border: '1px solid var(--border-subtle)',
  background: 'var(--bg-elevated)',
  overflow: 'hidden',
})

const fieldRow = css({
  display: 'flex',
  flexDirection: 'column',
  gap: 4,
  padding: '14px 20px',
  borderBottom: '1px solid var(--border-subtle)',
  '&:last-of-type': { borderBottom: 'none' },
})

const fieldLabel = css({
  fontSize: '0.6875rem',
  fontWeight: 600,
  textTransform: 'uppercase',
  letterSpacing: '0.08em',
  color: 'var(--text-disabled)',
})

const fieldValue = css({
  fontSize: '0.9375rem',
  fontWeight: 500,
  color: 'var(--text-primary)',
  wordBreak: 'break-all',
})

const actions = css({
  marginTop: 24,
  display: 'flex',
  gap: 10,
})

function getInitials(label: string): string {
  const parts = label.trim().split(/\s+/).filter(Boolean)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  if (parts.length === 1 && parts[0].length >= 2) return parts[0].slice(0, 2).toUpperCase()
  return label.slice(0, 2).toUpperCase() || '?'
}

export default function Profile() {
  const { userId, email, username } = useAppSelector((s) => s.auth)
  const { t } = useI18n()

  const displayName = username?.trim() || email?.trim() || userId || '—'
  const initials = getInitials(displayName)

  useEffect(() => {
    document.title = `${t.profilePage.title} — Flow`
  }, [t.profilePage.title])

  return (
    <div css={pageWrap}>
      <header css={{ marginBottom: 28 }}>
        <h1 css={pageTitle}>{t.profilePage.title}</h1>
        <p css={pageDesc}>{t.profilePage.subtitle}</p>
      </header>

      <div css={heroCard}>
        <div css={avatar} aria-hidden>{initials}</div>
        <div css={heroText}>
          <div css={displayNameStyle}>{displayName}</div>
          {email && <div css={heroSub}>{email}</div>}
        </div>
      </div>

      <div css={fieldList}>
        <div css={fieldRow}>
          <span css={fieldLabel}>{t.profilePage.username}</span>
          <span css={fieldValue}>{username?.trim() || '—'}</span>
        </div>
        <div css={fieldRow}>
          <span css={fieldLabel}>{t.profilePage.email}</span>
          <span css={fieldValue}>{email?.trim() || '—'}</span>
        </div>
        <div css={fieldRow}>
          <span css={fieldLabel}>{t.profilePage.userId}</span>
          <span css={fieldValue}>{userId || '—'}</span>
        </div>
      </div>

      <div css={actions}>
        <Link to="/" css={btnSecondary} style={{ textDecoration: 'none' }}>
          {t.nav.dashboard}
        </Link>
      </div>
    </div>
  )
}
