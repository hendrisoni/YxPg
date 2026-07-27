<template>
  <div class="h-full flex flex-col bg-navy-primary text-text-primary overflow-hidden">
    <!-- Header / Toolbar -->
    <div class="p-3 bg-navy-secondary border-b border-navy-border flex flex-wrap items-center justify-between gap-3 flex-shrink-0">
      <!-- Left side: Title & Connection/Schema Selectors -->
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-2">
          <div class="p-1.5 rounded bg-teal-accent/10 text-teal-accent">
            <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <line x1="18" y1="20" x2="18" y2="10" />
              <line x1="12" y1="20" x2="12" y2="4" />
              <line x1="6" y1="20" x2="6" y2="14" />
            </svg>
          </div>
          <div>
            <h1 class="text-sm font-semibold text-text-primary leading-none">Table & Index Size</h1>
            <p class="text-[11px] text-text-muted mt-0.5">Analisis ukuran penyimpanan tabel dan indeks (Urutan terbesar)</p>
          </div>
        </div>

        <!-- Connection Select -->
        <div class="flex items-center gap-1.5 ml-2">
          <label class="text-xs text-text-secondary">Conn:</label>
          <select
            v-model="selectedConnId"
            @change="loadData"
            class="text-xs bg-navy-tertiary border border-navy-border rounded px-2 py-1 text-text-primary focus:border-teal-accent focus:outline-none"
          >
            <option v-for="conn in connectionsStore.connections" :key="conn.id" :value="conn.id">
              {{ conn.name }} ({{ conn.database }})
            </option>
          </select>
        </div>

        <!-- Schema Select -->
        <div class="flex items-center gap-1.5">
          <label class="text-xs text-text-secondary">Schema:</label>
          <select
            v-model="selectedSchema"
            @change="loadData"
            class="text-xs bg-navy-tertiary border border-navy-border rounded px-2 py-1 text-text-primary focus:border-teal-accent focus:outline-none"
          >
            <option value="">All Schemas</option>
            <option v-for="s in availableSchemas" :key="s" :value="s">
              {{ s }}
            </option>
          </select>
        </div>
      </div>

      <!-- Right side: Color legend, Search input & Refresh button -->
      <div class="flex items-center gap-2">
        <!-- Color Legend Badge (Warning Level based on Size) -->
        <div class="hidden sm:flex items-center gap-1.5 text-[10px] bg-navy-tertiary/70 border border-navy-border/70 rounded px-2 py-1 font-mono">
          <span class="text-text-muted font-sans font-medium mr-0.5">Level:</span>
          <span class="text-emerald-400 font-medium" title="Dibawah 100 MB">&lt;100MB</span>
          <span class="text-text-muted">•</span>
          <span class="text-amber-400 font-semibold" title="100 MB - 500 MB">100-500MB</span>
          <span class="text-text-muted">•</span>
          <span class="text-orange-400 font-bold" title="500 MB - 1 GB">500MB-1GB</span>
          <span class="text-text-muted">•</span>
          <span class="text-red-400 font-bold" title="Diatas 1 GB">&gt;1GB</span>
        </div>


        <div class="relative">

          <svg class="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-text-muted" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="11" cy="11" r="8" />
            <path d="M21 21l-4.35-4.35" />
          </svg>
          <input
            v-model="searchFilter"
            type="text"
            placeholder="Cari tabel / indeks..."
            class="pl-8 pr-3 py-1 text-xs bg-navy-tertiary border border-navy-border rounded text-text-primary placeholder-text-muted focus:border-teal-accent focus:outline-none w-48"
          />
        </div>

        <button
          @click="loadData"
          :disabled="loading"
          class="flex items-center gap-1 px-2.5 py-1 text-xs bg-teal-accent text-navy-primary rounded font-medium hover:bg-teal-hover transition-colors disabled:opacity-50 cursor-pointer"
          title="Refresh Data"
        >
          <svg class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21.5 2v6h-6M2.5 22v-6h6" />
            <path d="M2 11.5a10 10 0 0 1 18.8-4.3L21.5 8M22 12.5a10 10 0 0 1-18.8 4.3L2.5 16" />
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 overflow-auto p-3">
      <!-- Loading State -->
      <div v-if="loading" class="h-full flex items-center justify-center">
        <div class="text-center">
          <div class="inline-block w-8 h-8 border-3 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
          <p class="text-xs text-text-muted mt-2">Memuat data ukuran tabel & indeks...</p>
        </div>
      </div>

      <!-- Error State -->
      <div v-else-if="errorMessage" class="h-full flex items-center justify-center p-6 text-center">
        <div class="bg-red-500/10 border border-red-500/30 rounded-lg p-6 max-w-md">
          <svg class="w-10 h-10 text-red-400 mx-auto mb-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" />
            <line x1="12" y1="8" x2="12" y2="12" />
            <line x1="12" y1="16" x2="12.01" y2="16" />
          </svg>
          <h3 class="text-sm font-semibold text-red-400">Gagal Memuat Ukuran</h3>
          <p class="text-xs text-text-muted mt-1">{{ errorMessage }}</p>
          <button @click="loadData" class="mt-3 px-3 py-1 text-xs bg-navy-tertiary hover:bg-navy-hover text-text-primary rounded border border-navy-border">
            Coba Lagi
          </button>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredItems.length === 0" class="h-full flex items-center justify-center text-center p-6">
        <div>
          <svg class="w-10 h-10 text-text-muted mx-auto mb-2" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M20 7H4a2 2 0 0 0-2 2v10a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2Z" />
            <path d="M16 21V5a2 2 0 0 0-2-2h-4a2 2 0 0 0-2 2v16" />
          </svg>
          <p class="text-xs text-text-muted">Tidak ada data ukuran yang ditemukan.</p>
        </div>
      </div>

      <!-- Data Table -->
      <div v-else class="border border-navy-border rounded-lg overflow-hidden bg-navy-secondary shadow-md">
        <table class="w-full text-left text-xs border-collapse">
          <thead>
            <tr class="bg-navy-tertiary/80 text-text-secondary border-b border-navy-border select-none">
              <th class="py-2.5 px-3 w-10 text-center">#</th>
              <th class="py-2.5 px-3 w-10"></th>
              <th @click="toggleSort('schema')" class="py-2.5 px-3 cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center gap-1">
                  Schema
                  <span v-if="sortBy === 'schema'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
                </div>
              </th>
              <th @click="toggleSort('table_name')" class="py-2.5 px-3 cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center gap-1">
                  Nama Tabel
                  <span v-if="sortBy === 'table_name'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
                </div>
              </th>
              <th @click="toggleSort('row_count')" class="py-2.5 px-3 text-right cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center justify-end gap-1">
                  Baris
                  <span v-if="sortBy === 'row_count'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
                </div>
              </th>
              <th @click="toggleSort('table_bytes')" class="py-2.5 px-3 text-right cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center justify-end gap-1">
                  Ukuran Tabel (KB/MB/GB)
                  <span v-if="sortBy === 'table_bytes'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
                </div>
              </th>
              <th @click="toggleSort('index_bytes')" class="py-2.5 px-3 text-right cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center justify-end gap-1">
                  Ukuran Indeks (KB/MB/GB)
                  <span v-if="sortBy === 'index_bytes'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
                </div>
              </th>
              <th @click="toggleSort('total_bytes')" class="py-2.5 px-3 text-right cursor-pointer hover:text-teal-accent transition-colors">
                <div class="flex items-center justify-end gap-1">
                  Total Ukuran
                  <span v-if="sortBy === 'total_bytes'">{{ sortDir === 'asc' ? '▲' : '▼' }}</span>
                </div>
              </th>
              <th class="py-2.5 px-3 w-40 text-center">Visual Ukuran</th>
              <th class="py-2.5 px-3 w-28 text-center">Aksi</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-navy-border/60 font-mono text-[11px]">
            <template v-for="(item, index) in filteredItems" :key="item.schema + '.' + item.table_name">
              <tr class="hover:bg-navy-hover/50 transition-colors" :class="{ 'bg-teal-accent/5': expandedRows.has(item.schema + '.' + item.table_name) }">
                <td class="py-2 px-3 text-center text-text-muted">{{ index + 1 }}</td>
                <td class="py-2 px-3 text-center">
                  <button
                    v-if="item.indexes && item.indexes.length > 0"
                    @click="toggleExpandRow(item.schema + '.' + item.table_name)"
                    class="p-0.5 rounded text-text-muted hover:text-teal-accent hover:bg-navy-tertiary transition-colors cursor-pointer"
                    :title="expandedRows.has(item.schema + '.' + item.table_name) ? 'Sembunyikan indeks' : 'Lihat rincian indeks'"
                  >
                    <svg class="w-3.5 h-3.5 transition-transform duration-150" :class="{ 'rotate-90': expandedRows.has(item.schema + '.' + item.table_name) }" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="9 18 15 12 9 6" />
                    </svg>
                  </button>
                </td>
                <td class="py-2 px-3 text-teal-accent font-sans font-medium">{{ item.schema }}</td>
                <td class="py-2 px-3 text-text-primary font-sans font-semibold">
                  {{ item.table_name }}
                  <span v-if="item.indexes && item.indexes.length > 0" class="ml-1 text-[10px] text-text-muted font-normal font-sans">
                    ({{ item.indexes.length }} idx)
                  </span>
                </td>
                <td class="py-2 px-3 text-right text-text-secondary">{{ formatNumber(item.row_count) }}</td>
                <td class="py-2 px-3 text-right text-text-primary font-semibold">{{ formatBytes(item.table_bytes) }}</td>
                <td class="py-2 px-3 text-right text-amber-400/90">{{ formatBytes(item.index_bytes) }}</td>
                <td class="py-2 px-3 text-right font-bold" :class="getSizeColorClass(item.total_bytes)">{{ formatBytes(item.total_bytes) }}</td>

                <td class="py-2 px-3">
                  <div class="w-full bg-navy-tertiary rounded-full h-2 overflow-hidden flex" title="Tabel (Teal) vs Indeks (Kuning)">
                    <div
                      class="bg-teal-accent h-full transition-all duration-300"
                      :style="{ width: getPercentage(item.table_bytes, maxTotalBytes) + '%' }"
                    ></div>
                    <div
                      class="bg-amber-400 h-full transition-all duration-300"
                      :style="{ width: getPercentage(item.index_bytes, maxTotalBytes) + '%' }"
                    ></div>
                  </div>
                </td>
                <td class="py-2 px-3 text-center">
                  <div class="flex items-center justify-center gap-1.5 font-sans">
                    <button
                      @click="openTable(item.schema, item.table_name)"
                      class="px-2 py-0.5 text-[10px] bg-navy-tertiary hover:bg-teal-accent hover:text-navy-primary text-text-secondary rounded border border-navy-border transition-colors cursor-pointer"
                      title="Buka Data Tabel"
                    >
                      Buka
                    </button>
                    <button
                      @click="openQuery(item.schema, item.table_name)"
                      class="px-2 py-0.5 text-[10px] bg-navy-tertiary hover:bg-teal-accent hover:text-navy-primary text-text-secondary rounded border border-navy-border transition-colors cursor-pointer"
                      title="Query Tabel ini"
                    >
                      Query
                    </button>
                  </div>
                </td>
              </tr>

              <!-- Expanded Index Breakdown Row -->
              <tr v-if="expandedRows.has(item.schema + '.' + item.table_name)" class="bg-navy-tertiary/40">
                <td colspan="10" class="p-3 pl-12 border-b border-navy-border/60">
                  <div class="bg-navy-primary/80 border border-navy-border rounded p-3 text-xs">
                    <h4 class="text-xs font-semibold text-amber-400 mb-2 flex items-center gap-1 font-sans">
                      <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                        <line x1="18" y1="20" x2="18" y2="10" />
                        <line x1="12" y1="20" x2="12" y2="4" />
                        <line x1="6" y1="20" x2="6" y2="14" />
                      </svg>
                      Rincian Indeks Tabel: {{ item.schema }}.{{ item.table_name }}
                    </h4>

                    <div v-if="!item.indexes || item.indexes.length === 0" class="text-text-muted text-[11px] font-sans">
                      Tidak ada indeks terpisah pada tabel ini.
                    </div>
                    <table v-else class="w-full text-left text-[11px] border-collapse font-mono">
                      <thead>
                        <tr class="text-text-muted border-b border-navy-border/80 font-sans">
                          <th class="py-1 px-2">Nama Indeks</th>
                          <th class="py-1 px-2">Tipe Indeks</th>
                          <th class="py-1 px-2 text-right">Ukuran (KB/MB/GB)</th>
                          <th class="py-1 px-2 text-right">% dari Total Indeks</th>
                        </tr>
                      </thead>
                      <tbody class="divide-y divide-navy-border/40">
                        <tr v-for="idx in item.indexes" :key="idx.index_name" class="hover:bg-navy-hover/30">
                          <td class="py-1.5 px-2 text-text-primary font-semibold">{{ idx.index_name }}</td>
                          <td class="py-1.5 px-2 text-text-secondary uppercase text-[10px] font-sans">{{ idx.index_type || 'btree' }}</td>
                          <td class="py-1.5 px-2 text-right text-amber-400 font-bold">{{ formatBytes(idx.index_bytes) }}</td>
                          <td class="py-1.5 px-2 text-right text-text-muted">
                            {{ item.index_bytes > 0 ? ((idx.index_bytes / item.index_bytes) * 100).toFixed(1) + '%' : '0%' }}
                          </td>
                        </tr>
                      </tbody>
                    </table>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Summary Bar Footer -->
    <div class="px-4 py-2.5 bg-navy-secondary border-t border-navy-border flex flex-wrap items-center justify-between text-xs text-text-secondary flex-shrink-0">
      <div class="flex items-center gap-4">
        <span>Total Tabel: <strong class="text-text-primary">{{ filteredItems.length }}</strong></span>
        <span>Total Baris: <strong class="text-text-primary">{{ formatNumber(summaryStats.totalRows) }}</strong></span>
      </div>
      <div class="flex items-center gap-5">
        <span class="flex items-center gap-1">
          <span class="w-2.5 h-2.5 rounded-full bg-teal-accent inline-block"></span>
          Ukuran Tabel: <strong class="text-text-primary">{{ formatBytes(summaryStats.totalTableBytes) }}</strong>
        </span>
        <span class="flex items-center gap-1">
          <span class="w-2.5 h-2.5 rounded-full bg-amber-400 inline-block"></span>
          Ukuran Indeks: <strong class="text-text-primary">{{ formatBytes(summaryStats.totalIndexBytes) }}</strong>
        </span>
        <span class="flex items-center gap-1 text-teal-accent font-semibold">
          Total DB Size: <strong class="text-teal-accent font-bold text-sm">{{ formatBytes(summaryStats.totalBytes) }}</strong>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useConnectionsStore } from '../stores/connections'
