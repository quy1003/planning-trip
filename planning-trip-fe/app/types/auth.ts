import type { ApiEnvelope } from '~/types/trip'

export interface AuthUser {
  id: string
  full_name: string
  email: string
  avatar_url?: string
  bio?: string
  created_at: string
  updated_at: string
}

export interface LoginResponseData {
  access_token: string
  token_type: string
  expires_in: number
  user: AuthUser
}

export type LoginResponse = ApiEnvelope<LoginResponseData>
export type RegisterResponse = ApiEnvelope<LoginResponseData>
