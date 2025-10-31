<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Tambah Karyawan</h1>
        <p class="mt-1 text-sm text-gray-500">Lengkapi data personal dan penempatan untuk menambahkan karyawan baru.</p>
      </div>
      <RouterLink
        :to="{ name: 'EmployeeData' }"
        class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 transition"
      >
        <i class="pi pi-arrow-left text-sm"></i>
        Kembali ke Data Karyawan
      </RouterLink>
    </div>

    <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
      <form class="grid gap-6 lg:grid-cols-2" @submit.prevent="handleSubmit">
        <section class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900">Informasi Personal</h2>
          <div class="grid gap-4 sm:grid-cols-2">
            <div class="sm:col-span-2">
              <label class="mb-2 block text-sm font-medium text-gray-700">NIK Karyawan <span class="text-red-500">*</span></label>
              <input
                v-model="form.nik"
                type="text"
                required
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
            <div class="sm:col-span-2">
              <label class="mb-2 block text-sm font-medium text-gray-700">Nama Lengkap <span class="text-red-500">*</span></label>
              <input
                v-model="form.fullName"
                type="text"
                required
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Tempat Lahir</label>
              <input
                v-model="form.birthPlace"
                type="text"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Tanggal Lahir</label>
              <Calendar
                v-model="form.birthDate"
                dateFormat="dd-mm-yy"
                showIcon
                iconDisplay="input"
                :inputClass="'w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100'"
                placeholder="dd-mm-yyyy"
              />
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Jenis Kelamin</label>
              <select
                v-model="form.gender"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              >
                <option value="">Pilih</option>
                <option value="Pria">Pria</option>
                <option value="Wanita">Wanita</option>
              </select>
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Golongan Darah</label>
              <select
                v-model="form.bloodType"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              >
                <option value="">Pilih</option>
                <option value="A">A</option>
                <option value="B">B</option>
                <option value="AB">AB</option>
                <option value="O">O</option>
              </select>
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Status Perkawinan</label>
              <select
                v-model="form.maritalStatus"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              >
                <option value="">Pilih</option>
                <option value="Belum Kawin">Belum Kawin</option>
                <option value="Kawin">Kawin</option>
                <option value="Cerai">Cerai</option>
              </select>
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Agama</label>
              <select
                v-model="form.religion"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              >
                <option value="">Pilih</option>
                <option value="Islam">Islam</option>
                <option value="Kristen">Kristen</option>
                <option value="Katolik">Katolik</option>
                <option value="Hindu">Hindu</option>
                <option value="Buddha">Buddha</option>
                <option value="Khonghucu">Khonghucu</option>
              </select>
            </div>
          </div>
        </section>

        <section class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900">Informasi Identitas & Kontak</h2>
          <div class="grid gap-4 sm:grid-cols-2">
            <div class="sm:col-span-2">
              <label class="mb-2 block text-sm font-medium text-gray-700">Foto Karyawan</label>
              <div class="flex flex-col gap-3 sm:flex-row sm:items-center">
                <div
                  class="flex h-20 w-20 items-center justify-center overflow-hidden rounded-full border border-gray-200 bg-gray-50 text-sm font-semibold text-gray-400"
                >
                  <img
                    v-if="form.photo"
                    :src="form.photo"
                    alt="Foto Karyawan"
                    class="h-full w-full object-cover"
                  />
                  <span v-else>Foto</span>
                </div>
                <div class="flex flex-col gap-2">
                  <label class="inline-flex w-max cursor-pointer items-center gap-2 rounded-lg border border-blue-200 px-4 py-2 text-sm font-medium text-blue-600 hover:bg-blue-50 transition">
                    <i class="pi pi-upload text-xs"></i>
                    <span>Browse Foto</span>
                    <input type="file" accept="image/*" class="hidden" @change="handlePhotoChange" />
                  </label>
                  <button
                    v-if="form.photo"
                    type="button"
                    class="inline-flex w-max items-center gap-2 rounded-lg border border-gray-200 px-3 py-1 text-xs font-medium text-gray-600 hover:bg-gray-50 transition"
                    @click="clearPhoto"
                  >
                    Bersihkan Foto
                  </button>
                </div>
              </div>
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Jenis Identitas</label>
              <select
                v-model="form.identityType"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              >
                <option value="">Pilih</option>
                <option value="KTP">KTP</option>
                <option value="SIM">SIM</option>
                <option value="Paspor">Paspor</option>
              </select>
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Nomor Identitas</label>
              <input
                v-model="form.identityNumber"
                type="text"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              />
            </div>
            <div class="sm:col-span-2">
              <label class="mb-2 block text-sm font-medium text-gray-700">Alamat KTP</label>
              <textarea
                v-model="form.addressKtp"
                rows="3"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              ></textarea>
            </div>
            <div class="sm:col-span-2 flex items-center gap-3">
              <input
                id="sync-address"
                v-model="form.syncAddress"
                type="checkbox"
                class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
              <label for="sync-address" class="text-sm text-gray-600">Alamat domisili sama dengan alamat KTP</label>
            </div>
            <div class="sm:col-span-2">
              <label class="mb-2 block text-sm font-medium text-gray-700">Alamat Domisili</label>
              <textarea
                v-model="form.addressDomicile"
                rows="3"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              ></textarea>
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Nomor HP</label>
              <input
                v-model="form.phone"
                type="tel"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
                placeholder="08xxxxxxxxxx"
              />
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Email</label>
              <input
                v-model="form.email"
                type="email"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
                placeholder="nama@perusahaan.co.id"
              />
            </div>
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
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Tanggal Masuk</label>
              <Calendar
                v-model="form.joinDate"
                dateFormat="dd-mm-yy"
                showIcon
                iconDisplay="input"
                :inputClass="'w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100'"
                placeholder="dd-mm-yyyy"
              />
            </div>
          </div>
        </section>

        <section class="space-y-4 lg:col-span-2">
          <h2 class="text-lg font-semibold text-gray-900">Penempatan Organisasi</h2>
          <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Divisi</label>
              <select
                v-model="form.divisionId"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              >
                <option value="">Pilih divisi</option>
                <option v-for="division in divisionOptions" :key="division.id" :value="division.id">{{ division.name }}</option>
              </select>
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Jabatan</label>
              <select
                v-model="form.positionId"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              >
                <option value="">Pilih jabatan</option>
                <option v-for="position in positionOptions" :key="position.id" :value="position.id">
                  {{ position.title }}
                </option>
              </select>
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Status Kepegawaian</label>
              <select
                v-model="form.status"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              >
                <option value="Aktif">Aktif</option>
                <option value="Cuti">Cuti</option>
                <option value="Resign">Resign</option>
              </select>
            </div>
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Jenis Kepegawaian</label>
              <select
                v-model="form.employmentType"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              >
                <option value="">Pilih jenis</option>
                <option value="Tetap">Tetap</option>
                <option value="Kontrak">Kontrak</option>
                <option value="Magang">Magang</option>
              </select>
            </div>
          </div>
        </section>

        <div class="lg:col-span-2 flex justify-end gap-3 border-t border-gray-100 pt-6">
          <RouterLink
            :to="{ name: 'EmployeeData' }"
            class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 transition"
          >
            Batal
          </RouterLink>
          <button
            type="submit"
            :disabled="!canCreateEmployee || isSubmitting"
            :class="[
              'inline-flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium text-white shadow transition',
              (!canCreateEmployee || isSubmitting) ? 'bg-blue-400 cursor-not-allowed opacity-70' : 'bg-blue-600 hover:bg-blue-700'
            ]"
          >
            <i v-if="isSubmitting" class="pi pi-spin pi-spinner"></i>
            <span>{{ isSubmitting ? 'Menyimpan...' : 'Simpan Karyawan' }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useHRStore } from '@/stores/hr'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const hrStore = useHRStore()
