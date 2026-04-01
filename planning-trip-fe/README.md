# Planning Trip Frontend (Nuxt 3/4)

Frontend cho dự án Planning Trip, build bằng NuxtJS.

## Run

```bash
npm install
npm run dev
```

App chạy mặc định tại `http://localhost:3000`.

## Env

File `.env.example`:

```bash
NUXT_PUBLIC_API_BASE=http://localhost:8080
```

Runtime config sẽ dùng biến này để gọi backend Go.

## Source structure

```text
app/
  app.vue
  layouts/
    default.vue
  pages/
    index.vue
    trips/
      index.vue
      [id].vue
  components/
    common/
      AppHeader.vue
    trip/
      TripCard.vue
  services/
    http.ts
    trip.service.ts
  composables/
    useTrips.ts
  types/
    trip.ts
  assets/
    styles/
      main.css
```

## Notes

- `services/http.ts`: wrapper `$fetch` + base URL từ runtime config.
- `services/trip.service.ts`: service layer để gọi API trips.
- `composables/useTrips.ts`: logic dùng lại cho pages.
- `pages/trips/*`: module trips (list/detail) để mở rộng dần.
