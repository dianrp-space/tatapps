<template>
  <div class="p-6">
    <div class="flex items-center gap-3 mb-6">
      <button
        type="button"
        @click="activeTab = 'users'"
        :class="[
          'px-4 py-2 rounded-lg border transition-colors',
          activeTab === 'users'
            ? 'bg-blue-600 text-white border-blue-600'
            : 'bg-white text-gray-600 border-gray-200 hover:bg-gray-100'
        ]"
      >
        User Management
      </button>
      <button
        type="button"
        @click="activeTab = 'roles'"
        :class="[
          'px-4 py-2 rounded-lg border transition-colors',
          activeTab === 'roles'
            ? 'bg-blue-600 text-white border-blue-600'
            : 'bg-white text-gray-600 border-gray-200 hover:bg-gray-100'
        ]"
      >
        Role Management
      </button>
    </div>

    <div v-if="activeTab === 'users'">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-900">User Management</h2>
        <button
          @click="openAddModal"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:ring-4 focus:ring-blue-300"
        >
          <i class="pi pi-plus mr-2"></i>
          Add User
        </button>
      </div>

      <div class="bg-white rounded-lg shadow overflow-hidden">
        <fwb-table hoverable>
          <fwb-table-head>
            <fwb-table-head-cell>Name</fwb-table-head-cell>
            <fwb-table-head-cell>Email</fwb-table-head-cell>
            <fwb-table-head-cell>Phone</fwb-table-head-cell>
            <fwb-table-head-cell>Role</fwb-table-head-cell>
            <fwb-table-head-cell>Warehouses</fwb-table-head-cell>
            <fwb-table-head-cell>Status</fwb-table-head-cell>
            <fwb-table-head-cell>Actions</fwb-table-head-cell>
          </fwb-table-head>
          <fwb-table-body>
            <fwb-table-row v-for="user in users" :key="user.id">
              <fwb-table-cell>{{ user.full_name }}</fwb-table-cell>
              <fwb-table-cell>{{ user.email }}</fwb-table-cell>
              <fwb-table-cell>{{ user.phone || '-' }}</fwb-table-cell>
            <fwb-table-cell>
              <Tag
                class="border-none rounded-full px-3 py-1 shadow-sm"
                :style="roleTagStyle(user.role)"
              >
                <span class="flex items-center gap-2">
                  <span
                    class="w-2.5 h-2.5 rounded-full"
                    :style="roleColorDotStyle(user.role)"
                  ></span>
                  {{ user.role?.name || 'N/A' }}
                </span>
              </Tag>
            </fwb-table-cell>
            <fwb-table-cell>
              <div class="flex flex-wrap gap-1">
                <Tag
                  v-for="warehouse in user.warehouses || []"
                  :key="warehouse.id || `${user.id}-${warehouse.warehouse_id}`"
                  class="text-xs border-none px-2 py-1 rounded-full"
                  :style="warehouseTagStyle(warehouse.warehouse)"
                >
                  <span class="flex items-center gap-2">
                    <span
                      class="w-2 h-2 rounded-full"
                      :style="warehouseDotStyle(warehouse.warehouse)"
                    ></span>
                    {{ warehouse.warehouse?.name || warehouse.warehouse?.code || 'Warehouse' }}
                  </span>
                </Tag>
                <span v-if="!user.warehouses || user.warehouses.length === 0" class="text-xs text-gray-400">-</span>
              </div>
            </fwb-table-cell>
            <fwb-table-cell>
              <Tag
                :value="user.is_active ? 'Active' : 'Inactive'"
                :severity="user.is_active ? 'success' : 'danger'"
              />
              </fwb-table-cell>
              <fwb-table-cell>
                <div class="flex gap-2">
                  <button
                    @click="openEditModal(user)"
                    class="text-blue-600 hover:text-blue-800"
                    title="Edit"
                  >
                    <i class="pi pi-pencil text-lg"></i>
                  </button>
                  <button
                    @click="toggleUserStatus(user)"
                    :class="user.is_active ? 'text-red-600 hover:text-red-800' : 'text-green-600 hover:text-green-800'"
                    :title="user.is_active ? 'Deactivate' : 'Activate'"
                  >
                    <i :class="user.is_active ? 'pi pi-ban' : 'pi pi-check-circle'" class="text-lg"></i>
                  </button>
                  <button
                    @click="deleteUser(user)"
                    class="text-red-600 hover:text-red-800"
                    title="Delete"
                  >
                    <i class="pi pi-trash text-lg"></i>
                  </button>
                </div>
              </fwb-table-cell>
            </fwb-table-row>
          </fwb-table-body>
        </fwb-table>

        <div v-if="users.length === 0" class="text-center py-8 text-gray-500">
          No users found.
        </div>
      </div>
    </div>

    <div v-else>
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-900">Role Management</h2>
        <button
          @click="openRoleModal()"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 focus:ring-4 focus:ring-blue-300"
        >
          <i class="pi pi-plus mr-2"></i>
          Add Role
        </button>
      </div>

      <div class="bg-white rounded-lg shadow overflow-hidden">
        <fwb-table hoverable>
          <fwb-table-head>
            <fwb-table-head-cell>Role</fwb-table-head-cell>
            <fwb-table-head-cell>Description</fwb-table-head-cell>
            <fwb-table-head-cell>Visible Menus</fwb-table-head-cell>
            <fwb-table-head-cell>Permissions</fwb-table-head-cell>
            <fwb-table-head-cell>Actions</fwb-table-head-cell>
          </fwb-table-head>
          <fwb-table-body>
            <fwb-table-row v-for="role in roles" :key="role.id">
              <fwb-table-cell class="capitalize">
                <Tag
                  class="border-none rounded-full px-3 py-1 shadow-sm"
                  :style="roleTagStyle(role)"
                >
                  <span class="flex items-center gap-2">
                    <span
                      class="w-2.5 h-2.5 rounded-full"
                      :style="roleColorDotStyle(role)"
                    ></span>
                    {{ role.name }}
                  </span>
                </Tag>
              </fwb-table-cell>
              <fwb-table-cell>{{ role.description || '-' }}</fwb-table-cell>
              <fwb-table-cell>
                <div class="flex flex-wrap gap-2">
                  <span
                    v-for="menu in role.menus || []"
                    :key="menu.id || menu.menu_key"
                    class="px-2 py-1 text-xs bg-blue-50 text-blue-600 rounded"
                  >
                    {{ menuLabelMap[menu.menu_key || menu.menuKey] || menu.menu_key || menu.menuKey }}
                  </span>
                  <span v-if="!role.menus || role.menus.length === 0" class="text-xs text-gray-400">No menus assigned</span>
                </div>
              </fwb-table-cell>
              <fwb-table-cell>
                <div class="flex flex-col gap-1">
                  <span
                    v-for="permission in rolePermissionSummaries(role.permissions)"
                    :key="permission"
                    class="text-xs text-gray-600"
                  >
                    {{ permission }}
                  </span>
                  <span v-if="!role.permissions || role.permissions.length === 0" class="text-xs text-gray-400">No permissions</span>
                </div>
              </fwb-table-cell>
              <fwb-table-cell>
                <div class="flex gap-2">
                  <button
                    class="text-blue-600 hover:text-blue-800"
                    @click="openRoleModal(role)"
                    title="Edit"
                  >
                    <i class="pi pi-pencil"></i>
                  </button>
                  <button
                    class="text-red-600 hover:text-red-800 disabled:text-gray-400 disabled:hover:text-gray-400"
                    @click="deleteRole(role)"
                    :disabled="role.name?.toLowerCase() === 'admin'"
                    title="Delete"
                  >
                    <i class="pi pi-trash"></i>
                  </button>
                </div>
              </fwb-table-cell>
            </fwb-table-row>
          </fwb-table-body>
        </fwb-table>

        <div v-if="roles.length === 0" class="text-center py-8 text-gray-500">
          No roles found.
        </div>
      </div>
    </div>

    <!-- Add User Modal -->
    <div
      v-if="showAddModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click.self="closeAddModal"
    >
      <div class="bg-white rounded-lg p-6 w-full max-w-md">
        <h3 class="text-xl font-bold mb-4">{{ isEditingUser ? 'Edit User' : 'Add New User' }}</h3>

        <form @submit.prevent="submitUser">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Full Name <span class="text-red-500">*</span>
              </label>
              <input
                v-model="newUser.full_name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Enter full name"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Email <span class="text-red-500">*</span>
              </label>
              <input
                v-model="newUser.email"
                type="email"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="user@example.com"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Phone (WhatsApp)
              </label>
              <input
                v-model="newUser.phone"
                type="tel"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="62812XXXXXXXX"
              />
              <p class="text-xs text-gray-500 mt-1">Format: 62812XXXXXXXX (without +)</p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Role <span class="text-red-500">*</span>
              </label>
              <select
                v-model="newUser.role_id"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="">Select Role</option>
                <option v-for="role in roles" :key="role.id" :value="role.id">
                  {{ role.name }}
                </option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Warehouse Access
              </label>
              <div class="border border-gray-200 rounded-lg divide-y max-h-48 overflow-y-auto">
                <label
                  v-for="warehouse in warehouseOptions"
                  :key="warehouse.id"
                  class="flex items-center gap-3 px-4 py-2 text-sm hover:bg-gray-50 cursor-pointer"
                >
                  <input
                    type="checkbox"
                    :value="String(warehouse.id)"
                    v-model="newUser.warehouse_ids"
                    class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  />
                  <div class="flex-1">
                    <p class="font-medium text-gray-700">{{ warehouse.name }}</p>
                    <p class="text-xs text-gray-500">Kode: {{ warehouse.code }}</p>
                  </div>
                </label>
                <div v-if="!warehouseOptions.length" class="px-4 py-3 text-sm text-gray-500">
                  Tidak ada data gudang tersedia.
                </div>
              </div>
              <p class="mt-1 text-xs text-gray-500">
                Cocok untuk role employee. Kosongkan bila ingin mencabut akses gudang.
              </p>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Password
              </label>
              <input
                v-model="newUser.password"
                type="password"
                :required="!isEditingUser"
                :minlength="!isEditingUser ? 6 : undefined"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                :placeholder="isEditingUser ? 'Optional - fill to reset password' : 'Minimum 6 characters'"
              />
            </div>

            <div class="flex items-center" v-if="!isEditingUser">
              <input
                v-model="newUser.send_welcome"
                type="checkbox"
                id="sendWelcome"
                :disabled="!newUser.phone"
                class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <label for="sendWelcome" class="ml-2 text-sm text-gray-700">
                Send welcome message via WhatsApp
                <span v-if="!newUser.phone" class="text-gray-400">(Phone required)</span>
              </label>
            </div>
          </div>

          <div class="flex justify-end gap-3 mt-6">
            <button
              type="button"
              @click="closeAddModal"
              class="px-4 py-2 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="isSubmitting"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              <span v-if="isSubmitting">{{ isEditingUser ? 'Saving...' : 'Creating...' }}</span>
              <span v-else>{{ isEditingUser ? 'Save Changes' : 'Create User' }}</span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Role Modal -->
    <div
      v-if="showRoleModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click.self="closeRoleModal"
    >
      <div class="bg-white rounded-lg p-6 w-full max-w-2xl max-h-[90vh] overflow-y-auto">
        <h3 class="text-xl font-bold mb-4">{{ roleForm.id ? 'Edit Role' : 'Add Role' }}</h3>

        <form @submit.prevent="submitRole">
          <div class="space-y-4">
            <div class="grid gap-4 md:grid-cols-2">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Role Name <span class="text-red-500">*</span>
                </label>
                <input
                  v-model="roleForm.name"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="e.g. supervisor"
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Description
                </label>
                <input
                  v-model="roleForm.description"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                  placeholder="Describe this role"
                />
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Badge Color
              </label>
              <div class="flex flex-wrap items-center gap-3">
                <span
                  class="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-semibold border"
                  :style="roleTagStyle({ color: roleForm.color })"
                >
                  <span
                    class="w-2.5 h-2.5 rounded-full"
                    :style="roleColorDotStyle({ color: roleForm.color })"
                  ></span>
                  {{ roleForm.color }}
                </span>
                <div class="flex flex-wrap gap-2">
                  <button
                    v-for="color in ROLE_COLOR_OPTIONS"
                    :key="color"
                    type="button"
                    class="w-10 h-10 rounded-full border-2 transition-transform"
                    :class="roleForm.color === resolveRoleColorValue(color) ? 'border-gray-900 scale-110' : 'border-transparent hover:scale-105'"
                    :style="{ backgroundColor: color }"
                    @click="roleForm.color = resolveRoleColorValue(color)"
                  ></button>
                </div>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Accessible Menus <span class="text-red-500">*</span>
              </label>
              <div class="grid gap-4 md:grid-cols-2">
                <div
                  v-for="group in menuOptionGroups"
                  :key="group.category"
                  class="p-4 border border-gray-200 rounded-lg"
                >
                  <p class="text-xs font-semibold uppercase text-gray-500 mb-3">
                    {{ categoryLabel(group.category) }}
                  </p>
                  <div class="space-y-2">
                    <label
                      v-for="option in group.items"
                      :key="option.key"
                      class="flex items-start gap-3 text-sm text-gray-700"
                    >
                      <input
                        type="checkbox"
                        class="mt-1 w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                        :value="option.key"
                        v-model="roleForm.menu_keys"
                      />
                      <span>{{ option.label }}</span>
                    </label>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <div class="flex items-center justify-between mb-2">
                <label class="block text-sm font-medium text-gray-700">
                  Permissions
                </label>
                <div v-if="permissionOptions.length" class="flex items-center gap-3">
                  <button
                    type="button"
                    class="text-xs text-blue-600 hover:underline"
                    @click="toggleAllPermissions(true)"
                  >
                    Select All
                  </button>
                  <button
                    type="button"
                    class="text-xs text-gray-500 hover:underline"
                    @click="toggleAllPermissions(false)"
                  >
                    Reset All
                  </button>
                </div>
              </div>
              <div v-if="permissionGroups.length" class="space-y-4 max-h-72 overflow-y-auto pr-1">
                <div
                  v-for="group in permissionGroups"
                  :key="group.module"
                  class="border border-gray-200 rounded-lg p-4"
                >
                  <div class="flex items-center justify-between mb-3">
                    <div>
                      <p class="text-sm font-semibold text-gray-800">{{ moduleLabel(group.module) }}</p>
                      <p class="text-xs text-gray-500">Kelola akses untuk modul ini</p>
                    </div>
                    <div class="flex items-center gap-2">
                      <button
                        type="button"
                        class="text-xs text-blue-600 hover:underline"
                        @click="toggleModulePermissions(group.module, true)"
                      >
                        Pilih semua
                      </button>
                      <button
                        type="button"
                        class="text-xs text-gray-500 hover:underline"
                        @click="toggleModulePermissions(group.module, false)"
                      >
                        Reset
                      </button>
                    </div>
                  </div>
                  <div class="grid grid-cols-2 gap-3">
                    <label
                      v-for="permission in group.permissions"
                      :key="permission.name"
                      class="flex items-center gap-2 text-sm text-gray-700"
                    >
                      <input
                        type="checkbox"
                        class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                        :value="permission.name"
                        v-model="roleForm.permission_keys"
                      />
                      <span class="capitalize">{{ actionLabel(permission.action) }}</span>
                    </label>
                  </div>
                </div>
              </div>
              <p v-else class="text-sm text-gray-500">
                Tidak ada data permission. Coba muat ulang halaman ini.
              </p>
            </div>
          </div>

          <div class="flex justify-end gap-3 mt-6">
            <button
              type="button"
              @click="closeRoleModal"
              class="px-4 py-2 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="isSavingRole"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed"
            >
              {{ isSavingRole ? 'Saving...' : roleForm.id ? 'Update Role' : 'Create Role' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { FwbTable, FwbTableBody, FwbTableCell, FwbTableHead, FwbTableHeadCell, FwbTableRow } from 'flowbite-vue'
import Tag from 'primevue/tag'
import { useToast } from 'primevue/usetoast'
import axios from '@/api/axios'

const toast = useToast()

const USERS_ENDPOINT = '/settings/users'
const ROLE_ENDPOINT = '/settings/roles'
const ROLES_LIST_ENDPOINT = '/roles'
const ROLE_MENU_OPTIONS_ENDPOINT = '/settings/roles/menu-options'
const ROLE_PERMISSION_OPTIONS_ENDPOINT = '/settings/roles/permission-options'

const activeTab = ref('users')

const users = ref([])
const roles = ref([])
const menuOptions = ref([])
const permissionOptions = ref([])
const warehouseOptions = ref([])

const showAddModal = ref(false)
const isSubmitting = ref(false)
const isEditingUser = ref(false)
const editingUserId = ref(null)
const showRoleModal = ref(false)
const isSavingRole = ref(false)

const newUser = ref({
  full_name: '',
  email: '',
  phone: '',
  password: '',
  role_id: '',
  send_welcome: false,
  warehouse_ids: []
})

const roleForm = ref({
  id: null,
  name: '',
  description: '',
  color: DEFAULT_ROLE_COLOR,
  menu_keys: [],
  permission_keys: []
})

const menuLabelMap = computed(() => {
  const map = {}
  menuOptions.value.forEach(option => {
    map[option.key] = option.label
  })
  return map
})

const menuOptionGroups = computed(() => {
  const groups = new Map()
  menuOptions.value
    .filter(option => option.key !== 'settings')
    .forEach(option => {
      const category = option.category || 'general'
      if (!groups.has(category)) {
        groups.set(category, [])
      }
      groups.get(category).push(option)
    })
  return Array.from(groups.entries()).map(([category, items]) => ({
    category,
    items
  }))
})

const capitalize = (value = '') => {
  if (!value || typeof value !== 'string') {
    return ''
  }
  return value.charAt(0).toUpperCase() + value.slice(1)
}

const MODULE_LABELS = {
  inventory: 'Inventory',
  warehouse: 'Warehouse',
  employee: 'Employees',
  lead: 'Leads',
  project: 'Projects',
  po: 'Purchase Orders'
}

const ACTION_LABELS = {
  view: 'View',
  create: 'Create',
  update: 'Update',
  delete: 'Delete',
  approve: 'Approve'
}

const ACTION_ORDER = ['view', 'create', 'update', 'delete', 'approve']

const ROLE_COLOR_OPTIONS = [
  '#EF4444', // red
  '#F97316', // orange
  '#F59E0B', // amber
  '#22C55E', // green
  '#0EA5E9', // sky
  '#6366F1', // indigo
  '#EC4899'  // pink
]

const DEFAULT_ROLE_COLOR = '#2563EB'
const NEUTRAL_ROLE_COLOR = '#64748B'
const HEX_COLOR_REGEX = /^#(?:[0-9A-Fa-f]{6}|[0-9A-Fa-f]{3})$/

const normalizeHexColor = (value = '') => {
  let color = (value || '').trim()
  if (!color) {
    return ''
  }
  if (!color.startsWith('#')) {
    color = `#${color}`
  }
  if (!HEX_COLOR_REGEX.test(color)) {
    return ''
  }
  if (color.length === 4) {
    color = `#${color
      .slice(1)
      .split('')
      .map((ch) => ch + ch)
      .join('')}`
  }
  return color.toUpperCase()
}

const resolveRoleColorValue = (value) => normalizeHexColor(value) || DEFAULT_ROLE_COLOR

const roleTextColor = (hexColor) => {
  const hex = normalizeHexColor(hexColor) || DEFAULT_ROLE_COLOR
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)
  const luminance = (0.299 * r + 0.587 * g + 0.114 * b) / 255
  return luminance > 0.65 ? '#1F2937' : '#FFFFFF'
}

const roleColorHex = (role) => {
  if (!role) return NEUTRAL_ROLE_COLOR
  return resolveRoleColorValue(role.color)
}

const roleTagStyle = (role) => {
  const hex = roleColorHex(role)
  return {
    backgroundColor: hex,
    color: roleTextColor(hex),
    border: 'none',
    boxShadow: '0 4px 10px rgba(15, 23, 42, 0.12)',
    fontWeight: 600
  }
}

const roleColorDotStyle = (role) => {
  const hex = roleColorHex(role)
  const outline = roleTextColor(hex) === '#FFFFFF' ? 'rgba(255, 255, 255, 0.6)' : 'rgba(15, 23, 42, 0.25)'
  return {
    backgroundColor: hex,
    boxShadow: `0 0 0 2px ${outline}`
  }
}

const moduleLabel = (module) => MODULE_LABELS[module] || capitalize(module)
const actionLabel = (action) => ACTION_LABELS[action] || capitalize(action)

const ensureWarehouseHex = (value) => {
  if (!value) return ''
  const color = value.trim()
  if (!color) return ''
  return color.startsWith('#') ? color.toUpperCase() : `#${color.toUpperCase()}`
}

const warehouseColor = (warehouse) => {
  if (!warehouse) return ''
  return ensureWarehouseHex(warehouse.color || warehouse.Color)
}

const warehouseTagStyle = (warehouse) => {
  const hex = warehouseColor(warehouse)
  if (!hex) {
    return {
      backgroundColor: '#E2E8F0',
      color: '#1F2937'
    }
  }
  return {
    backgroundColor: `${hex}1A`,
    color: hex,
    boxShadow: `0 4px 10px ${hex}26`
  }
}

const warehouseDotStyle = (warehouse) => {
  const hex = warehouseColor(warehouse) || '#1F2937'
  return {
    backgroundColor: hex,
    boxShadow: '0 0 0 2px rgba(255, 255, 255, 0.6)'
  }
}

const extractWarehouseIds = (warehouses = []) => {
  const ids = []
  if (!Array.isArray(warehouses)) {
    return ids
  }
  warehouses.forEach(entry => {
    if (!entry) return
    const candidate = entry.warehouse_id ?? entry.warehouse?.id
    const numeric = Number(candidate)
    if (!Number.isNaN(numeric) && numeric > 0) {
      const idString = String(numeric)
      if (!ids.includes(idString)) {
        ids.push(idString)
      }
    }
  })
  return ids
}

const withRoleColor = (role) => {
  if (!role) return role
  return {
    ...role,
    color: resolveRoleColorValue(role.color)
  }
}

const syncUserRoleColors = (availableRoles = []) => {
  if (!Array.isArray(availableRoles) || availableRoles.length === 0) {
    return
  }
  const roleColorMap = new Map(availableRoles.map(role => [role.id, role.color]))
  users.value = users.value.map(user => {
    const roleId = user.role?.id
    if (!roleId || !roleColorMap.has(roleId)) {
      return user
    }
    return {
      ...user,
      role: {
        ...user.role,
        color: roleColorMap.get(roleId)
      }
    }
  })
}

const permissionGroups = computed(() => {
  const groups = new Map()
  permissionOptions.value.forEach(permission => {
    if (!permission || !permission.module || !permission.name) {
      return
    }
    const moduleKey = permission.module
    if (!groups.has(moduleKey)) {
      groups.set(moduleKey, [])
    }
    groups.get(moduleKey).push(permission)
  })

  return Array.from(groups.entries())
    .map(([module, permissions]) => ({
      module,
      permissions: permissions
        .slice()
        .sort((a, b) => {
          const aIndex = ACTION_ORDER.indexOf(a.action)
          const bIndex = ACTION_ORDER.indexOf(b.action)
          const aScore = aIndex === -1 ? ACTION_ORDER.length : aIndex
          const bScore = bIndex === -1 ? ACTION_ORDER.length : bIndex
          return aScore - bScore
        })
    }))
    .sort((a, b) => moduleLabel(a.module).localeCompare(moduleLabel(b.module)))
})

const CATEGORY_LABELS = {
  main: 'Main Navigation',
  operations: 'Operations',
  support: 'Support',
  settings: 'Settings',
  general: 'General'
}

const categoryLabel = (category) => CATEGORY_LABELS[category] || category

const rolePermissionSummaries = (permissions = []) => {
  if (!Array.isArray(permissions) || permissions.length === 0) {
    return []
  }

  const grouped = new Map()
  permissions.forEach(permission => {
    if (!permission || !permission.module || !permission.action) {
      return
    }
    const moduleKey = permission.module
    if (!grouped.has(moduleKey)) {
      grouped.set(moduleKey, new Set())
    }
    grouped.get(moduleKey).add(actionLabel(permission.action))
  })

  return Array.from(grouped.entries())
    .map(([module, actions]) => `${moduleLabel(module)}: ${Array.from(actions).join(', ')}`)
    .sort((a, b) => a.localeCompare(b))
}

const toggleModulePermissions = (module, selectAll) => {
  const group = permissionGroups.value.find(entry => entry.module === module)
  if (!group) {
    return
  }

  const moduleKeys = group.permissions.map(permission => permission.name)
  if (selectAll) {
    const next = new Set(roleForm.value.permission_keys)
    moduleKeys.forEach(key => next.add(key))
    roleForm.value.permission_keys = Array.from(next)
  } else {
    roleForm.value.permission_keys = roleForm.value.permission_keys.filter(key => !moduleKeys.includes(key))
  }
}

const toggleAllPermissions = (selectAll) => {
  if (selectAll) {
    const next = new Set(roleForm.value.permission_keys)
    permissionOptions.value.forEach(permission => {
      if (permission?.name) {
        next.add(permission.name)
      }
    })
    roleForm.value.permission_keys = Array.from(next)
  } else {
    roleForm.value.permission_keys = []
  }
}

const loadUsers = async () => {
  try {
    const response = await axios.get(USERS_ENDPOINT)
    const fetchedUsers = (response.data.data || []).map(user => ({
      ...user,
      role: withRoleColor(user.role)
    }))
    users.value = fetchedUsers
    syncUserRoleColors(roles.value)
  } catch (error) {
    console.error('Error loading users:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Gagal memuat data user',
      life: 3000
    })
  }
}