import { useSchemaStore } from '../stores/schema'
import { useTabsStore } from '../stores/tabs'
import { useUiStore } from '../stores/ui'
import type { Tab, TableIndexSizeInfo } from '../types'
import * as App from '../../wailsjs/go/main/App'

const props = defineProps<{
  tab: Tab
}>()

const connectionsStore = useConnectionsStore()
const schemaStore = useSchemaStore()
const tabsStore = useTabsStore()
const uiStore = useUiStore()

const selectedConnId = ref<string>(props.tab.connectionId || connectionsStore.currentConnectionId || '')
const selectedSchema = ref<string>('')
const searchFilter = ref<string>('')
const items = ref<TableIndexSizeInfo[]>([])
const loading = ref<boolean>(false)
const errorMessage = ref<string>('')

// Sorting state: default total_bytes descending (largest size first)
const sortBy = ref<keyof TableIndexSizeInfo>('total_bytes')
const sortDir = ref<'asc' | 'desc'>('desc')

const expandedRows = ref<Set<string>>(new Set())

const availableSchemas = computed(() => {
  const set = new Set<string>()
  for (const item of items.value) {
    if (item.schema) set.add(item.schema)
  }
  return Array.from(set).sort()
})

const maxTotalBytes = computed(() => {
  let max = 0
  for (const item of items.value) {
    if (item.total_bytes > max) max = item.total_bytes
  }
  return max || 1
})

