<template>
  <div class="min-h-screen bg-gradient-to-b from-gray-50 via-white to-gray-100 text-gray-800">
    <div class="mx-auto flex min-h-screen max-w-5xl flex-col px-6 py-10">
      <header class="flex items-center justify-between">
        <RouterLink
          to="/"
          class="group flex items-center gap-3 rounded-lg px-3 py-2 transition-colors hover:bg-blue-50"
        >
          <img
            v-if="siteLogo"
            :src="siteLogo"
            :alt="`${appName} logo`"
            class="h-10 w-auto object-contain transition-transform group-hover:scale-105"
          />
          <span class="text-lg font-semibold text-gray-700 group-hover:text-blue-600">{{ appName }}</span>
        </RouterLink>
        <slot name="header-cta"></slot>
      </header>

      <main class="flex flex-1 flex-col items-center justify-center text-center">
        <p
          v-if="badgeText"
          class="text-xs font-semibold uppercase tracking-[0.3em] text-blue-600 md:text-sm"
        >
          {{ badgeText }}
        </p>

        <div class="mt-8 flex flex-col items-center gap-4 md:flex-row md:gap-6">
          <span class="text-6xl font-black leading-none text-gray-900 sm:text-7xl md:text-8xl">
            {{ code }}
          </span>
          <div class="text-left">
            <p class="text-sm font-medium uppercase tracking-[0.4em] text-gray-400">Error</p>
            <h1 class="mt-2 text-3xl font-bold text-gray-900 sm:text-4xl md:text-5xl">{{ title }}</h1>
          </div>
        </div>

        <p class="mt-6 max-w-2xl text-base text-gray-600 sm:text-lg">
          {{ description }}
        </p>

        <div class="mt-10 flex flex-wrap items-center justify-center gap-3">
          <template v-if="primaryAction">
            <RouterLink
              v-if="primaryAction.type === 'router'"
              :to="primaryAction.to"
              class="inline-flex items-center gap-2 rounded-lg border border-transparent bg-blue-600 px-5 py-2.5 text-sm font-semibold text-white shadow-sm transition-all hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            >
              <i v-if="primaryAction.icon" :class="primaryAction.icon"></i>
              <span>{{ primaryAction.label }}</span>
            </RouterLink>
            <a
              v-else-if="primaryAction.type === 'external'"
              :href="primaryAction.href"
              :target="primaryAction.external ? '_blank' : '_self'"
              :rel="primaryAction.external ? 'noopener noreferrer' : null"
              class="inline-flex items-center gap-2 rounded-lg border border-transparent bg-blue-600 px-5 py-2.5 text-sm font-semibold text-white shadow-sm transition-all hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            >
              <i v-if="primaryAction.icon" :class="primaryAction.icon"></i>
              <span>{{ primaryAction.label }}</span>
            </a>
            <button
              v-else
              type="button"
              class="inline-flex items-center gap-2 rounded-lg border border-transparent bg-blue-600 px-5 py-2.5 text-sm font-semibold text-white shadow-sm transition-all hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
              @click="handleAction(primaryAction)"
            >
              <i v-if="primaryAction.icon" :class="primaryAction.icon"></i>
              <span>{{ primaryAction.label }}</span>
            </button>
          </template>

          <template v-if="secondaryAction">
            <RouterLink
              v-if="secondaryAction.type === 'router'"
              :to="secondaryAction.to"
              class="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-semibold text-gray-700 shadow-sm transition-all hover:border-gray-300 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            >
              <i v-if="secondaryAction.icon" :class="secondaryAction.icon"></i>
              <span>{{ secondaryAction.label }}</span>
            </RouterLink>
            <a
              v-else-if="secondaryAction.type === 'external'"
              :href="secondaryAction.href"
              :target="secondaryAction.external ? '_blank' : '_self'"
              :rel="secondaryAction.external ? 'noopener noreferrer' : null"
              class="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-semibold text-gray-700 shadow-sm transition-all hover:border-gray-300 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            >
              <i v-if="secondaryAction.icon" :class="secondaryAction.icon"></i>
              <span>{{ secondaryAction.label }}</span>
            </a>
            <button
              v-else
              type="button"
              class="inline-flex items-center gap-2 rounded-lg border border-gray-200 bg-white px-5 py-2.5 text-sm font-semibold text-gray-700 shadow-sm transition-all hover:border-gray-300 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
              @click="handleAction(secondaryAction)"
            >
              <i v-if="secondaryAction.icon" :class="secondaryAction.icon"></i>
              <span>{{ secondaryAction.label }}</span>
            </button>
          </template>

          <button
            v-if="showBackButton"
            type="button"
            class="inline-flex items-center gap-2 rounded-lg border border-transparent px-5 py-2.5 text-sm font-semibold text-gray-500 transition-colors hover:text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
            @click="navigateBack"
          >
            <i class="pi pi-arrow-left text-sm"></i>
            <span>Go Back</span>
          </button>
        </div>
      </main>

      <footer class="mt-16 border-t border-gray-200 pt-6 text-center text-xs text-gray-400 sm:text-sm">
        &copy; {{ currentYear }} {{ appName }}. All rights reserved.
      </footer>
    </div>
  </div>
</template>

<script setup>
import { computed, watchEffect } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useSiteStore } from '@/stores/site'

const props = defineProps({
  code: {
    type: [String, Number],
    default: '404'
  },
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    default: ''
  },
  badge: {
    type: String,
    default: ''
  },
  primaryAction: {
    type: Object,
    default: null
  },
  secondaryAction: {
    type: Object,
    default: null
  },
  showBackButton: {
    type: Boolean,
    default: false
  }
})

const router = useRouter()
const siteStore = useSiteStore()

const appName = computed(() => siteStore.settings.app_name || 'TatApps')
const siteLogo = computed(() => siteStore.logoUrl())
const badgeText = computed(() => props.badge.trim())
const currentYear = new Date().getFullYear()

const primaryAction = computed(() => normalizeAction(props.primaryAction, 'primary'))
const secondaryAction = computed(() => normalizeAction(props.secondaryAction, 'secondary'))

watchEffect(() => {
  if (typeof document === 'undefined') return
  const title = props.title
    ? `${props.code} ${props.title} | ${appName.value}`
    : appName.value
  document.title = title
})

function normalizeAction(action, fallbackVariant) {
  if (!action || !action.label) return null
  if (action.to) {
    return {
      label: action.label,
      type: 'router',
      to: action.to,
      icon: action.icon || null,
      variant: action.variant || fallbackVariant
    }
  }
  if (action.href) {
    return {
      label: action.label,
      type: 'external',
      href: action.href,
      external: Boolean(action.external),
      icon: action.icon || null,
      variant: action.variant || fallbackVariant
    }
  }
  return {
    label: action.label,
    type: 'event',
    event: action.event || null,
    handler: typeof action.onClick === 'function' ? action.onClick : null,
    icon: action.icon || null,
    variant: action.variant || fallbackVariant
  }
}

function navigateBack() {
  if (typeof window !== 'undefined' && window.history.length > 1) {
    router.back()
  } else {
    router.push('/')
  }
}

function handleAction(action) {
  if (!action) return

  if (typeof action.handler === 'function') {
    action.handler()
    return
  }

  switch (action.event) {
    case 'back':
      navigateBack()
      break
    case 'reload':
      if (typeof window !== 'undefined') {
        window.location.reload()
      }
      break
    default:
      navigateBack()
  }
}
</script>
