<script setup lang="ts">
import {
  CalendarClock,
  Heart,
  Image as ImageIcon,
  MapPinned,
  MessageCircle,
  Route,
  ShieldCheck,
  Users,
} from 'lucide-vue-next'

const features = [
  {
    icon: MapPinned,
    title: 'Bản đồ tương tác',
    description: 'Thêm điểm đến và xem lộ trình trực quan trên bản đồ.',
  },
  {
    icon: CalendarClock,
    title: 'Lịch trình chi tiết',
    description:
      'Sắp xếp kế hoạch theo ngày/giờ và điều chỉnh nhanh trong nhóm.',
  },
  {
    icon: MessageCircle,
    title: 'Bình luận theo điểm',
    description:
      'Trao đổi trực tiếp ở từng địa điểm hoặc toàn bộ plan chuyến đi.',
  },
  {
    icon: Heart,
    title: 'Reaction cảm xúc',
    description: 'Thả tim, like, wow cho địa điểm, ảnh hoặc bình luận.',
  },
  {
    icon: ImageIcon,
    title: 'Album chung',
    description: 'Mọi thành viên cùng đóng góp ảnh, viết caption và lưu ký ức.',
  },
  {
    icon: Users,
    title: 'Quản lý thành viên',
    description: 'Phân quyền user linh hoạt theo từng trip.',
  },
  {
    icon: Route,
    title: 'Gợi ý lộ trình',
    description: 'Xem thứ tự hợp lý giữa các điểm đến để tối ưu di chuyển.',
  },
  {
    icon: ShieldCheck,
    title: 'Lịch sử thay đổi',
    description: 'Dễ theo dõi ai cập nhật lịch trình và nội dung vừa thay đổi.',
  },
]
</script>

<template>
  <section
    id="features"
    class="relative isolate overflow-hidden py-20 sm:py-24 mt-1"
  >
    <div
      class="absolute inset-0 z-0 bg-[url('/images/planning-trip-feature.png')] bg-cover bg-center"
    />
    <div
      class="absolute inset-0 z-0 bg-gradient-to-b from-[#EFE2BA]/34 via-[#EFE2BA]/45 to-[#C5CBE3]/58"
    />

    <div class="section-shell relative z-10">
      <div class="mx-auto max-w-3xl text-center">
        <h2 class="font-roboto text-balance text-4xl font-bold text-white sm:text-5xl">
          Mọi thứ bạn cần cho một chuyến đi hoàn hảo
        </h2>
      </div>

      <div class="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <LandingInfoCard
          v-for="item in features"
          :key="item.title"
          :icon="item.icon"
          :title="item.title"
          :description="item.description"
        />
      </div>
    </div>
  </section>
</template>
