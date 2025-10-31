<template>
  <div class="max-w-4xl">
    <h1 class="text-2xl font-bold mb-6">Notification Settings</h1>

    <!-- Low Stock Alerts Card -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
        <i class="pi pi-bell"></i>
        Low Stock Alerts
      </h2>

      <form @submit.prevent="saveSettings">
        <!-- Enable Low Stock Notifications -->
        <div class="flex items-center justify-between mb-6 p-4 bg-gray-50 rounded-lg">
          <div>
            <h3 class="font-medium text-gray-900">Enable Low Stock Notifications</h3>
            <p class="text-sm text-gray-500">Receive alerts when inventory items are running low</p>
          </div>
          <label class="relative inline-flex items-center cursor-pointer">
            <input
              type="checkbox"
              v-model="settings.enabled"
              class="sr-only peer"
            />
            <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
          </label>
        </div>

        <!-- Schedule -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Schedule
          </label>
          <div class="flex flex-wrap gap-3 mb-4">
            <button
              v-for="mode in scheduleModes"
              :key="mode.value"
              type="button"
              @click="setScheduleMode(mode.value)"
              :disabled="!settings.enabled"
              class="flex-1 min-w-[200px] px-4 py-3 rounded-lg border transition text-left disabled:opacity-60 disabled:cursor-not-allowed"
              :class="settings.schedule_mode === mode.value
                ? 'border-blue-500 bg-blue-50 shadow-sm'
                : 'border-gray-200 hover:border-blue-300'"
            >
              <div class="font-medium text-gray-900">{{ mode.label }}</div>
              <div class="text-xs text-gray-500 mt-1">{{ mode.description }}</div>
            </button>
          </div>

          <div v-if="!isCronMode" class="space-y-2">
            <select
              v-model="settings.check_frequency"
              class="w-full md:w-64 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :disabled="!settings.enabled"
            >
              <option value="hourly">Every Hour</option>
              <option value="daily">Daily (9:00 AM)</option>
              <option value="weekly">Weekly (Monday 9:00 AM)</option>
            </select>
            <p class="text-xs text-gray-500">Notifications will run according to the selected preset.</p>
          </div>

          <div v-else class="space-y-2">
            <input
              v-model="settings.cron_expression"
              type="text"
              placeholder="0 9 * * *"
              class="w-full md:w-80 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :disabled="!settings.enabled"
            />
            <p class="text-xs text-gray-500">
              Use standard cron format: minute hour day-of-month month day-of-week [year].
            </p>
          </div>

          <div class="mt-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Time Zone
            </label>
            <input
              v-model="settings.timezone"
              type="text"
              list="timezone-options"
              placeholder="e.g., Asia/Jakarta"
              class="w-full md:w-72 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              :disabled="!settings.enabled"
            />
            <datalist id="timezone-options">
              <option
                v-for="tz in timezoneOptions"
                :key="tz"
                :value="tz"
              />
            </datalist>
            <p class="text-xs text-gray-500 mt-1">
              Notifications follow the selected timezone.
            </p>
          </div>
        </div>

        <!-- Notification Channels -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-3">
            Notification Channels
          </label>
          
          <!-- WhatsApp Notification -->
          <div class="flex items-start gap-3 mb-4 p-4 border border-gray-200 rounded-lg">
            <input
              type="checkbox"
              v-model="settings.whatsapp_enabled"
              id="whatsapp"
              class="mt-1 w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              :disabled="!settings.enabled"
            />
            <div class="flex-1">
              <label for="whatsapp" class="font-medium text-gray-900 flex items-center gap-2 cursor-pointer">
                <i class="pi pi-whatsapp text-green-600"></i>
                WhatsApp Notifications
              </label>
              <p class="text-sm text-gray-500 mb-2">Receive alerts via WhatsApp</p>
              <div
                v-if="whatsappNumbers.length"
                class="flex flex-wrap gap-2 mb-2"
              >
                <span
                  v-for="number in whatsappNumbers"
                  :key="number"
                  class="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-green-50 text-green-700 text-xs border border-green-200"
                >
                  {{ number }}
                  <button
                    type="button"
                    class="text-green-600 hover:text-green-800"
                    :disabled="whatsappInputsDisabled"
                    @click="removeWhatsappNumber(number)"
                  >
                    <i class="pi pi-times text-xs"></i>
                  </button>
                </span>
              </div>
              <div class="flex flex-col md:flex-row md:items-center gap-2">
                <input
                  v-model="newWhatsappNumber"
                  type="text"
                  placeholder="Add WhatsApp number"
                  class="w-full md:w-64 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                  :disabled="whatsappInputsDisabled"
                  @keydown.enter.prevent="addWhatsappNumber"
                />
                <button
                  type="button"
                  class="px-3 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:opacity-60 disabled:cursor-not-allowed text-sm"
                  :disabled="whatsappInputsDisabled"
                  @click="addWhatsappNumber"
                >
                  Add Number
                </button>
              </div>
              <p class="text-xs text-gray-500 mt-1">
                Format: 628xxxxxxxxxx (no + or spaces). Multiple numbers are sent as {{ whatsappNumbers.length ? joinedWhatsappNumbers : '628111223|628569...' }}.
              </p>
            </div>
          </div>

          <!-- Email Notification -->
          <div class="flex items-start gap-3 p-4 border border-gray-200 rounded-lg">
            <input
              type="checkbox"
              v-model="settings.email_enabled"
              id="email"
              class="mt-1 w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              :disabled="!settings.enabled"
            />
            <div class="flex-1">
              <label for="email" class="font-medium text-gray-900 flex items-center gap-2 cursor-pointer">
                <i class="pi pi-envelope text-blue-600"></i>
                Email Notifications
              </label>
              <p class="text-sm text-gray-500 mb-2">Receive alerts via email</p>
              <input
                v-model="settings.email_address"
                type="email"
                placeholder="your.email@example.com"
                class="w-full md:w-80 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm"
                :disabled="!settings.enabled || !settings.email_enabled"
              />
            </div>
          </div>
        </div>

        <!-- Save Button -->
        <div class="flex justify-end gap-3">
          <button
            type="button"
            @click="testNotification"
            class="px-6 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors flex items-center gap-2"
            :disabled="!settings.enabled || isTesting"
          >
            <i class="pi pi-send" v-if="!isTesting"></i>
            <i class="pi pi-spin pi-spinner" v-else></i>
            {{ isTesting ? 'Sending...' : 'Send Test' }}
          </button>
          <button
            type="submit"
            class="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2"
            :disabled="isSaving"
          >
            <i class="pi pi-check" v-if="!isSaving"></i>
            <i class="pi pi-spin pi-spinner" v-else></i>
            {{ isSaving ? 'Saving...' : 'Save Settings' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Recent Notifications -->
    <div class="bg-white rounded-lg shadow-md p-6">
      <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
        <i class="pi pi-history"></i>
        Recent Notifications
      </h2>

      <!-- Loading State -->
      <div v-if="loadingHistory" class="flex justify-center py-8">
        <i class="pi pi-spin pi-spinner text-3xl text-blue-500"></i>
      </div>

      <!-- Notifications List -->
      <div v-else-if="notificationHistory.length > 0" class="space-y-3">
        <div
          v-for="notification in notificationHistory"
          :key="notification.id"
          class="flex items-start gap-3 p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
        >
          <div
            class="p-2 rounded-full"
            :class="{
              'bg-yellow-100 text-yellow-600': notification.type === 'low_stock',
              'bg-blue-100 text-blue-600': notification.type === 'test',
              'bg-red-100 text-red-600': notification.type === 'critical'
            }"
          >
            <i
              :class="{
                'pi pi-exclamation-triangle': notification.type === 'low_stock',
                'pi pi-info-circle': notification.type === 'test',
                'pi pi-times-circle': notification.type === 'critical'
              }"
            ></i>
          </div>
          <div class="flex-1">
            <div class="flex items-center justify-between">
              <h3 class="font-medium text-gray-900">{{ notification.title }}</h3>
              <span class="text-xs text-gray-500">{{ formatDate(notification.created_at) }}</span>
            </div>
            <p class="text-sm text-gray-600 mt-1">{{ notification.message }}</p>
            <div class="flex items-center gap-2 mt-2">
              <span
                v-if="notification.whatsapp_sent"
                class="text-xs px-2 py-1 bg-green-100 text-green-700 rounded"
              >
                <i class="pi pi-whatsapp mr-1"></i>WhatsApp
              </span>
              <span
                v-if="notification.email_sent"
                class="text-xs px-2 py-1 bg-blue-100 text-blue-700 rounded"
              >
                <i class="pi pi-envelope mr-1"></i>Email
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-8 text-gray-500">
        <i class="pi pi-inbox text-4xl mb-2"></i>
        <p>No notifications sent yet</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useToast } from 'primevue/usetoast'