const authStore = useAuthStore()
const toast = useToast()

onMounted(() => {
  hrStore.ensureHydrated().catch((error) => {
    const detail = error?.response?.data?.error || 'Gagal memuat referensi karyawan.'
    toast.add({ severity: 'error', summary: 'Load gagal', detail, life: 4000 })
  })
})

const canCreateEmployee = computed(() => authStore.hasPermission('employee.create'))
const isSubmitting = ref(false)

const form = ref({
  nik: '',
  fullName: '',
  birthPlace: '',
  birthDate: null,
  gender: '',
  bloodType: '',
  maritalStatus: '',
  religion: '',
  identityType: '',
  identityNumber: '',
  addressKtp: '',
  addressDomicile: '',
  syncAddress: true,
  phone: '',
  timezone: 'WIB (+7)',
  joinDate: new Date(),
  divisionId: '',
  positionId: '',
  status: 'Aktif',
  employmentType: '',
  email: '',
  photo: ''
})

watch(
  () => form.value.addressKtp,
  (value) => {
    if (form.value.syncAddress) {
      form.value.addressDomicile = value
    }
  }
)

watch(
  () => form.value.syncAddress,
  (sync) => {
    if (sync) {
      form.value.addressDomicile = form.value.addressKtp
    }
  }
)

