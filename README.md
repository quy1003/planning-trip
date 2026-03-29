# planning-trip

A collaborative trip planning workspace with interactive maps, daily schedules, and shared photo galleries.

## Features

- **Trips** – Create and manage multiple trips (name, destination, dates).
- **Daily Schedule** – Add activities (time, title, notes) to each day of a trip.
- **Photo Gallery** – Attach photos (via URL) with captions to any trip.
- **Interactive Map** – View your destination on an OpenStreetMap-powered Leaflet map.

## Getting Started

No build step required. Simply open `index.html` in a browser:

```bash
# Using Python's built-in server (recommended to avoid CORS issues)
python3 -m http.server 8080
# then open http://localhost:8080
```

Or just double-click `index.html` to open it directly in any modern browser.

## Technology

| Concern | Solution |
|---------|---------|
| UI | Vanilla HTML, CSS, JavaScript (ES6+) |
| Maps | [Leaflet.js](https://leafletjs.com/) + [OpenStreetMap](https://www.openstreetmap.org/) |
| Geocoding | [Nominatim API](https://nominatim.openstreetmap.org/) |
| Data persistence | `localStorage` (no backend needed) |

## Project Structure

```
planning-trip/
├── index.html       # App shell
├── css/
│   └── style.css    # All styles (CSS custom properties, responsive)
└── js/
    └── app.js       # All application logic (routing, rendering, data)
```
