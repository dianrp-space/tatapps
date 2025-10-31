<template>
  <div>
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <div class="flex items-center gap-4">
        <Button 
          icon="pi pi-arrow-left" 
          class="p-button-rounded p-button-text"
          @click="router.push('/inventory')"
        />
        <div>
          <h1 class="text-2xl font-bold text-gray-800">{{ item.name }}</h1>
          <p class="text-gray-600">SN: {{ item.sn }}</p>
        </div>
      </div>
      <div class="flex gap-2">
        <Button 
          label="Record Transaction" 
          icon="pi pi-plus-circle" 
          @click="openTransaction('in')"
          class="p-button-success"
          :disabled="!canRecordTransaction"
          v-tooltip.top="recordRestrictionMessage"
        />
        <Button 
          label="Edit Item" 
          icon="pi pi-pencil" 
          @click="router.push({ path: '/inventory', query: { edit: item.id } })"
          class="p-button-warning"
        />
      </div>
    </div>

    <div v-if="loading" class="flex justify-center items-center h-64">
      <ProgressSpinner />
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Item Details Card -->
      <div class="lg:col-span-2">
        <div class="bg-white rounded-lg shadow p-6 mb-6">
          <h2 class="text-xl font-bold mb-4">Item Information</h2>
          
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="text-sm text-gray-600">SN</label>
              <p class="font-semibold">{{ item.sn }}</p>
            </div>
            
            <div>
              <label class="text-sm text-gray-600">Warehouse</label>
              <p class="font-semibold">{{ item.warehouse?.name }}</p>
            </div>
            
            <div>
              <label class="text-sm text-gray-600">Category</label>
              <p><Tag :value="item.category" :style="getCategoryTagStyle(item.category)" /></p>
            </div>
            
            <div>
              <label class="text-sm text-gray-600">Unit</label>
              <p class="font-semibold">{{ item.unit }}</p>
            </div>
            
            <div class="col-span-2">
              <label class="text-sm text-gray-600">Description</label>
              <p class="text-gray-800">{{ item.description || '-' }}</p>
            </div>
          </div>
        </div>

        <!-- Transaction History -->
        <div class="bg-white rounded-lg shadow p-6">
          <h2 class="text-xl font-bold mb-4">Transaction History</h2>
          
          <DataTable 
            :value="transactions" 
            :loading="loadingTransactions"
            :paginator="true" 
            :rows="10"
            responsiveLayout="scroll"
            stripedRows
          >
            <Column field="created_at" header="Date" sortable>
              <template #body="{ data }">
                {{ formatDate(data.created_at) }}
              </template>
            </Column>
            
            <Column field="type" header="Type" sortable>
              <template #body="{ data }">
                <Tag 
                  :value="data.type" 
                  :severity="getTransactionSeverity(data.type)"
                />
              </template>
            </Column>
            
            <Column field="quantity" header="Quantity" sortable>
              <template #body="{ data }">
                <span :class="getQuantityClass(data.type)">
                  {{ data.type === 'out' || data.type === 'transfer-out' ? '-' : '+' }}{{ data.quantity }}
                </span>
              </template>
            </Column>
            
            <Column field="reference" header="Reference">
              <template #body="{ data }">
                {{ data.reference || '-' }}
              </template>
            </Column>
            
            <Column field="from_warehouse.name" header="From">
              <template #body="{ data }">
                {{ data.from_warehouse?.name || '-' }}
              </template>
            </Column>
            
            <Column field="to_warehouse.name" header="To">
              <template #body="{ data }">
                {{ data.to_warehouse?.name || '-' }}
              </template>
            </Column>
            
            <Column field="created_by.name" header="By">
              <template #body="{ data }">
                {{ data.created_by?.name }}
              </template>
            </Column>
            
            <Column field="notes" header="Notes">
              <template #body="{ data }">
                <span class="text-sm">{{ data.notes || '-' }}</span>
              </template>
            </Column>
          </DataTable>
        </div>
      </div>

      <!-- Stock Info Sidebar -->
      <div class="lg:col-span-1">
        <!-- Current Stock Card -->
        <div class="bg-white rounded-lg shadow p-6 mb-6">
          <h3 class="text-lg font-bold mb-4">Stock Status</h3>
          
          <div class="text-center mb-4">
            <div class="text-5xl font-bold mb-2" :class="getStockColor(item)">
              {{ item.quantity }}
            </div>
            <div class="text-gray-600">{{ item.unit }}</div>
            
            <Tag 
              v-if="item.quantity <= item.min_stock"
              value="Low Stock Alert"
              severity="danger"
              class="mt-2"
            />
          </div>

          <Divider />

          <div class="space-y-3">
            <div class="flex justify-between">
              <span class="text-gray-600">Min Stock:</span>
              <span class="font-semibold">{{ item.min_stock }} {{ item.unit }}</span>
            </div>
            
            <div class="flex justify-between">
              <span class="text-gray-600">Max Stock:</span>
              <span class="font-semibold">{{ item.max_stock }} {{ item.unit }}</span>
            </div>
            
            <Divider />
            
            <div class="flex justify-between">
              <span class="text-gray-600">Unit Price:</span>
              <span class="font-semibold">{{ formatCurrency(item.unit_price) }}</span>
            </div>
            
            <div class="flex justify-between text-lg">
              <span class="text-gray-600">Total Value:</span>
              <span class="font-bold text-blue-600">
                {{ formatCurrency(item.quantity * item.unit_price) }}
              </span>
            </div>
          </div>
        </div>

        <!-- Quick Actions Card -->
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-bold mb-4">Quick Actions</h3>
          
          <div class="space-y-2">
            <Button 
              label="Stock In" 
              icon="pi pi-arrow-down" 
              @click="openTransaction('in')"
              class="w-full p-button-success"
              :disabled="!canRecordTransaction"
              v-tooltip.top="!canRecordTransaction ? recordRestrictionMessage : null"
            />
            
            <Button 
              label="Stock Out" 
              icon="pi pi-arrow-up" 
              @click="openTransaction('out')"
              class="w-full p-button-danger"
              :disabled="!canRecordTransaction || item.quantity <= 0"
              v-tooltip.top="!canRecordTransaction ? recordRestrictionMessage : (item.quantity <= 0 ? 'No stock available' : null)"
            />
            
            <Button 
              label="Transfer" 
              icon="pi pi-arrow-right-arrow-left" 
              @click="openTransaction('transfer')"
              class="w-full p-button-info"
              :disabled="!canRecordTransaction || !hasTransferDestination"
              v-tooltip.top="!canRecordTransaction ? recordRestrictionMessage : (!hasTransferDestination ? 'No available destination warehouse' : null)"
            />
            
            <Button 
              label="Adjust Stock" 
              icon="pi pi-sync" 
              @click="openTransaction('adjustment')"
              class="w-full p-button-warning"
              :disabled="!canRecordTransaction"
              v-tooltip.top="!canRecordTransaction ? recordRestrictionMessage : null"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Transaction Dialog -->
    <Dialog 
      v-model:visible="showTransactionDialog" 
      :header="getTransactionTitle()"
      :modal="true"
      :style="{ width: '500px' }"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Transaction Type *</label>
          <Dropdown 
            v-model="transactionForm.type" 
            :options="transactionTypes"
            optionLabel="label"
            optionValue="value"
            placeholder="Select Type"
            class="w-full"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Quantity *</label>
          <InputNumber 
            v-model="transactionForm.quantity" 
            :min="1"
            :max="transactionForm.type === 'out' ? item.quantity : undefined"
            class="w-full"
          />
          <small v-if="transactionForm.type === 'out'" class="text-gray-500">
            Available: {{ item.quantity }} {{ item.unit }}
          </small>
        </div>

        <div v-if="transactionForm.type === 'transfer'">
          <label class="block text-sm font-medium text-gray-700 mb-2">To Warehouse *</label>
          <Dropdown 
            v-model="transactionForm.to_warehouse_id" 
            :options="transferCandidates"
            optionLabel="name"
            optionValue="id"
            placeholder="Select Warehouse"
            class="w-full"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Reference</label>
          <InputText 
            v-model="transactionForm.reference" 
            placeholder="e.g., PO-001, WO-123"
            class="w-full"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Notes</label>
          <Textarea 
            v-model="transactionForm.notes" 
            rows="3"
            placeholder="Transaction notes..."
            class="w-full"
          />
        </div>
      </div>

      <template #footer>
        <Button 
          label="Cancel" 
          icon="pi pi-times" 
          @click="showTransactionDialog = false"
          class="p-button-text"
        />
        <Button 
          label="Record Transaction" 
          icon="pi pi-check" 
          @click="saveTransaction"
          :loading="savingTransaction"
        />
      </template>
    </Dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import axios from '@/api/axios'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const toast = useToast()

