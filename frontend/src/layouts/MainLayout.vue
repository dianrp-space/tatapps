<template>
  <div class="min-h-screen bg-gray-100">
    <!-- Mobile overlay -->
    <div
      v-if="isMobileSidebarOpen"
      class="fixed inset-0 bg-black bg-opacity-40 z-30 lg:hidden"
      @click="closeMobileSidebar"
    ></div>

    <!-- Sidebar -->
    <aside
      :class="[
        'fixed inset-y-0 left-0 bg-white shadow-lg overflow-y-auto z-40 transition-all duration-300',
        isSidebarCollapsed ? 'w-20' : 'w-64',
        isMobileSidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'
      ]"
    >
      <div class="flex items-center h-16 border-b px-4">
        <button
          class="lg:hidden mr-2 text-gray-500 hover:text-gray-700"
          @click="closeMobileSidebar"
        >
          <i class="pi pi-times text-xl"></i>
        </button>
        <div class="flex items-center gap-3 max-w-full">
          <img
            v-if="siteLogo"
            :src="siteLogo"
            alt="App logo"
            :class="[isSidebarCollapsed ? 'h-10 w-10 mx-auto object-contain' : 'h-10 w-auto object-contain']"
          />
          <h1
            v-if="!isSidebarCollapsed"
            class="text-2xl font-bold text-blue-600 truncate"
          >
            {{ appName }}
          </h1>
        </div>
      </div>
      
      <nav class="p-4">
        <template v-for="item in menuItems" :key="item.path">
          <!-- Menu with submenu -->
          <div v-if="item.children" class="mb-2">
            <button
              @click="toggleSubmenu(item.path)"
              class="flex items-center w-full px-4 py-3 rounded-lg hover:bg-blue-50 transition-colors"
              :class="{ 'bg-blue-50 text-blue-600': isActiveMenu(item) }"
            >
              <div
                class="flex items-center gap-3"
                :class="isSidebarCollapsed ? 'justify-center w-full' : ''"
              >
                <i :class="item.icon"></i>
                <span v-if="!isSidebarCollapsed">{{ item.label }}</span>
              </div>
              <i 
                v-if="!isSidebarCollapsed"
                class="pi transition-transform duration-200"
                :class="openSubmenus.includes(item.path) ? 'pi-chevron-down' : 'pi-chevron-right'"
              ></i>
            </button>
            
            <!-- Submenu -->
            <div 
              v-if="!isSidebarCollapsed"
              v-show="openSubmenus.includes(item.path)"
              class="ml-4 mt-1 space-y-1"
            >
              <router-link
                v-for="child in item.children"
                :key="child.path"
                :to="child.path"
                class="flex items-center gap-3 px-4 py-2 rounded-lg hover:bg-gray-50 transition-colors text-sm"
                active-class="bg-blue-50 text-blue-600 font-medium"
                @click="closeMobileSidebar"
              >
                <i :class="child.icon" class="text-xs"></i>
                <span>{{ child.label }}</span>
              </router-link>
            </div>
          </div>
          
          <!-- Regular menu -->
          <router-link
            v-else
            :to="item.path"
            class="flex items-center gap-3 px-4 py-3 mb-2 rounded-lg hover:bg-blue-50 transition-colors"
            :class="isSidebarCollapsed ? 'justify-center' : ''"
            active-class="bg-blue-100 text-blue-600"
            @click="closeMobileSidebar"
          >
            <i :class="item.icon"></i>
            <span v-if="!isSidebarCollapsed">{{ item.label }}</span>
          </router-link>
        </template>
      </nav>
    </aside>

    <!-- Main Content -->
    <div :class="[isSidebarCollapsed ? 'lg:ml-20' : 'lg:ml-64', 'transition-all duration-300 ml-0']">
      <!-- Header -->
      <header class="bg-white shadow-sm h-16 flex items-center justify-between px-4 lg:px-6">
        <div class="flex items-center gap-3">
          <button
            class="p-2 rounded-full hover:bg-gray-100 lg:hidden"
            @click="openMobileSidebar"
          >
            <i class="pi pi-bars text-gray-600"></i>
          </button>
          <button
            class="hidden lg:inline-flex p-2 rounded-full hover:bg-gray-100"
            @click="toggleCollapse"
          >
            <i :class="isSidebarCollapsed ? 'pi pi-angle-double-right' : 'pi pi-angle-double-left'" class="text-gray-600"></i>
          </button>
          <div class="flex flex-col">
            <nav
              v-if="breadcrumbs.length"
              class="flex items-center text-xs sm:text-sm text-gray-500 gap-2"
              aria-label="Breadcrumb"
            >
              <template v-for="(crumb, index) in breadcrumbs" :key="`${crumb.label}-${index}`">
                <router-link
                  v-if="crumb.to"
                  :to="crumb.to"
                  class="hover:text-gray-700 transition-colors"
                >
                  {{ crumb.label }}
                </router-link>
                <span
                  v-else
                  class="text-gray-700 font-medium"
                >
                  {{ crumb.label }}
                </span>
                <i
                  v-if="index < breadcrumbs.length - 1"
                  class="pi pi-chevron-right text-[10px] sm:text-xs text-gray-400"
                ></i>
              </template>
            </nav>
            <h2 class="text-xl font-semibold text-gray-800 mt-1">{{ currentPageTitle }}</h2>
          </div>
        </div>
        
        <div class="flex items-center gap-4">
          <button class="p-2 rounded-full hover:bg-gray-100">
            <i class="pi pi-bell text-gray-600"></i>
          </button>
          
          <div class="flex items-center gap-3">
            <div ref="profileMenuRef" class="relative">
              <button
                type="button"
                @click.stop="toggleProfileMenu"
                class="flex items-center gap-3 p-2 rounded-full hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <img
                  :src="avatarUrl"
                  alt="Avatar"
                  class="w-10 h-10 rounded-full object-cover"
                />
                <div class="text-left">
                  <p class="text-sm font-medium text-gray-800">{{ authStore.user?.full_name }}</p>
                  <p class="text-xs text-gray-500">{{ authStore.user?.role?.name }}</p>
                </div>
                <i class="pi pi-chevron-down text-gray-500 text-sm"></i>
              </button>
              <div
                v-if="isProfileMenuOpen"
                class="absolute right-0 mt-2 w-44 bg-white border border-gray-200 rounded-lg shadow-lg py-2 z-50"
              >
                <router-link
                  to="/settings/profile"
                  class="flex items-center gap-2 px-4 py-2 text-sm text-gray-700 hover:bg-gray-50"
                  @click="closeProfileMenu"
                >
                  <i class="pi pi-user text-gray-500"></i>
                  <span>Profile Settings</span>
                </router-link>
                <button
                  type="button"
                  class="w-full flex items-center gap-2 px-4 py-2 text-sm text-left text-red-600 hover:bg-red-50"
                  @click="handleLogoutFromMenu"
                >
                  <i class="pi pi-sign-out"></i>
                  <span>Logout</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- Page Content -->
      <main class="p-6">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useSiteStore } from '@/stores/site'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const siteStore = useSiteStore()

