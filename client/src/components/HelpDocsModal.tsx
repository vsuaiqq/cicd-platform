import { useEffect, useRef, useState } from 'react'
import { css } from '@emotion/react'
import { useI18n } from '../i18n'

const backdrop = css({
  position: 'fixed',
  inset: 0,
  background: 'var(--overlay-backdrop)',
  backdropFilter: 'blur(4px)',
  zIndex: 8000,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  padding: 24,
  animation: 'modal-backdrop 0.15s ease both',
})

const dialog = css({
  width: '100%',
  maxWidth: 860,
  height: '85vh',
  background: 'var(--bg-elevated)',
  border: '1px solid var(--border-default)',
  borderRadius: 'var(--radius-xl)',
  boxShadow: 'var(--shadow-overlay), 0 0 0 1px var(--overlay-inset) inset',
  display: 'flex',
  flexDirection: 'column',
  overflow: 'hidden',
  animation: 'modal-slide 0.18s var(--ease-out) both',
})

const modalHeader = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: '16px 24px',
  borderBottom: '1px solid var(--border-subtle)',
  flexShrink: 0,
})

const modalTitle = css({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  fontSize: '1rem',
  fontWeight: 700,
  letterSpacing: '-0.02em',
  color: 'var(--text-primary)',
})

const closeBtn = css({
  width: 30,
  height: 30,
  borderRadius: 7,
  border: '1px solid var(--border-subtle)',
  background: 'transparent',
  color: 'var(--text-disabled)',
  cursor: 'pointer',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  transition: 'background 120ms, color 120ms',
  '&:hover': { background: 'var(--bg-hover)', color: 'var(--text-secondary)' },
})

const body = css({
  display: 'flex',
  flex: 1,
  overflow: 'hidden',
})

const sidebar = css({
  width: 200,
  flexShrink: 0,
  borderRight: '1px solid var(--border-subtle)',
  overflowY: 'auto',
  padding: '12px 8px',
  display: 'flex',
  flexDirection: 'column',
  gap: 2,
  '&::-webkit-scrollbar': { width: 0 },
})

const sideSection = css({
  fontSize: '0.6rem',
  fontWeight: 700,
  textTransform: 'uppercase',
  letterSpacing: '0.1em',
  color: 'var(--text-disabled)',
  padding: '8px 8px 3px',
  marginTop: 4,
  '&:first-of-type': { marginTop: 0 },
})

const sideItem = (active: boolean) => css({
  padding: '6px 10px',
  borderRadius: 6,
  fontSize: '0.8125rem',
  fontWeight: active ? 600 : 400,
  color: active ? 'var(--text-primary)' : 'var(--text-tertiary)',
  background: active ? 'var(--accent-muted)' : 'transparent',
  cursor: 'pointer',
  transition: 'background 100ms, color 100ms',
  '&:hover': {
    background: active ? 'var(--accent-muted)' : 'var(--bg-hover)',
    color: 'var(--text-primary)',
  },
})

const content = css({
  flex: 1,
  overflowY: 'auto',
  padding: '28px 32px',
  '&::-webkit-scrollbar': { width: 6 },
  '&::-webkit-scrollbar-track': { background: 'transparent' },
  '&::-webkit-scrollbar-thumb': { background: 'var(--border-default)', borderRadius: 4 },
})

const sectionTitle = css({
  fontSize: '1.125rem',
  fontWeight: 700,
  letterSpacing: '-0.02em',
  color: 'var(--text-primary)',
  marginBottom: 6,
  paddingTop: 4,
  display: 'flex',
  alignItems: 'center',
  gap: 8,
})

const sectionDesc = css({
  fontSize: '0.875rem',
  color: 'var(--text-tertiary)',
  lineHeight: 1.65,
  marginBottom: 20,
})

const divider = css({
  borderTop: '1px solid var(--border-subtle)',
  margin: '28px 0',
})

const fieldGrid = css({
  display: 'grid',
  gridTemplateColumns: '180px 1fr',
  gap: '1px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-subtle)',
  overflow: 'hidden',
  marginBottom: 20,
})

const fieldRow = css({
  display: 'contents',
})

