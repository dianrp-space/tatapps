<template>
  <div class="space-y-6">
    <div class="flex justify-between items-center">
      <div>
        <h1 class="text-2xl font-bold text-gray-800">Transaction History</h1>
        <p class="text-gray-600 mt-1">View and manage all inventory transactions</p>
      </div>
      <Button 
        label="New Transaction" 
        icon="pi pi-plus-circle"
        raised
        @click="openTransactionDialog()"
        :disabled="!canCreateTransaction"
        v-tooltip.top="createRestrictionMessage"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-sm md:text-base border-none shadow-md rounded-lg disabled:bg-gray-300 disabled:hover:bg-gray-300 disabled:border-gray-300 disabled:text-gray-600"
      />
    </div>

    <!-- Filters -->
    <div class="bg-white p-4 rounded-lg shadow-sm">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Date Range</label>
          <Calendar 
            v-model="filters.dateRange" 
            selectionMode="range" 
            dateFormat="dd/mm/yy"
            placeholder="Select date range"
            class="w-full"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Transaction Type</label>
          <Dropdown 
            v-model="filters.type" 
            :options="typeOptions"
            optionLabel="label"
            optionValue="value"
            placeholder="All Types"
            class="w-full"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Warehouse</label>
          <Dropdown 
            v-model="filters.warehouse_id" 
            :options="warehouses"
            optionLabel="name"
            optionValue="id"
            placeholder="All Warehouses"
            class="w-full"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Search</label>
          <InputText 
            v-model="filters.search" 
            placeholder="Search by item, reference..."
            class="w-full"
          />
        </div>
      </div>
      <div class="flex justify-end gap-2 mt-4">
        <Button 
          label="Clear Filters" 
          icon="pi pi-filter-slash" 
          @click="clearFilters"
          class="p-button-outlined"
        />
        <Button 
          label="Export CSV" 
          icon="pi pi-file" 
          @click="exportTransactions('csv')"
          class="p-button-outlined"
          :loading="exportingFormat === 'csv'"
          :disabled="exportingFormat && exportingFormat !== 'csv'"
        />
        <Button 
          label="Export PDF" 
          icon="pi pi-file-pdf" 
          @click="exportTransactions('pdf')"
          class="p-button-outlined"
          :loading="exportingFormat === 'pdf'"
          :disabled="exportingFormat && exportingFormat !== 'pdf'"
        />
      </div>
    </div>

    <!-- Transactions Table -->
    <div class="bg-white rounded-lg shadow-sm">
      <div class="overflow-x-auto">
        <fwb-table hoverable>
          <fwb-table-head>
            <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('created_at')">
              <div class="flex items-center gap-2">
                Date
                <i :class="getSortIcon('created_at')"></i>
              </div>
            </fwb-table-head-cell>
            <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('type')">
              <div class="flex items-center gap-2">
                Type
                <i :class="getSortIcon('type')"></i>
              </div>
            </fwb-table-head-cell>
            <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('item.name')">
              <div class="flex items-center gap-2">
                Item
                <i :class="getSortIcon('item.name')"></i>
              </div>
            </fwb-table-head-cell>
            <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('quantity')">
              <div class="flex items-center gap-2">
                Quantity
                <i :class="getSortIcon('quantity')"></i>
              </div>
            </fwb-table-head-cell>
            <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('from_warehouse.name')">
              <div class="flex items-center gap-2">
                From
                <i :class="getSortIcon('from_warehouse.name')"></i>
              </div>
            </fwb-table-head-cell>
            <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('to_warehouse.name')">
              <div class="flex items-center gap-2">
                To
                <i :class="getSortIcon('to_warehouse.name')"></i>
              </div>
            </fwb-table-head-cell>
            <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('reference')">
              <div class="flex items-center gap-2">
                Reference
                <i :class="getSortIcon('reference')"></i>
              </div>
            </fwb-table-head-cell>
            <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('notes')">
              <div class="flex items-center gap-2">
                Notes
                <i :class="getSortIcon('notes')"></i>
              </div>
            </fwb-table-head-cell>
            <fwb-table-head-cell class="cursor-pointer select-none" @click="toggleSort('created_by.full_name')">
              <div class="flex items-center gap-2">
                Created By
                <i :class="getSortIcon('created_by.full_name')"></i>
              </div>
            </fwb-table-head-cell>
            <fwb-table-head-cell>
              <span class="sr-only">Actions</span>
            </fwb-table-head-cell>
          </fwb-table-head>

          <fwb-table-body>
            <fwb-table-row v-if="loading">
              <fwb-table-cell :colspan="9">
                <div class="py-6 text-center text-gray-500">Loading transactions...</div>
              </fwb-table-cell>
            </fwb-table-row>

            <fwb-table-row v-else-if="!paginatedTransactions.length">
              <fwb-table-cell :colspan="9">
                <div class="py-8 text-center text-gray-500">
                  <i class="pi pi-inbox text-3xl text-gray-300 mb-2"></i>
                  <p>No transactions found</p>
                </div>
              </fwb-table-cell>
            </fwb-table-row>

            <fwb-table-row
              v-else
              v-for="transaction in paginatedTransactions"
              :key="transaction.id"
              class="align-middle"
            >
              <fwb-table-cell>
                <div class="text-sm">
                  <div class="font-medium text-gray-800">{{ formatDate(transaction.created_at) }}</div>
                  <div class="text-gray-500">{{ formatTime(transaction.created_at) }}</div>
                </div>
              </fwb-table-cell>
              <fwb-table-cell>
                <Tag 
                  :value="transaction.type.toUpperCase()" 
                  :severity="getTypeSeverity(transaction.type)"
                  :icon="getTypeIcon(transaction.type)"
                />
              </fwb-table-cell>
              <fwb-table-cell>
                <div class="text-sm">
                  <div class="font-medium text-gray-800">{{ transaction.item?.name || '-' }}</div>
                  <div class="text-gray-500">SN: {{ transaction.item?.sn || '-' }}</div>
                </div>
              </fwb-table-cell>
              <fwb-table-cell>
                <span 
                  class="font-semibold"
                  :class="getQuantityClass(transaction.type)"
                >
                  {{ transaction.type === 'out' ? '-' : '+' }}{{ transaction.quantity }} {{ transaction.item?.unit }}
                </span>
              </fwb-table-cell>
              <fwb-table-cell>
                <span class="text-sm text-gray-600">{{ transaction.from_warehouse?.name || '-' }}</span>
              </fwb-table-cell>
              <fwb-table-cell>
                <span class="text-sm text-gray-600">{{ transaction.to_warehouse?.name || '-' }}</span>
              </fwb-table-cell>
            <fwb-table-cell>
              <span class="text-sm text-gray-600">{{ transaction.reference || '-' }}</span>
            </fwb-table-cell>
            <fwb-table-cell>
              <span class="text-sm text-gray-600 whitespace-pre-line">{{ transaction.notes || '-' }}</span>
            </fwb-table-cell>
            <fwb-table-cell>
              <div class="flex items-center gap-2">
                <img 
                  :src="getUserAvatar(transaction.created_by)" 
                    class="w-8 h-8 rounded-full"
                    :alt="transaction.created_by?.full_name"
                  />
                  <span class="text-sm text-gray-700">{{ transaction.created_by?.full_name || '-' }}</span>
                </div>
              </fwb-table-cell>
              <fwb-table-cell>
                <div class="flex justify-end gap-2">
                  <Button
                    icon="pi pi-eye"
                    class="p-button-rounded p-button-text p-button-info"
                    @click="viewTransaction(transaction)"
                    v-tooltip.top="'View Details'"
                  />
                  <Button
                    icon="pi pi-trash"
                    class="p-button-rounded p-button-text p-button-danger"
                    @click="confirmDeleteTransaction(transaction)"
                    v-tooltip.top="'Delete Transaction'"
                  />
                </div>
              </fwb-table-cell>
            </fwb-table-row>
          </fwb-table-body>
        </fwb-table>
      </div>

      <div class="flex flex-col gap-4 border-t border-gray-100 px-4 py-3 md:flex-row md:items-center md:justify-between">
        <div class="text-sm text-gray-600">
          Showing
          <span class="font-semibold">{{ paginationStats.from }}</span>
          -
          <span class="font-semibold">{{ paginationStats.to }}</span>
          of
          <span class="font-semibold">{{ paginationStats.total }}</span>
          transactions
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <label class="text-sm text-gray-600">
            Rows per page
            <select
              v-model.number="rowsPerPage"
              class="ml-2 rounded border-gray-300 text-sm focus:border-blue-500 focus:ring-blue-500"
            >
              <option v-for="option in rowsPerPageOptions" :key="option" :value="option">
                {{ option }}
              </option>
            </select>
          </label>
          <div class="flex items-center gap-2">
            <button
              class="rounded border border-gray-200 px-3 py-1 text-sm text-gray-600 transition hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="currentPage === 1"
              @click="changePage(currentPage - 1)"
            >
              Prev
            </button>
            <span class="text-sm text-gray-600">
              Page
              <span class="font-semibold">{{ totalPages ? currentPage : 0 }}</span>
              /
              <span class="font-semibold">{{ totalPages }}</span>
            </span>
            <button
              class="rounded border border-gray-200 px-3 py-1 text-sm text-gray-600 transition hover:bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50"
              :disabled="currentPage === totalPages || totalPages === 0"
              @click="changePage(currentPage + 1)"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Transaction Dialog -->
    <Dialog 
      v-model:visible="transactionDialogVisible" 
      header="New Transaction"
      :modal="true"
      class="w-full max-w-2xl"
    >
      <div class="space-y-6 pt-4">
        <!-- Transaction Type -->
        <div class="grid grid-cols-3 gap-4">
          <div
            v-for="type in typeOptions"
            :key="type.value"
            @click="transactionForm.type = type.value"
            class="p-4 border-2 rounded-lg cursor-pointer transition-all text-center"
            :class="transactionForm.type === type.value 
              ? 'border-blue-600 bg-blue-50' 
              : 'border-gray-200 hover:border-gray-300'"
          >
            <i :class="type.icon" class="text-3xl mb-2" :style="{ color: type.color }"></i>
            <div class="font-semibold">{{ type.label }}</div>
            <div class="text-xs text-gray-500 mt-1">{{ type.description }}</div>
          </div>
        </div>

        <!-- Item Selection -->
        <div v-if="transactionForm.type !== 'transfer'">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Warehouse <span class="text-red-500">*</span>
          </label>
          <Dropdown 
            v-model="transactionForm.warehouse_id" 
            :options="warehouses"
            optionLabel="name"
            optionValue="id"
            placeholder="Select warehouse"
            class="w-full"
            showClear
          />
        </div>

        <!-- Transfer specific fields -->
        <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              From Warehouse <span class="text-red-500">*</span>
            </label>
            <Dropdown 
              v-model="transactionForm.from_warehouse_id" 
              :options="warehouses"
              optionLabel="name"
              optionValue="id"
              placeholder="Select warehouse"
              class="w-full"
              showClear
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              To Warehouse <span class="text-red-500">*</span>
            </label>
            <Dropdown 
              v-model="transactionForm.to_warehouse_id" 
              :options="destinationWarehouseOptions"
              optionLabel="name"
              optionValue="id"
              placeholder="Select warehouse"
              class="w-full"
              :disabled="!transactionForm.from_warehouse_id"
              showClear
            />
          </div>
        </div>

        <!-- Item Selection -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Item <span class="text-red-500">*</span>
          </label>
          <Dropdown 
            v-model="transactionForm.item_id" 
            :options="filteredItemOptions"
            optionLabel="name"
            optionValue="id"
            :placeholder="canSelectItem ? 'Select item' : 'Select warehouse first'"
            filter
            class="w-full"
            :disabled="!canSelectItem"
            @change="onItemSelect"
          >
            <template #option="{ option }">
              <div class="flex justify-between items-center">
                <div>
                  <div class="font-medium">{{ option.name }}</div>
                  <div class="text-sm text-gray-500">SN: {{ option.sn }}</div>
                </div>
                <Tag :value="`Stock: ${option.quantity}`" severity="info" />
              </div>
            </template>
            <template #empty>
              <div class="py-2 px-3 text-sm text-gray-500">
                {{ canSelectItem ? 'No items available for selected warehouse' : 'Select a warehouse first' }}
              </div>
            </template>
          </Dropdown>
        </div>

        <!-- Quantity -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Quantity <span class="text-red-500">*</span>
          </label>
          <InputNumber 
            v-model="transactionForm.quantity" 
            :min="0"
            placeholder="Enter quantity"
            class="w-full"
          />
          <small v-if="selectedItem" class="text-gray-500">
            Current stock: {{ selectedItem.quantity }} {{ selectedItem.unit }}
          </small>
        </div>

        <!-- Reference -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Reference Number
          </label>
          <InputText 
            v-model="transactionForm.reference" 
            placeholder="e.g., PO-2024-001"
            class="w-full"
          />
        </div>

        <!-- Notes -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Notes
          </label>
          <Textarea 
            v-model="transactionForm.notes" 
            rows="3"
            placeholder="Add notes about this transaction"
            class="w-full"
          />
        </div>
      </div>

      <template #footer>
        <Button 
          label="Cancel" 
          icon="pi pi-times" 
          severity="danger"
          class="font-bold bg-red-500 border-none text-white hover:bg-red-600 focus:bg-red-700"
          @click="transactionDialogVisible = false"
        />
        <Button 
          label="Save Transaction" 
          icon="pi pi-check" 
          severity="info"
          class="font-bold bg-blue-600 border-none text-white hover:bg-blue-700 focus:bg-blue-800"
          @click="saveTransaction"
          :loading="loading"
        />
      </template>
    </Dialog>

    <!-- View Transaction Dialog -->
    <Dialog 
      v-model:visible="viewDialogVisible" 
      header="Transaction Details"
      :modal="true"
      class="w-full max-w-md"
    >
      <div v-if="selectedTransaction" class="space-y-4 pt-4">
        <div class="flex items-center justify-between pb-4 border-b">
          <Tag 
            :value="selectedTransaction.type.toUpperCase()" 
            :severity="getTypeSeverity(selectedTransaction.type)"
            :icon="getTypeIcon(selectedTransaction.type)"
          />
          <span class="text-sm text-gray-500">{{ formatDateTime(selectedTransaction.created_at) }}</span>
        </div>

        <div class="space-y-3">
          <div>
            <label class="text-sm text-gray-500">Item</label>
            <p class="font-medium">{{ selectedTransaction.item?.name }}</p>
            <p class="text-sm text-gray-600">SN: {{ selectedTransaction.item?.sn }}</p>
          </div>

          <div>
            <label class="text-sm text-gray-500">Quantity</label>
            <p class="font-medium text-lg" :class="getQuantityClass(selectedTransaction.type)">
              {{ selectedTransaction.type === 'out' ? '-' : '+' }}{{ selectedTransaction.quantity }} {{ selectedTransaction.item?.unit }}
            </p>
          </div>

          <div v-if="selectedTransaction.from_warehouse">
            <label class="text-sm text-gray-500">From Warehouse</label>
            <p class="font-medium">{{ selectedTransaction.from_warehouse.name }}</p>
          </div>

          <div v-if="selectedTransaction.to_warehouse">
            <label class="text-sm text-gray-500">To Warehouse</label>
            <p class="font-medium">{{ selectedTransaction.to_warehouse.name }}</p>
          </div>

          <div v-if="selectedTransaction.reference">
            <label class="text-sm text-gray-500">Reference</label>
            <p class="font-medium">{{ selectedTransaction.reference }}</p>
          </div>

          <div v-if="selectedTransaction.notes">
            <label class="text-sm text-gray-500">Notes</label>
            <p class="text-gray-700">{{ selectedTransaction.notes }}</p>
          </div>

          <div>
            <label class="text-sm text-gray-500">Created By</label>
            <div class="flex items-center gap-2 mt-1">
              <img 
                :src="getUserAvatar(selectedTransaction?.created_by, 40)" 
                class="w-8 h-8 rounded-full"
              />
              <span class="font-medium">{{ selectedTransaction.created_by?.full_name }}</span>
            </div>
          </div>
        </div>
      </div>

      <template #footer>
        <Button 
          label="Close" 
          icon="pi pi-times" 
          severity="info"
          class="font-bold bg-gray-500 border-none text-white hover:bg-gray-600 focus:bg-gray-700"
          @click="viewDialogVisible = false"
        />
      </template>
    </Dialog>

    <!-- Delete Confirmation Dialog -->
    <Dialog
      v-model:visible="deleteDialogVisible"
      header="Delete Transaction"
      :modal="true"
      :style="{ width: '420px' }"
    >
      <div class="flex gap-3 py-2">
        <i class="pi pi-exclamation-triangle text-red-500 text-3xl"></i>
        <div class="space-y-1">
          <p class="font-semibold text-gray-800">
            Are you sure you want to delete this transaction?
          </p>
          <p class="text-sm text-gray-600">
            {{ transactionToDelete?.item?.name || 'Unknown item' }} —
            {{ transactionToDelete ? formatDateTime(transactionToDelete.created_at) : '' }}
          </p>
          <p class="text-xs text-gray-500">
            This action will revert the stock change made by the transaction.
          </p>
        </div>
      </div>

      <template #footer>
        <Button 
          label="Cancel" 
          icon="pi pi-times" 
          severity="danger"
          class="font-bold bg-red-500 border-none text-white hover:bg-red-600 focus:bg-red-700"
          :disabled="deletingTransaction"
          @click="deleteDialogVisible = false"
        />
        <Button 
          label="Delete" 
          icon="pi pi-trash" 
          severity="danger"
          class="font-bold bg-orange-500 border-none text-white hover:bg-orange-600 focus:bg-orange-700"
          :loading="deletingTransaction"
          @click="deleteTransaction"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import {
  FwbTable,
  FwbTableBody,
  FwbTableCell,
  FwbTableHead,
  FwbTableHeadCell,
  FwbTableRow,
} from 'flowbite-vue'
import Tag from 'primevue/tag'
import axios from '@/api/axios'
import { useAuthStore } from '@/stores/auth'

