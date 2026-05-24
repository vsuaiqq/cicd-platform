import { Outlet } from 'react-router-dom'
import { css } from '@emotion/react'
import AuthHeroPanel from './auth/AuthHeroPanel'
import ThemeToggle from './ThemeToggle'
import { card } from '../styles/theme'

const shell = css({
  minHeight: '100vh',
  display: 'grid',
  gridTemplateColumns: 'minmax(0, 1fr) minmax(320px, 1fr)',
  '@media (max-width: 900px)': {
    gridTemplateColumns: '1fr',
  },
})

const formColumn = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: 'clamp(20px, 4vw, 48px)',
  position: 'relative',
  background: 'var(--bg-base)',
})

const themeTogglePos = css({
  position: 'absolute',
  top: 16,
  right: 16,
  zIndex: 2,
})

const formCard = css([
  card,
  css({
    width: '100%',
    maxWidth: 420,
    padding: 'clamp(24px, 4vw, 36px) clamp(22px, 4vw, 32px)',
    marginBottom: 0,
    boxShadow: 'var(--shadow-sm)',
    animation: 'fade-in 0.3s var(--ease-out) both',
    position: 'relative',
    zIndex: 1,
    '&:hover': {
      borderColor: 'var(--border-subtle)',
    },
  }),
])

const formInner = css({
  position: 'relative',
  zIndex: 1,
})

export default function AuthLayout() {
  return (
    <div css={shell}>
      <AuthHeroPanel />
      <div css={formColumn}>
        <div css={themeTogglePos}>
          <ThemeToggle />
        </div>
        <div css={formCard}>
          <div css={formInner}>
            <Outlet />
          </div>
        </div>
      </div>
    </div>
  )
}
