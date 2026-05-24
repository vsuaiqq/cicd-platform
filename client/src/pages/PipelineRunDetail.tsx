import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { useParams, Link } from 'react-router-dom'
import { css } from '@emotion/react'
import { useAppSelector } from '../store'
import {
  useApproveJobMutation,
  useCancelRunMutation,
  useGetRunQuery,
  useRejectJobMutation,
} from '../store/api/apiSlice'
import {
  connectRunWS,
  type PipelineRun,
  type PipelineStep,
  type JobStatus,
  type RunStatus,
  type StepEvent,
  type WSEvent,
  durationMs,
  formatDuration,
  statusColor,
  statusBg,
} from '../api/pipelines'
import { analyzeFailure, type AnalysisResult } from '../api/ai'
import { Breadcrumb, RunDetailSkeleton, useToast } from '../components/ui'
import { btnSecondary } from '../styles/theme'
import { getRecentProjects } from '../lib/recentProjects'
import { formatStatus } from '../lib/formatStatus'
import { useI18n } from '../i18n'



const pageWrap = css({ animation: 'fade-in 0.25s var(--ease-out) both' })

const runHeader = css({
  marginBottom: 28,
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: 16,
  flexWrap: 'wrap',
})

const runTitle = css({
  fontSize: '1.5rem',
  fontWeight: 700,
  letterSpacing: '-0.03em',
  color: 'var(--text-primary)',
  marginBottom: 6,
})

const runMeta = css({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  flexWrap: 'wrap',
  fontSize: '0.8125rem',
  color: 'var(--text-tertiary)',
})

const metaChip = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 5,
  padding: '3px 9px',
  borderRadius: 'var(--radius-full)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-overlay)',
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.75rem',
})



const pipelineGraph = css({
  display: 'flex',
  alignItems: 'stretch',
  gap: 0,
  marginBottom: 28,
  overflowX: 'auto',
  paddingBottom: 8,
  paddingTop: 2,
})

const pipelineConnector = css({
  display: 'flex',
  alignItems: 'center',
  padding: '0 4px',
  color: 'var(--text-disabled)',
  flexShrink: 0,
  userSelect: 'none',
})

const jobNode = css({
  flex: '0 0 auto',
  width: 190,
  borderRadius: 'var(--radius-lg)',
  border: '1px solid var(--border-subtle)',
  background: 'var(--bg-card)',
  cursor: 'pointer',
  transition: 'border-color 200ms, box-shadow 200ms, transform 150ms',
  overflow: 'hidden',
  '&:hover': {
    borderColor: 'var(--border-default)',
    boxShadow: 'var(--shadow-md)',
    transform: 'translateY(-1px)',
  },
  '&:active': { transform: 'translateY(0)' },
})

const jobNodeActive = css({
  borderColor: 'var(--accent)',
  boxShadow: '0 0 0 2px var(--accent-muted), var(--shadow-md)',
})

const jobNodeColorBar = (statusClr: string) => css({
  height: 3,
  background: statusClr,
  opacity: 0.8,
})

const jobNodeHeader = css({
  padding: '10px 12px 8px',
  borderBottom: '1px solid var(--border-subtle)',
})

const jobNodeName = css({
  fontSize: '0.875rem',
  fontWeight: 600,
  color: 'var(--text-primary)',
  marginBottom: 4,
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

const jobNodeStatus = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 5,
  fontSize: '0.6875rem',
  fontWeight: 600,
  letterSpacing: '0.02em',
  textTransform: 'uppercase',
})

const jobNodeSteps = css({
  padding: '8px 12px 10px',
  display: 'flex',
  flexDirection: 'column',
  gap: 4,
})

const stepDot = css({
  display: 'flex',
  alignItems: 'center',
  gap: 6,
  fontSize: '0.75rem',
  color: 'var(--text-secondary)',
  overflow: 'hidden',
})

const stepDotName = css({
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})



const logPanel = css({
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-lg)',
  overflow: 'hidden',
  animation: 'fade-in 0.2s var(--ease-out) both',
})

const logPanelHeader = css({
  display: 'flex',
  alignItems: 'center',
  gap: 12,
  padding: '12px 18px',
  borderBottom: '1px solid var(--border-subtle)',
  background: 'var(--bg-overlay)',
})

const logPanelTitle = css({
  fontSize: '0.875rem',
  fontWeight: 600,
  color: 'var(--text-primary)',
  flex: 1,
})

const stepList = css({
  borderRight: '1px solid var(--border-subtle)',
  width: 200,
  flexShrink: 0,
})

const stepListItem = css({
  display: 'flex',
  alignItems: 'center',
  gap: 8,
  padding: '10px 14px',
  cursor: 'pointer',
  fontSize: '0.8125rem',
  color: 'var(--text-secondary)',
  borderBottom: '1px solid var(--border-subtle)',
  transition: 'background 150ms',
  '&:last-child': { borderBottom: 'none' },
  '&:hover': { background: 'var(--bg-hover)' },
})

const stepListItemActive = css({
  background: 'var(--accent-muted)',
  color: 'var(--text-primary)',
})

const logBody = css({
  display: 'flex',
  minHeight: 300,
  maxHeight: 600,
  flex: 1,
  position: 'relative',
})

