import { tripService } from "~/services/trip.service"
import type { TripSummary } from "~/types/trip"

export async function useTrips() {
  return useAsyncData<TripSummary[]>("trips", () => tripService.list(), {
    server: false,
    default: () => [],
  })
}
