<template>
  <div class="p-4">
    <div v-if="canViewInventory">
      <!-- Action toolbar -->
    <div class="flex items-center justify-between mb-4">
      <div>
        <h2 class="text-xl font-bold">Inventory Items</h2>
      </div>
      <div class="flex items-center space-x-2">
        <button @click="downloadTemplate" class="px-3 py-1 bg-gray-200 rounded">Template CSV</button>
        <button
          v-if="canCreateInventory"
          @click="triggerImport"
          class="px-3 py-1 bg-gray-200 rounded"
        >
          Import CSV
        </button>
        <button @click="exportCSV" class="px-3 py-1 bg-gray-200 rounded">Export CSV</button>
        <button @click="exportPDF" class="px-3 py-1 bg-gray-200 rounded">Export PDF</button>
        <Button
          v-if="canCreateInventory"
          label="Add Item"
          icon="pi pi-plus-circle"
          raised
          @click="openAddItem"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-sm md:text-base border-none shadow-md rounded-lg"
        />
        <button
          v-if="canDeleteInventory"
          @click="showBatchConfirm = true"
          :disabled="isDeleteDisabled"
          :class="deleteButtonClasses"
        >
          Delete Selected
        </button>
        <input type="file" ref="importInput" class="hidden" accept=".csv" @change="handleImportChange" />
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white p-4 rounded-lg shadow-sm mb-4">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Search</label>
          <input
            v-model="filters.search"
            placeholder="Search by name or SN..."
            class="w-full border p-2 rounded"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Warehouse</label>
          <select
            v-model="filters.warehouse_id"
            class="w-full border p-2 rounded"
          >
            <option value="">All Warehouses</option>
            <option v-for="wh in warehouses" :key="wh.id" :value="wh.id">{{ wh.name }}</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Category</label>
          <select
            v-model="filters.category"
            class="w-full border p-2 rounded"
          >
            <option value="">All Categories</option>
            <option v-for="cat in categories" :key="cat.id || cat.name" :value="cat.name">{{ cat.name }}</option>
          </select>
        </div>
        <div class="flex items-end">
          <button @click="clearFilters" class="px-4 py-2 border rounded">Clear Filters</button>
        </div>
      </div>
    </div>

    <!-- Data Table -->
    <fwb-table hoverable>
      <fwb-table-head>
        <fwb-table-head-cell>
          <input
            type="checkbox"
            :checked="allVisibleSelected"
            @change="toggleSelectAll($event.target.checked)"
            :disabled="!canDeleteInventory"
          />
        </fwb-table-head-cell>
        <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('name')">
          Name
          <i :class="getSortIcon('name')" class="ml-1"></i>
        </fwb-table-head-cell>
        <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('sn')">
          SN
          <i :class="getSortIcon('sn')" class="ml-1"></i>
        </fwb-table-head-cell>
        <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('category')">
          Category
          <i :class="getSortIcon('category')" class="ml-1"></i>
        </fwb-table-head-cell>
        <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('warehouse_name')">
          Warehouse
          <i :class="getSortIcon('warehouse_name')" class="ml-1"></i>
        </fwb-table-head-cell>
        <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('quantity')">
          Quantity
          <i :class="getSortIcon('quantity')" class="ml-1"></i>
        </fwb-table-head-cell>
        <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('min_stock')">
          Min Stock
          <i :class="getSortIcon('min_stock')" class="ml-1"></i>
        </fwb-table-head-cell>
        <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('is_active')">
          Status
          <i :class="getSortIcon('is_active')" class="ml-1"></i>
        </fwb-table-head-cell>
        <fwb-table-head-cell>
          <span class="sr-only">Actions</span>
        </fwb-table-head-cell>
      </fwb-table-head>
      <fwb-table-body>
        <fwb-table-row
          v-for="row in paginatedRows"
          :key="row.key"
          class="align-middle"
          :class="row.type === 'group' ? 'bg-gray-50' : ''"
        >
          <fwb-table-cell>
            <template v-if="row.type === 'group'">
              <span class="text-gray-400">—</span>
            </template>
            <template v-else>
              <input
                type="checkbox"
                :checked="isSelected(row.item)"
                @change="toggleSelection(row.item, $event.target.checked)"
                :disabled="!canDeleteInventory"
              />
            </template>
          </fwb-table-cell>
          <fwb-table-cell>
            <template v-if="row.type === 'group'">
              <div class="flex items-center gap-2">
                <button
                  type="button"
                  class="w-6 h-6 flex items-center justify-center border border-gray-300 rounded text-sm font-semibold"
                  @click="toggleGroup(row.group.key)"
                >
                  {{ isGroupExpanded(row.group.key) ? '−' : '+' }}
                </button>
                <div>
                  <div class="font-semibold text-gray-800">{{ row.group.displayName }}</div>
                  <div class="text-xs text-gray-500">{{ row.group.items.length }} item{{ row.group.items.length > 1 ? 's' : '' }}</div>
                </div>
              </div>
            </template>
            <template v-else>
              <div class="text-sm ml-8">
                <div class="font-medium text-gray-800">{{ row.item.name || '-' }}</div>
              </div>
            </template>
          </fwb-table-cell>
          <fwb-table-cell>
            <template v-if="row.type === 'group'">
              <span class="text-gray-400">—</span>
            </template>
            <template v-else>
              <div class="text-sm">
                <div class="font-medium text-gray-800">{{ row.item.sn || '-' }}</div>
              </div>
            </template>
          </fwb-table-cell>
          <fwb-table-cell>
            <template v-if="row.type === 'group'">
              <span class="text-gray-400 text-sm">—</span>
            </template>
            <template v-else>
              <Tag :value="row.item.category || '-'" :style="getCategoryTagStyle(row.item.category)" />
            </template>
          </fwb-table-cell>
          <fwb-table-cell>
            <template v-if="row.type === 'group'">
              <span class="text-gray-400 text-sm">—</span>
            </template>
            <template v-else>
              <div class="text-sm">
                <span class="text-gray-600">{{ row.item.warehouse?.name || row.item.Warehouse?.name || '-' }}</span>
              </div>
            </template>
          </fwb-table-cell>
          <fwb-table-cell>
            <template v-if="row.type === 'group'">
              <span class="text-sm text-gray-500">—</span>
            </template>
            <template v-else>
              <span
                class="font-semibold"
                :class="row.item.quantity <= row.item.min_stock ? 'text-red-600' : 'text-green-600'"
              >
                {{ row.item.quantity != null ? row.item.quantity + ' ' + (row.item.unit || 'pcs') : '-' }}
              </span>
            </template>
          </fwb-table-cell>
          <fwb-table-cell>
            <template v-if="row.type === 'group'">
              <span class="text-sm text-gray-500">—</span>
            </template>
            <template v-else>
              <span class="text-sm text-gray-600">
                {{ row.item.min_stock != null ? row.item.min_stock + ' ' + (row.item.unit || 'pcs') : '-' }}
              </span>
            </template>
          </fwb-table-cell>
          <fwb-table-cell>
            <template v-if="row.type === 'group'">
              <span class="text-gray-400 text-sm">—</span>
            </template>
            <template v-else>
              <span
                class="px-2 py-1 rounded-full text-xs font-semibold select-none"
                :class="[
                  row.item.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800',
                  canUpdateInventory ? 'cursor-pointer' : 'cursor-not-allowed opacity-70'
                ]"
                @click="toggleItemStatus(row.item)"
                :title="canUpdateInventory ? (row.item.is_active ? 'Klik untuk nonaktifkan' : 'Klik untuk aktifkan') : 'Tidak memiliki izin untuk mengubah status item'"
              >
                {{ row.item.is_active ? 'Aktif' : 'Tidak Aktif' }}
              </span>
            </template>
          </fwb-table-cell>
          <fwb-table-cell>
            <template v-if="row.type === 'group'">
              <span class="text-gray-400 text-sm">—</span>
            </template>
            <template v-else>
              <div class="flex items-center gap-2">
                <Button icon="pi pi-eye" class="p-button-text" @click="viewItem(row.item)" />
                <Button
                  v-if="canUpdateInventory"
                  icon="pi pi-pencil"
                  class="p-button-text p-button-warning"
                  @click="editItem(row.item)"
                />
                <Button
                  v-if="canDeleteInventory"
                  icon="pi pi-trash"
                  class="p-button-text p-button-danger"
                  @click="confirmDelete(row.item)"
                />
              </div>
            </template>
          </fwb-table-cell>
        </fwb-table-row>
      </fwb-table-body>
    </fwb-table>

    <!-- Pagination Controls -->
    <div class="flex items-center justify-between mt-4">
      <div class="flex items-center space-x-2">
        <span>Rows per page:</span>
        <select v-model.number="rowsPerPage" class="border p-1 rounded">
          <option v-for="opt in rowsPerPageOptions" :key="opt" :value="opt">{{ opt }}</option>
        </select>
      </div>
      <div class="flex items-center space-x-2">
        <button @click="changePage(currentPage - 1)" :disabled="currentPage === 1" class="px-2 py-1 border rounded">Prev</button>
        <span>Page <strong>{{ currentPage }}</strong> of <strong>{{ totalPages }}</strong></span>
        <button @click="changePage(currentPage + 1)" :disabled="currentPage === totalPages" class="px-2 py-1 border rounded">Next</button>
      </div>
    </div>

    <!-- Batch Delete Confirmation Modal -->
    <div v-if="showBatchConfirm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded shadow-lg w-80">
        <h3 class="text-lg font-semibold mb-4">Confirm Delete</h3>
        <p>Are you sure you want to delete {{ selectedItems.length }} items?</p>
        <div class="mt-6 flex justify-end space-x-2">
          <button @click="showBatchConfirm = false" class="px-4 py-2 rounded border">Cancel</button>
          <button @click="confirmBatchDelete" class="px-4 py-2 bg-red-600 text-white rounded">Delete</button>
        </div>
      </div>
    </div>

    <!-- Single Delete Confirmation Modal -->
    <div v-if="showDeleteConfirm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded shadow-lg w-80">
        <h3 class="text-lg font-semibold mb-4">Confirm Delete</h3>
        <p>Are you sure you want to delete this item?</p>
        <div class="text-sm text-gray-600 mt-2">
          <strong>{{ itemToDelete?.name }}</strong><br>
          SN: {{ itemToDelete?.sn }}
        </div>
        <div class="mt-6 flex justify-end space-x-2">
          <button @click="showDeleteConfirm = false; itemToDelete = null" class="px-4 py-2 rounded border">Cancel</button>
          <button @click="deleteItem" class="px-4 py-2 bg-red-600 text-white rounded">Delete</button>
        </div>
      </div>
    </div>

    <!-- Add/Edit Item Modal with PrimeVue Dialog -->
    <Dialog v-model:visible="showEditModal" modal :header="modalTitle" :style="{ width: '30rem' }">
      <span class="text-surface-500 dark:text-surface-400 block mb-8">
        {{ isEditing ? 'Update item information.' : 'Add a new item.' }}
      </span>
      <div class="flex flex-col gap-4 mb-4">
        <div class="flex items-center gap-4">
          <label for="snType" class="font-semibold w-32">Tipe SN</label>
          <div class="flex gap-4">
            <label class="flex items-center gap-2">
              <input type="radio" v-model="editForm.has_sn" :value="true" />
              <span>Punya SN</span>
            </label>
            <label class="flex items-center gap-2">
              <input type="radio" v-model="editForm.has_sn" :value="false" />
              <span>Tanpa SN</span>
            </label>
          </div>
        </div>
        <div v-if="editForm.has_sn" class="flex items-center gap-4">
          <label for="sn" class="font-semibold w-32">Serial Number</label>
          <InputText id="sn" v-model="editForm.sn" class="flex-auto" autocomplete="off" placeholder="Enter SN" />
        </div>
        <div class="flex items-center gap-4">
          <label for="name" class="font-semibold w-32">Name</label>
          <InputText id="name" v-model="editForm.name" class="flex-auto" autocomplete="off" placeholder="Enter name" />
        </div>
        <div class="flex items-center gap-4">
          <label for="category" class="font-semibold w-32">Category</label>
          <select id="category" v-model="editForm.category" class="flex-auto border p-2 rounded">
            <option value="">Select Category</option>
            <option v-for="cat in categories" :key="cat.id || cat.name" :value="cat.name">{{ cat.name }}</option>
          </select>
        </div>
        <div class="flex items-center gap-4">
          <label for="warehouse" class="font-semibold w-32">Warehouse</label>
          <select id="warehouse" v-model="editForm.warehouse_id" class="flex-auto border p-2 rounded">
            <option value="">Select Warehouse</option>
            <option v-for="wh in warehouses" :key="wh.id" :value="wh.id">{{ wh.name }}</option>
          </select>
        </div>
        <div class="flex items-center gap-4">
          <label for="quantity" class="font-semibold w-32">Quantity</label>
          <InputText id="quantity" v-model.number="editForm.quantity" class="flex-auto" type="number" autocomplete="off" placeholder="Enter quantity" />
        </div>
        <div class="flex items-center gap-4">
          <label for="min_stock" class="font-semibold w-32">Minimum Stock</label>
          <InputText id="min_stock" v-model.number="editForm.min_stock" class="flex-auto" type="number" min="0" autocomplete="off" placeholder="Enter minimum stock threshold" />
        </div>
        <div class="flex items-center gap-4">
          <label for="unit" class="font-semibold w-32">Unit</label>
          <InputText id="unit" v-model="editForm.unit" class="flex-auto" autocomplete="off" placeholder="e.g. pcs, meter, box" />
        </div>
        <div class="flex items-center gap-4">
          <label for="unit_price" class="font-semibold w-32">Unit Price</label>
          <InputText id="unit_price" v-model.number="editForm.unit_price" class="flex-auto" type="number" autocomplete="off" placeholder="Enter price" />
        </div>
        <div class="flex items-center gap-4">
          <label for="description" class="font-semibold w-32">Description</label>
          <InputText id="description" v-model="editForm.description" class="flex-auto" autocomplete="off" placeholder="Enter description" />
        </div>
        <div class="flex items-center gap-4">
          <label for="is_active" class="font-semibold w-32">Status</label>
          <select id="is_active" v-model="editForm.is_active" class="flex-auto border p-2 rounded">
            <option :value="true">Aktif</option>
            <option :value="false">Tidak Aktif</option>
          </select>
        </div>
      </div>
      <div class="flex justify-end gap-2 mt-6">
        <Button 
          type="button"
          label="Cancel"
          severity="danger"
          icon="pi pi-times"
          class="font-bold bg-red-500 border-none text-white hover:bg-red-600 focus:bg-red-700"
          @click="closeEditModal"
        />
        <Button 
          type="button"
          :label="submitButtonLabel"
          :disabled="saving"
          icon="pi pi-check"
          class="font-bold bg-blue-600 border-none text-white hover:bg-blue-700 focus:bg-blue-800"
          @click="submitItemForm"
        />
      </div>
    </Dialog>
    </div>
    <div v-else class="bg-white rounded-lg shadow p-6 text-center text-gray-600">
      <p class="text-lg font-semibold text-gray-800 mb-2">Akses Terbatas</p>
      <p class="text-sm">Anda tidak memiliki izin untuk melihat data inventory.</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '@/api/axios'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '@/stores/auth'