const item = ref({})
const transactions = ref([])
const warehouses = ref([])
const rawWarehouses = ref([])
const categories = ref([])
const loading = ref(false)
const loadingTransactions = ref(false)
const savingTransaction = ref(false)
const showTransactionDialog = ref(false)

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

const authStore = useAuthStore()

const transactionTypes = ref([
  { label: 'Stock In', value: 'in' },
  { label: 'Stock Out', value: 'out' },
  { label: 'Transfer', value: 'transfer' },
  { label: 'Adjustment', value: 'adjustment' }
])

const transactionForm = ref({
  type: 'in',
  quantity: 1,
  to_warehouse_id: null,
  reference: '',
  notes: ''
})

const normalizeId = (value) => {
  if (value === null || value === undefined || value === '') {
    return null
  }
  const parsed = Number(value)
  return Number.isNaN(parsed) ? value : parsed
}

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

const hasTransactionPermission = computed(() => {
  if (typeof authStore.hasAnyPermission !== 'function') {
    return false
  }
  return authStore.hasAnyPermission(['inventory.update', 'inventory.create'])
})

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

const isItemInAllowedWarehouse = computed(() => {
  const itemWarehouseId = toNumericId(item.value?.warehouse_id ?? item.value?.warehouse?.id)
  if (!isEmployee.value) {
    return true
  }
  if (itemWarehouseId === null) {
    return false
  }
  return isWarehouseAllowed(itemWarehouseId)
})

