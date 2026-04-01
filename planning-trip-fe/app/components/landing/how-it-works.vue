<script setup lang="ts">
import {
  CalendarPlus2,
  Image as ImageIcon,
  MapPin,
  UserPlus,
} from 'lucide-vue-next'

const steps = [
  {
    icon: CalendarPlus2,
    title: 'Tạo chuyến đi',
    description: 'Đặt tên chuyến đi, thời gian và mô tả mục tiêu chuyến đi.',
  },
  {
    icon: MapPin,
    title: 'Thêm điểm đến',
    description: 'Thêm địa điểm, thứ tự tham quan và timeline theo ngày.',
  },
  {
    icon: UserPlus,
    title: 'Mời thành viên',
    description: 'Mời bạn bè vào nhóm để cùng góp ý và bình luận.',
  },
  {
    icon: ImageIcon,
    title: 'Lưu ảnh kỷ niệm',
    description: 'Sau chuyến đi, cả nhóm cùng upload ảnh vào album chung.',
  },
]
</script>

<template>
  <section
    id="how-it-works"
    class="relative isolate overflow-hidden py-20 sm:py-24 mt-1"
  >
    <div
      class="absolute inset-0 z-0 bg-[url('/images/planning-trip-step.png')] bg-cover bg-center"
    />
    <div
      class="absolute inset-0 z-0 bg-gradient-to-b from-[#EFE2BA]/76 via-[#EFE2BA]/70 to-[#C5CBE3]/62"
    />

    <div class="section-shell relative z-10">
      <div class="mx-auto max-w-3xl text-center">
        <h2 class="font-display text-4xl font-bold text-white sm:text-5xl">
          4 bước đơn giản để lên kế hoạch chuyến đi
        </h2>
      </div>

      <div class="mt-12 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <LandingInfoCard
          v-for="(step, idx) in steps"
          :key="step.title"
          :icon="step.icon"
          :title="step.title"
          :description="step.description"
          :prefix="String(idx + 1)"
          card-class="bg-white"
          icon-wrap-class="h-10 bg-[#C5CBE3]/45"
          icon-class="h-5 w-5"
          title-class="mt-4 text-base font-bold text-[#344B93]"
          description-class="mt-2 text-sm text-[#344B93] leading-relaxed"
        />
      </div>
    </div>
  </section>
</template>
