<template>
  <div>
    <div class="flex justify-between items-center mb-6">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Warehouses</h1>
        <p class="text-gray-600 mt-1">Manage warehouse locations and information</p>
      </div>
      <Button
        v-if="authStore.isAdmin || authStore.isManager"
        label="Add Warehouse"
        icon="pi pi-plus-circle"
        raised
        @click="openDialog()"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-sm md:text-base border-none shadow-md rounded-lg"
      />
    </div>

    <div v-if="loading" class="text-center py-12">
      <ProgressSpinner />
    </div>

    <div v-else-if="warehouses.length === 0" class="text-center py-12 bg-white rounded-lg">
      <i class="pi pi-building text-6xl text-gray-300 mb-4"></i>
      <h3 class="text-lg font-medium text-gray-800 mb-2">No Warehouses Yet</h3>
      <p class="text-gray-600 mb-4">Start by creating your first warehouse</p>
      <Button
        v-if="authStore.isAdmin || authStore.isManager"
        label="Add Warehouse"
        icon="pi pi-plus-circle"
        raised
        @click="openDialog()"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-sm md:text-base border-none shadow-md rounded-lg"
      />
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div
        v-for="(warehouse, index) in warehouses"
        :key="warehouse.id"
        class="rounded-lg shadow-sm border hover:shadow-md transition-all overflow-hidden"
        :style="{
          backgroundColor: colorWithAlpha(getWarehouseColor(warehouse, index), 0.08),
          borderColor: colorWithAlpha(getWarehouseColor(warehouse, index), 0.35)
        }"
      >
        <div
          class="h-1"
          :style="{ backgroundColor: getWarehouseColor(warehouse, index) }"
        ></div>
        <div class="p-6">
          <div class="flex items-start justify-between mb-4">
            <div class="flex items-start gap-3 flex-1">
              <div
                class="w-12 h-12 rounded-full flex items-center justify-center text-white text-xl flex-shrink-0"
                :style="{ backgroundColor: getWarehouseColor(warehouse, index) }"
              >
                <i class="pi pi-building"></i>
              </div>
              <div>
                <h3 class="text-lg font-semibold text-gray-800">{{ warehouse.name }}</h3>
                <p class="text-sm text-gray-500 font-mono">{{ warehouse.code }}</p>
              </div>
            </div>
            <Tag
              :value="warehouse.is_active ? 'Active' : 'Inactive'"
              :severity="warehouse.is_active ? 'success' : 'danger'"
            />
          </div>

          <div class="space-y-2 text-sm mb-4">
            <div class="flex items-center text-gray-600">
              <i class="pi pi-map-marker mr-2 text-blue-500"></i>
              <span>{{ warehouse.city }}, {{ warehouse.province }}</span>
            </div>
            <div class="flex items-center text-gray-600">
              <i class="pi pi-phone mr-2 text-green-500"></i>
              <span>{{ warehouse.phone || '-' }}</span>
            </div>
            <div v-if="warehouse.manager" class="flex items-center text-gray-600">
              <i class="pi pi-user mr-2 text-purple-500"></i>
              <span>{{ warehouse.manager.full_name }}</span>
            </div>
          </div>

          <Divider />

          <div class="flex gap-2 pt-2">
            <Button
              label="View Detail"
              icon="pi pi-eye"
              @click="$router.push(`/warehouses/${warehouse.id}`)"
              class="flex-1 p-button-outlined"
              size="small"
            />
            <Button
              v-if="authStore.isAdmin || authStore.isManager"
              icon="pi pi-pencil"
              @click="openDialog(warehouse)"
              class="p-button-outlined p-button-warning"
              size="small"
              v-tooltip.top="'Edit'"
            />
            <Button
              v-if="authStore.isAdmin"
              icon="pi pi-trash"
              @click="confirmDelete(warehouse)"
              class="p-button-outlined p-button-danger"
              size="small"
              v-tooltip.top="'Delete'"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Add/Edit Dialog -->
    <Dialog
      v-model:visible="dialogVisible"
      :modal="true"
      :style="{ width: '800px' }"
      :dismissableMask="true"
    >
      <template #header>
        <div class="flex items-center gap-3">
          <div class="w-12 h-12 bg-blue-100 rounded-lg flex items-center justify-center">
            <i :class="editMode ? 'pi pi-pencil' : 'pi pi-plus'" class="text-blue-600 text-xl"></i>
          </div>
          <div>
            <h3 class="text-xl font-bold text-gray-800">
              {{ editMode ? 'Edit Warehouse' : 'Add New Warehouse' }}
            </h3>
            <p class="text-sm text-gray-500">
              {{ editMode ? 'Update warehouse information' : 'Fill in warehouse details' }}
            </p>
          </div>
        </div>
      </template>

      <div class="space-y-6 py-4">
        <!-- Basic Information -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="flex items-center gap-2 text-sm font-semibold text-gray-700 mb-2">
              <i class="pi pi-building text-gray-500"></i>
              Warehouse Name <span class="text-red-500">*</span>
            </label>
            <InputText
              v-model="form.name"
              placeholder="e.g., Main Warehouse Jakarta"
              class="w-full"
            />
          </div>

          <div>
            <label class="flex items-center gap-2 text-sm font-semibold text-gray-700 mb-2">
              <i class="pi pi-qrcode text-gray-500"></i>
              Warehouse Code <span class="text-red-500">*</span>
            </label>
            <InputText
              v-model="form.code"
              placeholder="e.g., WH-JKT-01"
              class="w-full"
            />
          </div>
        </div>

        <div>
          <label class="flex items-center gap-2 text-sm font-semibold text-gray-700 mb-2">
            <i class="pi pi-palette text-gray-500"></i>
            Warehouse Color
          </label>
          <div class="space-y-3">
            <div class="flex flex-wrap items-center gap-3">
              <input
                type="color"
                v-model="form.color"
                @change="form.color = normalizeHexColor(form.color) || form.color"
                class="w-12 h-12 border border-gray-300 rounded cursor-pointer"
              />
              <input
                v-model="form.color"
                type="text"
                maxlength="7"
                @blur="form.color = normalizeHexColor(form.color) || form.color"
                class="w-32 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono uppercase"
                placeholder="#FFFFFF"
              />
              <span class="text-xs text-gray-500">
                Use custom color or pick from presets below
              </span>
            </div>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="color in colorPalette"
                :key="color"
                type="button"
                class="w-10 h-10 rounded-lg border-2 transition-transform"
                :class="normalizeHexColor(form.color) === normalizeHexColor(color) ? 'border-gray-900 scale-105' : 'border-transparent hover:scale-105'"
                :style="{ backgroundColor: color }"
                @click="form.color = normalizeHexColor(color) || color"
              />
            </div>
          </div>
        </div>

        <!-- Address -->
        <div>
          <label class="flex items-center gap-2 text-sm font-semibold text-gray-700 mb-2">
            <i class="pi pi-map-marker text-gray-500"></i>
            Address <span class="text-red-500">*</span>
          </label>
          <Textarea
            v-model="form.address"
            rows="3"
            placeholder="Complete warehouse address"
            class="w-full"
          />
        </div>

        <!-- Location Details -->
        <div class="grid grid-cols-3 gap-4">
          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">
              City <span class="text-red-500">*</span>
            </label>
            <InputText
              v-model="form.city"
              placeholder="City"
              class="w-full"
            />
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">
              Province <span class="text-red-500">*</span>
            </label>
            <InputText
              v-model="form.province"
              placeholder="Province"
              class="w-full"
            />
          </div>

          <div>
            <label class="block text-sm font-semibold text-gray-700 mb-2">
              Postal Code
            </label>
            <InputText
              v-model="form.postal_code"
              placeholder="Postal Code"
              class="w-full"
            />
          </div>
        </div>

        <!-- Contact Information -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="flex items-center gap-2 text-sm font-semibold text-gray-700 mb-2">
              <i class="pi pi-phone text-gray-500"></i>
              Phone Number <span class="text-red-500">*</span>
            </label>
            <InputText
              v-model="form.phone"
              placeholder="e.g., +62 21 1234567"
              class="w-full"
            />
          </div>

          <div>
            <label class="flex items-center gap-2 text-sm font-semibold text-gray-700 mb-2">
              <i class="pi pi-envelope text-gray-500"></i>
              Email
            </label>
            <InputText
              v-model="form.email"
              type="email"
              placeholder="warehouse@company.com"
              class="w-full"
            />
          </div>
        </div>

        <!-- Manager & Status -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="flex items-center gap-2 text-sm font-semibold text-gray-700 mb-2">
              <i class="pi pi-user text-gray-500"></i>
              Warehouse Manager
            </label>
            <Dropdown
              v-model="form.manager_id"
              :options="managers"
              optionLabel="full_name"
              optionValue="id"
              placeholder="Select manager"
              class="w-full"
              :filter="true"
              showClear
            >
              <template #option="{ option }">
                <div class="flex items-center gap-2">
                  <i class="pi pi-user text-gray-400"></i>
                  <div>
                    <div class="font-medium">{{ option.full_name }}</div>
                    <div class="text-xs text-gray-500">{{ option.email }}</div>
                  </div>
                </div>
              </template>
            </Dropdown>
          </div>

          <div class="flex items-end">
            <div class="flex items-center gap-2 mb-2">
              <Checkbox v-model="form.is_active" inputId="is_active" binary />
              <label for="is_active" class="text-sm font-semibold text-gray-700">
                Active Warehouse
              </label>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <div class="flex justify-between items-center">
          <div class="text-sm text-gray-500">
            <i class="pi pi-info-circle mr-1"></i>
            Fields marked with <span class="text-red-500">*</span> are required
          </div>
          <div class="flex gap-3">
            <Button
              label="Cancel"
              icon="pi pi-times"
              @click="dialogVisible = false"
              class="p-button-outlined"
            />
            <Button
              :label="editMode ? 'Update' : 'Create'"
              :icon="editMode ? 'pi pi-check' : 'pi pi-plus'"
              @click="saveWarehouse"
              :loading="saving"
              class="px-6"
            />
          </div>
        </div>
      </template>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <Dialog
      v-model:visible="deleteDialogVisible"
      header="Confirm Delete"
      :modal="true"
      :style="{ width: '450px' }"
    >
      <div class="flex items-start gap-4">
        <i class="pi pi-exclamation-triangle text-3xl text-orange-500"></i>
        <div>
          <p class="text-gray-700 mb-2">
            Are you sure you want to delete warehouse <strong>{{ warehouseToDelete?.name }}</strong>?
          </p>
          <p class="text-sm text-gray-500">
            This action cannot be undone. All inventory items in this warehouse will need to be relocated.
          </p>
        </div>
      </div>

      <template #footer>
        <Button
          label="Cancel"
          icon="pi pi-times"
          @click="deleteDialogVisible = false"
          class="p-button-outlined"
        />
        <Button
          label="Delete"
          icon="pi pi-trash"
          @click="deleteWarehouse"
          class="p-button-danger"
          :loading="deleting"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useWarehouseStore } from '@/stores/warehouse'
