import { apiFetch } from '~/services/http'
import type { ApiEnvelope, PlaceCreateInput, PlaceInfo } from '~/types/trip'

function getAccessToken(): string {
  if (!import.meta.client) {
    return ''
  }
  return localStorage.getItem('access_token') || ''
}

function getAuthHeaders(): Record<string, string> {
  const token = getAccessToken()
  if (!token) {
    return {}
  }
  return {
    Authorization: `Bearer ${token}`,
  }
}

export const placeService = {
  async create(payload: PlaceCreateInput): Promise<PlaceInfo> {
    const result = await apiFetch<ApiEnvelope<PlaceInfo>>('/place', {
      method: 'POST',
      headers: {
        ...getAuthHeaders(),
      },
      body: payload,
    })

    if (!result?.success || !result?.data) {
      throw new Error(result?.message || 'Cannot create place')
    }

    return result.data
  },
}