import {
  FwbTable,
  FwbTableHead,
  FwbTableHeadCell,
  FwbTableBody,
  FwbTableRow,
  FwbTableCell,
  FwbA
} from 'flowbite-vue'
import Tag from 'primevue/tag'
import Button from 'primevue/button'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()
const items = ref([])
const canViewInventory = computed(() => authStore.hasPermission('inventory.view'))
const canCreateInventory = computed(() => authStore.hasPermission('inventory.create'))
const canUpdateInventory = computed(() => authStore.hasPermission('inventory.update'))
const canDeleteInventory = computed(() => authStore.hasPermission('inventory.delete'))
const selectedItems = ref([])
const isDeleteDisabled = computed(() => selectedItems.value.length === 0 || !canDeleteInventory.value)
const deleteButtonClasses = computed(() => [
  'px-3',
  'py-1',
  'rounded',
  'transition-colors',
  isDeleteDisabled.value
    ? 'bg-gray-300 text-gray-500 cursor-not-allowed'
    : 'bg-red-600 hover:bg-red-700 text-white'
])
const isEditing = ref(false)
const ensurePermission = (permission, detail) => {
  if (authStore.hasPermission(permission)) {
    return true
  }
  toast.add({
    severity: 'warn',
    summary: 'Akses ditolak',
    detail: detail || 'Anda tidak memiliki izin untuk melakukan aksi ini',
    life: 3000
  })
  return false
}
const modalTitle = computed(() => (isEditing.value ? 'Edit Item' : 'Add Item'))
const importInput = ref(null)
const showBatchConfirm = ref(false)
const showDeleteConfirm = ref(false)
const itemToDelete = ref(null)
const showEditModal = ref(false)
const defaultItemForm = () => ({
  id: null,
  has_sn: true,
  sn: '',
  name: '',
  category: '',
  warehouse_id: '',
  quantity: 0,
  min_stock: 0,
  unit: 'pcs',
  unit_price: 0,
  description: '',
  is_active: true
})
const editForm = ref({
  has_sn: true,
  sn: '',
  name: '',
  category: '',
  warehouse_id: '',
  quantity: 0,
  min_stock: 0,
  unit: 'pcs',
  unit_price: 0,
  description: '',
  is_active: true
})
const saving = ref(false)
const submitButtonLabel = computed(() => {
  if (saving.value) return isEditing.value ? 'Saving...' : 'Creating...'
  return isEditing.value ? 'Save' : 'Create'
})
const warehouses = ref([])
const categories = ref([])
const categoryColorMap = computed(() => {
  return categories.value.reduce((acc, cat) => {
    if (!cat) return acc
    if (typeof cat === 'string') {
      acc[cat] = null
    } else if (cat.name) {
      acc[cat.name] = cat.color || null
    }
    return acc
  }, {})
})
const ensureHexColor = (value) => {
  if (!value) return ''
  return value.startsWith('#') ? value : `#${value}`
}
const getCategoryTagStyle = (categoryName) => {
  const color = ensureHexColor(categoryColorMap.value[categoryName])
  if (!color) return {}
  return {
    backgroundColor: `${color}20`,
    color,
    border: `1px solid ${color}33`
  }
}
const FILTER_STORAGE_KEY = 'inventory_filters'

