import { useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useI18n } from '../i18n'
import { useAppDispatch, useAppSelector } from '../store'
import { register } from '../store/authSlice'
import { Alert } from '../components/ui'
import {
  form,
  formLabel,
  formInputAuth,
  btnAuthPrimary,
  authEyebrow,
  authPageTitle,
  authPageDesc,
  authFooter,
  authFooterLink,
  authPwdWrap,
  authPwdToggle,
  hintText,
} from '../styles/authPages'
import { IconEye, IconEyeOff } from '../components/auth/AuthIcons'

export default function Register() {
  const [email, setEmail] = useState('')
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [showPwd, setShowPwd] = useState(false)
  const dispatch = useAppDispatch()
  const { loading, error } = useAppSelector((s) => s.auth)
  const navigate = useNavigate()
  const { t } = useI18n()

  useEffect(() => {
    document.title = `${t.auth.createAccount} — Flow`
    return () => { document.title = 'Flow — CI/CD' }
  }, [t])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const result = await dispatch(register({ email, username: username || undefined, password }))
    if (register.fulfilled.match(result)) navigate('/', { replace: true })
  }

  return (
    <>
      <span css={authEyebrow}>{t.auth.eyebrowRegister}</span>
      <h1 css={authPageTitle}>{t.auth.createYourAccount}</h1>
      <p css={authPageDesc}>{t.auth.setupMinutes}</p>

      <form css={form} onSubmit={handleSubmit} noValidate>
        <label css={formLabel}>
          {t.auth.email}
          <input
            css={formInputAuth}
            type="email"
            autoComplete="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="you@company.com"
            required
            aria-invalid={!!error}
          />
        </label>
        <label css={formLabel}>
          {t.auth.displayName}
          <input
            css={formInputAuth}
            type="text"
            autoComplete="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            placeholder="e.g. Jane Smith"
          />
          <span css={hintText}>{t.auth.displayNameHint}</span>
        </label>
        <label css={formLabel}>
          {t.auth.password}
          <div css={authPwdWrap}>
            <input
              css={[formInputAuth, { paddingRight: 48, width: '100%' }]}
              type={showPwd ? 'text' : 'password'}
              autoComplete="new-password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              required
              aria-invalid={!!error}
            />
            <button
              type="button"
              css={authPwdToggle}
              onClick={() => setShowPwd(v => !v)}
              aria-label={showPwd ? t.auth.hidePassword : t.auth.showPassword}
            >
              {showPwd ? <IconEyeOff /> : <IconEye />}
            </button>
          </div>
        </label>
        {error && <Alert variant="error">{error}</Alert>}
        <button css={btnAuthPrimary} type="submit" disabled={loading} aria-busy={loading}>
          {loading ? t.auth.creatingAccount : t.auth.createAccount}
        </button>
      </form>

      <p css={authFooter}>
        {t.auth.alreadyHaveAccount}
        <Link to="/login" css={authFooterLink}>{t.auth.signIn}</Link>
      </p>
    </>
  )
}
