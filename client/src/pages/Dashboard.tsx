import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { css } from '@emotion/react'
import type { ProjectSummary } from '../api/projects'
import { useAppSelector } from '../store'
import { useDeleteProjectMutation, useListProjectsQuery } from '../store/api/apiSlice'
import { DashboardPageSkeleton, StatusBadge, EmptyState, ConfirmDialog, useToast } from '../components/ui'
import { useI18n } from '../i18n'
import {
  pageTitle,
  pageDesc,
  btnPrimary,
  statsRow,
  statCard,
  statLabel,
  statValue,
  statDesc,
  projectCard,
  projectCardName,
  projectCardMeta,
  branchChip,
  projectCardFooter,
  toolbar,
} from '../styles/theme'

const pageHeader = css({
  marginBottom: 32,
  animation: 'fade-in 0.25s var(--ease-out) both',
})

const grid = css({
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fill, minmax(320px, 1fr))',
  gap: 14,
})

const cardTopRow = css({
  display: 'flex',
  alignItems: 'flex-start',
  justifyContent: 'space-between',
  gap: 8,
})

const statAccentLine = css({
  height: 2,
  borderRadius: 2,
  marginTop: 'auto',
  marginBottom: 2,
})

const deleteBtn = css({
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  width: 28,
  height: 28,
  borderRadius: 'var(--radius-md)',
  border: 'none',
  background: 'transparent',
  color: 'var(--text-disabled)',
  cursor: 'pointer',
  transition: 'color 120ms, background 120ms',
  flexShrink: 0,
  '&:hover': {
    color: 'var(--danger)',
    background: 'var(--danger-muted)',
  },
})



function IconTrash() {
  return (
    <svg width="13" height="13" viewBox="0 0 16 16" fill="none" aria-hidden>
      <path d="M2 4h12" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
      <path d="M5 4V2.5a.5.5 0 01.5-.5h5a.5.5 0 01.5.5V4" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
      <path d="M3 4l.8 9a1 1 0 001 .9h6.4a1 1 0 001-.9L13 4" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M6.5 7v4M9.5 7v4" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" />
    </svg>
  )
}

function IconGitBranch() {
  return (
    <svg width="12" height="12" viewBox="0 0 16 16" fill="none" aria-hidden>
      <circle cx="4" cy="4" r="2" stroke="currentColor" strokeWidth="1.3" />
      <circle cx="4" cy="12" r="2" stroke="currentColor" strokeWidth="1.3" />
      <circle cx="12" cy="4" r="2" stroke="currentColor" strokeWidth="1.3" />
      <path d="M4 6v4M4 6c0 2 8 3 8-2" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  )
}

function IconRepo() {
  return (
    <svg width="12" height="12" viewBox="0 0 16 16" fill="none" aria-hidden>
      <rect x="2" y="2" width="12" height="12" rx="2" stroke="currentColor" strokeWidth="1.3" />
      <path d="M6 5v6M10 8H6" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  )
}

function StatIcon({ type }: { type: 'total' | 'active' | 'pending' }) {
  if (type === 'total') return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden>
      <rect x="1" y="1" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
      <rect x="9" y="1" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
      <rect x="1" y="9" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
      <rect x="9" y="9" width="6" height="6" rx="1.5" stroke="currentColor" strokeWidth="1.25" />
    </svg>
  )
  if (type === 'active') return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden>
      <circle cx="8" cy="8" r="6" stroke="currentColor" strokeWidth="1.25" />
      <path d="M5.5 8.5l2 2 3-4" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  )
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" aria-hidden>
      <circle cx="8" cy="8" r="6" stroke="currentColor" strokeWidth="1.25" />
      <path d="M8 5v4M8 11v.5" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
    </svg>
  )
}

const IconPlus = () => (
  <svg width="14" height="14" viewBox="0 0 16 16" fill="none" aria-hidden>
    <path d="M8 3v10M3 8h10" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" />
  </svg>
)



