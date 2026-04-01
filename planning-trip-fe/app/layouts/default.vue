<script setup lang="ts">
import Header from '../components/layout/header.vue'
import Footer from '../components/layout/footer.vue'
</script>

<template>
  <div class="min-h-screen bg-[#EFE2BA]/35">
    <Header />
    <slot />
    <Footer />
  </div>
</template>
