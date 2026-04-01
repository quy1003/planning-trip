<template>
  <div class="app-shell">
    <AppHeader />
    <main class="app-main">
      <slot />
    </main>
  </div>
</template>

<style scoped>
.app-shell {
  min-height: 100vh;
  background: radial-gradient(circle at top right, #eef6ff 0%, #f7fbff 35%, #ffffff 100%);
}

.app-main {
  max-width: 1120px;
  margin: 0 auto;
  padding: 24px 16px 40px;
}
</style>
