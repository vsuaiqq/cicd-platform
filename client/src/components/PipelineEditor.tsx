import { useMemo, useCallback } from 'react'
import { css } from '@emotion/react'
import CodeMirror, { type ReactCodeMirrorProps } from '@uiw/react-codemirror'
import { yaml } from '@codemirror/lang-yaml'
import { getCodeMirrorTheme } from '../lib/codeMirrorTheme'
import { useTheme } from '../lib/themeContext'
import { linter, lintGutter, type Diagnostic } from '@codemirror/lint'
import { EditorView } from '@codemirror/view'
import jsYaml from 'js-yaml'
import { useI18n, type Translations } from '../i18n'



export interface PipelineValidationError {

  line?: number
  message: string
  severity: 'error' | 'warning'
}




function findJobLine(src: string, jobName: string): number | undefined {
  const lines = src.split('\n')
  let inJobs = false
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    if (/^jobs\s*:/.test(line)) { inJobs = true; continue }
    if (inJobs) {

      if (/^[a-zA-Z_]/.test(line) && !/^\s/.test(line)) { inJobs = false; break }
      if (new RegExp(`^\\s+${escapeRegExp(jobName)}\\s*:`).test(line)) return i + 1
    }
  }
  return undefined
}

function findStepLine(src: string, jobName: string, stepIdx: number): number | undefined {
  const lines = src.split('\n')
  let inJob = false
  let stepCount = -1
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    if (new RegExp(`^\\s+${escapeRegExp(jobName)}\\s*:`).test(line)) { inJob = true; continue }
    if (inJob) {

      if (/^\s{0,2}\S/.test(line) && !/^\s{4,}/.test(line)) { inJob = false; break }
      if (/^\s*-\s/.test(line)) {
        stepCount++
        if (stepCount === stepIdx) return i + 1
      }
    }
  }
  return undefined
}

