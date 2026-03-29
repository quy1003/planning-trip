/* ============================================================
   PlanTrip – app.js
   Vanilla JavaScript single-page app.
   Data is stored in localStorage; no build step required.
   ============================================================ */

// ──────────────────────────────────────────────
// Storage helpers
// ──────────────────────────────────────────────
const Storage = {
  load() {
    return JSON.parse(localStorage.getItem('plantrip_trips') || '[]');
  },
  save(trips) {
    localStorage.setItem('plantrip_trips', JSON.stringify(trips));
  },
};

// ──────────────────────────────────────────────
// Utility helpers
// ──────────────────────────────────────────────
function generateId() {
  return Date.now().toString(36) + Math.random().toString(36).slice(2);
}

function formatDate(dateStr) {
  if (!dateStr) return '';
  const [year, month, day] = dateStr.split('-');
  return `${day}/${month}/${year}`;
}

/** Returns an array of {dayNum, date} objects between start and end (inclusive). */
function buildDays(startDate, endDate) {
  const days = [];
  if (!startDate || !endDate) return days;

  const current = new Date(startDate);
  const last = new Date(endDate);
  let dayNum = 1;

  while (current <= last) {
    days.push({ dayNum, date: current.toISOString().split('T')[0] });
    current.setDate(current.getDate() + 1);
    dayNum++;
  }
  return days;
}

// ──────────────────────────────────────────────
// App state
// ──────────────────────────────────────────────
let currentTripId = null;
let currentTab = 'schedule';
let leafletMap = null;

// ──────────────────────────────────────────────
// Router – hash-based navigation
// ──────────────────────────────────────────────
function navigate(path) {
  window.location.hash = path;
}

function handleRoute() {
  const hash = window.location.hash;

  if (hash.startsWith('#/trip/')) {
    const id = hash.replace('#/trip/', '');
    currentTripId = id;
    renderTripDetail(id);
  } else {
    currentTripId = null;
    renderHome();
  }
}

// ──────────────────────────────────────────────
// Modal helpers
// ──────────────────────────────────────────────
function openModal(html) {
  document.getElementById('modal').innerHTML = html;
  document.getElementById('modal-overlay').classList.remove('hidden');
}

function closeModal() {
  document.getElementById('modal-overlay').classList.add('hidden');
  document.getElementById('modal').innerHTML = '';
}

// ──────────────────────────────────────────────
// Home page – list of trips
// ──────────────────────────────────────────────
function renderHome() {
  const trips = Storage.load();

  const cards = trips.length
    ? trips
        .map(
          (t) => `
      <div class="trip-card" data-id="${t.id}">
        <div class="trip-card__name">${escapeHtml(t.name)}</div>
        <div class="trip-card__destination">📍 ${escapeHtml(t.destination)}</div>
        <div class="trip-card__dates">
          ${formatDate(t.startDate)} – ${formatDate(t.endDate)}
        </div>
        <div class="trip-card__footer">
          <span class="badge">${buildDays(t.startDate, t.endDate).length} days</span>
          <button class="btn btn-ghost btn-sm delete-trip" data-id="${t.id}">🗑</button>
        </div>
      </div>`
        )
        .join('')
    : `<div class="empty-state">
        <div class="icon">🌍</div>
        <p>No trips yet. Click <strong>+ New Trip</strong> to get started!</p>
       </div>`;

  document.getElementById('app').innerHTML = `
    <div class="page">
      <div class="page-title">My Trips</div>
      <div class="trips-grid">${cards}</div>
    </div>`;

  // Navigate to trip on card click (not delete button)
  document.querySelectorAll('.trip-card').forEach((card) => {
    card.addEventListener('click', (e) => {
      if (!e.target.closest('.delete-trip')) {
        navigate(`#/trip/${card.dataset.id}`);
      }
    });
  });

  // Delete trip
  document.querySelectorAll('.delete-trip').forEach((btn) => {
    btn.addEventListener('click', (e) => {
      e.stopPropagation();
      if (confirm('Delete this trip?')) {
        deleteTrip(btn.dataset.id);
      }
    });
  });
}