import axios from '@/api/axios'

const authStore = useAuthStore()
const toast = useToast()

const fallbackTimezones = ['Asia/Jakarta', 'Asia/Makassar', 'Asia/Jayapura', 'UTC']

const resolveDefaultTimezone = () => {
  if (typeof Intl !== 'undefined' && Intl.DateTimeFormat) {
    try {
      return Intl.DateTimeFormat().resolvedOptions().timeZone || 'Asia/Jakarta'
    } catch {
      return 'Asia/Jakarta'
    }
  }
  return 'Asia/Jakarta'
}

const resolveTimezoneOptions = (defaultTimezone) => {
  let options = []
  if (typeof Intl !== 'undefined' && typeof Intl.supportedValuesOf === 'function') {
    try {
      options = Intl.supportedValuesOf('timeZone') || []
    } catch {
      options = []
    }
  }
  if (!options.length) {
    options = [...fallbackTimezones]
  }
  const existing = new Set(options)
  if (defaultTimezone && !existing.has(defaultTimezone)) {
    options = [defaultTimezone, ...options]
  }
  return options
}

const defaultTimezone = resolveDefaultTimezone()
const timezoneOptions = ref(resolveTimezoneOptions(defaultTimezone))

const scheduleModes = [
  {
    label: 'Preset Options',
    value: 'preset',
    description: 'Use hourly, daily, or weekly intervals'
  },
  {
    label: 'Custom Cron',
    value: 'cron',
    description: 'Define your own cron expression'
  }
]

