<script setup lang="ts">
import { Star } from 'lucide-vue-next'

const items = [
  {
    name: 'Quý Nguyễn',
    role: 'Trip leader',
    text: 'Nhóm mình lên plan Nha Trang cực nhanh, ai cũng theo được timeline rõ ràng.',
  },
  {
    name: 'Minh Anh',
    role: 'Family planner',
    text: 'Mình thích nhất phần map + bình luận theo điểm đến, trao đổi rất tiện.',
  },
  {
    name: 'Linh Chi',
    role: 'Solo traveler',
    text: 'Dùng để lưu kế hoạch và album ảnh sau chuyến đi, xem lại rất sướng.',
  },
  {
    name: 'Đức Anh',
    role: 'Couple travel',
    text: 'UI sáng sủa, mobile cũng ổn, thêm địa điểm và sửa lịch trình rất mượt.',
  },
  {
    name: 'Thu Hà',
    role: 'Team building',
    text: 'Trip đông người nhưng vẫn quản lý được nhờ phân quyền thành viên rõ ràng.',
  },
  {
    name: 'Hoàng Long',
    role: 'Phượt thủ',
    text: 'Đã thay thế hẳn file excel trước đây, giờ mọi thứ tập trung một nơi.',
  },
]
</script>

<template>
  <section id="testimonials" class="bg-white py-20 sm:py-24">
    <div class="section-shell">
      <div class="mx-auto max-w-3xl text-center">
        <h2 class="font-roboto text-balance text-pt-orange text-4xl font-bold text-pt-ink sm:text-5xl">
          Được yêu thích bởi hàng nghìn người dùng
        </h2>
      </div>

      <div class="mt-10 grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <article
          v-for="item in items"
          :key="item.name + item.role"
          class="rounded-2xl border border-[#C5CBE3] bg-[#EFE2BA]/35 p-5"
        >
          <div class="mb-3 flex gap-1 text-pt-gold">
            <Star v-for="n in 5" :key="n" class="h-4 w-4 fill-current" />
          </div>
          <p class="text-sm leading-relaxed text-black">"{{ item.text }}"</p>
          <div class="mt-4 border-t border-[#C5CBE3] pt-3">
            <p class="text-sm font-bold text-pt-orange">{{ item.name }}</p>
            <p class="text-xs text-pt-orange">{{ item.role }}</p>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>
