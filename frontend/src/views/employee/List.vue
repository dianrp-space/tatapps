<template>
  <div class="space-y-6">
    <div class="flex justify-start" v-if="canCreateEmployee">
      <RouterLink
        :to="{ name: 'EmployeeCreate' }"
        class="inline-flex items-center gap-2 rounded-lg bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow hover:bg-blue-700 transition"
      >
        <i class="pi pi-plus text-sm"></i>
        Tambah Karyawan
      </RouterLink>
    </div>

    <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
      <div class="flex flex-col gap-4 md:flex-row md:items-end md:justify-between">
        <div class="grid gap-4 sm:grid-cols-2 md:grid-cols-4 md:gap-6">
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Pencarian</label>
            <input
              v-model="filters.search"
              type="text"
              placeholder="Masukkan nama, email, atau kode karyawan"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
            />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Divisi</label>
            <select
              v-model="filters.divisionId"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
            >
              <option value="">Semua Divisi</option>
              <option v-for="division in divisionOptions" :key="division.id" :value="division.id">
                {{ division.name }}
              </option>
            </select>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Status</label>
            <select
              v-model="filters.status"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
            >
              <option value="">Semua</option>
              <option value="Aktif">Aktif</option>
              <option value="Cuti">Cuti</option>
              <option value="Resign">Resign</option>
            </select>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Jenis Kepegawaian</label>
            <select
              v-model="filters.employmentType"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-100"
            >
              <option value="">Semua</option>
              <option value="Tetap">Tetap</option>
              <option value="Kontrak">Kontrak</option>
              <option value="Magang">Magang</option>
            </select>
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

    <div
      v-if="canDeleteEmployee"
      class="mt-4 flex flex-wrap items-center justify-between gap-3 text-xs text-gray-500"
    >
      <span>Dipilih: {{ selectedIds.length }}</span>
      <div class="flex items-center gap-2">
        <button
          type="button"
          class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-3 py-1.5 text-xs font-medium text-gray-600 hover:bg-gray-50 transition disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="!hasSelection || isBatchDeleting"
          @click="clearSelection"
        >
          <i class="pi pi-times text-xs"></i>
          Bersihkan Pilihan
        </button>
        <button
          type="button"
          class="inline-flex items-center gap-2 rounded-lg bg-red-500 px-3 py-1.5 text-xs font-medium text-white shadow hover:bg-red-600 transition disabled:cursor-not-allowed disabled:opacity-60"
          :disabled="!hasSelection || isBatchDeleting"
          @click="deleteSelected"
        >
          <i :class="isBatchDeleting ? 'pi pi-spin pi-spinner' : 'pi pi-trash'" class="text-xs"></i>
          Hapus Dipilih
        </button>
      </div>
    </div>

    <div class="mt-6 overflow-hidden rounded-xl border border-gray-100">
      <table class="min-w-full divide-y divide-gray-100 bg-white text-sm">
        <thead class="bg-gray-50 text-left text-xs font-semibold uppercase tracking-wide text-gray-500">
          <tr>
            <th v-if="canDeleteEmployee" class="w-12 px-4 py-3">
              <input
                ref="headerCheckbox"
                type="checkbox"
                :checked="allVisibleSelected"
                @change="toggleSelectAll($event.target.checked)"
                class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
            </th>
            <th class="px-6 py-3">NIK</th>
            <th class="px-6 py-3">Nama</th>
            <th class="px-6 py-3">Divisi</th>
            <th class="px-6 py-3">Jabatan</th>
            <th class="px-6 py-3">Status</th>
            <th class="px-6 py-3">Email</th>
              <th
                v-if="canUpdateEmployee || canDeleteEmployee"
                class="px-6 py-3 text-right"
              >
                Aksi
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="employee in filteredEmployees" :key="employee.id" class="hover:bg-gray-50/60">
              <td v-if="canDeleteEmployee" class="px-4 py-3">
                <input
                  type="checkbox"
                  :value="employee.id"
                  v-model="selectedIds"
                  class="h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
              </td>
              <td class="px-6 py-3 font-medium text-gray-700">{{ employee.nik || employee.employeeCode }}</td>
              <td class="px-6 py-3">
                <div class="flex items-center gap-3">
                  <div
                    class="flex h-10 w-10 items-center justify-center overflow-hidden rounded-full border border-gray-200 bg-gray-50 text-xs font-semibold text-gray-400"
                  >
                    <img
                      v-if="employee.photo"
                      :src="employee.photo"
                      :alt="employee.fullName"
                      class="h-full w-full object-cover"
                    />
                    <span v-else>{{ avatarInitial(employee.fullName) }}</span>
                  </div>
                  <div>
                    <p class="font-semibold text-gray-900">{{ employee.fullName }}</p>
                    <p class="text-xs text-gray-500">Masuk {{ formatJoinDate(employee.joinDate) }}</p>
                  </div>
                </div>
              </td>
              <td class="px-6 py-3 text-gray-700">{{ divisionName(employee.divisionId) }}</td>
              <td class="px-6 py-3 text-gray-700">{{ positionName(employee.positionId) }}</td>
              <td class="px-6 py-3">
                <span :class="statusBadgeClass(employee.status)">{{ employee.status }}</span>
              </td>
              <td class="px-6 py-3 text-gray-700">{{ employee.email || '-' }}</td>
              <td
                v-if="canUpdateEmployee || canDeleteEmployee"
                class="px-6 py-3 text-right"
              >
                <div class="flex justify-end gap-2">
                  <RouterLink
                    v-if="canUpdateEmployee"
                    :to="{ name: 'EmployeeEdit', params: { id: employee.id } }"
                    class="rounded-full border border-gray-200 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50"
                  >
                    Edit
                  </RouterLink>
                  <button
                    v-if="canDeleteEmployee"
                    class="rounded-full border border-red-200 px-3 py-1 text-xs text-red-600 hover:bg-red-50"
                    @click="deleteEmployee(employee.id)"
                  >
                    Hapus
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredEmployees.length === 0">
              <td
                :colspan="emptyTableColspan"
                class="px-6 py-10 text-center text-sm text-gray-500"
              >
                Belum ada data karyawan. Tambahkan data baru melalui tombol "Tambah Karyawan".
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted, watch } from 'vue'
import { RouterLink } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useHRStore } from '@/stores/hr'
import { useAuthStore } from '@/stores/auth'

