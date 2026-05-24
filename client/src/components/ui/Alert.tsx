import { css } from '@emotion/react'
import { alertError, alertSuccess } from '../../styles/theme'

type Variant = 'error' | 'success'

const iconWrap = css({
  flexShrink: 0,
  marginTop: 1,
})

function IconError() {
  return (
    <svg width="14" height="14" viewBox="0 0 16 16" fill="none" aria-hidden>
      <circle cx="8" cy="8" r="6.5" stroke="currentColor" strokeWidth="1.3" />
      <path d="M8 5v3.5M8 10.5v.5" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
    </svg>
  )
}

function IconSuccess() {
  return (
    <svg width="14" height="14" viewBox="0 0 16 16" fill="none" aria-hidden>
      <circle cx="8" cy="8" r="6.5" stroke="currentColor" strokeWidth="1.3" />
      <path d="M5.5 8.5l2 2 3-4" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export default function Alert({
  variant,
  children,
  ...props
}: { variant: Variant; children: React.ReactNode } & React.HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      css={variant === 'error' ? alertError : alertSuccess}
      role="alert"
      aria-live="polite"
      {...props}
    >
      <span css={iconWrap}>
        {variant === 'error' ? <IconError /> : <IconSuccess />}
      </span>
      <span>{children}</span>
    </div>
  )
}
