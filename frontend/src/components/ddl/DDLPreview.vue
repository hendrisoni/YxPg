<template>
  <div class="h-full flex flex-col overflow-hidden">
    <div class="px-3 py-2 border-b border-navy-border bg-navy-secondary">
      <span class="text-xs text-text-muted uppercase tracking-wider">DDL Preview</span>
    </div>
    <div class="flex-1 overflow-auto p-3">
      <pre class="text-xs font-mono text-text-secondary whitespace-pre-wrap">{{ ddl }}</pre>
    </div>
    <div class="px-3 py-2 border-t border-navy-border bg-navy-secondary flex items-center gap-2">
      <button
        @click="copyDDL"
        class="px-3 py-1 text-xs border border-navy-border rounded text-text-secondary hover:bg-navy-hover transition-colors"
      >
        Copy DDL
      </button>
      <button
        @click="applyDDL"
        :disabled="!connectionsStore.isConnected"
        class="px-3 py-1 text-xs bg-teal-accent text-navy-primary rounded font-medium hover:bg-teal-hover transition-colors disabled:opacity-50"
      >
        Apply DDL
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useConnectionsStore } from '../../stores/connections'
import { useUiStore } from '../../stores/ui'

defineProps<{
  ddl: string
}>()

const connectionsStore = useConnectionsStore()
const uiStore = useUiStore()

function copyDDL() {
  // Will be passed via props
}

async function applyDDL() {
  // Will be implemented with DDL execution
}
</script>
