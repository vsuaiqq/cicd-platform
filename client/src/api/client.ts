const getBaseUrl = (): string => {
  const url = import.meta.env.VITE_API_URL
  if (url) return url.replace(/\/$/, '')
  return ''
}

export const apiBase = getBaseUrl()

export interface RequestConfig extends RequestInit {
  token?: string | null
}

export async function apiRequest<T>(
  path: string,
  config: RequestConfig = {}
): Promise<T> {
  const { token, ...init } = config
  const url = apiBase ? `${apiBase}${path}` : path

  const headers: Record<string, string> = {
    ...(init.headers as Record<string, string>),
  }

  if (init.body != null) {
    headers['Content-Type'] = 'application/json'
  }
  if (token) {
    headers['Authorization'] = `Bearer ${token}`
  }

  const res = await fetch(url, { ...init, headers })


  if (res.status === 204) {
    return undefined as T
  }

  const data = await res.json().catch(() => null)
  if (!res.ok) {
    const message =
      (data as { message?: string } | null)?.message ??
      (data as { error?: string } | null)?.error ??
      res.statusText
    throw new Error(message)
  }
  return data as T
}
