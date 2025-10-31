<template>
  <div class="space-y-6">
    <div class="flex justify-start" v-if="canCreateEmployee">
      <button
        type="button"
        class="inline-flex items-center gap-2 rounded-lg bg-purple-600 px-4 py-2 text-sm font-medium text-white shadow hover:bg-purple-700 transition"
        @click="openCreateModal"
      >
        <i class="pi pi-plus text-sm"></i>
        Tambah Divisi
      </button>
    </div>

    <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
      <div class="flex items-center justify-between">
        <h2 class="text-lg font-semibold text-gray-800">Daftar Divisi</h2>
        <span class="rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-500">
          Total {{ hrStore.divisions.length }} divisi
        </span>
      </div>

      <div class="mt-4 flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
        <div class="grid gap-4 md:grid-cols-3 md:gap-6">
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Pencarian</label>
            <input
              v-model="filters.search"
              type="text"
              placeholder="Nama divisi atau penanggung jawab"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-100"
            />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Status Rekrutmen</label>
            <select
              v-model="filters.recruitmentStatus"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-100"
            >
              <option value="">Semua</option>
              <option value="Butuh segera">Butuh segera</option>
              <option value="Sedang berjalan">Sedang berjalan</option>
              <option value="Stabil">Stabil</option>
            </select>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Minimal Anggota</label>
            <input
              v-model.number="filters.minMember"
              type="number"
              min="0"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-100"
            />
          </div>
        </div>
        <button
          type="button"
          class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 transition"
          @click="resetFilters"
        >
          <i class="pi pi-filter-slash text-sm"></i>
          Reset Filter
        </button>
      </div>

      <div class="mt-6 overflow-hidden rounded-xl border border-gray-100">
        <table class="min-w-full divide-y divide-gray-100 bg-white text-sm">
          <thead class="bg-gray-50 text-left text-xs font-semibold uppercase tracking-wide text-gray-500">
            <tr>
              <th class="px-6 py-3">Nama Divisi</th>
              <th class="px-6 py-3">Penanggung Jawab</th>
              <th class="px-6 py-3">Jumlah Anggota</th>
              <th class="px-6 py-3">Status Rekrutmen</th>
              <th
                v-if="canUpdateEmployee || canDeleteEmployee"
                class="px-6 py-3 text-right"
              >
                Aksi
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="division in filteredDivisions" :key="division.id" class="hover:bg-gray-50/70">
              <td class="px-6 py-3">
                <p class="font-semibold text-gray-900">{{ division.name }}</p>
                <p class="text-xs text-gray-500">{{ division.description || 'Belum ada deskripsi' }}</p>
              </td>
              <td class="px-6 py-3 text-gray-700">
                <div class="flex items-center gap-3">
                  <div
                    class="flex h-10 w-10 items-center justify-center overflow-hidden rounded-full border border-gray-200 bg-gray-50 text-xs font-semibold text-gray-400"
                  >
                    <img
                      v-if="divisionHeadPhoto(division)"
                      :src="divisionHeadPhoto(division)"
                      :alt="divisionHeadName(division)"
                      class="h-full w-full object-cover"
                    />
                    <span v-else>{{ avatarInitial(divisionHeadName(division)) }}</span>
                  </div>
                  <div>
                    <p class="font-medium">{{ divisionHeadName(division) }}</p>
                    <p class="text-xs text-gray-400">{{ divisionHeadTitle(division) }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-3 text-gray-700">{{ divisionMemberCount(division.id) }} orang</td>
              <td class="px-6 py-3">
                <span :class="recruitmentBadgeClass(division.recruitmentStatus)">
                  {{ division.recruitmentStatus || 'Stabil' }}
                </span>
              </td>
              <td
                v-if="canUpdateEmployee || canDeleteEmployee"
                class="px-6 py-3 text-right"
              >
                <div class="flex justify-end gap-2">
                  <button
                    v-if="canUpdateEmployee"
                    class="rounded-full border border-gray-200 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50"
                    @click="openEditModal(division)"
                  >
                    Edit
                  </button>
                  <button
                    v-if="canDeleteEmployee"
                    class="rounded-full border border-red-200 px-3 py-1 text-xs text-red-600 hover:bg-red-50"
                    @click="deleteDivision(division.id)"
                  >
                    Hapus
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredDivisions.length === 0">
              <td
                :colspan="(canUpdateEmployee || canDeleteEmployee) ? 5 : 4"
                class="px-6 py-10 text-center text-sm text-gray-500"
              >
                Belum ada divisi yang terdaftar.
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <div
      v-if="showForm"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
      @click.self="closeModal"
    >
      <div class="w-full max-w-3xl rounded-2xl bg-white p-6 shadow-xl">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-xl font-semibold text-gray-900">{{ isEditing ? 'Edit Divisi' : 'Tambah Divisi' }}</h2>
            <p class="text-sm text-gray-500">Isi detail struktur organisasi sesuai dengan kebijakan perusahaan.</p>
          </div>
          <button class="text-gray-400 hover:text-gray-600" @click="closeModal">
            <i class="pi pi-times text-lg"></i>
          </button>
        </div>

        <form class="mt-6 grid gap-4 md:grid-cols-2" @submit.prevent="saveDivision">
          <div class="md:col-span-2">
            <label class="mb-2 block text-sm font-medium text-gray-700">Nama Divisi <span class="text-red-500">*</span></label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-100"
            />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Penanggung Jawab</label>
            <select
              v-model="form.headEmployeeId"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-100"
            >
              <option value="">Tidak ada</option>
              <option v-for="employee in employeeOptions" :key="employee.id" :value="employee.id">{{ employee.label }}</option>
            </select>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Jabatan Penanggung Jawab</label>
            <input
              type="text"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-100"
              :value="displayPosition(form.headPositionId)"
              disabled
            />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Status Rekrutmen</label>
            <select
              v-model="form.recruitmentStatus"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-100"
            >
              <option value="Stabil">Stabil</option>
              <option value="Sedang berjalan">Sedang berjalan</option>
              <option value="Butuh segera">Butuh segera</option>
            </select>
          </div>
          <div class="md:col-span-2">
            <label class="mb-2 block text-sm font-medium text-gray-700">Deskripsi</label>
            <textarea
              v-model="form.description"
              rows="3"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-purple-500 focus:outline-none focus:ring-2 focus:ring-purple-100"
            ></textarea>
          </div>
        </form>

        <div class="mt-6 flex justify-end gap-3">
          <button
            type="button"
            class="rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50"
            @click="closeModal"
          >
            Batal
          </button>
          <button
            type="button"
            :disabled="isSubmitting || (isEditing && !canUpdateEmployee) || (!isEditing && !canCreateEmployee)"
            :class="[
              'rounded-lg px-4 py-2 text-sm font-medium text-white shadow transition',
              isSubmitting || (isEditing && !canUpdateEmployee) || (!isEditing && !canCreateEmployee)
                ? 'bg-purple-400 cursor-not-allowed opacity-70'
                : 'bg-purple-600 hover:bg-purple-700'
            ]"
            @click="saveDivision"
          >
            <i v-if="isSubmitting" class="pi pi-spin pi-spinner"></i>
            <span>{{ isSubmitting ? 'Menyimpan...' : (isEditing ? 'Simpan Perubahan' : 'Simpan Divisi') }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useHRStore } from '@/stores/hr'