import { useAuthStore } from '@/stores/auth'
import { useToast } from 'primevue/usetoast'
import axios from '@/api/axios'

const router = useRouter()
const warehouseStore = useWarehouseStore()
const authStore = useAuthStore()
const toast = useToast()

const warehouses = ref([])
const managers = ref([])
const loading = ref(false)
const saving = ref(false)
const deleting = ref(false)
const dialogVisible = ref(false)
const deleteDialogVisible = ref(false)
const editMode = ref(false)
const warehouseToDelete = ref(null)

const defaultColor = '#6366F1'

const form = ref({
  code: '',
  name: '',
  address: '',
  city: '',
  province: '',
  postal_code: '',
  phone: '',
  email: '',
  manager_id: null,
  color: defaultColor,
  is_active: true
})

const colorPalette = [
  '#6366F1', // indigo
  '#F97316', // orange
  '#14B8A6', // teal
  '#F43F5E', // rose
  '#8B5CF6', // violet
  '#10B981', // emerald
  '#F59E0B', // amber
  '#3B82F6', // blue
  '#EC4899', // pink
  '#0EA5E9', // sky
  '#22D3EE', // cyan
  '#A855F7'  // purple
]

const HEX_COLOR_REGEX = /^#(?:[0-9A-Fa-f]{6}|[0-9A-Fa-f]{3})$/