const filters = ref(loadSavedFilters())

function loadSavedFilters() {
  try {
    const raw = localStorage.getItem(FILTER_STORAGE_KEY)
    if (!raw) {
      return {
        search: '',
        warehouse_id: '',
        category: ''
      }
    }
    const parsed = JSON.parse(raw)
    return {
      search: parsed.search || '',
      warehouse_id: parsed.warehouse_id || '',
      category: parsed.category || ''
    }
  } catch (error) {
    console.warn('Failed to load saved filters', error)
    return {
      search: '',
      warehouse_id: '',
      category: ''
    }
  }
}
// Pagination state
const rowsPerPageOptions = [10, 20, 50]
const rowsPerPage = ref(rowsPerPageOptions[0])
const currentPage = ref(1)
const sortField = ref('sn')
const sortAsc = ref(true)

const toggleSort = (field) => {
  if (sortField.value === field) sortAsc.value = !sortAsc.value
  else { sortField.value = field; sortAsc.value = true }
}

const getSortIcon = (field) => {
  if (sortField.value !== field) return 'pi pi-sort-alt'
  return sortAsc.value ? 'pi pi-sort-up' : 'pi pi-sort-down'
}

// Sorted and paginated items
const sortedItems = computed(() => {
  const arr = [...items.value]
  if (!sortField.value) return arr
  return arr.sort((a, b) => {
    let va = a[sortField.value] ?? ''
    let vb = b[sortField.value] ?? ''
    if (typeof va === 'string') { va = va.toLowerCase(); vb = vb.toLowerCase() }
    if (va < vb) return sortAsc.value ? -1 : 1
    if (va > vb) return sortAsc.value ? 1 : -1
    return 0
  })
})

