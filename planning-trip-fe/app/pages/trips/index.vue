<script setup lang="ts">
import type { TripSummary } from "~/types/trip"

const { data, pending, error } = await useTrips()

const trips = computed<TripSummary[]>(() => data.value ?? [])
</script>

<template>
  <section class="page">
    <h1>Trips</h1>

    <p v-if="pending">Loading trips...</p>
    <p v-else-if="error">Cannot load trips from API.</p>

    <div v-else class="grid">
      <TripCard v-for="trip in trips" :key="trip.id" :trip="trip" />
      <p v-if="trips.length === 0">No trips yet.</p>
    </div>
  </section>
</template>

<style scoped>
.page {
  display: grid;
  gap: 16px;
}

h1 {
  margin: 0;
  font-size: 28px;
}

.grid {
  display: grid;
  gap: 12px;
}
</style>