function normalizeHexColor(color) {
  if (typeof color !== 'string') {
    return ''
  }
  let value = color.trim()
  if (!value) {
    return ''
  }
  if (!value.startsWith('#')) {
    value = `#${value}`
  }
  if (!HEX_COLOR_REGEX.test(value)) {
    return ''
  }
  if (value.length === 4) {
    value = `#${value[1]}${value[1]}${value[2]}${value[2]}${value[3]}${value[3]}`
  }
  return value.toUpperCase()
}

function paletteColor(index) {
  if (!colorPalette.length) {
    return defaultColor
  }
  return colorPalette[index % colorPalette.length]
}

function getWarehouseColor(warehouse, index) {
  const candidate = typeof warehouse === 'string' ? warehouse : warehouse?.color
  const normalized = normalizeHexColor(candidate)
  if (normalized) {
    return normalized
  }
  return normalizeHexColor(paletteColor(index)) || defaultColor
}

function colorWithAlpha(hexColor, alpha) {
  const normalized = normalizeHexColor(hexColor) || normalizeHexColor(defaultColor) || defaultColor
  const hex = normalized.replace('#', '')

  const r = parseInt(hex.slice(0, 2), 16)
  const g = parseInt(hex.slice(2, 4), 16)
  const b = parseInt(hex.slice(4, 6), 16)

  return `rgba(${r}, ${g}, ${b}, ${alpha})`
}

