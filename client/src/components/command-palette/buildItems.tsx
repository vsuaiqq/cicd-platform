import type { ReactNode } from 'react'
import type { ProjectSummary } from '../../api/projects'
import { bestScore } from '../../lib/commandPaletteMatch'
import { projectSearchScore, shortRepo } from './utils'
import {
  IconClock,
  IconDashboard,
  IconPlus,
  IconProject,
  IconSettings,
  IconUser,
} from './icons'

export type IconTone = 'accent' | 'muted' | 'default'

export interface PaletteItem {
  id: string
  label: string
  sublabel?: string
  icon: ReactNode
  tone: IconTone
  onSelect: () => void
}

export type ListEntry =
  | { type: 'group'; id: string; label: string }
  | { type: 'item'; item: PaletteItem }

export interface PaletteCopy {
  navigate: string
  actions: string
  recent: string
  allProjects: string
  projects: string
  dashboardLabel: string
  dashboardSub: string
  newProjectLabel: string
  newProjectSub: string
  profileLabel: string
  profileSub: string
  settingsLabel: string
  settingsSub: string
}

interface BuildParams {
  projects: ProjectSummary[]
  recentIds: string[]
  query: string
  copy: PaletteCopy
  go: (path: string) => void
  onClose: () => void
  onOpenSettings?: () => void
}

export function buildPaletteItems({
  projects,
  recentIds,
  query,
  copy,
  go,
  onClose,
  onOpenSettings,
}: BuildParams): { entries: ListEntry[]; selectable: PaletteItem[] } {
  const q = query.trim()
  const recentSet = new Set(recentIds)
  const byId = new Map(projects.map((p) => [p.id, p]))

  const actionDefs: PaletteItem[] = [
    {
      id: 'action-dashboard',
      label: copy.dashboardLabel,
      sublabel: copy.dashboardSub,
      icon: <IconDashboard />,
      tone: 'accent',
      onSelect: () => go('/'),
    },
    {
      id: 'action-new',
      label: copy.newProjectLabel,
      sublabel: copy.newProjectSub,
      icon: <IconPlus />,
      tone: 'accent',
      onSelect: () => go('/projects/new'),
    },
    {
      id: 'action-profile',
      label: copy.profileLabel,
      sublabel: copy.profileSub,
      icon: <IconUser />,
      tone: 'accent',
      onSelect: () => go('/profile'),
    },
    {
      id: 'action-settings',
      label: copy.settingsLabel,
      sublabel: copy.settingsSub,
      icon: <IconSettings />,
      tone: 'accent',
      onSelect: () => {
        onClose()
        onOpenSettings?.()
      },
    },
  ]

  const actions = actionDefs.filter((a) => !q || bestScore([a.label, a.sublabel ?? ''], q) > 0)

  const toProjectItem = (p: ProjectSummary, tone: IconTone): PaletteItem => ({
    id: `project-${p.id}`,
    label: p.name,
    sublabel: shortRepo(p.repo_url),
    icon: tone === 'muted' ? <IconClock /> : <IconProject />,
    tone,
    onSelect: () => go(`/projects/${p.id}`),
  })

  let projectItems: PaletteItem[] = []

  if (!q) {
    const recentProjects = recentIds
      .map((id) => byId.get(id))
      .filter((p): p is ProjectSummary => p != null)

    const recentIdSet = new Set(recentProjects.map((p) => p.id))
    const rest = projects
      .filter((p) => !recentIdSet.has(p.id))
      .sort((a, b) => a.name.localeCompare(b.name))

    projectItems = [
      ...recentProjects.map((p) => toProjectItem(p, 'muted')),
      ...rest.map((p) => toProjectItem(p, 'default')),
    ]
  } else {
    projectItems = projects
      .map((p) => ({
        p,
        score: projectSearchScore(p, q, recentSet.has(p.id)),
      }))
      .filter(({ score }) => score > 0)
      .sort((a, b) => b.score - a.score)
      .slice(0, 12)
      .map(({ p }) => toProjectItem(p, recentSet.has(p.id) ? 'muted' : 'default'))
  }

  const entries: ListEntry[] = []
  const selectable: PaletteItem[] = []

  const pushGroup = (label: string, items: PaletteItem[]) => {
    if (!items.length) return
    entries.push({ type: 'group', id: `group-${label}`, label })
    for (const item of items) {
      entries.push({ type: 'item', item })
      selectable.push(item)
    }
  }

  if (actions.length) pushGroup(q ? copy.actions : copy.navigate, actions)

  if (!q && projectItems.length) {
    const recentOnly = projectItems.filter((i) => i.tone === 'muted')
    const restOnly = projectItems.filter((i) => i.tone !== 'muted')
    if (recentOnly.length) pushGroup(copy.recent, recentOnly)
    if (restOnly.length) pushGroup(copy.allProjects, restOnly)
  } else if (q && projectItems.length) {
    pushGroup(copy.projects, projectItems)
  } else if (!q && projects.length === 0) {

  }

  return { entries, selectable }
}
