import { useState } from 'react'
import { css } from '@emotion/react'
import type { PerformanceGateResult } from '../api/analytics'
import type { PerformanceGateJobConfig } from '../lib/pipelineYaml'
import type { Translations } from '../i18n'

const panel = css({
  borderBottom: '1px solid var(--border-subtle)',
  background: 'linear-gradient(to bottom, var(--bg-overlay), transparent)',
  animation: 'fade-in 0.25s var(--ease-out) both',
})

const accordionHeader = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 12,
  width: '100%',
  padding: '14px 20px',
  background: 'none',
  border: 'none',
  cursor: 'pointer',
  textAlign: 'left',
  userSelect: 'none',
  '&:hover': {
    background: 'var(--bg-hover)',
  },
})

const accordionBody = css({
  padding: '0 20px 18px',
})

const chevron = (open: boolean) =>
  css({
    flexShrink: 0,
    color: 'var(--text-disabled)',
    transition: 'transform 200ms var(--ease-out)',
    transform: open ? 'rotate(180deg)' : 'none',
  })

const titleRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  flexWrap: 'wrap',
})

const title = css({
  display: 'flex',
  alignItems: 'center',
  gap: 8,
  fontSize: '0.9375rem',
  fontWeight: 700,
  color: 'var(--text-primary)',
  letterSpacing: '-0.02em',
})

const verdictBanner = (passed: boolean) =>
  css({
    display: 'inline-flex',
    alignItems: 'center',
    gap: 6,
    padding: '5px 12px',
    borderRadius: 'var(--radius-full)',
    fontSize: '0.75rem',
    fontWeight: 700,
    letterSpacing: '0.04em',
    textTransform: 'uppercase',
    background: passed ? 'var(--success-muted)' : 'var(--danger-muted)',
    color: passed ? 'var(--success)' : 'var(--danger)',
    border: `1px solid ${passed ? 'var(--success-glow)' : 'var(--danger-glow)'}`,
  })

const badge = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 4,
  padding: '3px 9px',
  borderRadius: 'var(--radius-full)',
  fontSize: '0.6875rem',
  fontWeight: 600,
  background: 'var(--accent-muted)',
  color: 'var(--accent)',
  border: '1px solid var(--accent-glow)',
})

const coldBadge = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 4,
  padding: '3px 9px',
  borderRadius: 'var(--radius-full)',
  fontSize: '0.6875rem',
  fontWeight: 600,
  background: 'var(--warning-muted)',
  color: 'var(--warning)',
  border: '1px solid var(--warning-glow)',
})

const summaryBox = (passed: boolean) =>
  css({
    padding: '12px 14px',
    borderRadius: 'var(--radius-md)',
    background: passed ? 'var(--success-muted)' : 'var(--danger-muted)',
    border: `1px solid ${passed ? 'var(--success-glow)' : 'var(--danger-glow)'}`,
    color: 'var(--text-primary)',
    fontSize: '0.875rem',
    lineHeight: 1.6,
    marginBottom: 16,
  })

const metaGrid = css({
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fill, minmax(140px, 1fr))',
  gap: 10,
  marginBottom: 18,
})

const metaTile = css({
  padding: '10px 12px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
})

const metaLabel = css({
  fontSize: '0.6875rem',
  fontWeight: 600,
  color: 'var(--text-disabled)',
  textTransform: 'uppercase',
  letterSpacing: '0.06em',
  marginBottom: 4,
})

const metaValue = css({
  fontSize: '0.8125rem',
  fontWeight: 600,
  color: 'var(--text-secondary)',
  fontFamily: 'ui-monospace, monospace',
})

const metricsGrid = css({
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fill, minmax(280px, 1fr))',
  gap: 12,
})

const metricCard = (passed: boolean) =>
  css({
    padding: '14px 16px',
    borderRadius: 'var(--radius-lg)',
    background: 'var(--bg-card)',
    border: `1px solid ${passed ? 'var(--border-subtle)' : 'var(--danger-glow)'}`,
    boxShadow: passed ? 'none' : '0 0 0 1px var(--danger-muted)',
  })

const metricHeader = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  gap: 8,
  marginBottom: 10,
})

const metricName = css({
  fontSize: '0.8125rem',
  fontWeight: 600,
  color: 'var(--text-primary)',
  fontFamily: 'ui-monospace, monospace',
  wordBreak: 'break-all',
})