export default function Dashboard() {
  const { accessToken, sessionValidated } = useAppSelector((s) => s.auth)
  const toast = useToast()
  const { t } = useI18n()
  const { data: projects = [], isLoading, isFetching } = useListProjectsQuery(undefined, {
    skip: !accessToken || !sessionValidated,
  })
  const [deleteProject, { isLoading: deleting }] = useDeleteProjectMutation()
  const [deleteTarget, setDeleteTarget] = useState<ProjectSummary | null>(null)
  const loading = isLoading || (isFetching && projects.length === 0)
  useEffect(() => {
    document.title = `${t.dashboard.title} — Flow`
    return () => { document.title = 'Flow — CI/CD' }
  }, [t])

  const handleDeleteConfirm = async () => {
    if (!deleteTarget) return
    try {
      await deleteProject(deleteTarget.id).unwrap()
      toast.success(t.toast.projectDeleted.replace('{name}', deleteTarget.name))
      setDeleteTarget(null)
    } catch (err) {
      toast.error(err instanceof Error ? err.message : t.errors.failedToDeleteProject)
    }
  }

  const activeCount  = projects.filter((p) => p.status?.toLowerCase() === 'active').length
  const pendingCount = projects.filter((p) => p.status?.toLowerCase() === 'pending').length

  return (
    <>
      <header css={pageHeader}>
        <h1 css={pageTitle}>{t.dashboard.title}</h1>
        <p css={pageDesc}>{t.dashboard.desc}</p>
      </header>

      {loading ? (
        <DashboardPageSkeleton />
      ) : projects.length > 0 ? (
        <>

          <div css={statsRow}>
            <div css={[statCard, { borderTop: '2px solid var(--accent)' }]}>
              <div css={statLabel} style={{ display: 'flex', alignItems: 'center', gap: 6, color: 'var(--text-tertiary)' }}>
                <StatIcon type="total" />
                {t.dashboard.total}
              </div>
              <div css={statValue}>{projects.length}</div>
              <div css={statDesc}>{t.dashboard.totalDesc}</div>
              <div css={[statAccentLine, { background: 'var(--accent-muted)', width: '100%' }]} />
            </div>
            <div css={[statCard, activeCount > 0 ? { borderTop: '2px solid var(--success)' } : {}]}>
              <div css={statLabel} style={{ display: 'flex', alignItems: 'center', gap: 6, color: 'var(--success)' }}>
                <StatIcon type="active" />
                {t.dashboard.active}
              </div>
              <div css={[statValue, { color: activeCount > 0 ? 'var(--success)' : 'var(--text-primary)' }]}>
                {activeCount}
              </div>
              <div css={statDesc}>{t.dashboard.activeDesc}</div>
              <div css={[statAccentLine, { background: 'var(--success-muted)', width: `${Math.min(100, (activeCount / Math.max(projects.length, 1)) * 100)}%` }]} />
            </div>
            <div css={[statCard, pendingCount > 0 ? { borderTop: '2px solid var(--warning)' } : {}]}>
              <div css={statLabel} style={{ display: 'flex', alignItems: 'center', gap: 6, color: 'var(--warning)' }}>
                <StatIcon type="pending" />
                {t.dashboard.pending}
              </div>
              <div css={[statValue, { color: pendingCount > 0 ? 'var(--warning)' : 'var(--text-primary)' }]}>
                {pendingCount}
              </div>
              <div css={statDesc}>{t.dashboard.pendingDesc}</div>
              <div css={[statAccentLine, { background: 'var(--warning-muted)', width: `${Math.min(100, (pendingCount / Math.max(projects.length, 1)) * 100)}%` }]} />
            </div>
          </div>


          <div css={toolbar}>
            <span style={{ fontSize: '0.75rem', fontWeight: 700, color: 'var(--text-tertiary)', textTransform: 'uppercase', letterSpacing: '0.09em' }}>
              {t.dashboard.projectsLabel} · {projects.length}
            </span>
            <Link to="/projects/new" css={btnPrimary}>
              <IconPlus />
              {t.dashboard.newProject}
            </Link>
          </div>

          <div css={grid}>
            {projects.map((p, i) => (
              <div key={p.id} css={projectCard} style={{ animationDelay: `${i * 50}ms` }}>
                <div css={cardTopRow}>
                  <Link to={`/projects/${p.id}`} css={projectCardName}>
                    {p.name}
                  </Link>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 6, flexShrink: 0 }}>
                    <StatusBadge status={p.status} />
                    <button
                      type="button"
                      css={deleteBtn}
                      onClick={(e) => { e.preventDefault(); setDeleteTarget(p) }}
                      aria-label={`${t.dashboard.deleteProject} ${p.name}`}
                      title={t.dashboard.deleteProject}
                    >
                      <IconTrash />
                    </button>
                  </div>
                </div>

                <div css={projectCardMeta}>
                  <IconRepo />
                  <span title={p.repo_url}>{p.repo_url}</span>
                </div>

                <div css={projectCardFooter}>
                  <div css={branchChip}>
                    <IconGitBranch />
                    {p.default_branch}
                  </div>
                </div>
              </div>
            ))}
          </div>
        </>
      ) : (
        <>
          <div css={toolbar}>
            <span />
            <Link to="/projects/new" css={btnPrimary}>
              <IconPlus />
              {t.dashboard.newProject}
            </Link>
          </div>
          <EmptyState
            title={t.dashboard.noProjectsTitle}
            description={t.dashboard.noProjectsDesc}
            action={
              <Link to="/projects/new" css={btnPrimary}>
                {t.dashboard.createFirst}
              </Link>
            }
          />
        </>
      )}


      <ConfirmDialog
        open={deleteTarget !== null}
        title={t.dashboard.confirmDeleteTitle.replace('{name}', deleteTarget?.name ?? '')}
        description={t.dashboard.confirmDeleteDesc}
        confirmLabel={t.dashboard.deleteProject}
        loading={deleting}
        onConfirm={handleDeleteConfirm}
        onCancel={() => setDeleteTarget(null)}
      />
    </>
  )
}