const loadRoles = async () => {
  try {
    const response = await axios.get(ROLES_LIST_ENDPOINT)
    const fetchedRoles = (response.data.data || []).map(role => withRoleColor(role))
    roles.value = fetchedRoles
    syncUserRoleColors(fetchedRoles)
  } catch (error) {
    console.error('Error loading roles:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Gagal memuat data role',
      life: 3000
    })
  }
}

const loadMenuOptions = async () => {
  try {
    const response = await axios.get(ROLE_MENU_OPTIONS_ENDPOINT)
    menuOptions.value = response.data.data || []
  } catch (error) {
    console.error('Error loading menu options:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Gagal memuat daftar menu',
      life: 3000
    })
  }
}

const loadPermissionOptions = async () => {
  try {
    const response = await axios.get(ROLE_PERMISSION_OPTIONS_ENDPOINT)
    permissionOptions.value = response.data.data || []
  } catch (error) {
    console.error('Error loading permission options:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Gagal memuat daftar permission',
      life: 3000
    })
  }
}

const loadWarehouses = async () => {
  try {
    const response = await axios.get('/warehouses')
    warehouseOptions.value = response.data.data || []
  } catch (error) {
    console.error('Error loading warehouses:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Gagal memuat daftar gudang',
      life: 3000
    })
  }
}