const filteredItems = computed(() => {
  let res = items.value

  if (selectedSchema.value) {
    res = res.filter(i => i.schema === selectedSchema.value)
  }

  if (searchFilter.value.trim()) {
    const q = searchFilter.value.toLowerCase().trim()
    res = res.filter(i => {
      const matchTable = i.table_name.toLowerCase().includes(q) || i.schema.toLowerCase().includes(q)
      const matchIndex = i.indexes?.some(idx => idx.index_name.toLowerCase().includes(q))
      return matchTable || matchIndex
    })
  }

  // Sort results
  return res.slice().sort((a, b) => {
    let valA = a[sortBy.value]
    let valB = b[sortBy.value]

    if (typeof valA === 'string') {
      valA = (valA as string).toLowerCase()
      valB = (valB as string).toLowerCase()
    }

    if (valA < valB) return sortDir.value === 'asc' ? -1 : 1
    if (valA > valB) return sortDir.value === 'asc' ? 1 : -1
    return 0
  })
})

const summaryStats = computed(() => {
  let totalRows = 0
  let totalTableBytes = 0
  let totalIndexBytes = 0
  let totalBytes = 0

  for (const item of filteredItems.value) {
    totalRows += item.row_count || 0
    totalTableBytes += item.table_bytes || 0
    totalIndexBytes += item.index_bytes || 0
    totalBytes += item.total_bytes || 0
  }

  return {
    totalRows,
    totalTableBytes,
    totalIndexBytes,
    totalBytes
  }
})

