<script setup lang="ts">
interface MapPoint {
  id: string
  title: string
  lat: number
  lng: number
  day_index: number
}

const props = defineProps<{
  points: MapPoint[]
  selectedDay: number
  selectedPoint?: { lat: number; lng: number } | null
}>()

const emit = defineEmits<{
  (e: 'point-selected', payload: { lat: number; lng: number }): void
}>()

const containerRef = ref<HTMLElement | null>(null)
let map: any = null
let layerGroup: any = null
let lineLayer: any = null
let pickedMarker: any = null

function addLeafletCss() {
  if (document.getElementById('leaflet-css')) return
  const link = document.createElement('link')
  link.id = 'leaflet-css'
  link.rel = 'stylesheet'
  link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css'
  document.head.appendChild(link)
}

function loadLeafletScript(): Promise<any> {
  const existing = (window as any).L
  if (existing) return Promise.resolve(existing)

  return new Promise((resolve, reject) => {
    const current = document.getElementById('leaflet-js') as HTMLScriptElement | null
    if (current) {
      current.addEventListener('load', () => resolve((window as any).L))
      current.addEventListener('error', () => reject(new Error('Failed to load leaflet')))
      return
    }

    const script = document.createElement('script')
    script.id = 'leaflet-js'
    script.src = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.js'
    script.async = true
    script.onload = () => resolve((window as any).L)
    script.onerror = () => reject(new Error('Failed to load leaflet'))
    document.body.appendChild(script)
  })
}

function renderMap() {
  if (!map || !(window as any).L) return

  const L = (window as any).L
  const points = props.points
    .filter(point => point.day_index === props.selectedDay)
    .filter(point => Number.isFinite(point.lat) && Number.isFinite(point.lng))

  layerGroup.clearLayers()
  if (lineLayer) {
    map.removeLayer(lineLayer)
    lineLayer = null
  }

  if (points.length === 0) return

  const latlngs: [number, number][] = []
  points.forEach((point, index) => {
    const latlng: [number, number] = [point.lat, point.lng]
    latlngs.push(latlng)
    const marker = L.marker(latlng).bindPopup(`${index + 1}. ${point.title}`)
    layerGroup.addLayer(marker)
  })

  lineLayer = L.polyline(latlngs, {
    color: '#2a3570',
    weight: 4,
    opacity: 0.8,
  }).addTo(map)

  map.fitBounds(L.latLngBounds(latlngs), { padding: [30, 30] })
}

function onMapClick(event: any) {
  if (!map || !(window as any).L) return

  const L = (window as any).L
  const { lat, lng } = event.latlng

  if (pickedMarker) {
    map.removeLayer(pickedMarker)
  }

  pickedMarker = L.marker([lat, lng], {
    draggable: true,
  }).addTo(map)

  pickedMarker.bindPopup('Selected point').openPopup()
  pickedMarker.on('dragend', (dragEvent: any) => {
    const newPos = dragEvent.target.getLatLng()
    emit('point-selected', { lat: newPos.lat, lng: newPos.lng })
  })

  emit('point-selected', { lat, lng })
}

function syncPickedPoint(point?: { lat: number; lng: number } | null) {
  if (!map || !point || !(window as any).L) return
  const L = (window as any).L

  if (pickedMarker) {
    map.removeLayer(pickedMarker)
  }

  pickedMarker = L.marker([point.lat, point.lng], {
    draggable: true,
  }).addTo(map)

  pickedMarker.bindPopup('Selected point').openPopup()
  pickedMarker.on('dragend', (dragEvent: any) => {
    const newPos = dragEvent.target.getLatLng()
    emit('point-selected', { lat: newPos.lat, lng: newPos.lng })
  })

  map.setView([point.lat, point.lng], Math.max(map.getZoom(), 14))
}

onMounted(async () => {
  if (!containerRef.value) return
  addLeafletCss()

  try {
    const L = await loadLeafletScript()
    map = L.map(containerRef.value, {
      zoomControl: true,
    }).setView([16.047079, 108.20623], 11)

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      maxZoom: 19,
      attribution: '&copy; OpenStreetMap contributors',
    }).addTo(map)

    map.on('click', onMapClick)
    layerGroup = L.layerGroup().addTo(map)
    renderMap()
  } catch (error) {
    console.error(error)
  }
})

watch(
  () => [props.points, props.selectedDay],
  () => {
    renderMap()
  },
  { deep: true },
)

watch(
  () => props.selectedPoint,
  point => {
    syncPickedPoint(point)
  },
  { deep: true },
)

onBeforeUnmount(() => {
  if (map) {
    map.remove()
    map = null
  }
})
</script>

<template>
  <div ref="containerRef" class="h-[380px] w-full bg-slate-100" />
</template>
