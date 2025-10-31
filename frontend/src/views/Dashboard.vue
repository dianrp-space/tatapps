<template>
  <div>
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-3xl font-bold">Dashboard</h1>
        <p class="text-gray-500 mt-1">Ringkasan kondisi operasional terkini</p>
      </div>
      <div class="flex gap-3">
        <div class="relative group">
          <Button
            label="Cek Low Stock"
            icon="pi pi-bell"
            class="p-button-warning"
            :loading="checkingStock"
            @click="checkLowStock"
          />
          <div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-xs rounded-lg opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none z-10">
            Notifikasi akan dikirim ke nomor/email profile Anda
            <div class="absolute top-full left-1/2 transform -translate-x-1/2 -mt-1">
              <div class="border-4 border-transparent border-t-gray-900"></div>
            </div>
          </div>
        </div>
        <Button
          label="Refresh"
          icon="pi pi-refresh"
          class="p-button-outlined"
          :loading="refreshing"
          @click="refreshDashboard"
        />
      </div>
    </div>

    <div
      v-if="loading"
      class="bg-white border border-gray-200 rounded-xl p-6 mb-6 shadow-sm"
    >
      <div class="flex items-center gap-3 text-gray-500">
        <i class="pi pi-spin pi-spinner text-lg"></i>
        <span>Memuat data dashboard...</span>
      </div>
    </div>

    <div
      v-else
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6"
    >
      <div
        v-for="card in summaryCards"
        :key="card.title"
        class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 hover:shadow-md transition-shadow"
      >
        <div class="flex items-start justify-between">
          <div>
            <p class="text-gray-500 text-sm">{{ card.title }}</p>
            <h3 class="text-3xl font-bold mt-2 text-gray-800">
              {{ card.value }}
            </h3>
            <p v-if="card.caption" class="text-xs text-gray-400 mt-1">
              {{ card.caption }}
            </p>
          </div>
          <div :class="card.iconWrapper">
            <i :class="card.icon"></i>
          </div>
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
      <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 h-full">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800">Statistik Karyawan</h3>
          <span v-if="employeeStats.total" class="text-xs text-gray-400">
            Total {{ employeeStats.total }} orang
          </span>
        </div>
        <div
          v-if="employeeStats.total"
          class="grid grid-cols-1 sm:grid-cols-2 gap-4"
        >
          <div class="p-4 rounded-lg border border-gray-100 bg-gray-50">
            <p class="text-xs uppercase tracking-wide text-gray-500 mb-1">
              Aktif
            </p>
            <p class="text-2xl font-semibold text-emerald-600">
              {{ employeeStats.active }}
            </p>
          </div>
          <div class="p-4 rounded-lg border border-gray-100 bg-gray-50">
            <p class="text-xs uppercase tracking-wide text-gray-500 mb-1">
              Cuti
            </p>
            <p class="text-2xl font-semibold text-amber-600">
              {{ employeeStats.onLeave }}
            </p>
          </div>
          <div class="p-4 rounded-lg border border-gray-100 bg-gray-50">
            <p class="text-xs uppercase tracking-wide text-gray-500 mb-1">
              Resign
            </p>
            <p class="text-2xl font-semibold text-rose-600">
              {{ employeeStats.resigned }}
            </p>
          </div>
          <div class="p-4 rounded-lg border border-gray-100 bg-gray-50">
            <p class="text-xs uppercase tracking-wide text-gray-500 mb-1">
              Bergabung Bulan Ini
            </p>
            <p class="text-2xl font-semibold text-indigo-600">
              {{ employeeStats.newJoiners }}
            </p>
          </div>
        </div>
        <div v-else class="text-sm text-gray-500">
          Belum ada data karyawan yang tersimpan di modul Employees.
        </div>
        <div
          v-if="employeeStats.averageAge"
          class="mt-5 flex items-center gap-2 text-sm text-gray-600"
        >
          <i class="pi pi-chart-line text-indigo-500"></i>
          Rata-rata usia {{ employeeStats.averageAge }} tahun
        </div>
        <p class="mt-3 text-xs text-gray-400">
          Data diambil dari menu Employees.
        </p>
      </div>

      <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 h-full w-full relative">
        <div class="flex items-center justify-between mb-3">
          <h3 class="text-lg font-semibold text-gray-800">Kalender</h3>
          <span class="text-xs text-gray-400">{{ calendarReadableDate }}</span>
        </div>
        <div class="flex items-center justify-between mb-4 text-sm font-medium text-gray-600">
          <button
            type="button"
            class="p-2 rounded-full border border-gray-200 hover:bg-gray-50 transition-colors"
            @click="goToPrevMonth"
            aria-label="Bulan sebelumnya"
          >
            <i class="pi pi-chevron-left text-xs"></i>
          </button>
          <span class="text-gray-800 capitalize tracking-wide">{{ calendarMonthLabel }}</span>
          <button
            type="button"
            class="p-2 rounded-full border border-gray-200 hover:bg-gray-50 transition-colors"
            @click="goToNextMonth"
            aria-label="Bulan berikutnya"
          >
            <i class="pi pi-chevron-right text-xs"></i>
          </button>
        </div>
        <div class="overflow-x-auto">
          <div class="grid grid-cols-8 gap-2 text-center text-[11px] font-semibold uppercase tracking-wide text-gray-400 mb-2">
            <div>Wk</div>
            <div v-for="dayLabel in weekDayLabels" :key="dayLabel">{{ dayLabel }}</div>
          </div>
          <div class="space-y-2">
            <div
              v-for="week in calendarWeeks"
              :key="week.id"
              class="grid grid-cols-8 gap-2 items-center"
            >
              <div class="text-xs text-gray-400 font-medium">{{ week.weekNumber }}</div>
              <button
                v-for="day in week.days"
                :key="day.key"
                type="button"
                class="relative flex items-center justify-center h-10 rounded-full text-sm font-medium transition-colors duration-150"
                :class="[
                  day.inCurrentMonth ? 'text-gray-700 hover:bg-indigo-50' : 'text-gray-300 hover:bg-gray-100',
                  day.isToday && !day.isSelected ? 'border border-indigo-500' : '',
                  day.isSelected ? 'bg-indigo-500 text-white hover:bg-indigo-600 shadow-md' : ''
                ]"
                v-tooltip.top="day.birthdays.length ? birthdayTooltip(day.date) : ''"
                @click="handleDayClick(day)"
              >
                <span>{{ day.dayNumber }}</span>
                <span
                  v-if="day.birthdays.length"
                  class="absolute bottom-1 left-2 right-2 h-1 rounded-full"
                  :class="day.isSelected ? 'bg-white/80' : 'bg-pink-500'"
                ></span>
              </button>
            </div>
          </div>
        </div>
        <div v-if="upcomingBirthdays.length" class="mt-6">
          <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center gap-2">
            <i class="pi pi-gift text-pink-500"></i>
            Ulang Tahun Terdekat
          </h4>
          <ul class="space-y-2 text-sm text-gray-600">
            <li
              v-for="birthday in upcomingBirthdays"
              :key="birthday.id"
              class="flex items-center justify-between bg-gray-50 border border-gray-100 rounded-lg px-3 py-2"
            >
              <div>
                <p class="font-medium text-gray-800">{{ birthday.fullName }}</p>
                <p class="text-xs text-gray-400">{{ birthday.note }}</p>
              </div>
              <span class="text-xs text-gray-500">{{ birthday.formattedDate }}</span>
            </li>
          </ul>
        </div>
        <div v-else class="mt-6 text-sm text-gray-500">
          Belum ada jadwal ulang tahun dalam 30 hari ke depan.
        </div>
      </div>

      <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100 h-full">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800">Karyawan per Divisi</h3>
          <span v-if="divisionStats.length" class="text-xs text-gray-400">
            {{ divisionStats.length }} divisi
          </span>
        </div>
        <div v-if="divisionStats.length" class="space-y-4">
          <div
            v-for="division in divisionStats"
            :key="division.id"
            class="space-y-2"
          >
            <div class="flex items-center justify-between text-sm text-gray-600">
              <span class="font-medium text-gray-800">{{ division.name }}</span>
              <span class="text-xs text-gray-500">{{ division.count }} karyawan</span>
            </div>
            <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
              <div
                class="h-2 rounded-full bg-emerald-500 transition-all"
                :style="{ width: computeBarWidth(division.count, maxDivisionCount) }"
              ></div>
            </div>
          </div>
        </div>
        <div v-else class="text-sm text-gray-500">
          Belum ada data divisi atau karyawan untuk ditampilkan.
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
      <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800">Karyawan Berdasarkan Umur</h3>
          <span class="text-xs text-gray-400">Menggunakan data tanggal lahir</span>
        </div>
        <div v-if="employeeStats.total && maxAgeGroupCount" class="space-y-4">
          <div v-for="group in ageGroups" :key="group.key">
            <div class="flex items-center justify-between text-xs text-gray-500 mb-1">
              <span>{{ group.label }}</span>
              <span>{{ group.count }} karyawan</span>
            </div>
            <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
              <div
                class="h-2 rounded-full bg-indigo-500 transition-all"
                :style="{ width: computeBarWidth(group.count, maxAgeGroupCount) }"
              ></div>
            </div>
          </div>
        </div>
        <div v-else class="text-sm text-gray-500">
          Belum ada data tanggal lahir karyawan untuk ditampilkan.
        </div>
      </div>

      <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800">Karyawan Berdasarkan Gender</h3>
          <span v-if="employeeStats.total" class="text-xs text-gray-400">
            Total {{ employeeStats.total }} karyawan
          </span>
        </div>
        <div v-if="genderStats.length" class="space-y-4">
          <div
            v-for="gender in genderStats"
            :key="gender.key"
            class="space-y-2"
          >
            <div class="flex items-center justify-between text-sm text-gray-600">
              <div class="flex items-center gap-3">
                <span :class="['w-3 h-3 rounded-full', gender.color]"></span>
                <span class="font-medium text-gray-800">{{ gender.label }}</span>
              </div>
              <span class="text-xs text-gray-500">
                {{ gender.count }} karyawan
              </span>
            </div>
            <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
              <div
                class="h-2 rounded-full transition-all"
                :class="gender.color"
                :style="{ width: computeBarWidth(gender.percentage, 100) }"
              ></div>
            </div>
            <p class="text-xs text-gray-400 text-right">
              {{ gender.percentage }}%
            </p>
          </div>
        </div>
        <div v-else class="text-sm text-gray-500">
          Belum ada data gender karyawan.
        </div>
      </div>

      <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800">Karyawan Berdasarkan Generasi</h3>
          <span v-if="employeeStats.total" class="text-xs text-gray-400">
            Total {{ employeeStats.total }} karyawan
          </span>
        </div>
        <div v-if="generationStats.length" class="space-y-4">
          <div
            v-for="generation in generationStats"
            :key="generation.key"
            class="space-y-2"
          >
            <div class="flex items-center justify-between text-sm text-gray-600">
              <span class="font-medium text-gray-800">{{ generation.label }}</span>
              <span class="text-xs text-gray-500">{{ generation.count }} karyawan</span>
            </div>
            <div class="h-2 bg-gray-100 rounded-full overflow-hidden">
              <div
                class="h-2 rounded-full bg-cyan-500 transition-all"
                :style="{ width: computeBarWidth(generation.count, maxGenerationCount) }"
              ></div>
            </div>
            <p class="text-xs text-gray-400 text-right">
              {{ generation.percentage }}%
            </p>
          </div>
        </div>
        <div v-else class="text-sm text-gray-500">
          Belum ada data generasi karyawan.
        </div>
      </div>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800">Aktivitas Terbaru</h3>
          <span class="text-xs text-gray-400"
            >Menampilkan {{ recentTransactions.length }} transaksi
            terakhir</span
          >
        </div>

        <div
          v-if="recentTransactions.length === 0"
          class="text-gray-500 text-sm"
        >
          Tidak ada aktivitas terbaru.
        </div>

        <div v-else class="space-y-4">
          <div
            v-for="tx in recentTransactions"
            :key="tx.id"
            class="flex items-start gap-3 p-3 rounded-lg border border-gray-100 hover:border-blue-200 transition-colors"
          >
            <div :class="['p-2 rounded-full', transactionIconWrapper(tx.type)]">
              <i :class="transactionIcon(tx.type)"></i>
            </div>
            <div class="flex-1">
              <div class="flex items-center justify-between">
                <p class="text-sm font-medium text-gray-800">
                  {{ transactionTitle(tx) }}
                </p>
                <span class="text-xs text-gray-400">
                  {{ formatDateTime(tx.created_at) }}
                </span>
              </div>
              <p class="text-xs text-gray-500 mt-1">
                {{ transactionSubtitle(tx) }}
              </p>
              <div
                v-if="transactionActor(tx)"
                class="flex items-center gap-2 text-xs text-gray-400 mt-1"
              >
                <img
                  :src="transactionActorAvatar(tx, 28)"
                  alt=""
                  class="w-7 h-7 rounded-full object-cover border border-gray-200"
                />
                <span>Oleh: {{ transactionActor(tx) }}</span>
              </div>
              <p v-if="tx.notes" class="text-xs text-gray-400 mt-1">
                Catatan: {{ tx.notes }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-100">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-800">Low Stock Alerts</h3>
          <span class="text-xs text-gray-400"
            >{{ lowStockItems.length }} item butuh perhatian</span
          >
        </div>

        <div v-if="lowStockItems.length === 0" class="text-gray-500 text-sm">
          Semua stok berada di atas batas aman.
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="item in lowStockItems"
            :key="item.id || item.name"
            class="flex items-center justify-between p-3 rounded-lg border border-red-100 bg-red-50"
          >
            <div>
              <p class="text-sm font-medium text-red-700">
                {{ item.name || item.category || 'Item' }}
              </p>
              <p class="text-xs text-gray-500">
                {{ item.warehouse_name || item.warehouse?.name || 'Tidak ada gudang' }}
              </p>
              <p class="text-xs text-gray-400 mt-1">
                Tersedia:
                {{ formatQuantityDisplay(item.quantity) }}
                {{ item.unit || 'unit' }}
                • Minimum:
                {{ formatQuantityDisplay(item.min_stock) }}
              </p>
              <p
                v-if="item.aggregated"
                class="text-[11px] text-red-500 mt-1"
              >
                Akumulasi kategori {{ item.category || '-' }}
                ({{ item.item_ids?.length || 0 }} item SN)
              </p>
            </div>
            <i class="pi pi-exclamation-triangle text-red-600 text-lg"></i>
          </div>
        </div>
      </div>
  </div>

    <!-- Modal Pilih Channel Notifikasi -->
    <div
      v-if="showChannelModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click.self="showChannelModal = false"
    >
      <div class="bg-white rounded-xl shadow-2xl max-w-md w-full mx-4 p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-xl font-bold text-gray-800">Pilih Channel Notifikasi</h3>
          <button
            @click="showChannelModal = false"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <i class="pi pi-times text-xl"></i>
          </button>
        </div>

        <p class="text-gray-600 text-sm mb-4">
          Pilih channel untuk menerima notifikasi low stock alert:
        </p>

        <div class="space-y-3 mb-6">
          <!-- WhatsApp Option -->
          <label
            class="flex items-center p-4 border-2 rounded-lg cursor-pointer transition-colors"
            :class="selectedChannels.whatsapp ? 'border-green-500 bg-green-50' : 'border-gray-200 hover:border-green-300'"
          >
            <input
              type="checkbox"
              v-model="selectedChannels.whatsapp"
              class="w-5 h-5 text-green-600 border-gray-300 rounded focus:ring-green-500"
            />
            <div class="ml-3 flex-1">
              <div class="flex items-center gap-2">
                <i class="pi pi-whatsapp text-green-600 text-lg"></i>
                <span class="font-medium text-gray-800">WhatsApp</span>
              </div>
              <p class="text-xs text-gray-500 mt-1">
                Kirim ke: {{ userPhone || 'Nomor belum diset' }}
              </p>
            </div>
          </label>

          <!-- Email Option -->
          <label
            class="flex items-center p-4 border-2 rounded-lg cursor-pointer transition-colors"
            :class="selectedChannels.email ? 'border-blue-500 bg-blue-50' : 'border-gray-200 hover:border-blue-300'"
          >
            <input
              type="checkbox"
              v-model="selectedChannels.email"
              class="w-5 h-5 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
            />
            <div class="ml-3 flex-1">
              <div class="flex items-center gap-2">
                <i class="pi pi-envelope text-blue-600 text-lg"></i>
                <span class="font-medium text-gray-800">Email</span>
              </div>
              <p class="text-xs text-gray-500 mt-1">
                Kirim ke: {{ userEmail || 'Email belum diset' }}
              </p>
            </div>
          </label>
        </div>

        <div
          v-if="!selectedChannels.whatsapp && !selectedChannels.email"
          class="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg"
        >
          <p class="text-xs text-yellow-800">
            <i class="pi pi-exclamation-triangle mr-1"></i>
            Pilih minimal satu channel untuk melanjutkan
          </p>
        </div>

        <div class="flex gap-3">
          <button
            @click="showChannelModal = false"
            class="flex-1 px-4 py-2 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50 transition-colors"
          >
            Batal
          </button>
          <button
            @click="confirmSendNotification"
            :disabled="!selectedChannels.whatsapp && !selectedChannels.email || checkingStock"
            class="flex-1 px-4 py-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600 transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            <i class="pi pi-spin pi-spinner" v-if="checkingStock"></i>
            <i class="pi pi-send" v-else></i>
            {{ checkingStock ? 'Mengirim...' : 'Kirim Notifikasi' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Modal Detail Ulang Tahun -->
    <div
      v-if="showBirthdayModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 px-4"
      @click.self="closeBirthdayModal"
    >
      <div class="bg-white w-full max-w-md rounded-2xl shadow-2xl overflow-hidden">
        <div class="flex items-center justify-between px-5 py-4 border-b border-gray-100">
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Ulang Tahun</h3>
            <p class="text-xs text-gray-500 mt-1">
              {{ birthdayModalTitle }}
            </p>
          </div>
          <button
            type="button"
            class="p-2 rounded-full text-gray-400 hover:text-gray-600 hover:bg-gray-100 transition-colors"
            @click="closeBirthdayModal"
            aria-label="Tutup"
          >
            <i class="pi pi-times text-sm"></i>
          </button>
        </div>
        <div class="px-5 py-4 space-y-4">
          <div class="flex items-center justify-between text-sm text-gray-600">
            <span>Total Karyawan</span>
            <span class="font-semibold text-gray-900">{{ birthdayModalCount }} orang</span>
          </div>
          <ul class="divide-y divide-gray-100 rounded-xl border border-gray-100">
            <li
              v-for="(employee, index) in birthdayModal.employees"
              :key="employee.id || index"
              class="px-4 py-3 text-sm text-gray-700 flex items-center gap-3"
            >
              <span
                class="flex h-10 w-10 items-center justify-center overflow-hidden rounded-full border border-gray-200 bg-indigo-50 text-indigo-600 text-xs font-semibold"
              >
                <img
                  v-if="employee.photo"
                  :src="employee.photo"
                  :alt="employee.name"
                  class="h-full w-full object-cover"
                />
                <span v-else>{{ initials(employee.name) }}</span>
              </span>
              <div class="flex-1">
                <p class="font-medium text-gray-800">{{ employee.name }}</p>
              </div>
            </li>
          </ul>
        </div>
        <div class="px-5 py-4 bg-gray-50 flex justify-end">
          <button
            type="button"
            class="px-4 py-2 text-sm font-medium text-gray-600 rounded-lg border border-gray-200 hover:bg-white transition-colors"
            @click="closeBirthdayModal"
          >
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref, computed, watch } from "vue";
import { useToast } from "primevue/usetoast";
import { useAuthStore } from "@/stores/auth";
import { useHRStore } from "@/stores/hr";
import axios from "@/api/axios";

const toast = useToast();
const authStore = useAuthStore();
const hrStore = useHRStore();

const loading = ref(true);
const refreshing = ref(false);
const checkingStock = ref(false);
const showChannelModal = ref(false);
const selectedChannels = ref({
  whatsapp: true,
  email: true
});

const userPhone = computed(() => authStore.user?.phone || '');
const userEmail = computed(() => authStore.user?.email || '');
const calendarValue = ref(new Date());
const calendarViewDate = ref(startOfMonth(calendarValue.value));
const weekDayLabels = ["Sen", "Sel", "Rab", "Kam", "Jum", "Sab", "Min"];
const showBirthdayModal = ref(false);
const birthdayModal = ref({
  date: null,
  employees: [],
});

const summary = ref({
  warehouses: 0,
  items: 0,
  categories: 0,
  inventoryValue: 0,
});

const recentTransactions = ref([]);
const lowStockItems = ref([]);

const employees = computed(() => hrStore.employees || []);

const employeeStats = computed(() => {
  const total = employees.value.length;
  const statuses = employees.value.reduce(
    (acc, employee) => {
      const status = normalizeStatus(employee.status);
      if (status === "aktif") {
        acc.active += 1;
      } else if (status === "cuti") {
        acc.onLeave += 1;
      } else if (status === "resign" || status === "resigned") {
        acc.resigned += 1;
      }
      return acc;
    },
    { active: 0, onLeave: 0, resigned: 0 },
  );

  const now = new Date();
  const newJoiners = employees.value.filter((employee) => {
    const joinDate = parseDate(employee.joinDate);
    if (!joinDate) return false;
    return (
      joinDate.getMonth() === now.getMonth() &&
      joinDate.getFullYear() === now.getFullYear()
    );
  }).length;

  const ages = employees.value
    .map((employee) => calculateAge(employee.birthDate))
    .filter((age) => typeof age === "number");

  const averageAge =
    ages.length > 0
      ? Math.round(
          (ages.reduce((totalAge, age) => totalAge + age, 0) / ages.length) *
            10,
        ) / 10
      : null;

  return {
    total,
    active: statuses.active,
    onLeave: statuses.onLeave,
    resigned: statuses.resigned,
    newJoiners,
    averageAge,
  };
});

const ageGroups = computed(() => {
  const groups = [
    { key: "under25", label: "< 25 tahun", min: 0, max: 24, count: 0 },
    { key: "25to34", label: "25 - 34 tahun", min: 25, max: 34, count: 0 },
    { key: "35to44", label: "35 - 44 tahun", min: 35, max: 44, count: 0 },
    { key: "45to54", label: "45 - 54 tahun", min: 45, max: 54, count: 0 },
    { key: "55plus", label: "55+ tahun", min: 55, max: Infinity, count: 0 },
  ];

  employees.value.forEach((employee) => {
    const age = calculateAge(employee.birthDate);
    if (typeof age !== "number") return;
    const bucket = groups.find(
      (group) => age >= group.min && age <= group.max,
    );
    if (bucket) {
      bucket.count += 1;
    }
  });

  return groups;
});

const maxAgeGroupCount = computed(() =>
  ageGroups.value.reduce((max, group) => Math.max(max, group.count), 0),
);

const genderStats = computed(() => {
  const counters = { male: 0, female: 0, other: 0, unspecified: 0 };

  employees.value.forEach((employee) => {
    const normalized = normalizeGender(employee.gender);
    counters[normalized] += 1;
  });

  const total = employees.value.length || 0;
  const baseStats = [
    { key: "male", label: "Pria", color: "bg-blue-500" },
    { key: "female", label: "Wanita", color: "bg-pink-500" },
    { key: "other", label: "Lainnya", color: "bg-purple-500" },
    { key: "unspecified", label: "Belum Diisi", color: "bg-gray-400" },
  ];

  return baseStats
    .filter((item) => counters[item.key] > 0)
    .map((item) => ({
      ...item,
      count: counters[item.key],
      percentage: total > 0 ? Math.round((counters[item.key] / total) * 100) : 0,
    }));
});

const divisionStats = computed(() => {
  const counts = new Map();
  const divisionLookup = new Map(
    (hrStore.divisions || []).map((division) => [division.id, division.name]),
  );

  employees.value.forEach((employee) => {
    const divisionId = employee.divisionId || employee.division?.id || employee.division;
    if (!divisionId) return;
    const current = counts.get(divisionId) || { id: divisionId, name: divisionLookup.get(divisionId) || "Divisi Lainnya", count: 0 };
    current.count += 1;
    counts.set(divisionId, current);
  });

  return Array.from(counts.values()).sort((a, b) => b.count - a.count);
});

const maxDivisionCount = computed(() =>
  divisionStats.value.reduce((max, division) => Math.max(max, division.count), 0),
);

const generationStats = computed(() => {
  const generations = [
    { key: "babyBoomers", label: "Baby Boomers (1946-1964)", min: 1946, max: 1964 },
    { key: "genX", label: "Generasi X (1965-1980)", min: 1965, max: 1980 },
    { key: "millennial", label: "Millennial (1981-1996)", min: 1981, max: 1996 },
    { key: "genZ", label: "Generasi Z (1997-2012)", min: 1997, max: 2012 },
    { key: "genAlpha", label: "Generasi Alpha (2013+)", min: 2013, max: Infinity },
  ];

  const counters = generations.map((generation) => ({
    ...generation,
    count: 0,
  }));

  employees.value.forEach((employee) => {
    const birthDate = parseDate(employee.birthDate);
    if (!birthDate) return;
    const year = birthDate.getFullYear();
    const bucket = counters.find((generation) => year >= generation.min && year <= generation.max);
    if (bucket) {
      bucket.count += 1;
    }
  });

  const total = employees.value.length || 0;

  return counters
    .filter((generation) => generation.count > 0)
    .map((generation) => ({
      key: generation.key,
      label: generation.label,
      count: generation.count,
      percentage: total > 0 ? Math.round((generation.count / total) * 100) : 0,
    }));
});

const maxGenerationCount = computed(() =>
  generationStats.value.reduce((max, generation) => Math.max(max, generation.count), 0),
);

const birthdayLookup = computed(() => {
  const map = new Map();
  employees.value.forEach((employee) => {
    const birthDate = parseDate(employee.birthDate);
    if (!birthDate) return;
    const key = `${birthDate.getMonth()}-${birthDate.getDate()}`;
    if (!map.has(key)) {
      map.set(key, []);
    }
    map.get(key).push({
      id:
        employee.id ||
        employee.nik ||
        employee.employeeCode ||
        `${(employee.fullName || "karyawan").replace(/\s+/g, "-")}-${birthDate.getTime()}`,
      name: employee.fullName || employee.employeeCode || "Karyawan",
      photo: normalizeEmployeePhoto(employee.photo),
    });
  });
  return map;
});

const calendarMonthLabel = computed(() =>
  formatDateFriendly(calendarViewDate.value, {
    month: "long",
    year: "numeric",
  }),
);

const calendarWeeks = computed(() => {
  const base = startOfMonth(calendarViewDate.value);
  const month = base.getMonth();

  const start = new Date(base);
  const offset = (start.getDay() + 6) % 7;
  start.setDate(start.getDate() - offset);

  const weeks = [];
  for (let weekIndex = 0; weekIndex < 6; weekIndex += 1) {
    const weekStart = new Date(start);
    weekStart.setDate(start.getDate() + weekIndex * 7);

    const days = [];
    for (let dayIndex = 0; dayIndex < 7; dayIndex += 1) {
      const current = new Date(weekStart);
      current.setDate(weekStart.getDate() + dayIndex);

      const birthdayKeyValue = `${current.getMonth()}-${current.getDate()}`;
      const birthdays =
        birthdayLookup.value.get(birthdayKeyValue)?.map((entry) => ({ ...entry })) || [];

      days.push({
        date: current,
        dayNumber: current.getDate(),
        inCurrentMonth: current.getMonth() === month,
        isToday: isSameDate(current, new Date()),
        isSelected: isSameDate(current, calendarValue.value),
        birthdays,
        key: `${current.getFullYear()}-${current.getMonth()}-${current.getDate()}`,
      });
    }

    weeks.push({
      id: `${weekStart.getFullYear()}-${weekStart.getMonth()}-${weekStart.getDate()}`,
      weekNumber: String(getISOWeekNumber(weekStart)).padStart(2, "0"),
      days,
    });
  }

  return weeks.filter((week) => week.days.some((day) => day.inCurrentMonth));
});

const upcomingBirthdays = computed(() => {
  if (employees.value.length === 0) return [];
  const today = new Date();

  return employees.value
    .map((employee) => {
      const birthDate = parseDate(employee.birthDate);
      if (!birthDate) return null;

      const nextBirthdayYear =
        birthDate.getMonth() < today.getMonth() ||
        (birthDate.getMonth() === today.getMonth() &&
          birthDate.getDate() < today.getDate())
          ? today.getFullYear() + 1
          : today.getFullYear();

      const nextBirthday = new Date(
        nextBirthdayYear,
        birthDate.getMonth(),
        birthDate.getDate(),
      );

      const diffMilliseconds = nextBirthday.getTime() - today.getTime();
      const diffDays = Math.ceil(diffMilliseconds / (1000 * 60 * 60 * 24));

      return {
        id: employee.id,
        fullName: employee.fullName || employee.employeeCode || "Karyawan",
        nextBirthday,
        diffDays,
      };
    })
    .filter(
      (entry) =>
        entry && entry.diffDays >= 0 && entry.diffDays <= 30,
    )
    .sort((a, b) => a.diffDays - b.diffDays)
    .slice(0, 5)
    .map((entry) => ({
      ...entry,
      formattedDate: formatDateFriendly(entry.nextBirthday, {
        month: "long",
        day: "numeric",
      }),
      note:
        entry.diffDays === 0
          ? "Hari ini"
          : `${entry.diffDays} hari lagi`,
    }));
});

const calendarReadableDate = computed(() =>
  formatDateFriendly(calendarValue.value, { dateStyle: "full" }),
);

const birthdayModalTitle = computed(() =>
  birthdayModal.value.date
    ? formatDateFriendly(birthdayModal.value.date, {
        dateStyle: "full",
      })
    : "-",
);

const birthdayModalCount = computed(
  () => birthdayModal.value.employees.length || 0,
);

const apiBaseUrl = import.meta.env.VITE_API_URL || "http://localhost:8080/api/v1";

const fileBaseUrl = (() => {
  try {
    const url = new URL(apiBaseUrl);
    url.pathname = url.pathname.replace(/\/api\/v1\/?$/, "");
    const path = url.pathname.replace(/\/$/, "");
    return `${url.origin}${path === "" ? "" : path}`;
  } catch {
    return apiBaseUrl.replace(/\/api\/v1\/?$/, "");
  }
})();

const summaryCards = computed(() => [
  {
    title: "Employees",
    value: employeeStats.value.total,
    caption:
      employeeStats.value.total > 0
        ? `Aktif ${employeeStats.value.active} • Cuti ${employeeStats.value.onLeave}`
        : "",
    icon: "pi pi-users text-indigo-600 text-2xl",
    iconWrapper: "bg-indigo-100 p-3 rounded-full",
  },
  {
    title: "Warehouses",
    value: summary.value.warehouses,
    caption: "",
    icon: "pi pi-building text-blue-600 text-2xl",
    iconWrapper: "bg-blue-100 p-3 rounded-full",
  },
  {
    title: "Inventory Items",
    value: summary.value.items,
    caption:
      summary.value.inventoryValue > 0
        ? `Total value ${formatCurrency(summary.value.inventoryValue)}`
        : "",
    icon: "pi pi-box text-green-600 text-2xl",
    iconWrapper: "bg-green-100 p-3 rounded-full",
  },
  {
    title: "Categories",
    value: summary.value.categories,
    caption: "",
    icon: "pi pi-tags text-purple-600 text-2xl",
    iconWrapper: "bg-purple-100 p-3 rounded-full",
  },
]);

onMounted(() => {
  hrStore.ensureHydrated();
  fetchDashboard();
});

watch(calendarValue, (value) => {
  const parsed = parseDate(value);
  if (parsed) {
    calendarViewYearMonth(parsed);
  }
});

async function refreshDashboard() {
  refreshing.value = true;
  await fetchDashboard();
  refreshing.value = false;
}

function checkLowStock() {
  // Validasi apakah user sudah set phone atau email
  if (!userPhone.value && !userEmail.value) {
    toast.add({
      severity: 'warn',
      summary: 'Data Belum Lengkap',
      detail: 'Silakan lengkapi nomor WhatsApp atau Email di Profile terlebih dahulu',
      life: 5000
    });
    return;
  }

  // Set default channel berdasarkan data yang tersedia
  selectedChannels.value = {
    whatsapp: !!userPhone.value,
    email: !!userEmail.value
  };

  // Tampilkan modal pilih channel
  showChannelModal.value = true;
}

async function confirmSendNotification() {
  checkingStock.value = true;
  try {
    const response = await axios.post('/notifications/check-low-stock', {
      send_whatsapp: selectedChannels.value.whatsapp,
      send_email: selectedChannels.value.email
    });
    
    // Tutup modal
    showChannelModal.value = false;
    
    if (response.data.low_stock_count === 0) {
      toast.add({
        severity: 'success',
        summary: 'Stok Aman',
        detail: 'Semua item memiliki stok yang cukup',
        life: 5000
      });
    } else {
      const channels = [];
      if (response.data.whatsapp_sent) channels.push('WhatsApp');
      if (response.data.email_sent) channels.push('Email');
      
      const detail = `Ditemukan ${response.data.low_stock_count} item dengan stok rendah.\n` +
                     `Notifikasi telah dikirim ke ${channels.join(' dan ')} Anda.`;
      
      toast.add({
        severity: 'warn',
        summary: 'Low Stock Alert',
        detail: detail,
        life: 6000
      });
    }
    
    // Refresh data low stock
    const lowStockRes = await axios.get("/inventory/low-stock");
    lowStockItems.value = (lowStockRes.data?.data || []).sort(
      (a, b) => a.quantity - b.quantity
    );
  } catch (error) {
    console.error('Check low stock error:', error);
    const errorDetail = error.response?.data?.error || error.response?.data?.message || 'Gagal mengecek low stock';
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: errorDetail,
      life: 5000
    });
  } finally {
    checkingStock.value = false;
  }
}

