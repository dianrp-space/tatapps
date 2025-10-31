import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios'

export const useWarehouseStore = defineStore('warehouse', () => {
  const warehouses = ref([])
  const currentWarehouse = ref(null)
  const loading = ref(false)

  async function fetchWarehouses() {
    loading.value = true
    try {
      const response = await api.get('/warehouses')
      warehouses.value = response.data.data || []
      return warehouses.value
    } catch (error) {
      console.error('Fetch warehouses error:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function fetchWarehouse(id) {
    loading.value = true
    try {
      const response = await api.get(`/warehouses/${id}`)
      currentWarehouse.value = response.data.data || response.data
      return currentWarehouse.value
    } catch (error) {
      console.error('Fetch warehouse error:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  async function createWarehouse(data) {
    try {
      const response = await api.post('/warehouses', data)
      const newWarehouse = response.data.data || response.data
      warehouses.value.push(newWarehouse)
      return newWarehouse
    } catch (error) {
      throw error
    }
  }

  async function updateWarehouse(id, data) {
    try {
      const response = await api.put(`/warehouses/${id}`, data)
      const updatedWarehouse = response.data.data || response.data
      const index = warehouses.value.findIndex(w => w.id === id)
      if (index !== -1) {
        warehouses.value[index] = updatedWarehouse
      }
      return updatedWarehouse
    } catch (error) {
      throw error
    }
  }

  async function deleteWarehouse(id) {
    try {
      await api.delete(`/warehouses/${id}`)
      warehouses.value = warehouses.value.filter(w => w.id !== id)
    } catch (error) {
      throw error
    }
  }

  return {
    warehouses,
    currentWarehouse,
    loading,
    fetchWarehouses,
    fetchWarehouse,
    createWarehouse,
    updateWarehouse,
    deleteWarehouse
  }
})
