<script setup lang="ts">
import { CalendarDays, Clock3, MapPin, Users, Plus, Save, Send, Pencil, Trash2, ArrowUp, ArrowDown } from 'lucide-vue-next'
import { FetchError } from 'ofetch'
import { placeService } from '~/services/place.service'
import { tripService } from '~/services/trip.service'
import { tripPlaceService } from '~/services/tripplace.service'
import type { PlaceInfo, ScheduleItem, TripDetail } from '~/types/trip'

const route = useRoute()
const tripId = computed(() => String(route.params.id || ''))

const isLoggedIn = computed(() => {
  if (!import.meta.client) {
    return false
  }
  return Boolean(localStorage.getItem('access_token'))
})

const statusMessage = ref('')
const statusType = ref<'success' | 'error' | ''>('')
const isSavingDraft = ref(false)
const isPublishing = ref(false)
const isCreatingItem = ref(false)
const isReordering = ref(false)
const updatingItemId = ref<string | null>(null)
const deletingItemId = ref<string | null>(null)
const selectedDay = ref(1)
const editingItemId = ref<string | null>(null)

const createItemForm = reactive({
  title: '',
  description: '',
  start_time: '',
  end_time: '',
  place_name: '',
  place_address: '',
})

const editItemForm = reactive({
  title: '',
  description: '',
  start_time: '',
  end_time: '',
})

const pickedLocation = ref<{ lat: number; lng: number } | null>(null)
const placeSearchQuery = ref('')
const placeSearchLoading = ref(false)
const placeSearchResults = ref<Array<{
  name: string
  address: string
  lat: number
  lng: number
}>>([])

const { data, pending, error, refresh } = await useAsyncData(
  () => `trip-detail-${tripId.value}`,
  () => tripService.getById(tripId.value),
  {
    watch: [tripId],
  },
)

const trip = computed<TripDetail | null>(() => data.value ?? null)

function setStatus(type: 'success' | 'error', message: string) {
  statusType.value = type
  statusMessage.value = message
}

function clearStatus() {
  statusType.value = ''
  statusMessage.value = ''
}

function parseDate(value?: string): Date | null {
  if (!value) return null
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return null
  return date
}

function formatDate(value?: string): string {
  const date = parseDate(value)
  if (!date) return 'Not updated'
  return new Intl.DateTimeFormat('vi-VN').format(date)
}