const toast = useToast()

const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'

const fileBaseUrl = (() => {
  try {
    const url = new URL(apiBaseUrl)
    url.pathname = url.pathname.replace(/\/api\/v1\/?$/, '')
    const path = url.pathname.replace(/\/$/, '')
    return `${url.origin}${path === '' ? '' : path}`
  } catch {
    return apiBaseUrl.replace(/\/api\/v1\/?$/, '')
  }
})()

const getAvatarUrl = (avatarPath, size = 32) => {
  if (!avatarPath) {
    return `https://via.placeholder.com/${size}`
  }
  if (avatarPath.startsWith('http')) {
    return avatarPath
  }
  const base = fileBaseUrl || apiBaseUrl
  const normalizedBase = base.endsWith('/') ? base.slice(0, -1) : base
  const normalizedPath = avatarPath.startsWith('/') ? avatarPath.slice(1) : avatarPath
  return `${normalizedBase}/${normalizedPath}`
}

const getUserAvatar = (user, size = 32) => getAvatarUrl(user?.avatar, size)

const normalizeId = (value) => {
  if (value === null || value === undefined || value === '') {
    return null
  }
  const parsed = Number(value)
  return Number.isNaN(parsed) ? value : parsed
}

const getItemWarehouseId = (item) => {
  if (!item) {
    return null
  }
  if (item.warehouse_id !== undefined && item.warehouse_id !== null) {
    return normalizeId(item.warehouse_id)
  }
  return normalizeId(item.warehouse?.id ?? null)
}

