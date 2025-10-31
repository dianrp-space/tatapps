<template>
  <div class="max-w-4xl">
    <div class="mb-6 flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold">Profile Settings</h1>
        <p class="text-sm text-gray-500">Kelola profil pengguna dan preferensi akun.</p>
      </div>
      <nav class="hidden md:block text-sm text-blue-600">
        <RouterLink class="hover:underline" to="/settings/profile">Profil</RouterLink>
        <span class="mx-2 text-gray-400">·</span>
        <RouterLink class="hover:underline" to="/settings/company">Profil Perusahaan</RouterLink>
      </nav>
    </div>

    <!-- Profile Form Card -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-6">
      <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
        <i class="pi pi-user"></i>
        Personal Information
      </h2>

      <form @submit.prevent="saveProfile">
        <!-- Avatar Upload -->
        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">Profile Picture</label>
          <div class="flex items-center gap-4">
            <img
              :src="avatarUrl"
              alt="Avatar"
              class="w-24 h-24 rounded-full object-cover border-2 border-gray-200"
            />
            <div>
              <input
                type="file"
                ref="avatarInputRef"
                @change="handleAvatarChange"
                accept="image/*"
                class="hidden"
              />
              <button
                type="button"
                @click="() => avatarInputRef.click()"
                class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
              >
                <i class="pi pi-upload mr-2"></i>
                Upload Photo
              </button>
              <p class="text-xs text-gray-500 mt-2">JPG, PNG or GIF. Max 2MB</p>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Full Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Full Name</label>
            <input
              v-model="profileForm.full_name"
              type="text"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              placeholder="Enter your full name"
              required
            />
          </div>

          <!-- Phone (WhatsApp) -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Phone Number (WhatsApp)
            </label>
            <div class="relative">
              <span class="absolute inset-y-0 left-0 flex items-center pl-3 text-gray-500">
                <i class="pi pi-whatsapp"></i>
              </span>
              <input
                v-model="profileForm.phone"
                type="tel"
                class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="628xxxxxxxxxx"
                pattern="[0-9]{10,15}"
              />
            </div>
            <p class="text-xs text-gray-500 mt-1">Format: 628xxxxxxxxxx (no + or spaces)</p>
          </div>
        </div>

        <!-- Email (Full Width, Read-only) -->
        <div class="mt-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">Email (Login)</label>
          <input
            :value="authStore.user?.email"
            type="text"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg bg-gray-50 cursor-not-allowed"
            disabled
            readonly
          />
          <p class="text-xs text-gray-500 mt-1">Contact admin to change login email</p>
        </div>

        <!-- Save Profile Button -->
        <div class="flex justify-end mt-6">
          <button
            type="submit"
            class="px-6 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors flex items-center gap-2"
            :disabled="isSaving"
          >
            <i class="pi pi-check" v-if="!isSaving"></i>
            <i class="pi pi-spin pi-spinner" v-else></i>
            {{ isSaving ? 'Saving...' : 'Save Changes' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Change Password Card -->
    <div class="bg-white rounded-lg shadow-md p-6">
      <h2 class="text-lg font-semibold mb-4 flex items-center gap-2">
        <i class="pi pi-lock"></i>
        Change Password
      </h2>

      <form @submit.prevent="changePassword">
        <div class="grid grid-cols-1 gap-4 max-w-md">
          <!-- Current Password -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Current Password</label>
            <div class="relative">
              <input
                v-model="passwordForm.current_password"
                :type="showCurrentPassword ? 'text' : 'password'"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent pr-10"
                placeholder="Enter current password"
                required
              />
              <button
                type="button"
                @click="() => { showCurrentPassword = !showCurrentPassword }"
                class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-500"
              >
                <i :class="showCurrentPassword ? 'pi pi-eye-slash' : 'pi pi-eye'"></i>
              </button>
            </div>
          </div>

          <!-- New Password -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">New Password</label>
            <div class="relative">
              <input
                v-model="passwordForm.new_password"
                :type="showNewPassword ? 'text' : 'password'"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent pr-10"
                placeholder="Enter new password"
                minlength="6"
                required
              />
              <button
                type="button"
                @click="() => { showNewPassword = !showNewPassword }"
                class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-500"
              >
                <i :class="showNewPassword ? 'pi pi-eye-slash' : 'pi pi-eye'"></i>
              </button>
            </div>
          </div>

          <!-- Confirm New Password -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Confirm New Password</label>
            <div class="relative">
              <input
                v-model="passwordForm.confirm_password"
                :type="showConfirmPassword ? 'text' : 'password'"
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent pr-10"
                placeholder="Confirm new password"
                minlength="6"
                required
              />
              <button
                type="button"
                @click="() => { showConfirmPassword = !showConfirmPassword }"
                class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-500"
              >
                <i :class="showConfirmPassword ? 'pi pi-eye-slash' : 'pi pi-eye'"></i>
              </button>
            </div>
          </div>
        </div>

        <!-- Password Requirements -->
        <div class="mt-4 p-3 bg-blue-50 rounded-lg">
          <p class="text-xs text-gray-600 mb-1"><strong>Password Requirements:</strong></p>
          <ul class="text-xs text-gray-600 list-disc list-inside space-y-1">
            <li>Minimum 6 characters</li>
            <li>New password must match confirmation</li>
          </ul>
        </div>

        <!-- Change Password Button -->
        <div class="flex justify-end mt-6">
          <button
            type="submit"
            class="px-6 py-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600 transition-colors flex items-center gap-2"
            :disabled="isChangingPassword"
          >
            <i class="pi pi-key" v-if="!isChangingPassword"></i>
            <i class="pi pi-spin pi-spinner" v-else></i>
            {{ isChangingPassword ? 'Changing...' : 'Change Password' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useToast } from 'primevue/usetoast'
import axios from '@/api/axios'

const authStore = useAuthStore()
const toast = useToast()

// Profile form data
const profileForm = reactive({
  full_name: '',
  phone: ''
})

// Avatar handling
const avatarInputRef = ref(null)
const avatarPreview = ref(null)
const avatarFile = ref(null)

// Computed avatar URL
const avatarUrl = computed(() => {
  if (avatarPreview.value) {
    return avatarPreview.value
  }
  if (authStore.user?.avatar) {
    // If avatar path starts with http, use it directly, otherwise prepend backend URL
    if (authStore.user.avatar.startsWith('http')) {
      return authStore.user.avatar
    }
    return `http://localhost:8080/${authStore.user.avatar}`
  }
  return 'https://via.placeholder.com/120'
})

// Password form data
const passwordForm = reactive({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

// UI states
const isSaving = ref(false)
const isChangingPassword = ref(false)
const showCurrentPassword = ref(false)
const showNewPassword = ref(false)
const showConfirmPassword = ref(false)

// Load user data on mount
onMounted(() => {
  if (authStore.user) {
    profileForm.full_name = authStore.user.full_name || ''
    profileForm.phone = authStore.user.phone || ''
  }
})

// Handle avatar file selection
function handleAvatarChange(event) {
  const file = event.target.files[0]
  if (!file) return

  // Validate file size (max 2MB)
  if (file.size > 2 * 1024 * 1024) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'File size must be less than 2MB',
      life: 3000
    })
    return
  }

  // Validate file type
  if (!file.type.startsWith('image/')) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'File must be an image (JPG, PNG, or GIF)',
      life: 3000
    })
    return
  }

  avatarFile.value = file

  // Create preview
  const reader = new FileReader()
  reader.onload = (e) => {
    avatarPreview.value = e.target.result
  }
  reader.readAsDataURL(file)
}

