<template>
  <div id="app">
    <Toast />
    <router-view />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useSiteStore } from '@/stores/site'

const authStore = useAuthStore()
const siteStore = useSiteStore()

onMounted(() => {
  siteStore.fetchSiteSettings().catch(() => {
    // Ignore failures; branding will fall back to defaults
  })

  // Check if user is logged in
  authStore.checkAuth()
})
</script>
