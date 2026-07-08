<template>
  <Modal :show="show" title="Custom Query Runner" @close="$emit('close')">
    <div class="space-y-5 py-2">
      <!-- Description -->
      <p class="text-xs text-text-secondary leading-relaxed">
        Pilih bulan untuk mengeksekusi query. Variabel <code class="text-teal-accent">:date_from</code> akan digantikan oleh <strong>Tanggal 1</strong> dan <code class="text-teal-accent">:date_to</code> akan digantikan oleh <strong>Tanggal 2</strong> (hari terakhir di bulan tersebut) tanpa mengubah teks query asli.
      </p>

      <!-- Custom Dropdown Selector -->
      <div class="relative dropdown-container">
        <label class="block text-xs text-text-secondary mb-1.5">Pilih Periode Bulan</label>
        
        <!-- Dropdown trigger button -->
        <button
          type="button"
          @click="dropdownOpen = !dropdownOpen"
          class="w-full flex items-center justify-between bg-navy-tertiary border border-navy-border hover:border-teal-accent focus:border-teal-accent rounded-lg px-3 py-2 text-xs text-text-primary transition-all duration-150 cursor-pointer"
        >
          <div class="flex items-center gap-2">
            <!-- Calendar Icon -->
            <svg class="w-4 h-4 text-text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
              <line x1="16" y1="2" x2="16" y2="6" />
              <line x1="8" y1="2" x2="8" y2="6" />
              <line x1="3" y1="10" x2="21" y2="10" />
            </svg>
            <span class="font-medium">{{ selectedOption ? selectedOption.label : 'Pilih Bulan' }}</span>
          </div>
          <!-- Chevron Icon -->
          <svg
            class="w-3.5 h-3.5 text-text-muted transition-transform duration-150"
            :class="dropdownOpen ? 'rotate-180' : ''"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
          >
            <path d="m6 9 6 6 6-6" />
          </svg>
        </button>

        <!-- Dropdown Menu Overlay -->
        <div
          v-if="dropdownOpen"
          class="absolute left-0 right-0 mt-1.5 bg-navy-secondary border border-navy-border rounded-lg shadow-xl z-50 py-1 max-h-60 overflow-y-auto"
        >
          <button
            v-for="opt in monthOptions"
            :key="opt.label"
            type="button"
            @click="selectOption(opt)"
            class="w-full flex items-center justify-between px-3 py-2 text-xs text-left transition-colors duration-150 hover:bg-navy-hover cursor-pointer"
            :class="selectedOption?.label === opt.label ? 'bg-navy-hover/50 text-text-primary' : 'text-text-secondary hover:text-text-primary'"
          >
            <div class="flex items-center gap-2">
              <!-- Dot Indicator -->
              <span
                class="w-1.5 h-1.5 rounded-full flex-shrink-0"
                :class="selectedOption?.label === opt.label ? 'bg-teal-accent' : 'bg-text-muted/40'"
              ></span>
              <span class="font-medium">{{ opt.label }}</span>
              
              <!-- Ongoing Tag -->
              <span
                v-if="opt.isOngoing"
                class="text-[9px] font-semibold text-amber-500 bg-amber-500/10 px-1.5 py-0.5 rounded ml-1"
              >
                ongoing
              </span>
            </div>

            <!-- Checkmark Icon -->
            <svg
              v-if="selectedOption?.label === opt.label"
              class="w-3.5 h-3.5 text-teal-accent"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
            >
              <polyline points="20 6 9 17 4 12" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Preview of dates -->
      <div v-if="selectedOption" class="bg-navy-tertiary/40 border border-navy-border/60 rounded-lg p-3 space-y-2">
        <div class="text-[11px] text-text-secondary font-medium">Pratinjau Penggantian Parameter:</div>
        <div class="grid grid-cols-2 gap-4 text-xs">
          <div>
            <span class="text-text-muted mr-1.5">Tanggal 1 (:date_from) ➜</span>
            <code class="text-teal-accent font-semibold bg-navy-secondary/80 px-1.5 py-0.5 rounded border border-navy-border">'{{ selectedOption.dateFrom }}'::date</code>
          </div>
          <div>
            <span class="text-text-muted mr-1.5">Tanggal 2 (:date_to) ➜</span>
            <code class="text-teal-accent font-semibold bg-navy-secondary/80 px-1.5 py-0.5 rounded border border-navy-border">'{{ selectedOption.dateTo }}'::date</code>
          </div>
        </div>
      </div>

      <!-- Actions -->
      <div class="flex items-center justify-end gap-2 pt-2 border-t border-navy-border/40">
        <button
          type="button"
          @click="$emit('close')"
          class="px-3 py-1.5 text-xs text-text-secondary hover:bg-navy-hover rounded-md transition-colors"
        >
          Batal
        </button>
        <button
          type="button"
          @click="handleExecute"
          class="flex items-center gap-1.5 px-4 py-1.5 text-xs bg-teal-accent text-navy-primary rounded-md font-medium hover:bg-teal-hover transition-colors shadow-lg shadow-teal-accent/10"
        >
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="currentColor">
            <polygon points="5 3 19 12 5 21 5 3" />
          </svg>
          Jalankan
        </button>
      </div>
    </div>
  </Modal>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import Modal from '../shared/Modal.vue'

defineProps<{
  show: boolean
}>()

const emit = defineEmits(['close', 'execute'])

interface MonthOption {
  label: string
  dateFrom: string
  dateTo: string
  isOngoing: boolean
}

// Generate the last 12 months dynamically
const monthOptions = computed<MonthOption[]>(() => {
  const options: MonthOption[] = []
  const now = new Date()
  const currentYear = now.getFullYear()
  const currentMonthIdx = now.getMonth()

  const monthNames = [
    'January', 'February', 'March', 'April', 'May', 'June',
    'July', 'August', 'September', 'October', 'November', 'December'
  ]

  for (let i = 0; i < 12; i++) {
    // Generate dates going backwards from current month
    const d = new Date(currentYear, currentMonthIdx - i, 1)
    const year = d.getFullYear()
    const month = d.getMonth()

    const label = `${monthNames[month]} ${year}`

    const yStr = year.toString()
    const mStr = (month + 1).toString().padStart(2, '0')
    const dateFrom = `${yStr}${mStr}01`

    const lastDay = new Date(year, month + 1, 0).getDate()
    const dateTo = `${yStr}${mStr}${lastDay.toString().padStart(2, '0')}`

    options.push({
      label,
      dateFrom,
      dateTo,
      isOngoing: i === 0
    })
  }
  return options
})

const selectedOption = ref<MonthOption | null>(null)
const dropdownOpen = ref(false)

const selectOption = (opt: MonthOption) => {
  selectedOption.value = opt
  dropdownOpen.value = false
}

const handleExecute = () => {
  if (selectedOption.value) {
    emit('execute', {
      dateFrom: selectedOption.value.dateFrom,
      dateTo: selectedOption.value.dateTo
    })
  }
}

const closeDropdown = (e: MouseEvent) => {
  const target = e.target as HTMLElement
  if (!target.closest('.dropdown-container')) {
    dropdownOpen.value = false
  }
}

onMounted(() => {
  if (monthOptions.value.length > 0) {
    selectedOption.value = monthOptions.value[0]
  }
  document.addEventListener('click', closeDropdown)
})

onUnmounted(() => {
  document.removeEventListener('click', closeDropdown)
})
</script>