// ──────────────────────────────────────────────
// Trip detail page
// ──────────────────────────────────────────────
function renderTripDetail(id) {
  const trips = Storage.load();
  const trip = trips.find((t) => t.id === id);

  if (!trip) {
    navigate('#/');
    return;
  }

  document.getElementById('app').innerHTML = `
    <div class="page">
      <button class="back-btn" id="btn-back">← Back to trips</button>

      <div class="detail-header">
        <div class="detail-header__top">
          <div>
            <div class="detail-title">${escapeHtml(trip.name)}</div>
            <div class="detail-meta">
              📍 ${escapeHtml(trip.destination)} &nbsp;·&nbsp;
              ${formatDate(trip.startDate)} – ${formatDate(trip.endDate)}
            </div>
          </div>
        </div>
      </div>

      <div class="tabs">
        <button class="tab-btn ${currentTab === 'schedule' ? 'active' : ''}" data-tab="schedule">📅 Schedule</button>
        <button class="tab-btn ${currentTab === 'gallery' ? 'active' : ''}" data-tab="gallery">🖼 Gallery</button>
        <button class="tab-btn ${currentTab === 'map' ? 'active' : ''}" data-tab="map">🗺️ Map</button>
      </div>

      <div id="tab-content"></div>
    </div>`;

  document.getElementById('btn-back').addEventListener('click', () => navigate('#/'));

  document.querySelectorAll('.tab-btn').forEach((btn) => {
    btn.addEventListener('click', () => {
      currentTab = btn.dataset.tab;
      document.querySelectorAll('.tab-btn').forEach((b) => b.classList.remove('active'));
      btn.classList.add('active');
      renderActiveTab(trip);
    });
  });

  renderActiveTab(trip);
}

function renderActiveTab(trip) {
  if (currentTab === 'schedule') renderScheduleTab(trip);
  else if (currentTab === 'gallery') renderGalleryTab(trip);
  else if (currentTab === 'map') renderMapTab(trip);
}

// ──────────────────────────────────────────────
// Schedule tab
// ──────────────────────────────────────────────
function renderScheduleTab(trip) {
  const days = buildDays(trip.startDate, trip.endDate);

  if (!days.length) {
    document.getElementById('tab-content').innerHTML =
      '<p style="color:var(--text-muted)">No dates configured for this trip.</p>';
    return;
  }

  const dayBlocks = days
    .map((d) => {
      const activities = (trip.schedule || []).filter((a) => a.date === d.date);
      const activityItems = activities.length
        ? activities
            .map(
              (a, i) => `
          <li class="activity-item">
            <div class="activity-time">${escapeHtml(a.time || '')}</div>
            <div class="activity-info">
              <div class="activity-title">${escapeHtml(a.title)}</div>
              ${a.notes ? `<div class="activity-notes">${escapeHtml(a.notes)}</div>` : ''}
            </div>
            <button class="btn btn-ghost btn-sm delete-activity"
              data-date="${d.date}" data-index="${i}">✕</button>
          </li>`
            )
            .join('')
        : `<li class="no-activities">No activities yet.</li>`;

      return `
        <div class="day-block">
          <div class="day-block__header">
            <span>Day ${d.dayNum} – ${formatDate(d.date)}</span>
            <button class="btn btn-primary btn-sm add-activity" data-date="${d.date}">
              + Add Activity
            </button>
          </div>
          <ul class="day-block__activities">${activityItems}</ul>
        </div>`;
    })
    .join('');

  document.getElementById('tab-content').innerHTML = dayBlocks;

  document.querySelectorAll('.add-activity').forEach((btn) => {
    btn.addEventListener('click', () => openAddActivityModal(trip.id, btn.dataset.date));
  });

  document.querySelectorAll('.delete-activity').forEach((btn) => {
    btn.addEventListener('click', () => {
      deleteActivity(trip.id, btn.dataset.date, parseInt(btn.dataset.index, 10));
    });
  });
}