onMounted(() => {
  if (!selectedConnId.value && connectionsStore.connections.length > 0) {
    selectedConnId.value = connectionsStore.connections[0].id
  }
  loadData()
})

watch(() => props.tab.connectionId, (newConnId) => {
  if (newConnId && newConnId !== selectedConnId.value) {
    selectedConnId.value = newConnId
    loadData()
  }
})

async function loadData() {
  if (!selectedConnId.value) {
    errorMessage.value = 'Silakan pilih koneksi database terlebih dahulu.'
    return
  }

  loading.value = true
  errorMessage.value = ''

  try {
    const res = await App.GetTableIndexSizes(selectedConnId.value, selectedSchema.value)
    items.value = res || []
  } catch (err: any) {
    console.error('Failed to get table index sizes:', err)
    errorMessage.value = err?.message || String(err)
    uiStore.addNotification({
      type: 'error',
      title: 'Gagal Memuat Ukuran',
      message: errorMessage.value
    })
  } finally {
    loading.value = false
  }
}

function toggleSort(field: keyof TableIndexSizeInfo) {
  if (sortBy.value === field) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = field
    sortDir.value = field === 'table_name' || field === 'schema' ? 'asc' : 'desc'
  }
}

function toggleExpandRow(key: string) {
  if (expandedRows.value.has(key)) {
    expandedRows.value.delete(key)
  } else {
    expandedRows.value.add(key)
  }
}

