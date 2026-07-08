<template>
  <div class="p-4 space-y-4">
    <div class="flex items-center justify-between">
      <h3 class="text-sm font-semibold text-text-primary">Index Designer</h3>
      <button
        @click="$emit('add')"
        class="flex items-center gap-1 px-2 py-1 text-xs bg-teal-accent text-navy-primary rounded font-medium hover:bg-teal-hover transition-colors"
      >
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 5v14M5 12h14" />
        </svg>
        Add Index
      </button>
    </div>

    <div v-if="indexes.length === 0" class="text-center py-8 text-xs text-text-muted">
      No indexes defined
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="(idx, index) in indexes"
        :key="index"
        class="p-3 bg-navy-tertiary border border-navy-border rounded"
      >
        <div class="flex items-center gap-3 mb-2">
          <input
            v-model="idx.name"
            type="text"
            class="flex-1 px-2 py-1 text-xs bg-navy-primary border border-navy-border rounded"
            placeholder="Index name"
          />
          <select
            v-model="idx.index_type"
            class="w-24 px-2 py-1 text-xs bg-navy-primary border border-navy-border rounded"
          >
            <option value="btree">B-tree</option>
            <option value="hash">Hash</option>
            <option value="gin">GIN</option>
            <option value="gist">GiST</option>
            <option value="brin">BRIN</option>
          </select>
          <label class="flex items-center gap-1 text-xs text-text-secondary">
            <input type="checkbox" v-model="idx.is_unique" class="rounded" />
            Unique
          </label>
        </div>
        <input
          :value="idx.columns.join(', ')"
          @input="idx.columns = ($event.target as HTMLInputElement).value.split(',').map(s => s.trim())"
          type="text"
          class="w-full px-2 py-1 text-xs bg-navy-primary border border-navy-border rounded"
          placeholder="Columns (comma separated)"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { IndexDefinition } from '../../types'

defineProps<{
  indexes: IndexDefinition[]
}>()

defineEmits(['add'])
</script>
