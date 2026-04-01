export async function apiFetch<T>(path: string, options: Parameters<typeof $fetch<T>>[1] = {}): Promise<T> {
  const config = useRuntimeConfig()

  return $fetch<T>(path, {
    baseURL: config.public.apiBase,
    ...options,
  })
}