const logOutput = css({
  flex: 1,
  padding: '14px 18px',
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.8125rem',
  lineHeight: 1.7,
  color: 'var(--code-text)',
  background: 'var(--code-bg)',
  overflowY: 'auto',
  whiteSpace: 'pre-wrap',
  wordBreak: 'break-all',
  width: '100%',
})

const emptyLog = css({
  color: 'var(--text-disabled)',
  fontStyle: 'italic',
})

const btnCopyLog = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 5,
  padding: '3px 10px',
  borderRadius: 'var(--radius-md)',
  fontSize: '0.75rem',
  fontWeight: 500,
  border: '1px solid var(--border-default)',
  color: 'var(--text-secondary)',
  background: 'transparent',
  cursor: 'pointer',
  transition: 'border-color 120ms, color 120ms, background 120ms',
  '&:hover': {
    borderColor: 'var(--border-strong)',
    color: 'var(--text-primary)',
    background: 'var(--bg-hover)',
  },
  '&.copied': {
    borderColor: 'var(--success)',
    color: 'var(--success)',
    background: 'var(--success-muted)',
  },
})

const scrollBottomBtn = css({
  position: 'absolute',
  bottom: 12,
  right: 12,
  display: 'flex',
  alignItems: 'center',
  gap: 5,
  padding: '5px 12px',
  borderRadius: 'var(--radius-full)',
  fontSize: '0.75rem',
  fontWeight: 500,
  background: 'var(--bg-elevated)',
  border: '1px solid var(--border-default)',
  color: 'var(--text-secondary)',
  cursor: 'pointer',
  boxShadow: 'var(--shadow-md)',
  transition: 'border-color 120ms, color 120ms',
  animation: 'fade-in 0.15s var(--ease-out) both',
  '&:hover': {
    borderColor: 'var(--accent)',
    color: 'var(--accent)',
  },
})

const logOutputWrap = css({
  position: 'relative',
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  minHeight: 0,
})



const btnAnalyze = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 5,
  padding: '4px 11px',
  borderRadius: 'var(--radius-md)',
  fontSize: '0.8125rem',
  fontWeight: 500,
  border: '1px solid var(--ai-glow)',
  color: 'var(--ai)',
  background: 'var(--ai-muted)',
  cursor: 'pointer',
  transition: 'background 150ms, border-color 150ms',
  '&:hover:not(:disabled)': {
    background: 'var(--ai-muted)',
    borderColor: 'var(--ai-glow)',
  },
  '&:disabled': { opacity: 0.5, cursor: 'not-allowed' },
})

const analysisWrap = css({
  padding: '16px 18px',
  borderBottom: '1px solid var(--border-subtle)',
  background: 'linear-gradient(to bottom, var(--ai-muted), transparent)',
  animation: 'fade-in 0.2s var(--ease-out) both',
})

const thinkingWrap = css({
  position: 'relative',
  padding: '18px 20px',
  borderBottom: '1px solid var(--border-subtle)',
  background: 'linear-gradient(135deg, var(--ai-muted), transparent)',
  overflow: 'hidden',
  animation: 'fade-in 0.25s var(--ease-out) both',
  display: 'flex',
  alignItems: 'center',
  gap: 14,
})

const thinkingBar = css({
  position: 'absolute',
  top: 0,
  left: 0,
  width: '35%',
  height: 2,
  background: 'linear-gradient(90deg, transparent, var(--ai), transparent)',
  animation: 'ai-scan 1.8s ease-in-out infinite',
})

const thinkingLabel = css({
  flex: 1,
  fontSize: '0.875rem',
  fontWeight: 500,
  color: 'var(--ai)',
  display: 'flex',
  alignItems: 'center',
  gap: 8,
})

const thinkingDotsRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 5,
})

const dotBase = {
  width: 6,
  height: 6,
  borderRadius: '50%' as const,
  background: 'var(--ai)',
  display: 'inline-block' as const,
  animation: 'thinking-dot 1.3s ease-in-out infinite',
}
const thinkingDot1 = css({ ...dotBase })
const thinkingDot2 = css({ ...dotBase, animationDelay: '0.18s' })
const thinkingDot3 = css({ ...dotBase, animationDelay: '0.36s' })

const analysisHeader = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  marginBottom: 12,
})

const analysisTitle = css({
  display: 'flex',
  alignItems: 'center',
  gap: 7,
  fontSize: '0.875rem',
  fontWeight: 600,
  color: 'var(--ai)',
})

const analysisDismiss = css({
  padding: '2px 6px',
  borderRadius: 4,
  border: '1px solid var(--border-subtle)',
  background: 'transparent',
  color: 'var(--text-disabled)',
  cursor: 'pointer',
  fontSize: '0.75rem',
  '&:hover': { background: 'var(--bg-hover)', color: 'var(--text-secondary)' },
})

const analysisSummaryBox = css({
  padding: '10px 14px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--ai-muted)',
  border: '1px solid var(--ai-glow)',
  color: 'var(--text-primary)',
  fontSize: '0.875rem',
  lineHeight: 1.6,
  marginBottom: 12,
})

