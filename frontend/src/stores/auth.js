import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/api/axios'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token') || null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role?.name === 'admin')
  const isManager = computed(() => user.value?.role?.name === 'manager')
  const permissionSet = computed(() => {
    const set = new Set()
    const permissions = user.value?.role?.permissions || []
    permissions.forEach(permission => {
      if (permission?.name) {
        set.add(permission.name)
      }
    })
    return set
  })

  async function login(email, password) {
    try {
      const response = await api.post('/auth/login', { email, password })
      token.value = response.data.token
      user.value = response.data.user
      localStorage.setItem('token', response.data.token)
      return true
    } catch (error) {
      throw error
    }
  }

  async function register(userData) {
    try {
      const response = await api.post('/auth/register', userData)
      return response.data
    } catch (error) {
      throw error
    }
  }

  async function getProfile() {
    try {
      const response = await api.get('/auth/profile')
      user.value = response.data
      return response.data
    } catch (error) {
      throw error
    }
  }

  function logout() {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
  }

  async function checkAuth() {
    if (token.value) {
      try {
        await getProfile()
      } catch (error) {
        logout()
      }
    }
  }

  function hasPermission(permission) {
    if (!permission) return false
    if (isAdmin.value) return true
    return permissionSet.value.has(permission)
  }

  function hasAnyPermission(permissions = []) {
    if (!Array.isArray(permissions) || permissions.length === 0) {
      return false
    }
    if (isAdmin.value) return true
    return permissions.some(perm => permissionSet.value.has(perm))
  }

  function hasAllPermissions(permissions = []) {
    if (!Array.isArray(permissions) || permissions.length === 0) {
      return false
    }
    if (isAdmin.value) return true
    return permissions.every(perm => permissionSet.value.has(perm))
  }

  return {
    user,
    token,
    isAuthenticated,
    isAdmin,
    isManager,
    permissionSet,
    login,
    register,
    getProfile,
    logout,
    checkAuth,
    hasPermission,
    hasAnyPermission,
    hasAllPermissions
  }
})
