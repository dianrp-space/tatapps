import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: { guest: true }
  },
  {
    path: '/error/forbidden',
    name: 'ErrorForbidden',
    component: () => import('@/views/error/Forbidden.vue'),
    meta: {
      title: 'Akses Ditolak'
    }
  },
  {
    path: '/error/server',
    name: 'ErrorServer',
    component: () => import('@/views/error/ServerError.vue'),
    meta: {
      title: 'Kesalahan Sistem'
    }
  },
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: {
          title: 'Dashboard',
          breadcrumb: 'Dashboard'
        }
      },
      {
        path: '/warehouses',
        name: 'Warehouses',
        component: () => import('@/views/warehouse/List.vue'),
        meta: {
          title: 'Warehouses',
          breadcrumb: 'Warehouses'
        }
      },
      {
        path: '/warehouses/:id',
        name: 'WarehouseDetail',
        component: () => import('@/views/warehouse/Detail.vue'),
        meta: {
          title: 'Warehouse Detail',
          breadcrumb: route => `Warehouse ${route.params.id || ''}`.trim(),
          breadcrumbParents: [
            { label: 'Warehouses', to: '/warehouses' }
          ]
        }
      },
      {
        path: '/inventory',
        name: 'Inventory',
        component: () => import('@/views/inventory/List.vue'),
        meta: {
          title: 'Inventory',
          breadcrumb: 'Inventory'
        }
      },
      {
        path: '/inventory/:id',
        name: 'InventoryDetail',
        component: () => import('@/views/inventory/Detail.vue'),
        meta: {
          title: 'Inventory Detail',
          breadcrumb: route => `Item ${route.params.id || ''}`.trim(),
          breadcrumbParents: [
            { label: 'Inventory', to: '/inventory' }
          ]
        }
      },
      {
        path: '/inventory/categories',
        name: 'InventoryCategories',
        component: () => import('@/views/inventory/Categories.vue'),
        meta: {
          title: 'Categories',
          breadcrumb: 'Categories',
          permissions: ['category.view'],
          breadcrumbParents: [
            { label: 'Inventory', to: '/inventory' }
          ]
        }
      },
      {
        path: '/inventory/transactions',
        name: 'InventoryTransactions',
        component: () => import('@/views/inventory/Transactions.vue'),
        meta: {
          title: 'Transactions',
          breadcrumb: 'Transactions',
          breadcrumbParents: [
            { label: 'Inventory', to: '/inventory' }
          ]
        }
      },
      {
        path: '/employees',
        name: 'EmployeeData',
        component: () => import('@/views/employee/List.vue'),
        meta: {
          title: 'Data Karyawan',
          breadcrumb: 'Data Karyawan',
          permissions: ['employee.view']
        }
      },
      {
        path: '/employees/create',
        name: 'EmployeeCreate',
        component: () => import('@/views/employee/Create.vue'),
        meta: {
          title: 'Tambah Karyawan',
          breadcrumb: 'Tambah Karyawan',
          permissions: ['employee.create'],
          breadcrumbParents: [
            { label: 'Data Karyawan', to: '/employees' }
          ]
        }
      },
      {
        path: '/employees/:id/edit',
        name: 'EmployeeEdit',
        component: () => import('@/views/employee/Edit.vue'),
        meta: {
          title: 'Edit Karyawan',
          breadcrumb: 'Edit Karyawan',
          permissions: ['employee.update'],
          breadcrumbParents: [
            { label: 'Data Karyawan', to: '/employees' }
          ]
        }
      },
      {
        path: '/employees/divisions',
        name: 'EmployeeDivisions',
        component: () => import('@/views/employee/Divisions.vue'),
        meta: {
          title: 'Divisi',
          breadcrumb: 'Divisi',
          permissions: ['employee.view'],
          breadcrumbParents: [
            { label: 'Data Karyawan', to: '/employees' }
          ]
        }
      },
      {
        path: '/employees/positions',
        name: 'EmployeePositions',
        component: () => import('@/views/employee/Positions.vue'),
        meta: {
          title: 'Jabatan',
          breadcrumb: 'Jabatan',
          permissions: ['employee.view'],
          breadcrumbParents: [
            { label: 'Data Karyawan', to: '/employees' }
          ]
        }
      },
      {
        path: '/employees/positions/hierarchy',
        name: 'EmployeeHierarchy',
        component: () => import('@/views/employee/Hierarchy.vue'),
        meta: {
          title: 'Hierarki Jabatan',
          breadcrumb: 'Hierarki',
          permissions: ['employee.view'],
          breadcrumbParents: [
            { label: 'Data Karyawan', to: '/employees' },
            { label: 'Jabatan', to: '/employees/positions' }
          ]
        }
      },
      {
        path: '/leads',
        name: 'Leads',
        component: () => import('@/views/lead/List.vue'),
        meta: {
          title: 'Leads',
          breadcrumb: 'Leads'
        }
      },
      {
        path: '/projects',
        name: 'Projects',
        component: () => import('@/views/project/List.vue'),
        meta: {
          title: 'Projects',
          breadcrumb: 'Projects'
        }
      },
      {
        path: '/purchase-orders',
        name: 'PurchaseOrders',
        component: () => import('@/views/po/List.vue'),
        meta: {
          title: 'Purchase Orders',
          breadcrumb: 'Purchase Orders'
        }
      },
      {
        path: '/purchase-orders/:id',
        name: 'PODetail',
        component: () => import('@/views/po/Detail.vue'),
        meta: {
          title: 'Purchase Order Detail',
          breadcrumb: route => `PO ${route.params.id || ''}`.trim(),
          breadcrumbParents: [
            { label: 'Purchase Orders', to: '/purchase-orders' }
          ]
        }
      },
      {
        path: '/support',
        name: 'Support',
        component: () => import('@/views/support/Dashboard.vue'),
        meta: {
          title: 'Support',
          breadcrumb: 'Support'
        }
      },
      {
        path: '/settings/profile',
        name: 'Profile',
        component: () => import('@/views/settings/Profile.vue'),
        meta: {
          title: 'Profile',
          breadcrumb: 'Profile'
        }
      },
      {
        path: '/settings/company',
        name: 'CompanyProfile',
        component: () => import('@/views/settings/Company.vue'),
        meta: {
          title: 'Profil Perusahaan',
          breadcrumb: 'Profil Perusahaan',
          breadcrumbParents: [
            { label: 'Settings', to: '/settings/profile' }
          ]
        }
      },
      {
        path: '/settings/notifications',
        name: 'NotificationSettings',
        component: () => import('@/views/settings/Notifications.vue'),
        meta: {
          title: 'Notification Settings',
          breadcrumb: 'Notifications',
          breadcrumbParents: [
            { label: 'Settings', to: '/settings/profile' }
          ]
        }
      },
      {
        path: '/settings/sites',
        name: 'SiteSettings',
        component: () => import('@/views/settings/Sites.vue'),
        meta: {
          title: 'Site Settings',
          breadcrumb: 'Sites',
          breadcrumbParents: [
            { label: 'Settings', to: '/settings/profile' }
          ]
        }
      },
      {
        path: '/settings/users',
        name: 'UserManagement',
        component: () => import('@/views/users/List.vue'),
        meta: {
          title: 'User Management',
          breadcrumb: 'Users',
          breadcrumbParents: [
            { label: 'Settings', to: '/settings/profile' }
          ]
        }
      },
      // Redirect old profile route
      {
        path: '/profile',
        redirect: '/settings/profile'
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/error/NotFound.vue'),
    meta: {
      title: 'Halaman Tidak Ditemukan'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guard
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  const isAuthenticated = authStore.isAuthenticated

  // If authenticated but no user data, try to fetch profile
  if (isAuthenticated && !authStore.user) {
    try {
      await authStore.getProfile()
    } catch (error) {
      // If profile fetch fails, logout and redirect to login
      authStore.logout()
      if (to.meta.requiresAuth) {
        return next('/login')
      }
    }
  }

  if (to.meta.requiresAuth && !isAuthenticated) {
    next('/login')
    return
  }

  if (to.meta.guest && isAuthenticated) {
    next('/')
    return
  }

  const requiresPermissions = to.matched.some(record => {
    const meta = record.meta || {}
    const candidates = [
      meta.permissions,
      meta.allPermissions,
      meta.anyPermissions,
      meta.anyPermission
    ]
    return candidates.some(value => {
      if (!value) return false
      return Array.isArray(value) ? value.length > 0 : true
    })
  })

  if (requiresPermissions) {
    const allPermissions = new Set()
    const anyPermissions = new Set()
    const normalizePermissions = value => {
      if (!value) return []
      return Array.isArray(value) ? value : [value]
    }

    to.matched.forEach(record => {
      const meta = record.meta || {}
      normalizePermissions(meta.permissions).forEach(perm => allPermissions.add(perm))
      normalizePermissions(meta.allPermissions).forEach(perm => allPermissions.add(perm))
      const anyList = [
        ...normalizePermissions(meta.anyPermissions),
        ...normalizePermissions(meta.anyPermission)
      ]
      anyList.forEach(perm => anyPermissions.add(perm))
    })

    if (allPermissions.size && !authStore.hasAllPermissions([...allPermissions])) {
      next('/error/forbidden')
      return
    }

    if (anyPermissions.size && !authStore.hasAnyPermission([...anyPermissions])) {
      next('/error/forbidden')
      return
    }
  }

  next()
})

export default router