const filteredItems = computed(() => {
  let result = sortedItems.value
  
  // Search filter
  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    result = result.filter(item => 
      (item.name?.toLowerCase().includes(search)) ||
      (item.sn?.toLowerCase().includes(search))
    )
  }
  
  // Warehouse filter
  if (filters.value.warehouse_id) {
    result = result.filter(item => 
      item.warehouse_id == filters.value.warehouse_id ||
      item.Warehouse?.id == filters.value.warehouse_id
    )
  }
  
  // Category filter
  if (filters.value.category) {
    result = result.filter(item => item.category === filters.value.category)
  }
  
  return result
})

watch(filters, (newFilters) => {
  try {
    localStorage.setItem(FILTER_STORAGE_KEY, JSON.stringify(newFilters))
  } catch (error) {
    console.warn('Failed to persist inventory filters', error)
  }
}, { deep: true })

const expandedGroups = ref(new Set())

const groupKeyForItem = (item) => {
  const raw = item?.name?.trim().toLowerCase()
  return raw && raw.length ? raw : '__no_name__'
}

const groupedItems = computed(() => {
  const map = new Map()
  filteredItems.value.forEach(item => {
    const key = groupKeyForItem(item)
    if (!map.has(key)) {
      map.set(key, {
        key,
        displayName: item?.name?.trim() || 'Unnamed Item',
        items: []
      })
    }
    map.get(key).items.push(item)
  })
  return Array.from(map.values())
})

