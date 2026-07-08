<template>
  <div class="h-full flex flex-col overflow-hidden">
    <!-- Header -->
    <div class="flex items-center justify-between px-3 py-2 border-b border-navy-border">
      <h3 class="text-xs font-semibold text-text-primary">Query History</h3>
      <button
        @click="$emit('close')"
        class="p-1 rounded hover:bg-navy-hover text-text-muted hover:text-text-primary"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M18 6 6 18M6 6l12 12" />
        </svg>
      </button>
    </div>

    <!-- Search -->
    <div class="px-3 py-2 border-b border-navy-border">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search queries..."
        class="w-full px-2 py-1 text-xs bg-navy-tertiary border border-navy-border rounded"
      />
    </div>

    <!-- History list -->
    <div class="flex-1 overflow-y-auto">
      <div v-if="loading" class="flex items-center justify-center py-8">
        <div class="w-5 h-5 border-2 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
      </div>

      <div v-else-if="filteredHistory.length === 0" class="text-center py-8 text-xs text-text-muted">
        No query history
      </div>

      <div v-else>
        <div
          v-for="entry in filteredHistory"
          :key="entry.id"
          @click="$emit('select', entry.sql)"
          class="px-3 py-2 border-b border-navy-border cursor-pointer hover:bg-navy-hover transition-colors group"
        >
          <div class="flex items-center gap-2 mb-1">
            <span
              class="text-[10px] px-1.5 py-0.5 rounded font-medium"
              :class="getQueryTypeClass(entry.sql)"
            >
              {{ getQueryType(entry.sql) }}
            </span>
            <span class="text-[10px] text-text-muted">{{ formatDuration(entry.duration_ms) }}</span>
            <span v-if="entry.rows_returned > 0" class="text-[10px] text-text-muted">{{ entry.rows_returned }} rows</span>
          </div>
          <p class="text-xs text-text-secondary font-mono truncate">{{ truncateSQL(entry.sql) }}</p>
          <div class="flex items-center gap-2 mt-1">
            <span class="text-[10px] text-text-muted">{{ formatDate(entry.executed_at) }}</span>
            <span v-if="entry.error" class="text-[10px] text-accent-red">Error</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useConnectionsStore } from '../../stores/connections'
import type { QueryHistoryEntry } from '../../types'

const emit = defineEmits(['close', 'select'])

const connectionsStore = useConnectionsStore()

const history = ref<QueryHistoryEntry[]>([])
const loading = ref(false)
const searchQuery = ref('')

const filteredHistory = computed(() => {
  if (!searchQuery.value) return history.value
  const q = searchQuery.value.toLowerCase()
  return history.value.filter(e => e.sql.toLowerCase().includes(q))
})

onMounted(() => {
  loadHistory()
})

async function loadHistory() {
  if (!connectionsStore.currentConnectionId) return
  loading.value = true
  try {
    const bindings = connectionsStore.getWailsBindings()
    history.value = await bindings.GetQueryHistory(connectionsStore.currentConnectionId, 200)
  } catch {
    // Ignore
  } finally {
    loading.value = false
  }
}

function getQueryType(sql: string): string {
  const trimmed = sql.trim().toUpperCase()
  if (trimmed.startsWith('SELECT') || trimmed.startsWith('WITH')) return 'SELECT'
  if (trimmed.startsWith('INSERT')) return 'INSERT'
  if (trimmed.startsWith('UPDATE')) return 'UPDATE'
  if (trimmed.startsWith('DELETE')) return 'DELETE'
  if (trimmed.startsWith('CREATE') || trimmed.startsWith('ALTER') || trimmed.startsWith('DROP')) return 'DDL'
  return 'OTHER'
}

function getQueryTypeClass(sql: string): string {
  const type = getQueryType(sql)
  const classes: Record<string, string> = {
    SELECT: 'bg-accent-blue/20 text-accent-blue',
    INSERT: 'bg-accent-green/20 text-accent-green',
    UPDATE: 'bg-accent-amber/20 text-accent-amber',
    DELETE: 'bg-accent-red/20 text-accent-red',
    DDL: 'bg-purple-400/20 text-purple-400',
    OTHER: 'bg-text-muted/20 text-text-muted',
  }
  return classes[type] || classes.OTHER
}

function truncateSQL(sql: string): string {
  const oneLine = sql.replace(/\s+/g, ' ').trim()
  return oneLine.length > 120 ? oneLine.substring(0, 120) + '...' : oneLine
}

function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms}ms`
  return `${(ms / 1000).toFixed(1)}s`
}

function formatDate(dateStr: string): string {
  try {
    const date = new Date(dateStr)
    return date.toLocaleString()
  } catch {
    return dateStr
  }
}
</script>
