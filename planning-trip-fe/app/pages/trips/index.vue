<script setup lang="ts">
import type { TripSummary } from '~/types/trip'

const { data, pending, error } = await useTrips()

const trips = computed<TripSummary[]>(() => data.value ?? [])
const isLoggedIn = computed(() => {
  if (!import.meta.client) {
    return false
  }
  return Boolean(localStorage.getItem('access_token'))
})
</script>

<template>
  <section class="page">
    <div class="head">
      <h1 class="font-bold text-pt-orange">Danh sách chuyến đi</h1>
      <NuxtLink to="/trips/new" class="create-btn">+ Create Trip</NuxtLink>
    </div>

    <p v-if="!isLoggedIn">Login để thấy danh sách chuyến đi của bạn.</p>
    <p v-else-if="pending">Đang tải danh sách chuyến đi...</p>
    <p v-else-if="error">Không thể tải danh sách chuyến đi</p>

    <div v-else-if="isLoggedIn" class="grid">
      <TripCard v-for="trip in trips" :key="trip.id" :trip="trip" />
      <p v-if="trips.length === 0">Chưa có chuyến đi nào</p>
    </div>
  </section>
</template>

<style scoped>
.page {
  display: grid;
  gap: 16px;
  padding: 88px 16px 24px;
}

.head {
  align-items: center;
  display: flex;
  gap: 12px;
  justify-content: space-between;
}

h1 {
  margin: 0;
  font-size: 28px;
}

.create-btn {
  background: #e8883f;
  border-radius: 999px;
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  padding: 8px 14px;
  text-decoration: none;
}

.grid {
  display: grid;
  gap: 12px;
}
</style>
