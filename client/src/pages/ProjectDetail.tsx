import { useEffect, useRef, useState, useCallback } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { useI18n } from '../i18n'
import { css } from '@emotion/react'
import {
  AreaChart, Area, BarChart, Bar,
  XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, Cell,
} from 'recharts'
import { useAppSelector } from '../store'
import type { Project, EnvVar } from '../api/projects'
import { generatePipeline } from '../api/ai'
import { PipelineEditor } from '../components/PipelineEditor'
import type { Period } from '../api/analytics'
import {
  useDeleteProjectMutation,
  useGetDashboardQuery,
  useGetEnvVarsQuery,
  useGetPipelineYamlQuery,
  useGetProjectQuery,
  useInviteMemberMutation,
  useListMembersQuery,
  useListSecretsQuery,
  useRemoveMemberMutation,
  useSetPipelineYamlMutation,
  useSetSecretMutation,
  useDeleteSecretMutation,
  useUpdateEnvVarsMutation,
  useUpdateMemberRoleMutation,
  useUpdateProjectMutation,
} from '../store/api/apiSlice'
import { Breadcrumb, ConfirmDialog, useToast } from '../components/ui'
import { trackProjectVisit } from '../lib/recentProjects'
import {
  AnalyticsPanelSkeleton,
  ProjectPageSkeleton,
  SettingsSectionSkeleton,
} from '../components/ui'
import { ProjectView } from './project/ProjectView'
import { analyticsPeriodWrap, periodBar, periodBtn } from './project/styles'
import { btnSecondary } from '../styles/theme'



const fieldGroup = css({
  display: 'flex',
  flexDirection: 'column',
  gap: 4,
  marginBottom: 14,
  '&:last-child': { marginBottom: 0 },
})

const fieldLabel = css({
  fontSize: '0.75rem',
  fontWeight: 500,
  color: 'var(--text-secondary)',
  letterSpacing: '0.02em',
})

const fieldInput = css({
  padding: '8px 12px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-input, var(--bg-overlay))',
  color: 'var(--text-primary)',
  fontSize: '0.875rem',
  outline: 'none',
  width: '100%',
  boxSizing: 'border-box',
  transition: 'border-color 150ms',
  '&:focus': {
    outline: 'none',
    borderColor: 'var(--border-strong)',
    boxShadow: 'none',
  },
})

const btnPrimary = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  padding: '8px 16px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--accent)',
  color: 'var(--text-on-accent)',
  fontSize: '0.875rem',
  fontWeight: 600,
  border: 'none',
  cursor: 'pointer',
  transition: 'opacity 150ms',
  '&:disabled': { opacity: 0.5, cursor: 'not-allowed' },
})

const saveRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  marginTop: 16,
})

const savedMsg = css({
  fontSize: '0.8125rem',
  color: 'var(--success)',
  fontWeight: 500,
})

const errorMsg = css({
  fontSize: '0.8125rem',
  color: 'var(--danger)',
})



const envTableHead = css({
  fontSize: '0.75rem',
  color: 'var(--text-tertiary)',
  fontWeight: 500,
  textAlign: 'left',
  padding: '0 6px 8px',
  letterSpacing: '0.02em',
})

const envRow = css({
  display: 'grid',
  gridTemplateColumns: '1fr 1fr 32px',
  gap: 8,
  alignItems: 'center',
  marginBottom: 8,
})

const envInput = css({
  padding: '7px 10px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-input, var(--bg-overlay))',
  color: 'var(--text-primary)',
  fontSize: '0.8125rem',
  fontFamily: 'ui-monospace, monospace',
  outline: 'none',
  width: '100%',
  boxSizing: 'border-box',
  transition: 'border-color 150ms',
  '&:focus': {
    outline: 'none',
    borderColor: 'var(--border-strong)',
    boxShadow: 'none',
  },
  '&::placeholder': { color: 'var(--text-disabled)' },
})

const removeBtn = css({
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
  transition: 'color 150ms, background 150ms',
  '&:hover': { background: 'var(--danger-muted)', color: 'var(--danger)' },
})

const addVarBtn = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  padding: '6px 12px',
  borderRadius: 'var(--radius-md)',
  border: '1px dashed var(--border-default)',
  background: 'transparent',
  color: 'var(--text-tertiary)',
  fontSize: '0.8125rem',
  cursor: 'pointer',
  transition: 'border-color 150ms, color 150ms',
  '&:hover': { borderColor: 'var(--accent)', color: 'var(--accent)' },
})



interface EnvVarsEditorProps {
  projectId: string
  readOnly?: boolean
}

export function EnvVarsEditor({ projectId, readOnly = false }: EnvVarsEditorProps) {
  const toast = useToast()
  const { t } = useI18n()
  const { data, isLoading } = useGetEnvVarsQuery(projectId)
  const [updateEnvVars, { isLoading: saving }] = useUpdateEnvVarsMutation()
  const [vars, setVars] = useState<EnvVar[]>([])
  const [saved, setSaved] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const savedTimer = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    if (data) setVars(data)
  }, [data])

  useEffect(() => () => {
    if (savedTimer.current) clearTimeout(savedTimer.current)
  }, [])

  const handleChange = (i: number, field: 'key' | 'value', val: string) => {
    setVars((prev) => prev.map((v, idx) => idx === i ? { ...v, [field]: val } : v))
    setSaved(false)
  }

  const addVar = () => setVars((prev) => [...prev, { key: '', value: '' }])

  const removeVar = (i: number) => setVars((prev) => prev.filter((_, idx) => idx !== i))

  const handleSave = async () => {
    const trimmed = vars.filter((v) => v.key.trim() !== '')
    const keys = trimmed.map((v) => v.key.trim())
    const dup = keys.find((k, i) => keys.indexOf(k) !== i)
    if (dup) {
      setError(t.errors.duplicateKey.replace('{key}', dup))
      return
    }
    setError(null)
    try {
      await updateEnvVars({ id: projectId, vars: trimmed.map((v) => ({ key: v.key.trim(), value: v.value })) }).unwrap()
      setVars(trimmed)
      setSaved(true)
      toast.success(t.toast.envVarsSaved)
      savedTimer.current = setTimeout(() => setSaved(false), 3000)
    } catch (e) {
      const msg = e instanceof Error ? e.message : t.errors.failedToSave
      setError(msg)
      toast.error(msg)
    }
  }

  if (isLoading) {
    return <SettingsSectionSkeleton />
  }

  return (
    <div>
      {vars.length > 0 && (
        <div style={{ marginBottom: 4 }}>
          <div css={envRow}>
            <span css={envTableHead}>{t.project.key}</span>
            <span css={envTableHead}>{t.project.value}</span>
            {!readOnly && <span />}
          </div>
          {vars.map((v, i) => (
            <div key={i} css={envRow}>
              <input
                css={envInput}
                value={v.key}
                onChange={(e) => !readOnly && handleChange(i, 'key', e.target.value)}
                placeholder="KEY"
                spellCheck={false}
                readOnly={readOnly}
                disabled={readOnly}
                aria-label={`Env var ${i + 1} key`}
              />
              <input
                css={envInput}
                value={v.value}
                onChange={(e) => !readOnly && handleChange(i, 'value', e.target.value)}
                placeholder="value"
                spellCheck={false}
                readOnly={readOnly}
                disabled={readOnly}
                aria-label={`Env var ${i + 1} value`}
              />
              {!readOnly && (
                <button
                  type="button"
                  css={removeBtn}
                  onClick={() => removeVar(i)}
                  aria-label={t.aria.removeVariable}
                >
                  <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                    <path d="M2 2l8 8M10 2l-8 8" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
                  </svg>
                </button>
              )}
            </div>
          ))}
        </div>
      )}

      {vars.length === 0 && (
        <p style={{ fontSize: '0.8125rem', color: 'var(--text-tertiary)', marginBottom: 12 }}>
          {t.project.noEnvVarsText}
        </p>
      )}

      {!readOnly && (
        <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
          <button type="button" css={addVarBtn} onClick={addVar}>
            <svg width="12" height="12" viewBox="0 0 12 12" fill="none"><path d="M6 2v8M2 6h8" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" /></svg>
            {t.project.addVar}
          </button>
        </div>
      )}

      {!readOnly && (
        <div css={saveRow}>
          <button type="button" css={btnPrimary} onClick={handleSave} disabled={saving}>
            {saving ? t.project.saving : t.project.saveVars}
          </button>
          {saved && <span css={savedMsg}>{t.project.saved}</span>}
          {error && <span css={errorMsg}>{error}</span>}
        </div>
      )}
    </div>
  )
}



