<template>
  <div class="max-w-3xl">
    <h1 class="text-2xl font-bold mb-6">Site Settings</h1>

    <div class="bg-white rounded-lg shadow-md p-6">
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label for="app-name" class="block text-sm font-medium text-gray-700 mb-2">
            Application Name
          </label>
          <input
            id="app-name"
            v-model="appName"
            type="text"
            required
            maxlength="120"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Enter application name"
          />
          <p class="text-xs text-gray-500 mt-1">
            This name appears in the sidebar, document title, and other branding areas.
          </p>
        </div>

        <div class="grid gap-6 md:grid-cols-2">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Logo
            </label>
            <div class="border border-dashed border-gray-300 rounded-lg p-4 text-center">
              <div class="mb-3 flex justify-center">
                <img
                  v-if="logoPreview"
                  :src="logoPreview"
                  alt="Logo preview"
                  class="max-h-24 object-contain"
                />
                <div
                  v-else
                  class="h-24 w-24 flex items-center justify-center bg-gray-100 rounded"
                >
                  <i class="pi pi-image text-2xl text-gray-400"></i>
                </div>
              </div>
              <input
                ref="logoInput"
                type="file"
                accept="image/*"
                class="hidden"
                @change="onLogoSelected"
              />
              <button
                type="button"
                class="px-4 py-2 text-sm bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition-colors"
                @click="() => logoInput?.click()"
              >
                Choose Logo
              </button>
              <p class="text-xs text-gray-500 mt-2">
                Recommended: transparent PNG, at least 200px width.
              </p>
              <button
                v-if="logoFile"
                type="button"
                class="block mx-auto mt-3 text-xs text-red-500 hover:underline"
                @click="clearLogo"
              >
                Remove logo
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Favicon
            </label>
            <div class="border border-dashed border-gray-300 rounded-lg p-4 text-center">
              <div class="mb-3 flex justify-center">
                <img
                  v-if="faviconPreview"
                  :src="faviconPreview"
                  alt="Favicon preview"
                  class="h-16 w-16 object-contain rounded"
                />
                <div
                  v-else
                  class="h-16 w-16 flex items-center justify-center bg-gray-100 rounded"
                >
                  <i class="pi pi-star text-xl text-gray-400"></i>
                </div>
              </div>
              <input
                ref="faviconInput"
                type="file"
                accept="image/*"
                class="hidden"
                @change="onFaviconSelected"
              />
              <button
                type="button"
                class="px-4 py-2 text-sm bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition-colors"
                @click="() => faviconInput?.click()"
              >
                Choose Favicon
              </button>
              <p class="text-xs text-gray-500 mt-2">
                Recommended: square image (PNG/ICO), at least 64x64 pixels.
              </p>
              <button
                v-if="faviconFile"
                type="button"
                class="block mx-auto mt-3 text-xs text-red-500 hover:underline"
                @click="clearFavicon"
              >
                Remove favicon
              </button>
            </div>
          </div>
        </div>

        <div class="border-t border-gray-200 pt-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-4 flex items-center gap-2">
            <i class="pi pi-whatsapp text-green-500"></i>
            WhatsApp API
          </h2>
          <div class="grid gap-4 md:grid-cols-2">
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-1">API URL</label>
              <input
                v-model="waApiUrl"
                type="text"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="https://..."
              />
            </div>
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-1">API Key</label>
              <input
                v-model="waApiKey"
                type="password"
                autocomplete="new-password"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Enter API key"
              />
              <p class="text-xs text-gray-500 mt-1">Leave blank to keep the current key.</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Sender Number</label>
              <input
                v-model="waSender"
                type="text"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="628xxxxxxxxxx"
              />
              <p class="text-xs text-gray-500 mt-1">Use international format without + sign.</p>
            </div>
          </div>
        </div>

        <div class="border-t border-gray-200 pt-6">
          <h2 class="text-lg font-semibold text-gray-800 mb-4 flex items-center gap-2">
            <i class="pi pi-envelope text-blue-500"></i>
            SMTP Email
          </h2>
          <div class="grid gap-4 md:grid-cols-2">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">SMTP Host</label>
              <input
                v-model="smtpHost"
                type="text"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="smtp.example.com"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">SMTP Port</label>
              <input
                v-model="smtpPort"
                type="number"
                min="1"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="465"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Username</label>
              <input
                v-model="smtpUsername"
                type="text"
                autocomplete="username"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="email@example.com"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
              <input
                v-model="smtpPassword"
                type="password"
                autocomplete="new-password"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Enter password"
              />
              <p class="text-xs text-gray-500 mt-1">Leave blank to keep the current password.</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">From Email</label>
              <input
                v-model="smtpFromEmail"
                type="email"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="no-reply@example.com"
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">From Name</label>
              <input
                v-model="smtpFromName"
                type="text"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="TatApps Notification"
              />
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-3">
          <button
            type="button"
            class="px-6 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-colors"
            @click="resetForm"
            :disabled="siteStore.isSaving"
          >
            Reset
          </button>
          <button
            type="submit"
            class="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2 disabled:opacity-60 disabled:cursor-not-allowed"
            :disabled="siteStore.isSaving"
          >
            <i v-if="siteStore.isSaving" class="pi pi-spin pi-spinner"></i>
            <i v-else class="pi pi-save"></i>
            {{ siteStore.isSaving ? 'Saving...' : 'Save Changes' }}
          </button>
        </div>
      </form>
    </div>

    <div class="bg-white rounded-lg shadow-md p-6 mt-6">
      <div class="space-y-4">
        <div>
          <h2 class="text-lg font-semibold text-gray-800 flex items-center gap-2">
            <i class="pi pi-database text-blue-500"></i>
            Database Backup & Restore
          </h2>
          <p class="text-sm text-gray-500">
            Gunakan fitur ini untuk membuat cadangan database PostgreSQL atau memulihkan data dari cadangan saat migrasi server.
          </p>
        </div>

        <div class="grid gap-4 md:grid-cols-2">
          <div class="rounded-lg border border-gray-200 p-4">
            <h3 class="text-sm font-semibold text-gray-700 mb-2">Download Backup</h3>
            <p class="text-xs text-gray-500 mb-4">
              Cadangan akan diunduh sebagai file SQL. Simpan file ini di lokasi aman sebelum memindahkan server.
            </p>
            <button
              type="button"
              class="w-full inline-flex items-center justify-center gap-2 rounded-lg bg-blue-500 px-4 py-2 text-sm font-medium text-white hover:bg-blue-600 transition disabled:opacity-60"
              :disabled="isBackingUp"
              @click="handleDatabaseBackup"
            >
              <i :class="isBackingUp ? 'pi pi-spin pi-spinner' : 'pi pi-download'" class="text-sm"></i>
              {{ isBackingUp ? 'Menyiapkan Backup...' : 'Download Backup' }}
            </button>
          </div>

          <div class="rounded-lg border border-gray-200 p-4">
            <h3 class="text-sm font-semibold text-gray-700 mb-2">Restore Database</h3>
            <p class="text-xs text-gray-500 mb-4">
              Restore akan menimpa data saat ini. Pastikan Anda sudah membuat backup terbaru sebelum melanjutkan proses ini.
            </p>
            <div class="flex flex-col gap-2">
              <label class="inline-flex w-max cursor-pointer items-center gap-2 rounded-lg border border-blue-200 px-3 py-2 text-xs font-medium text-blue-600 hover:bg-blue-50 transition">
                <i class="pi pi-upload text-xs"></i>
                <span>Pilih File Backup (.sql)</span>
                <input ref="restoreInput" type="file" accept=".sql,.dump,.txt" class="hidden" @change="onRestoreFileSelected" />
              </label>
              <p v-if="restoreFileName" class="text-xs text-gray-600">Terpilih: {{ restoreFileName }}</p>
              <div class="flex gap-2">
                <button
                  type="button"
                  class="flex-1 inline-flex items-center justify-center gap-2 rounded-lg bg-emerald-500 px-4 py-2 text-sm font-medium text-white hover:bg-emerald-600 transition disabled:opacity-60"
                  :disabled="!restoreFile || isRestoring"
                  @click="handleRestore"
                >
                  <i :class="isRestoring ? 'pi pi-spin pi-spinner' : 'pi pi-check'" class="text-sm"></i>
                  {{ isRestoring ? 'Memulihkan...' : 'Restore Database' }}
                </button>
                <button
                  type="button"
                  class="px-4 py-2 rounded-lg border border-gray-200 text-sm text-gray-600 hover:bg-gray-50 transition"
                  :disabled="!restoreFile || isRestoring"
                  @click="clearRestoreFile"
                >
                  Bersihkan
                </button>
              </div>
            </div>
          </div>
        </div>

        <div class="rounded-md border border-amber-200 bg-amber-50 px-4 py-3 text-xs text-amber-700">
          <p class="font-semibold mb-1">Catatan Penting:</p>
          <ul class="list-disc pl-4 space-y-1">
            <li>Pastikan tidak ada pengguna lain yang aktif saat proses restore.</li>
            <li>Backup dan restore membutuhkan utilitas <code>pg_dump</code> dan <code>psql</code> tersedia di server.</li>
            <li>Disarankan melakukan pengujian di lingkungan staging sebelum menerapkan restore di produksi.</li>
          </ul>
        </div>
      </div>
    </div>

    <teleport to="body">
      <div
        v-if="showRestoreProgress"
        class="fixed inset-0 z-50 flex items-center justify-center bg-slate-900/60 px-4"
        @click.self="!isRestoreProcessing && closeRestoreProgress()"
      >
        <div class="w-full max-w-sm rounded-lg bg-white p-6 text-center shadow-xl">
          <i :class="restoreProgressIcon"></i>
          <p class="mt-4 text-base font-medium text-gray-800">{{ restoreProgressMessage }}</p>
          <p v-if="isRestoreProcessing" class="mt-2 text-xs text-gray-500">
            Mohon tunggu, proses ini mungkin memerlukan beberapa menit.
          </p>
          <p v-if="isRestoreError" class="mt-2 text-xs text-red-600">
            {{ restoreProgressError }}
          </p>
          <button
            v-if="isRestoreError"
            type="button"
            class="mt-5 inline-flex items-center justify-center rounded-lg bg-red-500 px-4 py-2 text-sm font-medium text-white hover:bg-red-600 transition"
            @click="closeRestoreProgress"
          >
            Tutup
          </button>
        </div>
      </div>
    </teleport>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, onBeforeUnmount, computed } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useSiteStore } from '@/stores/site'

