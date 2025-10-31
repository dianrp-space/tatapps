<template>
  <div>
    <!-- Header -->
    <div class="flex justify-between items-center mb-6">
      <div class="flex items-center gap-4">
        <Button 
          icon="pi pi-arrow-left" 
          class="p-button-rounded p-button-text"
          @click="$router.push('/warehouses')"
        />
        <div>
          <h1 class="text-2xl font-bold text-gray-800">{{ warehouse.name }}</h1>
          <p class="text-gray-600">{{ warehouse.code }}</p>
        </div>
      </div>
      <div class="flex gap-2">
        <Button 
          label="Edit" 
          icon="pi pi-pencil" 
          class="p-button-warning"
          v-if="authStore.isAdmin || authStore.isManager"
        />
      </div>
    </div>

    <div v-if="loading" class="flex justify-center items-center h-64">
      <ProgressSpinner />
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Warehouse Info Card -->
      <div class="lg:col-span-2">
        <div class="bg-white rounded-lg shadow p-6 mb-6">
          <h2 class="text-xl font-bold mb-4">Warehouse Information</h2>
          
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="text-sm text-gray-600">Code</label>
              <p class="font-semibold">{{ warehouse.code }}</p>
            </div>
            
            <div>
              <label class="text-sm text-gray-600">Status</label>
              <p>
                <Tag 
                  :value="warehouse.is_active ? 'Active' : 'Inactive'" 
                  :severity="warehouse.is_active ? 'success' : 'danger'"
                />
              </p>
            </div>
            
            <div class="col-span-2">
              <label class="text-sm text-gray-600">Address</label>
              <p class="text-gray-800">{{ warehouse.address }}</p>
            </div>
            
            <div>
              <label class="text-sm text-gray-600">City</label>
              <p class="font-semibold">{{ warehouse.city }}</p>
            </div>
            
            <div>
              <label class="text-sm text-gray-600">Province</label>
              <p class="font-semibold">{{ warehouse.province }}</p>
            </div>
            
            <div>
              <label class="text-sm text-gray-600">Postal Code</label>
              <p class="font-semibold">{{ warehouse.postal_code }}</p>
            </div>
            
            <div>
              <label class="text-sm text-gray-600">Phone</label>
              <p class="font-semibold">{{ warehouse.phone || '-' }}</p>
            </div>
            
            <div>
              <label class="text-sm text-gray-600">Email</label>
              <p class="font-semibold">{{ warehouse.email || '-' }}</p>
            </div>
            
            <div>
              <label class="text-sm text-gray-600">Manager</label>
              <p class="font-semibold">{{ warehouse.manager?.full_name || '-' }}</p>
            </div>
          </div>
        </div>

        <!-- Inventory Items -->
        <div class="bg-white rounded-lg shadow p-6">
          <div class="flex justify-between items-center mb-4">
            <h2 class="text-xl font-bold">Inventory Items</h2>
            <Button 
              label="Add Item" 
              icon="pi pi-plus-circle"
              raised
              class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 text-sm md:text-base border-none shadow-md rounded-lg"
              @click="$router.push('/inventory')"
            />
          </div>
          
          <DataTable 
            :value="inventoryItems" 
            :loading="loadingInventory"
            :paginator="true" 
            :rows="10"
            responsiveLayout="scroll"
            stripedRows
          >
            <Column field="sn" header="SN">
              <template #body="{ data }">
                <span class="font-mono text-sm">{{ data.sn }}</span>
              </template>
            </Column>
            
            <Column field="name" header="Item Name">
              <template #body="{ data }">
                <div>
                  <div class="font-semibold">{{ data.name }}</div>
                  <div class="text-sm text-gray-500">{{ data.category }}</div>
                </div>
              </template>
            </Column>
            
            <Column field="quantity" header="Stock">
              <template #body="{ data }">
                <div class="flex items-center gap-2">
                  <span class="font-bold" :class="getStockColor(data)">
                    {{ data.quantity }} {{ data.unit }}
                  </span>
                  <Tag
                    v-if="showLowTag(data)"
                    value="Low"
                    severity="danger"
                  />
                </div>
              </template>
            </Column>
            
            <Column field="unit_price" header="Unit Price">
              <template #body="{ data }">
                {{ formatCurrency(data.unit_price) }}
              </template>
            </Column>
            
            <Column header="Actions">
              <template #body="{ data }">
                <Button 
                  icon="pi pi-eye" 
                  class="p-button-rounded p-button-text p-button-sm"
                  @click="$router.push(`/inventory/${data.id}`)"
                  v-tooltip.top="'View Details'"
                />
              </template>
            </Column>
          </DataTable>
        </div>
      </div>

      <!-- Stats Sidebar -->
      <div class="lg:col-span-1">
        <!-- Total Items Card -->
        <div class="bg-white rounded-lg shadow p-6 mb-6">
          <h3 class="text-lg font-bold mb-4">Statistics</h3>
          
          <div class="space-y-4">
            <div class="border-b pb-3">
              <div class="text-sm text-gray-600">Total Items</div>
              <div class="text-2xl font-bold text-blue-600">{{ inventoryItems.length }}</div>
            </div>
            
            <div class="border-b pb-3">
              <div class="text-sm text-gray-600">Low Stock Items</div>
              <div class="text-2xl font-bold text-red-600">{{ lowStockCount }}</div>
            </div>
            
            <div class="border-b pb-3">
              <div class="text-sm text-gray-600">Total Value</div>
              <div class="text-2xl font-bold text-green-600">{{ formatCurrency(totalValue) }}</div>
            </div>
            
            <div>
              <div class="text-sm text-gray-600">Categories</div>
              <div class="text-2xl font-bold text-purple-600">{{ uniqueCategories }}</div>
            </div>
          </div>
        </div>

        <!-- Quick Actions -->
        <div class="bg-white rounded-lg shadow p-6">
          <h3 class="text-lg font-bold mb-4">Quick Actions</h3>
          
          <div class="space-y-2">
            <Button 
              label="View All Inventory" 
              icon="pi pi-box" 
              @click="$router.push('/inventory')"
              class="w-full"
            />
            
            <Button 
              label="Add New Item" 
              icon="pi pi-plus" 
              @click="$router.push('/inventory')"
              class="w-full p-button-success"
            />
            
            <Button 
              label="Generate Report" 
              icon="pi pi-file-pdf" 
              class="w-full p-button-secondary"
              :loading="generatingReport"
              @click="generateReport"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useWarehouseStore } from '@/stores/warehouse'
