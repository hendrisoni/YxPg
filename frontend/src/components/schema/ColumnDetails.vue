<template>
  <div class="h-full overflow-auto">
    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="w-5 h-5 border-2 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else-if="columns.length === 0" class="text-center py-8 text-xs text-text-muted">
      No columns loaded
    </div>

    <table v-else class="w-full text-xs">
      <thead class="bg-navy-secondary sticky top-0">
        <tr class="border-b border-navy-border">
          <th class="px-3 py-2 text-left text-text-muted font-medium">Column</th>
          <th class="px-3 py-2 text-left text-text-muted font-medium">Type</th>
          <th class="px-3 py-2 text-left text-text-muted font-medium">Nullable</th>
          <th class="px-3 py-2 text-left text-text-muted font-medium">Default</th>
          <th class="px-3 py-2 text-left text-text-muted font-medium">Key</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="col in columns"
          :key="col.column_name"
          class="border-b border-navy-border hover:bg-navy-hover transition-colors"
        >
          <td class="px-3 py-1.5 font-mono text-text-primary">{{ col.column_name }}</td>
          <td class="px-3 py-1.5 text-accent-blue">{{ formatType(col) }}</td>
          <td class="px-3 py-1.5">
            <span v-if="col.is_nullable" class="text-accent-green">YES</span>
            <span v-else class="text-accent-red">NO</span>
          </td>
          <td class="px-3 py-1.5 text-text-muted font-mono">{{ col.default_value || '—' }}</td>
          <td class="px-3 py-1.5">
            <span v-if="col.is_primary_key" class="px-1.5 py-0.5 bg-accent-amber/20 text-accent-amber rounded text-[10px] font-medium">PK</span>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { ColumnInfo } from '../../types'

defineProps<{
  columns: ColumnInfo[]
  loading?: boolean
}>()

function formatType(col: ColumnInfo): string {
  if (col.character_maximum_length) {
    return `${col.data_type}(${col.character_maximum_length})`
  }
  return col.data_type
}
</script>