function formatTime(value?: string): string {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  return new Intl.DateTimeFormat('vi-VN', {
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

function toLocalDateTimeInput(value?: string): string {
  if (!value) return ''
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return ''
  const tzOffset = date.getTimezoneOffset() * 60000
  return new Date(date.getTime() - tzOffset).toISOString().slice(0, 16)
}

function toISOStringFromLocalInput(value: string): string | undefined {
  if (!value.trim()) return undefined
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return undefined
  return date.toISOString()
}

function capitalize(value?: string): string {
  if (!value) return 'Unknown'
  return value.charAt(0).toUpperCase() + value.slice(1)
}

function roleLabel(role: string): string {
  const normalized = role.toLowerCase()
  if (normalized === 'owner') return 'Owner'
  if (normalized === 'editor') return 'Editor'
  return capitalize(normalized)
}

const totalDays = computed(() => {
  const scheduleMax = Math.max(
    1,
    ...(trip.value?.schedule?.map(item => item.day_index) || [1]),
  )
  const start = parseDate(trip.value?.start_date)
  const end = parseDate(trip.value?.end_date)
  if (!start || !end) return scheduleMax
  const diff = Math.floor((end.getTime() - start.getTime()) / (1000 * 60 * 60 * 24)) + 1
  return Math.max(scheduleMax, diff > 0 ? diff : 1)
})

const dayOptions = computed(() => Array.from({ length: totalDays.value }, (_, index) => index + 1))

watch(
  dayOptions,
  value => {
    if (!value.includes(selectedDay.value)) {
      selectedDay.value = value[0] || 1
    }
  },
  { immediate: true },
)

const selectedDaySchedule = computed<ScheduleItem[]>(() => {
  const source = trip.value?.schedule || []
  return source
    .filter(item => item.day_index === selectedDay.value)
    .sort((a, b) => {
      if (a.order_index !== b.order_index) return a.order_index - b.order_index
      const aTime = a.start_time ? new Date(a.start_time).getTime() : 0
      const bTime = b.start_time ? new Date(b.start_time).getTime() : 0
      return aTime - bTime
    })
})

const mapPoints = computed(() => {
  return (trip.value?.schedule || [])
    .filter(item => item.trip_place?.place)
    .map(item => ({
      id: item.id,
      title: item.title,
      lat: item.trip_place!.place.lat,
      lng: item.trip_place!.place.lng,
      day_index: item.day_index,
    }))
})

async function refreshTrip() {
  await refresh()
}

async function updateTripStatus(status: 'draft' | 'published') {
  if (!trip.value) return
  clearStatus()

  if (!isLoggedIn.value) {
    setStatus('error', 'You need to login to update trip status.')
    return
  }

  if (status === 'draft') {
    isSavingDraft.value = true
  } else {
    isPublishing.value = true
  }

  try {
    await tripService.updateStatus(trip.value.id, status)
    await refreshTrip()
    setStatus('success', status === 'draft' ? 'Draft saved.' : 'Trip published.')
  } catch (error) {
    const fetchError = error as FetchError<{ message?: string }>
    setStatus('error', fetchError?.data?.message || fetchError?.message || 'Cannot update trip status.')
  } finally {
    isSavingDraft.value = false
    isPublishing.value = false
  }
}

async function createScheduleItem() {
  if (!trip.value) return
  clearStatus()

  const title = createItemForm.title.trim()
  if (!title) {
    setStatus('error', 'Please enter schedule title.')
    return
  }

  isCreatingItem.value = true
  try {
    const orderIndex = selectedDaySchedule.value.length + 1
    let createdPlace: PlaceInfo | null = null
    let createdTripPlaceId: string | undefined

    if (pickedLocation.value) {
      const placeName = createItemForm.place_name.trim()
      if (!placeName) {
        throw new Error('Please enter place name for selected map point.')
      }

      createdPlace = await placeService.create({
        name: placeName,
        address: createItemForm.place_address.trim() || undefined,
        lat: Number(pickedLocation.value.lat.toFixed(7)),
        lng: Number(pickedLocation.value.lng.toFixed(7)),
      })

      const createdTripPlace = await tripPlaceService.create({
        trip_id: trip.value.id,
        place_id: createdPlace.id,
        title: placeName,
        note: createItemForm.description.trim() || undefined,
        day_index: selectedDay.value,
        order_index: orderIndex,
      })
      createdTripPlaceId = createdTripPlace.id
    }

    await tripService.createScheduleItem(trip.value.id, {
      title,
      description: createItemForm.description.trim() || undefined,
      start_time: toISOStringFromLocalInput(createItemForm.start_time),
      end_time: toISOStringFromLocalInput(createItemForm.end_time),
      day_index: selectedDay.value,
      order_index: orderIndex,
      trip_place_id: createdTripPlaceId,
    })

    createItemForm.title = ''
    createItemForm.description = ''
    createItemForm.start_time = ''
    createItemForm.end_time = ''
    createItemForm.place_name = ''
    createItemForm.place_address = ''
    pickedLocation.value = null

    await refreshTrip()
    setStatus('success', 'Schedule item created.')
  } catch (error) {
    const fetchError = error as FetchError<{ message?: string }>
    setStatus('error', fetchError?.data?.message || fetchError?.message || 'Cannot create schedule item.')
  } finally {
    isCreatingItem.value = false
  }
}

function onMapPointSelected(payload: { lat: number; lng: number }) {
  pickedLocation.value = payload
}

async function searchPlaces() {
  const query = placeSearchQuery.value.trim()
  if (query.length < 2) {
    placeSearchResults.value = []
    return
  }

  placeSearchLoading.value = true
  try {
    const result = await $fetch<Array<{
      display_name: string
      lat: string
      lon: string
      name?: string
    }>>('https://nominatim.openstreetmap.org/search', {
      query: {
        q: query,
        format: 'jsonv2',
        limit: 7,
        addressdetails: 1,
      },
      headers: {
        Accept: 'application/json',
      },
    })

    placeSearchResults.value = (result || [])
      .map(row => ({
        name: row.name || row.display_name.split(',')[0] || 'Selected place',
        address: row.display_name,
        lat: Number(row.lat),
        lng: Number(row.lon),
      }))
      .filter(row => Number.isFinite(row.lat) && Number.isFinite(row.lng))
  } catch {
    placeSearchResults.value = []
    setStatus('error', 'Cannot search place right now.')
  } finally {
    placeSearchLoading.value = false
  }
}

function chooseSearchPlace(place: { name: string; address: string; lat: number; lng: number }) {
  pickedLocation.value = {
    lat: place.lat,
    lng: place.lng,
  }
  createItemForm.place_name = place.name
  createItemForm.place_address = place.address
  placeSearchResults.value = []
}

function startEditScheduleItem(item: ScheduleItem) {
  editingItemId.value = item.id
  editItemForm.title = item.title
  editItemForm.description = item.description || ''
  editItemForm.start_time = toLocalDateTimeInput(item.start_time)
  editItemForm.end_time = toLocalDateTimeInput(item.end_time)
}

function cancelEditScheduleItem() {
  editingItemId.value = null
  editItemForm.title = ''
  editItemForm.description = ''
  editItemForm.start_time = ''
  editItemForm.end_time = ''
}

async function saveScheduleItem(item: ScheduleItem) {
  if (!trip.value) return
  clearStatus()

  const title = editItemForm.title.trim()
  if (!title) {
    setStatus('error', 'Title cannot be empty.')
    return
  }

  updatingItemId.value = item.id
  try {
    await tripService.updateScheduleItem(trip.value.id, item.id, {
      title,
      description: editItemForm.description.trim(),
      start_time: editItemForm.start_time ? toISOStringFromLocalInput(editItemForm.start_time) : null,
      end_time: editItemForm.end_time ? toISOStringFromLocalInput(editItemForm.end_time) : null,
      day_index: selectedDay.value,
      order_index: item.order_index,
    })

    cancelEditScheduleItem()
    await refreshTrip()
    setStatus('success', 'Schedule item updated.')
  } catch (error) {
    const fetchError = error as FetchError<{ message?: string }>
    setStatus('error', fetchError?.data?.message || fetchError?.message || 'Cannot update schedule item.')
  } finally {
    updatingItemId.value = null
  }
}

async function deleteScheduleItem(item: ScheduleItem) {
  if (!trip.value) return
  clearStatus()

  deletingItemId.value = item.id
  try {
    await tripService.deleteScheduleItem(trip.value.id, item.id)
    await refreshTrip()
    setStatus('success', 'Schedule item deleted.')
  } catch (error) {
    const fetchError = error as FetchError<{ message?: string }>
    setStatus('error', fetchError?.data?.message || fetchError?.message || 'Cannot delete schedule item.')
  } finally {
    deletingItemId.value = null
  }
}

async function moveScheduleItem(item: ScheduleItem, direction: -1 | 1) {
  if (!trip.value || isReordering.value) return

  const rows = [...selectedDaySchedule.value]
  const currentIndex = rows.findIndex(row => row.id === item.id)
  const targetIndex = currentIndex + direction
  if (currentIndex < 0 || targetIndex < 0 || targetIndex >= rows.length) {
    return
  }

  clearStatus()
  isReordering.value = true

  const swapped = [...rows]
  ;[swapped[currentIndex], swapped[targetIndex]] = [swapped[targetIndex], swapped[currentIndex]]

  const payload = {
    items: swapped.map((row, index) => ({
      id: row.id,
      day_index: selectedDay.value,
      order_index: index + 1,
    })),
  }

  try {
    await tripService.reorderScheduleItems(trip.value.id, payload)
    await refreshTrip()
  } catch (error) {
    const fetchError = error as FetchError<{ message?: string }>
    setStatus('error', fetchError?.data?.message || fetchError?.message || 'Cannot reorder schedule items.')
  } finally {
    isReordering.value = false
  }
}
</script>

<template>
  <section class="relative isolate min-h-screen overflow-hidden px-4 pb-16 pt-24 sm:px-6 lg:px-8">
    <div
      class="pointer-events-none absolute inset-0 -z-10 bg-[radial-gradient(circle_at_18%_16%,#f5a76e40,transparent_34%),radial-gradient(circle_at_84%_12%,#6fd4bf40,transparent_35%),linear-gradient(180deg,#f9f5e8,#f7edd8)]"
    />

    <div class="mx-auto w-full max-w-[1360px] space-y-4">
      <NuxtLink
        to="/trips"
        class="inline-flex items-center rounded-full border border-pt-lavender bg-white/85 px-4 py-2 text-sm font-semibold text-pt-blue transition hover:border-pt-blue"
      >
        Back to trips
      </NuxtLink>

      <div v-if="!isLoggedIn" class="rounded-2xl border border-amber-200 bg-amber-50 p-4 text-amber-800">
        You need to login to view and edit this trip plan.
      </div>

      <div v-else-if="pending" class="rounded-2xl border border-pt-lavender bg-white p-5 text-pt-ink-light">
        Loading trip data...
      </div>

      <div v-else-if="error" class="rounded-2xl border border-red-200 bg-red-50 p-5 text-red-700">
        Cannot load trip detail.
        <button class="ml-2 font-semibold underline" @click="refresh()">Retry</button>
      </div>

      <template v-else-if="trip">
        <div class="rounded-3xl border border-pt-lavender/70 bg-white/95 p-4 shadow-sm sm:p-5">
          <div class="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
            <div class="flex items-start gap-4">
              <img
                :src="trip.cover_image_url || 'https://images.unsplash.com/photo-1527631746610-bca00a040d60?w=800&q=80&auto=format&fit=crop'"
                :alt="trip.title"
                class="h-20 w-24 rounded-xl border border-pt-lavender object-cover sm:h-24 sm:w-32"
              />
              <div class="space-y-2">
                <h1 class="font-display text-2xl font-bold text-pt-blue sm:text-4xl">{{ trip.title }}</h1>
                <p class="text-sm text-pt-ink-light">
                  <CalendarDays class="mr-1 inline h-4 w-4" />
                  {{ formatDate(trip.start_date) }} - {{ formatDate(trip.end_date) }}
                  <span class="ml-1">({{ totalDays }} days)</span>
                </p>
                <div class="flex flex-wrap gap-4 text-sm text-pt-ink">
                  <span>
                    <MapPin class="mr-1 inline h-4 w-4" />
                    {{ mapPoints.length }} places
                  </span>
                  <span>
                    <Users class="mr-1 inline h-4 w-4" />
                    {{ trip.members.length }} members
                  </span>
                  <span class="rounded-full border px-2 py-0.5 text-xs font-semibold">
                    {{ capitalize(trip.visibility) }}
                  </span>
                  <span class="rounded-full border px-2 py-0.5 text-xs font-semibold">
                    {{ capitalize(trip.status) }}
                  </span>
                </div>
              </div>
            </div>

            <div class="flex flex-wrap gap-2">
              <button
                type="button"
                :disabled="isSavingDraft || isPublishing"
                class="inline-flex items-center gap-2 rounded-xl border border-pt-blue px-4 py-2 text-sm font-semibold text-pt-blue transition hover:bg-pt-blue hover:text-white disabled:cursor-not-allowed disabled:opacity-70"
                @click="updateTripStatus('draft')"
              >
                <Save class="h-4 w-4" />
                <span v-if="isSavingDraft" class="inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                Save draft
              </button>
              <button
                type="button"
                :disabled="isSavingDraft || isPublishing"
                class="inline-flex items-center gap-2 rounded-xl bg-pt-blue px-4 py-2 text-sm font-semibold text-white transition hover:bg-pt-blue-dark disabled:cursor-not-allowed disabled:opacity-70"
                @click="updateTripStatus('published')"
              >
                <Send class="h-4 w-4" />
                <span v-if="isPublishing" class="inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                Publish trip
              </button>
            </div>
          </div>

          <p
            v-if="statusMessage"
            class="mt-4 rounded-xl border px-4 py-3 text-sm font-semibold"
            :class="statusType === 'success' ? 'border-emerald-300 bg-emerald-50 text-emerald-700' : 'border-red-300 bg-red-50 text-red-700'"
          >
            {{ statusMessage }}
          </p>
        </div>

        <div class="grid gap-4 xl:grid-cols-[1.2fr_1.5fr_0.8fr]">
          <section class="rounded-2xl border border-pt-lavender bg-white/95 p-4 shadow-sm">
            <h2 class="text-2xl font-extrabold text-pt-blue">Detailed itinerary ({{ totalDays }} days)</h2>

            <div class="mt-3 flex flex-wrap gap-2">
              <button
                v-for="day in dayOptions"
                :key="day"
                type="button"
                class="rounded-xl border px-3 py-2 text-sm font-semibold transition"
                :class="day === selectedDay ? 'border-pt-blue bg-[#E8EEF9] text-pt-blue' : 'border-slate-200 bg-white text-slate-700 hover:border-pt-blue'"
                @click="selectedDay = day"
              >
                Day {{ day }}
              </button>
            </div>

            <div class="mt-4 rounded-xl border border-slate-200 bg-slate-50 p-3">
              <p class="mb-2 text-xs font-bold uppercase tracking-wide text-slate-500">Add schedule item</p>
              <div class="space-y-2">
                <input
                  v-model="createItemForm.title"
                  type="text"
                  placeholder="Title"
                  class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                >
                <textarea
                  v-model="createItemForm.description"
                  rows="2"
                  placeholder="Description"
                  class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                />
                <div class="grid grid-cols-2 gap-2">
                  <input
                    v-model="createItemForm.start_time"
                    type="datetime-local"
                    class="rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                  >
                  <input
                    v-model="createItemForm.end_time"
                    type="datetime-local"
                    class="rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                  >
                </div>
                <input
                  v-model="createItemForm.place_name"
                  type="text"
                  placeholder="Place name (required if map point selected)"
                  class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                >
                <input
                  v-model="createItemForm.place_address"
                  type="text"
                  placeholder="Place address (optional)"
                  class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                >
                <p v-if="pickedLocation" class="text-xs text-slate-600">
                  Picked location: {{ pickedLocation.lat.toFixed(6) }}, {{ pickedLocation.lng.toFixed(6) }}
                </p>
                <p v-else class="text-xs text-slate-500">
                  Click on map to pick location and save lat/lng.
                </p>
                <button
                  type="button"
                  :disabled="isCreatingItem"
                  class="inline-flex w-full items-center justify-center gap-2 rounded-lg bg-pt-blue px-3 py-2 text-sm font-semibold text-white transition hover:bg-pt-blue-dark disabled:cursor-not-allowed disabled:opacity-70"
                  @click="createScheduleItem"
                >
                  <Plus class="h-4 w-4" />
                  <span v-if="isCreatingItem" class="inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent" />
                  Add to Day {{ selectedDay }}
                </button>
              </div>
            </div>

            <div class="mt-4 space-y-3">
              <article
                v-for="(item, index) in selectedDaySchedule"
                :key="item.id"
                class="rounded-xl border border-slate-200 bg-white p-3"
              >
                <template v-if="editingItemId === item.id">
                  <div class="space-y-2">
                    <input
                      v-model="editItemForm.title"
                      type="text"
                      class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                    >
                    <textarea
                      v-model="editItemForm.description"
                      rows="2"
                      class="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                    />
                    <div class="grid grid-cols-2 gap-2">
                      <input
                        v-model="editItemForm.start_time"
                        type="datetime-local"
                        class="rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                      >
                      <input
                        v-model="editItemForm.end_time"
                        type="datetime-local"
                        class="rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                      >
                    </div>
                    <div class="flex gap-2">
                      <button
                        type="button"
                        :disabled="updatingItemId === item.id"
                        class="rounded-lg bg-pt-blue px-3 py-2 text-xs font-semibold text-white"
                        @click="saveScheduleItem(item)"
                      >
                        Save
                      </button>
                      <button
                        type="button"
                        class="rounded-lg border border-slate-300 px-3 py-2 text-xs font-semibold text-slate-700"
                        @click="cancelEditScheduleItem"
                      >
                        Cancel
                      </button>
                    </div>
                  </div>
                </template>

                <template v-else>
                  <div class="flex items-start justify-between gap-2">
                    <div>
                      <p class="text-base font-bold text-slate-900">{{ index + 1 }}. {{ item.title }}</p>
                      <p v-if="item.start_time || item.end_time" class="mt-1 text-sm text-slate-700">
                        <Clock3 class="mr-1 inline h-4 w-4" />
                        {{ formatTime(item.start_time) || '--:--' }} - {{ formatTime(item.end_time) || '--:--' }}
                      </p>
                      <p v-if="item.trip_place?.place.address" class="mt-1 text-sm text-slate-600">
                        <MapPin class="mr-1 inline h-4 w-4" />
                        {{ item.trip_place.place.address }}
                      </p>
                      <p v-if="item.description" class="mt-2 inline-block rounded-lg bg-slate-100 px-2 py-1 text-xs text-slate-700">
                        {{ item.description }}
                      </p>
                    </div>
                    <div class="flex gap-1">
                      <button type="button" class="rounded-md border p-1 text-slate-600" @click="moveScheduleItem(item, -1)">
                        <ArrowUp class="h-4 w-4" />
                      </button>
                      <button type="button" class="rounded-md border p-1 text-slate-600" @click="moveScheduleItem(item, 1)">
                        <ArrowDown class="h-4 w-4" />
                      </button>
                      <button type="button" class="rounded-md border p-1 text-slate-600" @click="startEditScheduleItem(item)">
                        <Pencil class="h-4 w-4" />
                      </button>
                      <button
                        type="button"
                        :disabled="deletingItemId === item.id"
                        class="rounded-md border p-1 text-red-600"
                        @click="deleteScheduleItem(item)"
                      >
                        <Trash2 class="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                </template>
              </article>

              <p v-if="selectedDaySchedule.length === 0" class="rounded-xl border border-dashed border-slate-300 p-4 text-sm text-slate-500">
                No schedule item in this day yet.
              </p>
            </div>
          </section>

          <section class="rounded-2xl border border-pt-lavender bg-white/95 p-4 shadow-sm">
            <h2 class="text-2xl font-extrabold text-pt-blue">Trip map and places</h2>
            <div class="mt-3">
              <div class="flex gap-2">
                <input
                  v-model="placeSearchQuery"
                  type="text"
                  placeholder="Search place (e.g. Hon Chong Nha Trang)"
                  class="w-full rounded-xl border border-slate-300 px-3 py-2 text-sm outline-none focus:border-pt-blue"
                  @keyup.enter="searchPlaces"
                >
                <button
                  type="button"
                  :disabled="placeSearchLoading"
                  class="rounded-xl bg-pt-blue px-4 py-2 text-sm font-semibold text-white transition hover:bg-pt-blue-dark disabled:cursor-not-allowed disabled:opacity-70"
                  @click="searchPlaces"
                >
                  <span
                    v-if="placeSearchLoading"
                    class="inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"
                  />
                  Search
                </button>
              </div>
              <div v-if="placeSearchResults.length > 0" class="mt-2 max-h-52 overflow-auto rounded-xl border border-slate-200 bg-white">
                <button
                  v-for="(result, index) in placeSearchResults"
                  :key="`${result.lat}-${result.lng}-${index}`"
                  type="button"
                  class="block w-full border-b border-slate-100 px-3 py-2 text-left text-sm hover:bg-slate-50"
                  @click="chooseSearchPlace(result)"
                >
                  <p class="font-semibold text-slate-800">{{ result.name }}</p>
                  <p class="text-xs text-slate-500">{{ result.address }}</p>
                </button>
              </div>
            </div>
            <div class="mt-3 overflow-hidden rounded-xl border border-slate-200 bg-slate-50">
              <ClientOnly>
                <TripLeafletMap
                  :points="mapPoints"
                  :selected-day="selectedDay"
                  :selected-point="pickedLocation"
                  @point-selected="onMapPointSelected"
                />
                <template #fallback>
                  <div class="grid h-[380px] place-items-center text-sm text-slate-500">Loading map...</div>
                </template>
              </ClientOnly>
            </div>
            <p class="mt-3 text-xs text-slate-500">
              Map uses Leaflet with OpenStreetMap and renders markers from schedule items that already have coordinates.
            </p>
          </section>

          <div class="space-y-4">
            <section class="rounded-2xl border border-pt-lavender bg-white/95 p-4 shadow-sm">
              <h2 class="text-2xl font-extrabold text-pt-blue">Members ({{ trip.members.length }})</h2>
              <div class="mt-3 space-y-3">
                <article v-for="member in trip.members" :key="member.user_id" class="flex items-center gap-3">
                  <img
                    v-if="member.avatar_url"
                    :src="member.avatar_url"
                    :alt="member.full_name"
                    class="h-11 w-11 rounded-full border border-slate-200 object-cover"
                  >
                  <div
                    v-else
                    class="grid h-11 w-11 place-content-center rounded-full bg-slate-200 text-sm font-bold text-slate-700"
                  >
                    {{ member.full_name?.charAt(0)?.toUpperCase() || '?' }}
                  </div>
                  <div>
                    <p class="font-semibold text-slate-900">{{ member.full_name }}</p>
                    <p class="text-sm text-slate-600">{{ roleLabel(member.role) }}</p>
                  </div>
                </article>
              </div>
              <p class="mt-4 text-xs text-slate-500">Invite member API is not available on BE yet.</p>
            </section>

            <section class="rounded-2xl border border-pt-lavender bg-white/95 p-4 shadow-sm">
              <h2 class="text-2xl font-extrabold text-pt-blue">Photo album</h2>
              <div class="mt-3 grid grid-cols-3 gap-2">
                <img
                  v-for="photo in trip.album_preview"
                  :key="photo.id"
                  :src="photo.thumbnail_url || photo.url"
                  :alt="photo.caption || 'Trip photo'"
                  class="h-20 w-full rounded-lg border border-slate-200 object-cover"
                >
              </div>
              <p v-if="trip.album_preview.length === 0" class="mt-3 text-sm text-slate-500">
                No photos in album yet.
              </p>
            </section>
          </div>
        </div>
      </template>

      <p v-else class="rounded-2xl border border-slate-200 bg-white p-5 text-slate-700">
        Trip not found.
      </p>
    </div>
  </section>
</template>