import { useAuthStore } from '@/stores/auth'

const hrStore = useHRStore()
const authStore = useAuthStore()
const toast = useToast()

onMounted(() => {
  hrStore.ensureHydrated().catch((error) => {
    const detail = error?.response?.data?.error || 'Gagal memuat data divisi.'
    toast.add({ severity: 'error', summary: 'Load gagal', detail, life: 4000 })
  })
})

const filters = ref({
  search: '',
  recruitmentStatus: '',
  minMember: 0
})

const formTemplate = {
  name: '',
  description: '',
  headEmployeeId: '',
  headPositionId: '',
  recruitmentStatus: 'Stabil'
}

const form = ref({ ...formTemplate })
const showForm = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const isSubmitting = ref(false)
const canCreateEmployee = computed(() => authStore.hasPermission('employee.create'))
const canUpdateEmployee = computed(() => authStore.hasPermission('employee.update'))
const canDeleteEmployee = computed(() => authStore.hasPermission('employee.delete'))

const filteredDivisions = computed(() => {
  return hrStore.divisions.filter((division) => {
    const matchesSearch = filters.value.search
      ? [division.name, divisionHeadName(division), divisionHeadTitle(division)]
          .filter(Boolean)
          .some((field) => field.toLowerCase().includes(filters.value.search.toLowerCase()))
      : true

    const matchesStatus = filters.value.recruitmentStatus
      ? (division.recruitmentStatus || 'Stabil') === filters.value.recruitmentStatus
      : true

    const matchesMember = filters.value.minMember
      ? hrStore.employeeCountByDivision(division.id) >= Number(filters.value.minMember)
      : true

    return matchesSearch && matchesStatus && matchesMember
  })
})

const recruitmentBadgeClass = (status = 'Stabil') => {
  switch (status) {
    case 'Butuh segera':
      return 'inline-flex items-center rounded-full bg-rose-50 px-3 py-1 text-xs font-medium text-rose-600'
    case 'Sedang berjalan':
      return 'inline-flex items-center rounded-full bg-amber-50 px-3 py-1 text-xs font-medium text-amber-600'
    default:
      return 'inline-flex items-center rounded-full bg-emerald-50 px-3 py-1 text-xs font-medium text-emerald-600'
  }
}

const resetFilters = () => {
  filters.value = { search: '', recruitmentStatus: '', minMember: 0 }
}

const openCreateModal = () => {
  if (!canCreateEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk menambah divisi.', life: 3500 })
    return
  }
  isEditing.value = false
  editingId.value = null
  form.value = { ...formTemplate }
  updateHeadPositionFromEmployee()
  showForm.value = true
}