const fieldKey = css({
  padding: '10px 14px',
  background: 'var(--bg-overlay)',
  borderBottom: '1px solid var(--border-subtle)',
  display: 'flex',
  flexDirection: 'column',
  gap: 3,
  '&:last-of-type': { borderBottom: 'none' },
})

const fieldName = css({
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.8125rem',
  fontWeight: 600,
  color: 'var(--accent)',
})

const fieldRequired = css({
  fontSize: '0.6rem',
  fontWeight: 700,
  letterSpacing: '0.06em',
  textTransform: 'uppercase',
  color: 'var(--danger)',
  background: 'var(--danger-muted)',
  padding: '1px 5px',
  borderRadius: 4,
  display: 'inline-block',
  width: 'fit-content',
})

const fieldOptional = css({
  fontSize: '0.6rem',
  fontWeight: 700,
  letterSpacing: '0.06em',
  textTransform: 'uppercase',
  color: 'var(--text-disabled)',
  background: 'var(--bg-overlay)',
  border: '1px solid var(--border-default)',
  padding: '1px 5px',
  borderRadius: 4,
  display: 'inline-block',
  width: 'fit-content',
})

const fieldVal = css({
  padding: '10px 14px',
  background: 'var(--bg-card)',
  borderBottom: '1px solid var(--border-subtle)',
  fontSize: '0.8125rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.6,
  '&:last-of-type': { borderBottom: 'none' },
})

const fieldType = css({
  display: 'inline-block',
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.6875rem',
  color: 'var(--ai)',
  background: 'var(--ai-muted)',
  border: '1px solid var(--ai-glow)',
  padding: '1px 6px',
  borderRadius: 4,
  marginBottom: 4,
})

const codeBlock = css({
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.8rem',
  lineHeight: 1.75,
  background: 'var(--code-bg)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-md)',
  padding: '16px 18px',
  color: 'var(--text-primary)',
  whiteSpace: 'pre',
  overflowX: 'auto',
  marginBottom: 20,
  position: 'relative',
})

const inlineCode = css({
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.875em',
  color: 'var(--accent)',
  background: 'var(--accent-muted)',
  padding: '1px 5px',
  borderRadius: 4,
})