const authStore = useAuthStore()

const isEmployee = computed(() => authStore.user?.role?.name === 'employee')

const allowedWarehouseIds = computed(() => {
  const entries = authStore.user?.warehouses || []
  const collected = []
  entries.forEach(entry => {
    if (!entry) return
    const candidate = entry.warehouse_id ?? entry.warehouse?.id
    const numeric = Number(normalizeId(candidate))
    if (!Number.isNaN(numeric) && numeric > 0 && !collected.includes(numeric)) {
      collected.push(numeric)
    }
  })
  return collected
})

const allowedWarehouseSet = computed(() => {
  const set = new Set()
  allowedWarehouseIds.value.forEach(id => {
    const numeric = Number(id)
    if (!Number.isNaN(numeric)) {
      set.add(numeric)
    }
  })
  return set
})

const hasWarehouseAccess = computed(() => !isEmployee.value || allowedWarehouseIds.value.length > 0)

const canCreateTransaction = computed(() => {
  const hasTxnPermission = typeof authStore.hasAnyPermission === 'function'
    ? authStore.hasAnyPermission(['inventory.update', 'inventory.create'])
    : false

  if (!hasWarehouseAccess.value) {
    return false
  }

  if (hasTxnPermission) {
    return true
  }

  return isEmployee.value
})

