import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api/axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1'
const FILE_BASE_URL = API_BASE_URL.replace(/\/api\/v1$/, '')
const DEFAULT_APP_NAME = 'TatApps'
const DEFAULT_FAVICON = '/vite.svg'

function buildAssetUrl(path) {
  if (!path) return ''
  if (path.startsWith('http')) {
    return path
  }
  return `${FILE_BASE_URL}/${path}`
}

function applyBranding(settings) {
  const title = settings?.app_name?.trim() || DEFAULT_APP_NAME
  document.title = title

  let link = document.querySelector("link[rel='icon']")
  if (!link) {
    link = document.createElement('link')
    link.rel = 'icon'
    document.head.appendChild(link)
  }
  link.href = settings?.favicon ? buildAssetUrl(settings.favicon) : DEFAULT_FAVICON
}

export const useSiteStore = defineStore('site', () => {
  const settings = ref({
    app_name: DEFAULT_APP_NAME,
    logo: '',
    favicon: '',
    whatsapp_api_url: '',
    whatsapp_api_key: '',
    whatsapp_sender: '',
    smtp_host: '',
    smtp_port: '',
    smtp_username: '',
    smtp_password: '',
    smtp_from_email: '',
    smtp_from_name: ''
  })
  const isLoading = ref(false)
  const isSaving = ref(false)

  async function fetchSiteSettings() {
    try {
      isLoading.value = true
      const response = await api.get('/settings/site')
      if (response.data?.data) {
        settings.value = { ...settings.value, ...response.data.data }
        applyBranding(settings.value)
      }
      return settings.value
    } finally {
      isLoading.value = false
    }
  }

  async function fetchSiteSettingsAdmin() {
    try {
      isLoading.value = true
      const response = await api.get('/settings/site/admin')
      if (response.data?.data) {
        settings.value = { ...settings.value, ...response.data.data }
        applyBranding(settings.value)
      }
      return settings.value
    } finally {
      isLoading.value = false
    }
  }

  async function updateSiteSettings(formData) {
    try {
      isSaving.value = true
      const response = await api.put('/settings/site', formData, {
        headers: {
          'Content-Type': 'multipart/form-data'
        }
      })
      if (response.data?.data) {
        settings.value = { ...settings.value, ...response.data.data }
        applyBranding(settings.value)
      }
      return response.data
    } finally {
      isSaving.value = false
    }
  }

  function logoUrl() {
    return buildAssetUrl(settings.value.logo)
  }

  function faviconUrl() {
    return buildAssetUrl(settings.value.favicon)
  }

  async function downloadDatabaseBackup() {
    const response = await api.get('/settings/database/backup', { responseType: 'blob' })
    let filename = 'database_backup.sql'
    const disposition = response.headers['content-disposition']
    if (disposition) {
      const match = /filename="?([^";]+)"?/i.exec(disposition)
      if (match && match[1]) {
        filename = match[1]
      }
    }
    return { blob: response.data, filename }
  }

  async function restoreDatabaseBackup(file) {
    const formData = new FormData()
    formData.append('backup', file)
    const response = await api.post('/settings/database/restore', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
    return response.data
  }

  return {
    settings,
    isLoading,
    isSaving,
    fetchSiteSettings,
    fetchSiteSettingsAdmin,
    updateSiteSettings,
    logoUrl,
    faviconUrl,
    downloadDatabaseBackup,
    restoreDatabaseBackup
  }
})
