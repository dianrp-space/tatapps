<template>
  <div class="min-h-[60vh] flex items-center justify-center px-4 py-12">
    <div class="max-w-2xl w-full bg-white border border-gray-100 shadow-sm rounded-2xl p-10 text-center">
      <div class="mx-auto mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-blue-50 text-blue-600">
        <i class="pi pi-users text-3xl"></i>
      </div>
      <h1 class="text-3xl font-bold text-gray-900">Leads Coming Soon</h1>
      <p class="mt-4 text-gray-600 leading-relaxed">
        Kami sedang menyiapkan fitur manajemen leads agar tim sales bisa meng-organisir prospek dengan rapi,
        mengatur tindak lanjut, dan berkolaborasi lebih mudah.
      </p>
      <p class="mt-3 text-sm text-gray-400">
        Update akan segera hadir — pantau terus untuk kabar selanjutnya.
      </p>
      <div class="mt-8 flex flex-wrap items-center justify-center gap-3">
        <router-link
          to="/"
          class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow hover:bg-blue-700 transition"
        >
          <i class="pi pi-home text-sm"></i>
          Kembali ke Dashboard
        </router-link>
        <button
          type="button"
          class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 transition"
          disabled
        >
          <i class="pi pi-clock text-sm"></i>
          Notifikasi segera tersedia
        </button>
      </div>
    </div>
  </div>
</template>