const totalPages = computed(() => Math.ceil(groupedItems.value.length / rowsPerPage.value) || 1)

const paginatedGroups = computed(() => {
  const start = (currentPage.value - 1) * rowsPerPage.value
  return groupedItems.value.slice(start, start + rowsPerPage.value)
})

const paginatedRows = computed(() => {
  const rows = []
  paginatedGroups.value.forEach(group => {
    rows.push({
      type: 'group',
      key: `group-${group.key}`,
      group
    })
    if (expandedGroups.value.has(group.key)) {
      group.items.forEach(item => {
        rows.push({
          type: 'item',
          key: `item-${item.id}`,
          item,
          groupKey: group.key
        })
      })
    }
  })
  return rows
})

const itemsOnCurrentPage = computed(() => {
  const items = []
  paginatedGroups.value.forEach(group => {
    items.push(...group.items)
  })
  return items
})

const toggleGroup = (key) => {
  const copy = new Set(expandedGroups.value)
  if (copy.has(key)) copy.delete(key)
  else copy.add(key)
  expandedGroups.value = copy
}

const isGroupExpanded = (key) => expandedGroups.value.has(key)

const changePage = (page) => {
  currentPage.value = Math.min(Math.max(page, 1), totalPages.value)
}

// Reset page on changes
watch(rowsPerPage, () => { currentPage.value = 1 })
watch(filteredItems, () => {
  if (currentPage.value > totalPages.value && totalPages.value > 0) {
    currentPage.value = totalPages.value
  }
  // Clear selections that are no longer visible after filter change
  selectedItems.value = selectedItems.value.filter(id => 
    filteredItems.value.some(item => item.id === id)
  )
})

