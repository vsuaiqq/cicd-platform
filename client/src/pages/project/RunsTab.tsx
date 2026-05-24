import { Link } from 'react-router-dom'
import type { PipelineRun } from '../../api/pipelines'
import { durationMs, formatDuration, timeAgo } from '../../api/pipelines'
import { EmptyState, RunStatusBadge, RunsTableSkeleton } from '../../components/ui'
import { btnPrimary } from '../../styles/theme'
import { useI18n } from '../../i18n'
import { IconPipeline, IconSetup } from './icons'
import { formatTriggerType } from './utils'
import {
  btnGhost,
  emptyPanelBody,
  panel,
  panelHeader,
  panelHeaderActions,
  panelMeta,
  panelTitle,
  runBranch,
  runCommit,
  runDuration,
  runRow,
  runsTableHead,
  runTime,
} from './styles'

interface RunsTabProps {
  projectId: string
  runs: PipelineRun[]
  loading: boolean
}

export function RunsTab({ projectId, runs, loading }: RunsTabProps) {
  const { t, lang } = useI18n()

  return (
    <section css={panel} aria-labelledby="runs-panel-title">
      <div css={panelHeader}>
        <h2 id="runs-panel-title" css={panelTitle}>
          <IconPipeline />
          {t.project.runsTitle}
        </h2>
        {!loading && runs.length > 0 && (
          <div css={panelHeaderActions}>
            <span css={panelMeta}>{t.project.runsCount.replace('{n}', String(runs.length))}</span>
            <Link to={`/projects/${projectId}/setup`} css={btnGhost}>
              <IconSetup />
              {t.project.setup}
            </Link>
          </div>
        )}
      </div>

      {loading ? (
        <RunsTableSkeleton />
      ) : runs.length === 0 ? (
        <div css={emptyPanelBody}>
          <EmptyState
            title={t.project.noRunsTitle}
            description={t.project.noRunsDesc}
            icon={<IconPipeline />}
            action={
              <Link to={`/projects/${projectId}/setup`} css={btnPrimary}>
                <IconSetup />
                {t.project.setupCardBtn}
              </Link>
            }
          />
        </div>
      ) : (
        <>
          <div css={runsTableHead} role="row">
            <span>{t.project.colStatus}</span>
            <span>{t.project.colRun}</span>
            <span style={{ textAlign: 'right' }}>{t.project.colDuration}</span>
            <span style={{ textAlign: 'right' }}>{t.project.colWhen}</span>
            <span />
          </div>
          {runs.map((run) => {
            const dur = durationMs(run)
            const triggerLabel = run.trigger_type
              ? formatTriggerType(run.trigger_type, t.project)
              : null
            return (
              <Link
                key={run.id}
                to={`/projects/${projectId}/runs/${run.id}`}
                css={runRow}
                role="row"
              >
                <div>
                  <RunStatusBadge status={run.status} />
                </div>
                <div>
                  <div css={runBranch}>{run.branch}</div>
                  <div css={runCommit} title={run.commit_sha}>
                    <span>{run.commit_sha.slice(0, 7)}</span>
                    {triggerLabel && (
                      <span style={{ opacity: 0.85 }}> · {triggerLabel}</span>
                    )}
                  </div>
                </div>
                <div css={runDuration}>{dur !== null ? formatDuration(dur) : '—'}</div>
                <div css={runTime}>{timeAgo(run.created_at, lang)}</div>
                <span />
              </Link>
            )
          })}
        </>
      )}
    </section>
  )
}