const secretsNotice = css({
  display: 'flex',
  alignItems: 'flex-start',
  gap: 8,
  padding: '10px 14px',
  borderRadius: 'var(--radius-md)',
  background: 'var(--warning-muted)',
  border: '1px solid var(--warning-glow)',
  fontSize: '0.8125rem',
  color: 'var(--warning)',
  marginBottom: 16,
  lineHeight: 1.5,
})

const secretRow = css({
  display: 'grid',
  gridTemplateColumns: '1fr auto auto',
  gap: 10,
  alignItems: 'center',
  padding: '9px 0',
  borderBottom: '1px solid var(--border-subtle)',
  '&:last-child': { borderBottom: 'none' },
})

const secretKey = css({
  fontFamily: 'ui-monospace, monospace',
  fontSize: '0.875rem',
  color: 'var(--text-primary)',
  fontWeight: 500,
})

const secretUpdated = css({
  fontSize: '0.75rem',
  color: 'var(--text-disabled)',
})

const secretEditArea = css({
  display: 'grid',
  gridTemplateColumns: '1fr 1fr',
  gap: 10,
  marginTop: 16,
  padding: '14px',
  background: 'var(--bg-overlay)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-md)',
})

const secretInput = css({
  padding: '8px 12px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-input, var(--bg-overlay))',
  color: 'var(--text-primary)',
  fontSize: '0.875rem',
  fontFamily: 'ui-monospace, monospace',
  outline: 'none',
  width: '100%',
  boxSizing: 'border-box',
  transition: 'border-color 150ms',
  '&:focus': { borderColor: 'var(--accent)' },
  '&::placeholder': { color: 'var(--text-disabled)' },
})

const btnSmall = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 5,
  padding: '6px 12px',
  borderRadius: 'var(--radius-md)',
  fontSize: '0.8125rem',
  fontWeight: 500,
  border: 'none',
  cursor: 'pointer',
  transition: 'opacity 150ms',
  '&:disabled': { opacity: 0.45, cursor: 'not-allowed' },
})

const btnSmallPrimary = css(btnSmall, {
  background: 'var(--accent)',
  color: 'var(--text-on-accent)',
})

const btnSmallDanger = css(btnSmall, {
  background: 'transparent',
  color: 'var(--danger)',
  border: '1px solid var(--danger-glow)',
  '&:hover': { background: 'var(--danger-muted)' },
})

const btnSmallGhost = css(btnSmall, {
  background: 'transparent',
  color: 'var(--text-secondary)',
  border: '1px solid var(--border-default)',
  '&:hover': { background: 'var(--bg-hover)' },
})

function timeAgoShort(iso: string): string {
  try {
    const diff = Date.now() - new Date(iso).getTime()
    const s = Math.floor(diff / 1000)
    if (s < 60) return 'just now'
    const m = Math.floor(s / 60)
    if (m < 60) return `${m}m ago`
    const h = Math.floor(m / 60)
    if (h < 24) return `${h}h ago`
    return `${Math.floor(h / 24)}d ago`
  } catch { return '' }
}



interface SecretsEditorProps {
  projectId: string
  readOnly?: boolean
}

