export interface TripSummary {
  id: string
  title: string
  description?: string
  status: string
  visibility: string
}

export interface ApiEnvelope<T> {
  success: boolean
  data: T
  message?: string
  timestamp: string
}