const metricStatus = (passed: boolean) =>
  css({
    fontSize: '0.625rem',
    fontWeight: 700,
    padding: '2px 7px',
    borderRadius: 'var(--radius-full)',
    background: passed ? 'var(--success-muted)' : 'var(--danger-muted)',
    color: passed ? 'var(--success)' : 'var(--danger)',
    flexShrink: 0,
  })

const valuesRow = css({
  display: 'grid',
  gridTemplateColumns: '1fr 1fr 1fr',
  gap: 8,
  marginBottom: 10,
})

const valueBlock = css({
  minWidth: 0,
})

const valueLabel = css({
  fontSize: '0.625rem',
  color: 'var(--text-disabled)',
  marginBottom: 2,
  textTransform: 'uppercase',
  letterSpacing: '0.05em',
})

const valueNum = css({
  fontSize: '0.8125rem',
  fontWeight: 600,
  color: 'var(--text-secondary)',
  fontVariantNumeric: 'tabular-nums',
  fontFamily: 'ui-monospace, monospace',
})

const barTrack = css({
  height: 6,
  borderRadius: 'var(--radius-full)',
  background: 'var(--bg-overlay)',
  overflow: 'hidden',
  position: 'relative',
  marginBottom: 8,
})

const barFill = (color: string, widthPct: number) =>
  css({
    height: '100%',
    width: `${Math.min(100, Math.max(0, widthPct))}%`,
    background: color,
    borderRadius: 'var(--radius-full)',
    transition: 'width 0.4s var(--ease-out)',
  })

const thresholdMarker = (pct: number) =>
  css({
    position: 'absolute',
    top: -2,
    left: `${Math.min(100, Math.max(0, pct))}%`,
    width: 2,
    height: 10,
    background: 'var(--warning)',
    borderRadius: 1,
    transform: 'translateX(-50%)',
  })

const reasonToggle = css({
  background: 'none',
  border: 'none',
  padding: 0,
  fontSize: '0.75rem',
  color: 'var(--accent)',
  cursor: 'pointer',
  fontWeight: 500,
  '&:hover': { textDecoration: 'underline' },
})

const reasonText = css({
  marginTop: 6,
  fontSize: '0.75rem',
  color: 'var(--text-tertiary)',
  lineHeight: 1.55,
  fontFamily: 'ui-monospace, monospace',
})

const pendingBox = css({
  padding: '20px',
  textAlign: 'center',
  color: 'var(--text-tertiary)',
  fontSize: '0.875rem',
})

const skeletonBar = css({
  height: 80,
  borderRadius: 'var(--radius-lg)',
  background: 'var(--bg-overlay)',
  animation: 'pulse-subtle 1.5s ease-in-out infinite',
})

function IconChevronDown() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden>
      <path d="M4 6l4 4 4-4" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function IconShield() {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden>
      <path
        d="M12 2l8 4v6c0 5-3.5 9.5-8 10-4.5-.5-8-5-8-10V6l8-4z"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinejoin="round"
      />
      <path d="M9 12l2 2 4-4" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
}

function fmtMetric(value: number, name: string): string {
  if (name.includes('rate') || name.includes('failed')) {
    return `${(value * 100).toFixed(2)}%`
  }
  if (value >= 1000) return value.toFixed(0)
  if (value >= 10) return value.toFixed(1)
  return value.toFixed(3)
}

function directionLabel(dir: string, t: Translations): string {
  if (dir === 'higher_is_better') return t.performanceGate.higherIsBetter
  return t.performanceGate.lowerIsBetter
}

function barScale(current: number, _baseline: number, threshold: number) {
  const max = Math.max(current, threshold, 1)
  return {
    currentPct: (current / max) * 100,
    thresholdPct: (threshold / max) * 100,
  }
}

function checkModeLabel(mode: string | undefined, t: Translations): string {
  switch (mode) {
    case 'both':
      return t.performanceGate.checkModeBoth
    case 'constant':
      return t.performanceGate.checkModeConstant
    case 'adaptive':
      return t.performanceGate.checkModeAdaptive
    case 'cold_start':
      return t.performanceGate.checkModeColdStart
    default:
      return mode ?? ''
  }
}