watch(groupedItems, (groups) => {
  const validKeys = new Set(groups.map(group => group.key))
  const current = expandedGroups.value
  expandedGroups.value = new Set(
    Array.from(current).filter(key => validKeys.has(key))
  )
})

watch(canViewInventory, async (allowed) => {
  if (allowed) {
    await fetchItems()
  } else {
    items.value = []
    selectedItems.value = []
  }
})

const fetchItems = async () => {
  if (!canViewInventory.value) {
    items.value = []
    return
  }

  try {
    const res = await axios.get('/inventory')
    // backend returns { data: items }
    items.value = Array.isArray(res.data.data) ? res.data.data : (res.data.data || [])
    syncCategoriesWithItems()
  } catch (error) {
    const detail = error.response?.data?.error || 'Failed to load items'
    toast.add({ severity: 'error', summary: 'Error', detail, life: 3000 })
  }
}

const fetchWarehouses = async () => {
  try {
    const res = await axios.get('/warehouses')
    warehouses.value = res.data.data || []
  } catch (err) {
    console.error('Failed to load warehouses:', err)
  }
}

const fetchCategoriesList = async () => {
  try {
    const res = await axios.get('/categories')
    const data = Array.isArray(res.data.data) ? res.data.data : (res.data.data ? [res.data.data] : [])
    categories.value = data
      .filter(cat => cat && cat.name)
      .map(cat => ({
        id: cat.id,
        name: cat.name,
        color: cat.color || ''
      }))
      .sort((a, b) => a.name.localeCompare(b.name))
    syncCategoriesWithItems()
  } catch (err) {
    console.error('Failed to load categories:', err)
  }
}

const syncCategoriesWithItems = () => {
  if (!Array.isArray(items.value)) return

  const itemCategories = Array.from(new Set(items.value
    .map(item => item?.category)
    .filter(Boolean)))

  if (!itemCategories.length) return

  const existingNames = new Set(categories.value.map(cat => cat.name))

  const additions = itemCategories
    .filter(name => !existingNames.has(name))
    .map(name => ({
      id: `item-${name}`,
      name,
      color: ''
    }))

  if (!additions.length) return

  categories.value = [...categories.value, ...additions]
    .sort((a, b) => a.name.localeCompare(b.name))
}

const clearFilters = () => {
  filters.value = {
    search: '',
    warehouse_id: '',
    category: ''
  }
  // Clear selections when clearing filters
  selectedItems.value = []
  try {
    localStorage.setItem(FILTER_STORAGE_KEY, JSON.stringify(filters.value))
  } catch (error) {
    console.warn('Failed to clear saved filters', error)
  }
}

const viewItem = (item) => {
  if (!ensurePermission('inventory.view', 'Anda tidak memiliki izin untuk melihat detail item inventori')) {
    return
  }
  // Navigate to detail page
  window.location.href = `/inventory/${item.id}`
}

const editItem = (item) => {
  if (!ensurePermission('inventory.update', 'Anda tidak memiliki izin untuk mengubah item inventori')) {
    return
  }

  isEditing.value = true
  editForm.value = {
    id: item.id,
    has_sn: !!item.sn,
    sn: item.sn || '',
    name: item.name || '',
    category: item.category || '',
    warehouse_id: item.warehouse_id || item.Warehouse?.id || '',
    quantity: item.quantity || 0,
    min_stock: item.min_stock || 0,
    unit: item.unit || 'pcs',
    unit_price: item.unit_price || 0,
    description: item.description || '',
    is_active: item.is_active ?? true
  }
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editForm.value = defaultItemForm()
  isEditing.value = false
}