export function SecretsEditor({ projectId, readOnly = false }: SecretsEditorProps) {
  const toast = useToast()
  const { t } = useI18n()
  const { data: secrets = [], isLoading: loading } = useListSecretsQuery(projectId)
  const [setSecretMut, { isLoading: saving }] = useSetSecretMutation()
  const [deleteSecretMut] = useDeleteSecretMutation()
  const [editKey, setEditKey]     = useState('')
  const [newKey, setNewKey]       = useState('')
  const [newValue, setNewValue]   = useState('')
  const [showValue, setShowValue] = useState(false)
  const [deleting, setDeleting]   = useState<string | null>(null)
  const [error, setError]         = useState<string | null>(null)
  const [formOpen, setFormOpen]   = useState(false)

  const openAdd = () => {
    setEditKey('')
    setNewKey('')
    setNewValue('')
    setShowValue(false)
    setError(null)
    setFormOpen(true)
  }

  const openEdit = (key: string) => {
    setEditKey(key)
    setNewKey(key)
    setNewValue('')
    setShowValue(false)
    setError(null)
    setFormOpen(true)
  }

  const closeForm = () => {
    setFormOpen(false)
    setError(null)
  }

  const handleSave = async () => {
    const k = (editKey || newKey).trim()
    if (!k) { setError(t.errors.keyRequired); return }
    if (!/^[A-Z_][A-Z0-9_]*$/.test(k)) {
      setError(t.errors.keyUpperCase)
      return
    }
    if (!newValue) { setError(t.errors.valueRequired); return }

    setError(null)
    try {
      await setSecretMut({ projectId, key: k, value: newValue }).unwrap()
      toast.success(editKey ? t.toast.secretUpdated.replace('{name}', k) : t.toast.secretAdded.replace('{name}', k))
      closeForm()
    } catch (e) {
      const msg = e instanceof Error ? e.message : t.errors.failedToSave
      setError(msg)
      toast.error(msg)
    }
  }

  const handleDelete = async (key: string) => {
    setDeleting(key)
    try {
      await deleteSecretMut({ projectId, key }).unwrap()
      toast.success(t.toast.secretDeleted.replace('{name}', key))
    } catch (e) {
      const msg = e instanceof Error ? e.message : t.errors.failedToDelete
      setError(msg)
      toast.error(msg)
    } finally {
      setDeleting(null)
    }
  }

  if (loading) return <SettingsSectionSkeleton />

  return (
    <div>
      <div css={secretsNotice}>
        <svg width="14" height="14" viewBox="0 0 16 16" fill="none" style={{ flexShrink: 0, marginTop: 1 }}>
          <circle cx="8" cy="8" r="6.5" stroke="currentColor" strokeWidth="1.2" />
          <path d="M8 5v4M8 11v.5" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
        </svg>
        <span>{t.project.secretsNoticeText}</span>
      </div>


      {secrets.length > 0 && (
        <div style={{ marginBottom: 16 }}>
          {secrets.map((s) => (
            <div key={s.key} css={secretRow}>
              <div>
                <div css={secretKey}>{s.key}</div>
                <div css={secretUpdated}>{t.project.updatedAt} {timeAgoShort(s.updated_at)}</div>
              </div>
              {!readOnly && (
                <button
                  type="button"
                  css={btnSmallGhost}
                  onClick={() => openEdit(s.key)}
                  style={{ fontSize: '0.75rem' }}
                >
                  {t.project.updateBtn}
                </button>
              )}
              {!readOnly ? (
                <button
                  type="button"
                  css={btnSmallDanger}
                  disabled={deleting === s.key}
                  onClick={() => handleDelete(s.key)}
                  style={{ fontSize: '0.75rem' }}
                  aria-label={t.aria.deleteSecret.replace('{name}', s.key)}
                  title={t.aria.deleteSecret.replace('{name}', s.key)}
                >
                  {deleting === s.key ? '…' : (
                    <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
                      <path d="M2 2l8 8M10 2l-8 8" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
                    </svg>
                  )}
                </button>
              ) : (
                <div style={{ color: 'var(--text-disabled)', fontSize: '0.75rem' }}>••••••••</div>
              )}
            </div>
          ))}
        </div>
      )}

      {secrets.length === 0 && (
        <p style={{ fontSize: '0.8125rem', color: 'var(--text-tertiary)', marginBottom: 14 }}>
          {t.project.noSecretsText}
        </p>
      )}


      {!readOnly && (formOpen ? (
        <div css={secretEditArea}>
          <div>
            <label css={fieldLabel} htmlFor="secret-key">Key</label>
            <input
              id="secret-key"
              css={secretInput}
              value={editKey || newKey}
              onChange={(e) => !editKey && setNewKey(e.target.value.toUpperCase())}
              readOnly={!!editKey}
              placeholder="DATABASE_URL"
              spellCheck={false}
              style={{ marginTop: 6, opacity: editKey ? 0.6 : 1 }}
            />
          </div>
          <div>
            <label css={fieldLabel} htmlFor="secret-value">
              {editKey ? t.project.newValueFor.replace('{key}', editKey) : t.project.secretValue}
            </label>
            <div style={{ position: 'relative', marginTop: 6 }}>
              <input
                id="secret-value"
                css={[secretInput, { paddingRight: 36 }]}
                type={showValue ? 'text' : 'password'}
                value={newValue}
                onChange={(e) => setNewValue(e.target.value)}
                placeholder="••••••••"
                spellCheck={false}
                autoComplete="new-password"
              />
              <button
                type="button"
                onClick={() => setShowValue((v) => !v)}
                style={{
                  position: 'absolute', right: 8, top: '50%', transform: 'translateY(-50%)',
                  background: 'none', border: 'none', cursor: 'pointer', color: 'var(--text-disabled)',
                  padding: 2, lineHeight: 1,
                }}
                aria-label={showValue ? t.aria.hideValue : t.aria.showValue}
              >
                {showValue ? (
                  <svg width="14" height="14" viewBox="0 0 16 16" fill="none"><path d="M1 8s2.5-5 7-5 7 5 7 5-2.5 5-7 5-7-5-7-5z" stroke="currentColor" strokeWidth="1.2"/><circle cx="8" cy="8" r="2" stroke="currentColor" strokeWidth="1.2"/><path d="M2 2l12 12" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round"/></svg>
                ) : (
                  <svg width="14" height="14" viewBox="0 0 16 16" fill="none"><path d="M1 8s2.5-5 7-5 7 5 7 5-2.5 5-7 5-7-5-7-5z" stroke="currentColor" strokeWidth="1.2"/><circle cx="8" cy="8" r="2" stroke="currentColor" strokeWidth="1.2"/></svg>
                )}
              </button>
            </div>
          </div>
          <div style={{ gridColumn: '1/-1', display: 'flex', alignItems: 'center', gap: 10, marginTop: 4 }}>
            <button type="button" css={btnSmallPrimary} onClick={handleSave} disabled={saving}>
              {saving ? t.project.saving : (editKey ? t.project.updateSecretBtn : t.project.addSecret)}
            </button>
            <button type="button" css={btnSmallGhost} onClick={closeForm}>{t.project.cancel}</button>
            {error && <span css={errorMsg}>{error}</span>}
          </div>
        </div>
      ) : (
        <button type="button" css={addVarBtn} onClick={openAdd}>
          <svg width="12" height="12" viewBox="0 0 12 12" fill="none">
            <path d="M6 2v8M2 6h8" stroke="currentColor" strokeWidth="1.4" strokeLinecap="round" />
          </svg>
          {t.project.addSecret}
        </button>
      ))}
    </div>
  )
}



const memberRow = css({
  display: 'grid',
  gridTemplateColumns: '1fr auto auto',
  alignItems: 'center',
  gap: 12,
  padding: '10px 0',
  borderBottom: '1px solid var(--border-subtle)',
  '&:last-child': { borderBottom: 'none' },
})

const memberInfo = css({
  display: 'flex',
  flexDirection: 'column',
  gap: 2,
  overflow: 'hidden',
})

const memberEmail = css({
  fontSize: '0.875rem',
  fontWeight: 500,
  color: 'var(--text-primary)',
  overflow: 'hidden',
  textOverflow: 'ellipsis',
  whiteSpace: 'nowrap',
})

const memberMeta = css({
  fontSize: '0.75rem',
  color: 'var(--text-tertiary)',
})

const roleBadge = (role: string) => css({
  display: 'inline-flex',
  alignItems: 'center',
  padding: '2px 8px',
  borderRadius: 4,
  fontSize: '0.75rem',
  fontWeight: 600,
  background: role === 'owner'
    ? 'var(--ai-muted)'
    : role === 'editor'
      ? 'var(--accent-muted)'
      : 'var(--bg-hover)',
  color: role === 'owner'
    ? 'var(--ai)'
    : role === 'editor'
      ? 'var(--accent)'
      : 'var(--text-secondary)',
})

const inviteForm = css({
  display: 'grid',
  gridTemplateColumns: '1fr auto auto',
  gap: 8,
  marginTop: 16,
  alignItems: 'flex-start',
})