const toast = useToast()
const siteStore = useSiteStore()

const appName = ref('')
const logoFile = ref(null)
const faviconFile = ref(null)
const logoPreview = ref('')
const faviconPreview = ref('')
const logoInput = ref(null)
const faviconInput = ref(null)

const restoreFile = ref(null)
const restoreFileName = ref('')
const restoreInput = ref(null)
const isBackingUp = ref(false)
const isRestoring = ref(false)
const showRestoreProgress = ref(false)
const restoreProgressStage = ref('idle')
const restoreProgressError = ref('')
const restoreProgressMessage = computed(() => {
  switch (restoreProgressStage.value) {
    case 'processing':
      return 'Sedang memulihkan database...'
    case 'success':
      return 'Database berhasil dipulihkan.'
    case 'error':
      return 'Restore gagal.'
    default:
      return ''
  }
})
const restoreProgressIcon = computed(() => {
  switch (restoreProgressStage.value) {
    case 'processing':
      return 'pi pi-spin pi-spinner text-3xl text-blue-500'
    case 'success':
      return 'pi pi-check-circle text-3xl text-emerald-500'
    case 'error':
      return 'pi pi-times-circle text-3xl text-red-500'
    default:
      return 'pi pi-info-circle text-3xl text-blue-500'
  }
})