const submitItemForm = async () => {
  const actionPermission = isEditing.value && editForm.value.id ? 'inventory.update' : 'inventory.create'
  const deniedMessage = actionPermission === 'inventory.update'
    ? 'Anda tidak memiliki izin untuk memperbarui item inventori'
    : 'Anda tidak memiliki izin untuk menambah item inventori'

  if (!ensurePermission(actionPermission, deniedMessage)) {
    return
  }

  const requiredFields = [
    { key: 'name', label: 'Name' },
    { key: 'warehouse_id', label: 'Warehouse' }
  ]
  if (editForm.value.has_sn) {
    requiredFields.unshift({ key: 'sn', label: 'Serial Number' })
  }

  const missingField = requiredFields.find(({ key }) => {
    const value = editForm.value[key]
    if (typeof value === 'string') return !value.trim()
    return value === null || value === undefined || value === ''
  })
  if (missingField) {
    toast.add({
      severity: 'warn',
      summary: 'Incomplete Form',
      detail: `${missingField.label} is required`,
      life: 3000
    })
    return
  }

    const payload = {
      sn: editForm.value.has_sn ? editForm.value.sn?.trim() : '',
      name: editForm.value.name?.trim(),
      category: editForm.value.category || '',
      warehouse_id: Number(editForm.value.warehouse_id),
      quantity: Number(editForm.value.quantity) || 0,
      min_stock: Number(editForm.value.min_stock) || 0,
      unit: editForm.value.unit?.trim() || 'pcs',
      unit_price: Number(editForm.value.unit_price) || 0,
      description: editForm.value.description || '',
      is_active: !!editForm.value.is_active
    }

  if (!payload.warehouse_id) {
    toast.add({
      severity: 'warn',
      summary: 'Incomplete Form',
      detail: 'Please select a warehouse',
      life: 3000
    })
    return
  }

  saving.value = true
  try {
    if (isEditing.value && editForm.value.id) {
      await axios.put(`/inventory/items/${editForm.value.id}`, payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'Item updated successfully', life: 3000 })
    } else {
      await axios.post('/inventory/items', payload)
      toast.add({ severity: 'success', summary: 'Success', detail: 'Item created successfully', life: 3000 })
    }

    closeEditModal()
    fetchItems()
  } catch (err) {
    const detail = err.response?.data?.error || err.response?.data?.message || 'Failed to save item'
    toast.add({ severity: 'error', summary: 'Error', detail, life: 3000 })
  } finally {
    saving.value = false
  }
}

const confirmDelete = (item) => {
  if (!ensurePermission('inventory.delete', 'Anda tidak memiliki izin untuk menghapus item inventori')) {
    return
  }

  itemToDelete.value = item
  showDeleteConfirm.value = true
}

const deleteItem = async () => {
  if (!ensurePermission('inventory.delete', 'Anda tidak memiliki izin untuk menghapus item inventori')) {
    return
  }

  if (!itemToDelete.value) return
  try {
    await axios.delete(`/inventory/items/${itemToDelete.value.id}`)
    toast.add({ severity: 'success', summary: 'Deleted', detail: 'Item deleted', life: 3000 })
    fetchItems()
  } catch {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete item', life: 3000 })
  } finally {
    showDeleteConfirm.value = false
    itemToDelete.value = null
  }
}

const downloadTemplate = async () => {
  if (!ensurePermission('inventory.view', 'Anda tidak memiliki izin untuk mengunduh template inventori')) {
    return
  }

  const res = await axios.get('/inventory/import/template', { responseType: 'blob' })
  const url = URL.createObjectURL(res.data)
  const a = document.createElement('a')
  a.href = url
  a.download = 'template.csv'
  a.click()
  URL.revokeObjectURL(url)
}

const triggerImport = () => {
  if (!ensurePermission('inventory.create', 'Anda tidak memiliki izin untuk mengimpor data inventori')) {
    return
  }
  importInput.value.click()
}

const handleImportChange = async (e) => {
  if (!ensurePermission('inventory.create', 'Anda tidak memiliki izin untuk mengimpor data inventori')) {
    if (e?.target) {
      e.target.value = ''
    }
    return
  }

  const file = e.target.files[0]
  if (!file) return
  const form = new FormData()
  form.append('file', file)
  try {
    await axios.post('/inventory/import/csv', form, { headers: { 'Content-Type': 'multipart/form-data' } })
    toast.add({ severity: 'success', summary: 'Imported', detail: 'CSV imported successfully', life: 3000 })
    fetchItems()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: err.response?.data?.error || 'Import failed', life: 3000 })
  }
}

const buildExportParams = () => {
  const params = {}
  const search = filters.value.search?.trim()
  if (search) params.search = search
  if (filters.value.warehouse_id) params.warehouse_id = filters.value.warehouse_id
  if (filters.value.category) params.category = filters.value.category
  return params
}

const exportCSV = async () => {
  if (!ensurePermission('inventory.view', 'Anda tidak memiliki izin untuk mengekspor data inventori')) {
    return
  }

  try {
    const res = await axios.get('/inventory/export/csv', {
      params: buildExportParams(),
      responseType: 'blob'
    })
    const url = URL.createObjectURL(res.data)
    const a = document.createElement('a')
    a.href = url
    a.download = 'inventory-items.csv'
    a.click()
    URL.revokeObjectURL(url)
    toast.add({ severity: 'success', summary: 'Success', detail: 'CSV exported successfully', life: 3000 })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Export CSV failed', life: 3000 })
  }
}

