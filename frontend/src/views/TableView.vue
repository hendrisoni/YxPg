<template>
  <div class="h-full flex flex-col overflow-hidden">
    <!-- Toolbar -->
    <div class="flex items-center gap-2 px-3 py-1.5 border-b border-navy-border bg-navy-secondary">
      <span class="text-xs text-text-secondary font-medium">
        {{ tab.schema }}.{{ tab.table }}
      </span>

      <div class="flex-1"></div>

      <!-- Export dropdown -->
      <div class="relative" v-if="result">
        <button
          @click="showExportMenu = !showExportMenu"
          class="flex items-center gap-1 px-2 py-1 text-xs text-text-secondary hover:bg-navy-hover rounded transition-colors"
        >
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4M7 10l5 5 5-5M12 15V3" />
          </svg>
          Export
        </button>
        <div
          v-if="showExportMenu"
          class="absolute right-0 top-full mt-1 bg-navy-secondary border border-navy-border rounded-lg shadow-xl py-1 min-w-[120px] z-10"
        >
          <button
            v-for="fmt in ['CSV', 'JSON', 'SQL']"
            :key="fmt"
            @click="handleExport(fmt.toLowerCase())"
            class="w-full px-3 py-1.5 text-xs text-left text-text-secondary hover:bg-navy-hover hover:text-text-primary"
          >
            As {{ fmt }}
          </button>
        </div>
      </div>

      <!-- Stats -->
      <div v-if="result" class="flex items-center gap-2 text-xs text-text-muted">
        <span>{{ result.row_count }} rows</span>
        <span>· {{ result.duration_ms }}ms</span>
      </div>
    </div>

    <!-- Grid -->
    <div class="flex-1 overflow-hidden">
      <ResultGrid :result="result" :tab="tab" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useConnectionsStore } from '../stores/connections'
import { useUiStore } from '../stores/ui'
import type { Tab, QueryResult, BrowseOptions } from '../types'
import ResultGrid from '../components/query/ResultGrid.vue'

const props = defineProps<{
  tab: Tab
}>()

const connectionsStore = useConnectionsStore()
const uiStore = useUiStore()

const result = ref<QueryResult | null>(null)
const loading = ref(false)
const showExportMenu = ref(false)
const currentPage = ref(1)
const pageSize = 100

onMounted(() => {
  loadData()
})

async function loadData() {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!connId || !props.tab.table) return

  loading.value = true
  try {
    const bindings = connectionsStore.getWailsBindings()
    const opts: BrowseOptions = {
      page: 1,
      page_size: 1, // Only fetch columns metadata & total count initially
    }
    result.value = await bindings.BrowseTable(
      connId,
      props.tab.schema || 'public',
      props.tab.table,
      opts as any
    )
  } catch (e: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Load Failed',
      message: e.message,
    })
  } finally {
    loading.value = false
  }
}

async function handleExport(format: string) {
  const connId = props.tab.connectionId || connectionsStore.currentConnectionId
  if (!result.value || !connId) return
  try {
    const bindings = connectionsStore.getWailsBindings()
    const data = await bindings.ExportData(
      result.value as any,
      format,
      props.tab.schema || 'public',
      props.tab.table!
    )

    // Copy to clipboard or download
    await navigator.clipboard.writeText(data)
    uiStore.addNotification({
      type: 'success',
      title: 'Exported',
      message: `Data exported as ${format.toUpperCase()} to clipboard`,
    })
  } catch (e: any) {
    uiStore.addNotification({
      type: 'error',
      title: 'Export Failed',
      message: e.message,
    })
  }
  showExportMenu.value = false
}
</script>
