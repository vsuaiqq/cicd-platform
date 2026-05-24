import { useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import type { Project } from '../../api/projects'
import { useAppSelector } from '../../store'
import {
  useListMembersQuery,
  useListRunsQuery,
} from '../../store/api/apiSlice'
import { Breadcrumb } from '../../components/ui'
import { useI18n } from '../../i18n'
import {
  EnvVarsEditor,
  MembersPanel,
  PipelineYAMLEditor,
  ProjectAnalytics,
  ProjectSettingsEditor,
  SecretsEditor,
} from '../ProjectDetail'
import { ProjectHero } from './ProjectHero'
import { ProjectTabNav, type ProjectTabId } from './ProjectTabNav'
import { RunsTab } from './RunsTab'
import { SettingsTab, type SettingsSectionId } from './SettingsTab'
import { IconSetup } from './icons'
import { pageWrap, projectNotice, projectNoticeLink, tabPanel } from './styles'
import { sortRunsByNewest } from './utils'

interface ProjectViewProps {
  project: Project
  onDelete: () => void
  onUpdated: (p: Project) => void
  deleting: boolean
}

export function ProjectView({ project, onDelete, onUpdated, deleting }: ProjectViewProps) {
  const { sessionValidated } = useAppSelector((s) => s.auth)
  const currentUserId = useAppSelector((s) => s.auth.userId) ?? ''
  const { t } = useI18n()

  const [activeTab, setActiveTab] = useState<ProjectTabId>('runs')
  const [settingsSection, setSettingsSection] = useState<SettingsSectionId>('general')

  const isOwner = project.user_id === currentUserId

  const { data: runsRaw = [], isLoading: runsLoading, isFetching: runsFetching } = useListRunsQuery(
    project.id,
    { skip: !sessionValidated },
  )

  const runs = useMemo(() => sortRunsByNewest(runsRaw), [runsRaw])

  const { data: membersData, isLoading: membersLoading } = useListMembersQuery(project.id, {
    skip: !sessionValidated || isOwner,
  })

  const userRole: 'owner' | 'editor' | 'viewer' | 'loading' = isOwner
    ? 'owner'
    : membersLoading
      ? 'loading'
      : (membersData?.requester_role ?? 'viewer')

  const runsBusy = runsLoading || (runsFetching && runs.length === 0)

  const settingsReady = sessionValidated

  return (
    <div css={pageWrap}>
      <Breadcrumb items={[{ label: t.nav.dashboard, to: '/' }, { label: project.name }]} />

      <ProjectHero
        project={project}
        userRole={userRole}
        deleting={deleting}
        onDelete={onDelete}
      />

      {project.status?.toLowerCase() === 'pending' && (
        <div css={projectNotice} role="status">
          <span>{t.project.pendingNotice}</span>
          <Link to={`/projects/${project.id}/setup`} css={projectNoticeLink}>
            <IconSetup />
            {t.project.setupCardBtn}
          </Link>
        </div>
      )}

      <ProjectTabNav
        active={activeTab}
        runsCount={runsBusy ? null : runs.length}
        onChange={setActiveTab}
      />

      {activeTab === 'runs' && (
        <div
          id="project-panel-runs"
          role="tabpanel"
          aria-labelledby="project-tab-runs"
          css={tabPanel}
        >
          <RunsTab projectId={project.id} runs={runs} loading={runsBusy} />
        </div>
      )}

      {activeTab === 'analytics' && settingsReady && (
        <div
          id="project-panel-analytics"
          role="tabpanel"
          aria-labelledby="project-tab-analytics"
          css={tabPanel}
        >
          <ProjectAnalytics projectId={project.id} />
        </div>
      )}

      {activeTab === 'settings' && settingsReady && (
        <div
          id="project-panel-settings"
          role="tabpanel"
          aria-labelledby="project-tab-settings"
          css={tabPanel}
        >
          <SettingsTab active={settingsSection} onChange={setSettingsSection}>
            {settingsSection === 'general' && (
              <ProjectSettingsEditor
                project={project}
                onUpdated={onUpdated}
                readOnly={userRole !== 'owner'}
              />
            )}
            {settingsSection === 'env' && (
              <EnvVarsEditor projectId={project.id} readOnly={userRole === 'viewer'} />
            )}
            {settingsSection === 'secrets' && (
              <SecretsEditor projectId={project.id} readOnly={userRole === 'viewer'} />
            )}
            {settingsSection === 'team' && (
              <MembersPanel projectId={project.id} currentUserId={currentUserId} />
            )}
            {settingsSection === 'pipeline' && (
              <PipelineYAMLEditor projectId={project.id} readOnly={userRole === 'viewer'} />
            )}
          </SettingsTab>
        </div>
      )}
    </div>
  )
}