const analysisSectionLabel = css({
  fontSize: '0.6875rem',
  fontWeight: 600,
  color: 'var(--text-disabled)',
  letterSpacing: '0.07em',
  textTransform: 'uppercase',
  marginBottom: 4,
  marginTop: 10,
  '&:first-of-type': { marginTop: 0 },
})

const analysisText = css({
  color: 'var(--text-secondary)',
  fontSize: '0.875rem',
  lineHeight: 1.65,
  whiteSpace: 'pre-line',
})

const analysisCodeBlock = css({
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.8rem',
  color: 'var(--danger)',
  background: 'var(--code-bg)',
  padding: '8px 12px',
  borderRadius: 'var(--radius-md)',
  lineHeight: 1.65,
  overflowX: 'auto',
  maxHeight: 150,
  overflowY: 'auto',
})

const analysisErrorBox = css({
  padding: '10px 14px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--danger-muted)',
  border: '1px solid var(--danger-glow)',
  color: 'var(--danger)',
  fontSize: '0.875rem',
})



const approvalBanner = css({
  display: 'flex',
  alignItems: 'center',
  gap: 12,
  padding: '12px 18px',
  background: 'var(--warning-muted)',
  borderBottom: '1px solid var(--warning-glow)',
})

const approvalLabel = css({
  flex: 1,
  fontSize: '0.875rem',
  color: 'var(--warning)',
  fontWeight: 600,
  display: 'flex',
  alignItems: 'center',
  gap: 8,
})

const btnApprove = css({
  padding: '6px 16px',
  borderRadius: 'var(--radius-md)',
  fontSize: '0.8125rem',
  fontWeight: 600,
  border: 'none',
  cursor: 'pointer',
  background: 'var(--success)',
  color: 'var(--text-on-accent)',
  transition: 'opacity 150ms',
  '&:disabled': { opacity: 0.5, cursor: 'default' },
  '&:hover:not(:disabled)': { opacity: 0.85 },
})

const btnReject = css({
  padding: '6px 16px',
  borderRadius: 'var(--radius-md)',
  fontSize: '0.8125rem',
  fontWeight: 600,
  border: '1px solid var(--danger)',
  cursor: 'pointer',
  background: 'transparent',
  color: 'var(--danger)',
  transition: 'background 150ms',
  '&:disabled': { opacity: 0.5, cursor: 'default' },
  '&:hover:not(:disabled)': { background: 'var(--danger-muted)' },
})

const btnCancelRun = css({
  padding: '6px 14px',
  borderRadius: 'var(--radius-md)',
  fontSize: '0.8125rem',
  fontWeight: 600,
  border: '1px solid var(--border-default)',
  cursor: 'pointer',
  background: 'var(--bg-overlay)',
  color: 'var(--text-secondary)',
  display: 'flex',
  alignItems: 'center',
  gap: 6,
  transition: 'border-color 150ms, color 150ms',
  '&:disabled': { opacity: 0.5, cursor: 'default' },
  '&:hover:not(:disabled)': { borderColor: 'var(--danger)', color: 'var(--danger)' },
})




function stripAnsi(str: string): string {

  return str.replace(/\u001b\[[0-9;]*[a-zA-Z]/g, '')
}

function LogLines({ content }: { content: string }) {
  const lines = stripAnsi(content).split('\n')
  return (
    <>
      {lines.map((line, i) => (
        <div key={i} style={{ display: 'flex', gap: 0 }}>
          <span style={{
            userSelect: 'none',
            color: 'var(--log-line-number)',
            minWidth: 42,
            paddingRight: 12,
            textAlign: 'right',
            flexShrink: 0,
            fontVariantNumeric: 'tabular-nums',
          }}>
            {i + 1}
          </span>
          <span style={{ wordBreak: 'break-all', flex: 1 }}>{line}</span>
        </div>
      ))}
    </>
  )
}



const configPanel = css({
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-lg)',
  overflow: 'hidden',
  marginTop: 16,
})

const configPanelHeader = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  padding: '11px 18px',
  background: 'var(--bg-overlay)',
  cursor: 'pointer',
  userSelect: 'none',
  '&:hover': { background: 'var(--bg-hover)' },
})

const configPanelTitle = css({
  fontSize: '0.875rem',
  fontWeight: 600,
  color: 'var(--text-primary)',
  display: 'flex',
  alignItems: 'center',
  gap: 8,
})

const configChevron = (open: boolean) => css({
  color: 'var(--text-disabled)',
  transition: 'transform 200ms',
  transform: open ? 'rotate(180deg)' : 'none',
})

const configBody = css({
  padding: '14px 18px',
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.8125rem',
  lineHeight: 1.7,
  color: 'var(--text-primary)',
  background: 'var(--code-bg)',
  overflowX: 'auto',
  whiteSpace: 'pre',
  maxHeight: 420,
  overflowY: 'auto',
})



function RunStatusBadge({ status }: { status: string }) {
  const { t } = useI18n()
  const label = formatStatus(status, t.status)
  const dot = css({
    width: 7,
    height: 7,
    borderRadius: '50%',
    background: statusColor(status),
    animation: status === 'running' ? 'pulse-dot 1.5s ease-in-out infinite' : 'none',
  })
  const badge = css({
    display: 'inline-flex',
    alignItems: 'center',
    gap: 6,
    padding: '5px 12px',
    borderRadius: 'var(--radius-full)',
    fontSize: '0.8125rem',
    fontWeight: 600,
    background: statusBg(status),
    color: statusColor(status),
    border: `1px solid ${statusColor(status)}30`,
  })
  return <span css={badge}><span css={dot} />{label}</span>
}



