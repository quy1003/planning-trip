import { apiFetch } from '~/services/http'
import type { LoginResponse, RegisterResponse } from '~/types/auth'

export interface LoginInput {
  email: string
  password: string
}

export interface RegisterInput {
  full_name: string
  email: string
  password: string
}

export const authService = {
  async login(payload: LoginInput): Promise<LoginResponse> {
    return apiFetch<LoginResponse>('/auth/login', {
      method: 'POST',
      body: payload,
    })
  },

  async register(payload: RegisterInput): Promise<RegisterResponse> {
    return apiFetch<RegisterResponse>('/auth/register', {
      method: 'POST',
      body: payload,
    })
  },
}
