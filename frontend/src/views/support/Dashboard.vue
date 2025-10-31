<template>
  <div class="min-h-[60vh] flex items-center justify-center px-4 py-12">
    <div class="max-w-2xl w-full bg-white border border-gray-100 shadow-sm rounded-2xl p-10 text-center">
      <div class="mx-auto mb-6 flex h-20 w-20 items-center justify-center rounded-full bg-sky-50 text-sky-600">
        <i class="pi pi-life-ring text-3xl"></i>
      </div>
      <h1 class="text-3xl font-bold text-gray-900">Support Center Coming Soon</h1>
      <p class="mt-4 text-gray-600 leading-relaxed">
        Kami sedang menyiapkan portal support terpadu untuk pencatatan tiket, panduan solusi,
        dan konsultasi langsung dengan tim kami agar setiap kendala dapat ditangani lebih cepat.
      </p>
      <p class="mt-3 text-sm text-gray-400">
        Sementara itu, gunakan saluran komunikasi yang tersedia — notifikasi fitur baru akan kami kirimkan segera.
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
      <div class="mt-10 grid gap-4 sm:grid-cols-3 text-left">
        <div class="space-y-2 rounded-xl border border-gray-100 bg-sky-50/60 p-4">
          <div class="flex h-10 w-10 items-center justify-center rounded-full bg-white text-sky-500 shadow">
            <i class="pi pi-ticket text-lg"></i>
          </div>
          <h3 class="font-semibold text-gray-800 text-sm">Ticketing System</h3>
          <p class="text-xs text-gray-500 leading-relaxed">
            Ajukan dan pantau status tiket dukungan secara real-time.
          </p>
        </div>
        <div class="space-y-2 rounded-xl border border-gray-100 bg-white p-4">
          <div class="flex h-10 w-10 items-center justify-center rounded-full bg-sky-100 text-sky-600 shadow">
            <i class="pi pi-book text-lg"></i>
          </div>
          <h3 class="font-semibold text-gray-800 text-sm">Knowledge Base</h3>
          <p class="text-xs text-gray-500 leading-relaxed">
            Kumpulan solusi dan panduan langkah demi langkah untuk tim Anda.
          </p>
        </div>
        <div class="space-y-2 rounded-xl border border-gray-100 bg-white p-4">
          <div class="flex h-10 w-10 items-center justify-center rounded-full bg-sky-100 text-sky-600 shadow">
            <i class="pi pi-comments text-lg"></i>
          </div>
          <h3 class="font-semibold text-gray-800 text-sm">Live Assistance</h3>
          <p class="text-xs text-gray-500 leading-relaxed">
            Hubungkan langsung dengan tim support saat prioritas tinggi.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>