function IconCheck() {
  return <svg width="10" height="10" viewBox="0 0 12 12" fill="none"><path d="M2 6.5l3 3 5-6" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" /></svg>
}
function IconX() {
  return <svg width="10" height="10" viewBox="0 0 12 12" fill="none"><path d="M3 3l6 6M9 3l-6 6" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" /></svg>
}
function IconClock() {
  return <svg width="10" height="10" viewBox="0 0 12 12" fill="none"><circle cx="6" cy="6" r="5" stroke="currentColor" strokeWidth="1.2" /><path d="M6 3.5V6l1.5 1.5" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" /></svg>
}
function IconSkip() {
  return <svg width="10" height="10" viewBox="0 0 12 12" fill="none"><path d="M2 3l5 3-5 3V3zm5 0l5 3-5 3V3z" fill="currentColor" /></svg>
}
function IconChevronDown() {
  return <svg width="14" height="14" viewBox="0 0 14 14" fill="none"><path d="M3 5l4 4 4-4" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" /></svg>
}
function IconConfig() {
  return <svg width="14" height="14" viewBox="0 0 16 16" fill="none"><rect x="2" y="3" width="12" height="10" rx="2" stroke="currentColor" strokeWidth="1.2" /><path d="M5 7h6M5 10h4" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" /></svg>
}
function IconBot() {
  return <svg width="14" height="14" viewBox="0 0 16 16" fill="none"><rect x="3" y="5" width="10" height="8" rx="2" stroke="currentColor" strokeWidth="1.2" /><circle cx="6" cy="9" r="1" fill="currentColor" /><circle cx="10" cy="9" r="1" fill="currentColor" /><path d="M6 3v2M10 3v2M1 9h2M13 9h2" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" /></svg>
}

function ThinkingPanel() {
  const { t } = useI18n()
  return (
    <div css={thinkingWrap}>
      <div css={thinkingBar} />
      <span css={thinkingLabel}>
        <span style={{ display: 'inline-flex', animation: 'spin 2.5s linear infinite' }}>
          <IconBot />
        </span>
        {t.run.aiAnalyzing}
      </span>
      <span css={thinkingDotsRow}>
        <span css={thinkingDot1} />
        <span css={thinkingDot2} />
        <span css={thinkingDot3} />
      </span>
    </div>
  )
}

function IconCopy() {
  return <svg width="12" height="12" viewBox="0 0 16 16" fill="none"><rect x="5" y="5" width="9" height="9" rx="1.5" stroke="currentColor" strokeWidth="1.2" /><path d="M11 5V3.5A1.5 1.5 0 009.5 2h-6A1.5 1.5 0 002 3.5v6A1.5 1.5 0 003.5 11H5" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" /></svg>
}
function IconCopyDone() {
  return <svg width="12" height="12" viewBox="0 0 16 16" fill="none"><path d="M2 8.5l3.5 3.5 8-8" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" /></svg>
}
function IconScrollDown() {
  return <svg width="11" height="11" viewBox="0 0 12 12" fill="none"><path d="M2 4l4 4 4-4" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" /></svg>
}

function IconApproval() {
  return <svg width="12" height="12" viewBox="0 0 16 16" fill="none"><path d="M8 2v5l3 3" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" /><circle cx="8" cy="8" r="6" stroke="currentColor" strokeWidth="1.3" /></svg>
}
function IconStopCircle() {
  return <svg width="12" height="12" viewBox="0 0 16 16" fill="none"><circle cx="8" cy="8" r="6" stroke="currentColor" strokeWidth="1.3" /><rect x="5.5" y="5.5" width="5" height="5" rx="1" fill="currentColor" /></svg>
}

function statusIcon(status: string) {
  switch (status) {
    case 'success':           return <IconCheck />
    case 'failed':            return <IconX />
    case 'skipped':           return <IconSkip />
    case 'cancelled':         return <IconX />
    case 'awaiting_approval': return <IconApproval />
    default:                  return <IconClock />
  }
}



function stepEventsToSteps(events: StepEvent[]): PipelineStep[] {
  return events.map((s) => ({
    index:       s.index,
    name:        s.name,
    status:      s.status as PipelineStep['status'],
    exit_code:   s.exit_code,
    log_output:  s.log_output,
    started_at:  s.started_at  ? new Date(s.started_at  * 1000).toISOString() : null,
    finished_at: s.finished_at ? new Date(s.finished_at * 1000).toISOString() : null,
  }))
}

const MAX_RETRIES = 5