const openEditModal = (division) => {
  if (!canUpdateEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk memperbarui divisi.', life: 3500 })
    return
  }
  isEditing.value = true
  editingId.value = division.id
  form.value = {
    name: division.name,
    description: division.description || '',
    headEmployeeId: division.headEmployeeId || '',
    headPositionId: division.headPositionId || '',
    recruitmentStatus: division.recruitmentStatus || 'Stabil'
  }
  updateHeadPositionFromEmployee()
  showForm.value = true
}

const closeModal = () => {
  showForm.value = false
  form.value = { ...formTemplate }
  editingId.value = null
  isSubmitting.value = false
}

const saveDivision = () => {
  if (isEditing.value) {
    if (!canUpdateEmployee.value) {
      toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk memperbarui divisi.', life: 3500 })
      return
    }
  } else if (!canCreateEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk menambah divisi.', life: 3500 })
    return
  }

  if (!form.value.name) {
    toast.add({ severity: 'warn', summary: 'Validasi', detail: 'Nama divisi wajib diisi.', life: 3500 })
    return
  }

  const headEmployeeId = form.value.headEmployeeId || null
  const headPositionId = headEmployeeId ? findEmployeePosition(headEmployeeId) : null
  const headName = headEmployeeId ? displayEmployee(headEmployeeId) : ''
  const headTitle = headPositionId ? displayPosition(headPositionId) : ''

  const payload = {
    name: form.value.name.trim(),
    description: form.value.description.trim(),
    headEmployeeId,
    headPositionId,
    head: headName,
    headTitle: headTitle,
    recruitmentStatus: form.value.recruitmentStatus || 'Stabil'
  }

  isSubmitting.value = true
  const action = isEditing.value && editingId.value
    ? hrStore.updateDivision(editingId.value, payload)
    : hrStore.createDivision(payload)

  action
    .then(() => {
      toast.add({
        severity: 'success',
        summary: 'Berhasil',
        detail: isEditing.value ? 'Divisi diperbarui.' : 'Divisi berhasil ditambahkan.',
        life: 3000
      })
      closeModal()
    })
    .catch((error) => {
      const detail = error?.response?.data?.error || 'Gagal menyimpan data divisi.'
      toast.add({ severity: 'error', summary: 'Gagal', detail, life: 4000 })
    })
    .finally(() => {
      isSubmitting.value = false
    })
}

const deleteDivision = async (id) => {
  if (!canDeleteEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk menghapus divisi.', life: 3500 })
    return
  }
  if (!confirm('Hapus data divisi ini?')) {
    return
  }
  try {
    await hrStore.deleteDivision(id)
    toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Divisi telah dihapus.', life: 3000 })
  } catch (error) {
    const detail = error?.response?.data?.error || 'Gagal menghapus data divisi.'
    toast.add({ severity: 'error', summary: 'Gagal', detail, life: 4000 })
  }
}

const divisionMemberCount = (divisionId) => hrStore.employeeCountByDivision(divisionId)

const employeeOptions = computed(() =>
  hrStore.employees.map((employee) => ({
    id: employee.id,
    label: employee.fullName || employee.employeeCode || 'Tanpa Nama'
  }))
)

const positionOptions = computed(() =>
  hrStore.positions.map((position) => ({
    id: position.id,
    label: position.title || 'Tanpa Jabatan'
  }))
)

const divisionHeadName = (division) => {
  if (!division) return '-'
  if (division.headEmployeeId) {
    const employee = hrStore.findEmployee(division.headEmployeeId)
    if (employee) return employee.fullName
  }
  return division.head || '-'
}

const divisionHeadTitle = (division) => {
  if (!division) return '-'
  if (division.headPositionId) {
    const position = hrStore.findPosition(division.headPositionId)
    if (position) return position.title
  }
  return division.headTitle || '-'
}

const displayEmployee = (id) => {
  if (!id) return ''
  const employee = hrStore.findEmployee(id)
  return employee ? employee.fullName || employee.employeeCode || '' : ''
}

const displayPosition = (id) => {
  if (!id) return ''
  const position = hrStore.findPosition(id)
  return position ? position.title || '' : ''
}

const findEmployeePosition = (employeeId) => {
  if (!employeeId) return null
  const employee = hrStore.findEmployee(employeeId)
  return employee ? employee.positionId || null : null
}

const updateHeadPositionFromEmployee = () => {
  const positionId = findEmployeePosition(form.value.headEmployeeId)
  form.value.headPositionId = positionId || ''
}

watch(
  () => form.value.headEmployeeId,
  () => updateHeadPositionFromEmployee()
)

const divisionHeadPhoto = (division) => {
  if (!division || !division.headEmployeeId) return ''
  const employee = hrStore.findEmployee(division.headEmployeeId)
  return employee ? employee.photo || '' : ''
}

const avatarInitial = (name) => {
  if (!name || name === '-') return '?'
  return name.trim().charAt(0).toUpperCase()
}
</script>