function MetricCard({
  metric,
  t,
}: {
  metric: PerformanceGateResult['metrics'][number]
  t: Translations
}) {
  const [showReason, setShowReason] = useState(false)
  const hasConstant = metric.constant_threshold != null
  const hasAdaptive = metric.adaptive_threshold != null && !metric.adaptive_skipped
  const compareThreshold = metric.adaptive_threshold ?? metric.constant_threshold ?? metric.threshold
  const { currentPct, thresholdPct } = barScale(
    metric.current,
    metric.baseline_mean,
    compareThreshold,
  )
  const barColor = metric.passed ? 'var(--success)' : 'var(--danger)'

  return (
    <div css={metricCard(metric.passed)}>
      <div css={metricHeader}>
        <span css={metricName}>{metric.name}</span>
        <span css={metricStatus(metric.passed)}>
          {metric.passed ? t.performanceGate.pass : t.performanceGate.fail}
        </span>
      </div>

      {metric.check_mode && (
        <div style={{ fontSize: '0.6875rem', color: 'var(--text-disabled)', marginBottom: 8 }}>
          {checkModeLabel(metric.check_mode, t)}
        </div>
      )}

      <div css={valuesRow}>
        <div css={valueBlock}>
          <div css={valueLabel}>{t.performanceGate.current}</div>
          <div css={valueNum}>{fmtMetric(metric.current, metric.name)}</div>
        </div>
        {hasAdaptive && (
          <div css={valueBlock}>
            <div css={valueLabel}>{t.performanceGate.baseline}</div>
            <div css={valueNum}>{fmtMetric(metric.baseline_mean, metric.name)}</div>
          </div>
        )}
      </div>

      {hasConstant && (
        <div style={{ fontSize: '0.8125rem', marginBottom: 6, color: metric.constant_passed ? 'var(--success)' : 'var(--danger)' }}>
          {metric.constant_passed ? '✓' : '✗'} {t.performanceGate.constantCheck}{' '}
          {metric.direction === 'higher_is_better' ? '≥' : '≤'}{' '}
          {fmtMetric(metric.constant_threshold!, metric.name)}
        </div>
      )}

      {hasAdaptive && (
        <div style={{ fontSize: '0.8125rem', marginBottom: 6, color: metric.adaptive_passed ? 'var(--success)' : 'var(--danger)' }}>
          {metric.adaptive_passed ? '✓' : '✗'} {t.performanceGate.adaptiveCheck}{' '}
          {metric.direction === 'higher_is_better' ? '≥' : '≤'}{' '}
          {fmtMetric(metric.adaptive_threshold!, metric.name)}
        </div>
      )}

      {metric.adaptive_skipped && (
        <p style={{ fontSize: '0.8125rem', color: 'var(--warning)', margin: '0 0 8px' }}>
          {metric.cold_start ? t.performanceGate.metricColdStart : t.performanceGate.adaptiveDisabled}
        </p>
      )}

      {(hasConstant || hasAdaptive) && (
        <>
          <div css={barTrack} title={t.performanceGate.barHint}>
            <div css={barFill(barColor, currentPct)} />
            <div css={thresholdMarker(thresholdPct)} />
          </div>
          <div style={{ fontSize: '0.6875rem', color: 'var(--text-disabled)', marginBottom: 4 }}>
            {directionLabel(metric.direction, t)}
            {hasAdaptive && metric.baseline_stddev > 0 && (
              <span> · σ={fmtMetric(metric.baseline_stddev, metric.name)}</span>
            )}
          </div>
        </>
      )}

      {metric.reason && (
        <>
          <button type="button" css={reasonToggle} onClick={() => setShowReason((v) => !v)}>
            {showReason ? t.performanceGate.hideDetails : t.performanceGate.showDetails}
          </button>
          {showReason && <div css={reasonText}>{metric.reason}</div>}
        </>
      )}
    </div>
  )
}

export interface PerformanceGatePanelProps {
  config: PerformanceGateJobConfig
  result: PerformanceGateResult | undefined
  loading: boolean
  jobStatus: string
  t: Translations
}