export default function PipelineRunDetail() {
  const { id: projectId, runId } = useParams<{ id: string; runId: string }>()
  const token = useAppSelector((s) => s.auth.accessToken)
  const toast = useToast()
  const { t, lang } = useI18n()

  const [run, setRun]               = useState<PipelineRun | null>(null)
  const [loading, setLoading]       = useState(true)
  const { data: runData, isLoading: runQueryLoading } = useGetRunQuery(runId!, {
    skip: !runId || !token,
  })
  const [cancelRunMut] = useCancelRunMutation()
  const [approveJobMut] = useApproveJobMutation()
  const [rejectJobMut] = useRejectJobMutation()
  const [selectedJob, setSelectedJob]   = useState<string | null>(null)
  const [selectedStep, setSelectedStep] = useState<number>(0)
  const [configOpen, setConfigOpen] = useState(false)


  const [serverTimeMs, setServerTimeMs] = useState<number>(() => Date.now())


  const [analysis, setAnalysis]           = useState<AnalysisResult | null>(null)
  const [analyzing, setAnalyzing]         = useState(false)
  const [analysisError, setAnalysisError] = useState<string | null>(null)


  const [cancelling, setCancelling] = useState(false)
  const [approving, setApproving]   = useState<string | null>(null)


  const [logCopied, setLogCopied]     = useState(false)
  const [atLogBottom, setAtLogBottom] = useState(true)

  const wsRef             = useRef<WebSocket | null>(null)
  const logRef            = useRef<HTMLDivElement>(null)
  const retryCountRef     = useRef(0)
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const isLiveRef         = useRef(false)

  const runIdRef  = useRef(runId)
  const tokenRef  = useRef(token)
  runIdRef.current  = runId
  tokenRef.current  = token



  const cachedProjectName = useMemo(() => {
    if (!projectId) return null
    return getRecentProjects().find(p => p.id === projectId)?.name ?? null
  }, [projectId])



  useEffect(() => {
    const label = run ? t.run.runNumber.replace('{id}', run.id.slice(0, 8)) : t.run.runLabel
    document.title = `${label} — Flow`
    return () => { document.title = 'Flow — CI/CD' }
  }, [run?.id, t.run.runLabel, t.run.runNumber])



  const openWS = (currentRunId: string, currentToken: string) => {
    wsRef.current?.close()
    wsRef.current = connectRunWS(currentRunId, {
      token: currentToken,



      onOpen: () => { retryCountRef.current = 0 },

      onEvent: (event: WSEvent) => {
        if (event.type === 'heartbeat') {


          if (event.server_time_ms) setServerTimeMs(event.server_time_ms)
          return
        }

        if (event.type === 'job_updated') {
          if (!event.job_id) return
          setRun((prev) => {
            if (!prev) return prev
            return {
              ...prev,
              jobs: prev.jobs?.map((j) => {
                if (j.id !== event.job_id) return j
                return {
                  ...j,
                  status: event.status as JobStatus,


                  ...(event.job_started_at_ms
                    ? { started_at: new Date(event.job_started_at_ms).toISOString() }
                    : {}),
                  ...(event.job_finished_at_ms
                    ? { finished_at: new Date(event.job_finished_at_ms).toISOString() }
                    : {}),


                  ...(event.steps && event.steps.length > 0
                    ? { steps: stepEventsToSteps(event.steps) }
                    : {}),
                }
              }),
            }
          })
        } else if (event.type === 'job_awaiting_approval') {
          if (!event.job_id) return
          setRun((prev) => {
            if (!prev) return prev
            return {
              ...prev,
              jobs: prev.jobs?.map((j) =>
                j.id === event.job_id
                  ? { ...j, status: 'awaiting_approval' as JobStatus }
                  : j
              ),
            }
          })
        } else if (event.type === 'run_updated') {
          const terminal = ['success', 'failed', 'cancelled'].includes(event.status ?? '')
          if (terminal) {
            isLiveRef.current = false
            wsRef.current?.close()
            wsRef.current = null
          }
          setRun((prev) => {
            if (!prev) return prev
            return {
              ...prev,
              status: event.status as RunStatus,


              ...(terminal && event.run_finished_at_ms
                ? { finished_at: new Date(event.run_finished_at_ms).toISOString() }
                : {}),
            }
          })
        }
      },

      onClose: () => {
        wsRef.current = null
        if (!isLiveRef.current) return
        if (retryCountRef.current >= MAX_RETRIES) return
        retryCountRef.current += 1
        reconnectTimerRef.current = setTimeout(() => {
          const rid = runIdRef.current
          const tok = tokenRef.current
          if (rid && tok && isLiveRef.current) openWS(rid, tok)
        }, Math.min(1000 * 2 ** retryCountRef.current, 30_000))
      },
      onError: () => {
        wsRef.current = null
      },
    })
  }



  useEffect(() => {
    if (!runId || !token) return
    retryCountRef.current = 0
    return () => {
      isLiveRef.current = false
      if (reconnectTimerRef.current) {
        clearTimeout(reconnectTimerRef.current)
        reconnectTimerRef.current = null
      }
      wsRef.current?.close()
      wsRef.current = null
    }

  }, [runId, token])

  useEffect(() => {
    if (!runData || !runId || !token) return
    setRun(runData)
    setSelectedJob((prev) => prev ?? runData.jobs?.[0]?.id ?? null)
    const live = runData.status === 'pending' || runData.status === 'running'
    isLiveRef.current = live
    if (live) openWS(runId, token)
    setLoading(false)
  }, [runData, runId, token])

  useEffect(() => {
    if (!runId || runQueryLoading) return
    if (!runData) setLoading(false)
  }, [runId, runQueryLoading, runData])

  const isLive = run?.status === 'pending' || run?.status === 'running'



  useEffect(() => {
    setAnalysis(null)
    setAnalysisError(null)
  }, [selectedJob])



  const handleCancel = useCallback(async () => {
    if (!runId || !token || cancelling) return
    setCancelling(true)
    try {
      await cancelRunMut(runId).unwrap()
      toast.info(t.toast.cancellationRequested)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : t.errors.failedToCancelRun)
    } finally {
      setCancelling(false)
    }
  }, [runId, token, cancelling, toast, t])

  const handleApprove = useCallback(async (jobId: string) => {
    if (!runId || !token || approving) return
    setApproving(jobId)
    try {
      await approveJobMut({ runId, jobId }).unwrap()
      toast.success(t.toast.jobApproved)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : t.errors.failedToApproveJob)
    } finally {
      setApproving(null)
    }
  }, [runId, token, approving, toast, t])

  const handleReject = useCallback(async (jobId: string) => {
    if (!runId || !token || approving) return
    setApproving(jobId)
    try {
      await rejectJobMut({ runId, jobId }).unwrap()
      toast.warning(t.toast.jobRejected)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : t.errors.failedToRejectJob)
    } finally {
      setApproving(null)
    }
  }, [runId, token, approving, toast, t])



  const handleAnalyze = useCallback(async (jobId: string) => {
    if (!runId || !token || analyzing) return
    setAnalyzing(true)
    setAnalysis(null)
    setAnalysisError(null)
    try {
      const result = await analyzeFailure(runId, jobId, token, lang)
      setAnalysis(result)
    } catch (err) {
      const msg = err instanceof Error ? err.message : t.errors.analysisFailed
      setAnalysisError(msg)
      toast.error(msg)
    } finally {
      setAnalyzing(false)
    }
  }, [runId, token, analyzing, toast, t])



  useEffect(() => {
    const el = logRef.current
    if (!el) return
    if (atLogBottom) el.scrollTop = el.scrollHeight
  }, [selectedJob, selectedStep, run, atLogBottom])


  useEffect(() => {
    setAtLogBottom(true)
  }, [selectedJob, selectedStep])

  const handleLogScroll = useCallback(() => {
    const el = logRef.current
    if (!el) return
    const distFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight
    setAtLogBottom(distFromBottom < 40)
  }, [])

  const scrollLogToBottom = () => {
    if (!logRef.current) return
    logRef.current.scrollTo({ top: logRef.current.scrollHeight, behavior: 'smooth' })
    setAtLogBottom(true)
  }

  const displayJobs = run?.jobs ?? []



  const handleCopyLog = useCallback((logContent: string | null | undefined) => {
    if (!logContent) return
    navigator.clipboard.writeText(stripAnsi(logContent)).then(() => {
      setLogCopied(true)
      setTimeout(() => setLogCopied(false), 2200)
    })
  }, [])



  if (loading) {
    return <RunDetailSkeleton />
  }

  if (!run) {
    return (
      <>
        <Breadcrumb items={[{ label: t.nav.dashboard, to: '/' }]} />
        <p style={{ color: 'var(--danger)', marginBottom: 16 }}>{t.run.notFound}</p>
        <Link to={`/projects/${projectId}`} css={btnSecondary}>{t.run.backToProject}</Link>
      </>
    )
  }

  const duration   = durationMs(run, isLive ? serverTimeMs : undefined)
  const activeJob  = displayJobs.find((j) => j.id === selectedJob) ?? displayJobs[0] ?? null
  const activeStep = activeJob?.steps?.[selectedStep] ?? null
  const hasLogContent = !!(activeStep?.log_output)

  return (
    <div css={pageWrap}>
      <Breadcrumb
        items={[
          { label: t.nav.dashboard, to: '/' },
          { label: cachedProjectName ?? t.breadcrumb.project, to: `/projects/${projectId}` },
          { label: t.run.runNumber.replace('{id}', run.id.slice(0, 8)) },
        ]}
      />

      <div css={runHeader}>
        <div>
          <h1 css={runTitle}>
            {t.run.title}
            {isLive && (
              <span style={{ marginLeft: 12, fontSize: '0.75rem', color: 'var(--accent)', fontWeight: 500, animation: 'pulse-subtle 1.5s ease-in-out infinite' }}>
                ● {t.status.live}
              </span>
            )}
          </h1>
          <div css={runMeta}>
            <span css={metaChip}>
              <svg width="11" height="11" viewBox="0 0 16 16" fill="none"><circle cx="4" cy="4" r="2" stroke="currentColor" strokeWidth="1.2" /><circle cx="4" cy="12" r="2" stroke="currentColor" strokeWidth="1.2" /><circle cx="12" cy="4" r="2" stroke="currentColor" strokeWidth="1.2" /><path d="M4 6v4M4 6c0 2 8 3 8-2" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" /></svg>
              {run.branch}
            </span>
            <span css={metaChip} title={run.commit_sha}>
              {run.commit_sha.slice(0, 7)}
            </span>
            {duration !== null && (
              <span style={{ color: 'var(--text-tertiary)', fontVariantNumeric: 'tabular-nums' }}>
                {formatDuration(duration)}
              </span>
            )}
          </div>
        </div>
        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          {isLive && (
            <button
              type="button"
              css={btnCancelRun}
              disabled={cancelling}
              onClick={handleCancel}
              title={t.aria.abortRun}
            >
              <IconStopCircle />
              {cancelling ? t.run.cancelling : t.run.cancelRun}
            </button>
          )}
          <RunStatusBadge status={run.status} />
        </div>
      </div>


      {displayJobs.length > 0 && (
        <div css={pipelineGraph}>
          {displayJobs.map((job, idx) => (
            <div key={job.id} style={{ display: 'flex', alignItems: 'center' }}>
              {idx > 0 && (
                <div css={pipelineConnector}>
                  <svg width="20" height="12" viewBox="0 0 20 12" fill="none" aria-hidden>
                    <path d="M0 6h14M10 1l6 5-6 5" stroke="var(--border-strong)" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round" />
                  </svg>
                </div>
              )}
              <div
                css={[jobNode, selectedJob === job.id && jobNodeActive]}
                onClick={() => { setSelectedJob(job.id); setSelectedStep(0) }}
                role="button"
                tabIndex={0}
                onKeyDown={(e) => e.key === 'Enter' && setSelectedJob(job.id)}
                aria-pressed={selectedJob === job.id}
              >

                <div css={jobNodeColorBar(statusColor(job.status))} />

                <div css={jobNodeHeader}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 6, marginBottom: 5 }}>
                    <div css={jobNodeName}>{job.display_name || job.name}</div>
                    {(job.attempt ?? 1) > 1 && (
                      <span style={{ fontSize: '0.625rem', fontWeight: 700, padding: '1px 6px', borderRadius: 'var(--radius-full)', background: 'var(--warning-muted)', color: 'var(--warning)', fontFamily: 'ui-monospace, monospace', flexShrink: 0 }}>
                        ×{job.attempt ?? 1}
                      </span>
                    )}
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 6 }}>
                    <div css={[jobNodeStatus, { color: statusColor(job.status) }]}>
                      <span css={{
                        width: 6, height: 6, borderRadius: '50%',
                        background: statusColor(job.status),
                        animation: job.status === 'running' ? 'pulse-dot 1.5s ease-in-out infinite' : 'none',
                        display: 'inline-block',
                      }} />
                      {formatStatus(job.status, t.status)}
                    </div>
                    {job.started_at && (
                      <span style={{ fontSize: '0.6rem', color: 'var(--text-disabled)', fontVariantNumeric: 'tabular-nums', fontFamily: 'ui-monospace, monospace' }}>
                        {formatDuration(
                          durationMs(
                            { started_at: job.started_at, finished_at: job.finished_at, status: job.status } as PipelineRun,
                            job.status === 'running' ? serverTimeMs : undefined
                          ) ?? 0
                        )}
                      </span>
                    )}
                  </div>
                </div>
                <div css={jobNodeSteps}>
                  {job.steps.slice(0, 4).map((step) => (
                    <div key={step.index} css={stepDot}>
                      <span style={{ color: statusColor(step.status), flexShrink: 0 }}>
                        {statusIcon(step.status)}
                      </span>
                      <span css={stepDotName}>{step.name}</span>
                    </div>
                  ))}
                  {job.steps.length > 4 && (
                    <div css={[stepDot, { color: 'var(--text-disabled)' }]}>
                      {t.run.moreSteps.replace('{count}', String(job.steps.length - 4))}
                    </div>
                  )}
                </div>
              </div>
            </div>
          ))}
        </div>
      )}


      {activeJob && (
        <div css={logPanel}>
          <div css={logPanelHeader}>
            <span css={logPanelTitle}>
              {activeJob.display_name || activeJob.name}
              {activeJob.status === 'running' && (
                <span style={{ marginLeft: 10, fontSize: '0.75rem', color: 'var(--accent)', fontWeight: 400 }}>{t.run.running}</span>
              )}
              {(activeJob.attempt ?? 1) > 1 && (
                <span style={{
                  marginLeft: 10,
                  fontSize: '0.6875rem',
                  fontWeight: 700,
                  padding: '2px 8px',
                  borderRadius: 'var(--radius-full)',
                  background: 'var(--warning-muted)',
                  color: 'var(--warning)',
                  fontFamily: 'ui-monospace, monospace',
                }}>
                  {t.run.attempt} {activeJob.attempt ?? 1}
                </span>
              )}
            </span>
            <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
              {hasLogContent && (
                <button
                  type="button"
                  css={[btnCopyLog, logCopied ? { borderColor: 'var(--success)', color: 'var(--success)', background: 'var(--success-muted)' } : {}]}
                  onClick={() => handleCopyLog(activeStep?.log_output)}
                  title={t.aria.copyLog}
                >
                  {logCopied ? <IconCopyDone /> : <IconCopy />}
                  {logCopied ? t.run.copied : t.run.copyLog}
                </button>
              )}
              {activeJob.status === 'failed' && (
                <button
                  type="button"
                  css={btnAnalyze}
                  disabled={analyzing}
                  onClick={() => handleAnalyze(activeJob.id)}
                  title={t.aria.analyzeFailure}
                >
                  <span style={analyzing ? { display: 'inline-flex', animation: 'spin 2.5s linear infinite' } : {}}>
                    <IconBot />
                  </span>
                  {analyzing ? t.run.analyzing : t.run.analyze}
                </button>
              )}
              <span style={{ fontSize: '0.75rem', color: statusColor(activeJob.status), fontWeight: 600 }}>
                {activeJob.status === 'awaiting_approval' ? t.status.awaitingApproval : formatStatus(activeJob.status, t.status)}
              </span>
            </span>
          </div>


          {activeJob.status === 'awaiting_approval' && (
            <div css={approvalBanner}>
              <span css={approvalLabel}>
                <IconApproval />
                {t.run.approvalRequired}
              </span>
              <button
                type="button"
                css={btnReject}
                disabled={approving === activeJob.id}
                onClick={() => handleReject(activeJob.id)}
              >
                {t.run.reject}
              </button>
              <button
                type="button"
                css={btnApprove}
                disabled={approving === activeJob.id}
                onClick={() => handleApprove(activeJob.id)}
              >
                {approving === activeJob.id ? t.run.approving : t.run.approve}
              </button>
            </div>
          )}


          {analyzing && <ThinkingPanel />}


          {!analyzing && (analysis || analysisError) && (
            <div css={analysisWrap}>
              <div css={analysisHeader}>
                <span css={analysisTitle}>
                  <IconBot />
                  {t.run.aiAnalysis}
                </span>
                <button
                  type="button"
                  css={analysisDismiss}
                  onClick={() => { setAnalysis(null); setAnalysisError(null) }}
                >
                  {t.run.dismiss}
                </button>
              </div>

              {analysisError && (
                <div css={analysisErrorBox}>{analysisError}</div>
              )}

              {analysis && (
                <>
                  <div css={analysisSummaryBox}>{analysis.summary}</div>

                  <div css={analysisSectionLabel}>{t.run.rootCause}</div>
                  <p css={analysisText}>{analysis.root_cause}</p>

                  <div css={analysisSectionLabel}>{t.run.howToFix}</div>
                  <p css={analysisText}>{analysis.fix}</p>

                  {analysis.relevant_lines?.length > 0 && (
                    <>
                      <div css={analysisSectionLabel}>{t.run.relevantLines}</div>
                      <pre css={analysisCodeBlock}>{analysis.relevant_lines.join('\n')}</pre>
                    </>
                  )}
                </>
              )}
            </div>
          )}

          <div css={logBody}>

            {activeJob.steps.length > 0 && (
              <div css={stepList}>
                {activeJob.steps.map((step, i) => (
                  <div
                    key={step.index}
                    css={[stepListItem, selectedStep === i && stepListItemActive]}
                    onClick={() => setSelectedStep(i)}
                    role="button"
                    tabIndex={0}
                    onKeyDown={(e) => e.key === 'Enter' && setSelectedStep(i)}
                  >
                    <span style={{ color: statusColor(step.status), flexShrink: 0 }}>
                      {statusIcon(step.status)}
                    </span>
                    <span style={{ overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap', fontSize: '0.8125rem' }}>
                      {step.name}
                    </span>
                  </div>
                ))}
              </div>
            )}


            <div css={logOutputWrap}>
              <div css={logOutput} ref={logRef} onScroll={handleLogScroll}>
                {activeStep ? (
                  activeStep.log_output ? (
                    <LogLines content={activeStep.log_output} />
                  ) : (
                    <span css={emptyLog}>{t.run.noOutput}</span>
                  )
                ) : activeJob.steps.length === 0 ? (
                  <span css={emptyLog}>{t.run.noSteps}</span>
                ) : (
                  <span css={emptyLog}>{t.run.selectStep}</span>
                )}
              </div>
              {!atLogBottom && hasLogContent && (
                <button
                  type="button"
                  css={scrollBottomBtn}
                  onClick={scrollLogToBottom}
                  title={t.aria.scrollToBottom}
                >
                  <IconScrollDown />
                  {t.run.scrollToBottom}
                </button>
              )}
            </div>
          </div>
        </div>
      )}

      {!run.jobs?.length && (
        <div style={{ padding: '40px 0', textAlign: 'center', color: 'var(--text-tertiary)', fontSize: '0.9375rem' }}>
          {t.run.noJobs}
        </div>
      )}


      {run.pipeline_yaml && (
        <div css={configPanel}>
          <div
            css={configPanelHeader}
            onClick={() => setConfigOpen((o) => !o)}
            role="button"
            tabIndex={0}
            onKeyDown={(e) => e.key === 'Enter' && setConfigOpen((o) => !o)}
            aria-expanded={configOpen}
          >
            <span css={configPanelTitle}>
              <IconConfig />
              {t.run.pipelineConfig}
              <span style={{ fontWeight: 400, fontSize: '0.75rem', color: 'var(--text-disabled)', fontFamily: 'ui-monospace, monospace' }}>
                {t.run.pipelineConfigSub}
              </span>
            </span>
            <span css={configChevron(configOpen)}><IconChevronDown /></span>
          </div>
          {configOpen && (
            <div css={configBody}>{run.pipeline_yaml}</div>
          )}
        </div>
      )}
    </div>
  )
}