const isRestoreProcessing = computed(() => restoreProgressStage.value === 'processing')
const isRestoreError = computed(() => restoreProgressStage.value === 'error')

const waApiUrl = ref('')
const waApiKey = ref('')
const waSender = ref('')

const smtpHost = ref('')
const smtpPort = ref('')
const smtpUsername = ref('')
const smtpPassword = ref('')
const smtpFromEmail = ref('')
const smtpFromName = ref('')

let logoObjectUrl = null
let faviconObjectUrl = null
let restoreProgressTimer = null

function revokeObjectUrl(current) {
  if (current) {
    URL.revokeObjectURL(current)
  }
}

function closeRestoreProgress() {
  if (restoreProgressTimer) {
    clearTimeout(restoreProgressTimer)
    restoreProgressTimer = null
  }
  showRestoreProgress.value = false
  restoreProgressStage.value = 'idle'
  restoreProgressError.value = ''
}

function syncFromStore() {
  appName.value = siteStore.settings.app_name || ''
  if (!logoFile.value) {
    logoPreview.value = siteStore.logoUrl()
  }
  if (!faviconFile.value) {
    faviconPreview.value = siteStore.faviconUrl()
  }

  waApiUrl.value = siteStore.settings.whatsapp_api_url || ''
  waApiKey.value = siteStore.settings.whatsapp_api_key || ''
  waSender.value = siteStore.settings.whatsapp_sender || ''

  smtpHost.value = siteStore.settings.smtp_host || ''
  smtpPort.value = siteStore.settings.smtp_port ? String(siteStore.settings.smtp_port) : ''
  smtpUsername.value = siteStore.settings.smtp_username || ''
  smtpPassword.value = ''
  smtpFromEmail.value = siteStore.settings.smtp_from_email || ''
  smtpFromName.value = siteStore.settings.smtp_from_name || ''
}