import axios from '@/api/axios'
import { useToast } from 'primevue/usetoast'

const route = useRoute()
const authStore = useAuthStore()
const warehouseStore = useWarehouseStore()
const toast = useToast()

const warehouse = ref({})
const inventoryItems = ref([])
const loading = ref(false)
const loadingInventory = ref(false)
const generatingReport = ref(false)

function normalizeQuantity(value) {
  const numeric = Number(value ?? 0)
  return Number.isFinite(numeric) ? numeric : 0
}

function normalizeMinStock(item) {
  if (!item) return 0
  const raw = item.min_stock ?? item.minStock ?? 0
  const numeric = Number(raw)
  return Number.isFinite(numeric) ? numeric : 0
}

function isSerialItem(item) {
  if (!item) return false
  return typeof item.sn === 'string' && item.sn.trim().length > 0
}

const lowStockEntries = computed(() => {
  const aggregated = new Map()
  const lowStock = []

  inventoryItems.value.forEach((rawItem) => {
    if (!rawItem) return
    const item = {
      id: rawItem.id,
      sn: (rawItem.sn || '').trim(),
      name: rawItem.name || '',
      category: (rawItem.category || '').trim(),
      unit: (rawItem.unit || '').trim(),
      quantity: normalizeQuantity(rawItem.quantity),
      minStock: normalizeMinStock(rawItem)
    }

    const hasSerial = item.sn !== ''

    if (hasSerial) {
      const categoryKey = item.category !== '' ? item.category : 'Uncategorized'
      const key = categoryKey.toLowerCase()

      if (!aggregated.has(key)) {
        aggregated.set(key, {
          category: categoryKey,
          quantity: 0,
          minStock: 0,
          unit: item.unit || 'unit',
          itemIds: []
        })
      }

      const entry = aggregated.get(key)
      entry.quantity += item.quantity
      if (item.minStock > entry.minStock) {
        entry.minStock = item.minStock
      }
      if (!entry.unit && item.unit) {
        entry.unit = item.unit
      }
      entry.itemIds.push(item.id)
    } else {
      if (item.minStock > 0 && item.quantity <= item.minStock) {
        lowStock.push({
          id: item.id,
          category: item.category,
          quantity: item.quantity,
          minStock: item.minStock,
          unit: item.unit || 'unit',
          aggregated: false
        })
      }
    }
  })

  aggregated.forEach((entry) => {
    if (entry.minStock > 0 && entry.quantity <= entry.minStock) {
      lowStock.push({
        category: entry.category,
        quantity: entry.quantity,
        minStock: entry.minStock,
        unit: entry.unit || 'unit',
        aggregated: true,
        itemIds: entry.itemIds
      })
    }
  })

  return lowStock
})