const inviteInput = css({
  padding: '8px 12px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-input, var(--bg-overlay))',
  color: 'var(--text-primary)',
  fontSize: '0.875rem',
  outline: 'none',
  width: '100%',
  boxSizing: 'border-box',
  transition: 'border-color 150ms',
  '&:focus': { borderColor: 'var(--accent)' },
  '&::placeholder': { color: 'var(--text-disabled)' },
})

const roleSelect = css({
  padding: '8px 10px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-input, var(--bg-overlay))',
  color: 'var(--text-primary)',
  fontSize: '0.875rem',
  outline: 'none',
  cursor: 'pointer',
  '&:focus': { borderColor: 'var(--accent)' },
})

interface MembersPanelProps {
  projectId: string
  currentUserId: string
}

export function MembersPanel({ projectId, currentUserId }: MembersPanelProps) {
  const { t } = useI18n()
  const toast = useToast()
  const { data, isLoading: loading, isError } = useListMembersQuery(projectId)
  const [inviteMember, { isLoading: inviting }] = useInviteMemberMutation()
  const [updateMemberRole] = useUpdateMemberRoleMutation()
  const [removeMember] = useRemoveMemberMutation()
  const [inviteEmail, setInviteEmail] = useState('')
  const [inviteRole, setInviteRole] = useState<'editor' | 'viewer'>('viewer')
  const [removing, setRemoving] = useState<string | null>(null)

  const members = data?.members ?? []
  const requesterRole = data?.requester_role ?? 'viewer'
  const ownerUserId = data?.owner_user_id ?? ''

  useEffect(() => {
    if (isError) toast.error(t.members.loadError)
  }, [isError, toast, t])

  const handleInvite = async () => {
    if (!inviteEmail.trim()) return
    try {
      await inviteMember({ projectId, email: inviteEmail.trim(), role: inviteRole }).unwrap()
      setInviteEmail('')
      toast.success(t.members.memberAdded)
    } catch (e) {
      const msg = e instanceof Error ? e.message : ''
      if (msg.includes('not found')) toast.error(t.members.userNotFound)
      else if (msg.includes('already')) toast.error(t.members.alreadyMember)
      else toast.error(t.members.inviteError)
    }
  }

  const handleRoleChange = async (userId: string, role: 'editor' | 'viewer') => {
    try {
      await updateMemberRole({ projectId, userId, role }).unwrap()
      toast.success(t.members.roleUpdated)
    } catch {
      toast.error(t.members.roleError)
    }
  }

  const handleRemove = async (userId: string) => {
    setRemoving(userId)
    try {
      await removeMember({ projectId, userId }).unwrap()
      toast.success(t.members.memberRemoved)
    } catch {
      toast.error(t.members.removeError)
    } finally {
      setRemoving(null)
    }
  }

  const roleLabel = (role: string) => {
    if (role === 'owner') return t.members.owner
    if (role === 'editor') return t.members.editor
    return t.members.viewer
  }

  if (loading) {
    return <SettingsSectionSkeleton />
  }

  const isOwner = requesterRole === 'owner'
  const iAmOwner = currentUserId === ownerUserId


  const ownerLabel = iAmOwner ? `${t.members.you} — ${t.members.owner}` : t.members.owner


  const selfEntry = !iAmOwner ? members.find((m) => m.user_id === currentUserId) : null
  const otherMembers = members.filter((m) => m.user_id !== currentUserId)

  return (
    <div>

      <div css={memberRow}>
        <div css={memberInfo}>
          <div css={memberEmail}>{ownerLabel}</div>
          <div css={memberMeta}>{t.members.ownerNote}</div>
        </div>
        <div css={roleBadge('owner')}>{t.members.owner}</div>
        <div style={{ width: iAmOwner ? 70 : 0 }} />
      </div>


      {selfEntry && (
        <div css={memberRow}>
          <div css={memberInfo}>
            <div css={memberEmail}>{selfEntry.email} <span style={{ color: 'var(--text-tertiary)', fontWeight: 400 }}>— {t.members.you}</span></div>
            {selfEntry.display_name && <div css={memberMeta}>{selfEntry.display_name}</div>}
          </div>
          <div css={roleBadge(selfEntry.role)}>{roleLabel(selfEntry.role)}</div>
          <div style={{ width: 0 }} />
        </div>
      )}


      {otherMembers.length === 0 && !selfEntry && (
        <p style={{ fontSize: '0.8125rem', color: 'var(--text-tertiary)', marginTop: 8, marginBottom: 0 }}>
          {t.members.noMembers}
        </p>
      )}
      {otherMembers.map((m) => (
        <div key={m.user_id} css={memberRow}>
          <div css={memberInfo}>
            <div css={memberEmail}>{m.email}</div>
            {m.display_name && <div css={memberMeta}>{m.display_name}</div>}
          </div>
          {isOwner ? (
            <select
              css={roleSelect}
              value={m.role}
              onChange={(e) => handleRoleChange(m.user_id, e.target.value as 'editor' | 'viewer')}
              style={{ fontSize: '0.75rem', padding: '4px 8px' }}
            >
              <option value="editor">{t.members.editor}</option>
              <option value="viewer">{t.members.viewer}</option>
            </select>
          ) : (
            <div css={roleBadge(m.role)}>{roleLabel(m.role)}</div>
          )}
          {isOwner ? (
            <button
              type="button"
              css={btnSmallDanger}
              disabled={removing === m.user_id}
              onClick={() => handleRemove(m.user_id)}
              style={{ fontSize: '0.75rem' }}
            >
              {removing === m.user_id ? t.members.removing : t.members.removeBtn}
            </button>
          ) : (
            <div style={{ width: 0 }} />
          )}
        </div>
      ))}


      {isOwner && (
        <div css={inviteForm}>
          <input
            css={inviteInput}
            type="email"
            placeholder={t.members.emailPlaceholder}
            value={inviteEmail}
            onChange={(e) => setInviteEmail(e.target.value)}
            onKeyDown={(e) => { if (e.key === 'Enter' && !inviting) handleInvite() }}
            disabled={inviting}
          />
          <select
            css={roleSelect}
            value={inviteRole}
            onChange={(e) => setInviteRole(e.target.value as 'editor' | 'viewer')}
            disabled={inviting}
          >
            <option value="viewer">{t.members.viewer}</option>
            <option value="editor">{t.members.editor}</option>
          </select>
          <button
            type="button"
            css={btnSmallPrimary}
            onClick={handleInvite}
            disabled={inviting || !inviteEmail.trim()}
          >
            {inviting ? t.members.inviting : t.members.inviteBtn}
          </button>
        </div>
      )}
    </div>
  )
}



const toggleRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 10,
  marginBottom: 14,
})

const toggleLabel = css({
  fontSize: '0.875rem',
  color: 'var(--text-primary)',
  cursor: 'pointer',
  userSelect: 'none',
})


const dividerLine = css({
  borderTop: '1px solid var(--border-subtle)',
  margin: '20px 0',
})

