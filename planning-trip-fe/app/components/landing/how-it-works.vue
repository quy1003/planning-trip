<script setup lang="ts">
import { Camera, MapPin, Plus, Users } from 'lucide-vue-next'

const steps = [
  {
    number: '01',
    icon: Plus,
    title: 'Tạo chuyến đi',
    description:
      'Đặt tên, chọn ngày và mô tả chuyến đi của bạn. Chỉ mất vài giây để bắt đầu.',
  },
  {
    number: '02',
    icon: MapPin,
    title: 'Thêm địa điểm',
    description:
      'Thêm các địa điểm muốn đến, sắp xếp theo thứ tự và xem trên bản đồ tương tác.',
  },
  {
    number: '03',
    icon: Users,
    title: 'Mời thành viên',
    description:
      'Mời bạn bè tham gia, cùng thảo luận và góp ý cho kế hoạch chuyến đi.',
  },
  {
    number: '04',
    icon: Camera,
    title: 'Lưu giữ kỷ niệm',
    description:
      'Sau chuyến đi, mọi người cùng đóng góp ảnh vào album chung để lưu giữ kỷ niệm.',
  },
]
</script>

<template>
  <section
    id="how-it-works"
    class="relative isolate overflow-hidden bg-background py-20 lg:py-32 mt-1"
  >
    <div
      class="absolute inset-0 z-0 bg-[url('/images/planning-trip-step.png')] bg-cover bg-center"
    />
    <div class="absolute inset-0 z-0 bg-background/85" />

    <div class="relative z-10 mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
      <div class="mx-auto mb-16 max-w-3xl text-center">
        <p
          class="mb-3 text-sm font-semibold uppercase tracking-wider text-primary"
        >
          Cách hoạt động
        </p>
        <h2
          class="text-balance font-roboto text-3xl font-bold text-white sm:text-4xl lg:text-5xl"
        >
          4 bước đơn giản để lên kế hoạch chuyến đi
        </h2>
        <p class="mt-4 text-lg text-white text-pretty">
          Không cần phức tạp, Planning Trip giúp bạn và nhóm lên kế hoạch nhanh
          chóng và hiệu quả.
        </p>
      </div>

      <div class="grid gap-8 md:grid-cols-2 lg:grid-cols-4">
        <div v-for="(step, index) in steps" :key="step.number" class="relative">
          <div
            v-if="index < steps.length - 1"
            class="absolute left-1/2 top-12 hidden h-0.5 w-full bg-border lg:block"
          />

          <div class="relative flex flex-col items-center text-center">
            <div class="relative z-10 mb-6">
              <div
                class="flex h-24 w-24 items-center justify-center rounded-2xl bg-white"
              >
                <component :is="step.icon" class="h-10 w-10 text-pt-orange" />
              </div>
              <div
                class="absolute -right-3 -top-3 z-20 flex h-8 w-8 items-center justify-center rounded-full bg-pt-orange text-white ring-2 ring-white shadow-sm"
              >
                <span class="text-xs font-bold leading-none">{{
                  step.number
                }}</span>
              </div>
            </div>

            <h3 class="mb-2 text-xl font-semibold text-white">
              {{ step.title }}
            </h3>
            <p class="text-white/90">{{ step.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>