const createRestrictionMessage = computed(() => {
  if (canCreateTransaction.value) {
    return null
  }
  if (!hasWarehouseAccess.value) {
    return 'No warehouse access assigned to your account.'
  }
  if (isEmployee.value) {
    return 'You are not allowed to record transactions for the available warehouses.'
  }
  return 'You do not have permission to manage inventory transactions.'
})

const toNumericId = (value) => {
  if (value === null || value === undefined || value === '') {
    return null
  }
  const normalized = normalizeId(value)
  if (normalized === null || normalized === undefined || normalized === '') {
    return null
  }
  const numeric = Number(normalized)
  return Number.isNaN(numeric) ? null : numeric
}

const isWarehouseAllowed = (warehouseId) => {
  if (!isEmployee.value) {
    return true
  }
  const numeric = toNumericId(warehouseId)
  if (numeric === null) {
    return false
  }
  return allowedWarehouseSet.value.has(numeric)
}

const applyWarehouseFilter = (list = []) => {
  if (!Array.isArray(list)) {
    return []
  }
  if (!isEmployee.value) {
    return list.slice()
  }
  if (!hasWarehouseAccess.value) {
    return []
  }
  return list.filter(warehouse => isWarehouseAllowed(warehouse?.id))
}

const applyItemFilter = (list = []) => {
  if (!Array.isArray(list)) {
    return []
  }
  if (!isEmployee.value) {
    return list.slice()
  }
  if (!hasWarehouseAccess.value) {
    return []
  }
  return list.filter(item => isWarehouseAllowed(getItemWarehouseId(item)))
}

