<template>
  <ErrorState
    code="403"
    title="Akses ditolak"
    description="Anda tidak memiliki izin untuk mengakses halaman ini. Jika Anda yakin ini adalah sebuah kesalahan, hubungi administrator untuk mendapatkan akses."
    badge="Restricted Area"
    :primary-action="primaryAction"
    :secondary-action="secondaryAction"
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
  label: 'Hubungi Dukungan',
  to: '/support',
  icon: 'pi pi-comments'
}
</script>
