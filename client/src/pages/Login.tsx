import { useEffect, useState } from 'react'
import { Link, useNavigate, useLocation } from 'react-router-dom'
import { useAppDispatch, useAppSelector } from '../store'
import { login } from '../store/authSlice'
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
} from '../styles/authPages'
import { IconEye, IconEyeOff } from '../components/auth/AuthIcons'
import { useI18n } from '../i18n'

export default function Login() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPwd, setShowPwd] = useState(false)
  const dispatch = useAppDispatch()
  const { loading, error } = useAppSelector((s) => s.auth)
  const navigate = useNavigate()
  const location = useLocation()
  const from = (location.state as { from?: { pathname: string } })?.from?.pathname ?? '/'
  const { t } = useI18n()

  useEffect(() => {
    document.title = `${t.auth.signIn} — Flow`
    return () => { document.title = 'Flow — CI/CD' }
  }, [t])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const result = await dispatch(login({ email, password }))
    if (login.fulfilled.match(result)) navigate(from, { replace: true })
  }

  return (
    <>
      <span css={authEyebrow}>{t.auth.eyebrowLogin}</span>
      <h1 css={authPageTitle}>{t.auth.welcomeBack}</h1>
      <p css={authPageDesc}>{t.auth.signInDesc}</p>

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
          {t.auth.password}
          <div css={authPwdWrap}>
            <input
              css={[formInputAuth, { paddingRight: 48, width: '100%' }]}
              type={showPwd ? 'text' : 'password'}
              autoComplete="current-password"
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
          {loading ? t.auth.signingIn : t.auth.signIn}
        </button>
      </form>

      <p css={authFooter}>
        {t.auth.noAccount}
        <Link to="/register" css={authFooterLink}>{t.auth.createAccount}</Link>
      </p>
    </>
  )
}
