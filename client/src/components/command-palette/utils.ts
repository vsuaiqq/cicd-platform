import { bestScore } from '../../lib/commandPaletteMatch'

export function shortRepo(url: string): string {
  try {
    const u = new URL(url)
    const path = u.pathname.replace(/\/$/, '')
    return path.length > 1 ? `${u.host}${path}` : u.host
  } catch {
    return url
  }
}


export function projectSearchScore(
  project: { name: string; repo_url: string; default_branch: string },
  query: string,
  recentBoost: boolean,
): number {
  const q = query.trim()
  if (!q) return recentBoost ? 100 : 50
  const base = bestScore([project.name, project.repo_url, project.default_branch], q)
  return base > 0 ? base + (recentBoost ? 15 : 0) : 0
}