const aiGenSection = css({
  background: 'var(--bg-overlay)',
  border: '1px solid var(--border-subtle)',
  borderRadius: 'var(--radius-md)',
  padding: '16px',
  marginBottom: 16,
})

const aiGenTitle = css({
  fontSize: '0.8125rem',
  fontWeight: 600,
  color: 'var(--text-secondary)',
  marginBottom: 10,
  display: 'flex',
  alignItems: 'center',
  gap: 6,
})

const aiGenInput = css({
  width: '100%',
  padding: '9px 12px',
  borderRadius: 'var(--radius-md)',
  border: '1px solid var(--border-default)',
  background: 'var(--bg-input, var(--bg-overlay))',
  color: 'var(--text-primary)',
  fontSize: '0.875rem',
  outline: 'none',
  boxSizing: 'border-box',
  marginBottom: 10,
  transition: 'border-color 150ms',
  '&:focus': { borderColor: 'var(--accent)' },
  '&::placeholder': { color: 'var(--text-disabled)' },
})

const aiThinkingWrap = css({
  position: 'relative',
  marginTop: 12,
  padding: '14px 16px',
  borderRadius: 'var(--radius-md)',
  background: 'linear-gradient(135deg, var(--ai-muted), transparent)',
  border: '1px solid var(--ai-glow)',
  overflow: 'hidden',
  animation: 'fade-in 0.25s var(--ease-out) both',
  display: 'flex',
  alignItems: 'center',
  gap: 12,
})

const aiThinkingBar = css({
  position: 'absolute',
  top: 0,
  left: 0,
  width: '35%',
  height: 2,
  background: 'linear-gradient(90deg, transparent, var(--ai), var(--ai-glow), transparent)',
  animation: 'ai-scan 1.8s ease-in-out infinite',
})

const aiThinkingLabel = css({
  flex: 1,
  fontSize: '0.8125rem',
  fontWeight: 500,
  color: 'var(--ai)',
  display: 'flex',
  alignItems: 'center',
  gap: 8,
})

const aiThinkingDotsRow = css({
  display: 'flex',
  alignItems: 'center',
  gap: 5,
})

const aiDotBase = {
  width: 5,
  height: 5,
  borderRadius: '50%' as const,
  background: 'var(--ai)',
  display: 'inline-block' as const,
  animation: 'thinking-dot 1.3s ease-in-out infinite',
}
const aiThinkingDot1 = css({ ...aiDotBase })
const aiThinkingDot2 = css({ ...aiDotBase, animationDelay: '0.18s' })
const aiThinkingDot3 = css({ ...aiDotBase, animationDelay: '0.36s' })

const btnAI = css({
  display: 'inline-flex',
  alignItems: 'center',
  gap: 6,
  padding: '8px 16px',
  borderRadius: 'var(--radius-md)',
  background: 'linear-gradient(135deg, var(--ai), var(--indigo))',
  color: 'var(--text-on-accent)',
  fontSize: '0.875rem',
  fontWeight: 600,
  border: 'none',
  cursor: 'pointer',
  transition: 'opacity 150ms',
  '&:disabled': { opacity: 0.5, cursor: 'not-allowed' },
})



interface PipelineYAMLEditorProps {
  projectId: string
  readOnly?: boolean
}

function IconAI() {
  return (
    <svg width="14" height="14" viewBox="0 0 16 16" fill="none" aria-hidden>
      <path d="M8 2 L10 6 L14 8 L10 10 L8 14 L6 10 L2 8 L6 6 Z" stroke="currentColor" strokeWidth="1.2" fill="none" strokeLinejoin="round" />
    </svg>
  )
}

function GeneratingPanel({ label }: { label: string }) {
  return (
    <div css={aiThinkingWrap}>
      <div css={aiThinkingBar} />
      <span css={aiThinkingLabel}>
        <span style={{ display: 'inline-flex', animation: 'spin 2.5s linear infinite' }}>
          <IconAI />
        </span>
        {label}
      </span>
      <span css={aiThinkingDotsRow}>
        <span css={aiThinkingDot1} />
        <span css={aiThinkingDot2} />
        <span css={aiThinkingDot3} />
      </span>
    </div>
  )
}

export function PipelineYAMLEditor({ projectId, readOnly = false }: PipelineYAMLEditorProps) {
  const toast = useToast()
  const { t } = useI18n()
  const token = useAppSelector((s) => s.auth.accessToken) ?? ''
  const { data: yamlData, isLoading: loading } = useGetPipelineYamlQuery(projectId)
  const [setPipelineYaml, { isLoading: saving }] = useSetPipelineYamlMutation()
  const [draft, setDraft]         = useState('')
  const [enabled, setEnabled]     = useState(false)
  const [saved, setSaved]         = useState(false)
  const [saveError, setSaveError] = useState<string | null>(null)
  const savedTimer                = useRef<ReturnType<typeof setTimeout> | null>(null)
  const [yamlSynced, setYamlSynced] = useState(false)


  const [description, setDescription]   = useState('')
  const [generating, setGenerating]     = useState(false)
  const [genError, setGenError]         = useState<string | null>(null)

  useEffect(() => {
    setYamlSynced(false)
  }, [projectId])

  useEffect(() => {
    if (yamlData !== undefined && !yamlSynced) {
      setDraft(yamlData)
      setEnabled(yamlData !== '')
      setYamlSynced(true)
    }
  }, [yamlData, yamlSynced])

  useEffect(() => () => { if (savedTimer.current) clearTimeout(savedTimer.current) }, [])

  const handleSave = useCallback(async () => {
    setSaveError(null)
    try {
      const toSave = enabled ? draft.trim() : ''
      await setPipelineYaml({ id: projectId, yaml: toSave }).unwrap()
      setSaved(true)
      toast.success(t.toast.pipelineConfigSaved)
      savedTimer.current = setTimeout(() => setSaved(false), 3000)
    } catch (e) {
      const msg = e instanceof Error ? e.message : t.errors.failedToSave
      setSaveError(msg)
      toast.error(msg)
    }
  }, [projectId, enabled, draft, toast, t, setPipelineYaml])

  const handleGenerate = useCallback(async () => {
    if (!description.trim() || !token) return
    setGenerating(true)
    setGenError(null)
    try {
      const yaml = await generatePipeline(description.trim(), token)
      setDraft(yaml)
      setEnabled(true)
      toast.success(t.toast.pipelineGenerated)
    } catch (e) {
      const msg = e instanceof Error ? e.message : t.errors.generationFailed
      setGenError(msg)
      toast.error(msg)
    } finally {
      setGenerating(false)
    }
  }, [description, token, toast, t])

  if (loading) {
    return <SettingsSectionSkeleton />
  }

  return (
    <div>

      {!readOnly && (
        <div css={aiGenSection}>
          <div css={aiGenTitle}>
            <IconAI />
            {t.project.aiGenTitle}
          </div>
          <p style={{ fontSize: '0.8125rem', color: 'var(--text-tertiary)', marginBottom: 10 }}>
            {t.project.aiGenDesc}
          </p>
          <input
            css={aiGenInput}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            onKeyDown={(e) => { if (e.key === 'Enter' && !generating) handleGenerate() }}
            placeholder="e.g. Go REST API with PostgreSQL — lint, tests with -race, build Docker image"
            disabled={generating}
          />
          <div style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
            <button type="button" css={btnAI} onClick={handleGenerate} disabled={generating || !description.trim()}>
              <span style={generating ? { display: 'inline-flex', animation: 'spin 2.5s linear infinite' } : {}}>
                <IconAI />
              </span>
              {generating ? t.project.generating : t.project.generatePipelineBtn}
            </button>
            {genError && <span css={errorMsg}>{genError}</span>}
          </div>
          {generating && <GeneratingPanel label={t.project.generatingPipeline} />}
        </div>
      )}

      {!readOnly && <div css={dividerLine} />}


      <div css={toggleRow}>
        <input
          id="yaml-override-toggle"
          type="checkbox"
          checked={enabled}
          onChange={(e) => !readOnly && setEnabled(e.target.checked)}
          disabled={readOnly}
          style={{ width: 16, height: 16, cursor: readOnly ? 'not-allowed' : 'pointer', accentColor: 'var(--accent)' }}
        />
        <label css={toggleLabel} htmlFor="yaml-override-toggle">
          {t.project.overrideToggleLabel}
        </label>
      </div>

      <PipelineEditor
        value={draft}
        onChange={(val) => { if (!readOnly) { setDraft(val); setSaved(false) } }}
        disabled={!enabled || readOnly}
        minHeight="300px"
      />
      {enabled && (
        <p style={{ fontSize: '0.75rem', color: 'var(--text-disabled)', marginTop: 10, marginBottom: 0 }}>
          {t.project.overrideNote}
        </p>
      )}

      {!readOnly && (
        <div css={saveRow}>
          <button type="button" css={btnPrimary} onClick={handleSave} disabled={saving}>
            {saving ? t.project.saving : t.project.saveConfigBtn}
          </button>
          {saved && <span css={savedMsg}>{t.project.saved}</span>}
          {saveError && <span css={errorMsg}>{saveError}</span>}
        </div>
      )}
    </div>
  )
}