onMounted(async () => {
  await fetchWarehouses()
  if (authStore.isAdmin || authStore.isManager) {
    await fetchManagers()
  }
})

async function fetchWarehouses() {
  loading.value = true
  try {
    await warehouseStore.fetchWarehouses()
    warehouses.value = warehouseStore.warehouses
  } catch (error) {
    console.error('Failed to fetch warehouses:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load warehouses',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

async function fetchManagers() {
  if (!(authStore.isAdmin || authStore.isManager)) {
    managers.value = []
    return
  }
  try {
    // Fetch users with Manager or Admin role
    const response = await axios.get('/settings/users?role=Manager,Admin')
    managers.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch managers:', error)
    // If endpoint doesn't exist, just set empty array
    managers.value = []
    if (error.response?.status === 403) {
      toast.add({
        severity: 'warn',
        summary: 'Access limited',
        detail: 'You do not have permission to view manager list.',
        life: 3000
      })
    }
  }
}

function openDialog(warehouse = null) {
  if (warehouse) {
    editMode.value = true
    form.value = {
      id: warehouse.id,
      code: warehouse.code,
      name: warehouse.name,
      address: warehouse.address,
      city: warehouse.city,
      province: warehouse.province,
      postal_code: warehouse.postal_code || '',
      phone: warehouse.phone,
      email: warehouse.email || '',
      manager_id: warehouse.manager_id,
      color: normalizeHexColor(warehouse.color) || defaultColor,
      is_active: warehouse.is_active
    }
  } else {
    editMode.value = false
    form.value = {
      code: '',
      name: '',
      address: '',
      city: '',
      province: '',
      postal_code: '',
      phone: '',
      email: '',
      manager_id: null,
      color: defaultColor,
      is_active: true
    }
  }
  dialogVisible.value = true
}

async function saveWarehouse() {
  // Validation
  if (!form.value.code || !form.value.name || !form.value.address || 
      !form.value.city || !form.value.province || !form.value.phone) {
    toast.add({
      severity: 'warn',
      summary: 'Validation Error',
      detail: 'Please fill in all required fields',
      life: 3000
    })
    return
  }

  saving.value = true
  const normalizedColor = normalizeHexColor(form.value.color) || defaultColor
  const payload = {
    ...form.value,
    color: normalizedColor
  }
  delete payload.id

  try {
    if (editMode.value) {
      await axios.put(`/warehouses/${form.value.id}`, payload)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Warehouse updated successfully',
        life: 3000
      })
    } else {
      await axios.post('/warehouses', payload)
      toast.add({
        severity: 'success',
        summary: 'Success',
        detail: 'Warehouse created successfully',
        life: 3000
      })
    }
    form.value.color = normalizedColor
    
    dialogVisible.value = false
    await fetchWarehouses()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.error || 'Failed to save warehouse',
      life: 3000
    })
  } finally {
    saving.value = false
  }
}

function confirmDelete(warehouse) {
  warehouseToDelete.value = warehouse
  deleteDialogVisible.value = true
}

async function deleteWarehouse() {
  deleting.value = true
  try {
    await axios.delete(`/warehouses/${warehouseToDelete.value.id}`)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Warehouse deleted successfully',
      life: 3000
    })
    deleteDialogVisible.value = false
    await fetchWarehouses()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.error || 'Failed to delete warehouse',
      life: 3000
    })
  } finally {
    deleting.value = false
  }
}
</script>
