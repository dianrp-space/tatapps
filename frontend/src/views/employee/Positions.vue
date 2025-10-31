<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div class="flex flex-col gap-2 sm:flex-row sm:items-center">
        <button
          v-if="canCreateEmployee"
          type="button"
          class="inline-flex items-center gap-2 rounded-lg bg-emerald-600 px-4 py-2 text-sm font-medium text-white shadow hover:bg-emerald-700 transition"
          @click="openCreateModal"
        >
          <i class="pi pi-plus text-sm"></i>
          Tambah Jabatan
        </button>
        <RouterLink
          :to="{ name: 'EmployeeHierarchy' }"
          class="inline-flex items-center gap-2 rounded-lg border border-emerald-200 px-4 py-2 text-sm font-medium text-emerald-600 transition hover:bg-emerald-50"
        >
          <i class="pi pi-sitemap text-sm"></i>
          Lihat Visualisasi
        </RouterLink>
      </div>
    </div>

    <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
      <div class="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
        <div class="grid gap-4 md:grid-cols-3 md:gap-6">
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Pencarian</label>
            <input
              v-model="filters.search"
              type="text"
              placeholder="Nama jabatan atau kode"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
            />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Divisi</label>
            <select
              v-model="filters.divisionId"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
            >
              <option value="">Semua Divisi</option>
              <option v-for="division in divisionOptions" :key="division.id" :value="division.id">{{ division.name }}</option>
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

      <div class="mt-6 overflow-hidden rounded-xl border border-gray-100">
        <table class="min-w-full divide-y divide-gray-100 bg-white text-sm">
          <thead class="bg-gray-50 text-xs font-semibold uppercase tracking-wide text-gray-500">
            <tr>
              <th class="px-6 py-3 text-left">Nama Jabatan</th>
              <th class="px-6 py-3 text-left">Divisi</th>
              <th class="px-6 py-3 text-left">Catatan</th>
              <th
                v-if="canUpdateEmployee || canDeleteEmployee"
                class="px-6 py-3 text-right"
              >
                Aksi
              </th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            <tr v-for="position in filteredPositions" :key="position.id" class="hover:bg-gray-50/70">
              <td class="px-6 py-3 font-semibold text-gray-900">{{ position.title }}</td>
              <td class="px-6 py-3 text-gray-700">{{ divisionName(position.divisionId) }}</td>
              <td class="px-6 py-3 text-gray-600">{{ position.notes || '-' }}</td>
              <td
                v-if="canUpdateEmployee || canDeleteEmployee"
                class="px-6 py-3 text-right"
              >
                <div class="flex justify-end gap-2">
                  <button
                    v-if="canUpdateEmployee"
                    class="rounded-full border border-gray-200 px-3 py-1 text-xs text-gray-600 hover:bg-gray-50"
                    @click="openEditModal(position)"
                  >
                    Edit
                  </button>
                  <button
                    v-if="canDeleteEmployee"
                    class="rounded-full border border-red-200 px-3 py-1 text-xs text-red-600 hover:bg-red-50"
                    @click="deletePosition(position.id)"
                  >
                    Hapus
                  </button>
                </div>
              </td>
            </tr>
            <tr v-if="filteredPositions.length === 0">
              <td
                :colspan="(canUpdateEmployee || canDeleteEmployee) ? 4 : 3"
                class="px-6 py-10 text-center text-sm text-gray-500"
              >
                Belum ada data jabatan. Tambahkan data baru untuk melengkapi matriks jabatan.
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
          <h2 class="text-xl font-semibold text-gray-900">{{ isEditing ? 'Edit Jabatan' : 'Tambah Jabatan' }}</h2>
          <p class="text-sm text-gray-500">Isi informasi jabatan sesuai struktur organisasi.</p>
        </div>
          <button class="text-gray-400 hover:text-gray-600" @click="closeModal">
            <i class="pi pi-times text-lg"></i>
          </button>
        </div>

        <form class="mt-6 grid gap-4 md:grid-cols-2" @submit.prevent="savePosition">
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Nama Jabatan <span class="text-red-500">*</span></label>
            <input
              v-model="form.title"
              type="text"
              required
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
            />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Divisi</label>
            <select
              v-model="form.divisionId"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
            >
              <option value="">Pilih divisi</option>
              <option v-for="division in divisionOptions" :key="division.id" :value="division.id">{{ division.name }}</option>
            </select>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700">Induk Jabatan</label>
            <select
              v-model="form.parentId"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
            >
              <option value="">Tanpa induk (posisi puncak)</option>
              <option v-for="option in parentOptions" :key="option.id" :value="option.id">{{ option.title }}</option>
            </select>
            <p class="mt-1 text-xs text-gray-500">Gunakan opsi ini untuk menghubungkan jabatan anak ke induknya.</p>
          </div>
          <div class="md:col-span-2">
            <label class="mb-2 block text-sm font-medium text-gray-700">Catatan / Deskripsi</label>
            <textarea
              v-model="form.notes"
              rows="3"
              class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:border-emerald-500 focus:outline-none focus:ring-2 focus:ring-emerald-100"
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
                ? 'bg-emerald-400 cursor-not-allowed opacity-70'
                : 'bg-emerald-600 hover:bg-emerald-700'
            ]"
            @click="savePosition"
          >
            <i v-if="isSubmitting" class="pi pi-spin pi-spinner"></i>
            <span>{{ isSubmitting ? 'Menyimpan...' : (isEditing ? 'Simpan Perubahan' : 'Simpan Jabatan') }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useHRStore } from '@/stores/hr'
import { useAuthStore } from '@/stores/auth'