const whatsappNumbers = ref([])
const newWhatsappNumber = ref('')

// Settings data
const settings = reactive({
  enabled: false,
  threshold: 10,
  check_frequency: 'daily',
  schedule_mode: 'preset',
  cron_expression: '0 9 * * *',
  timezone: defaultTimezone,
  whatsapp_enabled: false,
  whatsapp_number: '',
  email_enabled: false,
  email_address: ''
})

// UI states
const isSaving = ref(false)
const isTesting = ref(false)
const loadingHistory = ref(false)
const notificationHistory = ref([])

const isCronMode = computed(() => settings.schedule_mode === 'cron')
const whatsappInputsDisabled = computed(() => !settings.enabled || !settings.whatsapp_enabled)
const joinedWhatsappNumbers = computed(() => whatsappNumbers.value.join('|'))

const ensureTimezoneListed = (value) => {
  if (!value) return
  if (!timezoneOptions.value.includes(value)) {
    timezoneOptions.value = [value, ...timezoneOptions.value]
  }
}

const parseWhatsappNumbers = (value) => {
  if (!value) return []
  return value
    .split(/[\n,|;]+/)
    .map(entry => entry.replace(/\s+/g, '').replace(/[^0-9]/g, ''))
    .filter(Boolean)
    .filter((item, index, arr) => arr.indexOf(item) === index)
}