const openAddModal = () => {
  isEditingUser.value = false
  editingUserId.value = null
  newUser.value = {
    full_name: '',
    email: '',
    phone: '',
    password: '',
    role_id: '',
    send_welcome: false,
    warehouse_ids: []
  }
  showAddModal.value = true
}

const closeAddModal = () => {
  isEditingUser.value = false
  editingUserId.value = null
  showAddModal.value = false
}

const openEditModal = (user) => {
  isEditingUser.value = true
  editingUserId.value = user.id
  newUser.value = {
    full_name: user.full_name,
    email: user.email,
    phone: user.phone || '',
    password: '',
    role_id: user.role_id,
    send_welcome: false,
    warehouse_ids: extractWarehouseIds(user.warehouses)
  }
  showAddModal.value = true
}

const submitUser = async () => {
  if (isSubmitting.value) return

  const payload = {
    full_name: newUser.value.full_name.trim(),
    email: newUser.value.email.trim(),
    phone: newUser.value.phone.trim(),
    role_id: newUser.value.role_id
  }

  const selectedWarehouseIds = Array.isArray(newUser.value.warehouse_ids)
    ? newUser.value.warehouse_ids
        .map(id => Number(id))
        .filter(id => !Number.isNaN(id) && id > 0)
    : []
  payload.warehouse_ids = Array.from(new Set(selectedWarehouseIds))

  if (!payload.full_name || !payload.email || !payload.role_id) {
    toast.add({
      severity: 'warn',
      summary: 'Validasi',
      detail: 'Lengkapi nama, email, dan role user',
      life: 3000
    })
    return
  }

  const isPasswordProvided = newUser.value.password && newUser.value.password.trim().length > 0
  if (!isEditingUser.value) {
    if (!isPasswordProvided || newUser.value.password.length < 6) {
      toast.add({
        severity: 'warn',
        summary: 'Validasi',
        detail: 'Password minimal 6 karakter',
        life: 3000
      })
      return
    }
  } else if (isPasswordProvided && newUser.value.password.length < 6) {
    toast.add({
      severity: 'warn',
      summary: 'Validasi',
      detail: 'Password minimal 6 karakter',
      life: 3000
    })
    return
  }

  isSubmitting.value = true

  try {
    if (isEditingUser.value && editingUserId.value) {
      const updatePayload = {
        ...payload
      }
      if (isPasswordProvided) {
        updatePayload.password = newUser.value.password
      }

      const response = await axios.put(`${USERS_ENDPOINT}/${editingUserId.value}`, updatePayload)

      toast.add({
        severity: 'success',
        summary: 'Berhasil',
        detail: response.data.message || 'User berhasil diperbarui',
        life: 3000
      })
    } else {
      const createPayload = {
        ...payload,
        password: newUser.value.password,
        send_welcome: newUser.value.send_welcome
      }

      const response = await axios.post(USERS_ENDPOINT, createPayload)

      toast.add({
        severity: 'success',
        summary: 'Berhasil',
        detail: response.data.message || 'User berhasil ditambahkan',
        life: 3000
      })
    }

    closeAddModal()
    await loadUsers()
  } catch (error) {
    console.error('Error saving user:', error)

    let errorMessage = 'Gagal menyimpan user'
    if (error.response?.data?.error) {
      errorMessage = error.response.data.error
    } else if (error.response?.status === 400) {
      errorMessage = 'Data tidak valid atau email sudah terdaftar'
    } else if (error.response?.status === 500) {
      errorMessage = 'Terjadi kesalahan server'
    }

    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: errorMessage,
      life: 5000
    })
  } finally {
    isSubmitting.value = false
  }
}

