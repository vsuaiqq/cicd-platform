

export interface RegisterRequest {
  email: string
  username?: string
  password: string
}

export interface RegisterResponse {
  success: boolean
  message: string
  user_id?: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user_id: string
  token_type: string
}

export interface ValidateTokenResponse {
  valid: boolean
  user_id?: string
  email?: string
  username?: string
  expires_at?: number
}

export interface RefreshTokenRequest {
  refresh_token: string
}

export interface RefreshTokenResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  token_type: string
}
