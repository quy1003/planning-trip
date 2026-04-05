<script setup lang="ts">
import { FetchError } from 'ofetch'
import { MapPin, UserPlus } from 'lucide-vue-next'
import { authService } from '~/services/auth.service'

const router = useRouter()

const fullName = ref('')
const email = ref('')
const password = ref('')
const acceptedTerms = ref(false)
const isSubmitting = ref(false)
const status = ref<'success' | 'error' | ''>('')
const statusMessage = ref('')

async function onSubmit() {
  if (isSubmitting.value) {
    return
  }

  status.value = ''
  statusMessage.value = ''

  if (!fullName.value.trim() || !email.value.trim() || !password.value.trim()) {
    status.value = 'error'
    statusMessage.value = 'Đăng ký thất bại: vui lòng nhập đầy đủ thông tin.'
    return
  }

  if (password.value.trim().length < 6) {
    status.value = 'error'
    statusMessage.value = 'Đăng ký thất bại: mật khẩu tối thiểu 6 ký tự.'
    return
  }

  if (!acceptedTerms.value) {
    status.value = 'error'
    statusMessage.value = 'Đăng ký thất bại: vui lòng đồng ý điều khoản sử dụng.'
    return
  }

  isSubmitting.value = true

  try {
    const response = await authService.register({
      full_name: fullName.value.trim(),
      email: email.value.trim(),
      password: password.value,
    })

    if (response?.success && response?.data?.access_token) {
      localStorage.setItem('access_token', response.data.access_token)
      localStorage.setItem('current_user', JSON.stringify(response.data.user))
      await router.push('/trips')
      status.value = 'success'
      statusMessage.value = 'Đăng ký thành công.'
      return
    }

    status.value = 'error'
    statusMessage.value = 'Đăng ký thất bại.'
  } catch (error) {
    const fetchError = error as FetchError<{ message?: string }>
    status.value = 'error'
    statusMessage.value = fetchError?.data?.message
      ? `Đăng ký thất bại: ${fetchError.data.message}`
      : 'Đăng ký thất bại.'
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
              Bắt đầu khám phá
            </p>
            <h1 class="mt-3 font-display text-5xl font-bold text-white">
              Tạo tài khoản và bắt đầu lên kế hoạch chuyến đi của bạn ngay hôm nay
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

          <h2 class="text-3xl font-bold text-pt-orange">Đăng ký</h2>
          <p class="mt-2 text-sm text-black">Tạo tài khoản mới trong vài bước đơn giản</p>

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
              <label class="mb-1 block text-sm font-semibold text-pt-orange">Họ và tên</label>
              <input
                v-model="fullName"
                type="text"
                placeholder="Nguyễn Văn A"
                class="w-full rounded-xl border border-pt-lavender bg-white px-4 py-3 text-sm text-black outline-none transition focus:border-pt-orange"
              />
            </div>

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
              <label class="mb-1 block text-sm font-semibold text-pt-orange">Mật khẩu</label>
              <input
                v-model="password"
                type="password"
                placeholder="Tối thiểu 6 ký tự"
                class="w-full rounded-xl border border-pt-lavender bg-white px-4 py-3 text-sm text-black outline-none transition focus:border-pt-orange"
              />
            </div>

            <label class="flex items-start gap-2 text-sm text-black-light">
              <input
                v-model="acceptedTerms"
                type="checkbox"
                class="mt-0.5 rounded border border-pt-lavender"
              />
              <span>
                Tôi đồng ý với các điều khoản sử dụng và chính sách bảo mật của Planning Trip.
              </span>
            </label>

            <button
              type="submit"
              :disabled="isSubmitting || !acceptedTerms"
              class="inline-flex w-full items-center justify-center gap-2 rounded-xl bg-pt-orange px-4 py-3 text-sm font-semibold text-white shadow-sm transition hover:bg-pt-orange-dark"
              :class="isSubmitting || !acceptedTerms ? 'cursor-not-allowed opacity-70' : ''"
            >
              <UserPlus class="h-4 w-4" />
              {{ isSubmitting ? 'Đang tạo tài khoản...' : 'Tạo tài khoản' }}
            </button>
          </form>

          <p class="mt-5 text-center text-sm text-black-light">
            Đã có tài khoản?
            <NuxtLink to="/login" class="font-semibold text-pt-orange">
              Đăng nhập ngay
            </NuxtLink>
          </p>
        </div>
      </div>
    </div>
  </section>
</template>