watch(
  () => siteStore.settings,
  () => {
    syncFromStore()
  },
  { deep: true, immediate: true }
)

onMounted(() => {
  siteStore.fetchSiteSettingsAdmin().catch(() => {
    siteStore.fetchSiteSettings().catch(() => {
      toast.add({
        severity: 'error',
        summary: 'Failed to load',
        detail: 'Unable to retrieve current site settings.',
        life: 4000
      })
    })
  })
})

onBeforeUnmount(() => {
  closeRestoreProgress()
  revokeObjectUrl(logoObjectUrl)
  revokeObjectUrl(faviconObjectUrl)
})

function setLogoFile(file) {
  logoFile.value = file
  revokeObjectUrl(logoObjectUrl)
  if (file) {
    logoObjectUrl = URL.createObjectURL(file)
    logoPreview.value = logoObjectUrl
  } else {
    logoObjectUrl = null
    logoPreview.value = siteStore.logoUrl()
  }
}

function setFaviconFile(file) {
  faviconFile.value = file
  revokeObjectUrl(faviconObjectUrl)
  if (file) {
    faviconObjectUrl = URL.createObjectURL(file)
    faviconPreview.value = faviconObjectUrl
  } else {
    faviconObjectUrl = null
    faviconPreview.value = siteStore.faviconUrl()
  }
}

function onLogoSelected(event) {
  const [file] = event.target.files || []
  setLogoFile(file || null)
}

function onFaviconSelected(event) {
  const [file] = event.target.files || []
  setFaviconFile(file || null)
}

function clearLogo() {
  setLogoFile(null)
  if (logoInput.value) {
    logoInput.value.value = ''
  }
}

function clearFavicon() {
  setFaviconFile(null)
  if (faviconInput.value) {
    faviconInput.value.value = ''
  }
}

function resetForm() {
  clearLogo()
  clearFavicon()
  syncFromStore()
}

function clearRestoreFile() {
  restoreFile.value = null
  restoreFileName.value = ''
  if (restoreInput.value) {
    restoreInput.value.value = ''
  }
}

function onRestoreFileSelected(event) {
  const [file] = event.target.files || []
  if (file) {
    restoreFile.value = file
    restoreFileName.value = file.name
  } else {
    clearRestoreFile()
  }
}