function escapeRegExp(s: string): string {
  return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

type EditorMsg = Translations['editor']

function fill(template: string, vars: Record<string, string | number>): string {
  return Object.entries(vars).reduce(
    (s, [k, v]) => s.replaceAll(`{${k}}`, String(v)),
    template,
  )
}

export function validatePipelineYAML(content: string, m: EditorMsg): PipelineValidationError[] {
  if (!content.trim()) return []

  let doc: unknown
  try {
    doc = jsYaml.load(content)
  } catch (e: unknown) {
    const ye = e as { mark?: { line?: number }; reason?: string; message?: string }
    const line = ye.mark?.line != null ? ye.mark.line + 1 : undefined
    return [{ line, message: ye.reason ?? (ye as Error).message ?? m.yamlParseError, severity: 'error' }]
  }

  if (doc === null || typeof doc !== 'object' || Array.isArray(doc)) {
    return [{ message: m.mustBeMapping, severity: 'error' }]
  }

  const pipeline = doc as Record<string, unknown>
  const errors: PipelineValidationError[] = []


  if (pipeline.env !== undefined) {
    if (typeof pipeline.env !== 'object' || Array.isArray(pipeline.env)) {
      errors.push({ message: m.globalEnvMapping, severity: 'error' })
    }
  }


  if (pipeline.on !== undefined) {
    const on = pipeline.on as Record<string, unknown>
    if (on.push) {
      const push = on.push as Record<string, unknown>
      if (push.branches !== undefined && !Array.isArray(push.branches)) {
        errors.push({ message: m.branchesMustBeList, severity: 'error' })
      }
    }
  }


  if (!pipeline.jobs) {
    errors.push({ message: m.missingJobs, severity: 'error' })
    return errors
  }
  if (typeof pipeline.jobs !== 'object' || Array.isArray(pipeline.jobs)) {
    errors.push({ message: m.jobsMustBeMapping, severity: 'error' })
    return errors
  }

  const jobs = pipeline.jobs as Record<string, unknown>
  const jobNames = Object.keys(jobs)

  if (jobNames.length === 0) {
    errors.push({ message: m.jobsMinOne, severity: 'error' })
    return errors
  }

  for (const [jobName, rawJob] of Object.entries(jobs)) {
    const jobLine = findJobLine(content, jobName)

    if (!rawJob || typeof rawJob !== 'object' || Array.isArray(rawJob)) {
      errors.push({ line: jobLine, message: fill(m.jobMustBeMapping, { name: jobName }), severity: 'error' })
      continue
    }
    const job = rawJob as Record<string, unknown>

    if (job.image !== undefined && typeof job.image !== 'string') {
      errors.push({ line: jobLine, message: fill(m.jobImageString, { name: jobName }), severity: 'error' })
    }
    if (job.env !== undefined && (typeof job.env !== 'object' || Array.isArray(job.env))) {
      errors.push({ line: jobLine, message: fill(m.jobEnvMapping, { name: jobName }), severity: 'error' })
    }


    if (job.needs !== undefined) {
      if (!Array.isArray(job.needs)) {
        errors.push({ line: jobLine, message: fill(m.jobNeedsList, { name: jobName }), severity: 'error' })
      } else {
        for (const dep of job.needs as unknown[]) {
          if (typeof dep !== 'string') {
            errors.push({ line: jobLine, message: fill(m.jobNeedsEntryString, { name: jobName }), severity: 'error' })
          } else if (dep === jobName) {
            errors.push({ line: jobLine, message: fill(m.jobSelfDependency, { name: jobName }), severity: 'error' })
          } else if (!jobNames.includes(dep)) {
            errors.push({ line: jobLine, message: fill(m.jobUnknownDependency, { name: jobName, dep }), severity: 'error' })
          }
        }
      }
    }


    if (!job.steps) {
      errors.push({ line: jobLine, message: fill(m.jobMissingSteps, { name: jobName }), severity: 'error' })
      continue
    }
    if (!Array.isArray(job.steps)) {
      errors.push({ line: jobLine, message: fill(m.jobStepsList, { name: jobName }), severity: 'error' })
      continue
    }
    if (job.steps.length === 0) {
      errors.push({ line: jobLine, message: fill(m.jobStepsMinOne, { name: jobName }), severity: 'warning' })
      continue
    }

    job.steps.forEach((step: unknown, idx: number) => {
      const stepLine = findStepLine(content, jobName, idx) ?? jobLine
      const stepNum = idx + 1
      if (!step || typeof step !== 'object' || Array.isArray(step)) {
        errors.push({ line: stepLine, message: fill(m.stepMustBeMapping, { name: jobName, index: stepNum }), severity: 'error' })
        return
      }
      const s = step as Record<string, unknown>
      if (!s.name) {
        errors.push({ line: stepLine, message: fill(m.stepMissingName, { name: jobName, index: stepNum }), severity: 'error' })
      }
      if (!s.run) {
        errors.push({ line: stepLine, message: fill(m.stepMissingRun, { name: jobName, index: stepNum }), severity: 'error' })
      }
    })
  }


  const cycles = detectCycles(jobs)
  for (const cycle of cycles) {
    const firstJob = cycle[0]
    errors.push({
      line: findJobLine(content, firstJob),
      message: fill(m.circularDependency, { cycle: cycle.join(' → ') }),
      severity: 'error',
    })
  }

  return errors
}

function detectCycles(jobs: Record<string, unknown>): string[][] {
  const cycles: string[][] = []
  const visited = new Set<string>()
  const inStack = new Set<string>()
  const path: string[] = []

  function dfs(name: string) {
    if (inStack.has(name)) {
      const start = path.indexOf(name)
      cycles.push([...path.slice(start), name])
      return
    }
    if (visited.has(name)) return
    visited.add(name)
    inStack.add(name)
    path.push(name)
    const job = jobs[name]
    if (job && typeof job === 'object' && !Array.isArray(job)) {
      const needs = (job as Record<string, unknown>).needs
      if (Array.isArray(needs)) {
        for (const dep of needs) {
          if (typeof dep === 'string' && jobs[dep]) dfs(dep)
        }
      }
    }
    path.pop()
    inStack.delete(name)
  }

  for (const name of Object.keys(jobs)) {
    if (!visited.has(name)) dfs(name)
  }
  return cycles
}



const editorWrap = css({
  borderRadius: 'var(--radius-md)',
  overflow: 'hidden',
  border: '1px solid var(--border-default)',
  transition: 'border-color 150ms',
  '&:focus-within': { borderColor: 'var(--accent)' },

  '.cm-editor': {
    backgroundColor: 'var(--editor-bg) !important',
  },
  '.cm-scroller': {
    fontFamily: 'ui-monospace, SFMono-Regular, Menlo, monospace !important',
    fontSize: '0.8rem !important',
  },
  '.cm-gutters': {
    backgroundColor: 'var(--editor-bg) !important',
    borderRight: '1px solid var(--editor-gutter-border) !important',
  },
  '.cm-activeLineGutter': {
    backgroundColor: 'var(--editor-line-active) !important',
  },
  '.cm-activeLine': {
    backgroundColor: 'var(--editor-line-active) !important',
  },
})

const errorPanel = css({
  marginTop: 10,
  display: 'flex',
  flexDirection: 'column',
  gap: 4,
})

const errorItem = (severity: 'error' | 'warning') => css({
  display: 'flex',
  alignItems: 'flex-start',
  gap: 8,
  padding: '7px 12px',
  borderRadius: 'var(--radius-md)',
  background: severity === 'error' ? 'var(--danger-muted)' : 'var(--warning-muted)',
  border: `1px solid ${severity === 'error' ? 'var(--danger-glow)' : 'var(--warning-glow)'}`,
  fontSize: '0.8125rem',
  color: severity === 'error' ? 'var(--danger)' : 'var(--warning)',
  fontFamily: 'ui-monospace, monospace',
  lineHeight: 1.5,
})

const errorDot = (severity: 'error' | 'warning') => css({
  width: 6,
  height: 6,
  borderRadius: '50%',
  background: severity === 'error' ? 'var(--danger)' : 'var(--warning)',
  marginTop: 5,
  flexShrink: 0,
})

const validMsg = css({
  display: 'flex',
  alignItems: 'center',
  gap: 6,
  padding: '6px 12px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--success-muted)',
  border: '1px solid var(--success-glow)',
  fontSize: '0.8125rem',
  color: 'var(--success)',
})



interface PipelineEditorProps {
  value: string
  onChange: (value: string) => void
  disabled?: boolean
  minHeight?: string
}

export function PipelineEditor({ value, onChange, disabled = false, minHeight = '300px' }: PipelineEditorProps) {
  const { t } = useI18n()
  const { theme } = useTheme()
  const editorMsg = t.editor
  const cmTheme = useMemo(() => getCodeMirrorTheme(theme === 'dark'), [theme])


  const yamlLinter = useMemo(() => linter((view) => {
    const src = view.state.doc.toString()
    const errors = validatePipelineYAML(src, editorMsg)
    const diagnostics: Diagnostic[] = []

    for (const err of errors) {
      if (err.line == null) continue
      const lineNum = Math.max(1, Math.min(err.line, view.state.doc.lines))
      const line = view.state.doc.line(lineNum)
      diagnostics.push({
        from: line.from,
        to: line.to,
        severity: err.severity,
        message: err.message,
      })
    }
    return diagnostics
  }, { delay: 400 }), [editorMsg])

  const extensions = useMemo(() => [
    yaml(),
    yamlLinter,
    lintGutter(),
    EditorView.lineWrapping,
    EditorView.editable.of(!disabled),
  ], [yamlLinter, disabled])


  const errors = useMemo(() => validatePipelineYAML(value, editorMsg), [value, editorMsg])
  const panelErrors = useMemo(() => errors.filter((e) => e.line == null), [errors])
  const hasErrors = errors.some((e) => e.severity === 'error')
  const hasWarnings = !hasErrors && errors.some((e) => e.severity === 'warning')
  const isValid = value.trim() !== '' && errors.length === 0

  const handleChange = useCallback<NonNullable<ReactCodeMirrorProps['onChange']>>(
    (val) => onChange(val),
    [onChange],
  )

  return (
    <div>
      <div css={editorWrap} style={{ opacity: disabled ? 0.45 : 1 }}>
        <CodeMirror
          key={theme}
          value={value}
          onChange={handleChange}
          extensions={extensions}
          theme={cmTheme}
          height={minHeight}
          basicSetup={{
            lineNumbers: true,
            highlightActiveLine: !disabled,
            foldGutter: true,
            autocompletion: false,
            indentOnInput: true,
            tabSize: 2,
          }}
          readOnly={disabled}
        />
      </div>


      {value.trim() !== '' && (
        <div css={errorPanel}>
          {isValid && (
            <div css={validMsg}>
              <svg width="13" height="13" viewBox="0 0 16 16" fill="none" aria-hidden>
                <circle cx="8" cy="8" r="6.5" stroke="currentColor" strokeWidth="1.2" />
                <path d="M5 8l2 2 4-4" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
              </svg>
              {editorMsg.validYaml}
            </div>
          )}
          {panelErrors.map((e, i) => (
            <div key={i} css={errorItem(e.severity)}>
              <span css={errorDot(e.severity)} />
              <span>{e.message}</span>
            </div>
          ))}
          {(hasErrors || hasWarnings) && panelErrors.length === 0 && (
            <div css={errorItem(hasErrors ? 'error' : 'warning')}>
              <span css={errorDot(hasErrors ? 'error' : 'warning')} />
              <span>
                {errors.filter(e => e.severity === 'error').length > 0
                  ? fill(editorMsg.errorsInline, { count: errors.filter(e => e.severity === 'error').length })
                  : fill(editorMsg.warningsInline, { count: errors.filter(e => e.severity === 'warning').length })}
              </span>
            </div>
          )}
          {panelErrors.length > 0 && errors.some(e => e.line != null) && (
            <div css={errorItem(hasErrors ? 'error' : 'warning')}>
              <span css={errorDot(hasErrors ? 'error' : 'warning')} />
              <span>
                {fill(editorMsg.additionalErrors, { count: errors.filter(e => e.line != null && e.severity === 'error').length })}
              </span>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