async function fetchDashboard() {
  loading.value = true;
  try {
    const [
      warehousesRes,
      inventoryRes,
      categoriesRes,
      lowStockRes,
      transactionsRes,
    ] = await Promise.all([
      axios.get("/warehouses"),
      axios.get("/inventory"),
      axios.get("/categories"),
      axios.get("/inventory/low-stock"),
      axios.get("/inventory/transactions"),
    ]);

    const warehouses = warehousesRes.data?.data || [];
    const inventoryItems = inventoryRes.data?.data || [];
    const categories = categoriesRes.data?.data || [];
    const lowStock = lowStockRes.data?.data || [];
    const transactions = transactionsRes.data?.data || [];

    summary.value = {
      warehouses: warehouses.length,
      items: inventoryItems.length,
      categories: categories.length,
      inventoryValue: inventoryItems.reduce(
        (total, item) => total + (item.quantity || 0) * (item.unit_price || 0),
        0,
      ),
    };

    lowStockItems.value = lowStock.sort(
      (a, b) => (a.quantity || 0) - (b.quantity || 0),
    );

    recentTransactions.value = transactions
      .slice()
      .sort(
        (a, b) =>
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime(),
      )
      .slice(0, 5);
  } catch (error) {
    console.error("Dashboard load error:", error);
    toast.add({
      severity: "error",
      summary: "Gagal memuat dashboard",
      detail:
        error.response?.data?.message ||
        error.response?.data?.error ||
        error.message ||
        "Terjadi kesalahan tak terduga",
      life: 4000,
    });
  } finally {
    loading.value = false;
  }
}

