<script setup lang="ts">
import { Calendar, Image as ImageIcon } from 'lucide-vue-next'

const stats = [
  { value: '10K+', label: 'Chuyến đi đã tạo' },
  { value: '50K+', label: 'Người dùng' },
  { value: '4.9', label: 'Điểm người dùng' },
]

const stripItems = [
  'Bản đồ tương tác',
  'Lịch trình chi tiết',
  'Album ảnh chung',
  'Reaction tức thì',
]
</script>

<template>
  <section
    class="relative isolate overflow-hidden pb-14 pt-16 lg:min-h-[calc(100vh)]"
  >
    <div
      class="absolute inset-0 z-0 bg-[url('/images/planning-trip-carousel.jpeg')] bg-cover bg-center"
    />
    <div
      class="absolute inset-0 z-0 bg-gradient-to-b from-[#EFE2BA]/38 via-[#EFE2BA]/50 to-[#C5CBE3]/58"
    />

    <div
      class="section-shell relative z-10 grid gap-8 pb-8 pt-8 lg:grid-cols-2 lg:items-center lg:gap-12 lg:pb-6 lg:pt-6"
    >
      <div>
        <span
          class="inline-flex rounded-full border border-[#C5CBE3] bg-white/80 px-3 py-1 text-xs font-semibold text-pt-ink"
        >
          Lên kế hoạch chuyến đi dễ dàng
        </span>

        <h1
          class="font-display mt-4 text-5xl font-bold leading-[1.04] !text-white sm:text-6xl lg:text-7xl"
        >
          Lên kế hoạch
          <br />
          chuyến đi
          <span> cùng nhau</span>
        </h1>

        <p
          class="mt-4 max-w-xl text-base leading-relaxed text-white sm:text-lg"
        >
          Tạo itinerary theo nhóm, hiển thị điểm đến trên bản đồ và cộng tác
          theo thời gian thực qua bình luận, cảm xúc, ảnh chia sẻ.
        </p>

        <div class="mt-7 flex flex-wrap gap-3">
          <button
            class="rounded-[16px] !bg-[#344B93] px-6 py-3 text-sm font-semibold text-white shadow"
          >
            Tạo chuyến đi đầu tiên
          </button>
          <button
            class="rounded-[16px] bg-pt-orange px-6 py-3 text-sm font-semibold text-white"
          >
            Xem demo
          </button>
        </div>

        <div class="mt-6 grid max-w-md grid-cols-3 gap-3">
          <article
            v-for="item in stats"
            :key="item.label"
            class="rounded-xl border border-[#C5CBE3] bg-white/95 p-3 text-center shadow-sm"
          >
            <p class="text-2xl font-extrabold text-[#1F2340]">
              {{ item.value }}
            </p>
            <p class="mt-1 text-xs text-[#454B70]">{{ item.label }}</p>
          </article>
        </div>
      </div>

      <div class="relative">
        <article
          class="rounded-3xl border border-[#C5CBE3] bg-white/95 p-4 shadow-2xl backdrop-blur-sm sm:p-5"
        >
          <header
            class="flex items-center justify-between rounded-2xl bg-[#EFE2BA]/45 p-3"
          >
            <div>
              <p class="text-sm font-semibold text-[#1F2340]">
                Nha Trang Adventure
              </p>
              <p class="text-xs text-[#656B8D]">3 ngày • 5 thành viên</p>
            </div>
            <div class="flex -space-x-2">
              <span
                class="grid h-7 w-7 place-content-center rounded-full bg-white text-[11px] font-semibold ring-2 ring-[#C5CBE3]"
                >Q</span
              >
              <span
                class="grid h-7 w-7 place-content-center rounded-full bg-white text-[11px] font-semibold ring-2 ring-[#C5CBE3]"
                >L</span
              >
              <span
                class="grid h-7 w-7 place-content-center rounded-full bg-white text-[11px] font-semibold ring-2 ring-[#C5CBE3]"
                >M</span
              >
              <span
                class="grid h-7 w-7 place-content-center rounded-full bg-white text-[11px] font-semibold ring-2 ring-[#C5CBE3]"
                >+2</span
              >
            </div>
          </header>

          <div
            class="mt-3 rounded-2xl border border-[#C5CBE3] bg-gradient-to-br from-[#C5CBE3]/65 via-[#EFE2BA]/55 to-[#D79922]/25 p-3"
          >
            <div
              class="relative overflow-hidden rounded-xl border border-white/70 bg-white/70"
            >
              <img
                src="/images/planning-trip-map.png"
                alt="Bản đồ kế hoạch chuyến đi"
                class="block h-auto w-full max-h-[42vh] bg-contain lg:max-h-[40vh]"
                loading="lazy"
              />
            </div>
          </div>

          <div class="mt-3 grid gap-2 sm:grid-cols-2">
            <div class="flex items-center gap-2 rounded-xl bg-[#EFE2BA]/35 p-3">
              <Calendar class="h-4 w-4 text-[#4056A1]" />
              <div>
                <p class="text-xs font-semibold text-pt-ink">Tháp Bà Ponagar</p>
                <p class="text-[11px] text-[#656B8D]">Ngày 1 • 08:00 - 10:00</p>
              </div>
            </div>
            <div class="flex items-center gap-2 rounded-xl bg-[#EFE2BA]/35 p-3">
              <ImageIcon class="h-4 w-4 text-[#F13C20]" />
              <div>
                <p class="text-xs font-semibold text-pt-ink">VinWonders</p>
                <p class="text-[11px] text-[#656B8D]">Ngày 1 • 14:00 - 21:00</p>
              </div>
            </div>
          </div>
        </article>
      </div>
    </div>
  </section>
</template>