const hrStore = useHRStore()
const authStore = useAuthStore()
const toast = useToast()

onMounted(() => {
  hrStore.ensureHydrated().catch((error) => {
    const detail = error?.response?.data?.error || 'Gagal memuat data karyawan.'
    toast.add({ severity: 'error', summary: 'Load gagal', detail, life: 4000 })
  })
})

const canCreateEmployee = computed(() => authStore.hasPermission('employee.create'))
const canUpdateEmployee = computed(() => authStore.hasPermission('employee.update'))
const canDeleteEmployee = computed(() => authStore.hasPermission('employee.delete'))
const selectedIds = ref([])
const headerCheckbox = ref(null)
const isBatchDeleting = ref(false)

const filters = ref({
  search: '',
  divisionId: '',
  status: '',
  employmentType: ''
})

const divisionOptions = computed(() => hrStore.divisions.map((division) => ({
  id: division.id,
  name: division.name
})))

const filteredEmployees = computed(() => {
  return hrStore.employees.filter((employee) => {
    const matchesSearch = filters.value.search
      ? [employee.nik, employee.employeeCode, employee.fullName, employee.email]
          .filter(Boolean)
          .some((field) => field.toLowerCase().includes(filters.value.search.toLowerCase()))
      : true

    const matchesDivision = filters.value.divisionId
      ? employee.divisionId === filters.value.divisionId
      : true

    const matchesStatus = filters.value.status ? employee.status === filters.value.status : true
    const matchesType = filters.value.employmentType ? employee.employmentType === filters.value.employmentType : true

    return matchesSearch && matchesDivision && matchesStatus && matchesType
  })
})

const hasSelection = computed(() => selectedIds.value.length > 0)
const allVisibleSelected = computed(() => {
  if (!filteredEmployees.value.length) return false
  const visibleIds = filteredEmployees.value.map((employee) => employee.id)
  return visibleIds.every((id) => selectedIds.value.includes(id))
})

