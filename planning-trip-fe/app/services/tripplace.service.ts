import { apiFetch } from '~/services/http'
import type { ApiEnvelope, TripPlaceCreateInput, TripPlaceInfo } from '~/types/trip'

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

export const tripPlaceService = {
  async create(payload: TripPlaceCreateInput): Promise<TripPlaceInfo> {
    const result = await apiFetch<ApiEnvelope<TripPlaceInfo>>('/tripplace', {
      method: 'POST',
      headers: {
        ...getAuthHeaders(),
      },
      body: payload,
    })

    if (!result?.success || !result?.data) {
      throw new Error(result?.message || 'Cannot create trip place')
    }

    return result.data
  },
}