const exportPDF = async () => {
  if (!ensurePermission('inventory.view', 'Anda tidak memiliki izin untuk mengekspor data inventori')) {
    return
  }

  try {
    const res = await axios.get('/inventory/export/pdf', {
      params: buildExportParams(),
      responseType: 'blob'
    })
    const url = URL.createObjectURL(res.data)
    const a = document.createElement('a')
    a.href = url
    a.download = 'inventory-items.pdf'
    a.click()
    URL.revokeObjectURL(url)
    toast.add({ severity: 'success', summary: 'Success', detail: 'PDF exported successfully', life: 3000 })
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Export PDF failed', life: 3000 })
  }
}

const toggleSelection = (item, checked) => {
  if (!canDeleteInventory.value) {
    return
  }

  const id = item.id
  if (checked) {
    if (!selectedItems.value.includes(id)) {
      selectedItems.value.push(id)
    }
  } else {
    selectedItems.value = selectedItems.value.filter(i => i !== id)
  }
}
const isSelected = (item) => selectedItems.value.includes(item.id)
const toggleSelectAll = (checked) => {
  if (!canDeleteInventory.value) {
    return
  }

  if (checked) {
    const ids = new Set(selectedItems.value)
    itemsOnCurrentPage.value.forEach(item => ids.add(item.id))
    selectedItems.value = Array.from(ids)
  } else {
    const idsToRemove = new Set(itemsOnCurrentPage.value.map(item => item.id))
    selectedItems.value = selectedItems.value.filter(id => !idsToRemove.has(id))
  }
}

const allVisibleSelected = computed(() => {
  if (!canDeleteInventory.value) return false
  if (!itemsOnCurrentPage.value.length) return false
  return itemsOnCurrentPage.value.every(item => selectedItems.value.includes(item.id))
})

const confirmBatchDelete = async () => {
  if (!ensurePermission('inventory.delete', 'Anda tidak memiliki izin untuk menghapus item inventori')) {
    showBatchConfirm.value = false
    return
  }

  if (!selectedItems.value.length) {
    toast.add({ severity: 'warn', summary: 'Peringatan', detail: 'Tidak ada item yang dipilih', life: 3000 })
    showBatchConfirm.value = false
    return
  }

  try {
    await axios.delete('/inventory/items', { data: { ids: selectedItems.value } })
    toast.add({ severity: 'success', summary: 'Deleted', detail: `${selectedItems.value.length} items deleted`, life: 3000 })
    selectedItems.value = []
    fetchItems()
  } catch (err) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Batch delete failed', life: 3000 })
  } finally {
    showBatchConfirm.value = false
  }
}

const openAddItem = () => {
  if (!ensurePermission('inventory.create', 'Anda tidak memiliki izin untuk menambah item inventori')) {
    return
  }

  isEditing.value = false
  editForm.value = defaultItemForm()
  showEditModal.value = true
}

async function toggleItemStatus(item) {
  if (!ensurePermission('inventory.update', 'Anda tidak memiliki izin untuk memperbarui status item inventori')) {
    return
  }

  const newStatus = !item.is_active
  try {
    await axios.put(`/inventory/items/${item.id}`, { is_active: newStatus })
    item.is_active = newStatus
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: `Status item berhasil diubah menjadi ${newStatus ? 'Aktif' : 'Tidak Aktif'}`,
      life: 3000
    })
  } catch (err) {
    console.error('Gagal mengubah status item:', err)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Gagal mengubah status item. Silakan coba lagi.',
      life: 3000
    })
  }
}

onMounted(async () => {
  if (canViewInventory.value) {
    await fetchItems()
  } else {
    toast.add({
      severity: 'warn',
      summary: 'Akses inventori',
      detail: 'Anda tidak memiliki izin untuk melihat data inventory',
      life: 3000
    })
  }

  await fetchWarehouses()
  await fetchCategoriesList()
  
  // Check if there's an edit query parameter
  if (route.query.edit) {
    if (!canUpdateInventory.value) {
      toast.add({
        severity: 'warn',
        summary: 'Akses inventori',
        detail: 'Anda tidak memiliki izin untuk mengubah item inventori',
        life: 3000
      })
    } else {
      const itemId = parseInt(route.query.edit)
      const itemToEdit = items.value.find(item => item.id === itemId)
      
      if (itemToEdit) {
        editItem(itemToEdit)
        // Remove query parameter after opening modal
        router.replace({ path: '/inventory' })
      }
    }
  }
})
</script>

<style scoped>
/* Add component-specific styles if needed */
</style>
