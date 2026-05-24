import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import { authApi } from '../api/auth'
import type { LoginRequest, RegisterRequest } from '../types/auth'

const ACCESS_KEY   = 'cicd_access_token'
const REFRESH_KEY  = 'cicd_refresh_token'
const USER_ID_KEY  = 'cicd_user_id'
const USERNAME_KEY = 'cicd_username'
const EMAIL_KEY    = 'cicd_email'

function getStored() {
  return {
    accessToken:  localStorage.getItem(ACCESS_KEY),
    refreshToken: localStorage.getItem(REFRESH_KEY),
    userId:       localStorage.getItem(USER_ID_KEY),
    username:     localStorage.getItem(USERNAME_KEY),
    email:        localStorage.getItem(EMAIL_KEY),
  }
}

function resolveUsername(
  username: string | undefined,
  email: string | undefined,
  fallback?: string | null,
): string {
  const u = username?.trim()
  if (u) return u
  const e = email?.trim()
  if (e) {
    const local = e.split('@')[0]?.trim()
    if (local) return local
  }
  return fallback?.trim() ?? ''
}

function setStored(
  access: string | null,
  refresh: string | null,
  userId: string | null | undefined,
  profile?: { username?: string | null; email?: string | null },
) {
  if (access)  localStorage.setItem(ACCESS_KEY, access)
  else         localStorage.removeItem(ACCESS_KEY)
  if (refresh) localStorage.setItem(REFRESH_KEY, refresh)
  else         localStorage.removeItem(REFRESH_KEY)

  if (userId !== undefined) {
    if (userId) localStorage.setItem(USER_ID_KEY, userId)
    else        localStorage.removeItem(USER_ID_KEY)
  }
  if (profile?.username !== undefined) {
    if (profile.username) localStorage.setItem(USERNAME_KEY, profile.username)
    else                  localStorage.removeItem(USERNAME_KEY)
  }
  if (profile?.email !== undefined) {
    if (profile.email) localStorage.setItem(EMAIL_KEY, profile.email)
    else               localStorage.removeItem(EMAIL_KEY)
  }
}



export const login = createAsyncThunk(
  'auth/login',
  async (body: LoginRequest, { rejectWithValue }) => {
    try {
      const res = await authApi.login(body)
      setStored(res.access_token, res.refresh_token, res.user_id)
      return res
    } catch (e) {
      return rejectWithValue(e instanceof Error ? e.message : 'Login failed')
    }
  }
)

export const register = createAsyncThunk(
  'auth/register',
  async (body: RegisterRequest, { rejectWithValue }) => {
    try {
      await authApi.register(body)
      const res = await authApi.login({ email: body.email, password: body.password })
      setStored(res.access_token, res.refresh_token, res.user_id)
      return res
    } catch (e) {
      return rejectWithValue(e instanceof Error ? e.message : 'Registration failed')
    }
  }
)

interface ValidatedSession {
  userId: string
  email: string
  username: string
  accessToken: string
  refreshToken: string
}

export const validateSession = createAsyncThunk<
  ValidatedSession,
  void,
  { state: { auth: AuthState }; rejectValue: string }
>(
  'auth/validateSession',
  async (_, { getState, rejectWithValue }) => {
    const { auth } = getState()
    const accessToken  = auth.accessToken  ?? getStored().accessToken
    const refreshToken = auth.refreshToken ?? getStored().refreshToken

    if (!accessToken) return rejectWithValue('No token')

    try {
      const res = await authApi.validateToken(accessToken)
      if (!res.valid) throw new Error('invalid')
      const email = res.email ?? auth.email ?? ''
      return {
        userId:       res.user_id ?? auth.userId ?? '',
        email,
        username:     resolveUsername(res.username, email, auth.username),
        accessToken,
        refreshToken: refreshToken ?? '',
      }
    } catch {

      if (!refreshToken) {
        setStored(null, null, null, { username: null, email: null })
        return rejectWithValue('Session expired')
      }
      try {
        const res = await authApi.refreshToken({ refresh_token: refreshToken })

        const existingUserId = auth.userId ?? getStored().userId ?? null
        setStored(res.access_token, res.refresh_token, undefined)
        return {
          userId:       existingUserId ?? '',
          email:        auth.email ?? '',
          username:     auth.username ?? '',
          accessToken:  res.access_token,
          refreshToken: res.refresh_token,
        }
      } catch {
        setStored(null, null, null, { username: null, email: null })
        return rejectWithValue('Session expired')
      }
    }
  }
)

