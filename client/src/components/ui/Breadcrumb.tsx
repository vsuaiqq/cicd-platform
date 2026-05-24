import { Link } from 'react-router-dom'
import { css } from '@emotion/react'
import { breadcrumb } from '../../styles/theme'
import { useI18n } from '../../i18n'

export type BreadcrumbItem = { label: string; to?: string }

const sep = css({
  color: 'var(--text-disabled)',
  userSelect: 'none',
  lineHeight: 1,
})

export default function Breadcrumb({ items }: { items: BreadcrumbItem[] }) {
  const { t } = useI18n()
  return (
    <nav css={breadcrumb} aria-label={t.aria.breadcrumb}>
      {items.map((item, i) => (
        <span key={i} style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}>
          {i > 0 && (
            <span css={sep} aria-hidden>
              <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                <path d="M4.5 2.5L7.5 6l-3 3.5" stroke="currentColor" strokeWidth="1.1" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
            </span>
          )}
          {item.to ? (
            <Link to={item.to}>{item.label}</Link>
          ) : (
            <span style={{ color: 'var(--text-secondary)', fontWeight: 500 }}>{item.label}</span>
          )}
        </span>
      ))}
    </nav>
  )
}