const MENU_ITEMS = [
  { path: '/', label: 'Dashboard', icon: 'pi pi-home', key: 'dashboard' },
  { path: '/warehouses', label: 'Warehouses', icon: 'pi pi-building', key: 'warehouses' },
  {
    path: '/inventory',
    label: 'Inventory',
    icon: 'pi pi-box',
    key: 'inventory',
    children: [
      { path: '/inventory', label: 'Items', icon: 'pi pi-list', key: 'inventory.items' },
      { path: '/inventory/categories', label: 'Categories', icon: 'pi pi-tags', key: 'inventory.categories' },
      { path: '/inventory/transactions', label: 'Transactions', icon: 'pi pi-history', key: 'inventory.transactions' }
    ]
  },
  {
    path: '/employees',
    label: 'Employees',
    icon: 'pi pi-users',
    key: 'employees',
    children: [
      { path: '/employees', label: 'Data Karyawan', icon: 'pi pi-users', key: 'employees.data' },
      { path: '/employees/divisions', label: 'Divisi', icon: 'pi pi-sitemap', key: 'employees.divisions' },
      { path: '/employees/positions', label: 'Jabatan', icon: 'pi pi-briefcase', key: 'employees.positions' }
    ]
  },
  { path: '/leads', label: 'Leads', icon: 'pi pi-user-plus', key: 'leads' },
  { path: '/projects', label: 'Projects', icon: 'pi pi-briefcase', key: 'projects' },
  { path: '/purchase-orders', label: 'Purchase Orders', icon: 'pi pi-shopping-cart', key: 'purchase_orders' },
  { path: '/support', label: 'Support', icon: 'pi pi-question', key: 'support' },
  {
    path: '/settings',
    label: 'Settings',
    icon: 'pi pi-cog',
    key: 'settings',
    children: [
      { path: '/settings/profile', label: 'Profile', icon: 'pi pi-user', key: 'settings.profile' },
      { path: '/settings/company', label: 'Profil Perusahaan', icon: 'pi pi-building', key: 'settings.company' },
      { path: '/settings/notifications', label: 'Notifications', icon: 'pi pi-bell', key: 'settings.notifications' },
      { path: '/settings/users', label: 'User Management', icon: 'pi pi-users', key: 'settings.users' },
      { path: '/settings/sites', label: 'Sites', icon: 'pi pi-globe', key: 'settings.sites' }
    ]
  }
]

