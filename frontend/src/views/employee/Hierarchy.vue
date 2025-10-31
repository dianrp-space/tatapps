<template>
  <div class="space-y-6">
    <div class="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
      <div>
        <p class="text-sm font-semibold uppercase tracking-wide text-emerald-600">Visualisasi Jabatan</p>
        <h1 class="text-3xl font-bold text-gray-900">Hierarki Jabatan</h1>
        <p class="mt-1 text-sm text-gray-500">Gambaran struktur organisasi berdasarkan induk dan jabatan turunan.</p>
      </div>
      <RouterLink
        :to="{ name: 'EmployeePositions' }"
        class="inline-flex items-center gap-2 rounded-lg border border-gray-200 px-4 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 transition"
      >
        <i class="pi pi-arrow-left text-sm"></i>
        Kembali ke Matriks Jabatan
      </RouterLink>
    </div>

    <div class="rounded-2xl border border-gray-100 bg-white p-6 shadow-sm">
      <div class="flex flex-col gap-2 md:flex-row md:items-center md:justify-between">
        <div>
          <h2 class="text-lg font-semibold text-gray-900">Struktur Organisasi</h2>
          <p class="text-sm text-gray-500">
            Setiap kartu mewakili jabatan. Hubungkan jabatan induk dan anak melalui halaman matriks jabatan.
          </p>
        </div>
        <p v-if="positionTree.length" class="text-sm text-gray-400">
          Total jabatan dalam visualisasi: {{ totalPositions }}
        </p>
      </div>

      <div v-if="positionTree.length" class="mt-6 overflow-x-auto">
        <div class="min-w-[480px] pb-6">
          <PositionTree :nodes="positionTree" is-root />
        </div>
      </div>
      <div
        v-else
        class="mt-6 rounded-xl border border-dashed border-gray-200 bg-gray-50 p-6 text-center text-sm text-gray-500"
      >
        Belum ada jabatan yang memiliki induk. Tambahkan data melalui halaman matriks jabatan untuk melihat visualisasi.
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { useHRStore } from '@/stores/hr'
import PositionTree from '@/components/PositionTree.vue'

const hrStore = useHRStore()
onMounted(() => hrStore.ensureHydrated())

const positionTree = computed(() => {
  if (!hrStore.positions.length) return []

  const nodes = hrStore.positions.map((position) => ({
    ...position,
    children: []
  }))

  const lookup = new Map(nodes.map((node) => [node.id, node]))
  const roots = []

  nodes.forEach((node) => {
    if (node.parentId && lookup.has(node.parentId)) {
      lookup.get(node.parentId).children.push(node)
    } else {
      roots.push(node)
    }
  })

  const sortNodes = (list) => {
    list.sort((a, b) => a.title.localeCompare(b.title))
    list.forEach((child) => sortNodes(child.children))
  }

  sortNodes(roots)
  return roots
})

const totalPositions = computed(() => hrStore.positions.length)
</script>