const applyTransactionFilter = (list = []) => {
  if (!Array.isArray(list)) {
    return []
  }
  if (!isEmployee.value) {
    return list.slice()
  }
  if (!hasWarehouseAccess.value) {
    return []
  }
  return list.filter(transaction => {
    const itemData = transaction?.item || transaction?.Item
    const itemWarehouseId = toNumericId(itemData?.warehouse_id ?? itemData?.warehouse?.id ?? itemData?.Warehouse?.id)
    const fromWarehouseId = toNumericId(transaction?.from_warehouse_id ?? transaction?.from_warehouse?.id ?? transaction?.FromWarehouse?.id)
    const toWarehouseId = toNumericId(transaction?.to_warehouse_id ?? transaction?.to_warehouse?.id ?? transaction?.ToWarehouse?.id)
    return [itemWarehouseId, fromWarehouseId, toWarehouseId].some(id => id !== null && allowedWarehouseSet.value.has(id))
  })
}

const transactions = ref([])
const items = ref([])
const warehouses = ref([])
const rawTransactions = ref([])
const rawItems = ref([])
const rawWarehouses = ref([])
const loading = ref(false)
const transactionDialogVisible = ref(false)
const viewDialogVisible = ref(false)
const selectedTransaction = ref(null)
const selectedItem = ref(null)
const deleteDialogVisible = ref(false)
const transactionToDelete = ref(null)
const deletingTransaction = ref(false)

const filters = ref({
  dateRange: null,
  type: null,
  warehouse_id: null,
  search: ''
})

const rowsPerPageOptions = [10, 20, 50]
const rowsPerPage = ref(rowsPerPageOptions[0])
const currentPage = ref(1)
const exportingFormat = ref(null)

const sortState = ref({
  field: 'created_at',
  order: 'desc'
})

const transactionForm = ref({
  type: 'in',
  item_id: null,
  quantity: 0,
  warehouse_id: null,
  from_warehouse_id: null,
  to_warehouse_id: null,
  reference: '',
  notes: ''
})

const typeOptions = [
  { 
    label: 'Stock In', 
    value: 'in', 
    icon: 'pi pi-arrow-down',
    color: '#10B981',
    description: 'Add stock to warehouse'
  },
  { 
    label: 'Stock Out', 
    value: 'out', 
    icon: 'pi pi-arrow-up',
    color: '#EF4444',
    description: 'Remove stock from warehouse'
  },
  { 
    label: 'Transfer', 
    value: 'transfer', 
    icon: 'pi pi-arrows-h',
    color: '#3B82F6',
    description: 'Move between warehouses'
  },
  { 
    label: 'Adjustment', 
    value: 'adjustment', 
    icon: 'pi pi-sync',
    color: '#F59E0B',
    description: 'Adjust stock quantity'
  }
]

const filteredItemOptions = computed(() => {
  const sourceId = transactionForm.value.type === 'transfer'
    ? transactionForm.value.from_warehouse_id
    : transactionForm.value.warehouse_id

  const numericSourceId = toNumericId(sourceId)
  if (!numericSourceId) {
    return []
  }

  if (!isWarehouseAllowed(numericSourceId)) {
    return []
  }

  return (items.value || []).filter(item => {
    const itemWarehouseId = toNumericId(getItemWarehouseId(item))
    return itemWarehouseId !== null && itemWarehouseId === numericSourceId
  })
})

const canSelectItem = computed(() => {
  if (transactionForm.value.type === 'transfer') {
    return Boolean(transactionForm.value.from_warehouse_id)
  }
  return Boolean(transactionForm.value.warehouse_id)
})

const destinationWarehouseOptions = computed(() => {
  const excludeId = toNumericId(transactionForm.value.from_warehouse_id)
  const source = isEmployee.value ? rawWarehouses.value : warehouses.value
  const baseList = Array.isArray(source) ? source : []
  return baseList.filter(warehouse => {
    const warehouseId = toNumericId(warehouse?.id)
    if (warehouseId === null) {
      return false
    }
    if (!excludeId) {
      return true
    }
    return warehouseId !== excludeId
  })
})

