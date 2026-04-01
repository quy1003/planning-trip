import type { ApiEnvelope, TripSummary } from "~/types/trip"
import { apiFetch } from "~/services/http"

const fallbackTrips: TripSummary[] = [
  {
    id: "demo-nha-trang",
    title: "Nha Trang 4N3D",
    description: "Demo trip for UI structure",
    status: "published",
    visibility: "group",
  },
]

export const tripService = {
  async list(): Promise<TripSummary[]> {
    try {
      const result = await apiFetch<ApiEnvelope<TripSummary[]>>("/trips")
      if (result?.success && Array.isArray(result.data)) {
        return result.data
      }
      return fallbackTrips
    } catch {
      return fallbackTrips
    }
  },

  async getById(id: string): Promise<TripSummary | null> {
    const list = await this.list()
    const matched = list.find((item) => item.id === id)
    return matched ?? null
  },
}
