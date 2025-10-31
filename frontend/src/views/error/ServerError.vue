<template>
  <ErrorState
    code="500"
    title="Terjadi kesalahan pada sistem"
    description="Kami mengalami kendala tak terduga saat memproses permintaan Anda. Tim kami sudah menerima notifikasi dan sedang menelusurinya."
    badge="Internal Error"
    :primary-action="primaryAction"
    :secondary-action="secondaryAction"
    show-back-button
  />
</template>

<script setup>
import ErrorState from '@/components/ErrorState.vue'

const primaryAction = {
  label: 'Kembali ke Dashboard',
  to: '/',
  icon: 'pi pi-home'
}

const secondaryAction = {
  label: 'Coba Lagi',
  event: 'reload',
  icon: 'pi pi-refresh'
}
</script>