const filteredTransactions = computed(() => {
  let result = (transactions.value || []).slice()

  if (filters.value.type) {
    result = result.filter(t => t.type === filters.value.type)
  }

  if (filters.value.warehouse_id) {
    const targetId = toNumericId(filters.value.warehouse_id)
    result = result.filter(t => {
      const fromId = toNumericId(t.from_warehouse_id ?? t.from_warehouse?.id ?? t.FromWarehouse?.id)
      const toId = toNumericId(t.to_warehouse_id ?? t.to_warehouse?.id ?? t.ToWarehouse?.id)
      const itemData = t.item || t.Item
      const itemId = toNumericId(itemData?.warehouse_id ?? itemData?.warehouse?.id ?? itemData?.Warehouse?.id)
      return (fromId !== null && fromId === targetId) ||
        (toId !== null && toId === targetId) ||
        (itemId !== null && itemId === targetId)
    })
  }

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    result = result.filter(t => {
      const nameMatch = (t.item?.name || '').toLowerCase().includes(search)
      const snMatch = (t.item?.sn || '').toLowerCase().includes(search)
      const referenceMatch = (t.reference || '').toLowerCase().includes(search)
      const notesMatch = (t.notes || '').toLowerCase().includes(search)
      return nameMatch || snMatch || referenceMatch || notesMatch
    })
  }

  if (filters.value.dateRange && filters.value.dateRange[0]) {
    const [start, end] = filters.value.dateRange
    const startTime = start ? new Date(start).setHours(0, 0, 0, 0) : null
    const endTime = end ? new Date(end).setHours(23, 59, 59, 999) : null

    result = result.filter(t => {
      const created = new Date(t.created_at).getTime()
      if (Number.isNaN(created)) {
        return true
      }
      if (startTime && created < startTime) {
        return false
      }
      if (endTime && created > endTime) {
        return false
      }
      return true
    })
  }

  return result
})

function getFieldValue(transaction, field) {
  if (!field) {
    return null
  }
  return field.split('.').reduce((acc, key) => (acc == null ? acc : acc[key]), transaction)
}

const sortedTransactions = computed(() => {
  const data = filteredTransactions.value.slice()
  const { field, order } = sortState.value

  if (!field) {
    return data
  }

  return data.sort((a, b) => {
    const aVal = getFieldValue(a, field)
    const bVal = getFieldValue(b, field)

    if (aVal == null && bVal == null) return 0
    if (aVal == null) return order === 'asc' ? 1 : -1
    if (bVal == null) return order === 'asc' ? -1 : 1

    if (typeof aVal === 'number' && typeof bVal === 'number') {
      return order === 'asc' ? aVal - bVal : bVal - aVal
    }

    const aStr = String(aVal).toLowerCase()
    const bStr = String(bVal).toLowerCase()

    if (aStr === bStr) return 0
    return order === 'asc'
      ? (aStr > bStr ? 1 : -1)
      : (aStr < bStr ? 1 : -1)
  })
})

const totalPages = computed(() => {
  const total = sortedTransactions.value.length
  if (total === 0) {
    return 0
  }
  return Math.ceil(total / rowsPerPage.value)
})

const paginatedTransactions = computed(() => {
  if (!sortedTransactions.value.length) {
    return []
  }
  const safePage = totalPages.value === 0
    ? 1
    : Math.min(Math.max(currentPage.value, 1), totalPages.value)
  const start = (safePage - 1) * rowsPerPage.value
  return sortedTransactions.value.slice(start, start + rowsPerPage.value)
})

const paginationStats = computed(() => {
  const total = sortedTransactions.value.length
  if (!total) {
    return { from: 0, to: 0, total: 0 }
  }
  const from = (currentPage.value - 1) * rowsPerPage.value + 1
  const to = Math.min(from + rowsPerPage.value - 1, total)
  return { from, to, total }
})

onMounted(() => {
  fetchTransactions()
  fetchItems()
  fetchWarehouses()
})

watch(rowsPerPage, () => {
  currentPage.value = 1
})

watch(filteredTransactions, () => {
  currentPage.value = 1
})

watch(allowedWarehouseIds, () => {
  warehouses.value = applyWarehouseFilter(rawWarehouses.value)
  items.value = applyItemFilter(rawItems.value)
  transactions.value = applyTransactionFilter(rawTransactions.value)

  if (isEmployee.value && filters.value.warehouse_id && !isWarehouseAllowed(filters.value.warehouse_id)) {
    filters.value.warehouse_id = null
  }
})

watch(sortedTransactions, () => {
  if (totalPages.value > 0 && currentPage.value > totalPages.value) {
    currentPage.value = totalPages.value
  }
  if (totalPages.value === 0) {
    currentPage.value = 1
  }
})

watch(() => transactionForm.value.type, (newType, oldType) => {
  if (newType === oldType) {
    return
  }
  transactionForm.value.item_id = null
  selectedItem.value = null

  if (newType === 'transfer') {
    transactionForm.value.warehouse_id = null
  } else {
    transactionForm.value.from_warehouse_id = null
    transactionForm.value.to_warehouse_id = null
  }
})

watch(() => transactionForm.value.warehouse_id, (newValue, oldValue) => {
  if (newValue === oldValue) {
    return
  }
  transactionForm.value.item_id = null
  selectedItem.value = null
})

watch(() => transactionForm.value.from_warehouse_id, (newValue, oldValue) => {
  if (newValue === oldValue) {
    return
  }
  if (transactionForm.value.type === 'transfer') {
    transactionForm.value.item_id = null
    selectedItem.value = null
    if (transactionForm.value.to_warehouse_id && normalizeId(transactionForm.value.to_warehouse_id) === normalizeId(newValue)) {
      transactionForm.value.to_warehouse_id = null
    }
  }
})