const lowStockCount = computed(() => lowStockEntries.value.length)

const totalValue = computed(() => {
  return inventoryItems.value.reduce((total, item) => {
    return total + (item.quantity * item.unit_price)
  }, 0)
})

const uniqueCategories = computed(() => {
  const categories = new Set(inventoryItems.value.map(item => item.category))
  return categories.size
})

onMounted(() => {
  fetchWarehouse()
  fetchInventoryItems()
})

async function fetchWarehouse() {
  loading.value = true
  try {
    await warehouseStore.fetchWarehouse(route.params.id)
    warehouse.value = warehouseStore.currentWarehouse || {}
  } catch (error) {
    console.error('Failed to fetch warehouse:', error)
  } finally {
    loading.value = false
  }
}

async function fetchInventoryItems() {
  loadingInventory.value = true
  try {
    const response = await axios.get(`/inventory?warehouse_id=${route.params.id}`)
    inventoryItems.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch inventory items:', error)
  } finally {
    loadingInventory.value = false
  }
}

function getStockColor(item) {
  if (!item) return 'text-gray-700'
  if (isSerialItem(item)) {
    return 'text-gray-700'
  }

  const quantity = normalizeQuantity(item.quantity)
  const minStock = normalizeMinStock(item)

  if (minStock > 0) {
    if (quantity <= minStock) {
      return 'text-red-600'
    }
    if (quantity <= minStock * 2) {
      return 'text-orange-600'
    }
  } else if (quantity <= 0) {
    return 'text-red-600'
  }
  return 'text-green-600'
}

function showLowTag(item) {
  if (!item || isSerialItem(item)) {
    return false
  }
  const quantity = normalizeQuantity(item.quantity)
  const minStock = normalizeMinStock(item)
  return minStock > 0 && quantity <= minStock
}

function formatCurrency(value) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value)
}

function toCsvValue(value) {
  if (value === null || value === undefined) {
    return ''
  }
  const stringValue = String(value).replace(/"/g, '""')
  if (/[",\n]/.test(stringValue)) {
    return `"${stringValue}"`
  }
  return stringValue
}

function toCsvRow(values) {
  return values.map(toCsvValue).join(',')
}

async function generateReport() {
  if (generatingReport.value) return

  generatingReport.value = true
  try {
    if (!inventoryItems.value.length) {
      toast.add({
        severity: 'warn',
        summary: 'No Data',
        detail: 'Tidak ada item inventory untuk dilaporkan.',
        life: 3000
      })
      return
    }

    const headers = [
      'SN',
      'Item Name',
      'Category',
      'Unit',
      'Quantity',
      'Minimum Stock',
      'Unit Price (IDR)',
      'Total Value (IDR)'
    ]

    const rows = inventoryItems.value.map(item => {
      const total = Number(item.quantity || 0) * Number(item.unit_price || 0)
      return [
        item.sn || '-',
        item.name || '',
        item.category || '',
        item.unit || '',
        item.quantity ?? 0,
        item.min_stock ?? 0,
        item.unit_price ?? 0,
        total
      ]
    })

    const csvLines = [
      toCsvRow(headers),
      ...rows.map(toCsvRow),
      '',
      toCsvRow(['Summary']),
      toCsvRow(['Warehouse', warehouse.value.name || warehouse.value.code || route.params.id || '-']),
      toCsvRow(['Total Items', inventoryItems.value.length]),
      toCsvRow(['Low Stock Items', lowStockCount.value]),
      toCsvRow(['Total Inventory Value', formatCurrency(totalValue.value)])
    ]

    const csvContent = csvLines.join('\n')
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)

    const safeIdentifier = (warehouse.value.code || warehouse.value.name || 'warehouse')
      .toString()
      .toLowerCase()
      .replace(/[^a-z0-9-_]+/g, '-')
      .replace(/^-+|-+$/g, '') || 'warehouse'
    const timestamp = new Date().toISOString().slice(0, 10)
    const fileName = `${safeIdentifier}-report-${timestamp}.csv`

    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', fileName)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)

    toast.add({
      severity: 'success',
      summary: 'Report Generated',
      detail: `Laporan ${fileName} berhasil dibuat.`,
      life: 3000
    })
  } catch (error) {
    console.error('Failed to generate report:', error)
    toast.add({
      severity: 'error',
      summary: 'Gagal Membuat Laporan',
      detail: 'Terjadi kesalahan saat membuat laporan warehouse.',
      life: 4000
    })
  } finally {
    generatingReport.value = false
  }
}
</script>