function formatCurrency(value) {
  if (!value) return "Rp0";
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(value);
}

function formatQuantityDisplay(value) {
  if (value === undefined || value === null) {
    return 0;
  }

  const numeric = Number(value);
  if (Number.isNaN(numeric)) {
    return 0;
  }

  if (Math.abs(numeric-Math.round(numeric)) < 1e-9) {
    return Math.round(numeric);
  }

  return parseFloat(numeric.toFixed(2));
}

function formatDateTime(value) {
  if (!value) return "-";
  const date = new Date(value);
  return new Intl.DateTimeFormat("id-ID", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

function transactionIcon(type) {
  switch ((type || "").toLowerCase()) {
    case "in":
      return "pi pi-arrow-down text-green-600";
    case "out":
      return "pi pi-arrow-up text-red-600";
    case "transfer":
      return "pi pi-arrows-h text-blue-600";
    case "adjustment":
      return "pi pi-sliders-h text-purple-600";
    default:
      return "pi pi-info-circle text-gray-500";
  }
}

function transactionIconWrapper(type) {
  switch ((type || "").toLowerCase()) {
    case "in":
      return "bg-green-100";
    case "out":
      return "bg-red-100";
    case "transfer":
      return "bg-blue-100";
    case "adjustment":
      return "bg-purple-100";
    default:
      return "bg-gray-100";
  }
}

function transactionTitle(tx) {
  const base = tx.item?.name || "Item";
  switch ((tx.type || "").toLowerCase()) {
    case "in":
      return `Stock In - ${base}`;
    case "out":
      return `Stock Out - ${base}`;
    case "transfer":
      return `Transfer - ${base}`;
    case "adjustment":
      return `Adjustment - ${base}`;
    default:
      return base;
  }
}

function transactionSubtitle(tx) {
  const quantity = `${tx.quantity} ${tx.item?.unit || ""}`.trim();
  const warehouses = [
    tx.from_warehouse?.name ? `From: ${tx.from_warehouse.name}` : "",
    tx.to_warehouse?.name ? `To: ${tx.to_warehouse.name}` : "",
  ]
    .filter(Boolean)
    .join(" • ");

  return [quantity, warehouses].filter(Boolean).join(" | ");
}

function transactionActor(tx) {
  if (tx.created_by?.full_name) {
    return tx.created_by.full_name;
  }
  if (tx.created_by?.email) {
    return tx.created_by.email;
  }
  return "";
}

function getAvatarUrl(avatarPath, size = 32) {
  if (!avatarPath) {
    return `https://via.placeholder.com/${size}`;
  }
  if (avatarPath.startsWith("data:")) {
    return avatarPath;
  }
  if (avatarPath.startsWith("http")) {
    return avatarPath;
  }
  const base = fileBaseUrl || apiBaseUrl;
  const normalizedBase = base.endsWith("/") ? base.slice(0, -1) : base;
  const normalizedPath = avatarPath.startsWith("/")
    ? avatarPath.slice(1)
    : avatarPath;
  return `${normalizedBase}/${normalizedPath}`;
}

function transactionActorAvatar(tx, size = 32) {
  const avatarPath = tx.created_by?.avatar || "";
  return getAvatarUrl(avatarPath, size);
}

function normalizeStatus(value) {
  return (value || "").toString().trim().toLowerCase();
}

function parseDate(value) {
  if (!value) return null;
  if (value instanceof Date) {
    return Number.isNaN(value.getTime()) ? null : value;
  }
  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? null : date;
}

function calculateAge(birthDate) {
  const date = parseDate(birthDate);
  if (!date) return null;
  const today = new Date();
  let age = today.getFullYear() - date.getFullYear();
  const monthDiff = today.getMonth() - date.getMonth();
  if (
    monthDiff < 0 ||
    (monthDiff === 0 && today.getDate() < date.getDate())
  ) {
    age -= 1;
  }
  return age;
}

function normalizeGender(value) {
  const normalized = (value || "").toString().trim().toLowerCase();
  if (!normalized) return "unspecified";
  if (
    ["pria", "laki-laki", "laki laki", "male", "m"].some(
      (option) => option === normalized,
    )
  ) {
    return "male";
  }
  if (
    ["wanita", "perempuan", "female", "f"].some(
      (option) => option === normalized,
    )
  ) {
    return "female";
  }
  return "other";
}

function formatDateFriendly(value, options = { dateStyle: "full" }) {
  const date = parseDate(value);
  if (!date) return "-";
  return new Intl.DateTimeFormat("id-ID", options).format(date);
}

function birthdayKey(value) {
  const date = parseDate(value);
  if (!date) return null;
  return `${date.getMonth()}-${date.getDate()}`;
}

function birthdayTooltip(value) {
  const key = birthdayKey(value);
  if (!key) return "";
  const entries = birthdayLookup.value.get(key) || [];
  if (entries.length === 0) return "";
  return `(${entries.length}) karyawan ulang tahun`;
}

function computeBarWidth(count, max) {
  if (!max) return "0%";
  if (count === 0) return "0%";
  const raw = Math.round((count / max) * 100);
  const adjusted = Math.min(Math.max(raw, 8), 100);
  return `${adjusted}%`;
}

function normalizeEmployeePhoto(photoPath) {
  if (!photoPath || typeof photoPath !== "string") return "";
  if (photoPath.startsWith("data:")) return photoPath;
  if (photoPath.startsWith("http")) return photoPath;
  return getAvatarUrl(photoPath, 96);
}

function initials(value) {
  if (!value) return "??";
  return value
    .split(/\s+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0].toUpperCase())
    .join("");
}

function getISOWeekNumber(date) {
  const tempDate = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
  const dayNumber = tempDate.getUTCDay() || 7;
  tempDate.setUTCDate(tempDate.getUTCDate() + 4 - dayNumber);
  const yearStart = new Date(Date.UTC(tempDate.getUTCFullYear(), 0, 1));
  return Math.ceil(((tempDate - yearStart) / 86400000 + 1) / 7);
}

function isSameDate(a, b) {
  const dateA = parseDate(a);
  const dateB = parseDate(b);
  if (!dateA || !dateB) return false;
  return (
    dateA.getFullYear() === dateB.getFullYear() &&
    dateA.getMonth() === dateB.getMonth() &&
    dateA.getDate() === dateB.getDate()
  );
}

function startOfMonth(value) {
  const date = parseDate(value);
  if (!date) {
    const today = new Date();
    return new Date(today.getFullYear(), today.getMonth(), 1);
  }
  return new Date(date.getFullYear(), date.getMonth(), 1);
}

function calendarViewYearMonth(date) {
  calendarViewDate.value = startOfMonth(date);
}

function goToPrevMonth() {
  const prev = new Date(calendarViewDate.value);
  prev.setMonth(prev.getMonth() - 1);
  const normalized = startOfMonth(prev);
  calendarViewDate.value = normalized;
  calendarValue.value = normalized;
}

function goToNextMonth() {
  const next = new Date(calendarViewDate.value);
  next.setMonth(next.getMonth() + 1);
  const normalized = startOfMonth(next);
  calendarViewDate.value = normalized;
  calendarValue.value = normalized;
}

function selectDay(day) {
  calendarValue.value = new Date(
    day.date.getFullYear(),
    day.date.getMonth(),
    day.date.getDate(),
  );
}

function handleDayClick(day) {
  selectDay(day);
  if (day.birthdays.length) {
    birthdayModal.value = {
      date: new Date(
        day.date.getFullYear(),
        day.date.getMonth(),
        day.date.getDate(),
      ),
      employees: day.birthdays.map((entry) => ({ ...entry })),
    };
    showBirthdayModal.value = true;
  }
}

function closeBirthdayModal() {
  showBirthdayModal.value = false;
  birthdayModal.value = {
    date: null,
    employees: [],
  };
}
</script>
