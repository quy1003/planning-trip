import { apiFetch } from '~/services/http'
import type { ApiEnvelope } from '~/types/trip'

interface UploadImageResult {
  url: string
  public_id: string
}

function getAccessToken(): string {
  if (!import.meta.client) {
    return ''
  }
  return localStorage.getItem('access_token') || ''
}

export const uploadService = {
  async uploadImage(file: File): Promise<UploadImageResult> {
    const token = getAccessToken()
    const formData = new FormData()
    formData.append('file', file)

    const result = await apiFetch<ApiEnvelope<UploadImageResult>>('/upload/image', {
      method: 'POST',
      headers: token
        ? {
            Authorization: `Bearer ${token}`,
          }
        : undefined,
      body: formData,
    })

    if (!result?.success || !result?.data?.url) {
      throw new Error(result?.message || 'Cannot upload image')
    }

    return result.data
  },
}
