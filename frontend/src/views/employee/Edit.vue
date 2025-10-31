<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">Edit Karyawan</h1>
        <p class="mt-1 text-sm text-gray-500">Perbarui data personal dan penempatan karyawan.</p>
      </div>
      <RouterLink
        :to="{ name: 'EmployeeData' }"
        class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 transition"
      >
        <i class="pi pi-arrow-left text-sm"></i>
        Kembali ke Data Karyawan
      </RouterLink>
    </div>

    <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm" v-if="isReady">
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
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm foco-s:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
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
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm foco-s:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
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

        <section class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900">Penempatan Divisi</h2>
          <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
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

        <section class="space-y-4">
          <h2 class="text-lg font-semibold text-gray-900">Jabatan</h2>
          <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
            <div>
              <label class="mb-2 block text-sm font-medium text-gray-700">Jabatan</label>
              <select
                v-model="form.positionId"
                class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
              >
                <option value="">Pilih jabatan</option>
                <option v-for="position in hrStore.positions" :key="position.id" :value="position.id">
                  {{ position.title }}
                </option>
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
            :disabled="!canUpdateEmployee || isSubmitting"
            :class="[
              'inline-flex items-center gap-2 rounded-lg px-4 py-2 text-sm font-medium text-white shadow transition',
              (!canUpdateEmployee || isSubmitting) ? 'bg-blue-400 cursor-not-allowed opacity-70' : 'bg-blue-600 hover:bg-blue-700'
            ]"
          >
            <i v-if="isSubmitting" class="pi pi-spin pi-spinner"></i>
            <span>{{ isSubmitting ? 'Menyimpan...' : 'Simpan Perubahan' }}</span>
          </button>
        </div>
      </form>
    </div>

    <div v-else class="rounded-2xl border border-gray-100 bg-white p-8 text-center text-sm text-gray-500">
      Data karyawan tidak ditemukan.
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch, onMounted } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useHRStore } from '@/stores/hr'
import { useAuthStore } from '@/stores/auth'

const route = useRoute()
const router = useRouter()
const hrStore = useHRStore()
const authStore = useAuthStore()
const toast = useToast()

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
  email: '',
  timezone: 'WIB (+7)',
  joinDate: new Date(),
  divisionId: '',
  positionId: '',
  status: 'Aktif',
  employmentType: '',
  photo: ''
})

const isReady = ref(false)
const employeeId = Number(route.params.id)
const desiredPositionKey = ref('')
const desiredDivisionKey = ref('')
const canUpdateEmployee = computed(() => authStore.hasPermission('employee.update'))
const isSubmitting = ref(false)

onMounted(async () => {
  try {
    await hrStore.ensureHydrated()
  } catch (error) {
    const detail = error?.response?.data?.error || 'Gagal memuat data karyawan.'
    toast.add({ severity: 'error', summary: 'Load gagal', detail, life: 4000 })
    router.replace({ name: 'EmployeeData' })
    return
  }

  if (Number.isNaN(employeeId)) {
    toast.add({ severity: 'warn', summary: 'Tidak valid', detail: 'ID karyawan tidak valid.', life: 3500 })
    router.replace({ name: 'EmployeeData' })
    return
  }

  const employee = hrStore.findEmployee(employeeId)
  if (!employee) {
    toast.add({ severity: 'warn', summary: 'Tidak ditemukan', detail: 'Data karyawan tidak ditemukan.', life: 3500 })
    router.replace({ name: 'EmployeeData' })
    return
  }

  form.value.nik = employee.nik || employee.employeeCode || ''
  form.value.fullName = employee.fullName || ''
  form.value.birthPlace = employee.birthPlace || ''
  form.value.birthDate = parseDate(employee.birthDate)
  form.value.gender = employee.gender || ''
  form.value.bloodType = employee.bloodType || ''
  form.value.maritalStatus = employee.maritalStatus || ''
  form.value.religion = employee.religion || ''
  form.value.identityType = employee.identityType || ''
  form.value.identityNumber = employee.identityNumber || ''
  form.value.addressKtp = employee.addressKtp || ''
  form.value.addressDomicile = employee.addressDomicile || employee.address || ''
  form.value.phone = employee.phone || ''
  form.value.email = employee.email || ''
  form.value.timezone = employee.timezone || 'WIB (+7)'
  form.value.joinDate = parseDate(employee.joinDate) || new Date()

  const resolvedDivisionId = employee.divisionId || resolveDivisionId(employee.division)
  const resolvedPositionId = resolvePositionAny(employee.positionId || employee.position)

  form.value.divisionId = resolvedDivisionId || ''
  form.value.positionId = resolvedPositionId || ''

  if (!form.value.divisionId && resolvedPositionId) {
    const position = hrStore.findPosition(resolvedPositionId)
    if (position && position.divisionId) {
      form.value.divisionId = position.divisionId
    }
  }
  form.value.status = employee.status || 'Aktif'
  form.value.employmentType = employee.employmentType || ''
  form.value.syncAddress = form.value.addressDomicile === form.value.addressKtp
  form.value.photo = employee.photo || ''
  desiredPositionKey.value = employee.positionId || employee.position || ''
  desiredDivisionKey.value = employee.divisionId || employee.division || ''

  isReady.value = true
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

watch(
  () => hrStore.positions,
  () => {
    if (!desiredPositionKey.value) return
    const resolved = resolvePositionAny(desiredPositionKey.value)
    if (resolved) {
      form.value.positionId = resolved
    }
  },
  { deep: false }
)

const handleSubmit = async () => {
  if (!canUpdateEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk memperbarui karyawan.', life: 3500 })
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
    await hrStore.updateEmployee(employeeId, payload)
    toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Perubahan data karyawan disimpan.', life: 3000 })
    router.push({ name: 'EmployeeData' })
  } catch (error) {
    const detail = error?.response?.data?.error || 'Gagal menyimpan perubahan karyawan.'
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

function parseDate(value) {
  if (!value) return null
  const date = value instanceof Date ? value : new Date(value)
  return Number.isNaN(date.getTime()) ? null : date
}

function resolveDivisionId(divisionName) {
  if (!divisionName) return ''
  const match = hrStore.divisions.find((division) => division.name === divisionName)
  return match ? match.id : ''
}

function resolvePositionId(positionTitle) {
  if (!positionTitle) return ''
  const byId = hrStore.positions.find((position) => position.id === positionTitle)
  if (byId) return byId.id
  const match = hrStore.positions.find((position) => position.title === positionTitle)
  return match ? match.id : ''
}

function resolvePositionAny(key) {
  if (!key) return ''
  const byId = hrStore.positions.find((position) => position.id === key)
  if (byId) return byId.id
  const byTitle = hrStore.positions.find((position) => position.title === key)
  return byTitle ? byTitle.id : ''
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