// ──────────────────────────────────────────────
// Gallery tab
// ──────────────────────────────────────────────
function renderGalleryTab(trip) {
  const photos = trip.photos || [];

  const grid = photos.length
    ? photos
        .map(
          (p, i) => `
      <div class="photo-card">
        <img src="${escapeAttr(p.url)}" alt="${escapeAttr(p.caption)}"
             onerror="this.src='https://placehold.co/400x300?text=Image+not+found'" />
        <div class="photo-card__caption">
          <span>${escapeHtml(p.caption || '')}</span>
          <button class="btn btn-ghost btn-sm delete-photo" data-index="${i}">✕</button>
        </div>
      </div>`
        )
        .join('')
    : `<div class="empty-state"><div class="icon">📷</div><p>No photos yet.</p></div>`;

  document.getElementById('tab-content').innerHTML = `
    <div style="margin-bottom:16px;">
      <button class="btn btn-primary" id="btn-add-photo">+ Add Photo</button>
    </div>
    <div class="gallery-grid">${grid}</div>`;

  document.getElementById('btn-add-photo').addEventListener('click', () =>
    openAddPhotoModal(trip.id)
  );

  document.querySelectorAll('.delete-photo').forEach((btn) => {
    btn.addEventListener('click', () => deletePhoto(trip.id, parseInt(btn.dataset.index, 10)));
  });
}

// ──────────────────────────────────────────────
// Map tab – Leaflet + Nominatim geocoding
// ──────────────────────────────────────────────
function renderMapTab(trip) {
  document.getElementById('tab-content').innerHTML = `
    <div id="map"></div>
    <p style="margin-top:8px;font-size:0.85rem;color:var(--text-muted)">
      Showing map for: <strong>${escapeHtml(trip.destination)}</strong>
    </p>`;

  // Destroy previous map instance to avoid Leaflet re-init error
  if (leafletMap) {
    leafletMap.remove();
    leafletMap = null;
  }

  leafletMap = L.map('map').setView([20, 0], 2);

  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '© <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
    maxZoom: 18,
  }).addTo(leafletMap);

  // Geocode destination with Nominatim (User-Agent required by usage policy)
  fetch(
    `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(trip.destination)}&limit=1`,
    { headers: { 'User-Agent': 'PlanTrip/1.0 (https://github.com/quy1003/planning-trip)' } }
  )
    .then((res) => res.json())
    .then((data) => {
      if (data.length > 0) {
        const { lat, lon, display_name } = data[0];
        leafletMap.setView([lat, lon], 10);
        L.marker([lat, lon])
          .addTo(leafletMap)
          .bindPopup(`<strong>${escapeHtml(trip.destination)}</strong><br>${escapeHtml(display_name)}`)
          .openPopup();
      }
    })
    .catch(() => {
      // Map still visible without marker on geocoding failure
    });
}

// ──────────────────────────────────────────────
// Modals – New trip / Add activity / Add photo
// ──────────────────────────────────────────────
function openNewTripModal() {
  openModal(`
    <h2>New Trip</h2>
    <form id="form-new-trip">
      <div class="form-group">
        <label>Trip Name</label>
        <input name="name" placeholder="e.g. Summer in Japan" required />
      </div>
      <div class="form-group">
        <label>Destination</label>
        <input name="destination" placeholder="e.g. Tokyo, Japan" required />
      </div>
      <div class="form-row">
        <div class="form-group">
          <label>Start Date</label>
          <input type="date" name="startDate" required />
        </div>
        <div class="form-group">
          <label>End Date</label>
          <input type="date" name="endDate" required />
        </div>
      </div>
      <div class="form-actions">
        <button type="button" class="btn btn-ghost" id="btn-cancel">Cancel</button>
        <button type="submit" class="btn btn-primary">Create Trip</button>
      </div>
    </form>`);

  document.getElementById('btn-cancel').addEventListener('click', closeModal);
  document.getElementById('form-new-trip').addEventListener('submit', (e) => {
    e.preventDefault();
    const data = Object.fromEntries(new FormData(e.target));
    if (data.startDate > data.endDate) {
      alert('End date must be on or after start date.');
      return;
    }
    const id = createTrip(data);
    closeModal();
    navigate(`#/trip/${id}`);
  });
}