// Save profile changes
async function saveProfile() {
  isSaving.value = true

  try {
    const formData = new FormData()
    formData.append('full_name', profileForm.full_name)
    formData.append('email', authStore.user?.email || '')
    formData.append('phone', profileForm.phone)

    if (avatarFile.value) {
      formData.append('avatar', avatarFile.value)
    }

    const response = await axios.put('/users/profile', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    // Update store with new user data
    await authStore.getProfile()

    // Update form with latest data
    if (authStore.user) {
      profileForm.full_name = authStore.user.full_name || ''
      profileForm.phone = authStore.user.phone || ''
    }

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Profile updated successfully',
      life: 3000
    })

    // Clear avatar preview and file after successful upload
    avatarFile.value = null
    avatarPreview.value = null
  } catch (error) {
    console.error('Profile update error:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.message || 'Failed to update profile',
      life: 3000
    })
  } finally {
    isSaving.value = false
  }
}

// Change password
async function changePassword() {
  // Validate passwords match
  if (passwordForm.new_password !== passwordForm.confirm_password) {
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'New passwords do not match',
      life: 3000
    })
    return
  }

  isChangingPassword.value = true

  try {
    await axios.put('/users/change-password', {
      current_password: passwordForm.current_password,
      new_password: passwordForm.new_password
    })

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Password changed successfully',
      life: 3000
    })

    // Clear password form
    passwordForm.current_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''
  } catch (error) {
    console.error('Password change error:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: error.response?.data?.message || 'Failed to change password',
      life: 3000
    })
  } finally {
    isChangingPassword.value = false
  }
}
</script>
