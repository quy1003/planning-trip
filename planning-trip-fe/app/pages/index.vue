<script setup lang="ts">
import Header from '~/components/layout/header.vue'
import Hero from '~/components/landing/hero.vue'
import Feature from '~/components/landing/feature.vue'
import HowItWorks from '~/components/landing/how-it-works.vue'
import Testimonials from '~/components/landing/testimonials.vue'
import Cta from '~/components/landing/cta.vue'
import Footer from '~/components/layout/footer.vue'

definePageMeta({
  layout: 'landing',
})
</script>

<template>
  <div>
    <Header />
    <Hero />
    <Feature />
    <HowItWorks />

    <Cta />
    <Testimonials />
    <Footer />
  </div>
</template>