const toggleUserStatus = async (user) => {
  try {
    const newStatus = !user.is_active
    await axios.put(`${USERS_ENDPOINT}/${user.id}/status`, { is_active: newStatus })
    
    toast.add({
      severity: 'success',
      summary: 'Berhasil',
      detail: `User ${newStatus ? 'diaktifkan' : 'dinonaktifkan'} berhasil`,
      life: 3000
    })
    
    await loadUsers()
  } catch (error) {
    console.error('Error toggling user status:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.error || 'Gagal update status user',
      life: 3000
    })
  }
}

const deleteUser = async (user) => {
  if (!confirm(`Apakah Anda yakin ingin menghapus user "${user.full_name}"?`)) {
    return
  }
  
  try {
    await axios.delete(`${USERS_ENDPOINT}/${user.id}`)
    
    toast.add({
      severity: 'success',
      summary: 'Berhasil',
      detail: 'User berhasil dihapus',
      life: 3000
    })
    
    await loadUsers()
  } catch (error) {
    console.error('Error deleting user:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.error || 'Gagal menghapus user',
      life: 3000
    })
  }
}

const resetRoleForm = () => {
  roleForm.value = {
    id: null,
    name: '',
    description: '',
    color: DEFAULT_ROLE_COLOR,
    menu_keys: [],
    permission_keys: []
  }
}

