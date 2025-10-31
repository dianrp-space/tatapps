<template>
  <ErrorState
    code="404"
    title="Halaman tidak ditemukan"
    description="Kami tidak dapat menemukan halaman yang Anda cari. Periksa kembali alamat URL atau gunakan tombol di bawah untuk kembali ke alur kerja Anda."
    badge="Not Found"
    :primary-action="primaryAction"
    :secondary-action="secondaryAction"
  />
</template>

<script setup>
import { computed } from 'vue'
import ErrorState from '@/components/ErrorState.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const primaryAction = computed(() => {
  const isAuthenticated = authStore.isAuthenticated
  return {
    label: isAuthenticated ? 'Kembali ke Dashboard' : 'Masuk ke Aplikasi',
    to: isAuthenticated ? '/' : '/login',
    icon: isAuthenticated ? 'pi pi-chart-bar' : 'pi pi-sign-in'
  }
})

const secondaryAction = {
  label: 'Halaman Sebelumnya',
  event: 'back',
  icon: 'pi pi-arrow-left'
}
</script>