const addWhatsappNumber = () => {
  if (whatsappInputsDisabled.value) {
    return
  }
  const candidates = parseWhatsappNumbers(newWhatsappNumber.value)
  if (!candidates.length) {
    toast.add({
      severity: 'warn',
      summary: 'Invalid Number',
      detail: 'Please enter at least one valid WhatsApp number',
      life: 3000
    })
    return
  }
  let added = false
  candidates.forEach(number => {
    if (!whatsappNumbers.value.includes(number)) {
      whatsappNumbers.value.push(number)
      added = true
    }
  })
  if (added) {
    newWhatsappNumber.value = ''
  } else {
    toast.add({
      severity: 'info',
      summary: 'Duplicate Number',
      detail: 'All numbers you entered are already in the list',
      life: 2500
    })
  }
}

const removeWhatsappNumber = (number) => {
  whatsappNumbers.value = whatsappNumbers.value.filter(item => item !== number)
}

const setScheduleMode = (mode) => {
  if (mode === 'cron') {
    settings.schedule_mode = 'cron'
    if (!settings.cron_expression) {
      settings.cron_expression = '0 9 * * *'
    }
  } else {
    settings.schedule_mode = 'preset'
    if (!settings.check_frequency) {
      settings.check_frequency = 'daily'
    }
  }
}

const isCronExpressionValid = (expression) => {
  if (!expression) return false
  const segments = expression.trim().split(/\s+/)
  return segments.length === 5 || segments.length === 6
}

watch(joinedWhatsappNumbers, (value) => {
  settings.whatsapp_number = value
})

watch(
  () => settings.timezone,
  (value) => {
    if (!value) return
    ensureTimezoneListed(value.trim())
  }
)

// Load settings on mount
onMounted(async () => {
  await loadSettings()
  await loadNotificationHistory()
})

// Load current settings
async function loadSettings() {
  try {
    const response = await axios.get('/settings/notifications')
    if (response.data.data) {
      Object.assign(settings, {
        ...settings,
        ...response.data.data
      })
    }

    setScheduleMode(settings.schedule_mode)

    if (!settings.timezone) {
      settings.timezone = defaultTimezone
    }
    ensureTimezoneListed(settings.timezone)

    if (!settings.cron_expression) {
      settings.cron_expression = '0 9 * * *'
    }

    if (!settings.check_frequency) {
      settings.check_frequency = 'daily'
    }

    const parsedNumbers = parseWhatsappNumbers(settings.whatsapp_number)
    whatsappNumbers.value = parsedNumbers

    // Pre-fill with user data if empty
    if (!settings.whatsapp_number && authStore.user?.phone) {
      const initialNumbers = parseWhatsappNumbers(authStore.user.phone)
      whatsappNumbers.value = initialNumbers
    }
    if (!settings.email_address && authStore.user?.email) {
      settings.email_address = authStore.user.email
    }

    settings.whatsapp_number = joinedWhatsappNumbers.value
  } catch (error) {
    console.error('Failed to load settings:', error)
    // Pre-fill with user data on error
    if (authStore.user) {
      const fallback = parseWhatsappNumbers(authStore.user.phone || '')
      whatsappNumbers.value = fallback
      settings.whatsapp_number = joinedWhatsappNumbers.value
      settings.email_address = authStore.user.email || ''
    }
  }
}