async function handleSubmit() {
  if (!appName.value.trim()) {
    toast.add({
      severity: 'warn',
      summary: 'Validation',
      detail: 'Application name is required.',
      life: 3000
    })
    return
  }

  const formData = new FormData()
  formData.append('app_name', appName.value.trim())

  if (logoFile.value) {
    formData.append('logo', logoFile.value)
  }

  if (faviconFile.value) {
    formData.append('favicon', faviconFile.value)
  }

  if (waApiUrl.value?.trim()) {
    formData.append('whatsapp_api_url', waApiUrl.value.trim())
  } else {
    formData.append('whatsapp_api_url', '')
  }

  if (waApiKey.value?.trim()) {
    formData.append('whatsapp_api_key', waApiKey.value.trim())
  }

  if (waSender.value?.trim()) {
    formData.append('whatsapp_sender', waSender.value.trim())
  } else {
    formData.append('whatsapp_sender', '')
  }

  if (smtpHost.value?.trim()) {
    formData.append('smtp_host', smtpHost.value.trim())
  } else {
    formData.append('smtp_host', '')
  }

  if (smtpPort.value?.toString()) {
    formData.append('smtp_port', smtpPort.value.toString())
  } else {
    formData.append('smtp_port', '')
  }

  if (smtpUsername.value?.trim()) {
    formData.append('smtp_username', smtpUsername.value.trim())
  } else {
    formData.append('smtp_username', '')
  }

  if (smtpPassword.value?.trim()) {
    formData.append('smtp_password', smtpPassword.value.trim())
  }

  if (smtpFromEmail.value?.trim()) {
    formData.append('smtp_from_email', smtpFromEmail.value.trim())
  } else {
    formData.append('smtp_from_email', '')
  }

  if (smtpFromName.value?.trim()) {
    formData.append('smtp_from_name', smtpFromName.value.trim())
  } else {
    formData.append('smtp_from_name', '')
  }

  try {
    await siteStore.updateSiteSettings(formData)
    await siteStore.fetchSiteSettingsAdmin()
    toast.add({
      severity: 'success',
      summary: 'Updated',
      detail: 'Site settings saved successfully.',
      life: 3000
    })
    clearLogo()
    clearFavicon()
    syncFromStore()
  } catch (error) {
    const detail = error.response?.data?.error || 'Failed to update site settings.'
    toast.add({
      severity: 'error',
      summary: 'Save failed',
      detail,
      life: 4000
    })
  }
}

async function handleDatabaseBackup() {
  try {
    isBackingUp.value = true
    const { blob, filename } = await siteStore.downloadDatabaseBackup()
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-').slice(0, 19)
    link.href = url
    link.download = filename || `tatapps_backup_${timestamp}.sql`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    toast.add({ severity: 'success', summary: 'Backup Dibuat', detail: 'File backup berhasil diunduh.', life: 4000 })
  } catch (error) {
    const message = error?.response?.data?.error || 'Gagal membuat backup database.'
    toast.add({ severity: 'error', summary: 'Backup Gagal', detail: message, life: 5000 })
  } finally {
    isBackingUp.value = false
  }
}

async function handleRestore() {
  if (!restoreFile.value) return
  const confirmation = confirm('Restore akan menimpa seluruh data. Pastikan Anda memiliki backup terbaru. Lanjutkan?')
  if (!confirmation) return

  try {
    isRestoring.value = true
    restoreProgressError.value = ''
    restoreProgressStage.value = 'processing'
    showRestoreProgress.value = true
    if (restoreProgressTimer) {
      clearTimeout(restoreProgressTimer)
      restoreProgressTimer = null
    }

    await siteStore.restoreDatabaseBackup(restoreFile.value)

    restoreProgressStage.value = 'success'
    toast.add({ severity: 'success', summary: 'Restore Berhasil', detail: 'Database telah dipulihkan.', life: 4000 })
    clearRestoreFile()

    restoreProgressTimer = setTimeout(() => {
      closeRestoreProgress()
      setTimeout(() => {
        window.location.reload()
      }, 300)
    }, 1200)
  } catch (error) {
    const message = error?.response?.data?.error || 'Restore database gagal.'
    restoreProgressError.value = message
    restoreProgressStage.value = 'error'
    showRestoreProgress.value = true
    if (restoreProgressTimer) {
      clearTimeout(restoreProgressTimer)
      restoreProgressTimer = null
    }
    toast.add({ severity: 'error', summary: 'Restore Gagal', detail: message, life: 5000 })
  } finally {
    isRestoring.value = false
  }
}
</script>