const openSubmenus = ref(['/inventory', '/employees'])
const isSidebarCollapsed = ref(false)
const isMobileSidebarOpen = ref(false)
const isProfileMenuOpen = ref(false)
const profileMenuRef = ref(null)

// Computed avatar URL
const avatarUrl = computed(() => {
  if (!authStore.user?.avatar) {
    return 'https://via.placeholder.com/40'
  }
  if (authStore.user.avatar.startsWith('http')) {
    return authStore.user.avatar
  }
  return `http://localhost:8080/${authStore.user.avatar}`
})

const roleMenuKeys = computed(() => {
  const menus = authStore.user?.role?.menus || []
  return menus
    .map(menu => menu.menu_key || menu.menuKey)
    .filter(key => typeof key === 'string' && key.length > 0)
})

const allowedMenuSet = computed(() => new Set(roleMenuKeys.value))
const hasMenuRestrictions = computed(() => allowedMenuSet.value.size > 0)

function canAccessMenu(key) {
  if (!key) return true
  if (!hasMenuRestrictions.value) return true
  if (allowedMenuSet.value.has(key)) return true
  const segments = key.split('.')
  while (segments.length > 1) {
    segments.pop()
    const candidate = segments.join('.')
    if (allowedMenuSet.value.has(candidate)) {
      return true
    }
  }
  return false
}

function filterMenuTree(items) {
  const filtered = []
  items.forEach(item => {
    const entry = { ...item }
    if (entry.children && entry.children.length > 0) {
      const childItems = filterMenuTree(entry.children)
      const visible = canAccessMenu(entry.key) || childItems.length > 0
      if (visible) {
        entry.children = childItems
        filtered.push(entry)
      }
    } else if (canAccessMenu(entry.key)) {
      filtered.push(entry)
    }
  })
  return filtered
}

const menuItems = computed(() => filterMenuTree(MENU_ITEMS))

const appName = computed(() => siteStore.settings.app_name || 'TatApps')
const siteLogo = computed(() => siteStore.logoUrl())

const currentPageTitle = computed(() => {
  return route.meta.title || formatFromName(route.name) || 'Dashboard'
})

function formatFromName(name) {
  if (!name) return ''
  return name
    .toString()
    .replace(/([a-z0-9])([A-Z])/g, '$1 $2')
    .replace(/[_-]/g, ' ')
}

function normalizeRoutePath(path) {
  if (!path) return '/'
  if (path.startsWith('http')) return path
  let normalized = path.replace(/\/+/g, '/')
  if (!normalized.startsWith('/')) {
    normalized = `/${normalized}`
  }
  return normalized || '/'
}

function resolveRecordLink(record) {
  if (!record) return null
  if (record.meta?.disableBreadcrumbLink) return null
  if (record.name) {
    try {
      const location = router.resolve({
        name: record.name,
        params: route.params
      })
      return normalizeRoutePath(location.path)
    } catch {
      // Ignore resolution errors and fall back to static path
    }
  }
  if (record.path) {
    let target = record.path
    Object.entries(route.params).forEach(([key, value]) => {
      const pattern = new RegExp(`:${key}\\??`, 'g')
      target = target.replace(pattern, value ?? '')
    })
    return normalizeRoutePath(target.replace(/(:[^/]+)/g, ''))
  }
  return null
}

