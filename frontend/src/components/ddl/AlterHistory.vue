<template>
  <div class="p-4 space-y-4">
    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="w-5 h-5 border-2 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else-if="history.length === 0" class="text-center py-8 text-xs text-text-muted">
      No DDL history
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="entry in history"
        :key="entry.id"
        class="p-2 bg-navy-tertiary border border-navy-border rounded text-xs font-mono text-text-secondary"
      >
        <div class="flex items-center gap-2 mb-1">
          <span class="text-[10px] text-text-muted">{{ formatDate(entry.executed_at) }}</span>
          <span v-if="entry.error" class="text-accent-red">Error</span>
        </div>
        <pre class="whitespace-pre-wrap">{{ entry.sql }}</pre>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const loading = ref(false)
const history = ref<any[]>([])

function formatDate(dateStr: string): string {
  try {
    return new Date(dateStr).toLocaleString()
  } catch {
    return dateStr
  }
}
</script>
