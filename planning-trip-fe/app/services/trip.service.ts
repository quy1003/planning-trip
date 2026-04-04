import { apiFetch } from '~/services/http'
import type { ApiEnvelope, CreateTripInput, TripDetail, TripSummary } from '~/types/trip'

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

export const tripService = {
  async list(): Promise<TripSummary[]> {
    const token = getAccessToken()
    if (!token) {
      return []
    }

    try {
      const result = await apiFetch<ApiEnvelope<TripSummary[]>>('/trip', {
        headers: {
          ...getAuthHeaders(),
        },
      })
      if (result?.success && Array.isArray(result.data)) {
        return result.data
      }
      return []
    } catch (error) {
      throw error
    }
  },

  async getById(id: string): Promise<TripDetail | null> {
    try {
      const result = await apiFetch<ApiEnvelope<TripDetail>>(`/trip/${id}`, {
        headers: {
          ...getAuthHeaders(),
        },
      })
      if (result?.success && result?.data?.id) {
        return result.data
      }
    } catch {
      return null
    }
    return null
  },

  async create(payload: CreateTripInput): Promise<TripDetail> {
    const result = await apiFetch<ApiEnvelope<TripDetail>>('/trip', {
      method: 'POST',
      headers: {
        ...getAuthHeaders(),
      },
      body: payload,
    })

    if (!result?.success || !result?.data) {
      throw new Error(result?.message || 'Cannot create trip')
    }

    return result.data
  },
}