function resolveParentCrumbs(record) {
  const meta = record.meta || {}
  const parentsMeta = typeof meta.breadcrumbParents === 'function'
    ? meta.breadcrumbParents(route)
    : meta.breadcrumbParents || []
  const parentsArray = Array.isArray(parentsMeta) ? parentsMeta : [parentsMeta]
  return parentsArray
    .map(parent => {
      if (!parent) return null
      if (typeof parent === 'function') {
        return parent(route)
      }
      if (typeof parent === 'string') {
        return { label: parent, to: '/' }
      }
      return {
        label: parent.label || parent.title || '',
        to: parent.to || parent.path || null
      }
    })
    .filter(parent => parent && parent.label)
    .map(parent => ({
      label: parent.label,
      to: parent.to ? normalizeRoutePath(parent.to) : null
    }))
}

function resolveBreadcrumbLabel(record) {
  const meta = record.meta || {}
  if (typeof meta.breadcrumb === 'function') {
    return meta.breadcrumb(route)
  }
  if (typeof meta.breadcrumb === 'string') {
    return meta.breadcrumb
  }
  if (meta.title) {
    return meta.title
  }
  return formatFromName(record.name)
}

function addCrumb(list, crumb) {
  if (!crumb || !crumb.label) return
  const existingIndex = list.findIndex(
    item => item.label === crumb.label && (item.to || '') === (crumb.to || '')
  )
  if (existingIndex === -1) {
    list.push({ label: crumb.label, to: crumb.to || null })
  }
}

const breadcrumbs = computed(() => {
  const items = []
  const matchedRecords = route.matched.filter(record => record.meta?.breadcrumb !== false)

  matchedRecords.forEach((record, index) => {
    const parents = resolveParentCrumbs(record)
    parents.forEach(parent => addCrumb(items, parent))

    const label = resolveBreadcrumbLabel(record)
    if (!label) return

    const isLast = index === matchedRecords.length - 1
    const link = isLast ? null : resolveRecordLink(record)
    addCrumb(items, { label, to: link })
  })

  if (!items.length) {
    return [{ label: 'Dashboard', to: null }]
  }

  if (items[0].to !== '/' && route.path !== '/') {
    items.unshift({ label: 'Dashboard', to: '/' })
  }

  items[items.length - 1].to = null
  return items
})

// Auto-expand submenu when navigating to a child route
watch(() => route.path, (newPath) => {
  menuItems.value.forEach(item => {
    if (item.children) {
      const hasActiveChild = item.children.some(child => newPath.startsWith(child.path))
      if (hasActiveChild && !openSubmenus.value.includes(item.path)) {
        openSubmenus.value.push(item.path)
      }
    }
  })
}, { immediate: true })

function toggleSubmenu(path) {
  const index = openSubmenus.value.indexOf(path)
  if (index > -1) {
    openSubmenus.value.splice(index, 1)
  } else {
    openSubmenus.value.push(path)
  }
}

function isActiveMenu(item) {
  if (item.children) {
    return item.children.some(child => route.path.startsWith(child.path))
  }
  return route.path === item.path
}

function toggleCollapse() {
  isSidebarCollapsed.value = !isSidebarCollapsed.value
}

function openMobileSidebar() {
  isMobileSidebarOpen.value = true
}

function closeMobileSidebar() {
  isMobileSidebarOpen.value = false
}

function handleLogout() {
  closeProfileMenu()
  authStore.logout()
  router.push('/login')
}

function handleLogoutFromMenu() {
  handleLogout()
}

function toggleProfileMenu() {
  isProfileMenuOpen.value = !isProfileMenuOpen.value
}

function closeProfileMenu() {
  isProfileMenuOpen.value = false
}

function handleClickOutside(event) {
  if (profileMenuRef.value && !profileMenuRef.value.contains(event.target)) {
    closeProfileMenu()
  }
}

onMounted(() => {
  siteStore.fetchSiteSettings().catch(() => {
    // Ignore errors; UI will fallback to defaults
  })
  document.addEventListener('click', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', handleClickOutside)
})

watch(() => route.path, () => {
  isMobileSidebarOpen.value = false
  closeProfileMenu()
})
</script>