interface ProjectSettingsEditorProps {
  project: Project
  onUpdated: (p: Project) => void
  readOnly?: boolean
}

export function ProjectSettingsEditor({ project, onUpdated, readOnly = false }: ProjectSettingsEditorProps) {
  const toast = useToast()
  const { t } = useI18n()
  const [updateProject, { isLoading: saving }] = useUpdateProjectMutation()
  const [name, setName]               = useState(project.name)
  const [branch, setBranch]           = useState(project.default_branch)
  const [saved, setSaved]             = useState(false)
  const [error, setError]             = useState<string | null>(null)
  const savedTimer                    = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => () => {
    if (savedTimer.current) clearTimeout(savedTimer.current)
  }, [])

  useEffect(() => {
    setName(project.name)
    setBranch(project.default_branch)
  }, [project.id, project.name, project.default_branch])

  const handleSave = async () => {
    if (!name.trim()) {
      setError(t.errors.nameRequired)
      return
    }
    setError(null)
    try {
      const updated = await updateProject({
        id: project.id,
        body: { name: name.trim(), default_branch: branch.trim() || 'main' },
      }).unwrap()
      onUpdated(updated)
      setSaved(true)
      toast.success(t.toast.projectSettingsSaved)
      savedTimer.current = setTimeout(() => setSaved(false), 3000)
    } catch (e) {
      const msg = e instanceof Error ? e.message : t.errors.failedToSave
      setError(msg)
      toast.error(msg)
    }
  }

  return (
    <div>
      <div css={fieldGroup}>
        <label css={fieldLabel} htmlFor="proj-name">{t.project.projectName}</label>
        <input
          id="proj-name"
          css={fieldInput}
          value={name}
          onChange={(e) => { if (!readOnly) { setName(e.target.value); setSaved(false) } }}
          readOnly={readOnly}
          disabled={readOnly}
        />
      </div>
      <div css={fieldGroup}>
        <label css={fieldLabel} htmlFor="proj-branch">{t.project.defaultBranch}</label>
        <input
          id="proj-branch"
          css={fieldInput}
          value={branch}
          onChange={(e) => { if (!readOnly) { setBranch(e.target.value); setSaved(false) } }}
          placeholder="main"
          readOnly={readOnly}
          disabled={readOnly}
        />
      </div>
      <div css={fieldGroup}>
        <span css={fieldLabel}>{t.project.repoUrl}</span>
        <span style={{ fontFamily: 'ui-monospace, monospace', fontSize: '0.8125rem', color: 'var(--text-secondary)', padding: '8px 0' }}>
          {project.repo_url}
        </span>
      </div>

      {!readOnly && (
        <div css={saveRow}>
          <button type="button" css={btnPrimary} onClick={handleSave} disabled={saving}>
            {saving ? t.project.saving2 : t.project.saveProject}
          </button>
          {saved && <span css={savedMsg}>{t.project.saved}</span>}
          {error && <span css={errorMsg}>{error}</span>}
        </div>
      )}
    </div>
  )
}



const analyticsWrap = css({ animation: 'fade-in 0.3s var(--ease-out) both' })
const statGrid = css({
  display: 'grid', gridTemplateColumns: 'repeat(auto-fill, minmax(180px, 1fr))', gap: 14, marginBottom: 28,
})
const statTile = css({
  background: 'var(--bg-card)',
  border: '1px solid var(--border-default)',
  borderRadius: 10, padding: '16px 18px',
})
const statTileLabel = css({ fontSize: '0.75rem', color: 'var(--text-tertiary)', marginBottom: 4, fontWeight: 500 })
const statTileValue = css({ fontSize: '1.5rem', fontWeight: 700, color: 'var(--text-primary)' })
const statTileSub = css({ fontSize: '0.75rem', color: 'var(--text-tertiary)', marginTop: 2 })
const chartCard = css({
  background: 'var(--bg-card)', border: '1px solid var(--border-default)',
  borderRadius: 10, padding: '20px 22px', marginBottom: 18,
})
const chartTitle = css({ fontSize: '0.875rem', fontWeight: 600, color: 'var(--text-secondary)', marginBottom: 16 })
const twoCol = css({ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 18, marginBottom: 18,
  '@media (max-width: 860px)': { gridTemplateColumns: '1fr' } })
const jobTable = css({ width: '100%', borderCollapse: 'collapse' as const })
const jobTh = css({ textAlign: 'left' as const, fontSize: '0.75rem', color: 'var(--text-tertiary)',
  fontWeight: 600, paddingBottom: 8, borderBottom: '1px solid var(--border-default)' })
const jobTd = css({ fontSize: '0.8125rem', color: 'var(--text-secondary)',
  padding: '7px 0', borderBottom: '1px solid var(--border-faint)', verticalAlign: 'middle' as const })

