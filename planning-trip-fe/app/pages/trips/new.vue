<script setup lang="ts">
import { FetchError } from 'ofetch'
import { ArrowLeft, PlusCircle } from 'lucide-vue-next'
import { tripService } from '~/services/trip.service'
import { uploadService } from '~/services/upload.service'

const router = useRouter()

const form = reactive({
  title: '',
  description: '',
  start_date: '',
  end_date: '',
  visibility: 'private',
  status: 'draft',
  cover_image_url: '',
})

const isSubmitting = ref(false)
const isUploading = ref(false)
const status = ref<'success' | 'error' | ''>('')
const statusMessage = ref('')
const selectedFile = ref<File | null>(null)
const coverFileInputRef = ref<HTMLInputElement | null>(null)

const isLoggedIn = computed(() => {
  if (!import.meta.client) {
    return false
  }
  return Boolean(localStorage.getItem('access_token'))
})

async function onSubmit() {
  if (isSubmitting.value || isUploading.value) {
    return
  }

  status.value = ''
  statusMessage.value = ''

  if (!isLoggedIn.value) {
    status.value = 'error'
    statusMessage.value = 'You need to login before creating a trip.'
    return
  }

  if (!form.title.trim()) {
    status.value = 'error'
    statusMessage.value = 'Please enter trip title.'
    return
  }

  if (form.start_date && form.end_date && form.start_date > form.end_date) {
    status.value = 'error'
    statusMessage.value = 'Start date must be before or equal to end date.'
    return
  }

  isSubmitting.value = true

  try {
    const createdTrip = await tripService.create({
      title: form.title.trim(),
      description: form.description.trim(),
      start_date: form.start_date || undefined,
      end_date: form.end_date || undefined,
      visibility: form.visibility,
      status: form.status,
      cover_image_url: form.cover_image_url.trim() || undefined,
    })

    status.value = 'success'
    statusMessage.value = 'Trip created successfully.'
    await router.push(`/trips/${createdTrip.id}`)
  } catch (error) {
    const fetchError = error as FetchError<{ message?: string }>
    status.value = 'error'
    statusMessage.value =
      fetchError?.data?.message || fetchError?.message || 'Cannot create trip.'
  } finally {
    isSubmitting.value = false
  }
}

function onCoverFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  selectedFile.value = target.files?.[0] || null
  void uploadCoverImage()
}

async function uploadCoverImage() {
  if (isUploading.value) {
    return
  }

  status.value = ''
  statusMessage.value = ''

  if (!isLoggedIn.value) {
    status.value = 'error'
    statusMessage.value = 'You need to login before uploading image.'
    return
  }

  if (!selectedFile.value) {
    status.value = 'error'
    statusMessage.value = 'Please choose an image first.'
    return
  }

  isUploading.value = true
  try {
    const uploaded = await uploadService.uploadImage(selectedFile.value)
    form.cover_image_url = uploaded.url
  } catch (error) {
    const fetchError = error as FetchError<{ message?: string }>
    status.value = 'error'
    statusMessage.value =
      fetchError?.data?.message || fetchError?.message || 'Cannot upload image.'
  } finally {
    isUploading.value = false
    if (coverFileInputRef.value) {
      coverFileInputRef.value.value = ''
    }
  }
}

function selectAndUploadCoverImage() {
  if (isUploading.value) {
    return
  }
  coverFileInputRef.value?.click()
}
</script>

