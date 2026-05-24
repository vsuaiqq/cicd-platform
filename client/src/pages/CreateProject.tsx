import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { css } from '@emotion/react'
import { useCreateProjectMutation } from '../store/api/apiSlice'
import { Breadcrumb, useToast } from '../components/ui'
import { useI18n } from '../i18n'
import {
  pageTitle,
  pageDesc,
  form,
  formLabel,
  formInput,
  btnPrimary,
  btnSecondary,
  btnSecondaryLarge,
} from '../styles/theme'

const pageHeaderArea = css({
  marginBottom: 28,
  animation: 'fade-in 0.25s var(--ease-out) both',
})

const formActions = css({
  display: 'flex',
  alignItems: 'center',
  gap: 12,
  flexWrap: 'wrap',
})

const formCard = css({
  background: 'var(--bg-card)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-lg)',
  padding: '28px 28px',
  maxWidth: 520,
  boxShadow: 'var(--shadow-sm)',
  animation: 'fade-in 0.3s var(--ease-out) both',
})

const formHint = css({
  fontSize: '0.8125rem',
  color: 'var(--text-tertiary)',
  lineHeight: 1.5,
  marginTop: 4,
})

export default function CreateProject() {
  const toast = useToast()
  const { t } = useI18n()

  useEffect(() => {
    document.title = `${t.createProject.title} — Flow`
    return () => { document.title = 'Flow — CI/CD' }
  }, [t])

  const [name, setName] = useState('')
  const [repositoryUrl, setRepositoryUrl] = useState('')
  const [defaultBranch, setDefaultBranch] = useState('main')
  const [error, setError] = useState<string | null>(null)
  const [createProject, { isLoading: loading }] = useCreateProjectMutation()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    try {
      const project = await createProject({
        name: name.trim(),
        repository_url: repositoryUrl.trim(),
        default_branch: defaultBranch.trim() || 'main',
      }).unwrap()
      toast.success(t.toast.projectCreated.replace('{name}', project.name))
      navigate(`/projects/${project.id}/setup`, { replace: true })
    } catch (err) {
      const msg = err instanceof Error ? err.message : 'Failed to create project'
      setError(msg)
      toast.error(msg)
    }
  }

  return (
    <>
      <Breadcrumb
        items={[
          { label: t.nav.dashboard, to: '/' },
          { label: t.createProject.breadcrumb },
        ]}
      />

      <div css={pageHeaderArea}>
        <h1 css={pageTitle}>{t.createProject.title}</h1>
        <p css={pageDesc}>{t.createProject.desc}</p>
      </div>

      <div css={formCard}>
        <form css={form} onSubmit={handleSubmit} noValidate>
          <label css={formLabel}>
            {t.createProject.name}
            <input
              css={formInput}
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="my-app"
              required
              aria-invalid={!!error}
            />
          </label>

          <label css={formLabel}>
            {t.createProject.repoUrl}
            <input
              css={formInput}
              type="text"
              value={repositoryUrl}
              onChange={(e) => setRepositoryUrl(e.target.value)}
              placeholder="git@github.com:org/repo.git"
              required
              aria-invalid={!!error}
            />
            <span css={formHint}>{t.createProject.repoUrlHint}</span>
          </label>

          <label css={formLabel}>
            {t.createProject.defaultBranch}
            <input
              css={formInput}
              type="text"
              value={defaultBranch}
              onChange={(e) => setDefaultBranch(e.target.value)}
              placeholder="main"
            />
          </label>

          {error && (
            <p style={{ margin: 0, fontSize: '0.875rem', color: 'var(--danger)', lineHeight: 1.5 }} role="alert">
              {error}
            </p>
          )}

          <div css={formActions}>
            <button css={btnPrimary} type="submit" disabled={loading} aria-busy={loading}>
              {loading ? t.createProject.creating : t.createProject.create}
            </button>
            <Link to="/" css={[btnSecondary, btnSecondaryLarge]}>
              {t.createProject.cancel}
            </Link>
          </div>
        </form>
      </div>
    </>
  )
}
