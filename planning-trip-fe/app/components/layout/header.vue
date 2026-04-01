<script setup lang="ts">
import { Menu, X, MapPin } from 'lucide-vue-next'

const mobileMenuOpen = ref(false)
const navItems = [
  { href: '#features', label: 'Tính năng' },
  { href: '#how-it-works', label: 'Cách hoạt động' },
  { href: '#testimonials', label: 'Đánh giá' },
  { href: '#contact', label: 'Liên hệ' },
]

function closeMenu() {
  mobileMenuOpen.value = false
}
</script>

<template>
  <header
    class="fixed inset-x-0 top-0 z-50 border-b border-[#C5CBE3] bg-white/80 backdrop-blur-md"
  >
    <div class="section-shell flex h-16 items-center justify-between">
      <NuxtLink to="/" class="flex items-center gap-2" @click="closeMenu">
        <span
          class="grid h-8 w-8 place-content-center rounded-full bg-pt-orange text-white"
        >
          <MapPin class="h-4 w-4" />
        </span>
        <span class="text-base font-bold text-pt-orange tracking-wider">
          Planning Trip</span
        >
      </NuxtLink>

      <nav class="hidden items-center gap-7 md:flex">
        <a
          v-for="item in navItems"
          :key="item.href"
          :href="item.href"
          class="text-sm font-medium !text-pt-orange transition"
        >
          {{ item.label }}
        </a>
      </nav>

      <div class="hidden items-center gap-2 md:flex">
        <NuxtLink
          to="/login"
          class="rounded-full px-4 py-2 text-sm font-medium text-pt-orange font-bold"
        >
          Đăng nhập
        </NuxtLink>
        <NuxtLink
          to="/register"
          class="rounded-full !bg-pt-orange px-4 py-2 text-sm font-semibold text-white"
        >
          Bắt đầu miễn phí
        </NuxtLink>
      </div>

      <button
        class="inline-flex h-9 w-9 items-center justify-center rounded-lg border border-pt-orange text-pt-orange md:hidden"
        type="button"
        @click="mobileMenuOpen = !mobileMenuOpen"
      >
        <Menu v-if="!mobileMenuOpen" class="h-5 w-5" />
        <X v-else class="h-5 w-5" />
      </button>
    </div>

    <div
      v-if="mobileMenuOpen"
      class="border-t border-pt-orange bg-white p-4 md:hidden"
    >
      <nav class="flex flex-col gap-1">
        <a
          v-for="item in navItems"
          :key="item.href"
          :href="item.href"
          class="rounded-lg px-3 py-2 text-sm font-medium !text-pt-orange"
          @click="closeMenu"
        >
          {{ item.label }}
        </a>
      </nav>
      <div class="mt-3 grid gap-2">
        <NuxtLink
          to="/login"
          class="rounded-lg border border-pt-orange px-4 py-2 text-center text-sm font-medium text-pt-orange font-bold"
          @click="closeMenu"
        >
          Đăng nhập
        </NuxtLink>
        <NuxtLink
          to="/register"
          class="rounded-lg bg-pt-orange px-4 py-2 text-center text-sm font-semibold text-white"
          @click="closeMenu"
        >
          Bắt đầu miễn phí
        </NuxtLink>
      </div>
    </div>
  </header>
</template>