function formatBytes(bytes: number): string {
  if (bytes === 0 || !bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  const num = bytes / Math.pow(k, i)
  return `${num.toFixed(num >= 100 || i === 0 ? 0 : 2)} ${sizes[i]}`
}

function getSizeColorClass(bytes: number): string {
  if (!bytes || bytes <= 0) return 'text-text-muted font-normal'
  const mb = 1024 * 1024
  const gb = 1024 * 1024 * 1024

  if (bytes >= 5 * gb) return 'text-red-500 font-extrabold'     // >= 5 GB (Merah Pekat / Kritis)
  if (bytes >= 1 * gb) return 'text-red-400 font-bold'          // >= 1 GB (Merah / Warning Tinggi)
  if (bytes >= 500 * mb) return 'text-orange-400 font-bold'     // 500 MB - 1 GB (Oranye / Warning Sedang)
  if (bytes >= 100 * mb) return 'text-amber-400 font-semibold'  // 100 MB - 500 MB (Kuning / Warning Awal)
  return 'text-emerald-400 font-medium'                         // < 100 MB (Hijau / Aman)
}



function formatNumber(num: number): string {
  return new Intl.NumberFormat().format(num || 0)
}

function getPercentage(part: number, total: number): number {
  if (!total || total <= 0) return 0
  return Math.min(100, Math.max(0, (part / total) * 100))
}

function openTable(schema: string, table: string) {
  if (!selectedConnId.value) return
  tabsStore.createTab('table', {
    title: table,
    connectionId: selectedConnId.value,
    schema,
    table
  })
}

function openQuery(schema: string, table: string) {
  if (!selectedConnId.value) return
  tabsStore.createTab('query', {
    title: table,
    connectionId: selectedConnId.value,
    schema,
    table,
    sql: `SELECT * FROM ${schema}.${table} LIMIT 100`
  })
}
</script>