const hrStore = useHRStore()
const authStore = useAuthStore()
const toast = useToast()

onMounted(() => {
  hrStore.ensureHydrated().catch((error) => {
    const detail = error?.response?.data?.error || 'Gagal memuat data jabatan.'
    toast.add({ severity: 'error', summary: 'Load gagal', detail, life: 4000 })
  })
})

const filters = ref({
  search: '',
  divisionId: ''
})

const formTemplate = {
  title: '',
  divisionId: '',
  notes: '',
  parentId: ''
}

const form = ref({ ...formTemplate })
const showForm = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const canCreateEmployee = computed(() => authStore.hasPermission('employee.create'))
const canUpdateEmployee = computed(() => authStore.hasPermission('employee.update'))
const canDeleteEmployee = computed(() => authStore.hasPermission('employee.delete'))
const isSubmitting = ref(false)

function divisionName(divisionId) {
  const division = hrStore.findDivision(divisionId)
  return division ? division.name : '-'
}

const divisionOptions = computed(() => hrStore.divisions.map((division) => ({ id: division.id, name: division.name })))

const collectDescendants = (rootId) => {
  const descendants = new Set()
  if (!rootId) return descendants
  const traverse = (currentId) => {
    hrStore.positions
      .filter((position) => position.parentId === currentId)
      .forEach((child) => {
        if (!descendants.has(child.id)) {
          descendants.add(child.id)
          traverse(child.id)
        }
      })
  }
  traverse(rootId)
  return descendants
}

const parentOptions = computed(() => {
  const blocked = new Set()
  if (isEditing.value && editingId.value) {
    blocked.add(editingId.value)
    collectDescendants(editingId.value).forEach((id) => blocked.add(id))
  }
  return hrStore.positions
    .filter((position) => !blocked.has(position.id))
    .map((position) => ({ id: position.id, title: position.title }))
    .sort((a, b) => a.title.localeCompare(b.title))
})

const filteredPositions = computed(() => {
  return hrStore.positions.filter((position) => {
    const matchesSearch = filters.value.search
      ? [position.title, position.code, position.notes]
          .filter(Boolean)
          .some((field) => field.toLowerCase().includes(filters.value.search.toLowerCase()))
      : true

    const matchesDivision = filters.value.divisionId ? position.divisionId === filters.value.divisionId : true

    return matchesSearch && matchesDivision
  })
})

const resetFilters = () => {
  filters.value = { search: '', divisionId: '' }
}

const openCreateModal = () => {
  if (!canCreateEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk menambah jabatan.', life: 3500 })
    return
  }
  isEditing.value = false
  editingId.value = null
  form.value = { ...formTemplate }
  showForm.value = true
}

const openEditModal = (position) => {
  if (!canUpdateEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk memperbarui jabatan.', life: 3500 })
    return
  }
  isEditing.value = true
  editingId.value = position.id
  form.value = {
    title: position.title,
    divisionId: position.divisionId || '',
    notes: position.notes || '',
    parentId: position.parentId || ''
  }
  showForm.value = true
}

const closeModal = () => {
  showForm.value = false
  form.value = { ...formTemplate }
  editingId.value = null
  isSubmitting.value = false
}

const wouldCreateCircularHierarchy = (positionId, targetParentId) => {
  if (!positionId || !targetParentId) return false
  if (positionId === targetParentId) return true
  return collectDescendants(positionId).has(targetParentId)
}

const savePosition = async () => {
  if (isEditing.value) {
    if (!canUpdateEmployee.value) {
      toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk memperbarui jabatan.', life: 3500 })
      return
    }
  } else if (!canCreateEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk menambah jabatan.', life: 3500 })
    return
  }

  if (!form.value.title?.trim()) {
    toast.add({ severity: 'warn', summary: 'Validasi', detail: 'Nama jabatan wajib diisi.', life: 3500 })
    return
  }

  const payload = {
    title: form.value.title.trim(),
    divisionId: form.value.divisionId || null,
    notes: (form.value.notes || '').trim(),
    parentId: form.value.parentId || null
  }

  if (isEditing.value && editingId.value && wouldCreateCircularHierarchy(editingId.value, payload.parentId)) {
    toast.add({ severity: 'warn', summary: 'Validasi', detail: 'Induk yang dipilih membuat hierarki berulang. Silakan pilih jabatan lain.', life: 4000 })
    return
  }

  isSubmitting.value = true
  try {
    if (isEditing.value && editingId.value) {
      await hrStore.updatePosition(editingId.value, payload)
      toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Jabatan diperbarui.', life: 3000 })
    } else {
      await hrStore.createPosition(payload)
      toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Jabatan berhasil ditambahkan.', life: 3000 })
    }
    closeModal()
  } catch (error) {
    const detail = error?.response?.data?.error || 'Gagal menyimpan data jabatan.'
    toast.add({ severity: 'error', summary: 'Gagal', detail, life: 4000 })
  } finally {
    isSubmitting.value = false
  }
}

const deletePosition = async (id) => {
  if (!canDeleteEmployee.value) {
    toast.add({ severity: 'warn', summary: 'Tidak diizinkan', detail: 'Anda tidak memiliki izin untuk menghapus jabatan.', life: 3500 })
    return
  }
  if (!confirm('Hapus data jabatan ini?')) {
    return
  }
  try {
    await hrStore.deletePosition(id)
    toast.add({ severity: 'success', summary: 'Berhasil', detail: 'Jabatan telah dihapus.', life: 3000 })
  } catch (error) {
    const detail = error?.response?.data?.error || 'Gagal menghapus data jabatan.'
    toast.add({ severity: 'error', summary: 'Gagal', detail, life: 4000 })
  }
}
</script>
