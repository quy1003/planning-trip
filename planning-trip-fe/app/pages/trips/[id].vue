<script setup lang="ts">
import type { TripDetail } from '~/types/trip'
import { tripService } from '~/services/trip.service'

const route = useRoute()
const tripId = computed(() => String(route.params.id || ''))

const { data, pending, error } = await useAsyncData('trip-detail-' + tripId.value, () =>
  tripService.getById(tripId.value),
)

const trip = computed<TripDetail | null>(() => data.value ?? null)
</script>

<template>
  <section class="page">
    <NuxtLink to="/trips">Back to Trips</NuxtLink>

    <p v-if="pending">Loading trip detail...</p>
    <p v-else-if="error">Cannot load this trip.</p>

    <article v-else-if="trip" class="card">
      <h1>{{ trip.title }}</h1>
      <p>{{ trip.description || 'No description' }}</p>
      <small>Status: {{ trip.status }} | Visibility: {{ trip.visibility }}</small>
    </article>

    <p v-else>Trip not found.</p>
  </section>
</template>

<style scoped>
.page {
  display: grid;
  gap: 14px;
  padding: 88px 16px 24px;
}

.card {
  border: 1px solid #dbe7f3;
  border-radius: 14px;
  background: #ffffff;
  padding: 18px;
}

h1 {
  margin: 0 0 8px;
}

p {
  color: #334155;
}

small {
  color: #64748b;
}
</style>
