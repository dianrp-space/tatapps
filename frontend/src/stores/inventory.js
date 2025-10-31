import { defineStore } from 'pinia'
import axios from '@/api/axios'

export const useInventoryStore = defineStore('inventory', {
  state: () => ({
    items: [],
    currentItem: null,
    transactions: [],
    loading: false,
    error: null
  }),

  getters: {
    lowStockItems: (state) => {
      return state.items.filter(item => item.quantity <= item.min_stock)
    },

    itemsByWarehouse: (state) => {
      return (warehouseId) => {
        return state.items.filter(item => item.warehouse_id === warehouseId)
      }
    },

    itemsByCategory: (state) => {
      return (category) => {
        return state.items.filter(item => item.category === category)
      }
    },

    totalInventoryValue: (state) => {
      return state.items.reduce((total, item) => {
        return total + (item.quantity * item.unit_price)
      }, 0)
    }
  },

  actions: {
    async fetchItems() {
      this.loading = true
      this.error = null
      try {
        const response = await axios.get('/inventory')
        this.items = response.data.data || []
        return this.items
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to fetch items'
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchItemById(id) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.get(`/inventory/items/${id}`)
        this.currentItem = response.data.data
        return this.currentItem
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to fetch item'
        throw error
      } finally {
        this.loading = false
      }
    },

    async createItem(itemData) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.post('/inventory/items', itemData)
        this.items.push(response.data.data)
        return response.data.data
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to create item'
        throw error
      } finally {
        this.loading = false
      }
    },

    async updateItem(id, itemData) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.put(`/inventory/items/${id}`, itemData)
        const index = this.items.findIndex(item => item.id === id)
        if (index !== -1) {
          this.items[index] = response.data.data
        }
        if (this.currentItem?.id === id) {
          this.currentItem = response.data.data
        }
        return response.data.data
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to update item'
        throw error
      } finally {
        this.loading = false
      }
    },

    async deleteItem(id) {
      this.loading = true
      this.error = null
      try {
        await axios.delete(`/inventory/items/${id}`)
        this.items = this.items.filter(item => item.id !== id)
        if (this.currentItem?.id === id) {
          this.currentItem = null
        }
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to delete item'
        throw error
      } finally {
        this.loading = false
      }
    },

    async fetchTransactions(itemId) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.get(`/inventory/items/${itemId}/transactions`)
        this.transactions = response.data.data || []
        return this.transactions
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to fetch transactions'
        throw error
      } finally {
        this.loading = false
      }
    },

    async recordTransaction(transactionData) {
      this.loading = true
      this.error = null
      try {
        const response = await axios.post(`/inventory/items/${transactionData.item_id}/transactions`, transactionData)
        this.transactions.unshift(response.data.data)
        
        // Update item quantity in local state
        if (this.currentItem?.id === transactionData.item_id) {
          await this.fetchItemById(transactionData.item_id)
        }
        
        // Refresh items list to update quantities
        const itemIndex = this.items.findIndex(item => item.id === transactionData.item_id)
        if (itemIndex !== -1) {
          await this.fetchItems()
        }
        
        return response.data.data
      } catch (error) {
        this.error = error.response?.data?.message || 'Failed to record transaction'
        throw error
      } finally {
        this.loading = false
      }
    },

    async getLowStockAlerts() {
      await this.fetchItems()
      return this.lowStockItems
    },

    clearCurrentItem() {
      this.currentItem = null
    },

    clearError() {
      this.error = null
    }
  }
})
