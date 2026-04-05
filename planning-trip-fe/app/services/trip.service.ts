import { apiFetch } from '~/services/http'
import type {
  ApiEnvelope,
  CreateScheduleItemInput,
  CreateTripInput,
  ReorderScheduleInput,
  ScheduleItem,
  TripDetail,
  TripSummary,
  UpdateScheduleItemInput,
} from '~/types/trip'

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

  async updateStatus(id: string, status: 'draft' | 'published'): Promise<TripDetail> {
    const result = await apiFetch<ApiEnvelope<TripDetail>>(`/trip/${id}/status`, {
      method: 'PATCH',
      headers: {
        ...getAuthHeaders(),
      },
      body: { status },
    })

    if (!result?.success || !result?.data) {
      throw new Error(result?.message || 'Cannot update trip status')
    }

    return result.data
  },

  async listScheduleItems(tripId: string, dayIndex?: number): Promise<ScheduleItem[]> {
    const query = typeof dayIndex === 'number' ? `?day_index=${dayIndex}` : ''
    const result = await apiFetch<ApiEnvelope<ScheduleItem[]>>(`/trip/${tripId}/schedule-items${query}`, {
      headers: {
        ...getAuthHeaders(),
      },
    })

    if (!result?.success || !Array.isArray(result.data)) {
      throw new Error(result?.message || 'Cannot list schedule items')
    }

    return result.data
  },

  async createScheduleItem(tripId: string, payload: CreateScheduleItemInput): Promise<ScheduleItem> {
    const result = await apiFetch<ApiEnvelope<ScheduleItem>>(`/trip/${tripId}/schedule-items`, {
      method: 'POST',
      headers: {
        ...getAuthHeaders(),
      },
      body: payload,
    })

    if (!result?.success || !result?.data) {
      throw new Error(result?.message || 'Cannot create schedule item')
    }

    return result.data
  },

  async updateScheduleItem(
    tripId: string,
    itemId: string,
    payload: UpdateScheduleItemInput,
  ): Promise<ScheduleItem> {
    const result = await apiFetch<ApiEnvelope<ScheduleItem>>(`/trip/${tripId}/schedule-items/${itemId}`, {
      method: 'PUT',
      headers: {
        ...getAuthHeaders(),
      },
      body: payload,
    })

    if (!result?.success || !result?.data) {
      throw new Error(result?.message || 'Cannot update schedule item')
    }

    return result.data
  },

  async deleteScheduleItem(tripId: string, itemId: string): Promise<void> {
    const result = await apiFetch<ApiEnvelope<Record<string, never>>>(`/trip/${tripId}/schedule-items/${itemId}`, {
      method: 'DELETE',
      headers: {
        ...getAuthHeaders(),
      },
    })

    if (!result?.success) {
      throw new Error(result?.message || 'Cannot delete schedule item')
    }
  },

  async reorderScheduleItems(tripId: string, payload: ReorderScheduleInput): Promise<void> {
    const result = await apiFetch<ApiEnvelope<Record<string, never>>>(`/trip/${tripId}/schedule-items/reorder`, {
      method: 'PATCH',
      headers: {
        ...getAuthHeaders(),
      },
      body: payload,
    })

    if (!result?.success) {
      throw new Error(result?.message || 'Cannot reorder schedule items')
    }
  },
}