// Save settings
async function saveSettings() {
  // Validation
  if (settings.enabled) {
    if (settings.whatsapp_enabled && !whatsappNumbers.value.length) {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Please add at least one WhatsApp number',
        life: 3000
      })
      return
    }
    if (settings.email_enabled && !settings.email_address) {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Please enter email address',
        life: 3000
      })
      return
    }
    if (!settings.whatsapp_enabled && !settings.email_enabled) {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Please enable at least one notification channel',
        life: 3000
      })
      return
    }
    if (isCronMode.value && !isCronExpressionValid(settings.cron_expression)) {
      toast.add({
        severity: 'error',
        summary: 'Error',
        detail: 'Please provide a valid cron expression (5 or 6 segments)',
        life: 3000
      })
      return
    }
  }

  settings.timezone = (settings.timezone || '').trim() || defaultTimezone
  ensureTimezoneListed(settings.timezone)

  isSaving.value = true

  try {
    const payload = {
      enabled: settings.enabled,
      threshold: settings.threshold,
      check_frequency: settings.check_frequency,
      schedule_mode: settings.schedule_mode,
      cron_expression: settings.cron_expression,
      timezone: settings.timezone,
      whatsapp_enabled: settings.whatsapp_enabled,
      whatsapp_number: joinedWhatsappNumbers.value,
      email_enabled: settings.email_enabled,
      email_address: settings.email_address
    }

    const response = await axios.put('/settings/notifications', payload)

    if (response.data?.data) {
      Object.assign(settings, {
        ...settings,
        ...response.data.data
      })
      const syncedNumbers = parseWhatsappNumbers(response.data.data.whatsapp_number)
      whatsappNumbers.value = syncedNumbers
      setScheduleMode(settings.schedule_mode)
      if (!settings.cron_expression) {
        settings.cron_expression = '0 9 * * *'
      }
      if (!settings.check_frequency) {
        settings.check_frequency = 'daily'
      }
      settings.timezone = (settings.timezone || '').trim() || defaultTimezone
      ensureTimezoneListed(settings.timezone)
    }

    settings.whatsapp_number = joinedWhatsappNumbers.value

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Notification settings saved successfully',
      life: 3000
    })
  } catch (error) {
    console.error('Failed to save settings:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.message || 'Failed to save settings',
      life: 3000
    })
  } finally {
    isSaving.value = false
  }
}

// Send test notification
async function testNotification() {
  if (!settings.whatsapp_enabled && !settings.email_enabled) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Please enable at least one notification channel',
      life: 3000
    })
    return
  }

  if (settings.whatsapp_enabled && !whatsappNumbers.value.length) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Please add at least one WhatsApp number before sending a test',
      life: 3000
    })
    return
  }

  isTesting.value = true

  try {
    await axios.post('/notifications/test', {
      whatsapp_enabled: settings.whatsapp_enabled,
      whatsapp_number: joinedWhatsappNumbers.value,
      email_enabled: settings.email_enabled,
      email_address: settings.email_address
    })

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Test notification sent successfully',
      life: 3000
    })

    // Reload history
    await loadNotificationHistory()
  } catch (error) {
    console.error('Failed to send test notification:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.message || 'Failed to send test notification',
      life: 3000
    })
  } finally {
    isTesting.value = false
  }
}

// Load notification history
async function loadNotificationHistory() {
  loadingHistory.value = true

  try {
    const response = await axios.get('/notifications/history?limit=10')
    notificationHistory.value = response.data.data || []
  } catch (error) {
    console.error('Failed to load notification history:', error)
    notificationHistory.value = []
  } finally {
    loadingHistory.value = false
  }
}

// Format date
function formatDate(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  const now = new Date()
  const diff = now - date

  // Less than 1 minute
  if (diff < 60000) {
    return 'Just now'
  }

  // Less than 1 hour
  if (diff < 3600000) {
    const minutes = Math.floor(diff / 60000)
    return `${minutes} minute${minutes > 1 ? 's' : ''} ago`
  }

  // Less than 1 day
  if (diff < 86400000) {
    const hours = Math.floor(diff / 3600000)
    return `${hours} hour${hours > 1 ? 's' : ''} ago`
  }

  // More than 1 day
  const days = Math.floor(diff / 86400000)
  if (days < 7) {
    return `${days} day${days > 1 ? 's' : ''} ago`
  }

  // Format as date
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: date.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
  })
}
</script>
