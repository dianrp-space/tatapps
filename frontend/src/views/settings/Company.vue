<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Profil Perusahaan</h1>
        <p class="mt-1 text-sm text-gray-500">Perbarui identitas dan kontak utama perusahaan.</p>
      </div>
      <button
        type="button"
        class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 transition"
        @click="resetForm"
      >
        <i class="pi pi-undo text-sm"></i>
        Reset
      </button>
    </div>

    <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
      <form class="grid gap-6 lg:grid-cols-2" @submit.prevent="saveCompany">
        <section class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900">Identitas Perusahaan</h2>
          <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
            <div
              class="flex h-24 w-24 items-center justify-center overflow-hidden rounded-xl border border-gray-200 bg-gray-50 text-sm font-semibold text-gray-400 shadow-inner"
            >
              <img v-if="form.logo" :src="form.logo" alt="Logo Perusahaan" class="h-full w-full object-cover" />
              <span v-else>Logo</span>
            </div>
            <div class="flex flex-col gap-2">
              <label class="inline-flex w-max cursor-pointer items-center gap-2 rounded-lg border border-blue-200 px-4 py-2 text-sm font-medium text-blue-600 hover:bg-blue-50 transition">
                <i class="pi pi-upload text-xs"></i>
                <span>Browse Logo</span>
                <input type="file" accept="image/*" class="hidden" @change="handleLogoUpload" />
              </label>
              <button
                type="button"
                class="inline-flex w-max items-center gap-2 rounded-lg border border-gray-200 px-3 py-1 text-xs font-medium text-gray-600 hover:bg-gray-50 transition"
                @click="clearLogo"
                v-if="form.logo"
              >
                Hapus Logo
              </button>
            </div>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Nama Perusahaan</label>
            <input
              v-model="form.companyName"
              type="text"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              required
            />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Bidang Usaha</label>
            <input
              v-model="form.industry"
              type="text"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
            />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Alamat</label>
            <textarea
              v-model="form.address"
              rows="3"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
            ></textarea>
          </div>
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Provinsi</label>
              <input
                v-model="form.province"
                type="text"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Kabupaten / Kota</label>
              <input
                v-model="form.city"
                type="text"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
          </div>
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Kecamatan</label>
              <input
                v-model="form.district"
                type="text"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Kelurahan</label>
              <input
                v-model="form.subDistrict"
                type="text"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
          </div>
        </section>

        <section class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900">Kontak & Waktu Operasional</h2>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Zona Waktu</label>
            <select
              v-model="form.timezone"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
            >
              <option value="WIB (+7)">WIB (+7)</option>
              <option value="WITA (+8)">WITA (+8)</option>
              <option value="WIT (+9)">WIT (+9)</option>
            </select>
          </div>
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Kode Pos</label>
              <input
                v-model="form.postalCode"
                type="text"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Website</label>
              <input
                v-model="form.website"
                type="url"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
          </div>
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Nomor Telepon</label>
              <input
                v-model="form.phone"
                type="tel"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Nomor Fax</label>
              <input
                v-model="form.fax"
                type="tel"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
          </div>
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Email</label>
              <input
                v-model="form.email"
                type="email"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
          </div>
        </section>

        <div class="lg:col-span-2 flex justify-end gap-3 border-t border-gray-100 pt-6">
          <button
            type="button"
            class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 transition"
            @click="resetForm"
          >
            Batal
          </button>
          <button
            type="submit"
            class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow hover:bg-blue-700 transition"
          >
            Simpan Profil
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'

const STORAGE_KEY = 'tatapps_company_profile'

const defaultForm = {
  logo: '',
  companyName: '',
  industry: '',
  address: '',
  province: '',
  city: '',
  district: '',
  subDistrict: '',
  timezone: 'WIB (+7)',
  postalCode: '',
  phone: '',
  fax: '',
  email: '',
  website: ''
}

const form = reactive({ ...defaultForm })

const loadProfile = () => {
  if (typeof window === 'undefined') return
  const raw = window.localStorage.getItem(STORAGE_KEY)
  if (!raw) return
  try {
    const parsed = JSON.parse(raw)
    Object.assign(form, { ...defaultForm, ...parsed })
  } catch (error) {
    console.warn('Failed to load company profile', error)
  }
}

const persistProfile = () => {
  if (typeof window === 'undefined') return
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(form))
}

const resetForm = () => {
  Object.assign(form, defaultForm)
}

const saveCompany = () => {
  persistProfile()
  alert('Profil perusahaan disimpan.')
}

const handleLogoUpload = (event) => {
  const [file] = event.target.files || []
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    form.logo = typeof reader.result === 'string' ? reader.result : ''
  }
  reader.readAsDataURL(file)
}

const clearLogo = () => {
  form.logo = ''
}

onMounted(() => {
  loadProfile()
})
</script>