<template>
  <section
    class="relative isolate min-h-screen overflow-hidden px-4 pb-16 pt-24 sm:px-6 lg:px-8"
  >
    <div
      class="pointer-events-none absolute inset-0 -z-10 bg-[radial-gradient(circle_at_20%_20%,#f5a76e55,transparent_35%),radial-gradient(circle_at_80%_15%,#3db8a633,transparent_32%),linear-gradient(180deg,#efe2ba66,#f8f2df)]"
    />

    <div class="mx-auto w-full max-w-3xl">
      <NuxtLink
        to="/trips"
        class="mb-6 inline-flex items-center gap-2 rounded-full border border-pt-orange/50 bg-white/80 px-4 py-2 text-sm font-semibold text-pt-orange backdrop-blur"
      >
        <ArrowLeft class="h-4 w-4" />
        Back to trips
      </NuxtLink>

      <div
        class="rounded-3xl border border-pt-lavender/60 bg-white/95 p-6 shadow-[0_20px_60px_-24px_rgba(31,35,64,0.35)] sm:p-8"
      >
        <h1 class="mt-2 font-display text-4xl font-bold text-pt-orange">
          Create a new trip
        </h1>

        <div
          v-if="statusMessage"
          class="mt-5 rounded-xl border px-4 py-3 text-sm font-medium"
          :class="
            status === 'success'
              ? 'border-emerald-300 bg-emerald-50 text-emerald-700'
              : 'border-red-300 bg-red-50 text-red-700'
          "
        >
          {{ statusMessage }}
        </div>

        <form class="mt-6 grid gap-4" @submit.prevent="onSubmit">
          <div>
            <label class="mb-1 block text-sm font-semibold text-pt-orange"
              >Trip title *</label
            >
            <input
              v-model="form.title"
              type="text"
              placeholder="Example: Da Nang 3D2N"
              class="w-full rounded-xl border border-pt-lavender bg-white px-4 py-3 text-sm text-black outline-none transition focus:border-pt-orange"
            />
          </div>

          <div>
            <label class="mb-1 block text-sm font-semibold text-pt-orange"
              >Description *</label
            >
            <textarea
              v-model="form.description"
              rows="4"
              placeholder="Short description for your trip"
              class="w-full rounded-xl border border-pt-lavender bg-white px-4 py-3 text-sm text-black outline-none transition focus:border-pt-orange"
            />
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-1 block text-sm font-semibold text-pt-orange"
                >Start date *</label
              >
              <input
                v-model="form.start_date"
                type="date"
                class="w-full rounded-xl border border-pt-lavender bg-white px-4 py-3 text-sm text-black outline-none transition focus:border-pt-orange"
              />
            </div>
            <div>
              <label class="mb-1 block text-sm font-semibold text-pt-orange"
                >End date</label
              >
              <input
                v-model="form.end_date"
                type="date"
                class="w-full rounded-xl border border-pt-lavender bg-white px-4 py-3 text-sm text-black outline-none transition focus:border-pt-orange"
              />
            </div>
          </div>

          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-1 block text-sm font-semibold text-pt-orange"
                >Visibility *</label
              >
              <select
                v-model="form.visibility"
                class="w-full rounded-xl border border-pt-lavender bg-white px-4 py-3 text-sm text-black outline-none transition focus:border-pt-orange"
              >
                <option value="private">Private</option>
                <option value="group">Group</option>
                <option value="public">Public</option>
              </select>
            </div>

            <div>
              <label class="mb-1 block text-sm font-semibold text-pt-orange"
                >Status *</label
              >
              <select
                v-model="form.status"
                class="w-full rounded-xl border border-pt-lavender bg-white px-4 py-3 text-sm text-black outline-none transition focus:border-pt-orange"
              >
                <option value="draft">Draft</option>
                <option value="published">Published</option>
              </select>
            </div>
          </div>

          <div>
            <input
              ref="coverFileInputRef"
              type="file"
              accept="image/*"
              class="hidden"
              @change="onCoverFileChange"
            />
            <button
              type="button"
              :disabled="isUploading"
              class="mt-3 inline-flex items-center justify-center gap-2 rounded-xl border border-pt-orange px-4 py-3 text-sm font-semibold text-pt-orange transition hover:bg-pt-orange hover:text-white"
              :class="isUploading ? 'cursor-not-allowed opacity-70' : ''"
              @click="selectAndUploadCoverImage"
            >
              <span
                v-if="isUploading"
                class="inline-block h-4 w-4 animate-spin rounded-full border-2 border-current border-t-transparent"
              />
              {{ isUploading ? 'Uploading...' : 'Upload image' }}
            </button>

            <p
              v-if="form.cover_image_url"
              class="mt-2 text-xs text-pt-ink-light"
            >
              Image uploaded successfully and will be used as trip cover.
            </p>
            <img
              v-if="form.cover_image_url"
              :src="form.cover_image_url"
              alt="Cover preview"
              class="mt-3 h-40 w-full rounded-xl border border-pt-lavender object-cover"
            />
          </div>

          <button
            type="submit"
            :disabled="isSubmitting || isUploading"
            class="mt-2 inline-flex items-center justify-center gap-2 rounded-xl bg-pt-orange px-4 py-3 text-sm font-semibold text-white shadow-sm transition hover:bg-pt-orange-dark"
            :class="
              isSubmitting || isUploading ? 'cursor-not-allowed opacity-70' : ''
            "
          >
            <PlusCircle class="h-4 w-4" />
            {{ isSubmitting ? 'Creating...' : 'Create trip' }}
          </button>
        </form>
      </div>
    </div>
  </section>
</template>
