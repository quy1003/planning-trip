<script setup lang="ts">
import type { TripSummary } from "~/types/trip"

defineProps<{
  trip: TripSummary
}>()
</script>

<template>
  <NuxtLink :to="`/trips/${trip.id}`" class="card">
    <h3>{{ trip.title }}</h3>
    <p>{{ trip.description || "No description" }}</p>
    <small>{{ trip.status }} | {{ trip.visibility }}</small>
  </NuxtLink>
</template>

<style scoped>
.card {
  display: grid;
  gap: 8px;
  border: 1px solid #dbe7f3;
  border-radius: 12px;
  background: #ffffff;
  padding: 16px;
  text-decoration: none;
}

h3 {
  margin: 0;
  color: #0f172a;
}

p {
  margin: 0;
  color: #475569;
}

small {
  color: #64748b;
}
</style>
