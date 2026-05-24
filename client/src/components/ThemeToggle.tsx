import { css } from '@emotion/react'
import { useI18n } from '../i18n'
import { useTheme } from '../lib/themeContext'

const btn = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  padding: '6px 10px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-card)',
  color: 'var(--text-secondary)',
  fontSize: '0.75rem',
  fontWeight: 600,
  cursor: 'pointer',
  transition: 'background 120ms, color 120ms, border-color 120ms',
  '&:hover': {
    background: 'var(--bg-hover)',
    color: 'var(--text-primary)',
    borderColor: 'var(--border-strong)',
  },
})

function IconSun() {
  return (
    <svg width="13" height="13" viewBox="0 0 16 16" fill="none" aria-hidden>
      <circle cx="8" cy="8" r="3" stroke="currentColor" strokeWidth="1.3" />
      <path d="M8 1v2M8 13v2M1 8h2M13 8h2M3.2 3.2l1.4 1.4M11.4 11.4l1.4 1.4M3.2 12.8l1.4-1.4M11.4 4.6l1.4-1.4" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  )
}

function IconMoon() {
  return (
    <svg width="13" height="13" viewBox="0 0 16 16" fill="none" aria-hidden>
      <path d="M13.5 10A6 6 0 016 2.5a6.002 6.002 0 000 11 6 6 0 007.5-3.5z" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

type ThemeToggleProps = {
  className?: string
  style?: React.CSSProperties
}

export default function ThemeToggle({ className, style }: ThemeToggleProps) {
  const { theme, toggleTheme } = useTheme()
  const { t } = useI18n()
  const isDark = theme === 'dark'

  return (
    <button
      type="button"
      css={btn}
      className={className}
      style={style}
      onClick={toggleTheme}
      title={isDark ? t.prefs.light : t.prefs.dark}
      aria-label={isDark ? t.prefs.light : t.prefs.dark}
    >
      {isDark ? <IconMoon /> : <IconSun />}
      {isDark ? t.prefs.dark : t.prefs.light}
    </button>
  )
}
