<script setup lang="ts">
import { FetchError } from 'ofetch'
import { LogIn, MapPin } from 'lucide-vue-next'
import { authService } from '~/services/auth.service'

const email = ref('')
const password = ref('')
const isSubmitting = ref(false)
const status = ref<'success' | 'error' | ''>('')
const statusMessage = ref('')

async function onSubmit() {
  if (isSubmitting.value) {
    return
  }

  status.value = ''
  statusMessage.value = ''

  if (!email.value.trim() || !password.value.trim()) {
    status.value = 'error'
    statusMessage.value = 'Đăng nhập thất bại: vui lòng nhập đầy đủ email và mật khẩu.'
    return
  }

  isSubmitting.value = true

  try {
    const response = await authService.login({
      email: email.value.trim(),
      password: password.value,
    })

    if (response?.success && response?.data?.access_token) {
      localStorage.setItem('access_token', response.data.access_token)
      localStorage.setItem('current_user', JSON.stringify(response.data.user))
      status.value = 'success'
      statusMessage.value = 'Đăng nhập thành công.'
      return
    }

    status.value = 'error'
    statusMessage.value = 'Đăng nhập thất bại.'
  } catch (error) {
    const fetchError = error as FetchError<{ message?: string }>
    status.value = 'error'
    statusMessage.value = fetchError?.data?.message
      ? `Đăng nhập thất bại: ${fetchError.data.message}`
      : 'Đăng nhập thất bại.'
  } finally {
    isSubmitting.value = false
  }
}
</script>

<template>
  <section class="relative isolate min-h-screen overflow-hidden">
    <div
      class="absolute inset-0 z-0 bg-[url('/images/planning-trip-carousel.jpeg')] bg-cover bg-center"
    />
    <div class="absolute inset-0 z-0 bg-[#1F2340]/50" />

    <div
      class="relative z-10 mx-auto flex min-h-screen w-full max-w-7xl items-start px-4 pb-16 pt-28 sm:px-6 lg:px-8"
    >
      <div class="grid w-full gap-8 lg:grid-cols-2">
        <div class="hidden lg:flex lg:flex-col lg:justify-between">
          <NuxtLink to="/" class="inline-flex items-center gap-2 self-start">
            <span class="text-lg font-bold tracking-wider text-pt-orange"></span>
          </NuxtLink>

          <div class="max-w-lg">
            <p class="text-md font-semibold uppercase tracking-wider text-pt-orange">
              Chào mừng trở lại
            </p>
            <h1 class="mt-3 font-display text-5xl font-bold text-white">
              Tiếp tục lên kế hoạch chuyến đi
            </h1>
          </div>
        </div>

        <div
          class="mx-auto w-full max-w-md rounded-3xl border border-white/60 bg-white/95 p-6 shadow-2xl backdrop-blur-sm sm:p-8"
        >
          <NuxtLink
            to="/"
            class="mb-4 inline-flex items-center gap-2 text-sm font-semibold text-pt-orange lg:hidden"
          >
            <MapPin class="h-4 w-4" />
            Planning Trip
          </NuxtLink>

          <h2 class="text-3xl font-bold text-pt-orange">Đăng nhập</h2>
          <p class="mt-2 text-sm text-black">
            Nhập thông tin tài khoản để tiếp tục lên kế hoạch chuyến đi của bạn
          </p>

          <div
            v-if="statusMessage"
            class="mt-4 rounded-xl border px-4 py-3 text-sm font-medium"
            :class="
              status === 'success'
                ? 'border-emerald-300 bg-emerald-50 text-emerald-700'
                : 'border-red-300 bg-red-50 text-red-700'
            "
          >
            {{ statusMessage }}
          </div>

          <form class="mt-6 space-y-4" @submit.prevent="onSubmit">
            <div>
              <label class="mb-1 block text-sm font-semibold text-pt-orange">Email</label>
              <input
                v-model="email"
                type="email"
                placeholder="you@example.com"
                class="w-full rounded-xl border border-pt-lavender bg-white px-4 py-3 text-sm text-black outline-none transition focus:border-pt-orange"
              />
            </div>

            <div>
              <div class="mb-1 flex items-center justify-between">
                <label class="block text-sm font-semibold text-pt-orange">Mật khẩu</label>
                <a href="#" class="text-xs font-semibold text-pt-orange">Quên mật khẩu?</a>
              </div>
              <input
                v-model="password"
                type="password"
                placeholder="••••••••"
                class="w-full rounded-xl border border-pt-lavender bg-white px-4 py-3 text-sm text-black outline-none transition focus:border-pt-orange"
              />
            </div>

            <button
              type="submit"
              :disabled="isSubmitting"
              class="inline-flex w-full items-center justify-center gap-2 rounded-xl bg-pt-orange px-4 py-3 text-sm font-semibold text-white shadow-sm transition hover:bg-pt-orange-dark"
              :class="isSubmitting ? 'cursor-not-allowed opacity-70' : ''"
            >
              <LogIn class="h-4 w-4" />
              {{ isSubmitting ? 'Đang đăng nhập...' : 'Đăng nhập' }}
            </button>
          </form>

          <p class="mt-5 text-center text-sm text-black-light">
            Chưa có tài khoản?
            <NuxtLink to="/register" class="font-semibold text-pt-orange">
              Tạo tài khoản mới
            </NuxtLink>
          </p>
        </div>
      </div>
    </div>
  </section>
</template>