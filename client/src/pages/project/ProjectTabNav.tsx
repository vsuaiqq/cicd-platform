import { useI18n } from '../../i18n'
import { IconAnalytics, IconRuns, IconSettings } from './icons'
import { tabBtn, tabCount, tabList, tabShell } from './styles'

export type ProjectTabId = 'runs' | 'analytics' | 'settings'

const TAB_PANELS: Record<ProjectTabId, string> = {
  runs: 'project-panel-runs',
  analytics: 'project-panel-analytics',
  settings: 'project-panel-settings',
}

interface ProjectTabNavProps {
  active: ProjectTabId
  runsCount: number | null
  onChange: (tab: ProjectTabId) => void
}

export function ProjectTabNav({ active, runsCount, onChange }: ProjectTabNavProps) {
  const { t } = useI18n()

  return (
    <nav css={tabShell} aria-label={t.project.tabsLabel}>
      <div css={tabList} role="tablist">
        <button
          type="button"
          role="tab"
          id="project-tab-runs"
          aria-selected={active === 'runs'}
          aria-controls={TAB_PANELS.runs}
          css={tabBtn(active === 'runs')}
          onClick={() => onChange('runs')}
        >
          <IconRuns />
          {t.project.tabRuns}
          {runsCount !== null && runsCount > 0 && <span css={tabCount}>{runsCount}</span>}
        </button>
        <button
          type="button"
          role="tab"
          id="project-tab-analytics"
          aria-selected={active === 'analytics'}
          aria-controls={TAB_PANELS.analytics}
          css={tabBtn(active === 'analytics')}
          onClick={() => onChange('analytics')}
        >
          <IconAnalytics />
          {t.project.tabAnalytics}
        </button>
        <button
          type="button"
          role="tab"
          id="project-tab-settings"
          aria-selected={active === 'settings'}
          aria-controls={TAB_PANELS.settings}
          css={tabBtn(active === 'settings')}
          onClick={() => onChange('settings')}
        >
          <IconSettings />
          {t.project.tabSettings}
        </button>
      </div>
    </nav>
  )
}