const canRecordTransaction = computed(() => {
  if (!hasWarehouseAccess.value || !isItemInAllowedWarehouse.value) {
    return false
  }

  if (isEmployee.value) {
    return true
  }

  return hasTransactionPermission.value
})

const transferCandidates = computed(() => {
  return (warehouses.value || []).filter(w => w.id !== item.value?.warehouse_id)
})

const hasTransferDestination = computed(() => transferCandidates.value.length > 0)

const recordRestrictionMessage = computed(() => {
  if (canRecordTransaction.value) {
    return null
  }
  if (!hasWarehouseAccess.value) {
    return 'No warehouse access assigned to your account.'
  }
  if (!isItemInAllowedWarehouse.value) {
    return 'This item belongs to a warehouse outside your access.'
  }
  if (!hasTransactionPermission.value && !isEmployee.value) {
    return 'You do not have permission to manage inventory transactions.'
  }
  return 'You are not allowed to record transactions for this item.'
})

onMounted(() => {
  fetchItem()
  fetchTransactions()
  fetchWarehouses()
  fetchCategories()
})

async function fetchCategories() {
  try {
    const response = await axios.get('/categories')
    categories.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch categories:', error)
  }
}

async function fetchItem() {
  loading.value = true
  try {
    const response = await axios.get(`/inventory/items/${route.params.id}`)
    item.value = response.data.data || {}
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to fetch item details',
      life: 3000
    })
    router.push('/inventory')
  } finally {
    loading.value = false
  }
}

async function fetchTransactions() {
  loadingTransactions.value = true
  try {
    const response = await axios.get(`/inventory/items/${route.params.id}/transactions`)
    transactions.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch transactions:', error)
  } finally {
    loadingTransactions.value = false
  }
}

