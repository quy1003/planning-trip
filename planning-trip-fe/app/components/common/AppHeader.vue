<template>
  <header class="header">
    <div class="header-inner">
      <NuxtLink to="/" class="brand">Planning Trip</NuxtLink>
      <nav class="menu">
        <NuxtLink to="/trips">Trips</NuxtLink>
      </nav>
    </div>
  </header>
</template>

<style scoped>
.header {
  position: sticky;
  top: 0;
  z-index: 20;
  background: rgba(255, 255, 255, 0.9);
  border-bottom: 1px solid #e6edf5;
  backdrop-filter: blur(8px);
}

.header-inner {
  max-width: 1120px;
  margin: 0 auto;
  padding: 14px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.brand {
  font-size: 18px;
  font-weight: 700;
  color: #0f172a;
  text-decoration: none;
}

.menu {
  display: flex;
  gap: 16px;
}

.menu a {
  color: #334155;
  text-decoration: none;
  font-weight: 500;
}
</style>
