<template>
  <div class="p-4 space-y-4">
    <div v-if="loading" class="flex items-center justify-center py-8">
      <div class="w-5 h-5 border-2 border-teal-accent border-t-transparent rounded-full animate-spin"></div>
    </div>

    <div v-else-if="columns.length === 0" class="text-center py-8 text-xs text-text-muted">
      No columns to edit
    </div>

    <div v-else class="space-y-3">
      <div
        v-for="(col, idx) in columns"
        :key="idx"
        class="p-3 bg-navy-tertiary border border-navy-border rounded"
      >
        <div class="flex items-center gap-3 mb-2">
          <input
            v-model="col.name"
            type="text"
            class="flex-1 px-2 py-1 text-xs bg-navy-primary border border-navy-border rounded"
            placeholder="Column name"
          />
          <select
            v-model="col.data_type"
            class="w-40 px-2 py-1 text-xs bg-navy-primary border border-navy-border rounded"
          >
            <option v-for="type in pgTypes" :key="type.value" :value="type.value">{{ type.label }}</option>
          </select>
        </div>

        <div class="flex items-center gap-4 text-xs text-text-secondary">
          <label class="flex items-center gap-1">
            <input type="checkbox" v-model="col.is_nullable" class="rounded" />
            Nullable
          </label>
          <label class="flex items-center gap-1">
            <input type="checkbox" v-model="col.is_primary_key" class="rounded" />
            Primary Key
          </label>
          <label class="flex items-center gap-1">
            <input type="checkbox" v-model="col.is_unique" class="rounded" />
            Unique
          </label>
          <input
            v-model="col.default_value"
            type="text"
            class="flex-1 px-2 py-1 text-xs bg-navy-primary border border-navy-border rounded"
            placeholder="Default value"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { ColumnDefinition } from '../../types'
import { PG_TYPES } from '../../utils/type-mapper'

defineProps<{
  columns: ColumnDefinition[]
  loading?: boolean
}>()

const pgTypes = PG_TYPES
</script>
