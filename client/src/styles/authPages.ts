

import { css } from '@emotion/react'
import { btnPrimary, form, formInput, formLabel, pageDesc, pageTitle } from './theme'

export { form, formLabel }


export const formInputAuth = formInput

export const btnAuthPrimary = css([
  btnPrimary,
  css({
    width: '100%',
    marginTop: 4,
  }),
])


export const authEyebrow = css({
  fontSize: '0.6875rem',
  fontWeight: 600,
  color: 'var(--text-disabled)',
  textTransform: 'uppercase',
  letterSpacing: '0.09em',
  marginBottom: 12,
  display: 'block',
})


export const authPageTitle = pageTitle


export const authPageDesc = css([
  pageDesc,
  css({
    marginBottom: 28,
  }),
])

export const authPwdWrap = css({ position: 'relative', width: '100%' })

export const authPwdToggle = css({
  position: 'absolute',
  right: 10,
  top: '50%',
  transform: 'translateY(-50%)',
  background: 'transparent',
  border: 'none',
  padding: '4px 6px',
  cursor: 'pointer',
  color: 'var(--text-disabled)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  borderRadius: 'var(--radius-md)',
  transition: 'color var(--duration-fast) var(--ease-out), background var(--duration-fast) var(--ease-out)',
  '&:hover': {
    color: 'var(--text-secondary)',
    background: 'var(--bg-hover)',
  },
})

export const authFooter = css({
  marginTop: 28,
  paddingTop: 24,
  borderTop: '1px solid var(--border-subtle)',
  textAlign: 'center',
  fontSize: '0.875rem',
  color: 'var(--text-tertiary)',
  lineHeight: 1.5,
})

export const authFooterLink = css({
  color: 'var(--accent)',
  fontWeight: 500,
  marginLeft: 4,
  textDecoration: 'none',
  transition: 'color var(--duration-fast) var(--ease-out)',
  '&:hover': { color: 'var(--accent-hover)' },
})

export const hintText = css({
  fontSize: '0.75rem',
  color: 'var(--text-tertiary)',
  marginTop: 4,
  lineHeight: 1.5,
})