function DocText({ text }: { text: string }) {
  const parts = text.split(/(`[^`]+`)/g)
  return (
    <>
      {parts.map((part, i) =>
        part.startsWith('`') && part.endsWith('`')
          ? <code key={i} css={inlineCode}>{part.slice(1, -1)}</code>
          : part
      )}
    </>
  )
}

const tipBox = css({
  display: 'flex',
  gap: 10,
  padding: '12px 14px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--bg-active)',
  border: '1px solid var(--accent-glow)',
  marginBottom: 20,
  fontSize: '0.8125rem',
  color: 'var(--text-secondary)',
  lineHeight: 1.6,
})

const warnBox = css({
  display: 'flex',
  gap: 10,
  padding: '12px 14px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--warning-muted)',
  border: '1px solid var(--warning-glow)',
  marginBottom: 20,
  fontSize: '0.8125rem',
  color: 'var(--warning)',
  lineHeight: 1.6,
})

function Y({ children }: { children: string }) {
  const lines = children.split('\n').map((line, i) => {
    const commentMatch = line.match(/^(\s*)(#.*)$/)
    if (commentMatch) return (
      <div key={i}>
        <span>{commentMatch[1]}</span>
        <span style={{ color: 'var(--syntax-comment)' }}>{commentMatch[2]}</span>
      </div>
    )
    const keyMatch = line.match(/^(\s*)([^:]+)(:\s*)(.*)$/)
    if (keyMatch) {
      const val = keyMatch[4]
      let valEl: React.ReactNode = val
      if (/^".*"$/.test(val) || /^'.*'$/.test(val)) {
        valEl = <span style={{ color: 'var(--syntax-string)' }}>{val}</span>
      } else if (/^(true|false|null)$/.test(val)) {
        valEl = <span style={{ color: 'var(--syntax-boolean)' }}>{val}</span>
      } else if (/^\d/.test(val)) {
        valEl = <span style={{ color: 'var(--syntax-number)' }}>{val}</span>
      } else if (val.startsWith('-') || val === '') {
        valEl = <span style={{ color: 'var(--syntax-value)' }}>{val}</span>
      } else {
        valEl = <span style={{ color: 'var(--syntax-value)' }}>{val}</span>
      }
      return (
        <div key={i}>
          <span>{keyMatch[1]}</span>
          <span style={{ color: 'var(--syntax-key)' }}>{keyMatch[2]}</span>
          <span style={{ color: 'var(--syntax-punctuation)' }}>{keyMatch[3]}</span>
          {valEl}
        </div>
      )
    }
    const listMatch = line.match(/^(\s*- )(.*)$/)
    if (listMatch) return (
      <div key={i}>
        <span style={{ color: 'var(--syntax-punctuation)' }}>{listMatch[1]}</span>
        <span style={{ color: 'var(--syntax-string)' }}>{listMatch[2]}</span>
      </div>
    )
    return <div key={i}><span style={{ color: 'var(--syntax-punctuation)' }}>{line}</span></div>
  })
  return <>{lines}</>
}

function Field({
  name,
  type,
  required,
  children,
}: {
  name: string
  type: string
  required?: boolean
  children: React.ReactNode
}) {
  const { t } = useI18n()
  return (
    <div css={fieldRow}>
      <div css={fieldKey}>
        <span css={fieldName}>{name}</span>
        <span css={required ? fieldRequired : fieldOptional}>
          {required ? t.help.fieldRequired : t.help.fieldOptional}
        </span>
      </div>
      <div css={fieldVal}>
        <span css={fieldType}>{type}</span>
        <div>{children}</div>
      </div>
    </div>
  )
}

const STATIC_SECTIONS = [
  { id: 'overview',        labelKey: 'overview',        group: 'gettingStarted' },
  { id: 'global',          labelKey: 'global',          group: 'gettingStarted' },
  { id: 'on',              labelKey: 'on',              group: 'gettingStarted' },
  { id: 'jobs',            labelKey: 'jobs',            group: 'jobs' },
  { id: 'job-fields',      labelKey: 'jobFields',       group: 'jobs' },
  { id: 'steps',           labelKey: 'steps',           group: 'jobs' },
  { id: 'cache',           labelKey: 'cache',           group: 'jobs' },
  { id: 'artifacts',       labelKey: 'artifacts',       group: 'jobs' },
  { id: 'load-testing',    labelKey: 'loadTesting',     group: 'quality' },
  { id: 'performance-gate', labelKey: 'performanceGate', group: 'quality' },
  { id: 'secrets',         labelKey: 'secrets',         group: 'advanced' },
  { id: 'limitations',     labelKey: 'limitations',     group: 'advanced' },
  { id: 'example',         labelKey: 'example',         group: 'advanced' },
] as const

function SectionOverview() {
  const { t } = useI18n()
  const d = t.helpDocs.overview
  return (
    <>
      <h2 css={sectionTitle}>
        <IconFile />
        {d.title}
      </h2>
      <p css={sectionDesc}>
        <DocText text={d.desc} />
      </p>
      <div css={tipBox}>
        <span style={{ flexShrink: 0, marginTop: 1 }}>💡</span>
        <span>
          <DocText text={d.tip} />
        </span>
      </div>
      <p css={[sectionDesc, { marginBottom: 8 }]}>{d.minimalExample}</p>
      <div css={codeBlock}>
        <Y>{`name: My pipeline

on:
  push:
    branches:
      - main

jobs:
  build:
    image: node:20
    steps:
      - name: Install
        run: npm ci
      - name: Build
        run: npm run build`}
        </Y>
      </div>
    </>
  )
}

function SectionGlobal() {
  const { t } = useI18n()
  const d = t.helpDocs.global
  return (
    <>
      <h2 css={sectionTitle}><IconSliders />{d.title}</h2>
      <p css={sectionDesc}>{d.desc}</p>
      <div css={fieldGrid}>
        <Field name="name" type="string" required>
          <DocText text={d.fields.name} />
        </Field>
        <Field name="on" type="object" required>
          <DocText text={d.fields.on} />
        </Field>
        <Field name="env" type="map[string]string">
          <DocText text={d.fields.env} />
        </Field>
        <Field name="jobs" type="map[string]Job" required>
          <DocText text={d.fields.jobs} />
        </Field>
      </div>
    </>
  )
}

function SectionOn() {
  const { t } = useI18n()
  const d = t.helpDocs.on
  return (
    <>
      <h2 css={sectionTitle}><IconZap />{d.title}</h2>
      <p css={sectionDesc}>
        <DocText text={d.desc} />
      </p>
      <div css={fieldGrid}>
        <Field name="on.push" type="object">
          <DocText text={d.fields.push} />
        </Field>
        <Field name="on.push.branches" type="string[]">
          <DocText text={d.fields.branches} />
        </Field>
      </div>
      <div css={codeBlock}><Y>{`on:
  push:
    branches:
      - main
      - develop`}</Y></div>
      <div css={tipBox}>
        <span style={{ flexShrink: 0 }}>💡</span>
        <span><DocText text={d.exampleNote} /></span>
      </div>
    </>
  )
}

function SectionJobs() {
  const { t } = useI18n()
  const d = t.helpDocs.jobs
  return (
    <>
      <h2 css={sectionTitle}><IconPipeline />{d.title}</h2>
      <p css={sectionDesc}>
        <DocText text={d.desc} />
      </p>
      <div css={codeBlock}><Y>{`jobs:
  lint:          # job ID
    name: Lint   # display name
    image: golang:1.24
    steps: ...

  test:
    name: Test
    needs: [lint]   # runs after lint succeeds
    image: golang:1.24
    steps: ...

  build:
    needs: [test]
    image: golang:1.24
    steps: ...`}</Y></div>
      <div css={tipBox}>
        <span style={{ flexShrink: 0 }}>💡</span>
        <span>
          <DocText text={d.tip} />
        </span>
      </div>
    </>
  )
}

function SectionJobFields() {
  const { t } = useI18n()
  const d = t.helpDocs.jobFields
  return (
    <>
      <h2 css={sectionTitle}><IconSliders />{d.title}</h2>
      <p css={sectionDesc}>{d.desc}</p>
      <div css={fieldGrid}>
        <Field name="name" type="string">
          <DocText text={d.fields.name} />
        </Field>
        <Field name="image" type="string">
          <DocText text={d.fields.image} />
        </Field>
        <Field name="needs" type="string[]">
          <DocText text={d.fields.needs} />
        </Field>
        <Field name="env" type="map[string]string">
          <DocText text={d.fields.env} />
        </Field>
        <Field name="timeout" type="duration">
          <DocText text={d.fields.timeout} />
        </Field>
        <Field name="retry" type="integer">
          <DocText text={d.fields.retry} />
        </Field>
        <Field name="approval" type="boolean | required">
          <DocText text={d.fields.approval} />
        </Field>
        <Field name="cache" type="object">
          <DocText text={d.fields.cache} />
        </Field>
        <Field name="artifacts" type="object">
          <DocText text={d.fields.artifacts} />
        </Field>
        <Field name="steps" type="Step[]">
          <DocText text={d.fields.steps} />
        </Field>
        <Field name="performance_gate" type="PerformanceGateConfig">
          <DocText text={d.fields.performance_gate} />
        </Field>
      </div>
    </>
  )
}

function SectionSteps() {
  const { t } = useI18n()
  const d = t.helpDocs.steps
  return (
    <>
      <h2 css={sectionTitle}><IconList />{d.title}</h2>
      <p css={sectionDesc}>
        <DocText text={d.desc} />
      </p>
      <div css={fieldGrid}>
        <Field name="name" type="string" required>
          <DocText text={d.fields.name} />
        </Field>
        <Field name="run" type="string" required>
          <DocText text={d.fields.run} />
        </Field>
        <Field name="timeout" type="duration">
          <DocText text={d.fields.timeout} />
        </Field>
        <Field name="retry" type="integer">
          <DocText text={d.fields.retry} />
        </Field>
        <Field name="continue_on_error" type="boolean">
          <DocText text={d.fields.continue_on_error} />
        </Field>
      </div>
      <div css={codeBlock}><Y>{`steps:
  - name: Run tests
    run: go test -v ./...
    timeout: 10m
    retry: 1

  - name: Upload coverage
    run: bash <(curl -s https://codecov.io/bash)
    continue_on_error: true   # don't fail build if upload fails`}</Y></div>
    </>
  )
}

function SectionLoadTesting() {
  const { t } = useI18n()
  const d = t.helpDocs.loadTesting
  return (
    <>
      <h2 css={sectionTitle}><IconZap />{d.title}</h2>
      <p css={sectionDesc}><DocText text={d.desc} /></p>
      <div css={fieldGrid}>
        <Field name=".flow/perf-metrics.json" type="file">
          <DocText text={d.fields.metricsFile} />
        </Field>
        <Field name="format" type="JSON">
          <DocText text={d.fields.format} />
        </Field>
        <Field name="metrics.*" type="number">
          <DocText text={d.fields.exampleMetrics} />
        </Field>
      </div>
      <div css={codeBlock}><Y>{`{
  "version": 1,
  "tool": "curl-load",
  "metrics": {
    "http_req_duration_p95": 142.5,
    "http_req_duration_avg": 87.3,
    "http_reqs": 200,
    "http_req_failed_rate": 0.0
  }
}`}</Y></div>
      <div css={codeBlock}><Y>{`load-test:
  name: Load Test
  image: golang:1.25
  needs: [deploy-staging]
  steps:
    - name: Run load test
      run: sh scripts/load-test.sh`}</Y></div>
      <div css={tipBox}>
        <span style={{ flexShrink: 0 }}>💡</span>
        <span><DocText text={d.tip} /></span>
      </div>
      <p css={sectionDesc}><DocText text={d.exampleNote} /></p>
    </>
  )
}

function SectionPerformanceGate() {
  const { t } = useI18n()
  const d = t.helpDocs.performanceGate
  return (
    <>
      <h2 css={sectionTitle}><IconPipeline />{d.title}</h2>
      <p css={sectionDesc}><DocText text={d.desc} /></p>
      <div css={fieldGrid}>
        <Field name="performance_gate.source_job" type="string" required>
          <DocText text={d.fields.source_job} />
        </Field>
        <Field name="performance_gate.metrics[]" type="MetricRule[]">
          <DocText text={d.fields.metrics} />
        </Field>
        <Field name="performance_gate.baseline.window_days" type="integer">
          <DocText text={d.fields.baseline_window_days} />
        </Field>
        <Field name="performance_gate.baseline.min_samples" type="integer">
          <DocText text={d.fields.baseline_min_samples} />
        </Field>
        <Field name="performance_gate.baseline.branch" type="string">
          <DocText text={d.fields.baseline_branch} />
        </Field>
        <Field name="performance_gate.adaptive.enabled" type="boolean">
          <DocText text={d.fields.adaptive_enabled} />
        </Field>
        <Field name="performance_gate.adaptive.sigma_factor" type="number">
          <DocText text={d.fields.adaptive_sigma_factor} />
        </Field>
        <Field name="performance_gate.adaptive.max_regression_pct" type="number">
          <DocText text={d.fields.adaptive_max_regression_pct} />
        </Field>
      </div>
      <div css={codeBlock}><Y>{`performance-gate:
  name: Performance Quality Gate
  needs: [load-test]
  performance_gate:
    source_job: load-test
    metrics:
      - name: http_req_duration_p95
        direction: lower_is_better
        max: 500
      - name: http_reqs
        direction: higher_is_better
        min: 100
    baseline:
      window_days: 30
      min_samples: 3
    adaptive:
      enabled: true
      sigma_factor: 2.0
      max_regression_pct: 15`}</Y></div>
      <div css={tipBox}>
        <span style={{ flexShrink: 0 }}>💡</span>
        <span><DocText text={d.adaptiveFormula} /></span>
      </div>
      <div css={warnBox}>
        <span style={{ flexShrink: 0, marginTop: 1 }}>⚠️</span>
        <span><DocText text={d.coldStart} /></span>
      </div>
      <p css={sectionDesc}><DocText text={d.uiNote} /></p>
      <p css={sectionDesc}><DocText text={d.defaultMetrics} /></p>
    </>
  )
}

function SectionCache() {
  const { t } = useI18n()
  const d = t.helpDocs.cache
  return (
    <>
      <h2 css={sectionTitle}><IconBox />{d.title}</h2>
      <p css={sectionDesc}>{d.desc}</p>
      <div css={fieldGrid}>
        <Field name="key" type="string" required>
          <DocText text={d.fields.key} />
        </Field>
        <Field name="paths" type="string[]" required>
          <DocText text={d.fields.paths} />
        </Field>
      </div>
      <div css={codeBlock}><Y>{`cache:
  key: go-\${{ checksum "go.sum" }}
  paths:
    - /go/pkg/mod
    - /root/.cache/go-build`}</Y></div>
      <div css={tipBox}>
        <span style={{ flexShrink: 0 }}>💡</span>
        <span>
          <DocText text={d.tip} />
        </span>
      </div>
    </>
  )
}

function SectionArtifacts() {
  const { t } = useI18n()
  const d = t.helpDocs.artifacts
  return (
    <>
      <h2 css={sectionTitle}><IconPaperclip />{d.title}</h2>
      <p css={sectionDesc}>{d.desc}</p>
      <div css={fieldGrid}>
        <Field name="paths" type="string[]">
          <DocText text={d.fields.paths} />
        </Field>
        <Field name="download" type="ArtifactDownload[]">
          <DocText text={d.fields.download} />
        </Field>
      </div>
      <div css={codeBlock}><Y>{`build:
  artifacts:
    paths:
      - ./bin/

deploy:
  needs: [build]
  artifacts:
    download:
      - job: build
  steps:
    - name: Use binary
      run: ls -la ./bin/`}</Y></div>
    </>
  )
}

function SectionSecrets() {
  const { t } = useI18n()
  const d = t.helpDocs.secrets
  return (
    <>
      <h2 css={sectionTitle}><IconLock />{d.title}</h2>
      <p css={sectionDesc}>{d.desc}</p>
      <ol style={{ paddingLeft: 20, color: 'var(--text-secondary)', fontSize: '0.875rem', lineHeight: 2, marginBottom: 20 }}>
        {d.priorityList.map((item, i) => (
          <li key={i}><DocText text={item} /></li>
        ))}
      </ol>
      <p css={[sectionDesc, { fontWeight: 600, color: 'var(--text-primary)' }]}>{d.interpolationTitle}</p>
      <ul style={{ paddingLeft: 20, color: 'var(--text-secondary)', fontSize: '0.875rem', lineHeight: 1.85, marginBottom: 20 }}>
        {d.interpolationList.map((item, i) => (
          <li key={i}><DocText text={item} /></li>
        ))}
      </ul>
      <div css={warnBox}>
        <span style={{ flexShrink: 0, marginTop: 1 }}>⚠️</span>
        <span>
          <DocText text={d.warn} />
        </span>
      </div>
      <div css={codeBlock}><Y>{`# .cicd.yaml
env:
  NODE_ENV: production

jobs:
  deploy:
    env:
      DEPLOY_TOKEN: \${{ secrets.DEPLOY_TOKEN }}
    steps:
      - name: Deploy
        run: |
          curl -H "Authorization: Bearer $DEPLOY_TOKEN" \\
               https://api.example.com/deploy`}</Y></div>
    </>
  )
}

function SectionLimitations() {
  const { t } = useI18n()
  const d = t.helpDocs.limitations
  return (
    <>
      <h2 css={sectionTitle}><IconAlert />{d.title}</h2>
      <p css={sectionDesc}>{d.desc}</p>
      <ul style={{ paddingLeft: 20, color: 'var(--text-secondary)', fontSize: '0.875rem', lineHeight: 1.85, marginBottom: 20 }}>
        {d.items.map((item, i) => (
          <li key={i}><DocText text={item} /></li>
        ))}
      </ul>
    </>
  )
}

function SectionExample() {
  const { t } = useI18n()
  const d = t.helpDocs.example
  return (
    <>
      <h2 css={sectionTitle}><IconStar />{d.title}</h2>
      <p css={sectionDesc}>{d.desc}</p>
      <div css={codeBlock}><Y>{`name: hello-server CI/CD

on:
  push:
    branches:
      - main

env:
  APP_NAME: hello-server

jobs:
  lint:
    image: golang:1.25
    steps:
      - name: go vet
        run: go vet ./...

  test:
    needs: [lint]
    image: golang:1.25
    steps:
      - name: Unit tests
        run: go test -v ./...

  build:
    needs: [test]
    image: golang:1.25
    artifacts:
      paths: [./bin/]
    steps:
      - name: Compile
        run: go build -o bin/$APP_NAME ./cmd/server/

  deploy-staging:
    needs: [build]
    image: alpine:3.19
    steps:
      - name: Deploy staging
        run: echo "Deploying $APP_NAME to staging..."

  load-test:
    needs: [deploy-staging]
    image: golang:1.25
    steps:
      - name: Run load test
        run: sh scripts/load-test.sh

  performance-gate:
    needs: [load-test]
    performance_gate:
      source_job: load-test
      metrics:
        - name: http_req_duration_p95
          direction: lower_is_better
          max: 500
        - name: http_reqs
          direction: higher_is_better
          min: 100
      baseline:
        window_days: 30
        min_samples: 3
      adaptive:
        enabled: true
        sigma_factor: 2.0
        max_regression_pct: 15

  deploy-production:
    needs: [performance-gate]
    image: alpine:3.19
    steps:
      - name: Deploy production
        run: echo "Deploying $APP_NAME to production..."`}</Y></div>
    </>
  )
}

function IconFile() {
  return <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden><path d="M3 2h7l3 3v9H3V2z" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" /><path d="M10 2v3h3" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" /></svg>
}
function IconSliders() {
  return <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden><path d="M2 4h12M2 8h12M2 12h12" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" /><circle cx="5" cy="4" r="1.5" fill="var(--bg-elevated)" stroke="currentColor" strokeWidth="1.2" /><circle cx="11" cy="8" r="1.5" fill="var(--bg-elevated)" stroke="currentColor" strokeWidth="1.2" /><circle cx="7" cy="12" r="1.5" fill="var(--bg-elevated)" stroke="currentColor" strokeWidth="1.2" /></svg>
}
function IconZap() {
  return <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden><path d="M9 2L4 9h5l-2 5 7-7H9l2-5z" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" /></svg>
}
function IconPipeline() {
  return <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden><circle cx="3" cy="8" r="2" stroke="currentColor" strokeWidth="1.2" /><circle cx="8" cy="8" r="2" stroke="currentColor" strokeWidth="1.2" /><circle cx="13" cy="8" r="2" stroke="currentColor" strokeWidth="1.2" /><path d="M5 8h1M10 8h1" stroke="currentColor" strokeWidth="1.2" /></svg>
}
function IconList() {
  return <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden><path d="M6 4h8M6 8h8M6 12h8" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" /><circle cx="3" cy="4" r="1" fill="currentColor" /><circle cx="3" cy="8" r="1" fill="currentColor" /><circle cx="3" cy="12" r="1" fill="currentColor" /></svg>
}
function IconBox() {
  return <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden><path d="M2 5l6-3 6 3v6l-6 3-6-3V5z" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" /><path d="M2 5l6 3m0 0l6-3m-6 3v6" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" /></svg>
}
function IconPaperclip() {
  return <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden><path d="M13 7L7.5 12.5a3.5 3.5 0 01-5-5l7-7a2 2 0 013 3l-6.5 6.5a.5.5 0 01-.7-.7L10 5" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round" /></svg>
}
function IconLock() {
  return <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden><rect x="3" y="7" width="10" height="7" rx="1.5" stroke="currentColor" strokeWidth="1.2" /><path d="M5.5 7V5a2.5 2.5 0 015 0v2" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" /><circle cx="8" cy="11" r="1" fill="currentColor" /></svg>
}
function IconStar() {
  return <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden><path d="M8 2l1.5 3.5L13 6l-2.5 2.5.5 3.5L8 10.5 5 12l.5-3.5L3 6l3.5-.5L8 2z" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" /></svg>
}
function IconAlert() {
  return <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden><path d="M8 2l6.5 11H1.5L8 2z" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" /><path d="M8 6v3.5M8 11.5v.5" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" /></svg>
}

interface HelpDocsModalProps {
  open: boolean
  onClose: () => void
}

const SECTION_MAP: Record<string, React.ReactNode> = {
  overview:         <SectionOverview />,
  global:           <SectionGlobal />,
  on:               <SectionOn />,
  jobs:             <SectionJobs />,
  'job-fields':     <SectionJobFields />,
  steps:            <SectionSteps />,
  cache:            <SectionCache />,
  artifacts:        <SectionArtifacts />,
  'load-testing':   <SectionLoadTesting />,
  'performance-gate': <SectionPerformanceGate />,
  secrets:          <SectionSecrets />,
  limitations:      <SectionLimitations />,
  example:          <SectionExample />,
}

export default function HelpDocsModal({ open, onClose }: HelpDocsModalProps) {
  const { t } = useI18n()
  const [active, setActive] = useState('overview')
  const contentRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open) return
    const handler = (e: KeyboardEvent) => {
      if (e.key === 'Escape') onClose()
    }
    window.addEventListener('keydown', handler)
    return () => window.removeEventListener('keydown', handler)
  }, [open, onClose])

  useEffect(() => {
    if (contentRef.current) contentRef.current.scrollTop = 0
  }, [active])

  if (!open) return null

  const SECTIONS = STATIC_SECTIONS.map(s => ({
    id: s.id,
    label: t.helpSectionLabels[s.labelKey],
    group: t.helpGroups[s.group],
  }))

  const groups = Array.from(new Set(SECTIONS.map(s => s.group)))

  return (
    <div css={backdrop} onMouseDown={(e) => e.target === e.currentTarget && onClose()}>
      <div css={dialog} role="dialog" aria-modal aria-label={t.help.docAriaLabel}>
        <div css={modalHeader}>
          <div css={modalTitle}>
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden>
              <path d="M3 2h7l3 3v9H3V2z" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" />
              <path d="M10 2v3h3" stroke="currentColor" strokeWidth="1.2" strokeLinejoin="round" />
            </svg>
            {t.help.title}
          </div>
          <button type="button" css={closeBtn} onClick={onClose} aria-label={t.help.close}>
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
              <path d="M2 2l8 8M10 2l-8 8" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
            </svg>
          </button>
        </div>

        <div css={body}>
          <nav css={sidebar} aria-label={t.aria.docSections}>
            {groups.map(group => (
              <div key={group}>
                <div css={sideSection}>{group}</div>
                {SECTIONS.filter(s => s.group === group).map(s => (
                  <div
                    key={s.id}
                    css={sideItem(active === s.id)}
                    onClick={() => setActive(s.id)}
                    role="button"
                    tabIndex={0}
                    onKeyDown={e => e.key === 'Enter' && setActive(s.id)}
                  >
                    {s.label}
                  </div>
                ))}
              </div>
            ))}
          </nav>

          <div css={content} ref={contentRef}>
            {SECTION_MAP[active]}
            <div css={divider} />
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              {SECTIONS.findIndex(s => s.id === active) > 0 && (
                <button
                  type="button"
                  onClick={() => setActive(SECTIONS[SECTIONS.findIndex(s => s.id === active) - 1].id)}
                  style={{
                    background: 'none',
                    border: '1px solid var(--border-default)',
                    borderRadius: 8,
                    padding: '7px 14px',
                    color: 'var(--text-secondary)',
                    fontSize: '0.8125rem',
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    gap: 6,
                  }}
                >
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                    <path d="M9 6H1M4.5 2.5L1 6l3.5 3.5" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                  {SECTIONS[SECTIONS.findIndex(s => s.id === active) - 1].label}
                </button>
              )}
              {SECTIONS.findIndex(s => s.id === active) < SECTIONS.length - 1 && (
                <button
                  type="button"
                  onClick={() => setActive(SECTIONS[SECTIONS.findIndex(s => s.id === active) + 1].id)}
                  style={{
                    marginLeft: 'auto',
                    background: 'var(--accent)',
                    border: 'none',
                    borderRadius: 8,
                    padding: '7px 14px',
                    color: 'var(--text-on-accent)',
                    fontSize: '0.8125rem',
                    fontWeight: 600,
                    cursor: 'pointer',
                    display: 'flex',
                    alignItems: 'center',
                    gap: 6,
                  }}
                >
                  {SECTIONS[SECTIONS.findIndex(s => s.id === active) + 1].label}
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                    <path d="M2 6h8M6.5 2.5L10 6l-3.5 3.5" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                </button>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