const statusBadgeClass = (status) => {
  switch (status) {
    case 'Aktif':
      return 'inline-flex items-center gap-2 rounded-full bg-emerald-50 px-3 py-1 text-xs font-medium text-emerald-600'
    case 'Cuti':
      return 'inline-flex items-center gap-2 rounded-full bg-amber-50 px-3 py-1 text-xs font-medium text-amber-600'
    case 'Resign':
      return 'inline-flex items-center gap-2 rounded-full bg-rose-50 px-3 py-1 text-xs font-medium text-rose-600'
    default:
      return 'inline-flex items-center gap-2 rounded-full bg-gray-100 px-3 py-1 text-xs font-medium text-gray-500'
  }
}

const resetFilters = () => {
  filters.value = { search: '', divisionId: '', status: '', employmentType: '' }
}

const toggleSelectAll = (checked) => {
  const visibleIds = filteredEmployees.value.map((employee) => employee.id)
  if (!visibleIds.length) {
    return
  }
  if (checked) {
    const combined = new Set([...selectedIds.value, ...visibleIds])
    selectedIds.value = Array.from(combined)
  } else {
    const removeSet = new Set(visibleIds)
    selectedIds.value = selectedIds.value.filter((id) => !removeSet.has(id))
  }
}

const clearSelection = () => {
  selectedIds.value = []
}

const deleteEmployee = async (id) => {
  if (!canDeleteEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk menghapus karyawan.', life: 3500 })
    return
  }
  const confirmation = confirm('Hapus data karyawan ini?')
  if (!confirmation) return

  try {
    await hrStore.deleteEmployeesBatch([id])
    selectedIds.value = selectedIds.value.filter((selected) => selected !== id)
    toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Data karyawan telah dihapus.', life: 3000 })
  } catch (error) {
    const detail = error?.response?.data?.error || 'Gagal menghapus data karyawan.'
    toast.add({ severity: 'error', summary: 'Gagal', detail, life: 4000 })
  }
}

const deleteSelected = async () => {
  if (!canDeleteEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk menghapus karyawan.', life: 3500 })
    return
  }
  if (!selectedIds.value.length) return
  const confirmation = confirm(`Hapus ${selectedIds.value.length} karyawan terpilih?`)
  if (!confirmation) return

  isBatchDeleting.value = true
  const ids = [...selectedIds.value]
  try {
    await hrStore.deleteEmployeesBatch(ids)
    selectedIds.value = []
    toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Data karyawan terpilih telah dihapus.', life: 3000 })
  } catch (error) {
    const detail = error?.response?.data?.error || 'Gagal menghapus data karyawan.'
    toast.add({ severity: 'error', summary: 'Gagal', detail, life: 4000 })
  } finally {
    isBatchDeleting.value = false
  }
}

watch(
  () => hrStore.employees,
  (list) => {
    const existing = new Set(list.map((employee) => employee.id))
    selectedIds.value = selectedIds.value.filter((id) => existing.has(id))
  }
)

watch(
  [selectedIds, filteredEmployees],
  () => {
    if (!headerCheckbox.value) return
    const total = filteredEmployees.value.length
    const selectedCount = filteredEmployees.value.filter((employee) => selectedIds.value.includes(employee.id)).length
    headerCheckbox.value.indeterminate = selectedCount > 0 && selectedCount < total
    headerCheckbox.value.checked = total > 0 && selectedCount === total
  },
  { flush: 'post' }
)

const emptyTableColspan = computed(() => {
  let columns = 6
  if (canUpdateEmployee.value || canDeleteEmployee.value) {
    columns += 1
  }
  if (canDeleteEmployee.value) {
    columns += 1
  }
  return columns
})

const divisionName = (divisionId) => {
  const division = hrStore.findDivision(divisionId)
  return division ? division.name : '-'
}

const positionName = (positionId) => {
  const position = hrStore.findPosition(positionId)
  return position ? position.title : '-'
}

const formatJoinDate = (isoString) => {
  if (!isoString) return '-'
  try {
    const date = new Date(isoString)
    return date.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' })
  } catch {
    return isoString
  }
}

const avatarInitial = (name) => {
  if (!name) return '?'
  return name.trim().charAt(0).toUpperCase()
}
</script>
