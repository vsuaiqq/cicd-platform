import type { BaseQueryFn } from '@reduxjs/toolkit/query'
import { apiRequest } from '../../api/client'
import type { RootState } from '../index'

export type ApiArgs =
  | string
  | {
      url: string
      method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE'
      body?: unknown
    }

export type ApiError = { message: string; status?: number }

export const apiBaseQuery: BaseQueryFn<ApiArgs, unknown, ApiError> = async (args, { getState }) => {
  const token = (getState() as RootState).auth.accessToken ?? undefined
  const url = typeof args === 'string' ? args : args.url
  const method = typeof args === 'string' ? 'GET' : (args.method ?? 'GET')
  const body = typeof args === 'string' ? undefined : args.body

  try {
    const data = await apiRequest(url, {
      method,
      body: body !== undefined && body !== null ? JSON.stringify(body) : undefined,
      token,
    })
    return { data }
  } catch (e) {
    return { error: { message: e instanceof Error ? e.message : 'Request failed' } }
  }
}
