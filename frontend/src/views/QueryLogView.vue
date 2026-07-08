<template>
  <div class="h-full w-full flex flex-col bg-navy-primary overflow-hidden p-6">
    <div class="flex items-center justify-between mb-4 flex-shrink-0">
      <div>
        <h2 class="text-lg font-semibold text-text-primary flex items-center gap-2">
          <svg class="w-5 h-5 text-teal-accent" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M12 8v4l3 3" />
            <circle cx="12" cy="12" r="10" />
          </svg>
          Query History Log
        </h2>
        <p class="text-xs text-text-secondary mt-1">
          Menampilkan hingga 7 query terakhir yang pernah Anda jalankan. Anda dapat memuatnya kembali atau menyalin skripnya.
        </p>
      </div>
      <button
        v-if="logs.length > 0"
        @click="clearAllLogs"
        class="flex items-center gap-1.5 px-3 py-1.5 text-xs bg-red-950/20 hover:bg-red-950/40 border border-red-900/30 text-accent-red rounded-md transition-colors cursor-pointer"
      >
        Clear Log
      </button>
    </div>

    <!-- Empty State -->
    <div v-if="logs.length === 0" class="flex-1 flex flex-col items-center justify-center text-center p-8 bg-navy-secondary/20 border border-navy-border/40 rounded-xl">
      <svg class="w-12 h-12 text-text-muted mb-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M12 8v4l3 3" />
        <circle cx="12" cy="12" r="10" />
      </svg>
      <h3 class="text-sm font-semibold text-text-muted">Belum Ada Riwayat Query</h3>
      <p class="text-xs text-text-muted mt-1 max-w-xs">
        Jalankan query di editor untuk menyimpannya di dalam daftar riwayat log.
      </p>
    </div>

    <!-- Logs List -->
    <div v-else class="flex-1 overflow-y-auto space-y-4 pr-1">
      <div
        v-for="(query, idx) in logs"
        :key="idx"
        class="bg-navy-secondary border border-navy-border hover:border-navy-border-highlight rounded-xl p-4 transition-all duration-150 relative group"
      >
        <!-- Top Metadata Row -->
        <div class="flex items-center justify-between mb-2">
          <div class="flex items-center gap-2">
            <span class="text-xs font-bold text-teal-accent bg-teal-accent/10 px-2 py-0.5 rounded-md">
              #{{ logs.length - idx }}
            </span>
            <span class="text-[10px] text-text-muted">
              Query Terakhir ke-{{ idx + 1 }}
            </span>
          </div>

          <!-- Actions -->
          <div class="flex items-center gap-2">
            <button
              @click="copyQuery(query)"
              class="px-2.5 py-1 text-[11px] bg-navy-tertiary border border-navy-border hover:border-teal-accent text-text-secondary hover:text-teal-accent rounded-md transition-all cursor-pointer flex items-center gap-1"
              title="Salin SQL"
            >
              <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1" />
              </svg>
              Copy
            </button>

            <button
              @click="loadQuery(query)"
              class="px-2.5 py-1 text-[11px] bg-teal-accent text-navy-primary font-medium hover:bg-teal-hover rounded-md transition-all cursor-pointer flex items-center gap-1 shadow-md shadow-teal-accent/5"
              title="Buka di Tab Query Baru"
            >
              <svg class="w-3 h-3" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                <polygon points="5 3 19 12 5 21 5 3" />
              </svg>
              Open in Editor
            </button>

            <button
              @click="deleteLog(idx)"
              class="p-1 rounded bg-navy-tertiary border border-navy-border hover:bg-red-950/20 hover:border-red-900/30 text-text-muted hover:text-accent-red transition-all cursor-pointer"
              title="Hapus entri ini"
            >
              <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M3 6h18M19 6v14c0 1-1 2-2 2H7c-1 0-2-1-2-2V6M8 6V4c0-1 1-2 2-2h4c1 0 2 1 2 2v2" />
              </svg>
            </button>
          </div>
        </div>

        <!-- SQL Code Block -->
        <div class="relative mt-2">
          <pre class="text-xs font-mono text-text-primary bg-navy-tertiary/60 border border-navy-border/80 rounded-lg p-4 overflow-x-auto select-text leading-relaxed whitespace-pre-wrap break-all max-h-64 scrollbar-thin">{{ query }}</pre>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useTabsStore } from '../stores/tabs'
import { useUiStore } from '../stores/ui'

const tabsStore = useTabsStore()
const uiStore = useUiStore()

const logs = ref<string[]>([])

function loadLogs() {
  const stored = localStorage.getItem('yxpg:query_logs')
  if (stored) {
    try {
      logs.value = JSON.parse(stored)
    } catch {
      logs.value = []
    }
  }
}

onMounted(() => {
  loadLogs()
})

function copyQuery(query: string) {
  navigator.clipboard.writeText(query)
  uiStore.addNotification({
    type: 'success',
    title: 'Success',
    message: 'SQL query copied to clipboard',
  })
}

function loadQuery(query: string) {
  tabsStore.createTab('query', {
    title: 'Query',
    sql: query
  })
}

function deleteLog(index: number) {
  logs.value.splice(index, 1)
  localStorage.setItem('yxpg:query_logs', JSON.stringify(logs.value))
  uiStore.addNotification({
    type: 'success',
    title: 'Success',
    message: 'Query entry deleted from log',
  })
}

function clearAllLogs() {
  if (confirm('Apakah Anda yakin ingin menghapus seluruh riwayat query?')) {
    logs.value = []
    localStorage.removeItem('yxpg:query_logs')
    uiStore.addNotification({
      type: 'success',
      title: 'Success',
      message: 'All query history cleared',
    })
  }
}
</script>