const openRoleModal = (role = null) => {
  if (menuOptions.value.length === 0) {
    loadMenuOptions()
  }
  if (permissionOptions.value.length === 0) {
    loadPermissionOptions()
  }

  if (role) {
    const normalizedRole = withRoleColor(role)
    roleForm.value = {
      id: normalizedRole.id,
      name: normalizedRole.name,
      description: normalizedRole.description || '',
      color: normalizedRole.color,
      menu_keys: (normalizedRole.menus || []).map(menu => menu.menu_key || menu.menuKey),
      permission_keys: (normalizedRole.permissions || []).map(permission => permission.name).filter(Boolean)
    }
  } else {
    resetRoleForm()
  }
  showRoleModal.value = true
}

const closeRoleModal = () => {
  showRoleModal.value = false
  resetRoleForm()
}

const submitRole = async () => {
  if (isSavingRole.value) return

  if (!roleForm.value.name.trim()) {
    toast.add({
      severity: 'warn',
      summary: 'Validasi',
      detail: 'Nama role wajib diisi',
      life: 3000
    })
    return
  }

  if (!roleForm.value.menu_keys.length) {
    toast.add({
      severity: 'warn',
      summary: 'Validasi',
      detail: 'Pilih minimal satu menu yang dapat diakses',
      life: 3000
    })
    return
  }

  isSavingRole.value = true

  try {
    const payload = {
      name: roleForm.value.name.trim(),
      description: roleForm.value.description.trim(),
      color: resolveRoleColorValue(roleForm.value.color),
      menu_keys: [...roleForm.value.menu_keys],
      permission_keys: [...roleForm.value.permission_keys]
    }

    if (roleForm.value.id) {
      await axios.put(`${ROLE_ENDPOINT}/${roleForm.value.id}`, payload)
      toast.add({
        severity: 'success',
        summary: 'Berhasil',
        detail: 'Role berhasil diperbarui',
        life: 3000
      })
    } else {
      await axios.post(ROLE_ENDPOINT, payload)
      toast.add({
        severity: 'success',
        summary: 'Berhasil',
        detail: 'Role baru berhasil dibuat',
        life: 3000
      })
    }

    closeRoleModal()
    await loadRoles()
  } catch (error) {
    console.error('Error saving role:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.error || 'Gagal menyimpan role',
      life: 4000
    })
  } finally {
    isSavingRole.value = false
  }
}

const deleteRole = async (role) => {
  if (!confirm(`Hapus role "${role.name}"?`)) {
    return
  }

  try {
    await axios.delete(`${ROLE_ENDPOINT}/${role.id}`)
    toast.add({
      severity: 'success',
      summary: 'Berhasil',
      detail: 'Role berhasil dihapus',
      life: 3000
    })
    await loadRoles()
  } catch (error) {
    console.error('Error deleting role:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.error || 'Gagal menghapus role',
      life: 4000
    })
  }
}

onMounted(() => {
  loadUsers()
  loadRoles()
  loadMenuOptions()
  loadPermissionOptions()
  loadWarehouses()
})
</script>
