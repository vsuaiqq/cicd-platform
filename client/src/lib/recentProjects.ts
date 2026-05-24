import type { ProjectSummary } from '../api/projects'

const STORAGE_KEY = 'cicd_recent_projects'
const MAX_RECENT = 5

export interface RecentProject {
  id: string
  name: string
  repo_url: string
}

export function getRecentProjects(): RecentProject[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return []
    const parsed = JSON.parse(raw)
    if (!Array.isArray(parsed)) return []
    return parsed.filter(
      (p): p is RecentProject =>
        typeof p?.id === 'string' &&
        typeof p?.name === 'string' &&
        typeof p?.repo_url === 'string',
    )
  } catch {
    return []
  }
}


export function syncRecentWithProjects(projects: ProjectSummary[]): string[] {
  const validIds = new Set(projects.map((p) => p.id))
  const stored = getRecentProjects()
  const pruned = stored.filter((p) => validIds.has(p.id))

  if (pruned.length !== stored.length) {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(pruned))
    } catch {

    }
  }

  return pruned.map((p) => p.id)
}

export function trackProjectVisit(project: RecentProject): void {
  const recent = getRecentProjects().filter((p) => p.id !== project.id)
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify([project, ...recent].slice(0, MAX_RECENT)))
  } catch {

  }
}