export const logout = createAsyncThunk('auth/logout', async () => {
  setStored(null, null, null, { username: null, email: null })
})



export interface AuthState {
  accessToken:      string | null
  refreshToken:     string | null
  userId:           string | null
  email:            string | null
  username:         string | null
  isAuthenticated:  boolean
  sessionValidated: boolean
  loading:          boolean
  error:            string | null
}

const stored = getStored()

const initialState: AuthState = {
  accessToken:      stored.accessToken,
  refreshToken:     stored.refreshToken,
  userId:           stored.userId,
  email:            stored.email,
  username:         stored.username,
  isAuthenticated:  !!stored.accessToken,
  sessionValidated: false,
  loading:          false,
  error:            null,
}



const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    clearError(state) {
      state.error = null
    },
  },
  extraReducers: (builder) => {
    builder

      .addCase(login.pending, (state) => {
        state.loading = true
        state.error   = null
      })
      .addCase(login.fulfilled, (state, action) => {
        state.loading        = false
        state.error          = null
        state.accessToken    = action.payload.access_token
        state.refreshToken   = action.payload.refresh_token
        state.userId         = action.payload.user_id
        state.email          = action.meta.arg.email
        state.username       = resolveUsername(undefined, action.meta.arg.email, null)
        state.isAuthenticated = true
        state.sessionValidated = false
      })
      .addCase(login.rejected, (state, action) => {
        state.loading = false
        state.error   = (action.payload as string) ?? 'Login failed'
      })


      .addCase(register.pending, (state) => {
        state.loading = true
        state.error   = null
      })
      .addCase(register.fulfilled, (state, action) => {
        state.loading          = false
        state.error            = null
        state.accessToken      = action.payload.access_token
        state.refreshToken     = action.payload.refresh_token
        state.userId           = action.payload.user_id
        state.email            = action.meta.arg.email
        state.username         = action.meta.arg.username?.trim() || action.meta.arg.email
        state.isAuthenticated  = true
        state.sessionValidated = false
      })
      .addCase(register.rejected, (state, action) => {
        state.loading = false
        state.error   = (action.payload as string) ?? 'Registration failed'
      })


      .addCase(validateSession.pending, (state) => {
        if (!state.sessionValidated) state.loading = true
        state.error = null
      })
      .addCase(validateSession.fulfilled, (state, action) => {
        state.loading          = false
        state.error            = null
        state.accessToken      = action.payload.accessToken
        state.refreshToken     = action.payload.refreshToken
        if (action.payload.userId) state.userId = action.payload.userId
        if (action.payload.email)  state.email  = action.payload.email
        state.username = action.payload.username
        state.isAuthenticated  = true
        state.sessionValidated = true
        setStored(action.payload.accessToken, action.payload.refreshToken, action.payload.userId, {
          username: action.payload.username,
          email: action.payload.email || state.email,
        })
      })
      .addCase(validateSession.rejected, (state) => {
        state.loading          = false
        state.isAuthenticated  = false
        state.sessionValidated = true
        state.accessToken      = null
        state.refreshToken     = null
        state.userId           = null
        state.email            = null
        state.username         = null
      })


      .addCase(logout.fulfilled, (state) => {
        state.accessToken      = null
        state.refreshToken     = null
        state.userId           = null
        state.email            = null
        state.username         = null
        state.isAuthenticated  = false
        state.sessionValidated = true
        state.error            = null
      })
  },
})

export const { clearError } = authSlice.actions
export default authSlice.reducer
