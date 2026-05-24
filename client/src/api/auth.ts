import { apiRequest } from './client'
import type {
  RegisterRequest,
  RegisterResponse,
  LoginRequest,
  LoginResponse,
  ValidateTokenResponse,
  RefreshTokenRequest,
  RefreshTokenResponse,
} from '../types/auth'

const AUTH_PREFIX = '/api/v1/auth'

export const authApi = {
  register(body: RegisterRequest) {
    return apiRequest<RegisterResponse>(`${AUTH_PREFIX}/register`, {
      method: 'POST',
      body: JSON.stringify(body),
    })
  },

  login(body: LoginRequest) {
    return apiRequest<LoginResponse>(`${AUTH_PREFIX}/login`, {
      method: 'POST',
      body: JSON.stringify(body),
    })
  },

  validateToken(token: string) {
    return apiRequest<ValidateTokenResponse>(`${AUTH_PREFIX}/validate`, {
      method: 'POST',
      body: JSON.stringify({ token }),
    })
  },

  refreshToken(body: RefreshTokenRequest) {
    return apiRequest<RefreshTokenResponse>(`${AUTH_PREFIX}/refresh`, {
      method: 'POST',
      body: JSON.stringify(body),
    })
  },
}