async function fetchTransactions() {
  loading.value = true
  try {
    const response = await axios.get('/inventory/transactions')
    rawTransactions.value = response.data.data || []
    transactions.value = applyTransactionFilter(rawTransactions.value)
    currentPage.value = 1
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load transactions',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

async function fetchItems() {
  try {
    const response = await axios.get('/inventory')
    rawItems.value = response.data.data || []
    items.value = applyItemFilter(rawItems.value)
  } catch (error) {
    console.error('Failed to load items:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load inventory items',
      life: 3000
    })
  }
}

async function fetchWarehouses() {
  try {
    const response = await axios.get('/warehouses')
    rawWarehouses.value = response.data.data || []
    warehouses.value = applyWarehouseFilter(rawWarehouses.value)
  } catch (error) {
    console.error('Failed to load warehouses:', error)
  }
}

function openTransactionDialog() {
  if (!canCreateTransaction.value) {
    toast.add({
      severity: 'warn',
      summary: 'Access Restricted',
      detail: createRestrictionMessage.value || 'You do not have permission to create transactions for any warehouse.',
      life: 3000
    })
    return
  }

  transactionForm.value = {
    type: 'in',
    item_id: null,
    quantity: 0,
    warehouse_id: null,
    from_warehouse_id: null,
    to_warehouse_id: null,
    reference: '',
    notes: ''
  }
  selectedItem.value = null

  if (isEmployee.value && allowedWarehouseIds.value.length > 0) {
    const defaultWarehouse = allowedWarehouseIds.value[0]
    transactionForm.value.warehouse_id = defaultWarehouse
    transactionForm.value.from_warehouse_id = defaultWarehouse
  }

  transactionDialogVisible.value = true
}

function onItemSelect(event) {
  const matchedItem = items.value.find(i => i.id === event.value) || null
  selectedItem.value = matchedItem

  if (!matchedItem) {
    return
  }

  const targetWarehouseId = transactionForm.value.type === 'transfer'
    ? normalizeId(transactionForm.value.from_warehouse_id)
    : normalizeId(transactionForm.value.warehouse_id)
  const itemWarehouseId = normalizeId(getItemWarehouseId(matchedItem))

  if (targetWarehouseId && itemWarehouseId && itemWarehouseId !== targetWarehouseId) {
    toast.add({
      severity: 'warn',
      summary: 'Validation Error',
      detail: 'Selected item does not belong to the chosen warehouse',
      life: 3000
    })
    transactionForm.value.item_id = null
    selectedItem.value = null
    return
  }

  console.log('Item selected:', selectedItem.value)
}

async function saveTransaction() {
  // Validation
  if (!transactionForm.value.item_id) {
    toast.add({
      severity: 'warn',
      summary: 'Validation Error',
      detail: 'Please select an item',
      life: 3000
    })
    return
  }

  if (!transactionForm.value.quantity || transactionForm.value.quantity <= 0) {
    toast.add({
      severity: 'warn',
      summary: 'Validation Error',
      detail: 'Please enter a valid quantity',
      life: 3000
    })
    return
  }

  if (transactionForm.value.type !== 'transfer' && !transactionForm.value.warehouse_id) {
    toast.add({
      severity: 'warn',
      summary: 'Validation Error',
      detail: 'Please select a warehouse',
      life: 3000
    })
    return
  }

  if (transactionForm.value.type === 'transfer') {
    if (!transactionForm.value.from_warehouse_id || !transactionForm.value.to_warehouse_id) {
      toast.add({
        severity: 'warn',
        summary: 'Validation Error',
        detail: 'Please select both source and destination warehouses',
        life: 3000
      })
      return
    }
  }

  if (isEmployee.value) {
    if (transactionForm.value.type === 'transfer') {
      if (!isWarehouseAllowed(transactionForm.value.from_warehouse_id)) {
        toast.add({
          severity: 'warn',
          summary: 'Access Restricted',
          detail: 'Selected source warehouse is not permitted for your account',
          life: 3000
        })
        return
      }
    } else if (!isWarehouseAllowed(transactionForm.value.warehouse_id)) {
      toast.add({
        severity: 'warn',
        summary: 'Access Restricted',
        detail: 'Selected warehouse is not permitted for your account',
        life: 3000
      })
      return
    }

    const selected = items.value.find(item => item.id === transactionForm.value.item_id)
    if (!selected || !isWarehouseAllowed(getItemWarehouseId(selected))) {
      toast.add({
        severity: 'warn',
        summary: 'Access Restricted',
        detail: 'Selected item does not belong to an allowed warehouse',
        life: 3000
      })
      return
    }
  }

  loading.value = true
  try {
    const payload = {
      type: transactionForm.value.type,
      quantity: parseInt(transactionForm.value.quantity),
      reference: transactionForm.value.reference || '',
      notes: transactionForm.value.notes || ''
    }

    if (transactionForm.value.type === 'transfer') {
      payload.from_warehouse_id = transactionForm.value.from_warehouse_id
      payload.to_warehouse_id = transactionForm.value.to_warehouse_id
    }

    const url = `/inventory/items/${transactionForm.value.item_id}/transactions`
    const response = await axios.post(url, payload)
    
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Transaction recorded successfully',
      life: 3000
    })
    
    transactionDialogVisible.value = false
    await fetchTransactions()
    await fetchItems()
  } catch (error) {
    console.error('=== TRANSACTION ERROR ===')
    console.error('Error:', error)
    console.error('Response:', error.response)
    console.error('Data:', error.response?.data)
    console.error('Status:', error.response?.status)
    console.error('Headers:', error.response?.headers)
    console.error('=== END ERROR ===')
    
    const errorMsg = error.response?.data?.error || error.response?.data?.message || 'Failed to record transaction'
    
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: errorMsg,
      life: 5000
    })
  } finally {
    loading.value = false
  }
}

