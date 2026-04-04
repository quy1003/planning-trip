<script setup lang="ts">
import type { TripSummary } from "~/types/trip"

defineProps<{
  trip: TripSummary
}>()

function formatDate(value?: string): string {
  if (!value) return "Not set"
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return date.toLocaleDateString("en-GB")
}

function statusClass(status?: string): string {
  switch ((status || "").toLowerCase()) {
    case "published":
      return "tag-status-published"
    case "draft":
      return "tag-status-draft"
    default:
      return "tag-default"
  }
}

function visibilityClass(visibility?: string): string {
  switch ((visibility || "").toLowerCase()) {
    case "private":
      return "tag-visibility-private"
    case "group":
      return "tag-visibility-group"
    case "public":
      return "tag-visibility-public"
    default:
      return "tag-default"
  }
}
</script>

<template>
  <NuxtLink :to="`/trips/${trip.id}`" class="card">
    <div class="thumb-wrap">
      <img
        v-if="trip.cover_image_url"
        :src="trip.cover_image_url"
        :alt="trip.title"
        class="thumb"
      >
      <div v-else class="thumb fallback">No cover</div>
    </div>
    <h3>{{ trip.title }}</h3>
    <p>{{ trip.description || "No description" }}</p>

    <div class="dates">
      <span>{{ formatDate(trip.start_date) }}</span>
      <span class="dot">-</span>
      <span>{{ formatDate(trip.end_date) }}</span>
    </div>

    <div class="tags">
      <span class="tag" :class="statusClass(trip.status)">{{ trip.status }}</span>
      <span class="tag" :class="visibilityClass(trip.visibility)">{{ trip.visibility }}</span>
    </div>
  </NuxtLink>
</template>

<style scoped>
.card {
  display: grid;
  gap: 10px;
  border: 1px solid #dbe7f3;
  border-radius: 12px;
  background: #ffffff;
  padding: 12px;
  text-decoration: none;
}

.thumb-wrap {
  border-radius: 10px;
  overflow: hidden;
}

.thumb {
  width: 100%;
  height: 148px;
  object-fit: cover;
  display: block;
}

.fallback {
  align-items: center;
  background: linear-gradient(120deg, #f8efe1, #f3d7bf);
  color: #9a5c2a;
  display: flex;
  font-size: 13px;
  font-weight: 600;
  justify-content: center;
  letter-spacing: 0.02em;
}

h3 {
  margin: 0;
  color: #0f172a;
}

p {
  margin: 0;
  color: #475569;
}

.dates {
  color: #64748b;
  display: flex;
  font-size: 13px;
  font-weight: 600;
  gap: 6px;
}

.dot {
  opacity: 0.6;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag {
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  padding: 4px 10px;
  text-transform: lowercase;
}

.tag-default {
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  color: #334155;
}

.tag-status-draft {
  background: #f8fafc;
  border: 1px solid #cbd5e1;
  color: #475569;
}

.tag-status-published {
  background: #ecfdf5;
  border: 1px solid #86efac;
  color: #166534;
}

.tag-visibility-private {
  background: #fff7ed;
  border: 1px solid #fdba74;
  color: #9a3412;
}

.tag-visibility-group {
  background: #eef2ff;
  border: 1px solid #a5b4fc;
  color: #3730a3;
}

.tag-visibility-public {
  background: #eff6ff;
  border: 1px solid #93c5fd;
  color: #1e40af;
}
</style>