const divisionOptions = computed(() => hrStore.divisions.map((division) => ({ id: division.id, name: division.name })))

const positionOptions = computed(() => {
  const all = hrStore.positions
  if (!form.value.divisionId) {
    return all
  }
  const filtered = all.filter((position) => position.divisionId === form.value.divisionId)
  return filtered.length ? filtered : all
})

watch(
  () => form.value.divisionId,
  () => {
    if (form.value.positionId) {
      const exists = hrStore.positions.some((position) => {
        if (position.id !== form.value.positionId) return false
        if (!form.value.divisionId) return true
        return position.divisionId === form.value.divisionId
      })
      if (!exists) {
        form.value.positionId = ''
      }
    }
  }
)

const handleSubmit = async () => {
  if (!canCreateEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk menambah karyawan.', life: 3500 })
    return
  }
  if (!form.value.nik.trim() || !form.value.fullName.trim()) {
    toast.add({ severity: 'warn', summary: 'Validasi', detail: 'NIK karyawan dan nama wajib diisi.', life: 3500 })
    return
  }

  const payload = {
    employeeCode: form.value.nik.trim(),
    nik: form.value.nik.trim(),
    fullName: form.value.fullName.trim(),
    birthPlace: form.value.birthPlace.trim(),
    birthDate: form.value.birthDate ? toISODate(form.value.birthDate) : null,
    gender: form.value.gender,
    bloodType: form.value.bloodType,
    maritalStatus: form.value.maritalStatus,
    religion: form.value.religion,
    identityType: form.value.identityType,
    identityNumber: form.value.identityNumber.trim(),
    addressKtp: form.value.addressKtp.trim(),
    addressDomicile: form.value.addressDomicile.trim(),
    address: form.value.addressDomicile.trim(),
    phone: form.value.phone.trim(),
    email: form.value.email.trim(),
    timezone: form.value.timezone,
    divisionId: form.value.divisionId || null,
    positionId: form.value.positionId || null,
    employmentType: form.value.employmentType || '',
    status: form.value.status || 'Aktif',
    joinDate: form.value.joinDate ? toISODate(form.value.joinDate) : new Date().toISOString(),
    photo: form.value.photo || ''
  }

  isSubmitting.value = true
  try {
    await hrStore.createEmployee(payload)
    toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Data karyawan berhasil ditambahkan.', life: 3000 })
    router.push({ name: 'EmployeeData' })
  } catch (error) {
    const detail = error?.response?.data?.error || 'Gagal menambah data karyawan.'
    toast.add({ severity: 'error', summary: 'Gagal', detail, life: 4000 })
  } finally {
    isSubmitting.value = false
  }
}

function toISODate(value) {
  if (!value) return null
  const date = value instanceof Date ? value : new Date(value)
  const utc = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()))
  return utc.toISOString()
}

const handlePhotoChange = (event) => {
  const [file] = event.target.files || []
  if (!file) return
  const reader = new FileReader()
  reader.onload = () => {
    form.value.photo = typeof reader.result === 'string' ? reader.result : ''
  }
  reader.readAsDataURL(file)
}

const clearPhoto = () => {
  form.value.photo = ''
}
</script>