function viewTransaction(transaction) {
  selectedTransaction.value = transaction
  viewDialogVisible.value = true
}

function clearFilters() {
  filters.value = {
    dateRange: null,
    type: null,
    warehouse_id: null,
    search: ''
  }
}

function formatDateForRange(value, boundary) {
  if (!value) return null
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return null
  if (boundary === 'start') {
    date.setHours(0, 0, 0, 0)
  } else {
    date.setHours(23, 59, 59, 999)
  }
  const pad = (num) => String(num).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

function buildExportParams() {
  const params = {}
  if (filters.value.type) {
    params.type = filters.value.type
  }
  if (filters.value.warehouse_id) {
    params.warehouse_id = filters.value.warehouse_id
  }
  if (filters.value.search) {
    params.search = filters.value.search
  }
  if (filters.value.dateRange && filters.value.dateRange[0]) {
    const start = formatDateForRange(filters.value.dateRange[0], 'start')
    if (start) params.start_date = start
  }
  if (filters.value.dateRange && filters.value.dateRange[1]) {
    const end = formatDateForRange(filters.value.dateRange[1], 'end')
    if (end) params.end_date = end
  }
  return params
}

async function exportTransactions(format) {
  if (exportingFormat.value) {
    return
  }

  exportingFormat.value = format
  const isPDF = format === 'pdf'
  const endpoint = isPDF ? '/inventory/transactions/export/pdf' : '/inventory/transactions/export/csv'
  const fileExtension = isPDF ? 'pdf' : 'csv'
  const mimeType = isPDF ? 'application/pdf' : 'text/csv'

  try {
    const response = await axios.get(endpoint, {
      params: buildExportParams(),
      responseType: 'blob'
    })

    const blob = new Blob([response.data], { type: mimeType })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    const timestamp = new Date().toISOString().slice(0, 19).replace(/[:T]/g, '-')
    link.href = url
    link.download = `inventory-transactions-${timestamp}.${fileExtension}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)

    toast.add({
      severity: 'success',
      summary: 'Export',
      detail: `Transactions exported as ${fileExtension.toUpperCase()}`,
      life: 3000
    })
  } catch (error) {
    console.error('Failed to export transactions:', error)
    const detail = error.response?.data?.error || error.response?.data?.message || 'Export failed'
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail,
      life: 4000
    })
  } finally {
    exportingFormat.value = null
  }
}

function changePage(page) {
  if (totalPages.value === 0) {
    return
  }
  const safePage = Math.min(Math.max(page, 1), totalPages.value)
  currentPage.value = safePage
}

function toggleSort(field) {
  if (sortState.value.field === field) {
    sortState.value = {
      field,
      order: sortState.value.order === 'asc' ? 'desc' : 'asc'
    }
  } else {
    sortState.value = { field, order: 'asc' }
  }
}

function getSortIcon(field) {
  if (sortState.value.field !== field) {
    return 'pi pi-sort-alt text-gray-400 text-xs'
  }
  return sortState.value.order === 'asc'
    ? 'pi pi-sort-amount-up text-blue-500 text-xs'
    : 'pi pi-sort-amount-down text-blue-500 text-xs'
}

function getTypeSeverity(type) {
  const severities = {
    in: 'success',
    out: 'danger',
    transfer: 'info',
    adjustment: 'warning'
  }
  return severities[type] || 'info'
}

function getTypeIcon(type) {
  const icons = {
    in: 'pi pi-arrow-down',
    out: 'pi pi-arrow-up',
    transfer: 'pi pi-arrows-h',
    adjustment: 'pi pi-sync'
  }
  return icons[type] || 'pi pi-circle'
}

function getQuantityClass(type) {
  if (type === 'out') {
    return 'text-red-600'
  }
  if (type === 'transfer') {
    return 'text-blue-600'
  }
  return 'text-green-600'
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleDateString('id-ID')
}

function formatTime(dateString) {
  return new Date(dateString).toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' })
}

function formatDateTime(dateString) {
  return new Date(dateString).toLocaleString('id-ID')
}

function confirmDeleteTransaction(transaction) {
  transactionToDelete.value = transaction
  deleteDialogVisible.value = true
}

async function deleteTransaction() {
  if (!transactionToDelete.value) {
    return
  }

  deletingTransaction.value = true
  try {
    await axios.delete(`/inventory/transactions/${transactionToDelete.value.id}`)
    toast.add({
      severity: 'success',
      summary: 'Deleted',
      detail: 'Transaction deleted successfully',
      life: 3000
    })
    deleteDialogVisible.value = false
    transactionToDelete.value = null
    await fetchTransactions()
    await fetchItems()
  } catch (error) {
    console.error('Delete transaction error:', error)
    const errorMsg = error.response?.data?.error || error.response?.data?.message || 'Failed to delete transaction'
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: errorMsg,
      life: 4000
    })
  } finally {
    deletingTransaction.value = false
  }
}
</script>