async function fetchWarehouses() {
  try {
    const response = await axios.get('/warehouses')
    rawWarehouses.value = response.data.data || []
    filterWarehouses()
  } catch (error) {
    console.error('Failed to fetch warehouses:', error)
  }
}

const filterWarehouses = () => {
  if (!Array.isArray(rawWarehouses.value)) {
    warehouses.value = []
    return
  }

  if (!isEmployee.value) {
    warehouses.value = rawWarehouses.value.slice()
    return
  }

  warehouses.value = rawWarehouses.value.filter(warehouse => isWarehouseAllowed(warehouse?.id))
}

watch(allowedWarehouseIds, () => {
  filterWarehouses()
})

watch(() => isEmployee.value, () => {
  filterWarehouses()
})

function openTransaction(type) {
  if (!canRecordTransaction.value) {
    toast.add({
      severity: 'warn',
      summary: 'Access Restricted',
      detail: 'You are not allowed to record transactions for this item.',
      life: 3000
    })
    return
  }

  transactionForm.value = {
    type,
    quantity: 1,
    to_warehouse_id: null,
    reference: '',
    notes: ''
  }

  if (type === 'transfer') {
    if (transferCandidates.value.length === 0) {
      toast.add({
        severity: 'warn',
        summary: 'Transfer Not Available',
        detail: 'No permitted destination warehouse found.',
        life: 3000
      })
      return
    }
    if (transferCandidates.value.length === 1) {
      transactionForm.value.to_warehouse_id = transferCandidates.value[0].id
    }
  }

  showTransactionDialog.value = true
}

function getTransactionTitle() {
  const typeLabels = {
    'in': 'Stock In',
    'out': 'Stock Out',
    'transfer': 'Transfer Stock',
    'adjustment': 'Adjust Stock'
  }
  return typeLabels[transactionForm.value.type] || 'Record Transaction'
}

async function saveTransaction() {
  if (!canRecordTransaction.value) {
    toast.add({
      severity: 'warn',
      summary: 'Access Restricted',
      detail: 'You are not allowed to record transactions for this item.',
      life: 3000
    })
    return
  }

  if (transactionForm.value.type === 'transfer') {
    if (!transactionForm.value.to_warehouse_id || !isWarehouseAllowed(transactionForm.value.to_warehouse_id)) {
      toast.add({
        severity: 'warn',
        summary: 'Validation Error',
        detail: 'Please choose a permitted destination warehouse.',
        life: 3000
      })
      return
    }
  }

  savingTransaction.value = true
  try {
    const payload = {
      item_id: item.value.id,
      type: transactionForm.value.type,
      quantity: transactionForm.value.quantity,
      reference: transactionForm.value.reference,
      notes: transactionForm.value.notes
    }

    if (transactionForm.value.type === 'transfer') {
      payload.from_warehouse_id = item.value.warehouse_id
      payload.to_warehouse_id = transactionForm.value.to_warehouse_id
    }

    await axios.post(`/inventory/items/${route.params.id}/transactions`, payload)
    
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Transaction recorded successfully',
      life: 3000
    })
    
    showTransactionDialog.value = false
    fetchItem()
    fetchTransactions()
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.message || 'Failed to record transaction',
      life: 3000
    })
  } finally {
    savingTransaction.value = false
  }
}

function getStockColor(item) {
  if (item.quantity <= item.min_stock) {
    return 'text-red-600'
  } else if (item.quantity <= item.min_stock * 2) {
    return 'text-orange-600'
  }
  return 'text-green-600'
}

function getTransactionSeverity(type) {
  const severities = {
    'in': 'success',
    'out': 'danger',
    'transfer': 'info',
    'transfer-in': 'success',
    'transfer-out': 'warning',
    'adjustment': 'warning'
  }
  return severities[type] || 'info'
}

function getQuantityClass(type) {
  if (type === 'out' || type === 'transfer-out') {
    return 'text-red-600 font-semibold'
  }
  return 'text-green-600 font-semibold'
}

function formatCurrency(value) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value)
}

function formatDate(date) {
  return new Date(date).toLocaleString('id-ID', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>
