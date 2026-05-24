import { css } from '@emotion/react'
import { FlowBrandLogo } from '../FlowBrand'
import { chip, pageDesc, pageTitle } from '../../styles/theme'
import { useI18n } from '../../i18n'

const root = css({
  position: 'relative',
  height: '100%',
  minHeight: 'min(100vh, 720px)',
  display: 'flex',
  flexDirection: 'column',
  justifyContent: 'space-between',
  padding: 'clamp(24px, 4vw, 48px)',
  overflow: 'hidden',
  backgroundColor: 'var(--bg-elevated)',
  borderRight: '1px solid var(--border-subtle)',
  backgroundImage: `
    radial-gradient(ellipse 85% 55% at 50% -5%, var(--accent-muted), transparent 58%),
    radial-gradient(ellipse 50% 40% at 100% 100%, var(--ai-muted), transparent 55%)
  `,
  '@media (max-width: 900px)': {
    minHeight: 'auto',
    borderRight: 'none',
    borderBottom: '1px solid var(--border-subtle)',
    padding: '24px 20px 28px',
  },
})

const brandRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  position: 'relative',
  zIndex: 1,
})

const headlineBlock = css({
  position: 'relative',
  zIndex: 1,
  marginTop: 'clamp(28px, 6vh, 56px)',
  maxWidth: 400,
  '@media (max-width: 900px)': {
    marginTop: 24,
    maxWidth: 'none',
  },
})

const heroTitle = css([
  pageTitle,
  css({
    fontSize: 'clamp(1.5rem, 2.2vw, 1.75rem)',
    marginBottom: 8,
  }),
])

const heroSub = css([
  pageDesc,
  css({
    marginBottom: 0,
  }),
])

const features = css({
  display: 'flex',
  flexDirection: 'column',
  gap: 12,
  marginTop: 28,
  '@media (max-width: 900px)': {
    marginTop: 20,
  },
})

const featureRow = css({
  display: 'flex',
  alignItems: 'flex-start',
  gap: 12,
})

const featureIcon = css({
  flexShrink: 0,
  width: 32,
  height: 32,
  borderRadius: 'var(--radius-md)',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: 'var(--bg-overlay)',
  border: '1px solid var(--border-default)',
  color: 'var(--accent)',
})

const featureText = css({
  fontSize: '0.875rem',
  fontWeight: 600,
  letterSpacing: '-0.01em',
  color: 'var(--text-primary)',
  lineHeight: 1.35,
})

const featureHint = css({
  fontSize: '0.8125rem',
  color: 'var(--text-tertiary)',
  marginTop: 2,
  lineHeight: 1.5,
})

const pipeline = css({
  position: 'relative',
  zIndex: 1,
  marginTop: 'auto',
  paddingTop: 32,
  '@media (max-width: 900px)': {
    paddingTop: 20,
    display: 'none',
  },
})

const pipelineInner = css({
  display: 'flex',
  alignItems: 'center',
  flexWrap: 'wrap',
  gap: 6,
  fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace',
  fontSize: '0.75rem',
  color: 'var(--text-tertiary)',
})

const pipeArrow = css({ color: 'var(--text-disabled)', userSelect: 'none', fontSize: '0.6875rem' })

function IconGitBranch() {
  return (
    <svg width="15" height="15" viewBox="0 0 16 16" fill="none" aria-hidden>
      <circle cx="4" cy="4" r="2" stroke="currentColor" strokeWidth="1.3" />
      <circle cx="4" cy="12" r="2" stroke="currentColor" strokeWidth="1.3" />
      <circle cx="12" cy="4" r="2" stroke="currentColor" strokeWidth="1.3" />
      <path d="M4 6v4M4 6c0 2 8 3 8-2" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  )
}

function IconRocket() {
  return (
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden>
      <path d="M4.5 16.5c-1-1-1.5-4.5-1.5-4.5s3.5-.5 4.5-1.5 4-6 4-6-4 2.5-5 4-1.5 4.5-1.5 4.5z" stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round" />
      <path d="M12 15l-3 3" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <circle cx="12" cy="11" r="1" fill="currentColor" />
    </svg>
  )
}

function IconActivity() {
  return (
    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" aria-hidden>
      <path d="M3 3v18h18" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" />
      <path d="M7 14l3-3 3 2 4-5" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

export default function AuthHeroPanel() {
  const { t } = useI18n()

  return (
    <div css={root}>
      <div>
        <div css={brandRow}>
          <FlowBrandLogo />
        </div>

        <div css={headlineBlock}>
          <h2 css={heroTitle}>{t.auth.heroHeadline}</h2>
          <p css={heroSub}>{t.auth.heroSub}</p>

          <div css={features}>
            <div css={featureRow}>
              <div css={featureIcon}>
                <IconGitBranch />
              </div>
              <div>
                <div css={featureText}>{t.auth.heroFeature1}</div>
                <div css={featureHint}>{t.auth.heroFeature1Hint}</div>
              </div>
            </div>
            <div css={featureRow}>
              <div css={featureIcon}>
                <IconRocket />
              </div>
              <div>
                <div css={featureText}>{t.auth.heroFeature2}</div>
                <div css={featureHint}>{t.auth.heroFeature2Hint}</div>
              </div>
            </div>
            <div css={featureRow}>
              <div css={featureIcon}>
                <IconActivity />
              </div>
              <div>
                <div css={featureText}>{t.auth.heroFeature3}</div>
                <div css={featureHint}>{t.auth.heroFeature3Hint}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div css={pipeline}>
        <div css={pipelineInner}>
          <span css={chip}>checkout</span>
          <span css={pipeArrow}>→</span>
          <span css={chip}>build</span>
          <span css={pipeArrow}>→</span>
          <span css={chip}>test</span>
          <span css={pipeArrow}>→</span>
          <span css={chip}>deploy</span>
        </div>
      </div>
    </div>
  )
}