function fmtDur(sec: number): string {
  if (sec < 60) return `${Math.round(sec)}s`
  if (sec < 3600) return `${Math.round(sec / 60)}m ${Math.round(sec % 60)}s`
  return `${Math.floor(sec / 3600)}h ${Math.round((sec % 3600) / 60)}m`
}

function fmtPct(rate: number): string {
  return `${(rate * 100).toFixed(1)}%`
}

function RateDot({ rate }: { rate: number }) {
  const color = rate > 0.5 ? 'var(--status-failed)' : rate > 0.2 ? 'var(--status-running)' : 'var(--status-success)'
  return <span style={{ display: 'inline-block', width: 8, height: 8, borderRadius: '50%', background: color, marginRight: 6 }} />
}

interface ProjectAnalyticsProps {
  projectId: string
}

export function ProjectAnalytics({ projectId }: ProjectAnalyticsProps) {
  const { t } = useI18n()
  const [period, setPeriod] = useState<Period>('30d')
  const { data, isLoading: loading, isError, error: queryError } = useGetDashboardQuery({
    projectId,
    period,
  })
  const error = isError
    ? ((queryError as { message?: string })?.message ?? t.errors.failedToLoadAnalytics)
    : null

  return (
    <div css={analyticsWrap}>

      <div css={analyticsPeriodWrap}>
        <div css={periodBar}>
          {(['7d', '30d', '90d'] as Period[]).map((p) => (
            <button key={p} type="button" css={periodBtn(period === p)} onClick={() => setPeriod(p)}>
              {p === '7d' ? t.project.period7d : p === '30d' ? t.project.period30d : t.project.period90d}
            </button>
          ))}
        </div>
      </div>

      {loading && <AnalyticsPanelSkeleton />}
      {error && (
        <div style={{ color: 'var(--status-failed)', fontSize: '0.875rem', padding: '12px 0' }}>{error}</div>
      )}

      {!loading && !error && data && (
        <>

          <div css={statGrid}>
            <div css={statTile}>
              <div css={statTileLabel}>{t.project.totalRuns}</div>
              <div css={statTileValue}>{data.total_runs}</div>
            </div>
            <div css={statTile}>
              <div css={statTileLabel}>{t.project.successRate}</div>
              <div css={statTileValue} style={{ color: data.success_rate >= 0.8 ? 'var(--status-success)' : data.success_rate >= 0.5 ? 'var(--status-running)' : 'var(--status-failed)' }}>
                {fmtPct(data.success_rate)}
              </div>
              <div css={statTileSub}>{data.success_count} {t.project.statPassed} · {data.failed_count} {t.project.statFailed} · {data.cancelled_count} {t.project.statCancelled}</div>
            </div>
            <div css={statTile}>
              <div css={statTileLabel}>{t.project.avgDuration}</div>
              <div css={statTileValue}>{fmtDur(data.avg_duration_sec)}</div>
            </div>
            <div css={statTile}>
              <div css={statTileLabel}>{t.project.p50Duration}</div>
              <div css={statTileValue}>{fmtDur(data.p50_duration_sec)}</div>
            </div>
            <div css={statTile}>
              <div css={statTileLabel}>{t.project.p95Duration}</div>
              <div css={statTileValue}>{fmtDur(data.p95_duration_sec)}</div>
            </div>
          </div>


          <div css={chartCard}>
            <div css={chartTitle}>{t.project.runsChartTitle} {period === '7d' ? t.project.periodLabel7d : period === '30d' ? t.project.periodLabel30d : t.project.periodLabel90d}</div>
            {data.trend.length === 0 ? (
              <div style={{ color: 'var(--text-tertiary)', fontSize: '0.8125rem' }}>{t.project.noRunData}</div>
            ) : (
              <ResponsiveContainer width="100%" height={200}>
                <AreaChart data={data.trend} margin={{ top: 4, right: 8, left: -20, bottom: 0 }}>
                  <defs>
                    <linearGradient id="gradSuccess" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="var(--status-success)" stopOpacity={0.25} />
                      <stop offset="95%" stopColor="var(--status-success)" stopOpacity={0} />
                    </linearGradient>
                    <linearGradient id="gradFailed" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="var(--status-failed)" stopOpacity={0.25} />
                      <stop offset="95%" stopColor="var(--status-failed)" stopOpacity={0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border-faint)" />
                  <XAxis dataKey="date" tick={{ fontSize: 11, fill: 'var(--text-tertiary)' }}
                    tickFormatter={(v: string) => v.slice(5)} />
                  <YAxis tick={{ fontSize: 11, fill: 'var(--text-tertiary)' }} allowDecimals={false} />
                  <Tooltip
                    contentStyle={{ background: 'var(--bg-card)', border: '1px solid var(--border-default)', borderRadius: 8, fontSize: 12, color: 'var(--text-primary)' }}
                    labelStyle={{ color: 'var(--text-secondary)', fontWeight: 600 }}
                    itemStyle={{ color: 'var(--text-secondary)' }}
                  />
                  <Area type="monotone" dataKey="success" name={t.chart.success}
                    stroke="var(--status-success)" fill="url(#gradSuccess)" strokeWidth={2} dot={false} />
                  <Area type="monotone" dataKey="failed" name={t.chart.failed}
                    stroke="var(--status-failed)" fill="url(#gradFailed)" strokeWidth={2} dot={false} />
                </AreaChart>
              </ResponsiveContainer>
            )}
          </div>


          <div css={twoCol}>

            <div css={chartCard} style={{ marginBottom: 0 }}>
              <div css={chartTitle}>{t.project.topFailingJobs}</div>
              {data.top_failing_jobs.length === 0 ? (
                <div style={{ color: 'var(--text-tertiary)', fontSize: '0.8125rem' }}>{t.project.notEnoughData}</div>
              ) : (
                <>
                  <ResponsiveContainer width="100%" height={160}>
                    <BarChart data={data.top_failing_jobs} layout="vertical" margin={{ top: 0, right: 8, left: 0, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" stroke="var(--border-faint)" horizontal={false} />
                      <XAxis type="number" domain={[0, 1]} tickFormatter={(v: number) => `${Math.round(v * 100)}%`}
                        tick={{ fontSize: 10, fill: 'var(--text-tertiary)' }} />
                      <YAxis type="category" dataKey="job_name" width={72}
                        tick={{ fontSize: 11, fill: 'var(--text-secondary)' }} />
                      <Tooltip
                        formatter={(v) => [`${(Number(v ?? 0) * 100).toFixed(1)}%`, t.chart.failureRate]}
                        contentStyle={{ background: 'var(--bg-card)', border: '1px solid var(--border-default)', borderRadius: 8, fontSize: 12, color: 'var(--text-primary)' }}
                        itemStyle={{ color: 'var(--text-secondary)' }}
                      />
                      <Bar dataKey="failure_rate" name={t.chart.failureRate} radius={[0, 4, 4, 0]}>
                        {data.top_failing_jobs.map((entry, i) => (
                          <Cell key={i} fill={entry.failure_rate > 0.5 ? 'var(--status-failed)' : 'var(--status-running)'} />
                        ))}
                      </Bar>
                    </BarChart>
                  </ResponsiveContainer>
                  <table css={jobTable} style={{ marginTop: 12 }}>
                    <thead>
                      <tr>
                        <th css={jobTh}>{t.project.jobCol}</th>
                        <th css={jobTh} style={{ textAlign: 'right' }}>{t.project.failRateCol}</th>
                        <th css={jobTh} style={{ textAlign: 'right' }}>{t.project.runsCol}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {data.top_failing_jobs.map((j) => (
                        <tr key={j.job_name}>
                          <td css={jobTd}><RateDot rate={j.failure_rate} />{j.job_name}</td>
                          <td css={jobTd} style={{ textAlign: 'right', fontVariantNumeric: 'tabular-nums' }}>{fmtPct(j.failure_rate)}</td>
                          <td css={jobTd} style={{ textAlign: 'right', color: 'var(--text-tertiary)' }}>{j.total_runs}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </>
              )}
            </div>


            <div css={chartCard} style={{ marginBottom: 0 }}>
              <div css={chartTitle}>{t.project.slowestJobs}</div>
              {data.top_slow_jobs.length === 0 ? (
                <div style={{ color: 'var(--text-tertiary)', fontSize: '0.8125rem' }}>{t.project.notEnoughData}</div>
              ) : (
                <>
                  <ResponsiveContainer width="100%" height={160}>
                    <BarChart data={data.top_slow_jobs} layout="vertical" margin={{ top: 0, right: 8, left: 0, bottom: 0 }}>
                      <CartesianGrid strokeDasharray="3 3" stroke="var(--border-faint)" horizontal={false} />
                      <XAxis type="number" tickFormatter={(v: number) => fmtDur(v)}
                        tick={{ fontSize: 10, fill: 'var(--text-tertiary)' }} />
                      <YAxis type="category" dataKey="job_name" width={72}
                        tick={{ fontSize: 11, fill: 'var(--text-secondary)' }} />
                      <Tooltip
                        formatter={(v) => [fmtDur(Number(v ?? 0)), t.chart.avgDuration]}
                        contentStyle={{ background: 'var(--bg-card)', border: '1px solid var(--border-default)', borderRadius: 8, fontSize: 12, color: 'var(--text-primary)' }}
                        itemStyle={{ color: 'var(--text-secondary)' }}
                      />
                      <Bar dataKey="avg_duration_sec" name={t.chart.avgDuration} fill="var(--accent)" radius={[0, 4, 4, 0]} />
                    </BarChart>
                  </ResponsiveContainer>
                  <table css={jobTable} style={{ marginTop: 12 }}>
                    <thead>
                      <tr>
                        <th css={jobTh}>{t.project.jobCol}</th>
                        <th css={jobTh} style={{ textAlign: 'right' }}>{t.project.avgCol}</th>
                        <th css={jobTh} style={{ textAlign: 'right' }}>{t.project.runsCol}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {data.top_slow_jobs.map((j) => (
                        <tr key={j.job_name}>
                          <td css={jobTd}>{j.job_name}</td>
                          <td css={jobTd} style={{ textAlign: 'right', fontVariantNumeric: 'tabular-nums' }}>{fmtDur(j.avg_duration_sec)}</td>
                          <td css={jobTd} style={{ textAlign: 'right', color: 'var(--text-tertiary)' }}>{j.total_runs}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </>
              )}
            </div>
          </div>


          {data.flaky_jobs.length > 0 && (
            <div css={chartCard}>
              <div css={chartTitle}>{t.project.flakyJobs}</div>
              <table css={jobTable}>
                <thead>
                  <tr>
                    <th css={jobTh}>{t.project.jobCol}</th>
                    <th css={jobTh} style={{ textAlign: 'right' }}>{t.project.avgAttemptsCol}</th>
                    <th css={jobTh} style={{ textAlign: 'right' }}>{t.project.failRateCol}</th>
                    <th css={jobTh} style={{ textAlign: 'right' }}>{t.project.runsCol}</th>
                  </tr>
                </thead>
                <tbody>
                  {data.flaky_jobs.map((j) => (
                    <tr key={j.job_name}>
                      <td css={jobTd}>
                        <span style={{ marginRight: 6, fontSize: '0.75rem' }}>⚡</span>
                        {j.job_name}
                      </td>
                      <td css={jobTd} style={{ textAlign: 'right', fontVariantNumeric: 'tabular-nums', color: 'var(--status-running)' }}>
                        ×{j.avg_attempts.toFixed(2)}
                      </td>
                      <td css={jobTd} style={{ textAlign: 'right' }}>{fmtPct(j.failure_rate)}</td>
                      <td css={jobTd} style={{ textAlign: 'right', color: 'var(--text-tertiary)' }}>{j.total_runs}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}
    </div>
  )
}




export default function ProjectDetail() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const toast = useToast()
  const { t } = useI18n()
  const { sessionValidated } = useAppSelector((s) => s.auth)
  const [confirmDelete, setConfirmDelete] = useState(false)
  const [deleteProject, { isLoading: deleting }] = useDeleteProjectMutation()

  const {
    data: project,
    isLoading,
    isError,
    isFetching,
  } = useGetProjectQuery(id!, { skip: !id || !sessionValidated })

  const loading = isLoading || (isFetching && !project)

  useEffect(() => {
    if (!project) return
    trackProjectVisit({ id: project.id, name: project.name, repo_url: project.repo_url })
    document.title = `${project.name} — Flow`
    return () => { document.title = 'Flow — CI/CD' }
  }, [project?.id, project?.name, project?.repo_url])

  const handleDeleteConfirm = async () => {
    if (!id || !project) return
    try {
      await deleteProject(id).unwrap()
      toast.success(t.toast.projectDeleted.replace('{name}', project.name))
      navigate('/', { replace: true })
    } catch (e) {
      toast.error(e instanceof Error ? e.message : t.errors.failedToDeleteProject)
      setConfirmDelete(false)
    }
  }

  if (loading) {
    return <ProjectPageSkeleton />
  }

  if (isError || !project) {
    return (
      <>
        <Breadcrumb items={[{ label: t.nav.dashboard, to: '/' }]} />
        <p style={{ color: 'var(--danger)', marginBottom: 16 }}>{t.project.notFound}</p>
        <Link to="/" css={btnSecondary}>{t.project.backToDashboard}</Link>
      </>
    )
  }

  return (
    <>
      <ProjectView
        project={project}
        onDelete={() => setConfirmDelete(true)}
        onUpdated={() => {  }}
        deleting={deleting}
      />
      <ConfirmDialog
        open={confirmDelete}
        title={t.project.confirmDeleteTitle.replace('{name}', project?.name ?? '')}
        description={t.project.confirmDeleteDesc}
        confirmLabel={t.project.deleteProject}
        loading={deleting}
        onConfirm={handleDeleteConfirm}
        onCancel={() => setConfirmDelete(false)}
      />
    </>
  )
}