export function PerformanceGatePanel({
  config,
  result,
  loading,
  jobStatus,
  t,
}: PerformanceGatePanelProps) {
  const [expanded, setExpanded] = useState(false)
  const isPending = jobStatus === 'running' || jobStatus === 'pending'
  const waiting = isPending || (loading && !result?.found)

  const headerBadges = (
    <>
      {result?.found && (
        <span css={verdictBanner(result.passed)}>
          {result.passed ? t.performanceGate.released : t.performanceGate.blocked}
        </span>
      )}
      {result?.cold_start && (
        <span css={coldBadge}>{t.performanceGate.coldStart}</span>
      )}
      <span css={badge}>{t.performanceGate.adaptiveGate}</span>
      {waiting && (
        <span style={{ fontSize: '0.75rem', color: 'var(--text-tertiary)', fontWeight: 500 }}>
          {t.performanceGate.evaluating}
        </span>
      )}
    </>
  )

  return (
    <div css={panel}>
      <button
        type="button"
        css={accordionHeader}
        onClick={() => setExpanded((open) => !open)}
        aria-expanded={expanded}
      >
        <div css={titleRow}>
          <span css={title}>
            <IconShield />
            {t.performanceGate.title}
          </span>
          {headerBadges}
        </div>
        <span css={chevron(expanded)}>
          <IconChevronDown />
        </span>
      </button>

      {expanded && (
        <div css={accordionBody}>
          {waiting && (
            <div css={pendingBox}>
              <div css={skeletonBar} style={{ marginBottom: 12 }} />
              {t.performanceGate.evaluating}
            </div>
          )}

          {!waiting && !result?.found && (
            <div css={pendingBox}>{t.performanceGate.noResult}</div>
          )}

          {!waiting && result?.found && (
            <>
              {result.evaluated_at && (
                <div style={{ fontSize: '0.75rem', color: 'var(--text-disabled)', marginBottom: 12 }}>
                  {t.performanceGate.evaluatedAt.replace('{time}', new Date(result.evaluated_at).toLocaleString())}
                </div>
              )}

              <div css={summaryBox(result.passed)}>{result.summary}</div>

              <div css={metaGrid}>
                <div css={metaTile}>
                  <div css={metaLabel}>{t.performanceGate.sourceJob}</div>
                  <div css={metaValue}>{config.sourceJob}</div>
                </div>
                <div css={metaTile}>
                  <div css={metaLabel}>{t.performanceGate.baselineSamples}</div>
                  <div css={metaValue}>{result.baseline_samples}</div>
                </div>
                {config.baseline?.min_samples != null && (
                  <div css={metaTile}>
                    <div css={metaLabel}>{t.performanceGate.minSamples}</div>
                    <div css={metaValue}>{config.baseline.min_samples}</div>
                  </div>
                )}
                {config.adaptive?.sigma_factor != null && (
                  <div css={metaTile}>
                    <div css={metaLabel}>{t.performanceGate.sigmaFactor}</div>
                    <div css={metaValue}>{config.adaptive.sigma_factor}σ</div>
                  </div>
                )}
                {config.adaptive?.max_regression_pct != null && (
                  <div css={metaTile}>
                    <div css={metaLabel}>{t.performanceGate.maxRegression}</div>
                    <div css={metaValue}>{config.adaptive.max_regression_pct}%</div>
                  </div>
                )}
                <div css={metaTile}>
                  <div css={metaLabel}>{t.performanceGate.adaptiveMode}</div>
                  <div css={metaValue}>
                    {config.adaptive?.enabled === false ? t.performanceGate.adaptiveOff : t.performanceGate.adaptiveOn}
                  </div>
                </div>
                {config.baseline?.window_days != null && (
                  <div css={metaTile}>
                    <div css={metaLabel}>{t.performanceGate.windowDays}</div>
                    <div css={metaValue}>{config.baseline.window_days}d</div>
                  </div>
                )}
              </div>

              {result.metrics.length > 0 && (
                <>
                  <div
                    css={{
                      fontSize: '0.6875rem',
                      fontWeight: 600,
                      color: 'var(--text-disabled)',
                      letterSpacing: '0.07em',
                      textTransform: 'uppercase',
                      marginBottom: 10,
                    }}
                  >
                    {t.performanceGate.metricsComparison}
                  </div>
                  <div css={metricsGrid}>
                    {result.metrics.map((m) => (
                      <MetricCard key={m.name} metric={m} t={t} />
                    ))}
                  </div>
                </>
              )}
            </>
          )}
        </div>
      )}
    </div>
  )
}

export function PerformanceGateJobBadge({ label }: { label: string }) {
  return (
    <span
      css={{
        fontSize: '0.5625rem',
        fontWeight: 800,
        padding: '1px 5px',
        borderRadius: 'var(--radius-full)',
        background: 'var(--accent-muted)',
        color: 'var(--accent)',
        letterSpacing: '0.06em',
        flexShrink: 0,
      }}
      title={label}
    >
      GATE
    </span>
  )
}
