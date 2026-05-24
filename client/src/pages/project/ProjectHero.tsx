import type { Project } from '../../api/projects'
import { StatusBadge } from '../../components/ui'
import { branchChip } from '../../styles/theme'
import { useI18n } from '../../i18n'
import { IconBranch, IconExternal } from './icons'
import {
  btnDangerGhost,
  hero,
  heroActions,
  heroMain,
  heroMeta,
  heroTitle,
  heroTitleRow,
  heroTop,
  projectAvatar,
  repoLink,
  repoText,
  roleBadge,
} from './styles'

function shortRepo(url: string): string {
  try {
    const u = new URL(url)
    const path = u.pathname.replace(/\/$/, '')
    return path.length > 1 ? `${u.host}${path}` : u.host
  } catch {
    return url
  }
}

function projectInitial(name: string): string {
  const t = name.trim()
  return (t[0] ?? '?').toUpperCase()
}

interface ProjectHeroProps {
  project: Project
  userRole: 'owner' | 'editor' | 'viewer' | 'loading'
  deleting: boolean
  onDelete: () => void
}

export function ProjectHero({ project, userRole, deleting, onDelete }: ProjectHeroProps) {
  const { t } = useI18n()
  const roleLabel =
    userRole === 'loading'
      ? t.project.roleLoading
      : userRole === 'owner'
        ? t.project.roleOwner
        : userRole === 'editor'
          ? t.project.roleEditor
          : t.project.roleViewer

  return (
    <header css={hero}>
      <div css={heroTop}>
        <div css={projectAvatar} aria-hidden>
          {projectInitial(project.name)}
        </div>
        <div css={heroMain}>
          <div css={heroTitleRow}>
            <h1 css={heroTitle}>{project.name}</h1>
            <StatusBadge status={project.status} />
            <span css={roleBadge}>{roleLabel}</span>
          </div>
          <div css={heroMeta}>
            <a
              href={project.repo_url}
              target="_blank"
              rel="noopener noreferrer"
              css={repoLink}
              title={project.repo_url}
            >
              <span css={repoText}>{shortRepo(project.repo_url)}</span>
              <IconExternal />
            </a>
            <span css={branchChip}>
              <IconBranch />
              {project.default_branch}
            </span>
          </div>
        </div>
        {userRole === 'owner' && (
          <div css={heroActions}>
            <button
              type="button"
              css={btnDangerGhost}
              disabled={deleting}
              onClick={onDelete}
              aria-label={`${t.project.delete} ${project.name}`}
            >
              {deleting ? `${t.common.delete}…` : t.project.delete}
            </button>
          </div>
        )}
      </div>
    </header>
  )
}
