export interface TripSummary {
  id: string
  title: string
  description?: string
  status: string
  visibility: string
  cover_image_url?: string
  start_date?: string
  end_date?: string
}

export interface TripDetail extends TripSummary {
  creator_id: string
  start_date?: string
  end_date?: string
  cover_image_url?: string
  created_at: string
  updated_at: string
}

export interface CreateTripInput {
  title: string
  description?: string
  start_date?: string
  end_date?: string
  visibility?: string
  status?: string
  cover_image_url?: string
}

export interface ApiEnvelope<T> {
  success: boolean
  data: T
  message?: string
  timestamp: string
}