function openAddActivityModal(tripId, date) {
  openModal(`
    <h2>Add Activity – ${formatDate(date)}</h2>
    <form id="form-add-activity">
      <div class="form-group">
        <label>Time</label>
        <input type="time" name="time" />
      </div>
      <div class="form-group">
        <label>Activity</label>
        <input name="title" placeholder="e.g. Visit temple" required />
      </div>
      <div class="form-group">
        <label>Notes (optional)</label>
        <textarea name="notes" rows="3" placeholder="Extra details…"></textarea>
      </div>
      <div class="form-actions">
        <button type="button" class="btn btn-ghost" id="btn-cancel">Cancel</button>
        <button type="submit" class="btn btn-primary">Add</button>
      </div>
    </form>`);

  document.getElementById('btn-cancel').addEventListener('click', closeModal);
  document.getElementById('form-add-activity').addEventListener('submit', (e) => {
    e.preventDefault();
    const data = Object.fromEntries(new FormData(e.target));
    addActivity(tripId, date, data);
    closeModal();
    const trips = Storage.load();
    const trip = trips.find((t) => t.id === tripId);
    if (trip) renderScheduleTab(trip);
  });
}

function openAddPhotoModal(tripId) {
  openModal(`
    <h2>Add Photo</h2>
    <form id="form-add-photo">
      <div class="form-group">
        <label>Image URL</label>
        <input name="url" type="url" placeholder="https://example.com/photo.jpg" required />
      </div>
      <div class="form-group">
        <label>Caption (optional)</label>
        <input name="caption" placeholder="Describe the photo…" />
      </div>
      <div class="form-actions">
        <button type="button" class="btn btn-ghost" id="btn-cancel">Cancel</button>
        <button type="submit" class="btn btn-primary">Add Photo</button>
      </div>
    </form>`);

  document.getElementById('btn-cancel').addEventListener('click', closeModal);
  document.getElementById('form-add-photo').addEventListener('submit', (e) => {
    e.preventDefault();
    const data = Object.fromEntries(new FormData(e.target));
    addPhoto(tripId, data);
    closeModal();
    const trips = Storage.load();
    const trip = trips.find((t) => t.id === tripId);
    if (trip) renderGalleryTab(trip);
  });
}

// ──────────────────────────────────────────────
// Data mutations
// ──────────────────────────────────────────────
function createTrip({ name, destination, startDate, endDate }) {
  const trips = Storage.load();
  const id = generateId();
  trips.push({ id, name, destination, startDate, endDate, schedule: [], photos: [] });
  Storage.save(trips);
  return id;
}

function deleteTrip(id) {
  const trips = Storage.load().filter((t) => t.id !== id);
  Storage.save(trips);
  renderHome();
}

function addActivity(tripId, date, { time, title, notes }) {
  const trips = Storage.load();
  const trip = trips.find((t) => t.id === tripId);
  if (!trip) return;
  trip.schedule.push({ date, time, title, notes });
  trip.schedule.sort((a, b) => {
    const keyA = a.date + (a.time || '');
    const keyB = b.date + (b.time || '');
    return keyA.localeCompare(keyB);
  });
  Storage.save(trips);
}

function deleteActivity(tripId, date, index) {
  const trips = Storage.load();
  const trip = trips.find((t) => t.id === tripId);
  if (!trip) return;

  const dayActivities = trip.schedule.filter((a) => a.date === date);
  const activity = dayActivities[index];
  trip.schedule = trip.schedule.filter((a) => a !== activity);
  Storage.save(trips);
  renderScheduleTab(trip);
}

function addPhoto(tripId, { url, caption }) {
  const trips = Storage.load();
  const trip = trips.find((t) => t.id === tripId);
  if (!trip) return;
  trip.photos.push({ url, caption });
  Storage.save(trips);
}

function deletePhoto(tripId, index) {
  const trips = Storage.load();
  const trip = trips.find((t) => t.id === tripId);
  if (!trip) return;
  trip.photos.splice(index, 1);
  Storage.save(trips);
  renderGalleryTab(trip);
}

// ──────────────────────────────────────────────
// Security helpers – prevent XSS
// ──────────────────────────────────────────────
function escapeHtml(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

function escapeAttr(str) {
  return String(str)
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}

// ──────────────────────────────────────────────
// Bootstrap
// ──────────────────────────────────────────────
window.addEventListener('hashchange', handleRoute);
window.addEventListener('DOMContentLoaded', () => {
  document.getElementById('btn-new-trip').addEventListener('click', openNewTripModal);
  // Close modal when clicking backdrop
  document.getElementById('modal-overlay').addEventListener('click', (e) => {
    if (e.target === document.getElementById('modal-overlay')) closeModal();
  });
  handleRoute();
});
