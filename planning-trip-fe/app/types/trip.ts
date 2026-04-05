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

export interface TripMember {
  user_id: string
  full_name: string
  email: string
  avatar_url?: string
  role: string
  joined_at: string
}

export interface PlaceInfo {
  id: string
  name: string
  address?: string
  lat: number
  lng: number
}

export interface PlaceCreateInput {
  name: string
  address?: string
  lat: number
  lng: number
  google_place_id?: string
}

export interface TripPlaceInfo {
  id: string
  title?: string
  note?: string
  day_index: number
  order_index: number
  place: PlaceInfo
}

export interface TripPlaceCreateInput {
  trip_id: string
  place_id: string
  title?: string
  note?: string
  day_index: number
  order_index: number
}

export interface ScheduleItem {
  id: string
  trip_id: string
  trip_place_id?: string
  title: string
  description?: string
  start_time?: string
  end_time?: string
  day_index: number
  order_index: number
  created_by: string
  trip_place?: TripPlaceInfo
}

export interface AlbumPhoto {
  id: string
  url: string
  thumbnail_url?: string
  caption?: string
  taken_at?: string
  created_at: string
}

export interface TripDetail extends TripSummary {
  creator_id: string
  start_date?: string
  end_date?: string
  cover_image_url?: string
  created_at: string
  updated_at: string
  members: TripMember[]
  schedule: ScheduleItem[]
  album_preview: AlbumPhoto[]
}

export interface CreateScheduleItemInput {
  trip_place_id?: string
  title: string
  description?: string
  start_time?: string
  end_time?: string
  day_index: number
  order_index: number
}

export interface UpdateScheduleItemInput {
  trip_place_id?: string | null
  title?: string
  description?: string
  start_time?: string | null
  end_time?: string | null
  day_index?: number
  order_index?: number
}

export interface ReorderScheduleInput {
  items: Array<{
    id: string
    day_index: number
    order_index: number
  }>
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
