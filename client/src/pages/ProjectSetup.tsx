import { useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { css } from '@emotion/react'
import { useAppSelector } from '../store'
import { useGetProjectQuery, useVerifyProjectMutation } from '../store/api/apiSlice'
import { Breadcrumb, CodeBlock, SetupPageSkeleton, useToast } from '../components/ui'
import { useI18n } from '../i18n'
import {
  pageTitle,
  pageDesc,
  stepCard,
  stepTitle,
  stepSub,
  instructions,
  btnPrimary,
  btnSecondary,
} from '../styles/theme'

const verifySection = css({
  marginTop: 24,
  paddingTop: 24,
  borderTop: '1px solid var(--border-subtle)',
  animation: 'fade-in 0.3s var(--ease-out) both',
})

const actionsRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 12,
  flexWrap: 'wrap',
  marginBottom: 16,
})

const stepNumber = css({
  display: 'inline-flex',
  alignItems: 'center',
  justifyContent: 'center',
  width: 22,
  height: 22,
  borderRadius: 6,
  background: 'var(--accent-muted)',
  border: '1px solid var(--border-accent)',
  fontSize: '0.6875rem',
  fontWeight: 700,
  color: 'var(--accent)',
  flexShrink: 0,
})

const pageHeaderArea = css({
  marginBottom: 28,
  animation: 'fade-in 0.25s var(--ease-out) both',
})

export default function ProjectSetup() {
  const { id } = useParams<{ id: string }>()
  const toast = useToast()
  const { t } = useI18n()
  const { sessionValidated } = useAppSelector((s) => s.auth)
  const {
    data: project,
    isLoading,
    isError,
    isFetching,
  } = useGetProjectQuery(id!, { skip: !id || !sessionValidated })
  const [verifyProject, { isLoading: verifyLoading }] = useVerifyProjectMutation()

  const loading = isLoading || (isFetching && !project)

  useEffect(() => {
    document.title = project ? `Setup — ${project.name} — Flow` : 'Repository setup — Flow'
    return () => { document.title = 'Flow — CI/CD' }
  }, [project?.name])

  const handleVerify = async () => {
    if (!id) return
    try {
      const res = await verifyProject(id).unwrap()
      if (res.success) {
        toast.success(t.toast.repoConnected)
      } else {
        toast.error(res.error ?? t.errors.verificationFailed)
      }
    } catch (err) {
      toast.error(err instanceof Error ? err.message : t.errors.verificationFailed)
    }
  }

  if (loading) {
    return <SetupPageSkeleton />
  }

  if (isError || !project) {
    return (
      <>
        <Breadcrumb items={[{ label: t.nav.dashboard, to: '/' }]} />
        <p style={{ color: 'var(--danger)', marginBottom: 16 }}>{t.setup.notFound}</p>
        <Link to="/" css={btnSecondary}>{t.setup.back}</Link>
      </>
    )
  }

  return (
    <>
      <Breadcrumb
        items={[
          { label: t.nav.dashboard, to: '/' },
          { label: project.name, to: `/projects/${project.id}` },
          { label: t.setup.title },
        ]}
      />

      <div css={pageHeaderArea}>
        <h1 css={pageTitle}>{t.setup.repoSetupTitle}</h1>
        <p css={pageDesc}>{t.setup.repoSetupDesc}</p>
      </div>

      <div css={stepCard}>
        <h2 css={stepTitle}>
          <span css={stepNumber}>1</span>
          {t.setup.deployKeyTitle}
        </h2>
        <p css={stepSub}>{t.setup.deployKeySub}</p>
        <div css={instructions}>
          <ol>
            <li>{t.setup.deployKeyStep1}</li>
            <li>{t.setup.deployKeyStep2}</li>
            <li>{t.setup.deployKeyStep3}</li>
          </ol>
        </div>
        <CodeBlock value={project.public_key} label={t.setup.copyKey} />
      </div>

      <div css={stepCard}>
        <h2 css={stepTitle}>
          <span css={stepNumber}>2</span>
          {t.setup.webhookTitle}
        </h2>
        <p css={stepSub}>{t.setup.webhookSub}</p>
        <div css={instructions}>
          <p>{t.setup.webhookStep1}</p>
          <p>{t.setup.webhookPayloadUrl}</p>
        </div>
        <CodeBlock value={project.webhook_url || '(not configured)'} label={t.setup.copyUrl} />
        <div css={instructions} style={{ marginTop: 12 }}>
          <p>{t.setup.webhookSecret}</p>
        </div>
        <CodeBlock value={project.webhook_secret} label={t.setup.copySecret} />
        <div css={instructions} style={{ marginTop: 12 }}>
          {t.setup.webhookNote}
        </div>
      </div>

      <div css={verifySection}>
        <div css={actionsRow}>
          <button
            type="button"
            css={btnPrimary}
            onClick={handleVerify}
            disabled={verifyLoading || project.status === 'active'}
            aria-busy={verifyLoading}
          >
            {verifyLoading
              ? t.setup.verifying
              : project.status === 'active'
                ? t.setup.connected
                : t.setup.verify}
          </button>
          <Link to={`/projects/${project.id}`} css={btnSecondary}>
            {t.setup.skip}
          </Link>
        </div>
      </div>
    </>
  )
}
